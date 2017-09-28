# mock-api
This is a mock adhaar API for the Hackocracy Challenge built with Golang. This web API functions as a verfiying mechanism the voting system 
application can trust and autheticate its users. This system uses fingerprints as a the biometric to verify the user's authetcity.
This API is modelled as close as possible to the real adhaar API in terms of functionality and in some areas the features are lenient
to better suit the needs of the client application, the voting system.

## Requirements
1. Golang - https://www.digitalocean.com/community/tutorials/how-to-install-go-1-6-on-ubuntu-14-04
2. MySQL - https://www.digitalocean.com/community/tutorials/how-to-install-mysql-on-ubuntu-16-04

## Packages used
1. github.com/gorilla/handlers
2. github.com/julienschmidt/httprouter
3. github.com/gorilla/csrf
4. github.com/jteeuwen/imghash
5. github.com/go-sql-driver/mysql

## Setup Instructions
1. Install `go` on your server and clone the reposistory under `< $GOPATH >/src/` using the link-https://github.com/Vijayprasanna13/mock-api/ .
2. `cd` to the `config/` inside the repository and run `cp env.example.go env.go` and replace the fields with your MySQL database credentials. 
3. From the project root, run `go run migrate.go`. This should create the required tables in the database.
4. From the project root, run `go run routes.go` to start the server, check localhost:8000/ for the app.

## App routes
1. `/user/register` - create a new user. This is synonymous to adding an user in the adhaar database so that the client app
can use the users adhaar id to verify his/her authenticity.
2. `/mock-api/user/auth` - autheticate the user using fingerprint image as biometric.

## Screenshots
![Adhaar Registration with Fingerprints](/screenshot.png)
