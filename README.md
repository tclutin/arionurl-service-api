# ArionURL

Базовый url-shortener-api, созданный в целях обучения.

## Особенности

- **Сокращение URL**
- **Пользовательская установка времени жизни (seconds/hours)**
- **Пользовательская установка количества использований**

### API Endpoints

#### Редирект

- **URL:** `/<alias>`
- **Метод:** GET

#### Сокращение URL

- **URL:** `/alias`
- **Метод:** POST
- **Тело запроса:**

  ```json
  {
    "original_url": "http://example.com",
    "duration": "10m"
  }
  ```
  ```json
  {
    "original_url": "http://example.com",
    "duration": "1h"
  }
  ```
  ```json
  {
    "original_url": "http://example.com",
    "duration": "1h",
    "count_use": 5 
  }
  ```
  
- **Тело ответа:**
  ```json
  {
    "alias": "http://localhost:8080/<alias>"
  }
  ```
  #### запросы без указания "count_use" имеют бесконечное количество использований