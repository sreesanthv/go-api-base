# A Restful API Boilerplate for Go

Clone this repo and start your API development.

## Start Application
### Configuration
- Please refer .config.json.example
- You can modify .config.json the file or set in Env
- To use config file, duplicate .config.json.example and name it .config.json
- Please correct ```db_addr, redis_addr``` to run in local.
### Docker
- Clone this repository
- Build Docker Images ```docker-compose build api```
- Run Docker Container ```docker-compose up api```
### Local
- Start dependencies ```docker-compose up -d db redis```
- Start API with *serve*: ```go run main.go serve```

### Dev Dependencies
- Install https://github.com/codegangsta/gin for live reloading.
Run ```gin serve```
- Start Redis UI by ```docker-compose up -d redis-commander```
Open http://localhost:8081 in browser.

## Database
- Integrated with PostgreSQL
- Integrated with migration package.
- Run migration: ```go run main.go migrate```
- Reset migration: ```go run main.go migrate --reset```

## Dev Guidelines
- In order to prevent circular dependency follow the below rules.
- ```api``` package can import ```service``` package. ```service``` can't import ```api```.
- ```service``` package can import ```database``` package. database can't import ```service```.
- Do not import ```database``` package directly from api.
- Preferred import flow: ```api -> services -> database```

## Integrated Packages
- HTTP Router - https://github.com/go-chi/chi
- PostgreSQL - https://github.com/jackc/pgx
- Logging - github.com/sirupsen/logrus
- JWT Token - github.com/dgrijalva/jwt-go
- Cli - https://github.com/spf13/cobra
- Configuration - https://github.com/spf13/viper
- Validation - https://github.com/go-playground/validator
- Redis - https://github.com/go-redis/redis