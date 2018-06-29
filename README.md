# geetest
极验 (geetest) golang sdk

最近用到了极验，发现没有golang的sdk，于是乎自己撸了个轮子


## usage ##
```golang

    g := geetest.New("48a6ebac4ebc6642d68c217fca33eb4d", "4f1c085290bec5afdc54df73535fc361")
    p := url.Values{
	"user_id":     {"test"},
	"client_type": {"web"},
	"ip_address":  {"127.0.0.1"},
    }
    
    //预请求
    g.PreProcess(p)

    //服务器正常 校验
    g.SuccessValidate("geetest_challenge", "geetest_validate", "geetest_seccode", p)

   //服务器宕机 走failback模式
    g.FailValidate("geetest_challenge", "geetest_validate", "geetest_seccode")


```
