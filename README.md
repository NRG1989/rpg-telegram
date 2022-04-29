# Go-Aut-Registration-User-Telegram

## Short description
```
Telegram-bot to recieve codes for registration and change password
```

## Dependencies
You need to find bot in telegram, by yourself
```
@RPGn_bot
```

## Initialization

## Config
```
./config.json
```
You can also specify your own config file by run app with flag `-config`

## Migrations
Migration creation:
1) Choose migration directory
2) Create migration:
```shell script
$ goose create name_of_migration sql
```
3) In created file write your migration
4) Apply migration:
```shell script
 $ goose postgres "postgres://rpguser:rpgpass@10.10.15.90:5432/RpgDB?sslmode=disable" up

```
5) Or cancel it:
```shell script
$ goose postgres "postgres://rpguser:rpgpass@10.10.15.90:5432/RpgDB?sslmode=disable" down

```