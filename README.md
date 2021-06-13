# A Restful API Boilerplate for Go

Clone this repo and start your API development.

## Start Application
- Clone this repository
- Start API with *serve*: ```go run main.go serve```

## Environment Variables
- Please refer .config.json.example
- You can modify .config.json the file or set in Env
- To use config file, duplicate .config.json.example and name it .config.json

## Database
- Integrated with PostgreSQL
- Integrated with migration package.
- Run migration: ```go run main.go migrate```
- Reset migration: ```go run main.go migrate --reset```

## Dev Guidelines
- In order to prevent circular dependency follow the below rules.
- ```app``` package can import ```service``` package. ```service``` can't import ```app```.
- ```service``` package can import ```database``` package. database can't import ```service```.
- Do not import ```database``` package directly from app.

## Integrated Packages
- HTTP Router - https://github.com/go-chi/chi
- PostgreSQL - https://github.com/jackc/pgx
- Logging - github.com/sirupsen/logrus
- JWT Token - github.com/dgrijalva/jwt-go
- Cli - https://github.com/spf13/cobra
- Configuration - https://github.com/spf13/viper
- Validation - https://github.com/go-playground/validator
- Redis - https://github.com/go-redis/redis