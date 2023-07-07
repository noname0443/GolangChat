package main

import (
	"GolangChat/DBMS"
	"GolangChat/Utility"
	"GolangChat/api"
	"GolangChat/sockets"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"sync"
)

func main() {
	log.Println("Connecting to PostgreSQL DBMS")
	err, configs := Utility.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	DBMSString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", configs["user"], configs["password"], configs["dbname"], configs["sslmode"])
	db, err := sql.Open("postgres", DBMSString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Successful connection to PostgreSQL DBMS")

	var chatConnections sync.Map

	go DBMS.PgGetNotify(DBMSString, &chatConnections)

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/resources", "./resources")
	router.POST("/api/login", func(context *gin.Context) { api.Login(context, db) })
	router.POST("/api/register", func(context *gin.Context) { api.Register(context, db) })
	router.GET("/api/chats", func(context *gin.Context) { api.GetChats(context, db) })
	router.POST("/api/chats/make", func(context *gin.Context) { api.MakeChat(context, db) })
	router.POST("/api/chats/invite", func(context *gin.Context) { api.InviteUserToChat(context, db) })
	router.GET("/ws", func(c *gin.Context) {
		sockets.WSHandler(c.Writer, c.Request, db, &chatConnections)
	})

	router.GET("/chat", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", gin.H{})
	})
	router.GET("/chats", func(c *gin.Context) {
		c.HTML(http.StatusOK, "rooms.html", gin.H{})
	})
	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "registration.html", gin.H{})
	})
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	err = router.Run(configs["socket"])
	if err != nil {
		log.Fatal(err)
	}
}
