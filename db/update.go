package db

import (
	"github.com/ethereum/api-in/types"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
)

func UpdateUser(engine *xorm.Engine, uid string) error {
	var user types.Users
	_, err := engine.Table("users").Where("f_uid=?", uid).Incr("f_inviteNumber").Update(&user)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
