package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry/consul"
	MictoGetArea "ihomegin/proto/getArea/proto/getArea"
	getImg "ihomegin/proto/getImg"
	Register "ihomegin/proto/register"
	UserInfo "ihomegin/proto/user"
	"ihomegin/utils"
	"image/png"
	"net/http"
)

//获取所以信息
func GetArea(ctx *gin.Context) {
	//resp :=make(map[string]interface{})
	//	Areas,err:=model.GetArea()
	//	if err != nil {
	//		resp["errno"]=utils.RECODE_DBERR
	//		resp["errmsg"]=utils.RecodeText(utils.RECODE_DBERR)
	//		ctx.JSON(http.StatusOK,resp)
	//		return
	//	}
	//	resp["errno"]=utils.RECODE_OK
	//	resp["errmsg"]=utils.RecodeText(utils.RECODE_OK)
	//	resp["data"]=Areas
	//	fmt.Println(resp)
	//	ctx.JSON(http.StatusOK,resp)
	GetAreaClient := MictoGetArea.NewGetAreaService("go.micro.srv.getArea", client.DefaultClient)
	resp, err := GetAreaClient.MicrogetArea(context.TODO(), &MictoGetArea.Request{})
	if err != nil {
		fmt.Println(err)
		fmt.Println("err:", resp)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	fmt.Println(resp)
	ctx.JSON(http.StatusOK, resp)
}

func GetSession(ctx *gin.Context) {
	resp := make(map[string]interface{})
	dataTmp:=make(map[string]string)
	//初始化
	session:=sessions.Default(ctx)
	userName:=session.Get("userName")

	if userName==nil{
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
	}else {
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
		dataTmp["name"]=userName.(string)
		resp["data"]=dataTmp
	}

	ctx.JSON(http.StatusOK, resp)
}

func GetImageCd(ctx *gin.Context) {

	uuid := ctx.Param("uuid")

	consulReg := consul.NewRegistry()
	MicroService := micro.NewService(
		micro.Registry(consulReg),
	)
	ImgService := getImg.NewGetImgService("go.micro.srv.getImg", MicroService.Client())
	resp, err := ImgService.MicroGetImg(context.TODO(), &getImg.Request{Uuid: uuid})
	if err != nil {
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	var imges captcha.Image
	json.Unmarshal(resp.Data, &imges)
	png.Encode(ctx.Writer, imges)
	ctx.JSON(http.StatusOK, resp)
}

//获取短信验证码
func GetSmscd(ctx *gin.Context) {
	mobile := ctx.Param("mobile")
	text := ctx.Query("text")
	id := ctx.Query("id")
	if mobile == "" || text == "" || id == "" {
		fmt.Println("传过来是空值错误:err")
		return
	}
	consulReg := consul.NewRegistry()
	MicroService := micro.NewService(
		micro.Registry(consulReg),
	)
	ImgService := Register.NewRegisterService("go.micro.srv.register", MicroService.Client())
	resp, err := ImgService.SmsCode(context.TODO(), &Register.Request{
		Mobile: mobile,
		Text:   text,
		Uuid:   id,
	})
	if err != nil {
		fmt.Println("发送短信失败", err)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

type RestReter struct {
	Mobile string `json:"mobile"`
	Password string `json:"password"`
	Sms_code string `json:"sms_code"`
}

//注册用户
func PostRet(ctx *gin.Context) {
	session:=sessions.Default(ctx)
	/*mobile := ctx.PostForm("mobile")
	password := ctx.PostForm("password")
	sms_code := ctx.PostForm("sms_code")
	这里怎么
		fmt.Println("获取数据",mobile,password,sms_code)
	 */
	var Restet RestReter
	err:=ctx.Bind(&Restet)
	if err != nil {
		fmt.Println("前端数据获取错误",err)
		return
	}
	fmt.Println(Restet)
	mobile:=Restet.Mobile
	password:=Restet.Password
	sms_code:=Restet.Sms_code
	consulReg := consul.NewRegistry()
	MicroService := micro.NewService(
		micro.Registry(consulReg),
	)
	RetService := Register.NewRegisterService("go.micro.srv.register", MicroService.Client())
	resp,err:=RetService.RegisCode(context.TODO(), &Register.Regrequest{
		Mobile:   mobile,
		Password: password,
		SmsCode:  sms_code,
	})
	if err != nil {
		fmt.Println("注册插入数据错误mysql err:",err)
		ctx.JSON(http.StatusOK,resp)
		return
	}
	fmt.Println("插入数据成功了")
	//注册成功就自己抵用到
	fmt.Println("mobile -----",mobile)
	session.Set("userName",mobile)
	session.Save()
	ctx.JSON(http.StatusOK,resp)
}

//展现用户登录信息
func GetUserInfo(ctx *gin.Context)  {
	//在session获取用户名
	session:=sessions.Default(ctx)
	userName:=session.Get("userName")
	fmt.Println("<============到这里了",userName.(string))
	//调用远程服务
	consulCof:=consul.NewRegistry()
	MicroClient:=micro.NewService(
		micro.Registry(consulCof),
		)
	UserService:=UserInfo.NewUserService("go.micro.srv.user",MicroClient.Client())
	resp,err:=UserService.MicroUser(context.TODO(),&UserInfo.Request{
		Name:userName.(string),
	})

	if err != nil {
		fmt.Println("用户登录信息错误:err",err)
		ctx.JSON(http.StatusOK,resp)
		return
	}
	ctx.JSON(http.StatusOK,resp)
}
type UserStu struct {
	Mobile string `json:"mobile"`
	Password string `json:"password"`
}
//登录用户实现
func PostLogin(ctx *gin.Context)  {
	session:=sessions.Default(ctx)
	resp:=make(map[string]interface{})
	//获取用户数据
	var loguser UserStu
	err:=ctx.Bind(&loguser)
	if err != nil {
		resp["errno"]=utils.RECODE_NODATA
		resp["errer"]=utils.RecodeText(utils.RECODE_NODATA)
		ctx.JSON(http.StatusOK,resp)
		return
	}
	consulcfg:=consul.NewRegistry()
	MicroService:=micro.NewService(
		micro.Registry(consulcfg),
		)
	loginserver:=Register.NewRegisterService("go.micro.srv.register",MicroService.Client())
	rsp,err:=loginserver.LoginPost(context.TODO(),&Register.LoginRequest{
		Mobile:loguser.Mobile,
		Password:loguser.Password,
	})//rsp.Name
	if err != nil {
		ctx.JSON(http.StatusOK,rsp)
		fmt.Println("调用远程服务错误:err",err)
		return
	}
	fmt.Println("<=================rsp.Name",rsp.Name)
	session.Set("userName",rsp.Name)
	session.Save()
	ctx.JSON(http.StatusOK,rsp)
}

//删除用户状态(session)
func DeleteSession(ctx *gin.Context)  {
	fmt.Println("调用这里")
	resp:=make(map[string]interface{})
	session:=sessions.Default(ctx)
	session.Delete("userName")
	err:=session.Save()
	if err != nil {
		resp["errno"]=utils.RECODE_SESSIONERR
		resp["errmsg"]=utils.RecodeText(utils.RECODE_SESSIONERR)
		ctx.JSON(http.StatusOK,resp)
		return
	}
	resp["errno"]=utils.RECODE_OK
	resp["errmsg"]=utils.RecodeText(utils.RECODE_OK)
	ctx.JSON(http.StatusOK,resp)
}

type Updatauser struct {
	Name string
}
//更新用户名
func PutUserInfo(ctx *gin.Context)  {
	session:=sessions.Default(ctx)
	userName:=session.Get("userName")
	var UserUpdata Updatauser
	err:=ctx.Bind(&UserUpdata)
	if userName.(string) ==""|| err!=nil{
		fmt.Println("userName 数据不完成")
	}

	consulcof:=consul.NewRegistry()
	MicroConsul:=micro.NewService(
		micro.Registry(consulcof),
		)

	UserService:=UserInfo.NewUserService("go.micro.srv.user",MicroConsul.Client())
	resp,err:=UserService.UpdataUser(context.TODO(),&UserInfo.UpdataUserRequest{
		Oldname:userName.(string),
		Xinname:UserUpdata.Name,
	})
	if err != nil {
		resp.Errno=utils.RECODE_DBERR
		resp.Errmsg=utils.RecodeText(utils.RECODE_DBERR)
		fmt.Println("调用远程服务错误:",err)
		ctx.JSON(http.StatusOK,resp)
		return
	}

	session.Set("userName",resp.Data.Name)
	session.Save()
	ctx.JSON(http.StatusOK,resp)
}