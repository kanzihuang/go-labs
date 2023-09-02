package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	conn, err := gorm.Open(mysql.Open("reader:123456@tcp(mysql.123sou.cn:3306)/lab"))
	if err != nil {
		log.Fatal("数据库连接失败, " + err.Error())
	}
	roster := NewDbRoster(conn)
	//roster := NewMapRoster()
	runServer(roster)
}

func runServer(roster Roster) {
	r := gin.Default()
	r.POST("/registry", func(context *gin.Context) {
		var person Person
		err := context.BindJSON(&person)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"errorMessage": "读取数据失败，" + err.Error(),
			})
			return
		}
		if err := roster.Registry(person); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"errorMessage": "注册失败，" + err.Error(),
			})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"message": "Success",
		})
	})
	r.GET("/query/:name", func(context *gin.Context) {
		name := context.Param("name")
		if name == "" {
			context.JSON(http.StatusBadRequest, gin.H{
				"errorMessage": "name参数未设置",
			})
			return
		}
		person, err := roster.Get(name)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"errorMessage": "获取信息失败，" + err.Error(),
			})
			return
		}
		context.JSON(http.StatusOK, person)
	})
	r.GET("/all", func(context *gin.Context) {
		persons, err := roster.All()
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"errorMessage": "无法获取数据，" + err.Error(),
			})
			return
		}
		context.JSON(http.StatusOK, persons)
	})
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
