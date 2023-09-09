package util

import (
	"encoding/json"
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/types"
	"github.com/go-xorm/xorm"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const base_tron_url = "https://api.trongrid.io"

func ModifyUserFundIn(session *xorm.Session, engine *xorm.Engine, fundInParam *types.FundInParam, userAddr *types.UserAddr) (decimal.Decimal, error) {
	//取最新余额
	url := base_tron_url + "/wallet/getaccount"

	accountParam := types.AccountParam{
		Address: userAddr.Addr,
		Visible: true,
	}

	bodyStr, err := json.Marshal(accountParam)
	if err != nil {
		logrus.Info(err)
		return decimal.Zero, err
	}

	str1, err := Post(url, bodyStr)
	if err != nil {
		logrus.Info(err)
		return decimal.Zero, err
	}
	balance := gjson.Get(str1, "balance")
	dec, err := decimal.NewFromString(balance.Str) // 目前链上余额
	if err != nil {
		logrus.Info(err)
		return decimal.Zero, err
	}

	//取出用户最近的充值记录
	userFundIn, err := db.GetUserFundIn(engine, fundInParam.Uid, fundInParam.Network)
	if err != nil {
		return decimal.Zero, err
	}

	if userFundIn == nil { //没充过值，这里就是链上余额
		fund := types.UserFundIn{
			Uid:          fundInParam.Uid,
			Network:      fundInParam.Network,
			Addr:         userAddr.Addr,
			FundInAmount: dec.String(),
		}
		_, err = session.Table("userFundIn").Insert(&fund)
		if err != nil {
			err := session.Rollback()
			if err != nil {
				return decimal.Zero, err
			}
		}
	} else {
		if userFundIn.IsCollect == true { //发生过归集 本次充值金额为 目前的链上余额-上次归集后的剩余金额
			dec1, err := decimal.NewFromString(userFundIn.CollectRemain)
			if err != nil {
				return decimal.Zero, err
			}
			dec3 := dec.Sub(dec1)
			userFundIn.FundInAmount = dec3.String()
			userFundIn.AfterFundBalance = dec.String()
		} else { //未发生归集 本次充值金额为 本次充值后链上余额-上次充值后链上余额
			dec1, err := decimal.NewFromString(userFundIn.AfterFundBalance)
			if err != nil {
				return decimal.Zero, err
			}
			dec3 := dec.Sub(dec1)
			userFundIn.FundInAmount = dec3.String()
			userFundIn.AfterFundBalance = dec.String()
		}

		userFundIn.Id = userFundIn.Id + 1

		_, err = session.Table("userFundIn").Insert(&userFundIn)
		if err != nil {
			err := session.Rollback()
			if err != nil {
				return decimal.Zero, err
			}
		}
	}
	return decimal.NewFromString(userFundIn.FundInAmount), nil
}
