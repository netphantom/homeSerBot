#HomeSerBot

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
