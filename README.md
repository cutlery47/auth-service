# auth-service
a tiny auth service written in go

---

Тестовое задание на позицию junior Golang разработчика в MEDODS

## Инструкция запуску / закрытию приложения
1) Клонируем репозиторий в произвольную директорию

```
git clone git@github.com:cutlery47/auth-service.git
cd auth-service
```

2) Настраиваем переменные окружения

В корневой директории проекта лежит файл example.env, следующего содержания:

```
ACCESS_DURATION             =30m                (время валидности access-токена)
REFRESH_DURATION            =24h                (время валидности refresh-токена)
HASHING_SECRET              =                   (секрет для генерации токенов)
HASHING_COST                =10                 (bcrypt cost параметр)

SMTP_USERNAME               =                   (почта для доступа к SMTP-серверу)
SMTP_PASSWORD               =                   (пароль для доступа к SMTP-серверу)
SMTP_RECEIVER               =                   (почта получателя сообщений с SMTP-сервера)
SMTP_HOSTNAME               =smtp.gmail.com     (адрес SMTP-сервера)
SMTP_PORT                   =587                (порт SMTP-сервера)

POSTGRES_USER               =postgres           (имя POSTGRES пользователя внутри контейнера)
POSTGRES_PASSWORD           =12345              (пароль POSTGRES пользователя внутри контейнера)
POSTGRES_HOST               =postgres           (адрес POSTGRES-сервера внутри контейнера)
POSTGRES_PORT               =5432               (порт POSTGRES-сервера внутри контейнера)
POSTGRES_DB                 =auth               (имя POSTGRES базы данных внутри контейнера)
POSTGRES_TIMEOUT            =5s                 (тайм-аут на подключение к БД)
POSTGRES_MIGRATIONS         =./migrations       (директория с миграциями для БД)

LOGGER_DIR                  =./logs             (директория для логов на локальной машине)
LOGGER_INFO_PATH            =./logs/info.log    (директория для логов информации внутри контейнера)
LOGGER_ERROR_PATH           =./logs/error.log   (директория для логов ошибок внутри контейнера)

SERVER_INTERFACE            =0.0.0.0            (сетевой интерфейс, на котором слушает приложение)
SERVER_PORT                 =8080               (порт, на котором слушает приложение)
SERVER_READ_TIMEOUT         =3s                 (тайм-аут на чтение сервера)
SERVER_WRITE_TIMEOUT        =3s                 (тайм-аут на запись сервера)
SERVER_SHUTDOWN_TIMEOUT     =3s                 (тайм-аут на закрытие сервера)
```

Настоятельно не рекоммендуется изменять уже заданные значения :)

Мануал для настройки SMTP лежит в файте HOWTO-SMTP.md

После заполнения пустых значений следует изменить название файла с example.env на .env

3) Запускаем приложение в Docker-контейнере

`make up_build`

4) Выходим из приложения

`make down`

5) Чистим Docker volumes

`make clean`

## Инструкция по работе с приложением

После запуска приложения, Вы можете проверить статус работы сервиса:

`curl http://localhost:{ВАШ_ПОРТ}/ping`

При корректной работе возвращается пустой ответ

Если приложение работает корректно, Вы можете открыть Swagger-документацию внутри браузера. Для этого нужно открыть страницу по адресу 

`http://localhost:{ВАШ_ПОРТ}/swagger/`

Приложение обрабатывает два эндпойнта:

Первый - создает пару access и refresh токенов для пользователя с заданным guid:

`http://localhost:{ВАШ_ПОРТ}/api/v1/auth?guid={ВАШ_GUID}`

(Для создания GUID, Вы можете воспользоваться любым онлайн-генератором, например https://www.uuidgenerator.net/guid)

Второй - рефрешит пару refresh и access токенов

`http://localhost:{ВАШ_ПОРТ}/api/v1/refresh?guid={ВАШ_GUID}&refresh={ВАШ_REFRESH}`
