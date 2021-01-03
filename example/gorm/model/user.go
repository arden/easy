package model

// Entity is the golang structure for table c_user.
type User struct {
	Id              uint    `gorm:"id,primary"        json:"id"`                //
	Uname           string  `gorm:"uname"             json:"uname"`             // 用户昵称
	Sex             int     `gorm:"sex"               json:"sex"`               // 用户性别 1男 2女
	Money           float64 `gorm:"money"             json:"money"`             // 用户余额
	Openid          string  `gorm:"openid"            json:"openid"`            // 微信openid
	LotteryCount    int     `gorm:"lottery_count"     json:"lottery_count"`     // 参与的抽奖数
	IssueCount      int     `gorm:"issue_count"       json:"issue_count"`       // 发布的抽奖数
	Unionid         string  `gorm:"unionid"           json:"unionid"`           // 微信unionid
	Phone           string  `gorm:"phone"             json:"phone"`             // 手机号码
	Avatar          string  `gorm:"avatar"            json:"avatar"`            // 用户头像
	UserType        int     `gorm:"user_type"         json:"user_type"`         // 用户类型，0为普通用户
	ParentId        uint    `gorm:"parent_id"         json:"parent_id"`         // 引流者UID
	USetting        string  `gorm:"u_setting"         json:"u_setting"`         // 保存用户不常修改的设置信息,使用json格式保存
	Province        string  `gorm:"province"          json:"province"`          //
	City            string  `gorm:"city"              json:"city"`              //
	RegisterTime    int     `gorm:"register_time"     json:"register_time"`     // 初次登陆时间
	LoginTime       int     `gorm:"login_time"        json:"login_time"`        // 登陆时间
	UStatus         int     `gorm:"u_status"          json:"u_status"`          // 用户状态，0为正常状态，1为冻结状态, 2为无效状态
	SessionKey      string  `gorm:"session_key"       json:"session_key"`       // 微信授权返回的session_key
	Ip              string  `gorm:"ip"                json:"ip"`                // 登陆IP
	AwardCount      int     `gorm:"award_count"       json:"award_count"`       // 中奖记录数字
	LuckyCount      int     `gorm:"lucky_count"       json:"lucky_count"`       // 幸运币数量
	GzhOpenid       string  `gorm:"gzh_openid"        json:"gzh_openid"`        //
	ReflashDataTime int     `gorm:"reflash_data_time" json:"reflash_data_time"` // 用户数据更新时间，0为未更新过，当为0时使用register_time判断是否更新用户数据
	OtherMoney      float64 `gorm:"other_money"       json:"other_money"`       //
}


// TableName 会将 User 的表名重写为 `profiles`
func (User) TableName() string {
	return "c_user"
}

