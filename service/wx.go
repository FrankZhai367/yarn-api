package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WxService struct {
	db *gorm.DB
}

type GetOpenIdRes struct {
	Openid string `json:"openid"`
	Errmsg string `json:"errmsg"`
}

func NewWxService(db *gorm.DB) *WxService {
	return &WxService{
		db: db,
	}
}
func (service *WxService) RegisterHandler(e *gin.Engine) {
	r := e.Group("/api/wx")

	r.GET("/openid", func(c *gin.Context) {
		JSCODE := c.Query("JSCODE")
		var res Resp = service.GetOpenId(JSCODE)
		c.JSON(http.StatusOK, gin.H{"code": res.Code, "data": res.Data, "msg": res.Msg})
	})
}

func (service *WxService) GetOpenId(JSCODE string) Resp {
	const APPID = "wx06817db9a249009b"
	const SECRET = "48bf8af2c3e6c295faefc5efd5c65725"
	const authorization_code = "authorization_code"

	var url = "https://api.weixin.qq.com/sns/jscode2session?appid=" + APPID + "&secret=" + SECRET + "&js_code=" + JSCODE + "&grant_type=" + authorization_code

	// 请求微信服务器
	resp, err := http.Get(url)
	if err != nil {
		return Resp{Code: 1, Data: "", Msg: "请求微信服务器失败"}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("xxxxxxxxx", string(body))
	var res GetOpenIdRes
	_ = json.Unmarshal(body, &res)

	if res.Errmsg != "" {
		return Resp{Code: 1, Data: "", Msg: res.Errmsg}
	}

	return Resp{Code: 0, Data: res.Openid, Msg: ""}
}
