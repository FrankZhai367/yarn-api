package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yarn/backend/db"
	"gorm.io/gorm"
)

type FinishedService struct {
	db *gorm.DB
}

func NewFinishedService(db *gorm.DB) *FinishedService {
	return &FinishedService{
		db: db,
	}
}
func (service *FinishedService) RegisterHandler(e *gin.Engine) {
	r := e.Group("/api/finished")

	r.GET("", func(c *gin.Context) {
		openid := c.GetHeader("openid")
		var res Resp = service.GetFinishedList(openid)
		c.JSON(http.StatusOK, gin.H{"code": res.Code, "data": res.Data, "msg": res.Msg})
	})

	r.POST("/finish/:objectId", func(c *gin.Context) {
		var objectId = c.Param("objectId")
		openid := c.GetHeader("openid")

		var res Resp = service.Finish(objectId, openid)
		c.JSON(http.StatusOK, gin.H{"code": res.Code, "data": res.Data, "msg": res.Msg})
	})

}

func (service *FinishedService) GetFinishedList(openid string) Resp {
	var finishedIds []string
	err := service.db.Model(&db.Finished{}).Where("openid = ?", openid).Pluck("object_id", &finishedIds).Error

	if err != nil {
		return Resp{Code: 1, Data: "", Msg: "获取失败"}
	}
	return Resp{Code: 0, Data: finishedIds, Msg: ""}
}

func (service *FinishedService) Finish(objectId string, openid string) Resp {
	var finished db.Finished
	err := service.db.First(&finished, "object_id = ? AND openid = ?", objectId, openid).Error

	if err != nil {
		// 不存在则创建
		finished = db.Finished{
			ObjectId: objectId,
			Openid:   openid,
		}
		service.db.Create(&finished)

	} else {
		// 存在 则提示已经完成
		return Resp{Code: 1, Data: "", Msg: "已经完成"}
	}

	return Resp{Code: 0, Data: finished, Msg: "操作成功"}

}
