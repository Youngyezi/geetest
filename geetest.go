package geetest

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	version    = "golang_3.0.0"
	host       = "https://api.geetest.com"
	newcaptcha = "1"
)

type (
	//Geetest struct
	Geetest struct {
		CaptchaID  string
		PrivateKey string
	}

	//Response define response type
	Response map[string]interface{}
)

//New make a clien
func New(id string, key string) *Geetest {
	return &Geetest{
		CaptchaID:  id,
		PrivateKey: key,
	}
}

//PreProcess 判断极验服务器是否down机
func (g Geetest) PreProcess(p url.Values) Response {
	p.Add("gt", g.CaptchaID)
	p.Add("new_captcha", newcaptcha)

	url := host + "/register.php?" + p.Encode()

	challenge := GET(url)

	if challenge == "" || len(challenge) != 32 {
		return g.failbackProcess()
	}

	return g.successProcess(challenge)
}

//SuccessValidate 验证码校验
func (g *Geetest) SuccessValidate(challenge string, validate string, seccode string, data url.Values) int {

	if g.checkValidate(challenge, validate) == false {
		return 0
	}

	data.Add("seccode", seccode)
	data.Add("sdk", version)
	url := host + "/validate.php"

	codevalidate := POST(url, data)

	if codevalidate == "" {
		return 0
	}

	if MD5(seccode) == codevalidate {
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

	return 1
}

//SuccessProcess success process
func (g Geetest) successProcess(c string) Response {

	return Response{
		"success":     1,
		"gt":          g.CaptchaID,
		"challenge":   MD5(c + g.PrivateKey),
		"new_captcha": newcaptcha,
	}
}

//FailbackProcess  failback process
func (g Geetest) failbackProcess() Response {

	rnd1 := MD5(strconv.Itoa(MakeRand(100)))
	rnd2 := MD5(strconv.Itoa(MakeRand(100)))

	challenge := rnd1 + Substr(rnd2, 0, 2)

	return Response{
		"success":     0,
		"gt":          g.CaptchaID,
		"challenge":   challenge,
		"new_captcha": newcaptcha,
	}
}

func (g *Geetest) checkValidate(challenge string, validate string) bool {

	if len(validate) != 32 {
		return false
	}

	if MD5(g.PrivateKey+"geetest"+challenge) != validate {
		return false
	}

	return true
}

//GET make a http request
func GET(url string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println(err)
		return ""
	}

	res, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return ""
	}

	if err != nil {
		log.Println(err)
		return ""
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println(err)
		return ""
	}

	return string(body)
}

//POST make a http request
func POST(url string, data url.Values) string {

	res, err := http.PostForm(url, data)

	if err != nil {
		log.Println(err)
		return ""
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println(err)
		return ""
	}

	res.Body.Close()

	return string(body)
}

//MD5  hash md5
func MD5(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

//Substr  String intercept
func Substr(str string, start int, length int) string {
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

//MakeRand create a rand num
func MakeRand(num int) int {

	rand.Seed(time.Now().UnixNano())

	var mu sync.Mutex

	mu.Lock()

	v := rand.Intn(num)

	mu.Unlock()

	return v
}
