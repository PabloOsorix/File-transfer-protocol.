# File transfer protocol

Protocol that allow us join a channels to send and receive files, also you can send files directly to an user if it's registered in the hub.

You can execute the server with the following command in the main directory:

```
go run .
```

and the client in the client directory (/client/main.go) with the same command

```
cd client:
/client$
```
```
go run .
```

The client command interface show you the following commands to use to communicate with the server
| Command | Description |
| --- | --- |
| `REG` <@USERNAME> | Registers you in the server hub. |
| `JOIN` <#CHANNELNAME> | Use to join an existing channel in the server hub. |
| `LEAVE` <#CHANNELNAME> | Use it to leave from a channel you signed up for before. |
| `D_FILE` <nameOfFile> | Download a file from the server. |
| `S_FILE` <nameOfFile> -> <#CHANNELNAME> | Send a file to an user or channel. |
| `L_FILES`|  List all files existing on the server. |
| `USRS` |  List all users registerd on the server hub. |
| `CHNS` |  List all channels registered on the server hub. |
| `CHNS_FILE` <#CHANNELNAME> | List all files linked to a given channel. |
| `QUIT` | finish the execution of program, it means that exit from the server. |


Examples of use:

To send files you need to register
```
REG @NEW USER
OK, successfully register
```
To join a channel.
```
JOIN #GENERAL
User was registered in new channel
or user was registered in existing channel.
```

To leave from a channel
```
LEAVE #GENERAL
You leave of #GENERAL channel
```

You have to use L_FILES to know what files are in the server and use the name to download it.
```
D_FILE <name of the file iin the server> text.txt
```

NOTE: right arrow is mandatory.
```
S_FILE context.txt -> #GENERAL or (@name of user registered).

Receiving file...

Write directory of the file:
./

file was received and created
```

In the side of the server it print:
```
sending file completed 
File was received and shared in #GENERAL
```

To list all files from server
```
L_FILES

-client.txt -receipst.txt -testfile.txt -wakein.txt -wakeup.txt
```

To list all users on the server
```
USRS
@NEW_USER, @TAYLOR, @DUCK
```
To list all channels on the server.
```
CHNS
#GENERAL, #B3, #ANIME, #PROGRAMMING
```
To list all files registered on a certain channel.
```
CHNS_FILE #GENERAL
receipst.txt
```
Finish the execution
```
QUIT
```
