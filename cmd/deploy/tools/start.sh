RootDir=`pwd`
cd $RootDir
cd ..
PGM="mining-monitoring"
echo "Start mining-monitoring Server ..."

for i in $PGM
do
    ./$i
	ps -aux |grep -v "grep"|grep "$i"|awk '{print "[start] ",$0}';
done
