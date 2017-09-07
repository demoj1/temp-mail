# API сервис для создание временной почты
Создает временную почту с помощью сервиса `https://temp-mail.ru/`.

## 1. Описание работы
Использует API сервиса `https://temp-mail.ru/` для создания временной почты. Временная почта принимает сообщение и хранит их в течение _10_ минут.

## 2. Список возможных действий
### 2.1 Создание временной почты:
`POST` запрос на адрес http://SERVICE_URL/create_mail
Параметры:
* email - логин почтового адреса;
* domain - домен, на данный момент можно использовать только один (@p33.org).

Пример верного запроса:
```sh
curl -X POST \
  http://localhost:8080/create_mail \
  -F email=dddddddd \
  -F domain=@p33.org
```

Пример ответа:
```json
{
    "Login": "dddddddd",
    "Domain": "@p33.org",
    "EmailHash": "de9020ecb99588a4a8057940b19d4168"
}
```

В ответе поле `EmailHash` будет использоваться в дальнейшним для получения списка сообщений.

---

Пример запроса с ошибкой:
```sh
curl -X POST \
  http://localhost:8080/create_mail \
  -F email=dddddddd \
  -F domain=empty
```

Пример ответа:
```json
{
    "error": "домен empty не может быть выбран.\nПожалуйста выберите один из следующих [@p33.org @binka.me @doanart.com]"
}
```

### 2.2 Получение списка возможных доменов:
> На данный момент работает только один домен `@p33.org`

`GET` запрос на адрес http://SERVICE_URL/available_domains

Пример запроса:
```sh
curl -X GET \
  http://localhost:8080/available_domains \
```

Пример ответа:
```json
[
    "@p33.org",
    "@binka.me",
    "@doanart.com"
]
```

### 2.3 Получение списка писем:
`POST`запрос на адрес http://SERVICE_URL/messages
Параметры:
* email_hash - хеш почтового адреса полученый на шаге _2.1_;
> Хеш так же можно сгенерировать в ручную, для этого требуется взять _md5_ сумму от значения `email+domain` md5(ddddd@p33.org)

Пример ответа если письма присутствуют:
```json
[
    {
        "mail_from": "Дмитрий Дубина <dmitryd.prog@gmail.com>",
        "mail_subject": "Test theme",
        "mail_text": "Test body\n",
        "mail_timestamp": 1504797272.423
    }
]
```

Пример ответа если писем нет:
```json
null
```