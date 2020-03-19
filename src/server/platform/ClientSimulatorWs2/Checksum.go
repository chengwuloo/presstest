package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//  来自C++服务器框架协议

import "encoding/binary"

//GetChecksum CRC校验和
func GetChecksum(data []byte) uint16 {
	var sum uint16
	idx := 0
	size := len(data)
	for i := 0; i < size/2; i++ {
		//读取uint16，2字节
		sum += binary.LittleEndian.Uint16(data[idx:])
		idx += 2
	}
	if size%2 != 0 {
		//读取uint8，1字节
		sum += uint16(data[idx])
	}
	return sum
}
