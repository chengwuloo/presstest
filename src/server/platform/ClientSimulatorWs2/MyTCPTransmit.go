package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import (
	"encoding/binary"
	"errors"
	"log"
	"net"

	"github.com/davyxu/cellnet/codec"
)

//
type MyTCPTransmit struct {
}

//
func NewMyTCPTransmit() MsgTransmit {
	return &MyTCPTransmit{}
}

//
func (s *MyTCPTransmit) OnRecvMessage(peer Session) (msg interface{}, err error) {
	conn, ok := peer.Conn().(net.Conn)
	if !ok || conn == nil {
		return nil, nil
	}
	//len+CRC，4字节
	buf := make([]byte, 4)
	err = ReadFull(conn, buf)
	if err != nil {
		log.Fatalln("OnRecvMessage: ", err)
		return nil, err
	}
	//len，2字节
	len := binary.LittleEndian.Uint16(buf[:2])
	//CRC，2字节
	checksum := binary.LittleEndian.Uint16(buf[2:])
	// 读取剩余大小
	remain := make([]byte, len-4)
	err = ReadFull(conn, remain)
	if err != nil {
		log.Fatalln("OnRecvMessage: ", err)
		return nil, err
	}
	//CRC校验
	crc := GetChecksum(remain)
	if crc != checksum {
		log.Fatalln("OnRecvMessage: RecvPacket GetChecksum error")
		return nil, errors.New("RecvPacket GetChecksum error")
	}
	//版本0x0001
	ver := binary.LittleEndian.Uint16(remain[:2])
	//标记0x5F5F
	sign := binary.LittleEndian.Uint16(remain[2:4])
	//主命令ID
	mainID := uint8(remain[4])
	//子命令ID
	subID := uint8(remain[5])
	//加密类型
	encryptTy := uint8(remain[6])
	//预留字段
	reserve := uint8(remain[7])
	//请求ID
	reqID := binary.LittleEndian.Uint32(remain[8:12])
	//实际大小
	realSize := binary.LittleEndian.Uint16(remain[12:14])
	log.Printf("ver:%#x\nsign:%#x\nmainID:%d\nsubID:%d\nencTy:%#x\nreserv:%d\nreqID:%d\nrealSize:%d\n",
		ver, sign, mainID, subID, encryptTy, reserve, reqID, realSize)
	//实际protobuf数据
	data := remain[14:]
	msg, _, err = codec.DecodeMessage(int(subID), data)
	return msg, err
}

//
func (s *MyTCPTransmit) OnSendMessage(peer Session, msg interface{}) error {
	conn, ok := peer.Conn().(net.Conn)
	if !ok || conn == nil {
		return nil
	}
	log.Println("MyTCPTransmit::OnSendMessage\n", msg)
	h, ok := msg.(*Msg)
	if !ok || h == nil {
		return nil
	}
	data, _, err := codec.EncodeMessage(h.msg, nil)
	if err != nil {
		log.Fatalln("EncodeMessage : ", err)
		return err
	}
	buf := make([]byte, 18+len(data))
	//len，2字节
	length := 18 + len(data)
	binary.LittleEndian.PutUint16(buf[0:], uint16(length))
	//CRC，2字节
	//binary.LittleEndian.PutUint16(buf[2:], h.crc)
	//版本0x0001
	binary.LittleEndian.PutUint16(buf[4:], uint16(h.ver))
	//标记0x5F5F
	binary.LittleEndian.PutUint16(buf[6:], uint16(h.sign))
	//主命令ID
	buf[8] = byte(h.mainID)
	//子命令ID
	buf[9] = byte(h.subID)
	//加密类型
	buf[10] = byte(h.encType)
	//预留字段
	buf[11] = byte(0x01) //
	//请求ID
	binary.LittleEndian.PutUint32(buf[12:], uint32(0)) //
	//实际大小
	binary.LittleEndian.PutUint16(buf[16:], uint16(len(data)))
	//实际数据
	copy(buf[18:], data)
	//CRC，2字节
	crc := GetChecksum(buf[4:])
	binary.LittleEndian.PutUint16(buf[2:], crc)
	return WriteFull(conn, buf)
}
