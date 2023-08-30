package db

import (
	"fmt"
	"github.com/ethereum/BGService/types"
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

func GetUser(engine *xorm.Engine, uid string) (*types.Users, error) {
	var user types.Users
	has, err := engine.Where("f_uid=?", uid).Get(&user)
	if err != nil {
		return nil, err
	}
	if has {
		return &user, nil
	}
	return nil, nil
}

func GetUserBindInfos(engine *xorm.Engine, uid string) (*types.UserBindInfos, error) {
	var userBindInfos types.UserBindInfos
	has, err := engine.Table("userBindInfos").Where("f_uid=?", uid).Get(&userBindInfos)
	if err != nil {
		return nil, err
	}
	if has {
		return &userBindInfos, nil
	}

	return nil, nil
}

func GetUserExperience(engine *xorm.Engine, uid string) (*types.UserExperience, error) {
	var userExperience types.UserExperience
	has, err := engine.Table("userExperience").Where("f_uid=?", uid).Get(&userExperience)
	if err != nil {
		return nil, err
	}
	if has {
		return &userExperience, nil
	}
	return nil, nil
}

func GetPlatformExperience(engine *xorm.Engine) (*types.PlatformExperience, error) {
	var platformExperience types.PlatformExperience
	has, err := engine.Table("platformExperience").Get(&platformExperience)
	if err != nil {
		return nil, err
	}
	if has {
		return &platformExperience, nil
	}
	return nil, nil
}
