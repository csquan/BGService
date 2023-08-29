package model

/*
  struct 太多了，这里只放输出
*/

type GenGaResponse struct {
	Image string `json:"image" example:"base64"`
	Text  string `json:"text" example:"please write down this code"`
}

// TokenUser 放到 jwt 中的用户结构体
type TokenUser struct {
	BID       string `json:"bid"`
	Nick      string `json:"nick"`
	Mobile    string `json:"mobile"`
	Email     string `json:"email"`
	AppID     string `json:"app_id"`
	Eth       string `json:"addr"`
	FirmName  string `json:"firm_name"`
	KycTime   int64  `json:"kyc_time"`
	MAuthTime int64  `json:"m_auth_time"`
	EAuthTime int64  `json:"e_auth_time"`
	GAuthTime int64  `json:"g_auth_time"`
	Auth2     bool   `json:"auth2"`
	SK        string `json:"sk"`
	Admin     bool   `json:"admin"`
}

// ForgetResponse user validation status 1/0
type ForgetResponse struct {
	MaskedEmail  string
	MaskedMobile string
	E            bool
	M            bool
	G            bool
}

type IDResp struct {
	BID   string
	AppID string
}

// FirmResp user info with firm
type FirmResp struct {
	Uid          string
	Nick         string
	Email        string
	Mobile       string
	FirmName     string
	Country      string
	FirmVerified int64
	Admin        bool
	Status       int8
	Created      int64
	Fid          string
	AppID        string
	Eth          string
}

// UserPage user list pagination
type UserPage struct {
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
	Rows  []FirmResp `json:"rows"`
}

type FirmUser struct {
	BID    string
	Nick   string
	Status int8
}
