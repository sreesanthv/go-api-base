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