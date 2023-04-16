package service

import (
	"IM/models"
	"IM/utils"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// GetUserList
// @Summary 查询所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/GetUserList [get]
func GetUserList(ctx *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()

	ctx.JSON(http.StatusOK, gin.H{
		"message": "用户列表",
		"code":    0, // 0 ok -1 false
		"data":    data,
	})
}

// CreateUser
// @Summary 添加用户
// @Tags 用户模块
// @Param name query string false "用户名"
// @Param password query string false "密码"
// @Param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/CreateUser [post]
func CreateUser(ctx *gin.Context) {
	user := models.UserBasic{}
	user.Name = ctx.Query("name")
	PassWord := ctx.Query("password")
	RePassWord := ctx.Query("repassword")

	salt := fmt.Sprintf("%06d", rand.Int31())

	user.LoginTime = time.Now()
	user.HeartbeatTime = time.Now()
	user.LoginOutTime = time.Now()
	// user.Phone = phone
	// user.Email = email
	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		ctx.JSON(-1, gin.H{
			"message": "用户名已注册",
			"code":    -1, // 0 ok -1 false
			"data":    user,
		})
		return
	}

	if PassWord != RePassWord {
		ctx.JSON(-1, gin.H{
			"message": "两次密码不一致",
			"code":    -1, // 0 ok -1 false
			"data":    user,
		})
		return
	}

	user.Password = utils.MakePassword(PassWord, salt)
	user.Salt = salt
	models.CreateUser(user)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "新增用户成功",
		"code":    0, // 0 ok -1 false
		"data":    user,
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @Param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/DeleteUser [delete]
func DeleteUser(ctx *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(ctx.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0, // 0 ok -1 false
		"message": "删除成功",
		"data":    user,
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @Param id formData string false "id"
// @Param name formData string false "name"
// @Param password formData string false "password"
// @Param phone formData string false "phone"
// @Param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/UpdateUser [post]
func UpdateUser(ctx *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(ctx.PostForm("id"))
	user.ID = uint(id)

	user.Name = ctx.PostForm("name")
	user.Password = ctx.PostForm("password")
	user.Phone = ctx.PostForm("phone")
	user.Email = ctx.PostForm("email")
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(200, gin.H{
			"code":    -1, // 0 ok -1 false
			"message": "修改失败,不匹配",
			"data":    user,
		})
	} else {
		models.UpdateUser(user)
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0, // 0 ok -1 false
			"message": "修改用户成功",
			"data":    user,
		})
	}

}

// FindUserByNameAndPwd
// @Summary 用户登录
// @Tags 用户模块
// @Param name query string false "name"
// @Param password query string false "password"
// @Success 200 {string} json{"code","message"}
// @Router /user/FindUserByNameAndPwd [post]
func FindUserByNameAndPwd(ctx *gin.Context) {

	data := models.UserBasic{}

	name := ctx.Query("name")
	password := ctx.Query("password")
	user := models.FindUserByName(name)

	if user.Name == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1, // 0 ok -1 false
			"message": "该用户不存在",
			"data":    data,
		})
		return
	}

	flag := utils.ValidPassword(password, user.Salt, user.Password)
	if !flag {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1, // 0 ok -1 false
			"message": "密码不正确",
			"data":    data,
		})
		return
	}

	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPassWord(name, pwd)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0, // 0 ok -1 false
		"message": "登录成功",
		"data":    data,
	})
}

// 防止跨域站点的伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func (r *http.Request) bool {
		return true
	},
}

func SendMsg(ctx *gin.Context) {
	ws, err := upGrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, ctx)
}

func MsgHandler(ws *websocket.Conn, ctx *gin.Context) {
	for {
		msg, err := utils.Subscribe(ctx, utils.PublishKey)
		if err != nil {
			fmt.Println(err)
		}
		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}
	
}