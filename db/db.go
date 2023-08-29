package db

import (
	"fmt"
	"github.com/ethereum/api-in/config"
	"github.com/go-xorm/xorm"

	_ "github.com/lib/pq"
	"log"
)

// 连接
func GetDBEngine(con *config.Config) *xorm.Engine {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", con.Db.Ip, con.Db.Port, con.Db.Name, con.Db.Password, con.Db.Database)
	engine, err := xorm.NewEngine("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	engine.ShowSQL()
	err = engine.Ping()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fmt.Println("connect postgresql success")
	return engine
}
