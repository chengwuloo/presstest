package mysql

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"server/platform/util"
	"strconv"
	/* mysql 驱动初始化 */
	_ "github.com/go-sql-driver/mysql"
)

// SHOW VARIABLES LIKE '%max_allowed_packet%';
// SET GLOBAL max_allowed_packet = 10*1024*1024;

// GRANT ALL PRIVILEGES ON *.* TO root@"%" IDENTIFIED BY '123456'  WITH GRANT OPTION;
// FLUSH PRIVILEGES;

// Mysql ...
type Mysql struct {
	conf  *Config
	mysql *sql.DB
}

// newMysql ...
func newMysql(conf *Config) *Mysql {
	return &Mysql{conf: conf}
}

// Dsn ...
func (s *Mysql) Dsn() string {
	if s.conf != nil {
		// user:passwd@tcp(host:port)/db?charset=utf8mb4
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", s.conf.user, s.conf.passwd, s.conf.host, s.conf.port, s.conf.db, s.conf.charset)
	}
	return ""
}

// Connect ...
func (s *Mysql) Connect() error {
	s.Disconnect()
	mysql, err := sql.Open("mysql", s.Dsn())
	if err == nil {
		//mysql.SetConnMaxLifetime(3600 * 24 * 365 * time.Second)
		//mysql.SetMaxIdleConns(10)
		//mysql.SetMaxOpenConns(10)
		s.mysql = mysql
	} else {
		log.Printf("Mysql error: %v", err.Error())
	}
	return err
}

// TryConnect ...
func (s *Mysql) TryConnect() error {
	var err error
	limit := 5
	for {
		limit--
		if s.mysql == nil && limit > 0 {
			err = s.Connect()
		} else {
			break
		}
	}
	return err
}

// Disconnect ...
func (s *Mysql) Disconnect() {
	if s.mysql != nil {
		s.mysql.Close()
		s.mysql = nil
	}
}

// Ping ...
func (s *Mysql) Ping() error {
	return s.mysql.Ping()
}

// SetAutoCommit ...
func (s *Mysql) SetAutoCommit(auto bool) {
	if auto == true {
		s.exec("SET AUTOCOMMIT = 1")
	} else {
		s.exec("SET AUTOCOMMIT = 0")
	}
}

// AutoCommit ...
func (s *Mysql) AutoCommit() bool {
	var auto int
	row := s.queryRow("SELECT @@AUTOCOMMIT")
	row.Scan(&auto)
	return auto == 1
}

// Commit ...
func (s *Mysql) Commit() {
	s.exec("COMMIT")
}

// Rollback ...
func (s *Mysql) Rollback() {
	s.exec("ROLLBACK")
}

//
// DECLARE CONTINUE HANDLER FOR SQLEXCEPTION SET error = 1;
// START TRANSACTION;
// ...
// IF error = 1 THEN ROLLBACK; ELSE COMMIT; END IF;
//
// DECLARE EXIT HANDLER FOR SQLEXCEPTION ROLLBACK;
// START TRANSACTION;
// ...
// COMMIT;
//

// StartTransaction ...
func (s *Mysql) StartTransaction() {
	s.exec("START TRANSACTION")
}

// IsConnError ...
func (s *Mysql) IsConnError(err error) bool {
	errno, _ := s.GetError(err)
	return 2003 == errno || 2006 == errno || 2013 == errno
}

// GetError ...
func (s *Mysql) GetError(err error) (int32, string) {
	if err != nil {
		b := util.Str2Byte(err.Error())
		i := bytes.IndexByte(b, ':')
		if -1 == i {
			return 2006, err.Error() // "bad connection"
		}
		c := b[:i]
		d := b[i+1:]
		j := bytes.IndexByte(c, ' ')
		e := c[j+1:]
		errno, _ := strconv.Atoi(util.Byte2Str(e))
		return int32(errno), util.Byte2Str(d)
	}
	return 0, ""
}

//
// row, err := db.Query()
// row.Close()
//

// query ...
func (s *Mysql) query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.mysql.Query(query, args...)
}

// querySQL ...
func (s *Mysql) querySQL(query string, args ...interface{}) (*sql.Rows, error) {
	flag := false
	rows, err := s.query(query, args...)
	if err == nil {
		if flag && s.AutoCommit() == false {
			s.Commit()
		}
	} else {
		errno, _ := s.GetError(err)
		if errno >= 2000 && errno <= 2018 {
			// fatal error, disconnect
			s.Disconnect()
			// error: gone away
			if s.IsConnError(err) && s.TryConnect() == nil {
				rows, err = s.querySQL(query, args...)
			}
		} else /*if s.mysql != nil*/ {
			// error
			if flag && s.AutoCommit() == false {
				s.Rollback()
			}
		}
	}
	return rows, err
}

// exec ...
func (s *Mysql) exec(query string, args ...interface{}) (sql.Result, error) {
	return s.mysql.Exec(query, args...)
}

// execSQL ...
func (s *Mysql) execSQL(query string, args ...interface{}) (sql.Result, error) {
	flag := true
	result, err := s.exec(query, args...)
	if err == nil {
		if flag && s.AutoCommit() == false {
			s.Commit()
		}
	} else {
		errno, _ := s.GetError(err)
		if errno >= 2000 && errno <= 2018 {
			// fatal error, disconnect
			s.Disconnect()
			// error: gone away
			if s.IsConnError(err) && s.TryConnect() == nil {
				result, err = s.execSQL(query, args...)
			}
		} else /*if s.mysql != nil*/ {
			// error
			if flag && s.AutoCommit() == false {
				s.Rollback()
			}
		}
	}
	return result, err
}

// queryRow ...
func (s *Mysql) queryRow(query string, args ...interface{}) *sql.Row {
	return s.mysql.QueryRow(query, args...)
}

// prepare ...
func (s *Mysql) prepare(query string) (*sql.Stmt, error) {
	return s.mysql.Prepare(query)
}

// Query ...
func (s *Mysql) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if err := s.TryConnect(); err != nil {
		return nil, err
	}
	return s.querySQL(query, args...)
}

// Exec ...
func (s *Mysql) Exec(query string, args ...interface{}) (sql.Result, error) {
	if err := s.TryConnect(); err != nil {
		return nil, err
	}
	return s.execSQL(query, args...)
}

// QueryRow ...
func (s *Mysql) QueryRow(query string, args ...interface{}) *sql.Row {
	if err := s.TryConnect(); err != nil {
		return nil
	}
	return s.queryRow(query, args...)
}

// Prepare ...
func (s *Mysql) Prepare(query string) (*sql.Stmt, error) {
	if err := s.TryConnect(); err != nil {
		return nil, err
	}
	return s.prepare(query)
}
