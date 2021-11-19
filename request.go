package geetest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type (
	ValidateResp struct {
		*ValidateRespFailed
		*ValidateRespSuccessed
	}

	// ValidateRespFailed
	ValidateRespFailed struct {
		Code   string     `json:"code"`
		Msg    string     `json:"msg"`
		Status string     `json:"status"`
		Desc   FailedDesc `json:"desc"`
	}

	FailedDesc struct {
		Type string `json:"type"`
	}

	ValidateRespSuccessed struct {
		Result      string      `json:"result"`
		Reason      string      `json:"reason"`
		CaptchaArgs interface{} `json:"captcha_args"`
	}
)

func buildRequest(url string, data url.Values) (resp ValidateResp, err error) {
	client := http.Client{Timeout: time.Second * 3}
	result, err := client.PostForm(url, data)

	if err != nil || result.StatusCode != http.StatusOK {
		err = errors.New("geetest service request fail")
		return
	}

	defer result.Body.Close()

	body, _ := ioutil.ReadAll(result.Body)

	resp = ValidateResp{}
	json.Unmarshal(body, &resp)
	return
}
