# UCenter

golang版本的ucenter接口


# 简介
接口文档参照UCenter[接口函数](http://faq.comsenz.com/library/UCenter/interface/interface_user.htm)

### 功能
1. 用户注册
2. 用户登录
3. 获取用户信息
4. 同步登陆
5. 同步退出


# 使用

1. 安装discuz
2. 登录ucenter，添加应用
![](http://i.imgur.com/zC9aamv.png)
3. 输入应用名称和其他填写参数
![](http://i.imgur.com/PECNUdN.png)
4. 点击添加提交按钮，记住最后的ID
![](http://i.imgur.com/4DCJQvf.png)


### 使用代码
    	import (
			"github.com/fromiuan/ucenter"
			"fmt"	
		)
	
		const (
			url   = "http://localhost/dis/uc_server" //discuz的通讯地址(一般为http+域名+uc_server)
			appid = "3"                              //应用的id
			key   = "12345678"                       //通讯密钥
		)

    	var uc *ucenter.UC
		
		func init() {
			uc = ucenter.NewUC(url, appid, key)
		}
	
		func main() {
			userid, err := uc.Register("username", "password", "email@ihaoyue.com")
			if err != nil {
				fmt.Println(err.Error())
			}else{
				fmt.Println("user id is %d",userid)
			}
		}
		