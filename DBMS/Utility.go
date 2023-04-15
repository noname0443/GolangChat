package DBMS

import "database/sql"

func readChatData(rows *sql.Rows) (error, []Chat) {
	var chats []Chat
	for rows.Next(){
		p := Chat{}
		err := rows.Scan(&p.Name, &p.Id)
		if err != nil{
			return err, nil
		}
		chats = append(chats, p)
	}
	return nil, chats
}

func readMessagesData(rows *sql.Rows) (error, []Message) {
	var messages []Message
	for rows.Next(){
		p := Message{}
		err := rows.Scan(&p.Time, &p.Content, &p.ChatId, &p.UserId)
		if err != nil{
			return err, nil
		}
		messages = append(messages, p)
	}
	return nil, messages
}