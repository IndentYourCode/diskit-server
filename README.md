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


### Executing program

* How to run the program
* Step-by-step bullets
```
docker build -t diskit-server .
docker compose build
docker compose up -d
```

The server should now be running and accessible at http://localhost:8080/ you can test this by curling http://localhost:8080/health. You should receive the response "healthy"