package model

/*
  struct 太多了，这里只放输入参数
*/

type EmailInput struct {
	Email string `json:"email" binding:"required" example:"alice@example.com"`
}

type MobileInput struct {
	Mobile string `json:"mobile" binding:"required,max=20" example:"+8613901009988"`
}

type CodeInput struct {
	Code string `json:"code" example:"123456"`
}

type VerifyEmailInput struct {
	EmailInput
	CodeInput
}

type VerifyMobileInput struct {
	MobileInput
	CodeInput
}

type App struct {
	AppID string `json:"app-id" binding:"required" example:"888888"`
}

type SignUpInput struct {
	App
	Email    string `json:"email" example:"alice@gmail.com"`
	Mobile   string `json:"mobile" binding:"max=20" example:"+86 13901009988"`
	Password string `json:"password" binding:"required,min=8"`
	FirmName string `json:"firm-name" binding:"required"`
	FirmType uint   `json:"firm-type" binding:"required"`
	Country  string `json:"country" binding:"required"`
}

type SignInInput struct {
	App
	User     string `json:"user"  binding:"required" example:"alice@gmail.com/+8613901009988"`
	Password string `json:"password"  binding:"required" example:"********"`
}

type PresetPasswordInput struct {
	BindGAInput
	App
	User string `json:"user" binding:"required" example:"alice@gmail.com/+8613901009988"`
}

type ResetPasswordInput struct {
	Password string `json:"password" binding:"required,min=8"`
}

type ChangeMobileInput struct {
	ChangeInput
	MobileInput
}

type BindEmailInput struct {
	BindGAInput
	EmailInput
}

type UnbindGAInput struct {
	ECode string `json:"e-code,omitempty"`
	MCode string `json:"m-code,omitempty"`
}

type BindGAInput struct {
	UnbindGAInput
	GCode string `json:"g-code,omitempty"`
}

type ChangeInput struct {
	OldCode string `json:"old-code" binding:"required"`
	NewCode string `json:"new-code" binding:"required"`
}

type ChangeEmailInput struct {
	ChangeInput
	EmailInput
}

type BindMobileInput struct {
	BindGAInput
	MobileInput
}

type Verify struct {
	Verified string `json:"verified" binding:"required" example:"BASE198964"`
}

// FirmConfirmed input by admin
// FirmName + Country = uniq key
type FirmConfirmed struct {
	Verify
	BID      string `json:"bid" binding:"required" example:"1233210123"`
	FirmName string `json:"firm-name" binding:"required" example:"一地鸡毛蒜皮小公司"`
	FirmType uint   `json:"firm-type" binding:"required" example:"2"`
	Country  string `json:"country" binding:"required" example:"+86"`
}

type FirmQuery struct {
	Verify
	Uid    string `json:"uid,omitempty"`
	Email  string `json:"email,omitempty"`
	Mobile string `json:"mobile,omitempty"`
}

type AddrInput struct {
	Verify
	Addr string `json:"addr" binding:"required" example:"0x8cbf3d676bab7e93e94a9a2de153aff1e2f3124c"`
	Type string `json:"type" binding:"required" example:"eth"`
}

type KycUserStatusInput struct {
	Verify
	UidInput
	Status int8 `json:"status" example:"-1"`
}

type KycUserListInput struct {
	Verify
	Start  int64 `json:"start" binding:"required" example:"1672728135"`
	End    int64 `json:"end" binding:"required" example:"1672728135"`
	Page   int   `json:"page" binding:"required" example:"1"`
	Limit  int   `json:"limit" binding:"required" example:"20"`
	Status int8  `json:"status" example:"0"`
}

type NickInput struct {
	Nick string `json:"nick" binding:"required" example:"MadDog"`
}

type UidInput struct {
	Uid string `json:"uid" binding:"required" example:"569923450059"`
}

type PreSignInput struct {
	ECode string      `json:"e-code" binding:"required"`
	MCode string      `json:"m-code" binding:"required"`
	GCode string      `json:"g-code" binding:"required"`
	Data  interface{} `json:"data" binding:"required"`
}

type IUserQueryMultiInput struct {
	Verify
	AppID   string   `json:"app-id" binding:"required" example:"888888"`
	UidList []string `json:"uid-list" binding:"required" example:"[569923450059]"`
}
