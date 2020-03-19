package mysql

// Config ...
type Config struct {
	host    string /* "127.0.0.1" */
	user    string /* "root" */
	passwd  string /* "admin@2018" */
	db      string /* "xxoodb" */
	charset string /* "utf8" "utf8mb4" "gbk" */
	port    uint16 /* 3306 */
}

// NewConfig ...
func NewConfig(host, user, passwd, db, charset string, port uint16) *Config {
	return &Config{host: host, user: user, passwd: passwd, db: db, charset: charset, port: port}
}
