package mysql

import (
	"log"
	"server/platform/util"
)

// Pool ...
type Pool struct {
	conf *Config
	pool *util.FreeValues
}

// NewPool ...
func NewPool() *Pool {
	s := &Pool{pool: util.NewFreeValues()}
	s.pool.SetNew(s.newMysql)
	return s
}

// pool ...
var pool *Pool

// var x int32

// Instance ...
func Instance() *Pool {
	// if atomic.CompareAndSwapInt32(&x, 0, 1) == true {
	// 	pool = newPool()
	// }
	if pool == nil {
		pool = NewPool()
	}
	return pool
}

// InitConfig ...
func (s *Pool) InitConfig(conf *Config) {
	s.conf = conf
}

// ping ...
func (s *Pool) ping(i int32, value interface{}) {
	mysql := value.(*Mysql)
	if err := mysql.Ping(); err != nil {
		errno, errmsg := mysql.GetError(err)
		log.Printf("<%d>:%s", errno, errmsg)
		mysql.Connect()
	} else {
		log.Printf("ping[%d][%d] %s OK", util.GoroutineID(), i, mysql.Dsn())
	}
}

// Ping ...
func (s *Pool) Ping() {
	s.pool.Visit(s.ping)
}

// Alloc ...
func (s *Pool) Alloc() *Mysql {
	return s.allocMysql()
}

// Free ...
func (s *Pool) Free(mysql *Mysql) {
	s.freeMysql(mysql)
}

// newMysql ...
func (s *Pool) newMysql() interface{} {
	return newMysql(s.conf)
}

// allocMysql ...
func (s *Pool) allocMysql() *Mysql {
	return s.pool.Alloc().(*Mysql)
}

// freeMysql ...
func (s *Pool) freeMysql(mysql *Mysql) {
	s.pool.Free(mysql)
}

// reset ...
func (s *Pool) reset(value interface{}) {
	value.(*Mysql).Disconnect()
}

// Reset ...
func (s *Pool) Reset() {
	s.pool.ResetValues(s.reset)
}
