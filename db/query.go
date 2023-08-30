package db

import (
	"fmt"
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
	engine.Where("f_uid=?", uid).Get(&user)
	return &user
}

func QueryEmail(engine *xorm.Engine, email string) *types.Users {
	var user types.Users
	engine.Where("f_mailBox=?", email).Get(user)
	fmt.Println(user)
	return &user
}

func QueryInviteCode(engine *xorm.Engine, InviteCode string) *types.Users {
	var user types.Users
	engine.Where("f_invitationCode=?", InviteCode).Get(user)
	fmt.Println(user)
	return &user
}

func GetUser(engine *xorm.Engine, uid string) *types.Users {
	var user types.Users
	_, err := engine.Where("users.uid=?", uid).Get(user)
	if err != nil {
		return nil
	}
	return &user
}

func GetUserExperience(engine *xorm.Engine, uid string) *types.UserExperience {
	var userExperience types.UserExperience
	_, err := engine.Where("uid=?", uid).Get(userExperience)
	if err != nil {
		return nil
	}
	return &userExperience
}

func GetTotalRevenue(engine *xorm.Engine) *types.TotalRevenueInfo {
	var totalRevenueInfo types.TotalRevenueInfo
	_, err := engine.Get(totalRevenueInfo)
	if err != nil {
		return nil
	}
	return &totalRevenueInfo
}
