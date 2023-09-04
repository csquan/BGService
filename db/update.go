package db

import (
	"github.com/ethereum/BGService/types"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
)

func AddUserInvite(engine *xorm.Engine, uid string) error {
	var user types.Users
	_, err := engine.Table("users").Where("f_uid=?", uid).Incr("`f_inviteNumber`").Update(&user)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func UpdateUser(engine *xorm.Engine, uid string, user *types.Users) error {
	_, err := engine.Table("users").Where("f_uid=?", uid).Update(user)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
