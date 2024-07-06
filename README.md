# IM-System

IM-System is a simple instant messaging system based on the Go language, including server and client implementations.
The server side uses the `net` package to handle network connections and supports basic public chat functionality, user
list viewing, and username modification.

## Features

-Users can connect to the server for public chat.

-Support users to modify their usernames.

-Automatically detect the user's connection status and handle it automatically when the connection is disconnected.

-Provide basic chat commands such as exiting and listing online users.

## Installation and Running

### Clone the Repository

```bash
git clone https://github.com/luolayo/im-System.git
cd im-System
```

### Build the Project

```bash
go build -o im-system main.go
```

## Running the Server

Start the server with the following command:

```bash
./im-system
```

The server will start and listen on the default port 127.0.0.1:30001.

## Running the Client

Start the client with the following command:

```bash
./im-system
```

After starting, the client will connect to the server and prompt the user to enter a username. Once the username is
entered, the user can start chatting.

## Usage

### Server

```bash
2024/07/04 10:00:00 INFO Starting server on 127.0.0.1:8080
```

### Client

```bash
Enter server IP: 127.0.0.1
Enter server port: 8080

Menu:
1. Exit
2. Rename User
3. Enter Public Chat Mode
Enter your choice: 
```

### Public chat mode

```bash
Entering public chat mode. Type '/exit' to return to the menu.
> Hello everyone!
> /exit
```

## TODO
- [x] Add private chat function
- [ ] Add user registration and login functions
- [ ] Improve the log module to support different levels of log output
- [ ] Add more commands and functions