package types

import (
	"github.com/shopspring/decimal"
	"math/big"
	"time"
)

type Users struct {
	Uid                 string    `xorm:"f_uid"`
	UserName            string    `xorm:"f_userName"`
	Password            string    `xorm:"f_password"`
	InvitationCode      string    `xorm:"f_invitationCode"`
	InvitatedCode       string    `xorm:"f_invitatedCode"`
	MailBox             string    `xorm:"f_mailBox"`
	CreateTime          time.Time `xorm:"f_createTime"`
	IsBindGoogle        bool      `xorm:"f_isBindGoogle "`
	Secret              string    `xorm:"f_secret"`
	IsIDVerify          bool      `xorm:"f_isIDVerify "`
	Mobile              string    `xorm:"f_mobile"`
	InviteNumber        int       `xorm:"f_inviteNumber"`
	ClaimRewardNumber   int       `xorm:"f_claimRewardNumber "`
	ConcernCoinList     string    `xorm:"f_concernCoinList"`
	CollectStragetyList string    `xorm:"f_collectStragetyList"`
	UpdateTime          time.Time `xorm:"f_updateTime"`
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
	Addr             string    `xorm:"f_addr"`
	FundInAmount     string    `xorm:"f_fundInAmount"`
	AfterFundBalance string    `xorm:"f_balance"`
	IsCollect        bool      `xorm:"f_isCollect"`
	CollectAmount    string    `xorm:"f_collectAmount"`
	CollectTime      time.Time `xorm:"f_collectTime"`
	CollectRemain    string    `xorm:"f_collectRemain"`
	CreateTime       time.Time `xorm:"f_createTime"`
	UpdateTime       time.Time `xorm:"f_updateTime"`
}

type FundOutParam struct {
	Uid    string
	ToAddr string
	Amount string
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
	Uid     string
	Network string
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
	Permission      bool      `xorm:"f_permission"`
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
	Cover      string    `xorm:"f_cover "`
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
	Permission      bool      `xorm:"f_permission"`
}

type StrategyInput struct {
	PageSize             string `json:"pageSize" binding:"required"`
	PageIndex            string `json:"pageIndex" binding:"required"`
	Strategy             string `json:"strategy"`
	Currency             string `json:"currency"`
	StrategySource       string `json:"strategySource"`
	ProductCategory      string `json:"productCategory"`
	RunTime              string `json:"runTime"`
	ExpectedYield        string `json:"expectedYield"`
	MaxWithdrawalRate    string `json:"maxWithdrawalRate"`
	ComprehensiveSorting string `json:"comprehensiveSorting"`
	Keywords             string `json:"keywords"`
}

type UserCodeInfos struct {
	Uid  string `json:"uid"`
	Code string `json:"code"`
}

type UserConcern struct {
	Uid      string `json:"uid"`
	CoinPair string `json:"coinPair"`
	Method   string `json:"method"`
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
	ID        string `json:"id" binding:"required"`
	ProductId string `json:"productId" binding:"required"`
	IsBreak   string `json:"isBreak" binding:"required"`
}

type UserBindInfoInput struct {
	Cex        string `json:"cex" binding:"required"`
	ApiKey     string `json:"apiKey" binding:"required"`
	ApiSecret  string `json:"secretKey" binding:"required"`
	Passphrase string `json:"passphrase" binding:"required"`
	Alias      string `json:"alias"`
	Account    string `json:"account" binding:"required"`
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

// 用户体验金
type UserExperience struct {
	Uid          string          `xorm:"f_uid"`
	ReceiveSum   int64           `xorm:"f_receiveSum"`
	BenefitSum   decimal.Decimal `xorm:"f_benefitSum"`
	BenefitRatio decimal.Decimal `xorm:"f_benefitRatio"`
	ReceiveDays  int             `xorm:"f_receiveDays"`
}

// 平台体验金信息
type PlatformExperience struct {
	TotalSum       int64 `xorm:"f_totalSum"`
	PerSum         int64 `xorm:"f_perSum"`
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
	AccountTotalAssets map[string]string `json:"accountTotalAssets"`
	InitAssets         map[string]string `json:"initAssets"`
	CurBenefit         map[string]string `json:"curBenefit"`
	WithdrawalSum      map[string]string `json:"withdrawalSum"`
	InDays             string            `json:"inDays"`
	Source             string            `json:"source"`
	Type               string            `json:"type"`
	ShareRatio         string            `json:"shareRatio"`
	DividePeriod       string            `json:"dividePeriod"`
	AgreementPeriod    string            `json:"agreementPeriod"`
	ProductID          string            `json:"productID"`
}

type StrategyStats struct {
	TotalBenefit string
	TotalRatio   string
	RunTime      string
}

type UserBenefits struct {
	Date    string
	Benefit string
}
type UserBenefit30Days struct {
	BenefitSum   decimal.Decimal
	BenefitRatio string
	WinRatio     string
	Huiche       string
}

type Balance_Erc20 struct {
	Id             string `xorm:"id"`
	Addr           string `xorm:"addr"`
	ContractAddr   string `xorm:"contract_addr"`
	Balance        string `xorm:"balance"`
	Height         string `xorm:"height"`
	Balance_Origin string `xorm:"balance_origin"`
}

type Tx struct {
	Id                   string `xorm:"id"`
	TxType               string `xorm:"tx_type"`
	From                 string `xorm:"addr_from"`
	To                   string `xorm:"addr_to"`
	Hash                 string `xorm:"tx_hash"`
	Index                string `xorm:"tx_index"`
	Value                string `xorm:"tx_value"`
	Input                string `xorm:"input"`
	Nonce                string `xorm:"nonce"`
	GasPrice             string `xorm:"gas_price"`
	GasLimit             string `xorm:"gas_limit"`
	GasUsed              string `xorm:"gas_used"`
	IsContract           string `xorm:"is_contract"`
	IsContractCreate     string `xorm:"is_contract_create"`
	BlockTime            string `xorm:"block_time"`
	BlockNum             string `xorm:"block_num"`
	BlockHash            string `xorm:"block_hash"`
	ExecStatus           string `xorm:"exec_status"`
	CreateTime           string `xorm:"create_time"`
	BlockState           string `xorm:"block_state"`
	MaxFeePerGas         string `xorm:"max_fee_per_gas"`
	BaseFee              string `xorm:"base_fee"`
	MaxPriorityFeePerGas string `xorm:"max_priority_fee_per_gas"`
	BurntFees            string `xorm:"burnt_fees"`
}

type Erc20Tx struct {
	Id                string `xorm:"id"`
	TxHash            string `xorm:"tx_hash"`
	Addr              string `xorm:"addr"`
	Sender            string `xorm:"sender"`
	Receiver          string `xorm:"receiver"`
	Tokens_Cnt        string `xorm:"token_cnt"`
	Log_Index         string `xorm:"log_index"`
	Tokens_Cnt_Origin string `xorm:"token_cnt_origin"`
	Create_Time       string `xorm:"create_time"`
	Block_State       string `xorm:"block_state"`
	Block_Num         string `xorm:"block_num"`
	Block_Time        string `xorm:"block_time"`
}

type Erc20Info struct {
	Id                   string `xorm:"id"`
	Addr                 string `xorm:"addr"`
	Name                 string `xorm:"name"`
	Symbol               string `xorm:"symbol"`
	Decimals             string `xorm:"decimals"`
	Totoal_Supply        string `xorm:"total_supply"`
	Totoal_Supply_Origin string `xorm:"total_supply_origin"`
	Create_Time          string `xorm:"create_time"`
}

type CoinInfo struct {
	BaseInfo    Erc20Info
	HolderCount int
	Status      uint8
}

type StatusInfo struct {
	IsBlack          bool
	IsBlackIn        bool
	IsBlackOut       bool
	NowFrozenAmount  *big.Int
	WaitFrozenAmount *big.Int
}

type TxRes struct {
	Hash      string
	OpParams  *OpParam
	Amount    uint64
	TxGeneral *Tx
}

//type ContractAbi struct {
//	Contract_addr string
//	Abi_data      string
//}
//
//type ContractReceiver struct {
//	Contract_Addr string `xorm:"contract_addr"`
//	Receiver_Addr string `xorm:"reciver_addr"`
//}

type HttpRes struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"body"`
}

type CoinData struct {
	InitCoinSupply string `json:"init_coin_supply"`
	AddCoinHistory string `json:"add_coin_history"`
}

type TxData struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Data      string `json:"data"`
	ChainId   string `json:"chainId"`
	Value     string `json:"value"`
	RequestID string `json:"requestId"`
	UID       string `json:"uid"`
	UUID      string `json:"uuid"`
}

type TxLog struct {
	Id         uint64 `xorm:"id"`
	Hash       string `xorm:"tx_hash"`
	Addr       string `xorm:"addr"`
	AddrFrom   string `xorm:"addr_from"`
	AddrTo     string `xorm:"addr_to"`
	Topic0     string `xorm:"topic0"`
	Topic1     string `xorm:"topic1"`
	Topic2     string `xorm:"topic2"`
	Topic3     string `xorm:"topic3"`
	Data       string `xorm:"log_data"`
	Index      uint   `xorm:"log_index"`
	BlockState uint8  `xorm:"block_state"`
	BlockNum   uint64 `xorm:"block_num"`
	BlockTime  uint64 `xorm:"block_time"`
}

type EventHash struct {
	Op        string `xorm:"op"`
	EventHash string `xorm:"eventhash"`
}

type OpParam struct {
	Op     string `json:"op"`
	Addr1  string `json:"addr1"`
	Addr2  string `json:"addr2"`
	Value1 string `json:"value1"`
	Value2 string `json:"value2"`
	Value3 string `json:"value3"`
}

type Mechanism struct {
	Name      string `xorm:"f_name"`
	ApiKey    string `xorm:"f_key"`
	ApiSecret string `xorm:"f_secret"`
}

func (t *Mechanism) TableName() string {
	return "t_mechanism"
}

type Transfer struct {
	FromAccount string `json:"fromAccount"`
	ToAccount   string `json:"toAccount"`
	ThirdId     string `json:"thirdId"`
	Token       string `json:"token"`
	Amount      string `json:"amount"`
	CallBack    string `json:"callBack"`
	Ext         string `json:"ext"`
}

type TransferRecord struct {
	FromAccount string `xorm:"f_fromAccount"`
	ToAccount   string `xorm:"f_toAccount"`
	ThirdId     string `xorm:"f_thirdId"`
	Token       string `xorm:"f_token"`
	Amount      string `xorm:"f_amount"`
	CallBack    string `xorm:"f_callBack"`
	Ext         string `xorm:"f_ext"`
}

func (t *TransferRecord) TableName() string {
	return "t_transfer"
}

type Withdraw struct {
	Handshake
	ThirdId string `json:"thirdId"`
	Account string `json:"account"`
	Token   string `json:"token"`
	Amount  string `json:"amount"`
	Chain   string `json:"chain"`
	Addr    string `json:"addr"`
	IsSync  bool   `json:"isSync"`
}

type Handshake struct {
	ApiKey   string `json:"apiKey"`
	Nonce    string `json:"nonce"`
	Verified string `json:"verified"`
}

type InternalTransfer struct {
	Handshake
	Transfers     []Transfer `json:"transfers"`
	IsSync        bool       `json:"IsSync"`
	IsTransaction bool       `json:"isTransaction"`
}
