#修改RootDir为实际存放目录
RootDir=`pwd`
cd $RootDir
cd ..
MonitorProc=( [1]=mining-monitoring )
for i in ${!MonitorProc[@]};
do
	ProcessKey=${MonitorProc[$i]}
	ProcCnt=`ps -ef | grep $ProcessKey |grep -v grep|wc -l`
	if [ $ProcCnt -eq 0 ];then		
		./${MonitorProc[$i]}
	fi	
done
