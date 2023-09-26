package types

import (
	"github.com/shopspring/decimal"
	"time"
)

type UserUpdate struct {
	IsBindGoogle string `xorm:"f_isBindGoogle"`
	Secret       string `xorm:"f_secret"`
}

type Users struct {
	Uid                 string    `xorm:"f_uid"`
	UserName            string    `xorm:"f_userName"`
	Password            string    `xorm:"f_password"`
	InvitationCode      string    `xorm:"f_invitationCode"`
	MailBox             string    `xorm:"f_mailBox"`
	CreateTime          time.Time `xorm:"f_createTime"`
	IsBindGoogle        bool      `xorm:"f_isBindGoogle"`
	Secret              string    `xorm:"f_secret"`
	IsIDVerify          bool      `xorm:"f_isIDVerify"`
	Mobile              string    `xorm:"f_mobile"`
	InviteNumber        int       `xorm:"f_inviteNumber"`
	ConcernCoinList     string    `xorm:"f_concernCoinList"`
	CollectStragetyList string    `xorm:"f_collectStragetyList"`
	UpdateTime          time.Time `xorm:"f_updateTime"`
}

type Invitation struct {
	Uid        string    `xorm:"f_uid"`
	SonUid     string    `xorm:"f_sonUid"`
	CreateTime time.Time `xorm:"f_createTime"`
	UpdateTime time.Time `xorm:"f_updateTime"`
	Level      string    `xorm:"f_level"`
}

type UserStrategy struct {
	Uid          string    `xorm:"f_uid"`
	StrategyID   string    `xorm:"f_strategyID"`
	JoinTime     time.Time `xorm:"f_joinTime"`
	ActualInvest string    `xorm:"f_actualInvest"`
	IsValid      bool      `xorm:"f_isValid"`
}

// 用户得策略量化收益表
type UserStrategyEarnings struct {
	Id           string    `xorm:"f_id"`
	Uid          string    `xorm:"f_uid"`
	StrategyID   string    `xorm:"f_strategyID"`
	DayBenefit   string    `xorm:"f_dayBenefit"`
	DayRatio     string    `xorm:"f_dayRatio"`
	TotalBenefit string    `xorm:"f_totalBenefit"`
	CreateTime   time.Time `xorm:"f_createTime"`
	UpdateTime   time.Time `xorm:"f_updateTime"`
}

// 用户资产表
type UserAsset struct {
	Uid        string    `xorm:"f_uid"`
	Network    string    `xorm:"f_network"`
	CoinName   string    `xorm:"f_coinName"`
	Available  string    `xorm:"f_available"`
	Lock       string    `xorm:"f_lock"`
	Total      string    `xorm:"f_total"`
	CreateTime time.Time `xorm:"f_createTime"`
	UpdateTime time.Time `xorm:"f_updateTime"`
}

// 用户私钥助记词表
type UserKey struct {
	Addr       string    `xorm:"f_addr"`
	Name       string    `xorm:"f_name"`
	PrivateKey string    `xorm:"f_privateKey"`
	CreateTime time.Time `xorm:"f_createTime"`
	UpdateTime time.Time `xorm:"f_updateTime"`
}

// 用户链上地址表
type UserAddr struct {
	Uid        string    `xorm:"f_uid"`
	Network    string    `xorm:"f_network"`
	Addr       string    `xorm:"f_addr"`
	CreateTime time.Time `xorm:"f_createTime"`
	UpdateTime time.Time `xorm:"f_updateTime"`
}

// 用户充值记录表
type UserFundIn struct {
	Id               int64     `xorm:"f_id"`
	Uid              string    `xorm:"f_uid"`
	Network          string    `xorm:"f_network"`
	Addr             string    `xorm:"f_address"`
	Coin             string    `xorm:"f_coinName"`
	FundInAmount     string    `xorm:"f_fundInAmount"`
	AfterFundBalance string    `xorm:"f_balance"`
	IsCollect        bool      `xorm:"f_isCollect"`
	CollectAmount    string    `xorm:"f_collectAmount"`
	CollectTime      time.Time `xorm:"f_collectTime"`
	CollectRemain    string    `xorm:"f_collectRemain"`
	CreateTime       time.Time `xorm:"f_createTime"`
	UpdateTime       time.Time `xorm:"f_updateTime"`
}

// 用户提币记录表
type UserFundOut struct {
	Id         int64     `xorm:"f_id"`
	FromAddr   string    `xorm:"f_from"`
	ToAddr     string    `xorm:"f_to"`
	CoinName   string    `xorm:"f_coinName"`
	Gas        string    `xorm:"f_gas"`
	Amount     string    `xorm:"f_amount"`
	CreateTime time.Time `xorm:"f_createTime"`
	UpdateTime time.Time `xorm:"f_updateTime"`
}

// 用户分佣记录表
type UserShare struct {
	Id         int64     `xorm:"f_id"`
	UId        int64     `xorm:"f_uid"`
	CoinName   string    `xorm:"f_coinName"`
	Amount     string    `xorm:"f_amount"`
	CreateTime time.Time `xorm:"f_createTime"`
	UpdateTime time.Time `xorm:"f_updateTime"`
}

// 用户体验金记录表
type UserExperience struct {
	Id             int64     `xorm:"f_id"`
	UId            string    `xorm:"f_uid"`
	CoinName       string    `xorm:"f_coinName"`
	Type           string    `xorm:"f_type"`
	ExpType        string    `xorm:"f_expType"`
	ReceiveSum     int64     `xorm:"f_receiverSum"`
	ValidTime      time.Time `xorm:"f_validTime"`
	ValidStartTime time.Time `xorm:"f_validStartTime"`
	ValidEndTime   time.Time `xorm:"f_validEndTime"`
	Status         bool      `xorm:"f_status"`
	CreateTime     time.Time `xorm:"f_createTime"`
	UpdateTime     time.Time `xorm:"f_updateTime"`
}

type AddrOutput struct {
	Uid        string    `json:"uid"`
	Network    string    `json:"network"`
	Addr       string    `json:"addr"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

type RecordOutput struct {
	Time   string `json:"time"`
	Addr   string `json:"addr"`
	Coin   string `json:"coin"`
	Type   string `json:"type"`
	Amount string `json:"amount"`
	Status string `json:"status"`
}

type RecordOutputAndGas struct {
	Time   string `json:"time"`
	Addr   string `json:"addr"`
	Coin   string `json:"coin"`
	Type   string `json:"type"`
	Amount string `json:"amount"`
	Gas    string `json:"gas"`
	Status string `json:"status"`
}

type ExpRecordOutput struct {
	Time   string `json:"time"`
	Coin   string `json:"coin"`
	Type   string `json:"type"`
	Amount string `json:"amount"`
	Valid  string `json:"valid"`
	Status string `json:"status"`
}

type UserAssetOutput struct {
	Uid       string `json:"uid"`
	Network   string `json:"network"`
	CoinName  string `json:"coinName"`
	Available string `json:"available"`
	Lock      string `json:"lock"`
	Total     string `json:"total"`
}

type FundOutParam struct {
	Uid     string
	ToAddr  string
	Amount  string
	Network string
}

type AccountIdentifier struct {
	Address string `json:"address"`
}
type BlockIdentifier struct {
	Hash   string `json:"hash"`
	Number int    `json:"number"`
}

type AccountParam struct {
	Address string `json:"address"`
	Visible bool   `json:"visible"`
}

type FundInParam struct {
	Currency        string `json:"currency"`
	Network         string `json:"Network"`
	RechargeAddress string `json:"rechargeAddress"`
}

type UserBindInfos struct {
	ID              int       `xorm:"f_id"`
	Uid             string    `xorm:"f_uid"`
	Cex             string    `xorm:"f_cex"`
	ApiKey          string    `xorm:"f_apiKey"`
	ApiSecret       string    `xorm:"f_apiSecret"`
	Passphrase      string    `xorm:"f_passphrase"`
	Alias           string    `xorm:"f_alias"`
	Account         string    `xorm:"f_account"`
	CreateTime      time.Time `xorm:"f_createTime"`
	UpdateTime      time.Time `xorm:"f_updateTime"`
	SynchronizeTime time.Time `xorm:"f_synchronizeTime"`
}

type News struct {
	ID         string    `xorm:"f_id"`
	Type       string    `xorm:"f_type"`
	Title      string    `xorm:"f_title"`
	Content    string    `xorm:"f_content"`
	Source     string    `xorm:"f_source"`
	CreateTime time.Time `xorm:"f_createTime"`
	UpdateTime time.Time `xorm:"f_updateTime"`
	Hotspot    bool      `xorm:"f_hotspot "`
	Cover      string    `xorm:"f_cover"`
}

type InsertUserBindInfo struct {
	Uid             string    `xorm:"f_uid"`
	Cex             string    `xorm:"f_cex"`
	ApiKey          string    `xorm:"f_apiKey"`
	ApiSecret       string    `xorm:"f_apiSecret"`
	Passphrase      string    `xorm:"f_passphrase"`
	Alias           string    `xorm:"f_alias"`
	Account         string    `xorm:"f_account"`
	CreateTime      time.Time `xorm:"f_createTime"`
	UpdateTime      time.Time `xorm:"f_updateTime"`
	SynchronizeTime time.Time `xorm:"f_synchronizeTime"`
}

type StrategyInput struct {
	PageSize             int    `json:"pageSize" binding:"required"`
	PageIndex            int    `json:"pageIndex" binding:"required"`
	Strategy             int    `json:"strategy"`
	Currency             string `json:"currency"`
	StrategySource       int    `json:"strategySource"`
	ProductCategory      int    `json:"productCategory"`
	RunTime              int    `json:"runTime"`
	ExpectedYield        int    `json:"expectedYield"`
	MaxWithdrawalRate    int    `json:"maxWithdrawalRate"`
	ComprehensiveSorting int    `json:"comprehensiveSorting"`
	Keywords             string `json:"keywords"`
}

type UserCodeInfos struct {
	Uid  string `json:"uid"`
	Code string `json:"code"`
}

type UserConcern struct {
	CoinPair string `json:"coinPair"`
	Method   int    `json:"method"`
}

type UserInput struct {
	UserName string `json:"username"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	//PasswordConfirm string `json:"passwordConfirm" binding:"required"`
	// Photo           string `json:"photo" binding:"required"`
	VerifyCode string `json:"verifyCode" binding:"required"`
	InviteCode string `json:"inviteCode"`
}

type ExecuteStrategyInput struct {
	ID         int  `json:"id"`
	ProductId  int  `json:"productId"`
	BreakValue int  `json:"breakValue"`
	IsBreak    bool `json:"isBreak"`
}

type UserBindInfoInput struct {
	Cex        string `json:"cex" binding:"required"`
	ApiKey     string `json:"apiKey" binding:"required"`
	ApiSecret  string `json:"secretKey" binding:"required"`
	Passphrase string `json:"passphrase" binding:"required"`
	Account    string `json:"account" binding:"required"`
	Alias      string `json:"alias"`
}

type UnbindingApiInput struct {
	Cex        string `json:"cex" binding:"required"`
	ApiKey     string `json:"apiKey" binding:"required"`
	ApiSecret  string `json:"secretKey" binding:"required"`
	Passphrase string `json:"passphrase" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 用户总收益
type UserRevenue struct {
	Id           string
	TotalBenefit float64
}

// 用户总投资
type UserInvest struct {
	Id          string  `json:"f_uid"`
	TotalInvest float64 `json:"totalInvest"`
}

type ForgotPasswordInput struct {
	Email      string `json:"email" binding:"required"`
	VerifyCode string `json:"verifyCode" binding:"required"`
	Password   string `json:"password" binding:"required,min=8"`
}

// 平台体验金信息
type PlatformExperience struct {
	TotalSum       int64 `xorm:"f_totalSum"`
	MaxPersons     int64 `xorm:"f_maxPersons"`
	ReceivePersons int64 `xorm:"f_receivePersons"`
	RecyclePersons int64 `xorm:"f_recyclePersons"`
}

// 平台体验金收益
type PlatformExperienceRevenue struct {
	Sid          string `xorm:"f_sid"`
	InvestSum    string `xorm:"f_investSum"`
	BenefitSum   string `xorm:"f_benefitSum"`
	BenefitRatio string `xorm:"f_benefitRatio"`
}

// 策略表
type Strategy struct {
	StrategyID      string    `xorm:"f_strategyID"`
	IsValid         bool      `xorm:"f_isValid"`
	RecommendRate   string    `xorm:"f_recommendRate"`
	ParticipateNum  string    `xorm:"f_participateNum"`
	TotalYield      string    `xorm:"f_totalYield"`
	TotalRevenue    string    `xorm:"f_totalRevenue"`
	StrategyName    string    `xorm:"f_strategyName"`
	Describe        string    `xorm:"f_describe"`
	Source          string    `xorm:"f_source"`
	Type            string    `xorm:"f_type"`
	CreateTime      time.Time `xorm:"f_createTime"`
	ExpectedBefenit string    `xorm:"f_expectedBefenit"`
	MaxDrawDown     string    `xorm:"f_maxDrawDown"`
	Cap             string    `xorm:"f_cap"`
	LeverageRatio   string    `xorm:"f_leverageRatio"`
	ControlLine     string    `xorm:"f_controlLine"`
	WinChance       string    `xorm:"f_winChance"`
	SharpRatio      string    `xorm:"f_sharpRatio"`
	TradableAssets  string    `xorm:"f_tradableAssets"`
	ShareRatio      string    `xorm:"f_shareRatio"`
	DividePeriod    string    `xorm:"f_dividePeriod"`
	AgreementPeriod string    `xorm:"f_agreementPeriod"`
	HostPlatform    string    `xorm:"f_hostPlatform"`
	MinInvest       string    `xorm:"f_minInvest"`
	CoinName        string    `xorm:"f_coinName"`
	UpdateTime      string    `xorm:"f_updateTime"`
}

// 交易记录表
type TransactionRecords struct {
	ID         int       `xorm:"f_id"`
	Uid        string    `xorm:"f_uid"`
	Address    string    `xorm:"f_address"`
	StrategyID string    `xorm:"f_strategyID"`
	Action     string    `xorm:"f_action"`
	Behavior   string    `xorm:"f_behavior"`
	Time       string    `xorm:"f_time"`
	CreateTime time.Time `xorm:"f_createTime"`
	UpdateTime time.Time `xorm:"f_updateTime"`
}

// 交易记录
type TransactionRecord struct {
	ID       int    `json:"id"`
	Action   string `json:"action"`
	Behavior string `json:"behavior"`
	Time     string `json:"time"`
}

// 发送交易上链接任务表--目前只考虑trx
type TransactionTask struct {
	ID          uint64    `xorm:"f_id"`
	From        string    `xorm:"f_from"`
	To          string    `xorm:"f_to"`
	Amount      string    `xorm:"f_amount"`
	SignHash    string    `xorm:"f_sign_hash"`
	TxHash      string    `xorm:"f_tx_hash"`
	State       int       `xorm:"f_state"`
	Sig         string    `xorm:"f_sig"`
	Error       string    `xorm:"f_error"`
	NetWork     string    `xorm:"f_network"`
	CreatedTime time.Time `xorm:"f_createdTime"`
	UpdatedTime time.Time `xorm:"f_updatedTime"`
}

type TradeDetails struct {
	AccountTotalAssets string `json:"accountTotalAssets"`
	InitAssets         string `json:"initAssets"`
	CurBenefit         string `json:"curBenefit"`
	WithdrawalSum      string `json:"withdrawalSum"`
	InDays             string `json:"inDays"`
	Source             string `json:"source"`
	Type               string `json:"type"`
	ShareRatio         string `json:"shareRatio"`
	DividePeriod       string `json:"dividePeriod"`
	AgreementPeriod    string `json:"agreementPeriod"`
	ProductID          string `json:"productID"`
	Name               string `json:"name"`
}

type StrategyStats struct {
	TotalBenefit string
	TotalRatio   string
	RunTime      string
}

type UserBenefits struct {
	Date    string `json:"time"`
	Benefit string `json:"price"`
	Ratio   string `json:"yield"`
}
type UserBenefitNDays struct {
	BenefitSum   decimal.Decimal `json:"statEarnings"`
	BenefitRatio string          `json:"statYield"`
	WinRatio     string          `json:"statWinRate"`
	Huiche       string          `json:"statMaxWithdrawalRate"`
	Benefitlist  []UserBenefits  `json:"list"`
}

type HttpRes struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"body"`
}

// 涨幅榜结果
type CoinStats struct {
	Symbol  string `json:"symbol"`
	Percent string `json:"percent"`
}

// PriceChangeStats define price change stats
type PriceChangeStats struct {
	Symbol             string          `json:"symbol"`
	PriceChange        string          `json:"priceChange"`
	PriceChangePercent decimal.Decimal `json:"priceChangePercent"`
	WeightedAvgPrice   string          `json:"weightedAvgPrice"`
	PrevClosePrice     string          `json:"prevClosePrice"`
	LastPrice          string          `json:"lastPrice"`
	LastQty            string          `json:"lastQty"`
	BidPrice           string          `json:"bidPrice"`
	BidQty             string          `json:"bidQty"`
	AskPrice           string          `json:"askPrice"`
	AskQty             string          `json:"askQty"`
	OpenPrice          string          `json:"openPrice"`
	HighPrice          string          `json:"highPrice"`
	LowPrice           string          `json:"lowPrice"`
	Volume             string          `json:"volume"`
	QuoteVolume        string          `json:"quoteVolume"`
	OpenTime           int64           `json:"openTime"`
	CloseTime          int64           `json:"closeTime"`
	FristID            int64           `json:"firstId"`
	LastID             int64           `json:"lastId"`
	Count              int64           `json:"count"`
}

type PriceChangeStatss []PriceChangeStats

func (s PriceChangeStatss) Len() int {
	return len(s)
}

func (s PriceChangeStatss) Less(i, j int) bool {
	return s[i].PriceChangePercent.GreaterThan(s[j].PriceChangePercent)
}

// Swap()
func (s PriceChangeStatss) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
