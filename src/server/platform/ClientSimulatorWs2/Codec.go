package main

//
// Created by YangZhi
// 			4/9/2019
//

import (
	"errors"
	"log"
	"server/pb/Game_Common"
	"time"

	"github.com/davyxu/cellnet/codec"
)

//RootMsg 便于框架处理
type RootMsg struct {
	Cmd  uint32 //消息ID
	Data []byte //protobuf流数据
}

//Decodec 解码
func Decodec(msg interface{}) (uint32, interface{}, error) {
	if msg == nil {
		return 0, nil, errors.New("Decodec msg == nil")
	}
	// 提取RootMsg
	p := msg.(*RootMsg)
	data, _, err := codec.DecodeMessage(int(p.Cmd), p.Data)
	return p.Cmd, data, err
}

//DecodecAddReadTask 解码再投递给任务处理单元
func DecodecAddReadTask(msg interface{}, peer Session) {
	cmd, data, err := Decodec(msg)
	if err == nil {
		if cmd == uint32(Game_Common.MESSAGE_CLIENT_TO_SERVER_SUBID_KEEP_ALIVE_REQ) {
			// 收到来自客户端SYN心跳包
		} else if cmd == uint32(Game_Common.MESSAGE_CLIENT_TO_SERVER_SUBID_KEEP_ALIVE_RES) {
			//收到来自服务端ACK心跳包
			if true {
				//扔到业务中处理，断点调试业务或队列太长会直接影响对端心跳检测
				peer.AddReadTask(cmd, data)
			} else {
				//直接处理
				client := peer.GetCtx(TagUserInfo).(*DefWSClient)
				//收到心跳包，更新桶元素
				client.Cursor = peer.GetCell().GetWorker().GetTimer().GetTimeWheel().UpdateBucket(
					client.Cursor, peer.ID(), int32(*timeout)/1000)
				//间隔时间继续发送心跳包，心跳超时清理在业务中处理
				time.AfterFunc(time.Duration(*heartbeat)*time.Millisecond, func() {
					client := peer.GetCtx(TagUserInfo).(*DefWSClient)
					sendKeepAlive(peer, client.Token)
				})
			}
		} else {
			//扔给任务处理
			peer.AddReadTask(cmd, data)
		}
	} else {
		mainID, subID := DEWORD(int(cmd))
		log.Printf("DecodecAddReadTask[cmd = %d mainID=%d subID=%d] ERR: %v 强制关闭连接 !!!", cmd, mainID, subID, err)
		// 消息解析出错(不明数据), 直接断开连接
		//log.Fatalln("DecodecAddReadTask ", err)
		//peer.Close()
	}
}
