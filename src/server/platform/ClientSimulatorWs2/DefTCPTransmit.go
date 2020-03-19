package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import (
	"encoding/binary"
	"log"
	"net"

	"github.com/davyxu/cellnet/codec"
)

//
type DefTCPTransmit struct {
}

//
func NewDefTCPTransmit() MsgTransmit {
	return &DefTCPTransmit{}
}

//
func (s *DefTCPTransmit) OnRecvMessage(peer Session) (msg interface{}, err error) {
	conn, ok := peer.Conn().(net.Conn)
	if !ok || conn == nil {
		return nil, nil
	}
	//len+'\r\r'，6字节
	buf := make([]byte, 6)
	err = ReadFull(conn, buf)
	if err != nil {
		log.Fatalln("OnRecvMessage: ", err)
		return nil, err
	}
	//len，4字节
	len := binary.LittleEndian.Uint32(buf[:4])
	//'\r\r'，2字节，校验len
	if buf[4] != '\r' && buf[5] != '\r' {
		return nil, nil
	}
	//剩余字节CRC+msgID+protobuf
	data := make([]byte, len-6)
	err = ReadFull(conn, data)
	if err != nil {
		log.Fatalln("OnRecvMessage: ", err)
		return nil, err
	}
	//CRC，2字节
	crc := binary.LittleEndian.Uint16(data[:2])
	//CRC校验msgID+protobuf
	chsum := GetChecksum(data[2:])
	if crc != chsum {
		log.Fatalln("OnRecvMessage: GetChecksum error")
		return nil, nil
	}
	//消息ID，4字节
	msgID := binary.LittleEndian.Uint32(data[2:])
	//protobuf数据
	msg, _, err = codec.DecodeMessage(int(msgID), data[6:])
	return
}

//
func (s *DefTCPTransmit) OnSendMessage(peer Session, msg interface{}) error {
	conn, ok := peer.Conn().(net.Conn)
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
	return WriteFull(conn, buf)
}
