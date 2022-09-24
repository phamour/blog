---
title: 文献总结 - [2015] Hidden Technical Debt in Machine Learning Systems
date: 2022-09-24 13:59:45
tags: [文献总结, MLOps, Machine Learning]
---

# 一、主旨
使用ML来支撑业务系统并不像看上去那么速赢（quick win），它往往会带来高昂的后期维护开销。本文旨在警示ML从业者（increase the community's awareness），从技术负债的角度，论述了在做ML系统设计时需要注意的容易产生负债的几个方面：系统边界、耦合、反馈循环、消费者、数据依赖、配置问题等等。
![](mlsystems.png)

# 二、内容
{% markmap 520px 3 %}
- ML系统技术负债
  - 边界侵蚀（boundary erosion）
    - 问题：纠缠（entanglement）
      - 牵一发动全身：特征、超参数、抽样方法
      - 解法1：多个独立模型 + ensemble
      - 解法2：可视化工具展示每个特征的贡献度
    - 问题：校正链（correction cascades）
      - 模型魔改适配应用问题
      - 解法1：在原模型上增强、加特征做训练
      - 解法2：重做独立模型
    - 问题：不明消费（undeclared consumers）
      - 模型输出不设限
      - 可见性负债（visibility debt）
      - 隐藏的紧耦合（hidden tight coupling）
      - 解法：严格的访问控制
  - 数据依赖（data dependencies）
    - 问题：不稳定数据依赖（unstable data dependencies）
      - 上游数据易变化
      - 解法：数据来源版本化、冻结（versioned copy / frozen version）
        - 要注意，版本化本身也有额外开销
    - 问题：未充分使用数据依赖（underutilized data dependencies）
      - 无脑灌入无用的特征
      - correlation
      - 解法：经常性的充分的leave-one-feature-out evaluation
    - 论点：对数据依赖的静态分析
  - 反馈循环（feedback loops）
    - 直接反馈（direct feedback loops）
      - 模型的输出可能直接影响到它未来的训练数据
    - 隐藏反馈（hidden feedback loops）
      - 例子1：电商的商品推荐和评论选择两个功能，improve了商品导致点击变化导致评论变化
      - 例子2：两个完全没有交互的基金公司的平台，其中一家improve了产品会影响他们的买入卖出，影响市场，间接影响另一家
  - 反模式（anti-patterns）——短期收益大长期成本高
    - 问题：胶水代码（glue code）
      - 解法：把黑盒封装成通用API
    - 问题：流水线丛林（pipeline jungles）
      - 解法：对数据处理流程和特征工程的整体设计乃至重新设计
    - 问题：僵尸代码
      - if-else扩展出的实验性新方案，短期能够快速试验，长期维护成本高
      - 解法：定期重新检查每个实验分支是否有用
    - 问题：缺乏抽象（abstraction debt）
    - 问题：常见的不好的设计气味（common smells）
      - 裸数据（Plain-Old-Data type smell）
      - 多语言（multiple-language smell）
      - 原型化（prototype smell）
  - 配置负债（configuration debt）
    - 论点：好的配置系统的几项原则
      - 一项配置应是对上一份配置的少量变更 / 增量修改
      - 应足够容忍人工的错误、省略、疏忽
      - 应便于可视对比两份配置的差异
      - 应便于自动化断言和检查
      - 应支持对未使用的或冗余的配置项的检测
      - 配置项应纳入源码仓库和代码评审
  - 外部环境变更（dealing with changes in external world）
    - 问题：动态系统中的静态阈值
      - 解法：学习阈值
    - 挑战：监控与测试
      - 监控点1：预测偏差（prediction bias）——基于模型性能的偏移检测
      - 监控点2：对自动动作设定上限（action limits），触发人工审核
      - 监控点3：上游数据提供商（up-stream producers）
      - 需求：实时性
  - 其他方面
    - 问题：数据测试负债（data testing debt）——数据质量、sanity check
    - 问题：可复现性负债（reproducibility debt）
    - 问题：过程管理负债（process management debt）
      - 背景：多模型
      - 背景：多构建过程
      - 挑战：资源分配
      - 挑战：任务优先级
      - 挑战：故障恢复
    - 问题：文化负债（cultural debt）——研究与工程的融合
      - 崇尚减负如同崇尚模型刷点
        - 删特征
        - 减复杂度
      - 崇尚系统优化如同崇尚模型刷点
        - 提升可复现性
        - 提升稳定性
        - 增加监控
{%endmarkmap%}

# 三、Key takeaways
1. 令人不安的趋势：开发和部署ML系统很快很便宜，但后期维护它们既难且贵
{% blockquote %}
[...] a wide-spread and uncomfortable trend [...] : developing and deploying ML systems is relatively fast and cheap, but maintaining them over time is difficult and expensive.
{% endblockquote %}

2. 并非所有负债都是不好的，但所有负债都要支付利息/都有代价；偿还技术负债的方式有很多，例如重构，其目的不在于增加新功能，而是在于使后续优化成为可能、减少错误、增强可维护性
{% blockquote %}
Not all debt is bad, but all debt needs to be serviced.
{% endblockquote %}
{% blockquote %}
Technical debt may be paid down by refactoring code, improving unit tests, [...]. The goal is not to add new functionality, but to enable future improvements, reduce errors, and improve maintainability.
{% endblockquote %}

3. 在系统设计时考虑技术负债能够提前认识到加速工程落地会为系统维护带来的长期开销
{% blockquote %}
[...] help reason about the long term costs incurred by moving quickly in software engineering.
{% endblockquote %}

4. 相比一般的软件系统，ML系统更易带来技术负债，除了传统的软件系统维护问题之外，还有很多ML专有的问题（ML-specific issues），并且很多负债是系统层面上的，不是代码层面上的，所以传统的通过优化代码来减轻负债的方法是不够的

5. ML模型的黑盒化严重，容易导致大量的含有固定假设/写死的东西的胶水代码或适配校准模块
{% blockquote %}
ML packages may be treated as black boxes, resulting in large masses of "glue code" or calibration layers that can lock in assumptions.
{% endblockquote %}

6. 纠缠问题（entanglement），CACE principle: Changing Anything Changes Everything 牵一发动全身；输入n维特征，其中1维的数据分布变化会导致其他n-1维特征的权重变化；增加特征也是类似的问题；移除特征也是类似的问题；除了输入信号（input signals）以外，超参数、配置、抽样方法、收敛阈值、数据筛选等等都遵循CACE原则。

7. 依赖负债是增加传统软件系统代码复杂度的重要原因，而对于ML系统，数据依赖是依赖当中更难发现的一种

8. 未充分使用的数据依赖（underutilized data dependencies），还是依赖，因为不确定是否可以切断

9. 典型的未充分使用的数据依赖
    1. 遗留特征（legacy features），过去试验性使用过，随着模型迭代被其他特征替代，但未移除
    2. 捆绑特征（bundled features），一组特征放在一起使用，但迫于ddl，未深究任一特征均有用
    3. 微增值特征（epsilon-features），为了极小的准确率提升而引入的性价比低的特征
    4. 相关特征（correlated features）

10. 难以明晰的反馈循环会导致分析负债，即难以保证系统自身变更会从哪些方面影响自己

11. 使用通用的类库往往会导致大量的胶水代码，后者从长期角度看是高维护成本的
{% blockquote %}
Using generic packages often results in a glue code system design pattern, [...]. Glue code is costly in the long term [...]
{% endblockquote %}

12. 胶水代码（glue code）和流水线丛林（pipeline jungles）这类集成问题的根因往往是科研和工程的分离
{% blockquote %}
Glue code and pipeline jungles are symptomatic of integration issues that may have a root cause in overly separated "research" and "engineering" roles.
{% endblockquote %}

13. 通过增加if-else分支可以快速地在一套基准代码上扩展出新方案并投入实验，但长此以往会遗留很多条件分支，难以维护，难以端到端回归测试，进而难以保证上到生产系统后不会有意外情况进入不该进入的分支，Kight Capital血淋淋的教训，45分钟内损失4650万美元。

14. 应推崇研究与工程融合的文化，一碗水端平，模型层面准确率的提升应当与系统层面复杂度的降低有同等的重要性
{% blockquote %}
It is important to create team cultures that reward deletion of features, reduction of complexity, improvements in reproducibility, stability, and monitoring to the same degree that improvements in accuracy are valued.
{% endblockquote %}
{% blockquote %}
Paying down ML-related technical debt [...] often only achieved by a shift in team culture. Recognizing, prioritizing, and rewarding this effort is important for the long term health of successful ML teams.
{% endblockquote %}

15. 在快速发展中的ML系统团队对于减少负债和采用好实践是不自知的或不屑一顾的，但负债和代价是随时间而逐步显现的

16. 在发展过程当中多问自己的问题：
    1. 以生产尺度测试一个全新的算法有多容易？
    2. 数据依赖的完整传导过程是什么？
    3. 一次变更对系统的影响可以被细致观测到什么程度？
    4. 改进模型或输入信号（signal）会使其他模型效果变差吗？
    5. 团队的新成员多久可以习得实战能力？

17. Maintainable ML

18. 为了微小的准确率提升而付出大幅增加系统复杂度的代价，这种研究方案是不可取的（Reasonable ML？）
