#!/usr/bin/env bash
PIDS=`ps -ef | grep deploy | awk '{print $2}'`
for pid in $PIDS
do
  kill -9 $pid
done

nohup ./deploy >> deploy.log 2>&1 &
echo "done"
