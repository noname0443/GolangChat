package sockets

import (
	"GolangChat/DBMS"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"html"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func WSHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, chatConnections *sync.Map) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade:", err)
		return
	}

	go func() {
		defer conn.Close()

		_, handshakeMessage, err := conn.ReadMessage()

		stringValues := strings.Split(string(handshakeMessage), " ")
		if len(stringValues) < 3 {
			return
		}
		sessionPass := stringValues[0]
		chat_name := stringValues[1]
		messageTimeStr := stringValues[2]
		messageTime, err := time.Parse(time.RFC3339, messageTimeStr)
		if err != nil {
			return
		}
		fmt.Println(sessionPass)
		err, user := DBMS.GetUserBySession(db, sessionPass)
		if err != nil {
			return
		}
		err, chat := DBMS.GetChat(db, user, chat_name)
		if err != nil {
			return
		}
		var connections []*websocket.Conn
		connectionsT, ok := chatConnections.LoadOrStore(chat.Id, []*websocket.Conn{conn})
		if ok {
			connections = connectionsT.([]*websocket.Conn)
			connections = append(connections, conn)
			chatConnections.Store(chat.Id, connections) // TODO: Remove Message duplication race Condition
		}

		err, allMessages := DBMS.GetMessagesWithAuthors(db, chat, 20, messageTime)
		if err != nil {
			return
		}
		for i := 0; i < len(allMessages); i++ {
			allMessages[i].Username = html.EscapeString(allMessages[i].Username)
			allMessages[i].Content = html.EscapeString(allMessages[i].Content)
		}
		marshal, err := json.Marshal(allMessages)
		err = conn.WriteMessage(1, marshal)
		if err != nil {
			return
		}

		for { // TODO: Checks that user session is correct
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Failed to read message from client:", err)
				break
			}
			err = DBMS.SendMessage(db, user, chat, string(message))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}()
}