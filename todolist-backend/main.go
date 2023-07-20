package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Tasklist struct {
	TaskID     int64  `gorm:"primary_key;auto_increment;not null;"`
	TaskName   string `json:"task_name"`
	TaskDetail string `json:"task_detail"`
	TaskState  bool   `json:"task_state"`
}

var (
	db *gorm.DB
)

func InitDB() (err error) {
	dsn := "root:12345678@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("init error: ", err)
		return err
	}
	return nil
}

func main() {
	if err := InitDB(); err != nil {
		panic("conect error")
	}
	db.AutoMigrate(&Tasklist{})

	r := gin.Default()
	r.Static("/static", "static")
	r.LoadHTMLGlob("./template/index.html")

	v1 := r.Group("/")
	v1.Use(Cors())
	{
		v1.GET("/", GetHomePage)
		v1.GET("/tasklist", GetTaskList)
		v1.POST("/tasklist", AddTask)
		v1.POST("/tasklist/delete", DeleteTask)
		v1.POST("/tasklist/state", UpdateTaskState)
	}

	r.Run(":8081")
}

// 设置后端跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}

// 返回首页
func GetHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

//获取任务列表
func GetTaskList(c *gin.Context) {
	var QueryTask []Tasklist
	db.Find(&QueryTask)
	c.JSON(http.StatusOK, QueryTask)
}

// 添加任务
func AddTask(c *gin.Context) {
	bodydata := make(map[string]string)
	c.BindJSON(&bodydata)
	InsertTask := Tasklist{TaskName: bodydata["task_name"], TaskDetail: bodydata["task_detail"], TaskState: false}
	db.Create(&InsertTask)
}

// 删除任务
func DeleteTask(c *gin.Context) {
	bodydata := make(map[string]int)
	c.BindJSON(&bodydata)
	db.Where("task_id = ?", bodydata["TaskID"]).Delete(Tasklist{})
}

// 更新任务状态
func UpdateTaskState(c *gin.Context) {
	bodydata := make(map[string]string)
	c.BindJSON(&bodydata)
	db.Model(&Tasklist{}).Where("task_id = ?", bodydata["TaskID"]).Update("task_state", bodydata["task_state"])
}
