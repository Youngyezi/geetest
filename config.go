package geetest

const (
	// Host Service host
	Host = "http://gcaptcha4.geetest.com"
)

// CaptchaConfig
type CaptchaConfig struct {
	ID  string
	Key string
}

// ConfigDefault  test config
var ConfigDefault = CaptchaConfig{
	ID:  "647f5ed2ed8acb4be36784e01556bb71",
	Key: "b09a7aafbfd83f73b35a9b530d0337bf",
}
