package ucenter

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type UserAction struct {
}

// 用户信息
type UserInfo struct {
	Data []string `xml:"item"`
}

//userName:用户名
//password:密码
//email:邮箱
//返回值
//如果err为Nil那么userId为用户ID
//否则err将返回错误信息
func (us *UserAction) Register(username, password, email string) (userId int, err error) {
	//大于 0:返回用户 ID，表示用户注册成功
	//-1:用户名不合法
	//-2:包含不允许注册的词语
	//-3:用户名已经存在
	//-4:Email 格式有误
	//-5:Email 不允许注册
	//-6:该 Email 已经被注册
	input := fmt.Sprintf("username=%v&password=%v&email=%v&questionid=&answer=&regip=", username, password, email)
	data := UCConfiger.Post("user", "register", input)
	i, err := strconv.Atoi(data)
	if err == nil {
		switch i {
		case -1:
			{
				err = errors.New("用户名不合法")
			}
		case -2:
			{
				err = errors.New("包含不允许注册的词语")
			}
		case -3:
			{
				err = errors.New("用户名已经存在")
			}
		case -4:
			{
				err = errors.New("Email 格式有误")
			}
		case -5:
			{
				err = errors.New("Email 不允许注册")
			}
		case -6:
			{
				err = errors.New("该 Email 已经被注册")
			}
		default:
			{
				userId = i
			}
		}
	}
	return
}

//登陆方法
//username:用户名 / 用户 ID / 用户 E-mail
//password:密码
//isuid:
// 是否使用用户 ID登录
//1:使用用户 ID登录
//2:使用用户 E-mail登录
//0:使用用户名登录
//checkques:
//是否验证安装提问
//1:验证安全提问
//0:不验证安全提问
//questionid:提问索引
//answer:回答内容
//返回值
//re:是否重复
//info[0]:用户名
//info[1]:密码
//info[2]:Email
func (us *UserAction) Login(username, password string, isuid int, checkques int, questionid int, answer string) (userId int, info []string, re bool, err error) {
	input := fmt.Sprintf("username=%v&password=%v&isuid=%v&checkques=%v&questionid=&answer=%v", username, password, isuid, checkques, answer)
	data := UCConfiger.Post("user", "login", input)
	dst := UserInfo{}
	d := xml.NewDecoder(strings.NewReader(data))
	d.CharsetReader = CharsetReader
	err = d.Decode(&dst)
	if err == nil {
		var i int
		i, err = strconv.Atoi(dst.Data[0])
		switch {
		case i > 0:
			{
				userId = i
				info = dst.Data[1:4]
				re, err = strconv.ParseBool(dst.Data[4])
			}
		case i == -1:
			{
				err = errors.New("用户不存在")
			}
		case i == -2:
			{
				err = errors.New("密码错")
			}
		case i == -3:
			{
				err = errors.New("安全提问错")
			}
		}
	}
	return
}

//获取用户信息
//username:用户名
//isuid:1使用ID获取,0使用用户名获取
//返回值
//userId:用户ID
//st[0]:用户名
//st[1]:邮箱
func (us *UserAction) Get(username string, isuid int) (userId int, st []string, err error) {
	arg := fmt.Sprintf("username=%s&isuid=%v", username, isuid)
	data := UCConfiger.Post("user", "get_user", arg)
	data = strings.Replace(data, "ISO-8859-1", "utf-8", -1) //有中文
	dst := UserInfo{}
	d := xml.NewDecoder(strings.NewReader(data))
	err = d.Decode(&dst)
	if err == nil {
		userId, err = strconv.Atoi(dst.Data[0])
		st = dst.Data[1:]
	}
	return
}

//修改用户信息
//username:用户名
//oldpw:旧密码
//newpw:新密码,不修改为空
//email:邮箱,不修改留空
//ignoreoldpw:1忽略旧密码,0不忽略旧密码
//questionid:问题索引
//answer:回答
//返回值
//成功或者资料无变更为nil否则返回错误信息
func (us *UserAction) Edit(username, oldpw, newpw, email string, ignoreoldpw bool, questionid int, answer string) (err error) {
	return
}

//TODO:delete对应uc_user_delete
//TODO:deleteavatar对应uc_user_deleteavatar

//同步登陆
//返回值
//err为nil则html表示的同步登陆代码
func (us *UserAction) SynLogin(userId int) (html string) {
	arg := fmt.Sprintf("uid=%v", userId)
	data := UCConfiger.Post("user", "synlogin", arg)
	return data
}

//同步退出
//返回值
//err为nil则html表示同步退出的代码
func (us *UserAction) SynLogout() (html string) {
	data := UCConfiger.Post("user", "synlogout", "")
	return data
}

//TODO:checkemail对应uc_user_checkemail
//TODO:checkname对应uc_user_checkname
//TODO:addprotected对应uc_user_addprotected
//TODO:deleteprotected对应uc_user_deleteprotected
//TODO:getprotected对应uc_user_getprotected
//TODO:merge对应uc_user_merge
//TODO:merge_remove对应uc_user_merge_remove
