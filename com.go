package ucenter

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	header *http.Header
)

type UCconfig struct {
	UcUrl   string
	UCAppid string
	UcKey   string
}

const (
	useragent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.22 (KHTML, like Gecko) Chrome/25.0.1364.160 Safari/537.22"
)

func NewUCconfig(url, appid, key string) *UCconfig {
	return &UCconfig{UcUrl: url, UCAppid: appid, UcKey: key}
}

// post
func (uc *UCconfig) Post(moudle, action, arg string) string {
	arg = url.QueryEscape(uc.DiscuzEncode(fmt.Sprintf("%s&agent=%s&time=%d", arg, MD5(useragent), time.Now().Unix())))
	parseAge := fmt.Sprintf("m=%s&a=%s&inajax=2&release=20110501&input=%s&appid=%s", moudle, action, arg, uc.UCAppid)
	return uc.httpDo(parseAge)
}

// httpdo
func (uc *UCconfig) httpDo(arge string) string {
	client := &http.Client{}
	req, err := http.NewRequest("POST", uc.UcUrl+"/index.php?__times__=1", strings.NewReader(arge))
	req.Header = uc.setHeader()
	resp, err := client.Do(req)
	if err == nil {
		data, err := ioutil.ReadAll(resp.Body) //取出主体的内容
		if err != nil {
			return ""
		}
		return string(data)
	} else {
		return ""
	}
}

// 设置头文件
func (uc *UCconfig) setHeader() http.Header {
	if header == nil {
		header = &http.Header{}
		header.Set("Accept", "*/*")
		header.Set("Connection", "Close")
		header.Set("Accept-Language", "zh-cn")
		header.Set("Cache-Control", "no-cache")
		header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	header.Set("User-Agent", useragent)
	return *header
}

//加密字符串
func (uc *UCconfig) DiscuzEncode(source string) (str string) {
	return uc.DiscuzAuthcode(source, 0)
}

//解密字符串
func (uc *UCconfig) DiscuzDecode(source string) (str string) {
	return uc.DiscuzAuthcode(source, 1)
}

// 接口函数中必须包含以下函数，否则无法正确对 UCenter 传递过来的数据进行解码。
func (uc *UCconfig) DiscuzAuthcode(source string, t int) (str string) {
	ckeyLength := 4
	key := MD5(uc.UcKey)
	keya := MD5(key[:16])
	keyb := MD5(key[16:])
	var keyc string
	if t == 1 {
		keyc = source[:ckeyLength]
	} else {
		keyc = RandomString(ckeyLength)
	}
	cryptkey := keya + MD5(keya+keyc)
	key_length := len(cryptkey)
	var st []byte
	if t == 1 {
		var err error
		st, err = base64.StdEncoding.DecodeString(source[ckeyLength:])
		if err != nil {
			st, err = base64.StdEncoding.DecodeString(source[ckeyLength:] + "==")
			if err != nil {
				st, err = base64.StdEncoding.DecodeString(source[ckeyLength:] + "=")
				if err != nil {
					return ""
				}
			}
		}
	} else {
		st = []byte("0000000000" + MD5(source + keyb)[:16] + source)
	}
	st_length := len(st)
	result := bytes.Buffer{}
	box := make([]int, 256)
	for i := 0; i < 256; i++ {
		box[i] = i
	}
	rndkey := make([]int, 256)
	for i := 0; i < 256; i++ {
		l := i % key_length
		rndkey[i] = int([]byte(cryptkey[l : l+1])[0])
	}
	j := 0
	for i := 0; i < 256; i++ {
		j = (j + box[i] + rndkey[i]) % 256
		tmp := box[i]
		box[i] = box[j]
		box[j] = tmp
	}

	a := 0
	j = 0
	for i := 0; i < st_length; i++ {
		a = (a + 1) % 256
		j = (j + box[a]) % 256
		tmp := box[a]
		box[a] = box[j]
		box[j] = tmp

		fir := st[i]
		nex := byte(box[(box[a]+box[j])%256])
		result.WriteByte(fir ^ nex)
	}
	if t == 1 {
		res := result.String()
		if res[10:26] == MD5(res[26:] + keyb)[:16] {
			return res[26:]
		} else {
			return ""
		}
	} else {
		return keyc + strings.Replace(base64.StdEncoding.EncodeToString(result.Bytes()), "=", "", -1)
	}
	return
}

func RandomString(len int) (str string) {
	ch := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	sb := bytes.Buffer{}
	rand.Seed(time.Now().Unix())
	for i := 0; i < len; i++ {
		sb.WriteRune(ch[rand.Intn(len)])
	}
	return sb.String()
}

// md5
func MD5(str string) (st string) {
	m := md5.New()
	io.WriteString(m, str)
	return fmt.Sprintf("%x", m.Sum(nil))
}
