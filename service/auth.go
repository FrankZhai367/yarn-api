package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var dbMock = make(map[string]string)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		db: db,
	}
}
func (service *AuthService) RegisterHandler(e *gin.Engine) {
	r := e.Group("/api/auth")

	// Authorized group (uses gin.BasicAuth() middleware)
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar11", // user:foo password:bar
		"manu": "123",   // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	Zm9vOmJhcg== is base64("foo:bar")

	curl -X POST \
		http://localhost:8080/admin \
		-H 'authorization: Basic Zm9vOmJhcg==' \
		-H 'content-type: application/json' \
		-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		fmt.Println("json: ", json)
		fmt.Println("user ", user)

		if c.Bind(&json) == nil {
			dbMock[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})
}
