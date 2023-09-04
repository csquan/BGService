package db

import (
	"github.com/ethereum/BGService/types"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
)

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

func QueryInviteNum(engine *xorm.Engine, InviteCode string) (error, []types.Users) {
	var users []types.Users
	err := engine.Table("users").Where("`f_invitatedCode`=?", InviteCode).Find(&users)
	if err != nil {
		logrus.Error(err)
		return err, nil
	}
	return nil, users
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

func GetApiKeyUserBindInfos(engine *xorm.Engine, apiKey string) (*types.UserBindInfos, error) {
	var userBindInfos types.UserBindInfos
	has, err := engine.Table("userBindInfos").Where("`f_apiKey`=?", apiKey).Get(&userBindInfos)
	if err != nil {
		return nil, err
	}
	if has {
		return &userBindInfos, nil
	}
	return nil, nil
}

func GetIdUserBindInfos(engine *xorm.Engine, uid string, apiId string) (*types.UserBindInfos, error) {
	var userBindInfos types.UserBindInfos
	has, err := engine.Table("userBindInfos").
		Where("`f_uid`=? and `f_id`=?", uid, apiId).Get(&userBindInfos)
	if err != nil {
		return nil, err
	}
	if has {
		return &userBindInfos, nil
	}
	return nil, nil
}

func GetAllUserBindInfos(engine *xorm.Engine, uid string) ([]types.UserBindInfos, error) {
	var userBindInfos []types.UserBindInfos
	err := engine.Table("userBindInfos").Where("f_uid=?", uid).Find(&userBindInfos)
	if err != nil {
		return nil, err
	}
	return userBindInfos, nil
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

func GetConcernList(engine *xorm.Engine, uid string) (tags []string) {
	// the select query, returning 1 column of array type
	url := "SELECT 'f_concernCoinList' FROM users WHERE f_uid=$1"

	ret, err := engine.Query(url, uid)
	// wrap the output parameter in pq.Array for receiving into it
	//has, err := engine.SQL(url, uid).Get(pq.Array(&foo.concernCoinList))

	if err != nil {
		logrus.Info(ret)
		logrus.Info(err)
	}
	logrus.Info(ret)
	return
}
