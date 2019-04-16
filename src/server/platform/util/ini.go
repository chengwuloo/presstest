package util

//
// Created by YangZhi
// 			4/9/2019
//

import (
	"gopkg.in/ini.v1"
)

//
type Ini struct {
	reader *ini.File
}

//
// type IniParserError struct {
// 	err string
// }

//
//func (s *IniParserError) Error() string { return s.err }

//
func (s *Ini) Load(filename string) error {
	conf, err := ini.Load(filename)
	if err != nil {
		s.reader = nil
		return err
	}
	s.reader = conf
	return nil
}

//
func (s *Ini) GetString(section string, key string) string {
	if s.reader == nil {
		return ""
	}

	x := s.reader.Section(section)
	if x == nil {
		return ""
	}

	return x.Key(key).String()
}

//
func (s *Ini) GetInt(section string, key string) int {
	if s.reader == nil {
		return 0
	}

	x := s.reader.Section(section)
	if x == nil {
		return 0
	}

	val, _ := x.Key(key).Int()

	return val
}

//
func (s *Ini) GetInt32(section string, key string) int32 {
	if s.reader == nil {
		return 0
	}

	x := s.reader.Section(section)
	if x == nil {
		return 0
	}

	val, _ := x.Key(key).Int()

	return int32(val)
}

//
func (s *Ini) GetUint32(section string, key string) uint32 {
	if s.reader == nil {
		return 0
	}

	x := s.reader.Section(section)
	if x == nil {
		return 0
	}

	val, _ := x.Key(key).Uint()

	return uint32(val)
}

//
func (s *Ini) GetInt64(section string, key string) int64 {
	if s.reader == nil {
		return 0
	}

	x := s.reader.Section(section)
	if x == nil {
		return 0
	}

	val, _ := x.Key(key).Int64()
	return val
}

//
func (s *Ini) GetUint64(section string, key string) uint64 {
	if s.reader == nil {
		return 0
	}

	x := s.reader.Section(section)
	if x == nil {
		return 0
	}

	val, _ := x.Key(key).Uint64()
	return val
}

//
func (s *Ini) GetFloat32(section string, key string) float32 {
	if s.reader == nil {
		return 0
	}

	x := s.reader.Section(section)
	if x == nil {
		return 0
	}

	val, _ := x.Key(key).Float64()
	return float32(val)
}

//
func (s *Ini) GetFloat64(section string, key string) float64 {
	if s.reader == nil {
		return 0
	}

	x := s.reader.Section(section)
	if x == nil {
		return 0
	}

	val, _ := x.Key(key).Float64()
	return val
}
