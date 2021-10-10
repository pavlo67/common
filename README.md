# common

## Ой, мамо, шо це?

Це приватний інструментарій для створення програмних сервісів на Golang. Він містить:

* засоби організації системи (з деякою простою версією Dependency Injection);
* інтерфейс http-сервера з підтримкою Swaggerʼа (v2) і шаблоном API-тестів;
* інтерфейс авторизації і декілька його імплементацій;
* деякі инші часто вживані елементи.

Все — виключно на смак автора. Як в анекдоті.

Пятачок (возбужденно-радостно): "Винни! Винни! Смотри, вот мой портрет!!!"
Винни-Пух (с сомнением): "А почему он весь такой на куски порезаный и пронумерованный?"
Пятачок (гордо): "Мясник рисовал. Он так видит."

## Запуск тесту для auth_stub

    go test -v github.com/pavlo67/common/common/auth/auth_stub

## Запуск демо-сервісу

    cp _environments/env.yaml_example local.yaml
    go run apps/demo/main.go

## Запуск тесту для auth_stub/auth_server_http

    # запустити сервер, як описано вище
    cp _environments/env.yaml_example test.yaml
    go test -v github.com/pavlo67/common/common/auth/auth_http


## Swagger

При дефолтних настройках доступний за адресою http://localhost:3001/backend/api-docs

Що туди передавати: дивіться запити в логах при запуску тесту для auth_stub/auth_server_http
