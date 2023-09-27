package db

import (
	"fmt"
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

func UpdateUserGoogle(engine *xorm.Engine, uid string, user *types.UserUpdate) error {
	_, err := engine.Table("users").Where("f_uid=?", uid).Update(user)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func UpdateAddCollectProduct(engine *xorm.Engine, productId string, uid string) error {
	//var user types.Users
	updateSql := fmt.Sprintf("update users set `f_collectStragetyList`=array_append(`f_collectStragetyList`, '%s') where `f_uid`='%s'", productId, uid)
	_, err := engine.Exec(updateSql)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func UpdateDelCollectProduct(engine *xorm.Engine, productId string, uid string) error {
	updateSql := fmt.Sprintf("update users set `f_collectStragetyList`='%s' where `f_uid`='%s'", productId, uid)
	_, err := engine.Exec(updateSql)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func UpdatePlatformExp(engine *xorm.Engine, platExp *types.PlatformExperience) error {
	_, err := engine.Table("platformExperience").Update(platExp)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
