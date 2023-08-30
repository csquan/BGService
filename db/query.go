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

func QuerySecret(engine *xorm.Engine, uid string) (error, *types.Users) {
	var user types.Users
	has, err := engine.Where("f_uid=?", uid).Get(&user)
	if err != nil {
		return err, nil
	}
	if has {
		return nil, &user // 返回指向 user 的指针
	}
	return nil, nil
}

func QueryEmail(engine *xorm.Engine, email string) (error, bool) {
	var user types.Users
	has, err := engine.Where("f_mailBox=?", email).Get(user)
	if err != nil {
		return nil, has
	}
	fmt.Println(user)
	return nil, has
}

func QueryInviteCode(engine *xorm.Engine, InviteCode string) (error, *types.Users) {
	var user types.Users
	fmt.Println(InviteCode)
	has, err := engine.Table("users").Where("f_invitationcode=?", InviteCode).Get(&user)
	if err != nil {
		return err, nil
	}
	if has {
		return nil, &user // 返回指向 user 的指针
	}
	return nil, nil
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
