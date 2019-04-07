Variables
  - *proxy_url* - You can use any proxy tool if api telegram is not available from your country(do not fill out to keep default)
  - *token* - Your telegram bot token - Required
  - *chat_id* - Chat ID, which will be sent to the bot notifications - Required

Example pipeline
```yml
kind: pipeline
name: CI/CD mf

workspace:
  base: /go
  path: mod/github.com/user/service

steps:
  - name: tests
    image: golang:latest
    commands:
      - go test -v --cover ./...

  - name: linters
    image: golang:latest
    commands:
      - go get -u golang.org/x/lint/golint
      - golint ./...

  - name: telegram notify
    image: pocoz/drone-tg
    settings:
      proxy_url: "https://your.proxy.url"
      token:
        from_secret: telegram_token
      chat_id:
        from_secret: telegram_chat_id
    when:
      status: [ success, failure ]
```
Build image:

    docker build -t pocoz/drone-tg .

Push image:

    docker push pocoz/drone-tg