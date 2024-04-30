# Diskit Server

This is the server code for our CS-455 Networking Final Project

## Description

The Diskit Server is written using Golang's standard net/http package as well as the open-source httprouter that offers flexible but low-level control for routing endpoints. The server is effectively a REST api that performs CRUD opperations on the MongoDB backend.

## Getting Started

To run the server you will need Docker and Docker Compose. The server was built with the following versions of each.

```
Docker version 25.0.3, build 4debf41
Docker Compose version v2.22.0-desktop.2
```

If you wish to run without Docker, you can install the Go Compiler and have an instance of mongodb running on localhost. If both of these are true, then you can run the server by executing:

```
$ go run main.go
```

The program was made using Golang 1.18.1

### Executing program

* How to run the program
* Step-by-step bullets
```
$ docker build -t diskit-server .
$ docker compose build
$ docker compose up -d
```

The server should now be running and accessible at http://localhost:8080/ you can test this by curling http://localhost:8080/health. You should receive the response "healthy"

It's possible that you need to laod the database with values. I've prepared a CSV and python script for this reason. In the "add-courses" directory is a script called initialize-db.py.

```
$ python initialize-db.py
```

Will load the database with courses from the CSV file. Upon a success, you will receive the monogodb objects output as json. If you don't receive this, make sure the mongodb url is correct. 
