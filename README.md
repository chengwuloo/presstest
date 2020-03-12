# 压测工具
	go version go1.12.3 linux/amd64
	libprotoc 3.7.0
## 说明
	服务器压测程序框架组件(golang版本)
	MsgProc.go 消息处理
	MsgProcCell.go 消息处理单元
	Mailbox.go(MsgProcPool.go) 邮槽(消息处理池)
	Worker.go 工作类(逻辑业务类继承入口)
	TimeWheel.go 时间轮用来处理超时连接
	ScopedTimer.go 线程局部定时器
	MsgTransmit.go 消息协议解析入口
	MsgQueue.go 消息队列
	Session.go 连接会话
	SessionMgr.go 会话管理类
	Timestamp.go 时间戳
	
	握手/交互协议：websocket&protobuf
	可以根据需要自定义修改底层通信协议
	LoadClientSimulatorWs 负责加载 ClientSimulatorWs2 子进程
	配置在LoadClientSimulatorWs/conf.ini下

	一个worker线程处理N个conn，
	一个conn关联一个worker线程(hash运算或GetNextCell指定)，	
	没有其它worker线程竞争，没有线程切换开销，
	每个conn上逻辑业务处理有序，无锁高效
	
	1.listen/connect主线程 ///////////////////
	2.网络IO线程(M)，IO收发(recv/send) ///////
	3.worker线程(N)，处理游戏业务逻辑 ////////

## 特别注意
	如果改动框架协议后，需要使用win32_proto.bat/win64_proto.bat重新生成框架协议，
	然后手动修改HallServer.Message.pb.go和GameServer.Message.pb.go：
		import时添加 "server/pb/Game_Common"
		*Header 全部替换为 *Game_Common.Header
		&Header{} 全部替换为 &Game_Common.Header{}
	默认注释掉了win32_proto.bat/win64_proto.bat中生成框架协议脚本代码，
	除非框架协议改动，需要开启重新生成框架协议，并按照上面的方式修改相应协议文件
	
## 扩展方式
sdata.go handle.go msg.go player.go里面添加相应模块就可以了
* [sdata.go](http://192.168.2.210:12345/server/presstest/blob/master/src/server/platform/ClientSimulatorWs2/sdata.go): 添加子游戏类型(名称找ID，ID找名称两处)
* [msg.go](http://192.168.2.210:12345/server/presstest/blob/master/src/server/platform/ClientSimulatorWs2/Msg.go): 注册协议消息
* [handle.go](http://192.168.2.210:12345/server/presstest/blob/master/src/server/platform/ClientSimulatorWs2/handler.go): 发送消息处理
* [player.go](http://192.168.2.210:12345/server/presstest/blob/master/src/server/platform/ClientSimulatorWs2/player.go): 接收消息应答处理
