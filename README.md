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
| `REG` <@USERNAME> | List all *new or modified* files |
| `JOIN` <#CHANNELNAME> | Show file differences that **haven't been** staged |
| `LEAVE` <#CHANNELNAME> |
| `D_FILE` <nameOfFile> |
| `S_FILE` <nameOfFile> -> <#CHANNELNAME> |
| `L_FILES`|
| `USRS` |
| `CHNS` |
| `CHNS_FILE` <#CHANNELNAME> |

