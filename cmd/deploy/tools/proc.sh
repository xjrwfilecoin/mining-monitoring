#!/bin/bash
PGM="mining-monitoring"
echo "proc mining-monitoring Server ..."
for i in $PGM
do
        ps -ef |grep -v "grep"|grep "$i"|awk '{print "[proc] ",$0}';
done
