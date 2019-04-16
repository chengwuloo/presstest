package main

//
// Created by YangZhi
// 			4/9/2019
//

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"server/platform/util"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////
//.\ClientSimulatorWs2.exe -httpaddr= -wsaddr= -mailboxs= -clients= -baseTest= -deltaClients= -deltaTime= -interval= -timeout=

//HTTPAddr HTTP请求token地址
var httpaddr = flag.String("httpaddr", "192.168.2.214", "")

//wsaddr Websocket登陆地址
var wsaddr = flag.String("wsaddr", "192.168.2.211:10000", "")

//numMailbox 单进程邮槽数，最好等于clients 5000
var numMailbox = flag.Int("mailboxs", 1, "")

//numClient 单进程客户端数 5000
var numClient = flag.Int("clients", 1, "")

//BaseAccount 测试起始账号
var baseAccount = flag.Int64("baseTest", 6000000, "")

//deltaClients 间隔连接数检查时间戳
var deltaClients = flag.Int("deltaClients", 500, "")

//deltaTime 间隔毫秒数检查连接数
var deltaTime = flag.Int("deltaTime", 8000, "")

//heartbeat 心跳间隔毫秒数
var heartbeat = flag.Int("interval", 2000, "")

//timeout 心跳超时清理毫秒数 timeout>interval
var timeout = flag.Int("timeout", 5000, "")

//subGameID 测试子游戏，游戏类型
var subGameID = flag.Int("gameID", 210, "")

//subRoomID 测试子游戏，房间号
var subRoomID = flag.Int("roomID", 2101, "")

//tokenprefix 测试token，免http登陆
var tokenprefix = flag.String("prefix", "test_new2_", "")
var tokenstart = flag.Int("tokenstart", 6666, "")
var tokenend = flag.Int("tokenend", 99999, "")

//timestart 起始时间戳
var timestart, timestart2 Timestamp

//curConn 当前连接数
var curConn int64

//timenow 当前时间戳
var timenow Timestamp

//onInput 输入命令行参数 'q'退出 'c'清屏
func onInput(str string) int {
	switch str {
	case "q":
		{
			gSessMgr.CloseAll()
			gMailbox.Stop()
			return -1
		}
	case "c":
		{
			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()
			return 0
		}
	case "--help":
		{
			return 0
		}
	}
	return 0
}

//
func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(debug.Stack())
		}
	}()
	log.Printf("--- *** PID[%07d] %v\n", os.Getpid(), os.Args)
	//解析命令行
	flag.Parse()
	//worker工厂
	smain := NewSentryCreator()
	//启动10000个邮槽
	t1 := TimeNowMilliSec()
	gMailbox.Start(smain, *numMailbox, (*timeout)/1000+10)
	t2 := TimeNowMilliSec()
	log.Printf("--- *** PID[%07d] gMailbox.Start = [%03d] elapsed = %dms\n", os.Getpid(), *numMailbox, TimeDiff(t2, t1))
	go func() {
		//起始时间戳
		timestart = TimeNowMilliSec()
		timestart2 = timestart
		for i := 0; i < *numClient; i++ {
			//HTTP请求token
			// token, err := HTTPGetToken(*httpaddr, *baseAccount+int64(i))
			// if token == "" || err != nil {
			// 	continue
			// }
			//当前时间戳
			//timenow = TimeNowMilliSec()
			// timdiff := TimeDiff(timenow, timestart)
			// if timdiff >= int32(*deltaTime) {
			// 	timestart = timenow
			// 	c := gSessMgr.Count()
			// 	delteConn := c - curConn
			// 	curConn = c
			// 	log.Printf("--- *** detla = %dms deltaClients = %03d", timdiff, delteConn)
			// }
			//websocket客户端
			client := NewDefWSClient()
			token := *tokenprefix + fmt.Sprintf("%d", *tokenstart+i)
			client.(*DefWSClient).Token = token
			client.(*DefWSClient).Account = *baseAccount + int64(i)
			//连接游戏大厅
			client.ConnectTCP(*wsaddr)
		}
		t3 := TimeNowMilliSec()
		log.Printf("--- *** PID[%07d] clients = Succ[%03d] elapsed = %dms\n", os.Getpid(), gSessMgr.Count(), TimeDiff(t3, timestart))
	}()

	//控制台命令行输入 按'q'退出 'c'清屏
	util.ReadConsole(onInput)
	gSessMgr.Wait()
	gMailbox.Wait()
	log.Printf("--- *** PID[%07d] exit...", os.Getpid())
}
