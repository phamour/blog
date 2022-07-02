---
title: 文献总结 - [2016] Burns_Borg, Omega, and Kubernetes
date: 2022-04-16 13:29:25
tags: [文献总结, Borg, 分布式系统]
---

原文链接：[Borg, Omega, and Kubernetes: Lessons learned from three container-management systems over a decade](https://dl.acm.org/doi/abs/10.1145/2898442.2898444)

![](network.jpg)

# 一、主旨
Google内研或主导的三个集群管理/编排系统，三者各有不同的起始目标，因此其架构和特点也各有不同。但K8s作为Borg和Omega的形式后继，有着后发优势，发扬了前两者的一些优势，也吸取了前两者的一些教训，此文可从一定程度上深化我们对K8s一些decisions的理解。

{% blockquote %}
Lessons learned from three container-management systems over a decade
{% endblockquote %}

# 二、内容梳理
{% markmap 640px 2 %}
- 容器管理
  - Borg
    - 前身：Babysitter & Global Work Queue
    - 目标：资源利用率（20年前）
    - 清晰定义在线任务：latency-sensitive user-facing services
    - 清晰定义离线任务：CPU-hungry batch processes
    - 评价1：新增需求/功能=增加复杂度
    - 评价2：功能堆叠，成为ad-hoc collection of systems
    - 评价3：功能多而广，extreme robustness
  - Omega
    - Borg的后代 (an offspring of Borg)
    - 目标：improve the software engineering of the Borg ecosystem
    - 架构特性1：中心化的状态存储Paxos（类比ETCD之于K8s）
    - 架构特性2：乐观锁机制 (optimisitc concurrency control)
    - 架构特性3：去中心化的系统组件（即，无master或apiserver）
    - 部分特性被Borg吸收
  - K8s
    - 起因：Google外部的用户对容器管理的需求——社区
    - 开源
    - 目标1：简化复杂分布式系统的部署与管理 (make it easy to deploy and manage complex distributed systems)
    - 目标2：资源利用率 (benefiting from the improved utilization that containers enable)
    - 前人栽树1：中心化状态存储——ETCD
    - 前人栽树2：中心化总控逻辑封装——apiserver
- 容器技术（略）
- 基础设施变革
  - 容器技术很大程度上影响了数据中心/云计算的基础架构形态，更方便业务上云
  - Before: machine-oriented infra
  - After: application-oriented infra
  - 容器带来的对执行环境的封装
    - 屏蔽部署环境的细节
    - 抹平开发到部署的差距：部署可靠性 & 开发部署一致性 & 推动DevOps
    - Borg背后的容器是自己维护的一套东西，以机器为单位做base image，有很多问题
  - 容器带来的应用的隔离
    - 管理容器=管理应用
    - 为容器内部安插generic的API/probe: metrics, health, etc.
    - 在容器外部以容器为单位做monitoring (cAdvisor)
    - 多容器协同
- Consistency
  - 为了避免或减小系统发展带来的复杂度增加
  - K8s的API设计充分考虑了consistency
  - 设计特性1：API的模板化——apiVersion, kind, metadata, spec
  - 设计特性2：API分类分区解耦
    - 接口+实现，一类对象满足不同需求的不同实现
      - ReplicaSet/Deployment
      - DaemonSet
      - StatefulSet
      - Job
  - 设计特性3：调谐为设计模式的控制逻辑
- 教训
  - 节点共享的IP-ports vs. Pod IP-ports
  - 基于列表的容器组织 vs. 基于标签的容器组织
  - 强制的从属关系 vs. 松散的从属关系
  - 全中心化的总控API vs. 中心化存储和去中心化的控制 vs. 中心化存储和中心化控制封装API和去中心化的业务
- 挑战
  - 配置管理
  - 依赖管理
    - 工程实践中的重要问题就是依赖
    - 部署时，依赖的处理并非想象那么简单
{%endmarkmap%}

# 三、重点分析
1. 架构的扩展性和可持续性/一致性 (consistency) 对于高速发展中的系统非常重要，业务扩张必然带来需求增加，如果不能做到新功能实现与原有系统的一致/连贯，必然会对系统带来额外的复杂度，不论是开发、维护还是部署、使用。
2. 容器技术为系统架构、云计算、数据中心、PaaS均带来突破性变革（普遍共识）。
3. K8s的一致的控制循环——调谐：从顶层业务系统视角出发，拆分成不同的系统组件，基于容器，包装成不同的对象，每个对象有不同的控制器，但控制逻辑均为调谐，通过一个个小的控制循环达成整体业务系统的控制。
4. 原文举例说明了一个非常实用的在生产环境在线排查问题的方法：（带入Anylearn视角）后端出问题，卡死或响应异常，定位到某个pod后，可以直接扔掉该pod的标签（该pod还在正常存续中），让上层的ReplicaSet或其他replication controller的selector过滤不到这个pod，便会再启动一个pod（因为它认为丢失了一个）。这样一来，既可以“重启”后端的副本，恢复业务，同时可以留住pod，保存现场，乃至上到pod内排查问题，两不耽误。
5. 建模其实也非常重要。Google对于Borg最不满意的一点就是其基于index的容器组织逻辑，容器数量成规模后，中间丢一个index就可能会有问题，维护所有的index-container mapping也是个麻烦事。基于label的组织和过滤就要灵活得多。带入Anylearn视角，我们现在的数据建模可能需要三思的有：训练项目与训练任务的从属关系模型、资源与算法/数据集/模型/文件的继承关系模型、模型未拆分成细粒度对象、文件未拆分成细粒度对象，等等。

# 四、引申
做分布式锁（本文中提到选主相关应用）的Chubby： Burrows, M. 2006. The Chubby lock service for loosely coupled distributed systems. Symposium on Operating System Design and Implementation (OSDI), Seattle, WA.
