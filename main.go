package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"ihomegin/controller"
	"ihomegin/model"
	"ihomegin/utils"
	"net/http"
)

func lenater(ctx *gin.Context)  {
	resp:=make(map[string]interface{})
	session:=sessions.Default(ctx)
	userName:=session.Get("userName")
	if userName.(string)==""{
		fmt.Println("没有session 回去")
		ctx.Abort()
		resp["errno"]=utils.RECODE_DBERR
		resp["errmsg"]=utils.RECODE_DBERR
		ctx.JSON(http.StatusOK,resp)
	}
	resp["errno"]=utils.RECODE_OK
	resp["errmsg"]=utils.RecodeText(utils.RECODE_OK)
	ctx.Next()
}

func main() {
	router :=gin.Default()

	//初始化数据redis
	model.InitRedis()
	//请求分配

	//静态路由
	router.Static("/home","view")
	r1 :=router.Group("api/v1.0")

	//使用中间建
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	r1.Use(sessions.Sessions("mysession", store))

	{
		//路由规范  Restful  api
		r1.GET("/areas",controller.GetArea)

		r1.GET("/imagecode/:uuid/",controller.GetImageCd)
		r1.GET("/smscode/:mobile",controller.GetSmscd)


		r1.GET("/session",controller.GetSession)
		r1.POST("/users",controller.PostRet)
		r1.POST("/sessions",controller.PostLogin)
		r1.Use(lenater)
		r1.GET("/user",controller.GetUserInfo)
		r1.PUT("/user/name",controller.PutUserInfo)
		r1.DELETE("/session",controller.DeleteSession)
	}



	router.Run(":8081")

}
