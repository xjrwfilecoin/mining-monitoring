#!/usr/bin/env bash
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
