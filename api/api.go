package api

import (
	"GolangChat/DBMS"
	"database/sql"
	"github.com/gin-gonic/gin"
	"html"
)

// Register user in chat service and send their cookie
func Register(c *gin.Context, db *sql.DB) {
	var SimpleUser struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}
	err := c.ShouldBind(&SimpleUser)
	if err != nil {
		c.Redirect(301, "/register")
		return
	}
	err, user := DBMS.RegisterUser(db, SimpleUser.Username, SimpleUser.Password)
	if err != nil {
		c.Redirect(301, "/register")
		return
	}
	c.SetCookie("Session", *user.SessionPass, 86400, "/", "", false, false)
	c.Redirect(301, "/chats")
}

// MakeChat makes chat for the user on service
func MakeChat(c *gin.Context, db *sql.DB) {
	cookie, err := c.Cookie("Session")
	if err != nil {
		c.Redirect(301, "/chats")
		return
	}
	var SimpleChat struct {
		Name string `form:"name"`
	}
	err = c.ShouldBind(&SimpleChat)
	if err != nil || SimpleChat.Name == "" {
		c.Redirect(301, "/chats")
		return
	}
	err, user := DBMS.GetUserBySession(db, cookie)
	if err != nil {
		c.Redirect(301, "/chats")
		return
	}
	err, chat := DBMS.MakeChat(db, user, SimpleChat.Name)
	if err != nil {
		c.Redirect(301, "/chats")
		return
	}
	chat.Name = html.EscapeString(chat.Name)
	c.Redirect(301, "/chats")
}

// InviteUserToChat invites user to chat
func InviteUserToChat(c *gin.Context, db *sql.DB) {
	cookie, err := c.Cookie("Session")
	if err != nil {
		c.Redirect(301, "/chats")
		return
	}
	var SimpleInvite struct {
		Username string `form:"username"`
		Chat string `form:"chat"`
	}
	err = c.ShouldBind(&SimpleInvite)
	if err != nil || SimpleInvite.Username == "" || SimpleInvite.Chat == "" {
		c.Redirect(301, "/chats")
		return
	}
	err, user := DBMS.GetUserBySession(db, cookie)
	if err != nil {
		c.Redirect(301, "/chats")
		return
	}
	err, chat := DBMS.GetChat(db, user, SimpleInvite.Chat)
	if err != nil {
		c.Redirect(301, "/chats")
		return
	}
	err = DBMS.AddUserInChat(db, chat, SimpleInvite.Username)
	if err != nil {
		c.Redirect(301, "/chats")
		return
	}
	c.Redirect(301, "/chat#" + SimpleInvite.Chat)
}

// Login user in chat service and send their cookie
func Login(c *gin.Context, db *sql.DB) {
	var SimpleUser struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}
	err := c.ShouldBind(&SimpleUser)
	if err != nil {
		c.Redirect(301, "/login")
		return
	}
	err, user := DBMS.UpdateSession(db, SimpleUser.Username, SimpleUser.Password)
	if err != nil {
		c.Redirect(301, "/login")
		return
	}
	c.SetCookie("Session", *user.SessionPass, 86400, "/", "", false, false)
	c.Redirect(301, "/chats")
}

// GetChats Outputs all user chats
func GetChats(c *gin.Context, db *sql.DB) {
	cookie, err := c.Cookie("Session")
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}
	err, user := DBMS.GetUserBySession(db, cookie)
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}
	err, chats := DBMS.GetChats(db, user)
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}
	for i := 0; i < len(chats); i++ {
		chats[i].Name = html.EscapeString(chats[i].Name)
	}
	c.JSON(200, chats)
}