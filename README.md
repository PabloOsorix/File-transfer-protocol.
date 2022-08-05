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

