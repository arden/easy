package gorm

import (
	"github.com/gogf/gf/os/glog"
	"testing"
)

// Entity is the golang structure for table c_user.
type User struct {
	Id              uint    `orm:"id,primary"        json:"id"`                //
	Uname           string  `orm:"uname"             json:"uname"`             // 用户昵称
	Sex             int     `orm:"sex"               json:"sex"`               // 用户性别 1男 2女
	Money           float64 `orm:"money"             json:"money"`             // 用户余额
	Openid          string  `orm:"openid"            json:"openid"`            // 微信openid
	LotteryCount    int     `orm:"lottery_count"     json:"lottery_count"`     // 参与的抽奖数
	IssueCount      int     `orm:"issue_count"       json:"issue_count"`       // 发布的抽奖数
	Unionid         string  `orm:"unionid"           json:"unionid"`           // 微信unionid
	Phone           string  `orm:"phone"             json:"phone"`             // 手机号码
	Avatar          string  `orm:"avatar"            json:"avatar"`            // 用户头像
	UserType        int     `orm:"user_type"         json:"user_type"`         // 用户类型，0为普通用户
	ParentId        uint    `orm:"parent_id"         json:"parent_id"`         // 引流者UID
	USetting        string  `orm:"u_setting"         json:"u_setting"`         // 保存用户不常修改的设置信息,使用json格式保存
	Province        string  `orm:"province"          json:"province"`          //
	City            string  `orm:"city"              json:"city"`              //
	RegisterTime    int     `orm:"register_time"     json:"register_time"`     // 初次登陆时间
	LoginTime       int     `orm:"login_time"        json:"login_time"`        // 登陆时间
	UStatus         int     `orm:"u_status"          json:"u_status"`          // 用户状态，0为正常状态，1为冻结状态, 2为无效状态
	SessionKey      string  `orm:"session_key"       json:"session_key"`       // 微信授权返回的session_key
	Ip              string  `orm:"ip"                json:"ip"`                // 登陆IP
	AwardCount      int     `orm:"award_count"       json:"award_count"`       // 中奖记录数字
	LuckyCount      int     `orm:"lucky_count"       json:"lucky_count"`       // 幸运币数量
	GzhOpenid       string  `orm:"gzh_openid"        json:"gzh_openid"`        //
	ReflashDataTime int     `orm:"reflash_data_time" json:"reflash_data_time"` // 用户数据更新时间，0为未更新过，当为0时使用register_time判断是否更新用户数据
	OtherMoney      float64 `orm:"other_money"       json:"other_money"`       //
}

// TableName 会将 User 的表名重写为 `profiles`
func (User) TableName() string {
	return "c_user"
}

func TestDatabase(t *testing.T) {
	var user User
	gormDB := Database()
	db := gormDB.First(&user,23166898)
	if db.Error != nil {
		glog.Error(db.Error.Error())
		return
	}
	println(user.Uname)

	var user2 User
	gormDB2 := Database()
	db = gormDB2.First(&user2,  1120442)
	if db.Error != nil {
		glog.Error(db.Error.Error())
		return
	}
	println(user2.Uname)
}
