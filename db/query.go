package db

import (
	"github.com/ethereum/api-in/types"
	"github.com/go-xorm/xorm"
)

// 查询
/*
func selectAll(engine *xorm.Engine, name string) {
	/*var user []Users
	err := engine.Where("users.username=?", name).Find(&user)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(user)
}
*/

func QuerySecret(engine *xorm.Engine, uid string) *types.Users {
	var user types.Users
	engine.Where("users.uid=?", uid).Get(&user)
	return &user
}
