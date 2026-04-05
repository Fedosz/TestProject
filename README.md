# rates_project

gRPC сервис для получения курса USDT с биржи Grinex.

Сервис:
- ходит в API Grinex по HTTP через `resty`
- считает `ask` и `bid` по заданному методу
- сохраняет каждый полученный курс в PostgreSQL
- отдаёт результат через gRPC
- поднимает healthcheck
- поддерживает graceful shutdown
- экспортирует метрики Prometheus
- трассирует запросы через OpenTelemetry

---

## Что умеет сервис

### GetRates
Метод получает стакан с биржи Grinex и считает итоговые значения `ask` и `bid`.

Поддерживаются два способа расчёта:

- `topN` — взять значение с позиции `N`
- `avgNM` — посчитать среднее по диапазону `[N; M]`

`ask` считается по массиву `asks`, `bid` — по массиву `bids`.

Каждый успешный вызов `GetRates` сохраняется в БД вместе со временем получения курса.

### Health check
Метод проверяет доступность сервиса.  
Сейчас healthcheck опирается на доступность PostgreSQL.

## Запуск


```bash
make build 

docker-compose up -d

docker-compose run --rm app ./app
```

---

## Проверка сервиса

### Health check

```bash
grpcurl -plaintext -import-path api/proto -proto rates.proto \
  localhost:50051 rates.v1.HealthService/Check
```

---

### GetRates с методом topN

```bash
grpcurl -plaintext -import-path api/proto -proto rates.proto \
  -d '{
    "ask": {
      "method": "CALC_METHOD_TOP_N",
      "n": 1
    },
    "bid": {
      "method": "CALC_METHOD_TOP_N",
      "n": 1
    }
  }' \
  localhost:50051 rates.v1.RatesService/GetRates
```

Пример ответа:

```json
{
  "ask": 80.82,
  "bid": 80.71,
  "timestampUnix": 1775419112
}
```

---

### GetRates с методом avgNM

```bash
grpcurl -plaintext -import-path api/proto -proto rates.proto \
  -d '{
    "ask": {
      "method": "CALC_METHOD_AVG_NM",
      "n": 1,
      "m": 3
    },
    "bid": {
      "method": "CALC_METHOD_AVG_NM",
      "n": 1,
      "m": 3
    }
  }' \
  localhost:50051 rates.v1.RatesService/GetRates
```

### Метрики

```bash
curl http://localhost:2112/metrics
```

---

Запуск тестов:

```bash
make test
```

Запуск линтера:

```bash
make lint
```

Генерация контрактов:

```bash
make proto
```

Генерация моков:

```bash
make mocks
```
