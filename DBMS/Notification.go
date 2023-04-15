package DBMS

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
	"html"
	"log"
	"sync"
	"time"
)

// PgMakeListener makes listener for PostgreSQL notification about user's messages
func PgMakeListener(DBMSinfo string) *pq.Listener {
	listener := pq.NewListener(DBMSinfo, 10*time.Second, time.Minute,
		func(ev pq.ListenerEventType, err error) {
			{
				if err != nil {
					log.Fatal(err)
				}
			}
		})
	err := listener.Listen("MessageAdded")
	if err != nil {
		log.Fatal(err)
	}
	return listener
}

// PgGetNotify get message notification and sends message to connected to chat users
func PgGetNotify(DBMSinfo string, chatConnections *sync.Map) {
	listener := PgMakeListener(DBMSinfo)
	defer listener.Close()

	for {
		select {
		case n := <-listener.Notify:
			NotificationMessage := struct {
				ChatId int  `json:"chatid"`
				UserId int  `json:"userid"`
				Content string    `json:"content"`
				Time    time.Time `json:"time"`
				Username string `json:"username"`
			}{}

			err := json.Unmarshal([]byte(n.Extra), &NotificationMessage)
			if err != nil {
				continue
			}

			namedMessage := NamedMessage{
				Username: html.EscapeString(NotificationMessage.Username),
				Content: html.EscapeString(NotificationMessage.Content),
				Time: NotificationMessage.Time,
			}
			marshal, err := json.Marshal(&namedMessage)
			if err != nil {
				continue
			}

			l, ok := chatConnections.Load(NotificationMessage.ChatId)
			if !ok {
				continue
			}
			connectionList := l.([]*websocket.Conn)
			for i := 0; i < len(connectionList); i++ {
				connection := connectionList[i]
				err = (*connection).WriteMessage(1, marshal)
				if err != nil {
					connectionList = append(connectionList[:i], connectionList[i+1:]...)
				}
			}
		}
	}
}