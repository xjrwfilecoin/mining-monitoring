RootDir=`pwd`
cd $RootDir
cd ..
MonitorProc=( [1]=mining-monitoring )
for i in ${!MonitorProc[@]};
do
	ProcessKey=${MonitorProc[$i]}
	Suffix=".pid"
	Prefix="."
	ProcPid=`cat $Prefix$ProcessKey$Suffix`
	ProcCnt=`ps -ef | grep $ProcPid |grep -v grep|wc -l`
	if [ $ProcCnt -gt 0 ];then		
		echo "stop ${MonitorProc[$i]} $ProcPid"
		kill -USR1 $ProcPid
	fi	
done
