---
title: 文献总结 - [2015] Large-scale cluster management at Google with Borg
date: 2022-04-09 14:43:40
tags: [文献总结, Borg, 分布式系统]
---

原文链接：[Large-scale cluster management at Google with Borg](https://dl.acm.org/doi/abs/10.1145/2741948.2741964)

![](cluster.jpg)

# 一、主旨
Borg是Google内部从21世纪初就开始研发的大规模集群管理系统，属于底层基础架构，服务上层业务应用的部署和运行。
Borg对K8s有深远意义，从某种意义上可以说Borg是K8s的前身，但两者并不一样，且时至今日在Google内部Borg也并没有被K8s替代，应当也没有被替代的趋势。

# 二、内容
{% markmap 640px %}
- 业务语义
  - 负载类型
    - 服务/长活/在线任务
      - 对瞬时性能波动敏感
    - 批处理/离线/短时任务
      - 相对更容忍短时间内的性能波动
  - 集群拓扑
    - 物理集群cluster：单一数据中心中的单一物理厂房
      - 通常单一逻辑集群单元cell
      - 特殊性质的其他逻辑集群单元
      - ~10000服务器
      - 异构
    - 天生的多集群管理
  - 任务单元
    - job
      - 运行在同一个cell内
      - tasks
      - 1task = 1容器
      - job-tasks逻辑关系不清晰
  - 节点执行单元
    - alloc (类比pod)
      - 1 or n tasks
  - 集群执行单元
    - alloc set (类比三方batch job)
  - 调度配置
    - priority
      - 高优杀低优
      - 预置优先级区间
        - monitoring
        - production
          - 此区间不会发生争抢 (preemption cascade)
        - batch
        - best effort
    - quota
      - 配额检查发生在submission阶段，并非实际调度阶段
      - 超售
    - admission control
  - 服务间识别
    - BNS (内部DNS，与K8s机制类似)
  - 监控
    - WebUI——Sigma
      - why pending提示
      - 任务结束后的log短时保留
    - event记录——Infrastore
- 系统架构
  - Borgmaster (类比apiserver)
    - 单进程模型 5副本提高可用性
    - 选主
    - WebUI (Sigma的备选)
    - checkpoints/snapshots
    - 模拟器：Fauxmaster (using checkpoints)
    - 一个master副本只负责集群中的一部分节点 (见后Borglet)
  - scheduling
    - 任务队列
      - 被调度器监听
    - 调度器
      - round-robin扫队列，保证队头不被难以调度的大任务卡死
      - 可行性检查
        - 满足任务要求的节点
      - 评分
        - 用户定义的preference
        - 内置调度策略的体现
          - 最小化被高优抢占的任务数量
          - 已拥有任务所需packges (程序、数据、文件) 的节点优先
          - 分散可用区、分散节点失效风险
          - 混合高优在线任务和低优离线任务，保证高优任务的突发资源抢占
        - worst fit
          - E-PVM 尽可能分散调度任务以保持节点评分的稳定性 (minimizes the change in cost)
          - 优点：资源占用均衡，对在线任务的突发资源抢占有好处
          - 缺点：资源使用过于分散，fragmentation严重
        - best fit
          - 尽可能优先填满一台机器 (fill machines as tightly as possible)
          - 缺点：一些低优任务希望申请极少资源并占用机器上未分配的资源，投机取巧地使用更多资源，在best fit下这种任务很容易被调度但很难使用到更多的资源，导致任务其实根本跑不动
          - 缺点：突发性能必然面临强杀一些低优离线任务
          - 优点：通常会有不少机器属于空置状态，对大规模离线任务有好处
        - hybrid
          - 实际在Borg中使用的策略
          - 尽量减少某一种资源的“成块碎片”
            - 并非我们之前讲的资源碎片
            - 一个节点上其中一类资源有空闲但其他资源都占满导致这类资源浪费
          - 比best fit好3-5%
        - 因抢占资源而遭驱逐的任务将被放回任务队列重新调度
  - Borglet
    - 被Borgmaster轮询
    - 向Borgmaster汇报状态
    - 一个Borgmaster只轮询集群内的一部分节点的Borglet
    - 节点失效时Borgmaster会重调度该节点的所有任务，节点恢复后将收到Borgmaster指令驱逐掉这部分重调度的任务
  - 可伸缩性
    - 文字肌肉：Borg无已知上限，每次将要遇到瓶颈均被解决
    - 数字肌肉：1master可轻松管理几千节点，有的集群可以达到10000 tasks per minute
    - 乐观锁机制
    - 读写请求分离
    - 调度评分缓存
    - 任务的资源套餐 (同样资源需求的任务只计算一次调度评分)
      - 这里不确定是静态的套餐还是动态计算出的套餐
    - 调度评分时，节点顺序随机，且无需遍历所有节点
- 技术挑战
  - 可用性
    - 失效是大规模系统的日常 failures are the norm in large scale systems
    - Borgmaster或Borglet失效时，节点上的任务仍正常运行
    - 重调度被驱逐的任务
    - 分散一个job的不同tasks
    - 限定一个job中最多能down多少tasks (类比maxUnavailable)
    - 声明式API和幂等写操作
    - 节点失效等待，防抖
    - 99.99%
  - 资源利用率
    - 资源浪费=成本
    - 怎么测？减少机器 (cell compaction) 直至不能消化原有负载，以机器数量下限为指标
    - 混合在离线任务
      - 可变相理解为资源池化，即，不强行区分在离线任务使用的资源
      - 池化资源确实会带来近邻效应影响性能，但相比资源节约带来的收益，是可以接受的
    - 集群大小不能随便缩
    - 细粒度的资源请求
      - 限定资源套餐会带来资源浪费
    - 在线资源回收 (resource reclamation)
      - 请求的资源超出实际使用的资源
        - 需要细粒度的资源使用情况感知
        - 需要对在线服务的资源使用做预测
      - 超出部分拿出来分给可以接受不稳定资源的任务
      - 在线服务不使用回收资源
  - 资源隔离
    - 安全需求
    - 性能需求
{%endmarkmap%}

# 三、对K8s的影响
1. 应用级别的多构件关联
    1. Borg中只有job概念可以串联多构件
    2. K8s中可通过label来松散串联多构件
2. 网络地址空间
    1. Borg中所有任务共用节点IP和端口，容器端口直接映射节点端口，空间狭小
    2. K8s中pod自带网络地址空间，上层还有svc可以封装和控制是否映射到节点端口
3. 用户倾向
    1. Borg最初目的是优先服务大团队
    2. K8s社区决定玩家结构
4. 多容器的调度单元——pod
    1. Alloc in Borg: app + logsaver + dataloader 模式，分段开发维护
    2. Pod in K8s: init containers + helper containers/sidecars
5. 运维的可观测辅助
    1. Debugging information
    2. Events
6. Master即内核——Borgmaster in Borg = kube-apiserver in K8s

# 四、同类系统
| SysName                                                     | Ref                                                                                                                                                                                                                                                                                                                   |
|-------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Apache Mesos                                                | B. Hindman, A. Konwinski, M. Zaharia, A. Ghodsi,A. Joseph, R. Katz, S. Shenker, and I. Stoica.Mesos: aplatform for fine-grained resource sharing in the data center.InProc. USENIX Symp. on Networked Systems Design andImplementation (NSDI), 2011                                                                   |
| YARN                                                        | V. K. Vavilapalli, A. C. Murthy, C. Douglas, S. Agarwal,M. Konar, R. Evans, T. Graves, J. Lowe, H. Shah, S. Seth,B. Saha, C. Curino, O. O’Malley, S. Radia, B. Reed, andE. Baldeschwieler.Apache Hadoop YARN: Yet AnotherResource Negotiator.InProc. ACM Symp. on CloudComputing (SoCC), Santa Clara, CA, USA, 2013.  |
| Tupperware                                                  | A. Narayanan.Tupperware: containerized deployment atFacebook.http://www.slideshare.net/dotCloud/tupperware-containerized-deployment-at-facebook,June 2014.                                                                                                                                                            |
| Apache Aurora (retired)                                     | Apache Aurora.http://aurora.incubator.apache.org/, 2014.                                                                                                                                                                                                                                                              |
| Autopilot                                                   | https://aurora.apache.org/                                                                                                                                                                                                                                                                                            |
| Quincy (on Borg)                                            | M. Isard, V. Prabhakaran, J. Currey, U. Wieder, K. Talwar,and A. Goldberg.Quincy: fair scheduling for distributedcomputing clusters.InProc. ACM Symp. on OperatingSystems Principles (SOSP), 2009.                                                                                                                    |
| Cosmos                                                      | P. Helland.Cosmos: big data and big challenges.http://research.microsoft.com/en-us/events/fs2011/helland\_cosmos\_big\_data\_and\_big\_challenges.pdf, 2011.                                                                                                                                                          |
| Apollo                                                      | E. Boutin, J. Ekanayake, W. Lin, B. Shi, J. Zhou, Z. Qian,M. Wu, and L. Zhou.Apollo: scalable and coordinatedscheduling for cloud-scale computing.InProc. USENIXSymp. on Operating Systems Design and Implementation(OSDI), Oct. 2014.                                                                                |
| Fuxi                                                        | Z. Zhang, C. Li, Y. Tao, R. Yang, H. Tang, and J. Xu.Fuxi: afault-tolerant resource management and job schedulingsystem at internet scale.InProc. Int’l Conf. on Very LargeData Bases (VLDB), pages 1393–1404. VLDB EndowmentInc., Sept. 2014.                                                                        |
| Omega                                                       | M. Schwarzkopf, A. Konwinski, M. Abd-El-Malek, andJ. Wilkes.Omega: flexible, scalable schedulers for largecompute clusters.InProc. European Conf. on ComputerSystems (EuroSys), Prague, Czech Republic, 2013.                                                                                                         |
| Kubernetes                                                  | https://kubernetes.io/

# 五、文献引申
1. 大规模系统的关键因素：J. Hamilton.On designing and deploying internet-scaleservices.InProc. Large Installation System AdministrationConf. (LISA), pages 231–242, Dallas, TX, USA, Nov. 2007.
2. 系统性能评测的建模：D. G. Feitelson.Workload Modeling for Computer SystemsPerformance Evaluation.Cambridge University Press, 2014.
3. 资源利用率相关指标和实验：A. Verma, M. Korupolu, and J. Wilkes.Evaluating jobpacking in warehouse-scale computing.InIEEE Cluster,pages 48–56, Madrid, Spain, Sept. 2014.
4. 资源调度worst-fit：Y. Amir, B. Awerbuch, A. Barak, R. S. Borgstrom, andA. Keren.An opportunity cost approach for job assignmentin a scalable computing cluster.IEEE Trans. Parallel Distrib.Syst., 11(7):760–768, July 2000.
5. 真实工作负载记录（数据集）：J. Wilkes.More Google cluster data.http://googleresearch.blogspot.com/2011/11/more-google-cluster-data.html, Nov. 2011.
6. 上述数据集的使用：
    1. 数据集分析：C. Reiss, A. Tumanov, G. Ganger, R. Katz, and M. Kozuch.Heterogeneity and dynamicity of clouds at scale: Googletrace analysis.InProc. ACM Symp. on Cloud Computing(SoCC), San Jose, CA, USA, Oct. 2012.
    2. 使用：O. A. Abdul-Rahman and K. Aida.Towards understandingthe usage behavior of Google cloud users: the mice andelephants phenomenon.InProc. IEEE Int’l Conf. on CloudComputing Technology and Science (CloudCom), pages272–277, Singapore, Dec. 2014.
    3. 使用：S. Di, D. Kondo, and W. Cirne.Characterization andcomparison of cloud versus Grid workloads.InInternationalConference on Cluster Computing (IEEE CLUSTER), pages230–238, Beijing, China, Sept. 2012.
    4. 使用：S. Di, D. Kondo, and C. Franck.Characterizing cloudapplications on a Google data center.InProc. Int’l Conf. onParallel Processing (ICPP), Lyon, France, Oct. 2013.
    5. 使用：Z. Liu and S. Cho.Characterizing machines and workloadson a Google cluster.InProc. Int’l Workshop on Schedulingand Resource Management for Parallel and DistributedSystems (SRMPDS), Pittsburgh, PA, USA, Sept. 2012.
7. Borg后续
    1. [2016_Burns_Borg, Omega, and Kubernetes: Lessons learned from three container-management systems over a decade](https://dl.acm.org/doi/abs/10.1145/2898442.2898444)
    2. [2020_Tirmazi_Borg: the next generation](https://dl.acm.org/doi/abs/10.1145/3342195.3387517)
