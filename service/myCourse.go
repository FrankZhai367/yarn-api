package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yarn/backend/db"
	"gorm.io/gorm"
)

type MyCourseService struct {
	db *gorm.DB
}

func NewMyCourseService(db *gorm.DB) *MyCourseService {
	return &MyCourseService{
		db: db,
	}
}
func (service *MyCourseService) RegisterHandler(e *gin.Engine) {
	r := e.Group("/api/myCourse")

	r.GET("", func(c *gin.Context) {
		openid := c.GetHeader("openid")
		var res Resp = service.GetMyCourseList(openid)
		c.JSON(http.StatusOK, gin.H{"code": res.Code, "data": res.Data, "msg": res.Msg})
	})

	r.POST("/toggleMyCourse/:courseId", func(c *gin.Context) {
		var courseId = c.Param("courseId")
		openid := c.GetHeader("openid")

		var res Resp = service.ToggleMyCourse(courseId, openid)
		c.JSON(http.StatusOK, gin.H{"code": res.Code, "data": res.Data, "msg": res.Msg})
	})

}

func (service *MyCourseService) GetMyCourseList(openid string) Resp {

	var courseIds []string

	err := service.db.Model(&db.MyCourse{}).Where("openid = ?", openid).Pluck("course_id", &courseIds).Error

	if err != nil {
		return Resp{Code: 1, Data: "", Msg: "获取失败"}
	}
	return Resp{Code: 0, Data: courseIds, Msg: ""}
}

func (service *MyCourseService) ToggleMyCourse(courseId string, openid string) Resp {

	var myCourse db.MyCourse

	err := service.db.First(&myCourse, "course_id = ? AND openid = ?", courseId, openid).Error

	if err != nil {
		// 不存在则创建
		myCourse = db.MyCourse{
			CourseId: courseId,
			Openid:   openid,
		}
		newEerr := service.db.Create(&myCourse).Error
		if newEerr != nil {
			return Resp{Code: 1, Data: "", Msg: "创建失败"}
		}
		fmt.Println("不存在创建 myCourse", newEerr)
	} else {
		// 存在则删除
		service.db.Delete(&myCourse)
		fmt.Println("存在删除 myCourse", myCourse)
	}

	return Resp{Code: 0, Data: myCourse, Msg: "操作成功"}

}
