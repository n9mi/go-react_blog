# Go Echo Notes CRUD
Simple CRUD with JWT auth using golang Echo and javascript React

## **Packages used**
- github.com/spf13/viper
- github.com/labstack/echo/v4
- github.com/go-playground/validator/v10
- gorm.io/driver/postgres
- gorm.io/gorm 
- github.com/stretchr/testify
- github.com/redis/go-redis/v9
- github.com/golang-jwt/jwt

## **How to run the app**
```
docker compose build
docker compose up
```
If you're using linux and have running postgresql in your machine, 'db' port will conflict with the postgresql in your machine. Use this command to stop postgresql in your machine
```
sudo systemctl stop postgresql
```

## **Run the migration**
Migration will automatically run when the server starts, and resetting the migration on another run.
![image](https://github.com/naomigrain/httprouter-crud-notes/assets/113373725/2488a53e-3bf0-421c-be45-4faa2c87d66f)

## **Structure**
Based on repository pattern, this project use:
- Repository layer: For accessing db in the behalf of project to store/update/delete data
- Usecase layer: Contains set of logic/action needed to process data/orchestrate those data
- Entity: Contains set of database atribute
- Model: Contains set of data that will be parsed or send as request or response
- Controller layer: Acts to mapping users input/request and presented it back to user as relevant responses

## TODO
- React frontend :")

## **API Endpoints**


