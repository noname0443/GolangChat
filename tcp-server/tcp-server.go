package tcp_server

import (
	"ShadowChat/DBMS"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
)

func commandHandler(db *sql.DB, s string) (error, []byte) {
	strings.ReplaceAll(s, "", "")
	args := strings.SplitN(s, " ", 5)
	args[len(args) - 1] = strings.ReplaceAll(args[len(args) - 1], "\n", "")

	if args[0] == "GET" && args[1] == "MESS" {
		if len(args) != 4 {
			return errors.New("bad argument count"), nil
		}
		err, user := DBMS.GetUserBySession(db, args[2])
		if err != nil {
			return err, nil
		}
		err, chat := DBMS.GetChat(db, user, args[3])
		if err != nil {
			return err, nil
		}
		err, messages := DBMS.GetMessages(db, chat)
		if err != nil {
			return err, nil
		}

		var list [][]byte
		for _, message := range messages {
			list = append(list, []byte(message.Content))
		}
		if len(messages) == 0 {
			list = append(list, []byte("Nothing"))
		}
		return nil, bytes.Join(list, []byte(";"))
	} else if args[0] == "POST" && args[1] == "MESS" {
		if len(args) != 5 {
			return errors.New("bad argument count"), nil
		}
		err, user := DBMS.GetUserBySession(db, args[2])
		if err != nil {
			return err, nil // TODO: Remove SQLi
		}
		err, chat := DBMS.GetChat(db, user, args[3])
		if err != nil {
			return err, nil // TODO: Remove SQLi
		}
		err = DBMS.SendMessage(db, user, chat, args[4])
		if err != nil {
			return err, nil // TODO: Remove SQLi
		}
		return nil, []byte("Successful!")
	}
	return errors.New("Unkown command"), nil
}

func readFromConnection(conn net.Conn) (error, string) {
	all := make([]byte, 256)
	_, err := conn.Read(all)
	all = bytes.Trim(all, "\x00\n\r\t ")
	if err != nil {
		fmt.Println(err)
		return errors.New("Can't read data"), ""
	}
	return nil, string(all)
}

func handleConnection(db *sql.DB, conn net.Conn, chatConnections *sync.Map){
	defer conn.Close();

	err, data := readFromConnection(conn)
	if err != nil {
		return
	}
	stringValues := strings.Split(data, " ")
	sessionPass := stringValues[0]
	chat_name := stringValues[1]
	err, user := DBMS.GetUserBySession(db, sessionPass)
	if err != nil {
		return 
	}
	err, chat := DBMS.GetChat(db, user, chat_name)
	if err != nil {
		return
	}
	var connections []*net.Conn
	connectionsT, ok := chatConnections.LoadOrStore(chat.Id, []*net.Conn{&conn})
	if ok {
		connections = connectionsT.([]*net.Conn)
		connections = append(connections, &conn)
		chatConnections.Store(chat.Id, connections) // TODO: Remove Race Condition
	}

	err, allMessages := DBMS.GetMessagesWithAuthors(db, chat)
	if err != nil {
		return 
	}
	marshal, err := json.Marshal(allMessages)
	_, err = conn.Write(marshal)
	if err != nil {
		return
	}
	/*for _, message := range allMessages {
		marshal, err := json.Marshal(message)
		if err != nil {
			return
		}
		_, err = conn.Write(marshal)
		if err != nil {
			return 
		}
	}*/

	for {
		err, all := readFromConnection(conn)
		if err != nil {
			return
		}

		err = DBMS.SendMessage(db, user, chat, all)
		if err != nil {
			fmt.Println(err)
			return
		}
		//_, err = conn.Write(answer)
		//if err != nil {
		//	return
		//}
	}
}

func MessageHandler(db *sql.DB, chatConnections *sync.Map){
	fmt.Println("Launching server...")

	ln, _ := net.Listen("tcp", ":8081")

	for {
		conn, _ := ln.Accept()
		go handleConnection(db, conn, chatConnections)
	}
}
