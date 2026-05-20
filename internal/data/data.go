package data

import (
	"my_knowledges/internal/biz"
	"my_knowledges/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewMySQL,
	NewGreeterRepo,
	NewDomainRepo,
	NewKnowledgeRepo,
	NewArticleRepo,
)

// 创建一个实现了 gormlogger.Writer 接口的包装器
type gormLogWriter struct {
	help *log.Helper
}

// Printf 实现 gormlogger.Writer 接口的 Printf 方法
func (w *gormLogWriter) Printf(format string, args ...interface{}) {
	w.help.Infof(format, args...)
}

func NewMySQL(conf *conf.Data, logger log.Logger) (*gorm.DB, func(), error) {
	l := log.NewHelper(logger)

	// 创建 GORM 日志配置，启用 SQL 查询日志
	gormLogger := gormlogger.New(
		&gormLogWriter{help: l},
		gormlogger.Config{
			SlowThreshold:             time.Second,     // 慢查询阈值
			LogLevel:                  gormlogger.Info, // 日志级别，设置为 Info 以输出 SQL 查询
			IgnoreRecordNotFoundError: true,            // 忽略记录未找到错误
			Colorful:                  true,            // 启用彩色打印
		},
	)

	db, err := gorm.Open(mysql.Open(conf.Mysql.Dsn), &gorm.Config{
		Logger: gormLogger, // 设置 GORM 日志
	})
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	sqlDB.SetMaxOpenConns(int(conf.Mysql.MaxOpenConn))
	sqlDB.SetMaxIdleConns(int(conf.Mysql.MaxIdleConn))

	if conf.Mysql.ConnMaxLifetime != "" {
		d, _ := time.ParseDuration(conf.Mysql.ConnMaxLifetime)
		sqlDB.SetConnMaxLifetime(d)
	}

	// 自动迁移表结构
	if err := db.AutoMigrate(
		&biz.Domain{},
		&biz.Knowledge{},
		&biz.Article{},
	); err != nil {
		l.Errorf("Failed to auto migrate: %v", err)
		return nil, nil, err
	}

	cleanup := func() {
		l.Info("closing mysql db")
		_ = sqlDB.Close()
	}

	return db, cleanup, nil
}

// Data .
type Data struct {
	MySQL *gorm.DB
}

// NewData .
func NewData(c *conf.Data, mysqlDb *gorm.DB) (*Data, func(), error) {
	cleanup := func() {
		log.Info("closing the data resources")
	}

	return &Data{
		MySQL: mysqlDb,
	}, cleanup, nil
}
