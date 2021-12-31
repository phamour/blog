---
title: Linux传输工具rsync
date: 2021-12-31 13:13:23
tags: [DailyDev, Linux]
---

# 简介

远程传输工具rsync的命令语法与scp相似，原理与scp不同，最主要的特性是支持增量传输和断点续传（`-P`选项），先抄一段帮助文档的使用说明：

{% codeblock lang:shell %}
Usage: rsync [OPTION]... SRC [SRC]... DEST
  or   rsync [OPTION]... SRC [SRC]... [USER@]HOST:DEST
  or   rsync [OPTION]... SRC [SRC]... [USER@]HOST::DEST
  or   rsync [OPTION]... SRC [SRC]... rsync://[USER@]HOST[:PORT]/DEST
  or   rsync [OPTION]... [USER@]HOST:SRC [DEST]
  or   rsync [OPTION]... [USER@]HOST::SRC [DEST]
  or   rsync [OPTION]... rsync://[USER@]HOST[:PORT]/SRC [DEST]
The ':' usages connect via remote shell, while '::' & 'rsync://' usages connect
to an rsync daemon, and require SRC or DEST to start with a module name.
{% endcodeblock %}

![](rsync-logo.png)

在做目录级的传输时，需要特别注意的是，源目录路径末尾有没有斜杠是会有不同效果的：
- 带斜杠，如/data/，则同步/data目录下的所有文件到目标路径下
- 不带斜杠，如/data，则同步/data目录到目标路径下，作为子目录

# 一些有用的options
* `-a` 全家桶选项，组合了多个其他选项（`-rlptgoD`），基本无脑使用，并且对于目录来说`-a`是包含`-r`的
* `-z` 对于未压缩的文件做压缩传输（gzip）
* `-P` 断点续传
* `-e` 通过ssh连接时，如果需要走key，可以写成`-e "ssh -i /path/to/sshkey"`
* `--bwlimit` 限制网络传输速度，值的单位为KB/s

# 例子
## Push to remote
{% codeblock lang:shell %}
rsync -azP /data/ user@1.2.3.4:/data
{% endcodeblock %}
将本地/data目录下的所有文件推送到远程1.2.3.4机器的/data目录下。
## Pull from remote
{% codeblock lang:shell %}
rsync -e "ssh -i ~/.ssh/mykey" --bwlimit=3000 -azP user@1.2.3.4:/data/mydata.tar /data/
{% endcodeblock %}
以`~/.ssh/mykey`作为ssh key，限制速度约为3MB/s，将远程1.2.3.4机器的/data/mydata.tar文件拉取到本地/data目录下。
