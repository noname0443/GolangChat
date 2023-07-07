# GolangChat
## Description
GolangChat is real-time message-service on golang language which created for educational purposes. It's provide users to
exchange message, save message history and makes new chats for your new contacts.
It's not a social network just a chat.

Now users can
1. Make simple account (Register/Login)
2. Make chat to communicate
3. Invite other users by their name (count of users in chat isn't limited)
4. Send messages and get messages in real time
## Installation
It's very important to set up all environment variables before launch program. You can use prepared scripts after reboot or write it to .bashrc/windows variables.
### Ubuntu
```bash
sudo apt-get install postgresql
sudo apt-get install go
mkdir GolangChat
cd GolangChat
git clone https://github.com/noname0443/GolangChat.git
psql -U myuser -f FirstRun.SQL
source script.sh
go build .
```
### Windows
1. Install [golang compiler](https://go.dev/dl/) on your device.
2. Install [PostgreSQL](https://www.postgresql.org/download/).
3. Copy project from github
4. Launch ```psql.exe -U myuser -f FirstRun.sql```
5. ./script.bat
6. Compile it with ```go build .```

## Security
Positive:
1. All SQL queries uses prepared statements and input to it undergoing sanitizing.
2. Client side get clear from XSS messages and chat names from server.
3. All files that the client can interact with are in the folder templates.

Negative:
1. User's DBMS records have slow MD5 hash algorithm
2. Messages which users send to server isn't encrypted.

## To-Do
- Database session clearer
- Make client-side message encryption (RAS, AES-256).
- Give users ability to leave from chats, remove other users, delete chats
- Update front-end design
- Make one-click button that will install all dependencies and construct working project. (Docker)
- Websocket connection pool (Redis)
- Add message queue (RabbitMQ)
- Separate SQL queries from golang code to separate file (prepared procedures)
- Get free ssl from https://letsencrypt.org/
- Find html generator to improve front-end of chat
- Invites in real-time with accept/deny
- Use graphql with GIN
- Refactor
	- Clear sockets.go. It's too big and hard to read
	- Use SHA256 instead of MD5
	- Get database variables from environment instead of file
	- Use SQL interface
	- Use SQLX to make annotation for read data faster and in less code
	- Add more comments to generate full godoc
	- Use tests
- Email confirm
- Merge all separated connection processes in one with "pool of workers" and just use redis or RabbitMQ for messages. To store active sockets use map. Clear it's from closed sockets for every minute (or make it with events)
