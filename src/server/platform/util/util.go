package util

//
// Created by YangZhi
// 			4/9/2019
//

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//
func GoroutineID() uint32 {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("GoroutineID recover: %v\n", err)
		}
	}()
	// b := make([]byte, 64)
	// b = b[:runtime.Stack(b, false)]
	// b = bytes.TrimPrefix(b, []byte("goroutine "))
	// b = b[:bytes.IndexByte(b, ' ')]
	// n, err := strconv.ParseUint(string(b), 10, 64)
	// if err != nil {
	// 	panic(fmt.Sprintf("goroutineID panic: %v\n", err))
	// }
	// return int(n)
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	str := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	ID, err := strconv.Atoi(str)
	if err != nil {
		panic(fmt.Sprintf("GoroutineID panic: %v\n", err))
	}
	return uint32(ID)
}

//
func JSON2Byte(v interface{}) []byte {
	if b, err := json.Marshal(v); err == nil {
		return b
	}
	return nil
}

//
func Byte2JSON(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}

//
func Byte2Str(b []byte) string {
	return string(b[:])
}

//
func Str2Byte(str string) []byte {
	return []byte(str)
}

//
func JSON2Str(v interface{}) string {
	if b, err := json.Marshal(v); err == nil {
		return Byte2Str(b)
	}
	return ""
}

//
func Str2JSON(str string, v interface{}) error {
	return json.Unmarshal(Str2Byte(str), v)
}

//
func LoadJSON(filename string) interface{} {
	if b, err := ioutil.ReadFile(filename); err == nil {
		var v interface{}
		if err := Byte2JSON(b, &v); err == nil {
			return v
		}
	}
	return nil
}

// Encode interface{} 编码成 []byte
func Encode(data interface{}) ([]byte, error) {
	//b := new(bytes.Buffer)
	//err := binary.Write(b, binary.LittleEndian, data)
	// b := bytes.NewBuffer(nil)
	// err := gob.NewEncoder(b).Encode(data)
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(data)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), err
}

// Decode []byte解码成 interface{}
func Decode(stream []byte) (interface{}, error) {
	var data interface{}
	b := bytes.NewBuffer(stream)
	err := gob.NewDecoder(b).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}

// ToBytes interface{} 编码成 []byte
func ToBytes(data interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(data)
	if err != nil {
		log.Println("ToBytes Error:", err)
		return nil, err
	}
	return b.Bytes(), nil
}

// MD5 生成32位 MD5
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// randomseedInit 随机种子
func randomseedInit() int32 {
	rand.Seed(time.Now().UnixNano())
	return 1
}

//
var x = randomseedInit()

// RandomString 生成随机字符串
func RandomString(length int32) string {
	str := "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789_"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := int32(0); i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// RandomNumberStr 生成随机数字字符串
func RandomNumberStr(length int32) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := int32(0); i < length; i++ {
		x := r.Intn(len(bytes))
		if i == 0 && bytes[x] == '0' {
			i--
		} else {
			result = append(result, bytes[x])
		}
	}
	return string(result)
}

// SplitIPPort 拆分IP:PORT
func SplitIPPort(addr string) (string, uint32) {
	pos := strings.Index(addr, ":")
	port, _ := strconv.Atoi(addr[pos+1:])
	return addr[0:pos], uint32(port)
}

//
func CreateGUID() string {
	return RandomNumberStr(9)
}

//
func CreateToken() string {
	return RandomString(128)
}

//
func CreateGiftID() string {
	return RandomNumberStr(8)
}

//
func GetCurrentPath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return ""
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return ""
	}
	return string(path[0 : i+1])
}

// 区分 rune和byte
// https://blog.csdn.net/HaoDaWang/article/details/79971395
// substr ...
func substr(s string, pos, length int32) string {
	runes := []rune(s)
	l := pos + length
	if l > int32(len(runes)) {
		l = int32(len(runes))
	}
	return string(runes[pos:l])
}

//
func GetParentDirectory(dirctory string) string {
	return substr(dirctory, 0, int32(strings.LastIndex(dirctory, "/")))
}

//
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//ReadConsole 标准输入
func ReadConsole(callback func(string) int) {
	for {
		// 从标准输入读取字符串，以\n为分割
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			break
		}
		// 去掉读入内容的空白符
		text = strings.TrimSpace(text)
		rc := callback(text)
		if rc < 0 {
			break
		}
	}
}
