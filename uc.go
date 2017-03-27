package ucenter

var UCConfiger *UCconfig

type UC struct {
	UserAction
}

func NewUC(url, appid, key string) *UC {
	UCConfiger = NewUCconfig(url, appid, key)
	return &UC{}
}
