package geetest

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
)

func New(captchaID, captchaKey string) (config CaptchaConfig) {
	config = CaptchaConfig{ID: captchaID, Key: captchaKey}
	return
}

// NewDemo demo to test
func NewDemo() (config CaptchaConfig) {
	return ConfigDefault
}

type Req map[string]string

// Validate
func (config CaptchaConfig) Validate(req Req) (err error) {
	req["sign_token"] = config.sign(req["lot_number"])

	params := formatRequestParams(req)

	resp, err := buildRequest(config.validateAPI(), params)

	if err != nil || resp.ValidateRespFailed != nil {
		err = errors.New("validation failed")
		return
	}

	if resp.ValidateRespSuccessed.Result != "success" {
		err = errors.New(resp.ValidateRespSuccessed.Reason)
		return
	}

	return
}

func formatRequestParams(params Req) (p url.Values) {
	p = url.Values{}
	for key, val := range params {
		p.Add(key, val)
	}
	return
}

func (config CaptchaConfig) sign(lotNumber string) (token string) {
	token = HmacEncode(config.Key, lotNumber)
	return
}

func (config CaptchaConfig) validateAPI() (url string) {
	url = fmt.Sprintf("%s/validate?captcha_id=%s", Host, config.ID)
	return
}

// HmacEncode  hmac-sha256
func HmacEncode(key string, data string) string {
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
