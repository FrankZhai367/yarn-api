package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yarn/backend/db"
	"github.com/yarn/backend/request"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

type Resp struct {
	Code int
	Data interface{}
	Msg  string
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}
func (service *UserService) RegisterHandler(e *gin.Engine) {
	r := e.Group("/api/user")

	r.GET("", func(c *gin.Context) {
		openid := c.GetHeader("openid")
		var res Resp = service.GetUser(openid)
		c.JSON(http.StatusOK, gin.H{"code": res.Code, "data": res.Data, "msg": res.Msg})
	})

	r.POST("/update", func(c *gin.Context) {
		openid := c.GetHeader("openid")
		var payload *request.UpdateUserReq
		c.BindJSON(&payload)
		var res Resp = service.UpdateUser(openid, payload)
		c.JSON(http.StatusOK, gin.H{"code": res.Code, "data": res.Data, "msg": res.Msg})
	})
}

// TODO 添加加密方式，防止crsf攻击

func (service *UserService) UpdateUser(openid string, payload *request.UpdateUserReq) Resp {

	var user db.User

	err := service.db.First(&user, "openid = ?", openid).Error

	if err != nil {
		return Resp{Code: 1, Data: "", Msg: "用户不存在"}
	}

	if payload.NickName != "" {
		user.NickName = payload.NickName
	}
	if payload.AvatarUrl != "" {
		user.AvatarUrl = payload.AvatarUrl
	}

	service.db.Save(&user)

	return Resp{Code: 0, Data: user, Msg: ""}
}

func (service *UserService) GetUser(openid string) Resp {

	var user db.User

	err := service.db.First(&user, "openid = ?", openid).Error

	if user.Openid == "" || err != nil {
		if openid == "" {
			return Resp{Code: 1, Data: map[string]string{}, Msg: "用户不存在"}
		} else {
			// 创建用户，创建reward
			user = db.User{Openid: openid}
			service.db.Create(&user)
			reward := db.Reward{
				Openid:        openid,
				CrochetCount:  0,
				KnittingCount: 0,
				Lv1Count:      0,
				Lv2Count:      0,
				Lv3Count:      0,
				ShareCount:    0,
			}
			service.db.Create(&reward)
		}
	}
	return Resp{Code: 0, Data: user, Msg: ""}
}
