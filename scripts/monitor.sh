#pid=$(ps -ef |grep search_proxy | grep -v grep |awk '{print $2}')
pid=$(ps -ef |grep search_ | grep -v grep |awk '{print $2}')
echo $pid
top -pid $pid
