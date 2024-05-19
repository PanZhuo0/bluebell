package mysql

import (
	"backend/settings"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init(conf *settings.MySQLConfig) (err error) {
	// dsn := `root:123123@tcp(localhost:3306)/tset?parseTime=true,local`
	dsn := fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local`, conf.User, conf.Password, conf.Host, conf.Port, conf.DB)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("Connect MySQL failed", zap.Error(err))
		panic(err)
	}
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	zap.L().Info("Connect MySQL success")
	return
}
