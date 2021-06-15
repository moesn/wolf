package wolf

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db    *gorm.DB
	sqlDB *sql.DB
)

// 打开PostgreSQL数据库连接
func OpenPgDB(dsn string, maxIdleConns, maxOpenConns int, config *gorm.Config) (err error) {
	if config == nil {
		config = &gorm.Config{}
	}

	// 初始化postgres数据库
	if db, err = gorm.Open(postgres.Open(dsn), config); err != nil {
		logrus.Errorf("打开数据库连接失败: %s", err.Error())
		return
	}

	// 设置最大空闲连接数和最大打开连接数
	if sqlDB, err = db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(maxIdleConns)
		sqlDB.SetMaxOpenConns(maxOpenConns)
	} else {
		logrus.Error(err)
	}

	return
}

// 打开MySql数据库连接
func OpenMySqlDB(dsn string, maxIdleConns, maxOpenConns int, config *gorm.Config) (err error) {
	return
}

// 获取数据库链接
func DB() *gorm.DB {
	return db
}

// 关闭数据库连接
func CloseDB() {
	if sqlDB == nil {
		return
	}
	if err := sqlDB.Close(); nil != err {
		logrus.Errorf("关闭数据库连接失败: %s", err.Error())
	}
}
