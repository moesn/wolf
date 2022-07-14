package db

import (
	"database/sql"
	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	MYSQL    = "mysql"
	POSTGRES = "postgres"
)

var (
	db    *gorm.DB
	sqlDB *sql.DB
	logger func(iris.Context,map[string]interface{},string)
)

func Open(dbUrl, dbType string, config *gorm.Config, maxIdleConns, maxOpenConns int, tablePrefix string,
	recorder func(iris.Context, map[string]interface{}, string), models ...interface{}) (err error) {
	if config == nil {
		config = &gorm.Config{}
	}

	if config.NamingStrategy == nil {
		config.NamingStrategy = schema.NamingStrategy{
			TablePrefix:   tablePrefix,
			SingularTable: true,
		}
	}

	switch dbType {
	case MYSQL:
		if db, err = gorm.Open(mysql.Open(dbUrl), config); err != nil {
			logrus.Errorf("打开数据库连接失败: %s", err.Error())
			return
		}
		break
	case POSTGRES:
		if db, err = gorm.Open(postgres.Open(dbUrl), config); err != nil {
			logrus.Errorf("打开数据库连接失败: %s", err.Error())
			return
		}
		break
	}

	if sqlDB, err = db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(maxIdleConns)
		sqlDB.SetMaxOpenConns(maxOpenConns)
	} else {
		logrus.Error(err)
	}

	if err = db.AutoMigrate(models...); nil != err {
		logrus.Errorf("自动合并表失败: %s", err.Error())
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
		logrus.Errorf("关闭数据库连接失败: %s", err.Error())
	}
}
