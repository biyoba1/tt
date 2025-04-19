### Выдача JWT токенов

**URL:** `http://localhost:8002/auth/sign-up`  
**Метод:** `POST`  

Этот эндпоинт позволяет получить Access & Refresh токены по GUID пользователя. GUID передается в теле запроса.

**Тело запроса (JSON):**
```json
{
  "guid": "qweqwe"
}
```

**Пример запроса:**
```bash
curl -X POST http://localhost:8002/auth/sign-up \
     -H "Content-Type: application/json" \
     -d '{"guid": "qweqwe"}'
```

**Пример ответа:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "WlhsS2FHSkhZMmxQYVVwSlZYcEpNVTVwU1hOSmJsSTFZME5KTmtscmNGaFdRMG81..."
}
```

**Коды ответа:**  
- `200 OK` - Токены успешно выданы.  
- `400 Bad Request` - Некорректные данные в запросе.  
- `404 Not Found` - Пользователь с указанным GUID не найден.  
- `500 Internal Server Error` - Произошла ошибка на сервере.  

---

### Рефреш JWT токенов

**URL:** `http://localhost:8002/auth/refresh`  
**Метод:** `POST`  

Этот эндпоинт позволяет получить новый Access & Refresh токены по валидному Refresh-токену пользователя. Формат передачи токена - Base64.

**Тело запроса (JSON):**
```json
{
  "refresh_token": "WlhsS2FHSkhZMmxQYVVwSlZYcEpNVTVwU1hOSmJsSTFZME5KTmtscmNGaFdRMG81TG1WNVNteGxTRUZwVDJwRk0wNUVXWHBOVkUwMFRWUkpjMGx0YkhkSmFtOXBWM3B2TmsxV01EWk9WRTE1VG5wRmFVeERTbTVrVjJ4clNXcHZhV05ZWkd4aldHUnNTVzR3TGxsRFpWTTJhbk5QYkVacFZWQkNXRWt3Y2xOMGVsOVFhRVZwY3pCUmJITlFUMjlpWmkweFRteFlRbk09"
}
```

**Пример запроса:**
```bash
curl -X POST http://localhost:8002/auth/refresh \
     -H "Content-Type: application/json" \
     -d '{"refresh_token": "WlhsS2FHSkhZMmxQYVVwSlZYcEpNVTVwU1hOSmJsSTFZME5KTmtscmNGaFdRMG81TG1WNVNteGxTRUZwVDJwRk0wNUVXWHBOVkUwMFRWUkpjMGx0YkhkSmFtOXBWM3B2TmsxV01EWk9WRTE1VG5wRmFVeERTbTVrVjJ4clNXcHZhV05ZWkd4aldHUnNTVzR3TGxsRFpWTTJhbk5QYkVacFZWQkNXRWt3Y2xOMGVsOVFhRVZwY3pCUmJITlFUMjlpWmkweFRteFlRbk09"}'
```

**Пример ответа:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "WlhsS2FHSkhZMmxQYVVwSlZYcEpNVTVwU1hOSmJsSTFZME5KTmtscmNGaFdRMG81..."
}
```

**Коды ответа:**  
- `200/201 OK` - Новые токены успешно выданы.  
- `400 Bad Request` - Некорректные данные в запросе.  
- `500 Internal Server Error` - Произошла ошибка на сервере.  

---
