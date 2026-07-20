# ProxyGuard
## Цель сервиса
Снизить количество отправляемых спам-СМС и сэкономить бюджет

## Архитектура сервиса
Сервис установлен между внешним nginx на сервере ngc-ingress-ext-01 и бэкенд серверами сайтов и МП. Данный сервис является прозрачным прокси, который, если не сработали фильтры, отправляет запросы и получает ответы от бэкенд серверов. Если фильтры сработали, то сервис не отправляет запрос на отправку СМС на бэкенд сервера, и в качестве ответа подставляет свой ответ. Ранее ответы подставлялись успешные с полной аналогией ответов от бэков, чтобы злоумышленник не понял срабатывает ли защита или нет; но при этом клиенты тоже не могли понять ждать им СМС или нет, т.к. ошибки не возникало. На данный момент, при срабатывании фильтра, сервис возвращает ошибку 429, а клиенты видят ошибку на фронте.

## Фильтры
Проверки производятся сверху-вниз в указанном порядке
### Сайты
1. Ratelimit: 3 СМС за 100 секунд. В этом фильтре блокировку производим по телефону.
1. Проверка по сессии и телефону. Сессию берём из куки `shop_session`.
   Логика по данному фильтру следующая:
   * если один и тот же номер телефона с одной сессии, то применяем стандартный ratelimit - 3 СМС за 100 секунд
   * если сессия одна, а телефонов несколько, то применяем "тяжёлый" ratelimit - 3 СМС за сутки(86400 секунд)
   В этом фильтре блокировку производим по сессии
1. Блокировка по куки, если нет хотя бы одно из указанных кук:
   * `_ga`
   * `rrpvid`
   * `rcuid`
   * `_ym_uid`
   * `advcake_track_url`
   * `advcake_track_id`

### Мобильное приложение
Изначально устанавливается счётчик на любые GET-запросы. Для запроса фиксируется ip и id устройства клиента из хедера `X-Device-Id`.
Это сделано, чтобы отличить реального клиента от робота. Клиент, при загрузке приложения, обязательно должен получить какие-то данные методом GET - картинки, информацию о товарах и тд.. Счётчик по GET-запросам хранится 30 минут(1800 секунд).
Для мобильного приложения фильтр устанавливается на два URL:
* `/api/v8/manzana/registration` - регистрация клиента
* `/api/v8/ecom-auth/login-sms-prestep` - отправка кода в СМС при авторизации

### Регистрация клиента
Дополнительно к проверке количества GET-запросов при регистрации производится проерка по Email и телефону регистрации. Если указан Email из домена gmail.com, сохраняем ip-адрес и телефон, указанные в запросе
1. Блокируем запрос клиента, если количество предыдущих GET-запросов меньше двух. Блокировка производится по ключу `ip+deviceId`
1. Блокируем запрос клиента, если в запросе Email из домена gmail.com
1. Блокируем запросы клиента, если для одного телефона указаны разные Email, фильтр срабатывает если количество запросв с разными Email больше одного. Если емейлы одинаковые, то фильтр срабатывает, если количество запросов больше двух, т.о. мы даём клиенту одну попытку, в случае, если данные ввыедены некорректно в первом запросе
   Блокировка производится по телефону на 30 минут(1800 секунд)

### Отправка кода в СМС при авторизации
1. Блокируем запрос, если сработал один из фильтров при регистрации
1. Блокируем запрос, если было больше двух попыток регистрации для одного телефона за 1 минуту
1. Блокируем запрос, если для одного `X-Device-Id` было больше одной попытки на разные номера телефонов
   Блокировка производится и по телефону и по `X-Device-Id`, время блокировки - 30 минут

## Диагностика
Диагностика производится по лог-файлам. Есть три типа файлов:
  *all - все входящие запросы
  *handle - причина блокировки, по номеру телефона или ip
  *accepted - все запросы, которые были пропущены

Лог-файлы для МП разделены на бренды в имени.

Примеры команд:
1. Поиск всех входящих запросов по номеру телефона:
   ```bash
   cat logs/site/all.log|grep 9999999999
   ```
1. Сколько СМС было отправлено за определённый час:
   ```bash
   cat logs/site/accepted.log|grep "2026-07-08T14"|wc -l
   ```
Здесь мы можем увидеть сколько СМС отправлено 8 июля за 14 до 15 часов.
Норма отправленных СМС: 150-200 в час для бренда Озерки, если значение превышено, значит фильтры не срабатывают и необходимо корректировка.

# Стэк
* Go
* Redis

# Инфраструкутра
Сервер: ngc-proxy-guard-01.dev.erkapharm.ru(10.22.3.71), в названии dev, но используется как боевой.
Открытые порты:
 * 8090/tcp - порт пдля сайтов
 * 8080-8083/tcp - порты для мобильных приложений
 * 8070/tcp - порт для зеркального трафика, для подсчёта всех GET-запросов

# Развитие сервиса
* админка с поиском заблокированных телефонов
* черные/белые списки
* мониторинг эффективности фильтров, количества заблокированных/пропущенных запросов через exporter в Grafana
* логирование в общий Elastic

# Примеры конфигурационных файлов nginx
* конфигурация для сайта:
```nginx
location /api/customer/auth-sms {
   proxy_pass http://proxy-guard-01.dev.erkapharm.ru:8090;
   proxy_set_header X-Real-IP $remote_addr;
   proxy_set_header X-Forwarded-for $remote_addr;
   proxy_set_header X-Forwarded-Proto $scheme;
   proxy_set_header Host $host;
   proxy_set_header Connection close;
   proxy_pass_header Content-Type;
   proxy_pass_header Content-Disposition;
   proxy_pass_header Content-Length;
   client_max_body_size 250m;
   client_body_buffer_size 4m;
   fastcgi_buffers 16 16k;
   fastcgi_buffer_size 32k;
}

location = /mirror {
   internal;
   proxy_set_header X-Real-IP $remote_addr;
   proxy_set_header X-Forwarded-for $remote_addr;
   proxy_set_header X-Orig-Path $request_uri;
   proxy_pass http://proxy-guard-01.dev.erkapharm.ru:7070;
}
```
# Примерs запросов на регистрацию и авторизацию
* Отправка кода для сайта

  **POST /api/customer/auth-sms**
  ```json
  {
     "Cookie": [
      "PHPSESSID=kfejut6gct61j0b6b3bh6e4tp6; storeId=7329; __ddg1_=h4X2kys3Y2RudNII6GYo; client_ip=178.49.192.69%2C%2010.0.0.2; tmr_lvid=f4e976b6e991904331b4ec2f9331fd95; tmr_lvidTS=1783054667459; _gcl_au=1.1.1929435479.1783054668; uxs_uid=bd10cf80-769b-11f1-bf43-7f415668cbdb; gcui=; gcmi=; gcvi=19ZadvpP8Bl; fromCart=1; deduplication_cookie=lekvapteke.ru; deduplication_cookie=lekvapteke.ru; tt_deduplication_cookie=lekvapteke.ru; tt_deduplication_cookie=lekvapteke.ru; advcake_track_id=66442a42-1d8b-70a0-1e8d-431f7e2180b7; _ym_d=1783321980; _ym_uid=1783321980100356813; advcake_utm_webmaster=8qLA4W; advcake_click_id=_M_WAgEEAAB4AwAAAWUB; advcake_utm_partner=sovaunion; aprt_last_partner=advcake; _userGUID=0:mqq5xqgp:ZdRv0CYVTNoDqsDeDEYYCoTX8jbJFD22; flocktory-uuid=7056824b-e443-434a-8df0-e81abc74db0f-6; advcake_session_id=3d73d2e5-35f0-a758-1e1f-ef0e87b1979f; tds_session_ctx=1783334971385896; regionId=16; mindboxDeviceUUID=bce9f9bc-7dbe-44cc-a749-5f24c35ae34c; directCrm-session=%7B%22deviceGuid%22%3A%22bce9f9bc-7dbe-44cc-a749-5f24c35ae34c%22%7D; cartId=5736e654-0f1f-449d-8ff8-136443da25ae; _ga=GA1.1.1814590812.1783054812; domain_sid=_VjyCax896NPQiiycU4k2%3A1784456282498; shop_session=Y34KJwRtcUhZQUZTuvrNm3izvQBHAGoAtJsPRxHQ; evnt_s=Y34KJwRtcUhZQUZTuvrNm3izvQBHAGoAtJsPRxHQ; digi_uc=|v:178408:360082!178409:359863!178411:368325!178412:419364:302840!178426:16588!178445:367366; __ddg9_=178.49.192.69; _ym_isad=2; _ym_visorc=b; gcsi=ZOXMOCOdCw7; advcake_track_url=%3D20250113K6XXPR8RHPFiFebB82UR6ACabHdx9NMGcI6%2BTl0pzAkjr%2Fmi5zmmifEibegxmJWSymKGqa5X3A8ZV%2BrBFzO5%2FdI%2FbshxnnLAdiXBM0wqzqGz4c87dqIfVoWaTwa0MBhphjkV0d8%2FepVIEDSgTIJdkYBy6gKpjEcH11P7VpHkegenZltW%2FCmWiH63AswlAlak7FIVbAeCoVuIr5KamZnaaJeYUB5ZaBvvND3geoMqb%2Blya4VQP2f4tUz4S1l4QTfo%2FectaFsRjsHm4PoiDhJxa52AObObai1mlks5sIjUPGnlRvolAPG89uZvWFc9HR2RpkuTgGrvz7u1Y0jE56ZHgVxwyU2B17dWwMoGEoNnUF6oUAortSLMcXSiDhWVb5UNkC6Y6kVk7L19b0%2F5vmLQIOyGT9irYQoCh%2FtHps7z4SZ8P1WzgiCX%2BZ07Ou10ed3zEMAl06lOV4E1HB6gPHwU6cTTqUrTc%2B%2FcJ4G1Wmia5KOM7AE%2Bdz%2B9UvFh6yTuAQw%2FeWinK6APHMardo%2FJHmPzw7ITaALzSVpWhiPf46P%2BGqRFP9jsP3CyXVboXDi%2FJrZqGakfEHlyYI7TmlbpDtl4%2B5Yp4Ujq5AnqFq25%2Bg0d%2FDVDOJtQ01%2BNGPdfbUoTFjpfnwYLDpXI%2BlIpDdBBA7WAlf8y16E3%2BrTOdnz7%2Bc7BTaHI7pie6LBIi0w%3D; tmr_detect=0%7C1784528805267; dSesn=8865789e-84f1-02cc-b336-42688286e04a; _dvs=0:mrsudhat:8XKrEF4xuKlEa1e4v1nUMTD8Cgb1hmnf; DIGI_CARTID=40202856239; _ga_LPGWK50SVR=GS2.1.s1784528802$o14$g1$t1784528849$j13$l0$h0; __ddg8_=8nO0b9Tmw3kLNwhZ; __ddg10_=1784528849",
     "payload": {
        "g-recaptcha-response": "",
        "phоne": "79999999999",
        "verificationCode": "1234"
     }
  }
  ```
* авторизация для мобильного приложения:
  **POST /api/v8/ecom-auth/login-sms-prestep**
  ```json
  {
   "headers": {
     "Content-Type": [
       "application/json; charset=UTF-8"
     ],
     "User-Agent": [
       "okhttp/5.0.0-alpha.10"
     ],
     "X-Device-Id": [
       "11825ba1-6cd4-479f-9742-b158bfc3ff86"
     ],
     "X-Device-Id-Mb": [
       "a4b43f8a-c406-4104-972f-766d09441b0a"
     ],
     "X-Os": [
       "Android"
     ],
     "X-Real-Ip": [
       "91.79.122.19"
     ]
   },
   "body": {
     "phone": "79999999999"
   }
  }
  ```
* регистрация пользователя для мобильного приложения
  **POST /api/v8/manzana/registration**
  ```json
  {
    "headers": {
      "Content-Type": [
        "application/json; charset=UTF-8"
      ],
      "User-Agent": [
        "okhttp/5.0.0-alpha.10"
      ],
      "X-Device-Id": [
        "1ee46262-1498-8578-e629-39b771366149"
      ],
      "X-Forwarded-For": [
        "216.126.229.32"
      ],
      "X-Real-Ip": [
        "216.126.229.32"
      ]
    },
    "body": {
      "AgreeToTerms": true,
      "BirthDate": "01.15.1980",
      "DeviceId": "41b126b5-7456-3e18-771b-9527e8994963",
      "EmailAddress": "17Krm22CjS2B@gmail.com",
      "Firstname": "Арай",
      "GenderCode": 1,
      "Lastname": "Верховецкая",
      "MobilePhone": "+79222271391",
      "Password": "ZS3JRbS#t0a1"
    }
   }
   ```

## Запуск для разработки
```bash
docker compose run --rm --service-ports --build cli
```