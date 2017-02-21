# geetest
极验/geetest golang

最近用到了极验，发现没有golang的sdk，于是乎自己撸了个轮子，翻译了下。


#usage
```golang
    g := &geetest.Geetest{
		CaptchaID:  "",
		PrivateKey: "",
	}

    //预请求获取gt-server状态
    serverStatus := g.PreProcess(userID)
    //gt-server正常，向gt-server进行二次验证
    g.SuccessValidate(challenge, validate, seccode)
    //t-server非正常情况下，进行failback模式验证
    g.FailValidate(challenge, validate)
```