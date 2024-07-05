# IM-System

IM-System is a simple instant messaging system based on the Go language, including server and client implementations. The server side uses the `net` package to handle network connections and supports basic public chat functionality, user list viewing, and username modification.

## Features

- **Public Chat**: All users connected to the server can send messages to each other in the same chat room.
- **User List**: Users can view the current online user list with the `/list` command.
- **Username Modification**: Users can change their username with the `/rename <new_name>` command.
- **User Join/Leave Notifications**: The server sends notifications to all users when a user joins or leaves.
- **User Exit**: Users can exit the chat room with the `/exit` command.

## Installation and Running

### Clone the Repository

```bash
git clone https://github.com/yourusername/im-System.git
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

After starting, the client will connect to the server and prompt the user to enter a username. Once the username is entered, the user can start chatting.

## Usage

### Client Commands

- Send a message: Directly enter the message and press enter to send it to the public chat.
- Exit: Enter /exit and press enter to exit the client.
- View user list: Enter /list and press enter to view the current online users.
- Change username: Enter /rename <new_name> and press enter to change your username.

## 实例

```
Enter '1' to start the Server or '2' to start the client: 2
Enter Server address (default 127.0.0.1:30001):
Enter your name: luola
> Welcome to the chat!
You have successfully changed your name to luola
> Hello everyone!
> list
Online users: luola
> rename john
You have successfully changed your name to john
> exit
```