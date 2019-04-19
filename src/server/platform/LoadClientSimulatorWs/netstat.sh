#!/bin/sh
for((i=0;i<10000;i++))
do
{
	netstat -aon |grep "192.168.2.28" |grep ESTABLISHED |wc -l
	sleep 5;
}
done
