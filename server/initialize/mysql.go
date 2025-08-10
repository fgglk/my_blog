package initialize

import (
	"fmt"
	"os"
	"server/global"

	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitMysql() {
	mysqlConfig := global.Config.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		mysqlConfig.Username,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.DBName,
		mysqlConfig.Config)
	// fmt.Println(dsn)
	// fmt.Println("db_name:", mysqlConfig.DBName)

	logLevel := logger.Silent
	switch mysqlConfig.LogMode {
	case "info":
		logLevel = logger.Info
	case "warn":
		logLevel = logger.Warn
	case "error":
		logLevel = logger.Error
	case "silent":
		logLevel = logger.Silent
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logLevel)})
	if err != nil {
		global.ZapLog.Error("failed to connect database", zap.Error(err))
		os.Exit(1)
	}
	sqlDB, err := db.DB()
	if err != nil {
		global.ZapLog.Error("failed to get database instance", zap.Error(err))
		os.Exit(1)
	}

	sqlDB.SetMaxIdleConns(mysqlConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(mysqlConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(mysqlConfig.MaxLifeTime)
	sqlDB.SetConnMaxIdleTime(mysqlConfig.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		global.ZapLog.Error("database connection ping failed", zap.Error(err))
		os.Exit(1)
	}

	global.DB = db
	global.ZapLog.Info("database connection successful")
}
