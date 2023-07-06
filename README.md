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

## Plans

1. Make one-click button that will install all dependencies and construct working
project.
2. Make client-side message encryption (RAS, AES-256).
3. Give users ability to leave from chats, remove other users
   (with lower chat privileges)
4. Update front-end design
5. Update password hash algorithm and session control mechanism
6. Use Redis to store active connections
