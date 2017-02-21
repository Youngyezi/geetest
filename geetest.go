package geetest

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//GtSdkVersion GEETEST SDK Version
const GtSdkVersion = "go_1.0.0"

//Geetest q
type Geetest struct {
	CaptchaID  string
	PrivateKey string
	Res        *rep
}

type rep map[string]string

//PreProcess 判断极验服务器是否down机
func (g *Geetest) PreProcess(uid string) rep {
	url := "https://api.geetest.com/register.php?gt=" + g.CaptchaID +
		"&user_id=" + uid
	challenge, err := g.GET(url)

	if err != nil {
		return g.failbackProcess()
	}

	if len(challenge) != 32 {
		return g.failbackProcess()
	}

	return g.successProcess(challenge)
}

//SuccessValidate 验证码校验
func (g *Geetest) SuccessValidate(challenge string, validate string, seccode string) int {

	if g.checkValidate(challenge, validate) == false {
		return 0
	}

	data := url.Values{}
	data.Add("seccode", seccode)
	data.Add("sdk", GtSdkVersion)
	url := "https://api.geetest.com/validate.php"
	codevalidate, err := g.POST(url, data)

	if err != nil {
		return 0
	}

	if g.md5(seccode) == codevalidate {
		return 1
	}

	return 0
}

//FailValidate 验证码校验
func (g *Geetest) FailValidate(challenge string, validate string) int {
	v := strings.Split(validate, "_")

	if len(v) != 3 {
		return 0
	}
	ans := g.decode(challenge, v[0])
	bgIdx := g.decode(challenge, v[1])
	grpIdx := g.decode(challenge, v[2])
	xPos := g.getFailbackPicAns(bgIdx, grpIdx)

	answer := ans - xPos

	if answer < 0 {
		answer = -answer
	}

	if answer < 4 {
		return 1
	}
	return 0

}

func (g *Geetest) checkValidate(challenge string, validate string) bool {

	if len(validate) != 32 {
		return false
	}

	if g.md5(g.PrivateKey+"geetest"+challenge) != validate {
		return false
	}

	return true
}

func (g *Geetest) failbackProcess() rep {
	r1 := g.md5(strconv.Itoa(int(RandInt64(0, 100))))
	r2 := g.md5(strconv.Itoa(int(RandInt64(0, 100))))

	challenge := r1 + g.Substr(r2, 0, 2)

	res := map[string]string{
		"success":   "0",
		"gt":        g.CaptchaID,
		"challenge": challenge,
	}

	return res
}

func (g *Geetest) successProcess(challenge string) rep {

	challenge = g.md5(challenge)

	return map[string]string{
		"success":   "1",
		"gt":        g.CaptchaID,
		"challenge": challenge,
	}
}

func (g *Geetest) md5(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

//GET make a http request
func (g *Geetest) GET(url string) (rep string, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	res, _ := client.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	return string(body), err
}

//POST make a http request
func (g *Geetest) POST(url string, data url.Values) (rep string, err error) {
	res, err := http.PostForm(url, data)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	res.Body.Close()

	return string(body), nil
}

func (g *Geetest) decode(challenge string, str string) int {
	if len(str) > 100 {
		return 0
	}

	key := map[string]int{}
	aryChallenge := g.splitStr(challenge)
	aryValue := g.splitStr(str)
	num := []int{1, 2, 5, 10, 50}
	count := 0
	tmp := map[string]int{}
	for i := 0; i < len(challenge); i++ {

		item := aryChallenge[i]

		if tmp[item] == 1 {
			continue
		}

		key[item] = num[count%5]
		count++
		tmp[item] = 1
	}

	res := 0
	for j := 0; j < len(str); j++ {
		res += key[aryValue[j]]
	}

	res -= g.decodeRandBase(challenge)
	return res
}

func (g *Geetest) splitStr(str string) []string {
	ary := []string{}
	l := len(str)
	for i := 0; i < l; i++ {

		ary = append(ary, string(str[i]))
	}

	return ary
}

func (g *Geetest) decodeRandBase(challenge string) int {
	base := g.Substr(challenge, 32, 2)

	var res int
	tmp := []int{}
	for i := 0; i < len(base); i++ {
		as := int(base[i])
		if as > 57 {
			res = as - 87
		} else {
			res = as - 48
		}
		tmp = append(tmp, res)
	}

	return tmp[0]*36 + tmp[1]
}

func (g *Geetest) getFailbackPicAns(f int, i int) int {

	fullBgName := g.Substr(g.md5(fmt.Sprintf("%d", f)), 0, 9)
	bgName := g.Substr(g.md5(fmt.Sprintf("%d", i)), 10, 9)

	answer := ""
	for i := 0; i < 9; i++ {
		if i%2 == 0 {
			answer = answer + string(fullBgName[i])
		} else {
			answer = answer + string(bgName[i])
		}
	}

	xDecode := g.Substr(answer, 4, 5)
	xPos := g.getXPosFromStr(xDecode)

	return xPos
}

func (g *Geetest) getXPosFromStr(str string) int {
	if len(str) != 5 {
		return 0
	}

	sum, _ := strconv.ParseInt(str, 16, 32)

	res := sum % 200

	if res < 40 {
		res = 40
	}

	return int(res)
}

//Substr 字符串截取
func (g *Geetest) Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

//RandInt64 range min~max 随机数
func RandInt64(min, max int64) int64 {
	if min >= max || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}
