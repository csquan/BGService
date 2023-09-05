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

func QueryInviteNumLimit(engine *xorm.Engine, InviteCode string, total int) (error, []types.Users) {
	var users []types.Users
	err := engine.Table("users").Where("`f_invitatedCode`=?", InviteCode).Limit(total).Find(&users)
	if err != nil {
		logrus.Error(err)
		return err, nil
	}
	return nil, users
}

func QueryClaimRewardNumber(engine *xorm.Engine) (error, []types.Users) {
	var users []types.Users
	err := engine.Table("users").Where("`f_claimRewardNumber` > ?", 0).Desc("`f_claimRewardNumber`").Find(&users)
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

func GetMsg(engine *xorm.Engine, pageSizeInt int, pageIndexInt int, Type string) ([]types.News, error) {
	var news []types.News
	err := engine.Table("news").Where("f_type=?", Type).Limit(pageSizeInt, (pageIndexInt-1)*pageSizeInt).Find(&news)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func GetHotMsg(engine *xorm.Engine, total int, Type string) ([]types.News, error) {
	var news []types.News
	err := engine.Table("news").Where("f_type=?", Type).Limit(total).Find(&news)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func GetMsgDetail(engine *xorm.Engine, msgId string) (*types.News, error) {
	var news types.News
	has, err := engine.Table("news").Where("f_id=?", msgId).Get(&news)
	if err != nil {
		return nil, err
	}
	if has {
		return &news, nil
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

func GetStrategy(engine *xorm.Engine, sid string) (*types.Strategy, error) {
	var strategy types.Strategy
	has, err := engine.Where("`f_strategyID`=?", sid).Get(&strategy)
	if err != nil {
		return nil, err
	}
	if has {
		return &strategy, nil
	}
	return nil, nil
}

func GetUserStrategy(engine *xorm.Engine, uid string, sid string) (*types.UserStrategy, error) {
	var userStrategy types.UserStrategy
	has, err := engine.Where("f_uid=? and `f_strategyID`=?", uid, sid).Get(&userStrategy)
	if err != nil {
		return nil, err
	}
	if has {
		return &userStrategy, nil
	}
	return nil, nil
}

func GetUserStrategys(engine *xorm.Engine, uid string) ([]*types.UserStrategy, error) {
	var userStrategys []*types.UserStrategy
	err := engine.Table("userStrategy").Where("f_uid=?", uid).Find(&userStrategys)
	if err != nil {
		return nil, err
	}
	return userStrategys, nil
}

func GetStrategyTotalAssets(engine *xorm.Engine) (float64, error) {
	var userStrategy *types.UserStrategy
	total, err := engine.Table("userStrategy").Sum(userStrategy, "f_actualInvest")
	if err != nil {
		return 0, err
	}
	return total, nil
}

func GetStrategyUserCount(engine *xorm.Engine) (int64, error) {
	total, err := engine.Table("userStrategy").GroupBy("f_uid").Count("f_uid")
	if err != nil {
		return 0, err
	}
	return total, nil
}

func GetAllStrategy(engine *xorm.Engine) ([]*types.Strategy, error) {
	var Strategy []*types.Strategy
	err := engine.Table("userStrategy").Where("f_isValid=?", "t").Find(&Strategy)
	if err != nil {
		return nil, err
	}
	return Strategy, nil
}
