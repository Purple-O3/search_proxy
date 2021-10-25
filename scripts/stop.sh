pid=$(ps -ef |grep search_proxy | grep -v grep |awk '{print $2}')
kill $pid
