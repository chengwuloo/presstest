package main

//
// Created by YangZhi
// 			4/9/2019
//

import (
	"encoding/binary"
	"log"

	"github.com/davyxu/cellnet/codec"
	"github.com/gorilla/websocket"
)

//
type DefWsTransmit struct {
}

//
func NewDefWsTransmit() MsgTransmit {
	return &DefWsTransmit{}
}

//
func (s *DefWsTransmit) OnRecvMessage(peer Session) (msg interface{}, err error) {
	conn, ok := peer.Conn().(*websocket.Conn)
	if !ok || conn == nil {
		return nil, nil
	}
	//len+'\r\r'，6字节
	conn.SetReadLimit(6)
	msgType, buf, err := conn.ReadMessage()
	if err != nil {
		log.Fatalln("OnRecvMessage: ", err)
		return nil, err
	}
	if len(buf) != 6 {
		log.Fatalln("OnRecvMessage: checklen error")
		return nil, nil
	}
	//TextMessage/BinaryMessage
	if websocket.BinaryMessage != msgType {
		return nil, nil
	}
	//len，4字节
	length := binary.LittleEndian.Uint32(buf[:4])
	//'\r\r'，2字节，校验len
	if buf[4] != '\r' && buf[5] != '\r' {
		return nil, nil
	}
	//剩余字节CRC+msgID+protobuf
	conn.SetReadLimit(int64(length - 6))
	_, data, err := conn.ReadMessage()
	if err != nil {
		log.Fatalln("OnRecvMessage: ", err)
		return nil, err
	}
	if uint32(len(data)) != (length - 6) {
		log.Fatalln("OnRecvMessage: checklen error")
		return nil, nil
	}
	//CRC，2字节
	crc := binary.LittleEndian.Uint16(data[:2])
	//CRC校验msgID+protobuf
	chsum := GetChecksum(data[2:])
	if crc != chsum {
		log.Fatalln("OnRecvMessage: GetCheckSum error")
		return nil, nil
	}
	//消息ID，4字节
	msgID := binary.LittleEndian.Uint32(data[2:])
	//protobuf数据
	msg, _, err = codec.DecodeMessage(int(msgID), data[6:])
	return msg, err
}

//
func (s *DefWsTransmit) OnSendMessage(peer Session, msg interface{}) error {
	conn, ok := peer.Conn().(*websocket.Conn)
	if !ok || conn == nil {
		return nil
	}
	//protobuf数据
	data, meta, err := codec.EncodeMessage(msg, nil)
	if err != nil {
		log.Fatalln("EncodeMessage: ", err)
		return err
	}
	//消息ID，4字节
	msgID := meta.ID
	//len+'\r\r'+CRC+msgID+protobuf
	buf := make([]byte, 10+len(data))
	//len，4字节
	len := 10 + len(data)
	binary.LittleEndian.PutUint32(buf[0:], uint32(len))
	//'\r\r'，2字节，校验len
	copy(buf[4:], "\r\r")
	//CRC，2字节
	//binary.LittleEndian.PutUint16(buf[6:], crc)
	//消息ID，4字节
	binary.LittleEndian.PutUint32(buf[8:], uint32(msgID))
	//protobuf数据
	copy(buf[12:], data)
	//CRC，2字节
	crc := GetChecksum(buf[8:])
	binary.LittleEndian.PutUint16(buf[6:], crc)
	return conn.WriteMessage(websocket.BinaryMessage, buf)
}
