package main

//
// Created by YangZhi
// 			4/9/2019
//

import "time"

// Timestamp 时间戳
type Timestamp interface {
	Valid() bool
	Add(sec int32) Timestamp
	Less(t Timestamp) bool
	Equal(t Timestamp) bool
	Greater(t Timestamp) bool
	SinceUnixEpoch() int64
}

//
type timeStamp struct {
	val int64 // second/millisecond/microsecond/nanosecond
}

//
func NewTimestamp(val int64) Timestamp {
	return &timeStamp{val: val}
}

//
func (s *timeStamp) Valid() bool {
	return s.val > int64(0)
}

//
func (s *timeStamp) Add(sec int32) Timestamp {
	s.val = s.val + int64(sec)
	return s
}

//
func (s *timeStamp) Less(t Timestamp) bool {
	return s.val < t.SinceUnixEpoch()
}

//
func (s *timeStamp) Equal(t Timestamp) bool {
	return s.val == t.SinceUnixEpoch()
}

//
func (s *timeStamp) Greater(t Timestamp) bool {
	return s.val > t.SinceUnixEpoch()
}

//
func (s *timeStamp) SinceUnixEpoch() int64 {
	return s.val
}

//
func TimeAdd(t Timestamp, val int32) Timestamp {
	return NewTimestamp(t.SinceUnixEpoch() + int64(val))
}

// TimeNow 当前时间(秒)
func TimeNow() Timestamp {
	return NewTimestamp(time.Now().Unix())
}

// TimeNowMilliSec 当前时间(毫秒)
func TimeNowMilliSec() Timestamp {
	return NewTimestamp(time.Now().UnixNano() / 1e6)
}

// TimeNowMicroSec 当前时间(微秒)
func TimeNowMicroSec() Timestamp {
	return NewTimestamp(time.Now().UnixNano() / 1e3)
}

// TimeNowNanoSec 当前时间(纳秒)
func TimeNowNanoSec() Timestamp {
	return NewTimestamp(time.Now().UnixNano())
}

// TimeDiff 前后间隔时间差(second/millisecond/microsecond/nanosecond)
func TimeDiff(high, low Timestamp) int32 {
	diff := int32(high.SinceUnixEpoch() - low.SinceUnixEpoch())
	return diff
}
