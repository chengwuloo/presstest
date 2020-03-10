package util

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import (
	"log"
	"os"
	"reflect"
)

//
func Log(module string, class string, funcname string, v interface{}) {
	if v == nil {
		return
	}
	//log.Println("--- *** ", module, " - ", class, ":: ", funcname, "\n", reflect.TypeOf(v).Elem(), "\n", JSON2Str(v))
	//log.Printf("--- *** PID[%07d] %v - %v :: %v \n %v \n %v\n", os.Getpid(), module, class, funcname, reflect.TypeOf(v).Elem(), JSON2Str(v))
}

//
func Logx(module string, class string, userID, Account int64, funcname string, v interface{}) {
	if v == nil {
		log.Fatalf("Logx")
		return
	}
	//log.Println("--- *** ", module, " - ", class, ":: ", funcname, "\n", reflect.TypeOf(v).Elem(), "\n", JSON2Str(v))
	//log.Printf("--- *** PID[%07d][%d:%d] %v - %v :: %v \n %v \n %v\n", os.Getpid(), userID, Account, module, class, funcname, reflect.TypeOf(v).Elem(), JSON2Str(v))
}

//
func Logy(module string, class string, funcname string, v interface{}) {
	if v == nil {
		return
	}
	//log.Println("--- *** ", module, " - ", class, ":: ", funcname, "\n", reflect.TypeOf(v).Elem(), "\n", JSON2Str(v))
	//log.Printf("--- *** PID[%07d] %v - %v :: %v \n %v \n %v\n\n", os.Getpid(), module, class, funcname, reflect.TypeOf(v).Elem(), JSON2Str(v))
}

//
func Logz(module string, class string, userID, Account int64, funcname string, v interface{}) {
	if v == nil {
		log.Fatalf("Logz")
		return
	}
	//log.Println("--- *** ", module, " - ", class, ":: ", funcname, "\n", reflect.TypeOf(v).Elem(), "\n", JSON2Str(v))
	log.Printf("--- *** PID[%07d][%d:%d] %v - %v :: %v \n %v \n %v\n", os.Getpid(), userID, Account, module, class, funcname, reflect.TypeOf(v).Elem(), JSON2Str(v))
}
