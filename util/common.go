package util

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/types"
	"github.com/go-xorm/xorm"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const base_tron_url = "https://api.trongrid.io"

func ModifyUserFundIn(session *xorm.Session, engine *xorm.Engine, fundInParam *types.FundInParam, userAddr *types.UserAddr, uidFormatted string) (string, error) {
	//取最新余额
	url := base_tron_url + "/wallet/getaccount"

	accountParam := types.AccountParam{
		Address: userAddr.Addr,
		Visible: true,
	}

	bodyStr, err := json.Marshal(accountParam)
	if err != nil {
		logrus.Info(err)
		return "", err
	}

	str1, err := Post(url, bodyStr)
	if err != nil {
		logrus.Info(err)
		return "", err
	}
	balance := gjson.Get(str1, "balance")

	if balance.Raw == "" {
		logrus.Info("余额为0，没有充值")
		return "", errors.New("余额为0，没有充值")
	}

	dec, err := decimal.NewFromString(balance.Raw) // 目前链上余额
	if err != nil {
		logrus.Info(err)
		return "", err
	}
	//取出用户最近的充值记录
	userFundIn, err := db.GetUserFundIn(engine, uidFormatted, fundInParam.Network)
	if err != nil {
		return "", err
	}

	if userFundIn == nil { //没充过值，这里就是链上余额
		userFundIn = &types.UserFundIn{
			Id:           0,
			Uid:          uidFormatted,
			Network:      fundInParam.Network,
			Addr:         userAddr.Addr,
			FundInAmount: dec.String(),
		}
	} else {
		if userFundIn.IsCollect == true { //发生过归集 本次充值金额为 目前的链上余额-上次归集后的剩余金额
			dec1, err := decimal.NewFromString(userFundIn.CollectRemain)
			if err != nil {
				return "", err
			}
			dec3 := dec.Sub(dec1)
			userFundIn.FundInAmount = dec3.String()
			userFundIn.AfterFundBalance = dec.String()
		} else { //未发生归集 本次充值金额为 本次充值后链上余额-上次充值后链上余额
			dec1, err := decimal.NewFromString(userFundIn.AfterFundBalance)
			if err != nil {
				return "", err
			}
			dec3 := dec.Sub(dec1)
			userFundIn.FundInAmount = dec3.String()
			userFundIn.AfterFundBalance = dec.String()
		}
		userFundIn.Id = userFundIn.Id + 1
	}
	_, err = session.Table("fundIn").Insert(userFundIn)
	if err != nil {
		err := session.Rollback()
		if err != nil {
			return "", err
		}
	}
	return userFundIn.FundInAmount, nil
}

func CheckGetKlineParam(interval string, startTime int64, endTime int64, KlineType int) error {
	if interval != "5m" && interval != "1h" && interval != "6h" {
		logrus.Info("interval not right ")
		return errors.New("")
	}
	if KlineType != 1 && KlineType != 2 && KlineType != 3 {
		logrus.Info("KlineType not right ")
		return errors.New("")
	}
	return nil
}
