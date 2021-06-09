[![CodeQL](https://github.com/netphantom/homeSerBot/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/netphantom/homeSerBot/actions/workflows/codeql-analysis.yml)
[![Build Release](https://github.com/netphantom/homeSerBot/actions/workflows/release.yml/badge.svg)](https://github.com/netphantom/homeSerBot/actions/workflows/release.yml)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)


![](docs/img/logo.png)


Welcome to the HomeSerBot project page.
HomeSerBot is a tool that allows you to keep track on Telegram of the status, and the status code of your processes 
on your machine. 

It has been developed to give multiple users the possibility to register from which processes they want to receive
notifications and manage them through the integrated dashboard.

## Installation and usage

There are two ways to run HomeSerBot on your machine: download the pre-built packages or compile them on your machine, 
however as all the notifications are stored in a database, you must have access to any database source and provide 
HomeSerBot the correct credentials.

Currently, only the following database engines are supported: Mysql/MariaDb, Sqlite, Postgres, Sqlserver.

### Installation
If you prefer to use the pre-built packages, then go in the "release" section, then download the package corresponding 
to you system and extract the content in a folder.

Alternatively, you can clone this repo and compile HomeSerBot on your machine.
To do so, you need to install GO first, then run the `go build` command in the main folder.
Finally, you can run HomeSerBot with `./HomeSerBot`.

### Docker 
To run HomeSerBot using Docker, just use the following line, filling correctly all the required parameters.
Notice that the web dashboard of the container uses the port `4000`. 

```
docker run homeserbot ...
```

### Usage
The following parameters are required to correctly run HomeSerBot on your machine: 

`-tApi` is telegram API key you have created from the BotFather (https://t.me/botfather)

`-dbUserName` is the username to access to the DataBase 

`-dbPass` the user password  

`-dbIp` IP address of the database 

`-dbPort` the database port 

`-dbName` the name of the DataBase

You can access to the web dashboard by typing the address of the machine you have deployed (typically `localhost`) HomeSerBot, with the 
`4000` port

## Telegram commands

To be able to use the telegram bot you firstly need to register on the system.
This can be done  by typing the `/register` command on the Telegram bot. 
To validate the first user, you need to access to the web dashboard and validate your account.

Currently, few commands have been implemented:

`/pidList` provides with the process list registered on the service.

`/subscribe <item>` by appending the process ID from the previous command, you receive notifications when those are available from the system.

`/unsubscribe <item>` you can decide to stop receiving notifications from the given service.

`/subscriptions` you receive the list of all your process subscriptions.
