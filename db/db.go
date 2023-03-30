package db

import (
	"database/sql"
	"github.com/kataras/iris/v12"
	cnd "github.com/moesn/wolf/sql"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	MYSQL    = "mysql"
	POSTGRES = "postgres"
	SQLITE   = "sqlite"
)

var (
	db     *gorm.DB
	sqlDB  *sql.DB
	logger func(iris.Context, map[string]interface{}, string)
)

func Open(dsn, dbType string, tablePrefix string, maxIdleConns, maxOpenConns int, config *gorm.Config,
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
		if db, err = gorm.Open(mysql.Open(dsn), config); err != nil {
			logrus.Errorf("打开数据库连接失败: %s", err.Error())
			return
		}
		break
	case POSTGRES:
		if db, err = gorm.Open(postgres.Open(dsn), config); err != nil {
			logrus.Errorf("打开数据库连接失败: %s", err.Error())
			return
		}
		break
	case SQLITE:
		if db, err = gorm.Open(sqlite.Open(dsn), config); err != nil {
			logrus.Errorf("打开数据库连接失败: %s", err.Error())
			return
		}
		cnd.NoConcat = true
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

	logger = recorder
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
