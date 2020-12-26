#!/usr/bin/env bash
echo "start bak ..."
rm -rf bak
mkdir bak
cp poolweb timertask bak

echo "start compile ..."
cd ~/go/src/poolweb
git pull origin master
go build
\cp poolweb /usr/local/pooltest
cd timertask
go build
\cp timertask /usr/local/pooltest
echo "stop application ..."
PIDS=`ps -ef | grep poolweb | awk '{print $2}'`
for pid in $PIDS
do      
  kill -9 $pid
done
PIDS=`ps -ef | grep timertask | awk '{print $2}'`
for pid in $PIDS
do      
  kill -9 $pid
done
echo "start application ..."
cd /usr/local/pooltest
nohup ./poolweb >> poolweb.log 2>&1 &
sleep 1
nohup ./timertask >> timertask.log 2>&1 &
echo "done"
