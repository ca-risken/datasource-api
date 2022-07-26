package db

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	mimosasql "github.com/ca-risken/common/pkg/database/sql"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/vikyd/zero"
)

type Client struct {
	MasterDB *gorm.DB
	SlaveDB  *gorm.DB
	logger   logging.Logger
}

func NewClient(conf *Config, l logging.Logger) (*Client, error) {
	ctx := context.Background()
	m, err := connect(conf, true)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database of master: %w", err)
	}
	l.Info(ctx, "Connected to Database of master.")

	s, err := connect(conf, false)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database of slave: %w", err)
	}
	l.Info(ctx, "Connected to Database of slave.")

	return &Client{
		MasterDB: m,
		SlaveDB:  s,
		logger:   l,
	}, nil
}

type Config struct {
	MasterHost     string
	MasterUser     string
	MasterPassword string
	SlaveHost      string
	SlaveUser      string
	SlavePassword  string

	Schema        string
	Port          int
	LogMode       bool
	MaxConnection int
}

func connect(conf *Config, isMaster bool) (*gorm.DB, error) {
	var user, pass, host string
	if isMaster {
		user = conf.MasterUser
		pass = conf.MasterPassword
		host = conf.MasterHost
	} else {
		user = conf.SlaveUser
		pass = conf.SlavePassword
		host = conf.SlaveHost
	}

	dsn := fmt.Sprintf("%s:%s@tcp([%s]:%d)/%s?charset=utf8mb4&interpolateParams=true&parseTime=true&loc=Local",
		user, pass, host, conf.Port, conf.Schema)
	db, err := mimosasql.Open(dsn, conf.LogMode, conf.MaxConnection)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB. isMaster: %t, err: %+v", isMaster, err)
	}

	return db, nil
}

func convertZeroValueToNull(input interface{}) interface{} {
	if zero.IsZeroVal(input) {
		return gorm.Expr("NULL")
	}
	return input
}
