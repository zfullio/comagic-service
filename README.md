# Comagic

* gRPC сервер для сбора данных Comagic и записи в BigQuery
* скрипт (клиент) сбора данных по расписанию

## Реализованные методы

* PushCallsToBQ - Сбор звонков за указанный период
* PushOfflineMessagesToBQ - Сбор заявок за указанный период

Более подробно ознакомиться с протоколами можно в файле `api/grpc/comagic.proto`

## Алгоритм работы

* Получение данных по dataapi Comagic
* Генерация CSV файла в temp директории
* Проверка наличия / создание таблицы в BigQuery
* Сохранение файла в CloudStorage Bucket (Bucket должен быть уже создан)
* Удаление данных из таблицы BigQuery за период указаннй в настройках сбора
* Запись данных в таблицу BigQuery из CloudStorage

## Makefile

`gen_go` - Конверт-ия proto файлов в код Go на основании `api/grpc/comagic.proto`

```bash
make gen
```

`gen_python` - Конверт-ия proto файлов в исходный код Python на основании `api/grpc/comagic.proto`

```bash
make gen_python
```

`build_server` - Компиляция исполняемого файла (/bin/server_app)

```bash
make build_server
```

## Опции

`-f`: Выбор конфигурационного файла.

`--env`: Использование в качестве конфигурации переменных окружения.

`--trace`: Сделать уровень логирования по умолчанию равным "Trace". Так же добавит: pid процесса, функция вызова

## Конфигурация (Сервер)

На выбор несколько вариантов настройки:

* По умолчанию в качестве настроек используется config.yml в папке с исполняемым файлом
* С помощью флага `-f` укажите путь на конфигурационный файл
* С помощью флага `--env` в качестве настроек используются переменные окружения

### Файл конфигурации

По умолчанию для настроек используется config.yml в папке с исполняемым файлом.  
Для использования альтернативного файла используйте флаг `-f`

```yaml
# Пример конфигурационного файла
keys_dir: "/path/to/keys" // Путь к папке с сервисными ключами

grpc:
  ip: "0.0.0.0"           // Host
  port: 50051             // Порт, который будет прослушивать сервис

comagic:
  version: "2.0"          // Версия DataApi Comagic (текущая 2.0)

tg:
  token: 'TG Token'       // Токен для telegram бота
  chat: 0000000000        // ID чата в который будут отправляться уведомления
  is_enabled: false       // Статус уведомлений
```

### Использование переменных окружения

Для использования переменных окружения используйте флаг  `--env`

| Переменная         | Описание                                         |
|--------------------|--------------------------------------------------|
| `GRPC_IP`          | Host                                             |
| `GRPC_PORT`        | Порт, который будет прослушивать сервис          | 
| `TG_TOKEN`         | Токен для telegram бота                          |
| `TG_CHAT`          | ID чата в который будут отправляться уведомления |
| `TG_ENABLED`       | Статус уведомлений                               |
| `COMAGIC_VERSION ` | Версия DataApi Comagic (текущая 2.0)             |
| `KEYS_DIR `        | Путь к папке с сервисными ключами                |

## Конфигурация (Клиент)

* По умолчанию в качестве настроек используется schedule_config.yml в папке с исполняемым файлом
* С помощью флага `-f` укажите путь на конфигурационный файл

### Файл конфигурации

По умолчанию для настроек используется config.yml в папке с исполняемым файлом.  
Для использования альтернативного файла используйте флаг `-f`

```yaml
# Пример конфигурационного файла
time: "07:47" // Время ежедневного запуска

grpc:
  ip: "0.0.0.0"           // Host
  port: 50051             // Порт, который будет прослушивать сервис

reports:
  - object: "название клиента 1"
    comagic_token: "xxxxxxxxxx"
    google_service_key: "service_key.json"
    project_id: 'bq-project-id'
    dataset_id: 'bq-dataset-id'
    bucket_name: 'storage-bucket-name'
    offline_message_table: 'bq_offline_messages_table'
    calls_table: 'bq_calls_comagic_table'

  - object: "название клиента 2"
    comagic_token: "xxxxxxxxxx"
    google_service_key: "service_key.json"
    project_id: 'bq-project-id'
    dataset_id: 'bq-dataset-id'
    bucket_name: 'storage-bucket-name'
    offline_message_table: 'bq_offline_messages_table'
    calls_table: 'bq_calls_comagic_table'
```