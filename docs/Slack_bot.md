# Створення застосунку в Slack та бота на мові GO

- Переходимо на Slack website
- Обираємо: `From Scratch`
- Вказуємо назву застосунку для Slackbot

    ![](../assets/img/image-0.png)  

- Натискаємо "Create App"

    ![](../assets/img/image-1.png)

- Обираємо який саме застосунок нам потрібно створити: `Bots`.
- Потрапляємо на сторінку довідкової інформації де натискаємо `Review Scopes to Add`  
- Додамо чотири головних області застосування боту (scopes):

    ![](../assets/img/image-2.png)

- Тепер можемо перейти до встановлення застосунку "Install the application"

    ![](../assets/img/image-3.png)

- Обираємо канал, який буде використовуватись ботом

    ![](../assets/img/image-4.png)

- Натискаємо `Allow` та отримуємо OAth Token та Webhook URL. 

    ![](../assets/img/image-5.png)

- Далі переходимо в створений WorkSpace Slack та запрошуємо бота на потрібний канал.

    ![](../assets/img/image-6.png)

- Цю дію можна зробити наступною командою з використанням слеш `/`
  
    `/invite @IlonaBot`

## Golang Setup and Installation

```sh
✗ go mod init github.com/k3ilona/publisher-slack-bot
go: creating new go.mod: module github.com/k3ilona/publisher-slack-bot
go: to add module requirements and sums:
        go mod tidy

✗ cobra-cli init          
Your Cobra application is ready at
/root/publisher-slack-bot

✗ cobra-cli add ilonabot
✗ cobra-cli add list
✗ cobra-cli add diff
✗ cobra-cli add promote
✗ cobra-cli add rollback

✗ go get -u github.com/slack-go/slack
go: added github.com/gorilla/websocket v1.5.1
go: added github.com/slack-go/slack v0.12.3
go: added golang.org/x/net v0.20.0

✗ go get -u github.com/joho/godotenv 
go: downloading github.com/joho/godotenv v1.5.1
go: added github.com/joho/godotenv v1.5.1
```

You have to run the below command for the functioning of the program

```sh
✗ go get
go: downloading github.com/google/go-github v17.0.0+incompatible
go: downloading golang.org/x/oauth2 v0.16.0
go: downloading github.com/google/go-github/v38 v38.1.0
go: downloading google.golang.org/appengine v1.6.7
go: downloading github.com/google/go-querystring v1.0.0
go: downloading golang.org/x/crypto v0.18.0
go: downloading google.golang.org/protobuf v1.31.0
go: added github.com/golang/protobuf v1.5.3
go: added github.com/google/go-github/v38 v38.1.0
go: added github.com/google/go-querystring v1.0.0
go: added golang.org/x/oauth2 v0.16.0
go: added google.golang.org/appengine v1.6.7
go: added google.golang.org/protobuf v1.31.0

✗ go run main.go

✗ gofmt -s -w ./                         

```


## Використані матеріали:
[How to Develop SlackBot Using Golang?](https://www.technource.com/blog/how-to-create-a-slackbot-using-golang/#What_Is_Slack_Bot)
[Develop a Slack-bot using Golang](https://programmingpercy.tech/blog/develop-a-slack-bot-using-golang/)  
[GO Documentation](https://go.dev/doc/)  
[How to Distribute Go Modules](https://www.digitalocean.com/community/tutorials/how-to-distribute-go-modules)  

---
← [Повернутись до змісту](../README.md)  
