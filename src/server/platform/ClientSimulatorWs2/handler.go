package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"server/pb/GameServer"
	"server/pb/Game_Common"
	"server/pb/HallServer"
	"server/platform/util"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

//ParallGetTokenRequest 发起并发获取Token请求
//-------------------------------------------------------------
func ParallGetTokenRequest() {
	//*httpaddr = "192.168.2.214:8083"
	go func() {
		//起始时间戳
		timestart = TimeNowMilliSec()
		for i := 0; i < *totalClients; i++ {
			//获取token请求
			sendHTTPGetTokenRequest(*httpaddr, *baseAccount+int64(i), *agentID)
		}
	}()
}

//ParallLoginRequest 发起并发连接/登陆请求
//-------------------------------------------------------------
func ParallLoginRequest() {
	go func() {
		//起始时间戳
		timestart = TimeNowMilliSec()
		for i := 0; i < *totalClients; i++ {
			//进入访问资源
			gSemLogin.Enter()
			//HTTP请求token
			token, ipaddr, err := HTTPGetToken(*httpaddr, *baseAccount+int64(i), *agentID)
			if token == "" || ipaddr == "" || err != nil {
				continue
			}
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
			//token := *tokenprefix + fmt.Sprintf("%d", *tokenstart+i)
			client.(*DefWSClient).Token = token
			client.(*DefWSClient).Account = *baseAccount + int64(i)
			//连接游戏大厅
			if *dynamic == 0 {
				client.ConnectTCP(*wsaddr)
			} else {
				client.ConnectTCP(ipaddr)
			}
		}
	}()
}

//ParallEnterRoomRequest 发起并发进房间请求
//-------------------------------------------------------------
func ParallEnterRoomRequest() {
	go func() {
		//起始时间戳
		timestart = TimeNowMilliSec()
		//游戏类型和房间都有效，则进入房间
		p, ok := GGames.Exist(int32(*subGameID))
		if 0 != *subGameID && 0 != *subRoomID && ok && p.Exist(int32(*subRoomID)) {
			for i := 0; i < *totalClients; i++ {
				//进入访问资源
				gSemEnter.Enter()
				sesID := PopPeer()
				if sesID > 0 {
					peer := gSessMgr.Get(sesID)
					if peer != nil {
						//登陆成功，进入房间
						// client := peer.GetCtx(TagUserInfo).(*DefWSClient)
						// client.GameID = int32(*subGameID) //保存
						// client.RoomID = int32(*subRoomID)
						// reqEnterRoom(peer, client.GameID, client.RoomID, client.Pwd[:])
						//登陆成功，获取游戏列表
						reqGameListInfo(peer)
					}
				}
			}
		}
	}()
}

//ParallOrderRequest 发起并发上下分请求
//-------------------------------------------------------------
func ParallOrderRequest() {
	go func() {
		//起始时间戳
		timestart = TimeNowMilliSec()
		//为10000个账户开辟协程跑QPS，不然会阻塞for
		for _, item := range GLoginedUsers.Players {
			//进入访问资源
			//gSemOrder.Enter()
			//登陆成功，上下分请求
			sendHTTPOrderRequest(item.UserID, item.Account, item.AgentID)
		}
	}()
}

//
type HTTPAuthResult struct {
	Type     float64       `json:"type,omitempty"`
	Maintype string        `json:"maintype,omitempty"`
	Data     *HTTPAuthData `json:"data,omitempty"`
}

//
type HTTPAuthData struct {
	Code    float64 `json:"code,omitempty"`
	Account string  `json:"account,omitempty"`
	URL     string  `json:"url,omitempty"`
	IP      string  `json:"domain,omitempty"`
	Port    string  `json:"port,omitempty"`
	Token   string  `json:"token,omitempty"`
}

//HTTPGetTokenTest 查询token压测
//-------------------------------------------------------------
func HTTPGetTokenTest(httpaddr string) (e error) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln(debug.Stack())
		}
	}()
	requrl := fmt.Sprintf("http://%s/TxVoidHandle", httpaddr)
	//log.Printf("--- *** PID[%07d] HTTPGetToken >>> %v", os.Getpid(), requrl)
	//client := http.Client{Timeout: time.Duration(*httptimeout) * time.Second}
	//rsp, err := client.Get(requrl)
	rsp, err := http.Get(requrl)
	if err != nil {
		//log.Printf("--- *** PID[%07d] HTTPGetToken httpGet %v\n", os.Getpid(), err)
		e = err
		return
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		//log.Printf("--- *** PID[%07d] HTTPGetToken ReadAll %v\n", os.Getgid(), err)
		e = err
		return
	}
	/*str :=*/ util.Byte2Str(body)
	//log.Println(str)
	//log.Printf("--- *** PID[%07d] HTTPGetToken <<< %v", os.Getpid(), str)
	return
}

//sendHTTPGetTokenRequest 开辟协程获取token请求
func sendHTTPGetTokenRequest(httpaddr string, account int64, agentID int) {
	//为当前账户开辟一个协程，间隔1s发送HTTP请求
	go func(httpaddr string, account int64, agentID int) {
		i, t := 0, 5 //CPU分片调度
		//间隔1s发送HTTP请求
		for {
			if i > t {
				i = 0
				//进程多开时腾出CPU时间片给其它进程
				runtime.Gosched()
			}
			i++
			//间隔1s发送，不然OS内核发送缓存区被打满，引起I/O超时异常
			time.Sleep(100 * time.Millisecond)
			//跑QPS性能数据
			sendHTTPGetTokenRequestQPS(httpaddr, account, agentID)
		}
	}(httpaddr, account, agentID)
}

//sendHTTPGetTokenRequest 上下分请求
func sendHTTPGetTokenRequestQPS(httpaddr string, account int64, agentID int) (e error) {
	e = HTTPGetTokenTest(httpaddr)
	//_, _, e = HTTPGetToken(httpaddr, account, agentID)
	return
}

//HTTPGetToken 客户端 - 查询token
//HTTPGetToken ipaddr 网关ipaddr
//-------------------------------------------------------------
func HTTPGetToken(httpaddr string, account int64, agentID int) (token, ipaddr string, e error) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln(debug.Stack())
		}
	}()
	requrl := fmt.Sprintf("http://%s/GameHandle?testAccount=%d&agentid=%d", httpaddr, account, agentID)
	log.Printf("--- *** PID[%07d] HTTPGetToken >>> %v", os.Getpid(), requrl)
	//client := http.Client{Timeout: time.Duration(*httptimeout) * time.Second}
	//rsp, err := client.Get(requrl)
	rsp, err := http.Get(requrl)
	if err != nil {
		//log.Printf("--- *** PID[%07d] HTTPGetToken httpGet %v\n", os.Getpid(), err)
		e = err
		return
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		e = err
		//log.Printf("--- *** PID[%07d] HTTPGetToken ReadAll %v\n", os.Getgid(), err)
		return
	}
	str := util.Byte2Str(body)
	//log.Println(str)
	str = strings.Replace(str, "\\", "", -1)
	//str = str[1 : len(str)-1]
	body = util.Str2Byte(str)
	var authResult HTTPAuthResult
	if err := util.Byte2JSON(body, &authResult); err != nil {
		//log.Println("----->>>> ", str)
		//log.Printf("--- *** PID[%07d] HTTPGetToken Byte2JSON %v", os.Getpid(), err)
		e = err
		return
	}
	// values, _ := url.ParseQuery(authResult.Data.URL)
	// //log.Println(values)
	// for url := range values {
	// 	sub := url[strings.Index(url, "?")+1:]
	// 	dic := map[string]string{}
	// 	for {
	// 		s := strings.Index(sub, "=")
	// 		if s == -1 {
	// 			break
	// 		}
	// 		p := strings.Index(sub, "&")
	// 		if p == -1 {
	// 			dic[sub[0:s]] = sub[s+1:]
	// 			break
	// 		} else {
	// 			dic[sub[0:s]] = sub[s+1 : p]
	// 		}
	// 		sub = sub[p+1:]
	// 	}
	// 	authResult.Data.IP = dic["domain"]
	// 	authResult.Data.Port = dic["port"]
	// 	authResult.Data.Token = dic["token"]
	// 	break
	// }
	// //token
	// token = authResult.Data.Token
	// //网关ipaddr
	// ipaddr = authResult.Data.IP + ":" + authResult.Data.Port
	// //log.Printf("--- *** PID[%07d] token >>> %v", os.Getpid(), token)
	url := authResult.Data.URL
	sub := url[strings.Index(url, "?")+1:]
	dic := map[string]string{}
	for {
		s := strings.Index(sub, "=")
		if s == -1 {
			break
		}
		p := strings.Index(sub, "&")
		if p == -1 {
			dic[sub[0:s]] = sub[s+1:]
			break
		} else if p > s {
			dic[sub[0:s]] = sub[s+1 : p]
			sub = sub[p+1:]
		} else {
			break
		}
	}
	requrl2 := fmt.Sprintf("http://%s/TokenHandle?token=%s&descode=%s&white=%s&versions=%s&logintype=%s&gameid=%s",
		httpaddr, dic["token"], dic["descode"], dic["white"], dic["versions"], dic["logintype"], dic["gameid"])
	//client2 := http.Client{Timeout: time.Duration(*httptimeout) * time.Second}
	//rsp2, err := client2.Get(requrl2)
	rsp2, err := http.Get(requrl2)
	if err != nil {
		//log.Printf("--- *** PID[%07d] HTTPGetToken httpGet %v\n", os.Getpid(), err)
		e = err
		return
	}
	defer rsp2.Body.Close()
	body2, err := ioutil.ReadAll(rsp2.Body)
	if err != nil {
		e = err
		//log.Printf("--- *** PID[%07d] HTTPGetToken ReadAll %v\n", os.Getgid(), err)
		return
	}
	ipaddr = dic["domain"] + ":" + dic["port"]
	str2 := util.Byte2Str(body2)
	token = str2[strings.Index(str2, "=")+1:]
	if *dynamic == 0 {
		log.Printf("--- *** PID[%07d] HTTPGetToken <<< token[%v] wsaddr[%v]", os.Getpid(), token, *wsaddr)
	} else {
		log.Printf("--- *** PID[%07d] HTTPGetToken <<< token[%v] wsaddr[%v]", os.Getpid(), token, ipaddr)
	}
	return
}

//HTTPOrderRequest 上下分请求
//-------------------------------------------------------------
func HTTPOrderRequest(httpaddr string, Type int, agentID uint32, timestamp int64, paramVal, key string) (e error) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln(debug.Stack())
		}
	}()
	requrl := fmt.Sprintf("http://%s/GameHandle?type=%d&agentid=%d&timestamp=%v&paraValue=%v&key=%v",
		httpaddr, Type, agentID, timestamp, paramVal, key)
	//log.Printf("--- *** PID[%07d] HTTPOrderRequest >>> %v", os.Getpid(), requrl)
	client := http.Client{Timeout: time.Duration(*httptimeout) * time.Second}
	rsp, err := client.Get(requrl)
	//rsp, err := http.Get(requrl)
	if err != nil {
		//log.Printf("--- *** PID[%07d] HTTPOrderRequest httpGet %v\n", os.Getpid(), err)
		e = err
		//离开释放资源
		//gSemOrder.Leave()
		return
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		e = err
		//log.Printf("--- *** PID[%07d] HTTPOrderRequest ReadAll %v\n", os.Getgid(), err)
		//离开释放资源
		//gSemOrder.Leave()
		return
	}
	/*str :=*/ util.Byte2Str(body)
	//log.Printf("--- *** PID[%07d] HTTPOrderRequest <<< %s", os.Getpid(), str)
	//离开释放资源
	//gSemOrder.Leave()
	return
}

//HTTPOrderRequest2 上下分请求
//-------------------------------------------------------------
func HTTPOrderRequest2(httpaddr string, Type int, orderID string, agentID uint32, userID int64, account int64, score int64) (e error) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln(debug.Stack())
		}
	}()
	requrl := fmt.Sprintf("http://%s/GameHandle?type=%d&orderid=%s&agentid=%d&userid=%v&account=%v&score=%v",
		httpaddr, Type, orderID, agentID, userID, account, score)
	//log.Printf("--- *** PID[%07d] HTTPOrderRequest2 >>> %v", os.Getpid(), requrl)
	client := http.Client{Timeout: time.Duration(*httptimeout) * time.Second}
	rsp, err := client.Get(requrl)
	//rsp, err := http.Get(requrl)
	if err != nil {
		//log.Printf("--- *** PID[%07d] HTTPOrderRequest2 httpGet %v\n", os.Getpid(), err)
		e = err
		//离开释放资源
		//gSemOrder.Leave()
		return
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		e = err
		//log.Printf("--- *** PID[%07d] HTTPOrderRequest2 ReadAll %v\n", os.Getgid(), err)
		//离开释放资源
		//gSemOrder.Leave()
		return
	}
	/*str :=*/ util.Byte2Str(body)
	//log.Printf("--- *** PID[%07d] HTTPOrderRequest2 <<< %s", os.Getpid(), str)
	//离开释放资源
	//gSemOrder.Leave()
	return
}

//sendHTTPOrderRequestQPS 上下分请求
func sendHTTPOrderRequestQPS(userID int64, account int64, agentID uint32) (e error) {
	var opType int
	if (rand.Intn(100) % 2) == 0 {
		opType = 2 //上分
	} else {
		opType = 3 //下分
	}
	if *isdecrypt != 0 {
		//加密方式HTTP请求 ///
		timestamp := TimeNow().SinceUnixEpoch()
		orderID := util.RandomNumberStr(20)
		score := rand.Int63n(2) + 1
		//str := fmt.Sprintf("%v%v%v", agentID, timestamp, *md5code)
		//md5编码 bug???
		key := util.MD5(fmt.Sprintf("%v%v%v", agentID, timestamp, *md5code), true)
		//log.Printf("+++++++++++++++++++\nagentID=%v\ntimestamp=%v\nmd5code=%v\nstr=%v\nkey=%v\n\n", agentID, timestamp, *md5code, str, key)
		//AES加密
		rawParaValue := fmt.Sprintf("userid=%v&account=%v&orderid=%v&score=%v", userID, account, orderID, score)
		encrypt := util.AesEncrypt(rawParaValue, *descode)
		//Base64编码
		strBase64 := util.Base64Encode(encrypt)
		//b, _ := util.Base64Decode(strBase64)
		//decrypt := util.AESDecrypt(b, util.Str2Byte(*descode))
		//fmt.Print(util.Byte2Str(decrypt))
		//UrlEncode ///
		paraValue := util.URLEncode(strBase64)
		if atomic.LoadInt64(&gOrderTimeStart) == 0 {
			//起始时间戳(ms)
			atomic.StoreInt64(&gOrderTimeStart, TimeNowMilliSec().SinceUnixEpoch())
		}
		//本次请求开始时间戳(ms) ///
		timestart := TimeNowMilliSec()
		e = HTTPOrderRequest(*httpaddr1, opType, agentID, timestamp, paraValue, key)
		//本次请求结束时间戳(ms) ///
		timenow := TimeNowMilliSec()
		//统计请求次数 ///
		atomic.AddInt64(&gOrderRequestNum, 1)
		//历史请求次数 ///
		atomic.AddInt64(&gOrderRequestNumTotal, 1)
		//本次请求是否成功 ///
		resultStr := "[OK]"
		if e == nil {
			//统计成功次数 ///
			atomic.AddInt64(&gOrderRequestNumSucc, 1)
			//历史成功次数 ///
			atomic.AddInt64(&gOrderRequestNumTotalSucc, 1)
		} else {
			resultStr = "[ERR]"
			//统计失败次数 ///
			atomic.AddInt64(&gOrderRequestNumFailed, 1)
			//历史失败次数 ///
			atomic.AddInt64(&gOrderRequestNumTotalFailed, 1)
		}
		//统计间隔时间(ms) 间隔时间(ms)打印一次 ///
		totalTime := TimeDiff(timenow, NewTimestamp(gOrderTimeStart))
		if totalTime >= int32(*deltaTime) {
			//最近一次请求耗时(ms) ///
			timdiff := TimeDiff(timenow, timestart)
			//统计请求次数 ///
			requestNum := atomic.LoadInt64(&gOrderRequestNum)
			//统计成功次数 ///
			requestNumSucc := atomic.LoadInt64(&gOrderRequestNumSucc)
			//统计失败次数 ///
			requestNumFailed := atomic.LoadInt64(&gOrderRequestNumFailed)
			//统计命中率 ///
			ratio := (float32)(requestNumSucc) / (float32)(requestNum)
			//历史请求次数 ///
			requestNumTotal := atomic.LoadInt64(&gOrderRequestNumTotal)
			//历史成功次数 ///
			requestNumTotalSucc := atomic.LoadInt64(&gOrderRequestNumTotalSucc)
			//历史失败次数 ///
			requestNumTotalFailed := atomic.LoadInt64(&gOrderRequestNumTotalFailed)
			//历史命中率 ///
			ratioTotal := (float32)(requestNumTotalSucc) / (float32)(requestNumTotal)
			//平均请求耗时(ms) ///
			avgTime := (float32)(totalTime) / (float32)(requestNum)
			//每秒请求次数(QPS) ///
			avgNum := (int32)(requestNum) / (totalTime / 1000)
			pid := os.Getpid()
			log.Printf("\n\n--- *** ------------------------------------------------------\n")
			log.Printf(
				"\n--- *** PID[%07d][注单]本次统计间隔时间[%v]ms 超时值[%v]s\n"+
					"--- *** PID[%07d][注单]本次统计请求次数[%v] 成功[%v] 失败[%v] 命中率[%v]\n"+
					"--- *** PID[%07d][注单]最近一次请求耗时[%v]ms %v\n"+
					"--- *** PID[%07d][注单]平均请求耗时[%v]ms\n"+
					"--- *** PID[%07d][注单]每秒请求次数(QPS) = [%v]\n"+
					"--- *** PID[%07d][注单]历史请求次数[%v] 成功[%v] 失败[%v] 命中率[%v]\n\n",
				pid, totalTime, *httptimeout, pid, requestNum, requestNumSucc, requestNumFailed, ratio,
				pid, timdiff, resultStr, pid, avgTime, pid, avgNum, pid,
				requestNumTotal, requestNumTotalSucc, requestNumTotalFailed, ratioTotal)
			//重置起始时间戳(ms) ///
			atomic.StoreInt64(&gOrderTimeStart, timenow.SinceUnixEpoch())
			//重置统计请求次数 ///
			atomic.StoreInt64(&gOrderRequestNum, 0)
			//重置统计成功次数 ///
			atomic.StoreInt64(&gOrderRequestNumSucc, 0)
			//重置统计失败次数 ///
			atomic.StoreInt64(&gOrderRequestNumFailed, 0)
		}
	} else {
		//非加密方式HTTP请求 ///
		orderID := util.RandomNumberStr(20)
		score := rand.Int63n(2) + 1
		if atomic.LoadInt64(&gOrderTimeStart) == 0 {
			//起始时间戳(ms)
			atomic.StoreInt64(&gOrderTimeStart, TimeNowMilliSec().SinceUnixEpoch())
		}
		//本次请求开始时间戳(ms) ///
		timestart := TimeNowMilliSec()
		e = HTTPOrderRequest2(*httpaddr1, opType, orderID, agentID, userID, account, score)
		//本次请求结束时间戳(ms) ///
		timenow := TimeNowMilliSec()
		//统计请求次数 ///
		atomic.AddInt64(&gOrderRequestNum, 1)
		//历史请求次数 ///
		atomic.AddInt64(&gOrderRequestNumTotal, 1)
		//本次请求是否成功 ///
		resultStr := "[OK]"
		if e == nil {
			//统计成功次数 ///
			atomic.AddInt64(&gOrderRequestNumSucc, 1)
			//历史成功次数 ///
			atomic.AddInt64(&gOrderRequestNumTotalSucc, 1)
		} else {
			resultStr = "[ERR]"
			//统计失败次数 ///
			atomic.AddInt64(&gOrderRequestNumFailed, 1)
			//历史失败次数 ///
			atomic.AddInt64(&gOrderRequestNumTotalFailed, 1)
		}
		//统计间隔时间(ms) 间隔时间(ms)打印一次 ///
		totalTime := TimeDiff(timenow, NewTimestamp(gOrderTimeStart))
		if totalTime >= int32(*deltaTime) {
			//最近一次请求耗时(ms) ///
			timdiff := TimeDiff(timenow, timestart)
			//统计请求次数 ///
			requestNum := atomic.LoadInt64(&gOrderRequestNum)
			//统计成功次数 ///
			requestNumSucc := atomic.LoadInt64(&gOrderRequestNumSucc)
			//统计失败次数 ///
			requestNumFailed := atomic.LoadInt64(&gOrderRequestNumFailed)
			//统计命中率 ///
			ratio := requestNumSucc / requestNum
			//历史请求次数 ///
			requestNumTotal := atomic.LoadInt64(&gOrderRequestNumTotal)
			//历史成功次数 ///
			requestNumTotalSucc := atomic.LoadInt64(&gOrderRequestNumTotalSucc)
			//历史失败次数 ///
			requestNumTotalFailed := atomic.LoadInt64(&gOrderRequestNumTotalFailed)
			//历史命中率 ///
			ratioTotal := requestNumTotalSucc / requestNumTotal
			//平均请求耗时(ms) ///
			avgTime := (float32)(totalTime) / (float32)(requestNum)
			//每秒请求次数(QPS) ///
			avgNum := (int32)(requestNum) / (totalTime / 1000)
			pid := os.Getpid()
			log.Printf("\n\n--- *** ------------------------------------------------------\n")
			log.Printf(
				"\n--- *** PID[%07d][注单]本次统计间隔时间[%v]ms 超时值[%v]s\n"+
					"--- *** PID[%07d][注单]本次统计请求次数[%v] 成功[%v] 失败[%v] 命中率[%v]\n"+
					"--- *** PID[%07d][注单]最近一次请求耗时[%v]ms %v\n"+
					"--- *** PID[%07d][注单]平均请求耗时[%v]ms\n"+
					"--- *** PID[%07d][注单]每秒请求次数(QPS) = [%v]\n"+
					"--- *** PID[%07d][注单]历史请求次数[%v] 成功[%v] 失败[%v] 命中率[%v]\n\n",
				pid, totalTime, *httptimeout, pid, requestNum, requestNumSucc, requestNumFailed, ratio,
				pid, timdiff, resultStr, pid, avgTime, pid, avgNum, pid,
				requestNumTotal, requestNumTotalSucc, requestNumTotalFailed, ratioTotal)
			//重置起始时间戳(ms) ///
			atomic.StoreInt64(&gOrderTimeStart, timenow.SinceUnixEpoch())
			//重置统计请求次数 ///
			atomic.StoreInt64(&gOrderRequestNum, 0)
			//重置统计成功次数 ///
			atomic.StoreInt64(&gOrderRequestNumSucc, 0)
			//重置统计失败次数 ///
			atomic.StoreInt64(&gOrderRequestNumFailed, 0)
		}
	}
	return
}

//sendHTTPOrderRequest 开辟协程上下分请求
//-------------------------------------------------------------
func sendHTTPOrderRequest(userID int64, account int64, agentID uint32) {
	//为当前账户开辟一个协程，间隔1s发送HTTP上下分请求
	go func(userID int64, account int64, agentID uint32) {
		i, t := 0, 5 //CPU分片调度
		//间隔1s发送HTTP上下分请求
		for {
			if i > t {
				i = 0
				//进程多开时腾出CPU时间片给其它进程
				runtime.Gosched()
			}
			i++
			//间隔1s发送，不然OS内核发送缓存区被打满，引起I/O超时异常
			time.Sleep(time.Second)
			//跑QPS性能数据
			sendHTTPOrderRequestQPS(userID, account, agentID)
		}
	}(userID, account, agentID)
}

//sendPlayerLogin 登陆大厅
//-------------------------------------------------------------
func sendPlayerLogin(peer Session, token string) {
	reqdata := &HallServer.LoginMessage{}
	val, _ := strconv.ParseUint("F5F5F5F5", 16, 32)
	reqdata.Header = &Game_Common.Header{}
	reqdata.Header.Sign = int32(val)
	reqdata.Session = token
	msg := newMsg(
		uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_HALL),
		uint8(Game_Common.MESSAGE_CLIENT_TO_HALL_SUBID_CLIENT_TO_HALL_LOGIN_MESSAGE_REQ),
		reqdata)
	util.Log("UserClient", "DefWSClient", "sendPlayerLogin", reqdata)
	peer.Write(msg)
}

//sendKeepAlive 发送心跳包
//-------------------------------------------------------------
func sendKeepAlive(peer Session, token string) {
	reqdata := &Game_Common.KeepAliveMessage{}
	val, _ := strconv.ParseUint("F5F5F5F5", 16, 32)
	reqdata.Header = &Game_Common.Header{}
	reqdata.Header.Sign = int32(val)
	reqdata.Session = token //
	msg := newMsg(
		uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_HALL),
		uint8(Game_Common.MESSAGE_CLIENT_TO_SERVER_SUBID_KEEP_ALIVE_REQ),
		reqdata)
	//util.Log("UserClient", "Player", "sendKeepAlive", reqdata)
	peer.Write(msg)
}

//reqGameListInfo 取游戏信息
//-------------------------------------------------------------
func reqGameListInfo(peer Session) {
	reqdata := &HallServer.GetGameMessage{}
	val, _ := strconv.ParseUint("F5F5F5F5", 16, 32)
	reqdata.Header = &Game_Common.Header{}
	reqdata.Header.Sign = int32(val)
	msg := newMsg(
		uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_HALL),
		uint8(Game_Common.MESSAGE_CLIENT_TO_HALL_SUBID_CLIENT_TO_HALL_GET_GAME_ROOM_INFO_REQ),
		reqdata)
	util.Log("UserClient", "Player", "reqGameListInfo", reqdata)
	peer.Write(msg)
}

//reqGameserverInfo 获取游戏IP
//-------------------------------------------------------------
func reqGameserverInfo(peer Session, gameID, roomID int32) {
	reqdata := &HallServer.GetGameServerMessage{}
	val, _ := strconv.ParseUint("F5F5F5F5", 16, 32)
	reqdata.Header = &Game_Common.Header{}
	reqdata.Header.Sign = int32(val)
	reqdata.GameId = uint32(gameID) //游戏类型
	reqdata.RoomId = uint32(roomID) //游戏房间
	client := peer.GetCtx(TagUserInfo).(*DefWSClient)
	client.GameID = gameID //保存
	client.RoomID = roomID
	msg := newMsg(
		uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_HALL),
		uint8(Game_Common.MESSAGE_CLIENT_TO_HALL_SUBID_CLIENT_TO_HALL_GET_GAME_SERVER_MESSAGE_REQ),
		reqdata)
	util.Log("UserClient", "Player", "reqGameserverInfo", reqdata)
	log.Printf("--- *** PID[%07d] player[%d:%d:%s] %d:%s %d:%s \n",
		os.Getpid(),
		client.UserID,
		client.Account,
		client.Token,
		//
		gameID,
		GGames.ByID[gameID].Name,
		//
		roomID,
		GGames.ByID[gameID].ByID[roomID])
	peer.Write(msg)
}

//reqEnterRoom 进入房间
//-------------------------------------------------------------
func reqEnterRoom(peer Session, gameID, roomID int32, pwd []byte) {
	reqdata := &GameServer.MSG_C2S_UserEnterMessage{}
	val, _ := strconv.ParseUint("F5F5F5F5", 16, 32)
	reqdata.Header = &Game_Common.Header{}
	reqdata.Header.Sign = int32(val)
	reqdata.GameId = gameID          //游戏类型
	reqdata.RoomId = roomID          //游戏房间
	reqdata.DynamicPassword = pwd[:] //动态密码
	msg := newMsg(
		uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_GAME_SERVER),
		uint8(GameServer.SUBID_SUB_C2S_ENTER_ROOM_REQ),
		reqdata)
	util.Log("UserClient", "Player", "reqEnterRoom", reqdata)
	peer.Write(msg)
}

//reqPlayerReady 玩家就绪
//-------------------------------------------------------------
func reqPlayerReady(peer Session) {
	reqdata := &GameServer.MSG_C2S_UserReadyMessage{}
	val, _ := strconv.ParseUint("F5F5F5F5", 16, 32)
	reqdata.Header = &Game_Common.Header{}
	reqdata.Header.Sign = int32(val)
	msg := newMsg(
		uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_GAME_SERVER),
		uint8(GameServer.SUBID_SUB_C2S_USER_READY_REQ),
		reqdata)
	util.Log("UserClient", "Player", "reqPlayerReady", reqdata)
	peer.Write(msg)
}

//reqPlayerLeave 玩家离开
//-------------------------------------------------------------
func reqPlayerLeave(peer Session, userID int32, gameID, roomID, Type int32) {
	reqdata := &GameServer.MSG_C2S_UserLeftMessage{}
	val, _ := strconv.ParseUint("F5F5F5F5", 16, 32)
	reqdata.Header = &Game_Common.Header{}
	reqdata.Header.Sign = int32(val)
	reqdata.UserId = uint32(userID)
	reqdata.GameId = gameID
	reqdata.RoomId = roomID
	reqdata.Type = Type
	msg := newMsg(
		uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_GAME_SERVER),
		uint8(GameServer.SUBID_SUB_C2S_USER_LEFT_REQ),
		reqdata)
	//util.Log("UserClient", "Player", "reqPlayerLeave", reqdata)
	peer.Write(msg)
}
