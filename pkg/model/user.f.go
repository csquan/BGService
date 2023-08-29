package model

/*
  内容太多了 func 都放这里
*/

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/api-in/pkg/util"
	"github.com/ethereum/api-in/pkg/util/ecies"
	"github.com/nyaruka/phonenumbers"
	"github.com/rs/zerolog"
	"golang.org/x/exp/slices"
	"net/http"
	"net/mail"
	"strings"
	"time"
)

// Auth2 default true after user signed up
func (u *User) Auth2() bool {
	return u.BID != ""
}

// Validate the email format is correct.
func (ei *EmailInput) Validate() util.Err {
	_, err := mail.ParseAddress(ei.Email)
	if err != nil {
		return util.ErrEmailFormat(ei.Email)
	}
	return nil
}

// Validate the phone number is correct.
func (mi *MobileInput) Validate() util.Err {
	num, err := phonenumbers.Parse(mi.Mobile, "")
	if err != nil {
		return util.ErrMobileFormat(mi.Mobile)
	}
	if valid := phonenumbers.IsValidNumber(num); !valid {
		return util.ErrMobileFormat(mi.Mobile)
	}
	return nil
}

// Number must is validated
func (mi *MobileInput) Number() string {
	if mi.Validate() != nil {
		return ""
	}
	num, _ := phonenumbers.Parse(mi.Mobile, "")
	regionNumber := phonenumbers.GetRegionCodeForNumber(num)
	countryCode := phonenumbers.GetCountryCodeForRegion(regionNumber)
	nationalNumber := phonenumbers.GetNationalSignificantNumber(num)
	return fmt.Sprintf("+%d %s", countryCode, nationalNumber)
}

func (u *User) ToTU(a *AppIDAddress) *TokenUser {
	return &TokenUser{
		BID:       u.BID,
		Nick:      u.Nick,
		Mobile:    u.Mobile,
		Email:     u.Email,
		AppID:     a.AppID,
		Eth:       a.Eth,
		FirmName:  u.FirmName,
		KycTime:   u.FirmVerified,
		MAuthTime: u.LastMVTime,
		EAuthTime: u.LastEVTime,
		GAuthTime: u.LastGVTime,
		Auth2:     u.Auth2(),
		Admin:     u.Admin,
	}
}
func (u *User) ToTu(t *TokenUser) *TokenUser {
	return &TokenUser{
		BID:       u.BID,
		Nick:      u.Nick,
		Mobile:    u.Mobile,
		Email:     u.Email,
		AppID:     t.AppID,
		Eth:       t.Eth,
		FirmName:  u.FirmName,
		KycTime:   u.FirmVerified,
		MAuthTime: u.LastMVTime,
		EAuthTime: u.LastEVTime,
		GAuthTime: u.LastGVTime,
		Auth2:     u.Auth2(),
		Admin:     u.Admin,
		SK:        t.SK,
	}
}

func (tu *TokenUser) Name() string {
	if tu.Email != "" {
		return tu.Email
	}
	if tu.Mobile != "" {
		return tu.Mobile
	}
	return tu.BID
}

func (tu *TokenUser) ReSk() {
	tu.SK, _ = util.CryptoRandomString(8)
}

func passwordRule(password string) util.Err {
	e, n, u, _ := util.VerifyRule(password)
	if e && n && u {
		return nil
	} else {
		return util.ErrBrokenPasswordRule
	}
}

func (su *SignUpInput) Validate() util.Err {
	err := passwordRule(su.Password)
	if err != nil {
		return err
	}
	if su.Mobile != "" {
		mo := MobileInput{su.Mobile}
		return mo.Validate()
	}
	return nil
}

func (u *User) MaskedEmail() string {
	if u.Email == "" {
		return ""
	}
	two := strings.Split(u.Email, "@")
	if len(two) != 2 {
		return ""
	}
	if len(two[0]) <= 3 {
		return two[0] + "****@" + two[1]
	} else {
		return two[0][:3] + "****@" + two[1]
	}
}

func (u *User) MaskedMobile() string {
	if u.Mobile == "" {
		return ""
	}
	two := strings.SplitN(u.Mobile, " ", 2)
	if len(two) == 1 {
		return two[0][:5] + "****" + two[0][len(two[0])-2:]
	}
	return two[0] + " " + two[1][:2] + "****" + two[1][len(two[1])-2:]
}

func (ri *ResetPasswordInput) Validate() util.Err {
	return passwordRule(ri.Password)
}

func (ub *UnbindGAInput) Validate() util.Err {
	if ub.MCode == "" && ub.ECode == "" {
		return util.ErrInvalidArgument
	}
	return nil
}

func (be *BindEmailInput) Validate() util.Err {
	if err := be.BindGAInput.Validate(); err != nil {
		return err
	}
	return be.EmailInput.Validate()
}

func (bm *BindMobileInput) Validate() util.Err {
	if err := bm.BindGAInput.Validate(); err != nil {
		return err
	}
	return bm.MobileInput.Validate()
}

func (fc *FirmConfirmed) Validate(pk *ecies.PrivateKey, log *zerolog.Logger) util.Err {
	return fc.verify(pk, log)
}

func (v *Verify) verify(pk *ecies.PrivateKey, log *zerolog.Logger) util.Err {
	byt, err := hex.DecodeString(v.Verified)
	if err != nil {
		return util.ErrMsgDecode
	}
	dt, er := pk.Decrypt(byt)
	if er != nil {
		return er
	}
	t1, err := time.Parse(http.TimeFormat, string(dt))
	if err != nil {
		return util.ErrMsgDecrypt
	}
	span := time.Now().Sub(t1)
	if span > 3*time.Second || span < -3*time.Second {
		log.Debug().Dur("span", span).Send()
		return util.ErrMsgDecrypt
	}
	return nil
}

func (fq *FirmQuery) Number() string {
	mi := MobileInput{Mobile: fq.Mobile}
	if mi.Validate() != nil {
		return ""
	}
	return mi.Number()
}

func (fq *FirmQuery) Validate(pk *ecies.PrivateKey, log *zerolog.Logger) util.Err {
	if fq.Uid == "" && fq.Email == "" && fq.Mobile == "" {
		return util.ErrTokenInvalid
	}
	return fq.verify(pk, log)
}

func ethAddrValidate(addr string) util.Err {
	if len(addr) != 42 {
		return util.ErrAddr
	}
	if addr[:2] != "0x" {
		return util.ErrAddr
	}
	for i := 0; i < 40; i++ {
		if strings.IndexByte(HexChars, addr[i+2]) < 0 {
			return util.ErrAddr
		}
	}
	return nil
}

func (ai *AddrInput) Validate(pk *ecies.PrivateKey, log *zerolog.Logger) util.Err {
	if !slices.Contains(AllChainType, ai.Type) {
		return util.ErrChain
	}
	if err := ethAddrValidate(ai.Addr); err != nil {
		return err
	}
	return ai.verify(pk, log)
}

func (usi *KycUserStatusInput) Validate(pk *ecies.PrivateKey, log *zerolog.Logger) util.Err {
	if len(usi.Uid) != BidLen {
		return util.ErrInvalidArgument
	}
	return usi.verify(pk, log)
}

func (uli *KycUserListInput) Validate(pk *ecies.PrivateKey, log *zerolog.Logger) util.Err {
	return uli.verify(pk, log)
}

func (uli *KycUserListInput) GetLimit() int {
	if uli.Limit < 10 || uli.Limit > 150 {
		return 20
	} else {
		return uli.Limit
	}
}

func (uli *KycUserListInput) GetFrom() int {
	if uli.Page <= 0 {
		uli.Page = 1
	}
	return (uli.Page - 1) * uli.GetLimit()
}

func (uqm *IUserQueryMultiInput) Validate(pk *ecies.PrivateKey, log *zerolog.Logger) util.Err {
	if len(uqm.AppID) != 6 || len(uqm.UidList) == 0 {
		return util.ErrInvalidArgument
	}
	return nil
	// return uqm.verify(pk, log) TODO java
}
