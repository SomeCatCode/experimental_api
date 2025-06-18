#

## Build 
```bash
git clone ....
go install github.com/swaggo/swag/cmd/swag@latest
swag init
go build .
```

## Swagger
http://localhost:8080/swagger/index.html#/meters

## Enviroment File
PORT=8080
MONGO_URI=
MONGO_DB=

## ToDo
- Better Field Validation
- Better Error Responses
- Better Filtering/Search
- SideLoads
- AMPQ
- JWT
- Roles&Permissions
- Adjustable Log Level