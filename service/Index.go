package service

import (
	"html/template"
	_ "net/http"
	_ "text/template/parse"

	"github.com/gin-gonic/gin"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(ctx *gin.Context) {
	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(ctx.Writer, "index")
	// ctx.JSON(http.StatusOK, gin.H{
	// 	"message": "welcome !!",
	// })
}

// ToRegister
// @Tags 与createUSer 相同
// @Success 200 {string} welcome
// @Router /toRegister [post]
func ToRegister(ctx *gin.Context) {
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(ctx.Writer, "register")
	// ctx.JSON(http.StatusOK, gin.H{
	// 	"message": "welcome !!",
	// })
}