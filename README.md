# url-shortener
### Для запуска:
Запустите с помощью Docker Compose:
```bash
docker-compose up -d
```

### API:
Post запрос: /url-shorter
Пример body: 
```json
{
    "original-url": "http://asfasfs"
}
```
Get запрос: /url-shorter/{alias}

### Порт по умолчанию: 8090
