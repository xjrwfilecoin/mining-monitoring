#！/usr/bin bash
 sensors&&uptime&&free -h&&df -h&&sar -n DEV 1 2&& iotop -bn1|head -n 2&&nvidia-smi
