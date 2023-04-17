package service

import (
	"IM/models"
	"html/template"
	_ "net/http"
	"strconv"
	_ "text/template/parse"

	"github.com/gin-gonic/gin"
)

// GetIndex
func GetIndex(ctx *gin.Context) {
	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(ctx.Writer, "index")

}

// ToRegister

func ToRegister(ctx *gin.Context) {
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(ctx.Writer, "register")

}

// ToChat

func ToChat(ctx *gin.Context) {
	ind, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/foot.html",
		"views/chat/tabmenu.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/main.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
		"views/chat/concat.html")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(ctx.Query("userId"))
	token := ctx.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
	ind.Execute(ctx.Writer, user)
}

func Chat(ctx *gin.Context) {
	models.Chat(ctx.Writer, ctx.Request)
}