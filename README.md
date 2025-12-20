# MIPT'25 budget planner

Запуск/перезапуск в контейнере:

```
docker compose up --build -d
```

Приложение запускается по адресу:  
http://127.0.0.1:8090/

Swagger:  
http://127.0.0.1:8090/swagger

Особенности:
1) При добавлении транзакции с isIncome=false (расход) amount следует указывать отрицательным.  
2) Пример файла для импорта транзакций лежит в `example/transactions_for_import.csv`  
3) При загрузке списка бюджетов и отчетов дополнительно возвращается флаг hitCache
4) Логика обработки кеширования лежит в `backend\ledger\internal\domain\budget\usecase\transaction\queries.go` и `backend\ledger\internal\domain\budget\usecase\budget\queries.go`

P.S.:  
В случае проблем обращаться https://t.me/k_milano
