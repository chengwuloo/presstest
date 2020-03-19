package mysql

// Access ...
type Access struct {
	mysql *Mysql
}

// NewAccess ...
func NewAccess() *Access {
	return &Access{mysql: Instance().Alloc()}
}

// GetMysql ...
func (s *Access) GetMysql() *Mysql {
	return s.mysql
}

// Reset ...
func (s *Access) Reset() {
	Instance().Free(s.mysql)
}
