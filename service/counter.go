package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yarn/backend/db"
	"github.com/yarn/backend/request"
	"gorm.io/gorm"
)

type CounterService struct {
	db *gorm.DB
}

func NewCounterService(db *gorm.DB) *CounterService {
	return &CounterService{
		db: db,
	}
}
func (service *CounterService) RegisterHandler(e *gin.Engine) {
	r := e.Group("/api/counter")

	r.GET("", func(c *gin.Context) {
		// 从header中获取openid
		openid := c.GetHeader("openid")
		var res Resp = service.GetCounterList(openid)
		c.JSON(http.StatusOK, gin.H{"code": res.Code, "data": res.Data, "msg": res.Msg})
	})

	r.POST("/delete/:id", func(c *gin.Context) {
		openid := c.GetHeader("openid")
		var id = c.Param("id")
		var res Resp = service.DeleteCounter(id, openid)
		c.JSON(http.StatusOK, gin.H{"code": res.Code, "data": res.Data, "msg": res.Msg})
	})

	r.POST("/update/:id", func(c *gin.Context) {
		var id = c.Param("id")
		var payload *request.UpdateCounterReq
		c.BindJSON(&payload)
		var res Resp = service.UpdateCounter(payload, id)
		c.JSON(http.StatusOK, gin.H{"code": res.Code, "data": res.Data, "msg": res.Msg})
	})

	r.POST("/add", func(c *gin.Context) {
		var payload *request.AddCounterReq
		c.BindJSON(&payload)
		var res Resp = service.AddCounter(payload)
		c.JSON(http.StatusOK, gin.H{"code": res.Code, "data": res.Data, "msg": res.Msg})
	})
}

func (service *CounterService) GetCounterList(openid string) Resp {

	var counters []db.Counter

	err := service.db.Where("openid = ?", openid).Find(&counters).Error

	if err != nil {
		return Resp{Code: 1, Data: "", Msg: "获取失败"}
	}
	return Resp{Code: 0, Data: counters, Msg: ""}
}

func (service *CounterService) DeleteCounter(id string, openid string) Resp {

	var counter db.Counter

	err := service.db.Where("id = ? and openid = ?", id, openid).Delete(&counter).Error

	if err != nil {
		return Resp{Code: 1, Data: "", Msg: "删除失败"}
	}
	return Resp{Code: 0, Data: "", Msg: "删除成功"}
}

func (service *CounterService) UpdateCounter(data *request.UpdateCounterReq, id string) Resp {

	var counter db.Counter

	err := service.db.First(&counter, "id = ?", id).Error

	if err != nil {
		return Resp{Code: 1, Data: "", Msg: "计数器不存在"}
	}

	err = service.db.Model(&counter).Select("Count", "Name").Updates(db.Counter{
		Count: data.Count,
		Name:  data.Name,
	}).Error

	if err != nil {
		return Resp{Code: 1, Data: "", Msg: "更新失败"}
	}
	return Resp{Code: 0, Data: "", Msg: "更新成功"}
}

func (service *CounterService) AddCounter(data *request.AddCounterReq) Resp {

	var counter db.Counter = db.Counter{
		Name:   data.Name,
		Count:  data.Count,
		Openid: data.Openid,
	}

	err := service.db.Create(&counter).Error

	if err != nil {
		return Resp{Code: 1, Data: "", Msg: err.Error()}
	}
	return Resp{Code: 0, Data: counter, Msg: "添加成功"}
}
