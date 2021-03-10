[![CodeQL](https://github.com/netphantom/homeSerBot/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/netphantom/homeSerBot/actions/workflows/codeql-analysis.yml)
[![Generate release-artifacts](https://github.com/netphantom/homeSerBot/actions/workflows/build.yml/badge.svg)](https://github.com/netphantom/homeSerBot/actions/workflows/build.yml)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)


![](docs/img/logo.png)


Welcome to the HomeSerBot project page.
HomeSerBot is a tool that allows you to keep track on Telegram of the status, and the status code of your processes 
on your machine. 

It has been developed to give multiple users the possibility to register which processes they want to keep track and 
manage them through the integrated dashboard.

## Installation and usage

To install and run HomeSerBot you need two things:

1) Go installed on your machine

2) Access to a MySQL database

To run HomeSerBot, navigate in the folder `cmd/main` and run:`go build main.go`.
When the build is completed you can run it with `./main` appending all the needed parameters:

`-tApi` is your telegram API key

`-dbUserName` is the username to access to the DataBase 

`-dbPass` the user password  

`-dbIp` IP address of the database 

`-dbPort` the port to access 

`-dbName` the name of the DataBase

## Telegram commands

To be able to use the telegram bot you firstly need to register on the system.
This can be done  by typing the `/register` command on the Telegram bot. 
To validate the first user, you need to access to the web dashboard and validate your account.

Currently, few commands have been implemented:

`/pidList` provides with the process list registered on the service.

`/subscribe <item>` by appending the process ID from the previous command, you receive notifications when those are available from the system.

`/unsubscribe <item>` you can decide to stop receiving notifications from the given service.

`/subscriptions` you receive the list of all your process subscriptions.
