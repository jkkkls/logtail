#!/bin/bash

SERVERNAME="logtail"

start()  #运行程序
{
    echo "start $SERVERNAME"
    $SERVERNAME & #绝对路径拉起程序
    # 上面的 & 一定要保留 他的作用类似脱机执行，不然程序会阻塞
    echo "start $SERVERNAME ok!"
    exit 0;
}

stop() 	#停止程序
{
    echo "stop $SERVERNAME"
    killall $SERVERNAME
    echo "stop $SERVERNAME ok!"
}

case "$1" in
start)
    start
    ;;
stop)
    stop
    ;;
restart)
    stop
    start
    ;;
*)
    echo "usage: $0 start|stop|restart"
    exit 0;
esac
exit

