package util

import "fmt"

// Common Errors forming the base of our error system
//
// Many Errors returned can be tested against these errors
// using errors.Is.
var (
	ErrInvalidArgument    = NewErr(1, "输入参数有误")
	ErrResty              = NewErr(2, "resty err")
	ErrAlreadyExist       = NewErr(3, "资源已经存在")
	ErrNotExist           = NewErr(4, "资源不存在")
	ErrUnAuthed           = NewErr(5, "未授权的访问")
	ErrBrokenPasswordRule = NewErr(6, "不满足密码规则")
	ErrPrvKeyDecode       = NewErr(9, "私钥解码错❌")
	ErrPubKeyDecode       = NewErr(9, "公钥解码错❌")
	ErrMsgDecode          = NewErr(9, "解码错❌")
	ErrMsgDecrypt         = NewErr(9, "密文解码错误")

	ErrEccInvalidMessage             = NewErr(10, "椭圆曲线: 非可解密信息")
	ErrEccGenErr                     = NewErr(10, "椭圆曲线: 密钥生成失败")
	ErrEccImport                     = NewErr(10, "椭圆曲线: 密钥倒入失败")
	ErrEccInvalidCurve               = NewErr(10, "椭圆曲线: 不一致的曲线算法")
	ErrEccIVGen                      = NewErr(10, "椭圆曲线: 随机数生成失败")
	ErrEccKeySize                    = NewErr(10, "椭圆曲线: key长度不合法")
	ErrEccInvalidPublicKey           = NewErr(10, "椭圆曲线: 不合法的公钥")
	ErrEccSharedKeyIsPointAtInfinity = NewErr(10, "椭圆曲线: 共享公钥指向了无限远")
	ErrEccSharedKeyTooBig            = NewErr(10, "椭圆曲线: 共享密钥参数过大")
	ErrEccUnsupportedECDHAlgorithm   = NewErr(10, "椭圆曲线: 不支持的曲线算法")
	ErrEccUnsupportedECIESParameters = NewErr(10, "椭圆曲线: 不支持的曲线参数")
	ErrEccInvalidKeyLen              = NewErr(10, "椭圆曲线: key过大，大于512")

	ErrDB    = NewErr(11, "db内部错误，请稍后重试或联系管理员")
	ErrNoRec = NewErr(11, "该账户未注册")
	ErrRDB   = NewErr(12, "rdb内部错误，请稍后重试或联系管理员")

	ErrOldLogin = NewErr(-1107, "登录无效，您被新登录踢出")

	ErrCryptoRand        = NewErr(13, "加密随机数生成错误")
	ErrCryptoAesCipher   = NewErr(13, "加密密钥处理错误")
	ErrCryptoAesGcm      = NewErr(13, "加密过程处理错误")
	ErrDeCryptoAesCipher = NewErr(13, "解密密钥处理错误")
	ErrDeCryptoAesGcm    = NewErr(13, "解密过程处理错误")
	ErrDeCryptoAesDec    = NewErr(13, "解密过程处理错误")
	ErrAesSize           = NewErr(13, "密文过短")

	ErrTokenGen     = NewErr(15, "签发令牌出错")
	ErrTokenDec     = NewErr(15, "令牌解析出错")
	ErrTokenInvalid = NewErr(15, "令牌非法")
	ErrMobileNo     = NewErr(7, "手机无效")
	ErrMobileFirst  = NewErr(7, "您必须先校验手机📱")
	ErrEmailNo      = NewErr(8, "邮箱📮无效")
	ErrEmailFirst   = NewErr(8, "您必须先校验邮箱")
	ErrGaFirst      = NewErr(16, "您必须先校验谷歌验证")
	ErrGaGen        = NewErr(16, "谷歌验证生成错误")
	ErrGaInvalid    = NewErr(16, "谷歌验证错误")
	ErrGaNew        = NewErr(16, "您需要重新生成谷歌验证")
	ErrBcryptHash   = NewErr(17, "加密出错")
	ErrBcryptComp   = NewErr(17, "密码错误")

	ErrWalletSvr = NewErr(18, "钱包服务器出错")
	ErrSafeSvr   = NewErr(28, "多签服务器出错")
	ErrSafeNot   = NewErr(28, "必须先全部解除多签委派")

	ErrEmailByGa     = NewErr(19, "您可以通过谷歌验证来修改邮箱")
	ErrMobileByGa    = NewErr(19, "您可以通过谷歌验证来修改手机")
	ErrEmailGaNo     = NewErr(20, "您尚未认证邮箱和谷歌验证")
	ErrMobileGaNo    = NewErr(20, "您尚未认证手机和谷歌认证")
	ErrEmailMobileNo = NewErr(20, "您没有认证的邮箱和手机")

	ErrEmailSend = NewErr(21, "邮件发送出错")

	ErrAddr    = NewErr(22, "不是一个合法链地址")
	ErrChain   = NewErr(25, "不支持的链类型")
	ErrUserBan = NewErr(23, "用户被管理员禁用")

	Err2FaExpire = NewErr(24, "二次验证时间过久")
)

func ErrKycUser(name string) Err {
	return NewErr(26, "该用户已经认证，企业名称: "+name)
}
func ErrKycUser1(nick string) Err {
	return NewErr(26, "该用户已经认证为你企业用户: "+nick)
}
func ErrKycUser2(nick string) Err {
	return NewErr(26, "该用户已经认证: "+nick)
}

func ErrKycFirm(uid string) Err {
	return NewErr(27, "该企业已被认证，用户: "+uid)
}

func ErrNickExists(n string) Err {
	return NewErr(25, "昵称被占用: "+n)
}
func ErrTokenAlg(m string) Err {
	return NewErr(15, "令牌算法不支持: "+m)
}

func ErrMobileFormat(m string) Err {
	return NewErr(7, "不是一个正确的手机号码: "+m)
}
func ErrMobileExists(m string) Err {
	return NewErr(7, "手机号已经存在："+m)
}
func ErrMobileNotEq(m1, m2 string) Err {
	return NewErr(7, "手机与已验证的不一致："+m1+"!="+m2)
}
func ErrMobileAlready(m string) Err {
	return NewErr(7, "您已经有验证过的手机: "+m)
}
func ErrMobileCode(m string) Err {
	return NewErr(7, "短信验证码错误，请确认 "+m)
}

func ErrEmailFormat(e string) Err {
	return NewErr(8, "不是一个正确的邮箱格式: "+e)
}
func ErrEmailExists(e string) Err {
	return NewErr(8, "邮箱已经存在："+e)
}
func ErrEmailNotEq(e1, e2 string) Err {
	return NewErr(8, "邮箱与已验证的不一致："+e1+"!="+e2)
}
func ErrEmailAlready(e string) Err {
	return NewErr(8, "您已经有验证过的邮箱: "+e)
}
func ErrEmailCode(e string) Err {
	return NewErr(8, "邮箱验证码错误，请确认 "+e)
}
func ErrEmailByMobile(m string) Err {
	return NewErr(19, "您可以通过认证过的手机来复位邮箱-"+m)
}
func ErrMobileByEmail(e string) Err {
	return NewErr(19, "您可以通过认证过的邮箱来复位手机-"+e)
}

type Err interface {
	Code() int
	Msg() string
	LStr() string
}

type MyErr struct {
	code int
	msg  string
}

func NewErr(code int, msg string) Err {
	return &MyErr{code: code, msg: msg}
}

func (e *MyErr) Code() int {
	return e.code
}

func (e *MyErr) Msg() string {
	return e.msg
}

func (e *MyErr) LStr() string {
	return fmt.Sprintf("{%d,%s}", e.code, e.msg)
}
