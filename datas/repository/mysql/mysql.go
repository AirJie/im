package mysql

import (
	"database/sql"
	"github.com/aijie/michat/config"
	"github.com/aijie/michat/server/logger"
	_ "github.com/go-sql-driver/mysql"
)
var DBCli *sql.DB

func Init()  {
	var err error
	DBCli, err = sql.Open("mysql", config.LogicConf.MySQL)
	if err !=  nil {
		logger.Sugar.Error(err)
		panic(err)
		return
	}
}
