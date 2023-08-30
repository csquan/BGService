package db

import (
	"github.com/ethereum/api-in/types"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
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
		logrus.Error(err)
		return err, nil
	}
	if has {
		return nil, &user // 返回指向 user 的指针
	}
	return nil, nil
}

func QueryEmail(engine *xorm.Engine, email string) (error, *types.Users) {
	var user types.Users
	has, err := engine.Where("`f_mailBox`=?", email).Get(&user)
	if err != nil {
		logrus.Error(err)
		return err, nil
	}
	if has {
		return nil, &user // 返回指向 user 的指针
	}
	return nil, nil
}

func QueryInviteCode(engine *xorm.Engine, InviteCode string) (error, *types.Users) {
	var user types.Users
	has, err := engine.Table("users").Where("`f_invitationCode`=?", InviteCode).Get(&user)
	if err != nil {
		logrus.Error(err)
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
		logrus.Error(err)
		return nil
	}
	return &user
}

func GetUserExperience(engine *xorm.Engine, uid string) *types.UserExperience {
	var userExperience types.UserExperience
	_, err := engine.Where("uid=?", uid).Get(userExperience)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return &userExperience
}

func GetTotalRevenue(engine *xorm.Engine) *types.TotalRevenueInfo {
	var totalRevenueInfo types.TotalRevenueInfo
	_, err := engine.Get(totalRevenueInfo)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return &totalRevenueInfo
}
