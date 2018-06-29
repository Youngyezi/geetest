# geetest
极验/geetest golang

最近用到了极验，发现没有golang的sdk，于是乎自己撸了个轮子，翻译了下。


#usage
```golang
    g := geetest.New("CaptchaID", "PrivateKey")
	u := url.Values{
		"uid":         {"test"},
		"client_type": {"web"},
		"ip_address":  {"127.0.0.1"},
	}
```