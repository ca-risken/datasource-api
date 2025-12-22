package db

import (
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ca-risken/common/pkg/logging"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func newDBMock() (*Client, sqlmock.Sqlmock, error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open mock sql db, error: %+w", err)
	}
	if sqlDB == nil {
		return nil, nil, fmt.Errorf("failed to create mock db, db: %+v, mock: %+v", sqlDB, mock)
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open gorm, error: %+w", err)
	}
	return &Client{
		MasterDB: gormDB,
		SlaveDB:  gormDB,
		logger:   logging.NewLogger(),
	}, mock, nil
}
