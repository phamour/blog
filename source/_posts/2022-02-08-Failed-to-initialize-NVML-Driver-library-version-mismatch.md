---
title: 'Failed to initialize NVML: Driver/library version mismatch'
date: 2022-02-08 23:55:38
tags: [DailyDev, Linux, GPU]
---

![](error.png)

显卡驱动与NVML库版本不一致问题。

可到`/usr/lib/x86_64-linux-gnu/`路径下查找`libnvidia-ml*`相关软链和文件：

{% codeblock lang:shell %}
root@10-101-72-17:/usr/lib/x86_64-linux-gnu# ll | grep libnvidia-ml
lrwxrwxrwx   1 root root        17 Jan 29 01:08 libnvidia-ml.so -> libnvidia-ml.so.1
lrwxrwxrwx   1 root root        26 Jan 29 01:08 libnvidia-ml.so.1 -> libnvidia-ml.so.470.103.01
-rw-r--r--   1 root root   1828056 Jan  6 20:11 libnvidia-ml.so.470.103.01
-rwxr-xr-x   1 root root   1823960 Nov 15 15:58 libnvidia-ml.so.470.42.01*
{% endcodeblock %}

一般`libnvidia-ml.so`为软链指向`libnvidia-ml.so.1`又为软链指向一个`libnvidia-ml.so.xxx.yyy.zzz`的文件，`xxx.yyy.zzz`为NVML版本号。

如果该目录下存在多个不同版本号的NVML动态链接库文件（如上述代码段中所示），则有可能手动安装过显卡驱动并随驱动安装了一个版本的NVML，而系统又因各种原因新安装了另一个版本，导致驱动版本与NVML版本不匹配。

这里一个可能的原因是Ubuntu的系统自动更新，可通过查阅`/var/log/apt/history.log`日志，看是否自动更新过libnvidia或libcuda等可能与显卡相关的组件。

简单粗暴的解决方案：重启节点并重新安装驱动，手动安装或apt均可，保证将驱动和NVML刷成一致的版本即可。

另一个不想重启的解决思路（未验证）：如果是手动安装的驱动，由系统更新造成了版本问题，不重启时直接重装原来装过的驱动时会失败，或许可以下载与系统更新出的NVML版本完全一致的驱动文件再手动安装。
