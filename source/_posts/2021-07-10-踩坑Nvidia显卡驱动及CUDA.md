---
title: 踩坑Nvidia显卡驱动及CUDA
date: 2021-07-10 13:46:29
tags: [DailyDev, GPU]
---

这周新机器和显卡逐渐到位，开始装显卡驱动和CUDA，对过程中遇到的一些问题做个总结。

### Boot
显卡总共有4种：A100、RTX2080Ti、K80和K40（包括K40c和K40m两种，这里不做区分），系统是Ubuntu20.04。一上来直接尝试用apt装驱动，所有机器统一装了`nvidia-driver-460-server`，结果几乎所有机器都重启失败，只能重装系统。这里推测是apt里的n卡驱动会把`Xorg`桌面程序装上并带起来，但处理不善导致系统无法启动。重装后手动切了booting target到命令行，避免桌面GUI的启动：（{% link 相关链接 https://www.cyberciti.biz/faq/switch-boot-target-to-text-gui-in-systemd-linux/ %}）
{% codeblock %}
sudo systemctl set-default multi-user.target
{% endcodeblock %}

### A100
**驱动&gt;450，cuda&gt;11**
用apt带的nvidia-driver-460系列安装后有几率无法重启，原因未知。官网cuda安装包，按装440 cuda10.2会报错，报显卡驱动的安装错误，错误码256。cuda11以上的runfile安装完成后，提示可以通过/usr/local/cuda/bin/cuda-uninstaller进行卸载，但实际上这个卸载器并不存在。而10.2的runfile安装包是装了这个卸载程序的。cuda11安装完成后，需要进一步ln几个库才能正常使用：

{% codeblock %}
cd /usr/local/cuda-11.4/targets/x86_64-linux/lib && \
ln -s libcusolver.so.11 libcusolver.so.10 && \
ln -s ../../../extras/CUPTI/lib64/libcupti.so libcupti.so && \
ln -s ../../../extras/CUPTI/lib64/libcupti.so.11.4 libcupti.so.11.0
{% endcodeblock %}

### K80
可以装440+cuda10.2也可以装470+cuda11.4，相对来说兼容性更好，直接通过cuda runfile安装即可，另需按cuda11安装后所需的ln操作。

### K40c/K40m
用apt带的nvidia-driver-460系列安装后有几率无法重启，原因未知。用官网cuda runfile安装包会报错，同样是报显卡驱动的安装错误，错误码256。其实是470驱动不支持k40，k40最高只能到460（截至20210708）。按官方support matrix，460驱动也可以搭配cuda11.4，因此可以官网下载460纯驱动先进行安装，安装前要停掉ubuntu默认的Nouveau驱动：

{% codeblock %}
vi /etc/modprobe.d/blacklist-nouveau.conf

blacklist nouveau
options nouveau modeset=0

update-initramfs -u
{% endcodeblock %}

装完驱动后，再用cuda11.4的runfile装cuda，安装时去掉自带的470驱动安装，仅安装cuda。这样安装完后就是460驱动的k40c/k40m+cuda11.4了。但是这样装完，daemon-reload或重启之后，nvidia-smi会报driver/library version mismatch。可能还是需要装cuda11.2才能完美匹配。
