package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yarn/backend/db"
	"github.com/yarn/backend/request"
	"gorm.io/gorm"
)

type RewardService struct {
	db *gorm.DB
}

func NewRewardService(db *gorm.DB) *RewardService {
	return &RewardService{
		db: db,
	}
}
func (service *RewardService) RegisterHandler(e *gin.Engine) {
	r := e.Group("/api/reward")

	r.GET("", func(c *gin.Context) {
		// 从header中获取openid
		openid := c.GetHeader("openid")
		var res Resp = service.GetReward(openid)
		c.JSON(http.StatusOK, gin.H{"code": res.Code, "data": res.Data, "msg": res.Msg})
	})

	r.POST("/set", func(c *gin.Context) {
		// 从header中获取openid
		openid := c.GetHeader("openid")
		var payload *request.SetRewardReq
		if err := c.BindJSON(&payload); err != nil {
			c.Error(err)
			return
		}
		var res Resp = service.SetReward(payload, openid)
		c.JSON(http.StatusOK, gin.H{"code": res.Code, "data": res.Data, "msg": res.Msg})
	})
}

func (service *RewardService) GetReward(openid string) Resp {

	var reward db.Reward

	err := service.db.First(&reward, "openid = ?", openid).Error

	if err != nil {
		return Resp{Code: 1, Data: "", Msg: "用户奖励不存在"}
	}
	return Resp{Code: 0, Data: reward, Msg: ""}
}

func (service *RewardService) SetReward(data *request.SetRewardReq, openid string) Resp {

	var reward db.Reward

	err := service.db.First(&reward, "openid = ?", openid).Error

	if err != nil {
		return Resp{Code: 1, Data: "", Msg: "用户奖励不存在"}
	}

	err = service.db.Model(&reward).Updates(db.Reward{
		CrochetCount:  data.CrochetCount,
		KnittingCount: data.KnittingCount,
		ShareCount:    data.ShareCount,
		Lv1Count:      data.Lv1Count,
		Lv2Count:      data.Lv2Count,
		Lv3Count:      data.Lv3Count,
	}).Error

	if err != nil {
		return Resp{Code: 1, Data: "", Msg: "设置失败"}
	}

	return Resp{Code: 0, Data: "", Msg: "设置成功"}

}
