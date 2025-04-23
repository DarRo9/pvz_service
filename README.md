
### Запуск

Запуск сервиса

```git clone https://github.com/DarRo9/pvz_service
cd pvz_service
make create_env
go mod tidy
make run
```

Запуск unit-тестов
```make unit_test```

Запуск integration-теста (должен быть запущен сервис)
```make integration_test```


### Структура проекта
```
├── api
│   ├── openapi            # спецификация HTTP API
│   └── proto              # спецификация GRPC API
├── cmd
│   └── server             # Точка входа
├── config                 # Конфигурационный файл
├── integration_test
│   └── utils              # Интеграционный тест
├── internal
│   ├── db                 # Подключение к базе данных
│   ├── grpc               # GRPC сервер и хендлеры 
│   ├── handler            # HTTP сервер и хендлеры
│   ├── metrics            # Метрики prometheus
│   ├── middleware         # Middlewares для авторизации и метрик 
│   ├── repository         # Логика работы с базой данных
│   ├── service            # Бизнес логика
│   └── utils              # Утилиты для JWT токена
└── migrations
```
### Что было реализовано
- API
- Пользовательская авторизация по методам /register и /login 
- Unit и integration тесты
- gRPC метод для получения всех ПВЗ
- Метрики prometheus
- Логирование
- Кодогенерация DTO endpoint'ов
