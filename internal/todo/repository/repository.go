package repository

import (
	"ice/internal/adapter/mysql"
)

type Repository struct {
	mysql *mysql.MySQL
}

func NewRepository(mysql *mysql.MySQL) *Repository {
	return &Repository{mysql: mysql}
}
