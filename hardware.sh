#！/usr/bin bash
ssh -p 22 root@192.168.0.10 'sensors&&uptime&&free -h&&df -h&&sar -n DEV 1 2&& iotop -bn1|head -n 2&&nvidia-smi'
