package DBMS

import (
	"database/sql"
	"errors"
	"time"
)

type Chat struct{
	Id int `json:"-"`
	Name string `json:"name"`
}

type Message struct{
	ChatId int  `json:"chatid"`
	UserId int  `json:"userid"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

type NamedMessage struct{
	Username string  `json:"username"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

func GetChats(db *sql.DB, user User) (error, []Chat) {
	rows, err := db.Query(`SELECT * FROM chat_users WHERE user_id = $1;`, user.Id)
	if err != nil {
		return err, nil
	}
	var chats []Chat
	for rows.Next(){
		p := Chat{}
		var chat_id int
		var user_id int
		err := rows.Scan(&chat_id, &user_id)
		if err != nil{
			return err, nil
		}
		rows1, err := db.Query(`SELECT * FROM chats WHERE id = $1;`, chat_id)
		if err != nil {
			return err, nil
		}
		rows1.Next()
		err = rows1.Scan(&p.Name, &p.Id)
		if err != nil {
			return err, nil
		}
		chats = append(chats, p)
	}
	return err, chats
}

func GetChat(db *sql.DB, user User, chat_name string) (error, Chat) {
	err, chats := GetChats(db, user) // TODO: Make procedure to find only one chat
	if err != nil {
		return err, Chat{}
	}
	for _, chat := range chats {
		if chat_name == chat.Name {
			return nil, chat
		}
	}
	return errors.New("chat doesn't exist"), Chat{}
}

func GetMessages(db *sql.DB, chat Chat) (error, []Message) {
	rows, err := db.Query(`SELECT * FROM messages WHERE chat_id = $1;`, chat.Id)
	if err != nil {
		return err, nil
	}
	err, resultMessages := readMessagesData(rows)
	return err, resultMessages
}

func GetMessagesWithAuthors(db *sql.DB, chat Chat, count int, fromTime time.Time) (error, []NamedMessage) {
	rows, err := db.Query(`SELECT time, content, username FROM messages LEFT JOIN users ON messages.user_id = users.id WHERE messages.chat_id = $1 AND time < $2 ORDER BY time DESC LIMIT $3;`, chat.Id, fromTime, count)
	if err != nil {
		return err, nil
	}
	var messages []NamedMessage
	for rows.Next(){
		var message NamedMessage
		err := rows.Scan(&message.Time, &message.Content, &message.Username)
		if err != nil {
			return err, nil
		}
		messages = append(messages, message)
	}
	return err, messages
}

func SendMessage(db *sql.DB, user User, chat Chat, message string) error {
	_, err := db.Query(`INSERT INTO public.messages("time", chat_id, user_id, content) VALUES ($1, $2, $3, $4);`, time.Now(), chat.Id, user.Id, message)
	return err
}

func MakeChat(db *sql.DB, user User, chat string) (error, Chat) {
	rows, err := db.Query(`INSERT INTO public.chats(name) VALUES ($1) RETURNING name, id;`, chat)
	if err != nil {
		return err, Chat{} // TODO: if error occurs then delete chat
	}
	err, chats := readChatData(rows)
	if err != nil {
		return nil, Chat{}
	}
	if len(chats) != 1 {
		return errors.New("Can't create chat"), Chat{}
	}
	_, err = db.Query(`INSERT INTO public.chat_users(chat_id, user_id) VALUES ($1, $2);`, chats[0].Id, user.Id)
	if err != nil {
		return err, Chat{}
	}
	return err, chats[0]
}

func AddUserInChat(db *sql.DB, chat Chat, user2 string) error {
	rows, _ := db.Query(`SELECT * FROM users WHERE username = $1;`, user2)
	err, resultUsers := readUserData(rows)
	if err != nil {
		return err
	}
	if len(resultUsers) != 1 {
		return errors.New("Can't find second user")
	}
	_, err = db.Query(`INSERT INTO public.chat_users(chat_id, user_id) VALUES ($1, $2);`, chat.Id, resultUsers[0].Id)
	return err
}