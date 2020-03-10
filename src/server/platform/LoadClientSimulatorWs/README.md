#压测工具使用
LoadClientSimulatorWs 负责加载 ClientSimulatorWs2 子进程
配置在LoadClientSimulatorWs/conf.ini下

#压测需要2台机器：
1. 一台跑ProxyServer HallServer ，ApiServer （最好cpu 32核心以上）
2. 一台跑压测工具(最好cpu 32核心以上) ，比较吃cpu，本机压测工具并发cpu占用率200%以上
3. 足够的网络带宽

#带宽利用率
cat /proc/net/dev

sudo yum install -y lrzsz unzip
sudo ulimit -a / ulimit -n

sudo vim /etc/security/limits.d/20-nproc.conf
*          soft    nproc     204800
*          hard    nproc     204800
*          soft    nofile    65535
*          hard    nofile    65535
sudo vim /etc/security/limits.conf
*          soft    nproc     204800
*          hard    nproc     204800
*          soft    nofile    65535
*          hard    nofile    65535
//http://juduir.com/topic/1802500000000000075

//提高单机短连接QPS到20万
//https://blog.csdn.net/chdhust/article/details/53047602

sudo sysctl -w net.ipv4.tcp_fin_timeout=15
sudo sysctl -w net.ipv4.tcp_timestamps=1
sudo sysctl -w net.ipv4.tcp_tw_recycle=1

sudo sysctl -a |grep port_range
//net.ipv4.ip_local_port_range = 32768	60999
sudo vim /etc/sysctl.conf
//net.ipv4.ip_local_port_range = 10000	65000
sudo sysctl -p
http://blog.csdn.net/wenshuangzhu/article/details/44060901

//Too Many open files问题
//https://www.cnblogs.com/htkj/p/10932537.html

lsof -n|awk '{print $2}'|uniq -c|sort -nr
netstat -ant|grep 8080|wc -l

//tcp连接内核参数调优somaxconn
//https://blog.csdn.net/jackyechina/article/details/70992308
https://stackoverflow.com/questions/2569620/socket-accept-error-24-to-many-open-files
vim /etc/sysctl.conf
net.ipv4.tcp_max_syn_backlog = 8192
net.ipv4.tcp_synack_retries = 2
net.core.somaxconn=65535
net.core.netdev_max_backlog = 65535
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_tw_recycle = 1
net.ipv4.tcp_fin_timeout = 15
net.ipv4.tcp_timestamps = 1
fs.file-max=65535
sysctl -p

sysctl -w net.core.somaxconn=65535
sysctl -a

//https://www.cnblogs.com/seasonxin/p/8192020.html

