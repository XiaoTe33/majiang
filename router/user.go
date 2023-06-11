package router

import (
	"github.com/gin-gonic/gin"
	"majiang/errors"
	"majiang/model"
	"majiang/utils"
)

func register(c *gin.Context) {
	u := c.PostForm("username")
	p := c.PostForm("password")
	user := model.User{}

	if handleError(c, db.Where("username = ? ", u).Find(&user).Error) {
		return
	}
	if user != (model.User{}) {
		jsonError(c, errors.ErrUsernameExist)
		return
	}
	ctx := db.Begin()

	if handleError(c, ctx.Create(&model.User{Id: utils.GetUserId(), Username: u, Password: utils.Md5Encoded(p)}).Error) {
		ctx.Rollback()
		return
	}
	ctx.Commit()
	jsonSuccess(c)
}

func login(c *gin.Context) {
	u := c.PostForm("username")
	p := c.PostForm("password")
	user := model.User{}
	if handleError(c, db.Where("username = ? and password = ?", u, utils.Md5Encoded(p)).Find(&user).Error) {
		return
	}
	if user == (model.User{}) {
		jsonError(c, errors.ErrWrongUsernameOrPassword)
		return
	}
	id := user.Id
	jsonData(c, gin.H{
		"refreshToken": utils.GenerateRefreshToken(&utils.MyClaim{Id: id}),
		"accessToken":  utils.GenerateAccessToken(&utils.MyClaim{Id: id}),
	})

}

func refreshToken(c *gin.Context) {
	token := c.Query("refreshToken")
	id, err := utils.IsRefreshToken(token)
	if handleError(c, err) {
		return
	}
	jsonData(c, gin.H{
		"refreshToken": utils.GenerateRefreshToken(&utils.MyClaim{Id: id}),
		"accessToken":  utils.GenerateAccessToken(&utils.MyClaim{Id: id}),
	})
}
