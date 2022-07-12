package db

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type GormModel struct {
	Id int64 `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
}

var (
	db    *gorm.DB
	sqlDB *sql.DB
	logger func(iris.Context,map[string]interface{},string)
)

func Open(dsn string, config *gorm.Config, maxIdleConns, maxOpenConns int,tablePrefix string,
	recorder func(iris.Context,map[string]interface{},string), models ...interface{}) (err error) {
	if config == nil {
		config = &gorm.Config{}
	}

	if config.NamingStrategy == nil {
		config.NamingStrategy = schema.NamingStrategy{
			TablePrefix:   tablePrefix,
			SingularTable: true,
		}
	}

	if db, err = gorm.Open(mysql.Open(dsn), config); err != nil {
		logrus.Errorf("opens database failed: %s", err.Error())
		return
	}

	if sqlDB, err = db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(maxIdleConns)
		sqlDB.SetMaxOpenConns(maxOpenConns)
	} else {
		logrus.Error(err)
	}

	if err = db.AutoMigrate(models...); nil != err {
		logrus.Errorf("auto migrate tables failed: %s", err.Error())
	}

	logger=recorder
	return
}

func DB() *gorm.DB {
	return db
}

func Close() {
	if sqlDB == nil {
		return
	}
	if err := sqlDB.Close(); nil != err {
		logrus.Errorf("Disconnect from database failed: %s", err.Error())
	}
}
