# MIPT'25 budget planner

Приложение задеплоено по адресу:

```
https://79.174.83.151.sslip.io/
https://79.174.83.151.sslip.io/swagger
```

Запуск/перезапуск в контейнере:

```
docker compose up --build -d
```

Приложение запускается по адресу (доступен UI):  
http://127.0.0.1:8090/

Swagger:  
http://127.0.0.1:8090/swagger

Особенности:

1. При подгрузке отчетов следует помнить что используется кеширование и записи инвалидируются только по истечению TTL=30 сек. Кеш бюджета также имеет TTL=30 сек, но дополнительно инвалидируется при операциях записи.
2. При добавлении транзакции через API с isIncome=false (расход) amount следует указывать отрицательным.
3. Пример файла для импорта транзакций лежит в `example/transactions_for_import.csv`
4. При загрузке списка бюджетов и отчетов дополнительно возвращается флаг hitCache
5. Логика обработки кеширования лежит в `backend\ledger\internal\domain\budget\usecase\transaction\queries.go` и `backend\ledger\internal\domain\budget\usecase\budget\queries.go`

P.S.:  
В случае проблем обращаться https://t.me/k_milano
