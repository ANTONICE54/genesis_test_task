Відповідно до умови було реалізовано 3 запити:
- /rate
- /subscibe
- /sendEmails

Для створення API був використаний [Gin web framework]([url](https://github.com/gin-gonic/gin)).
В якості СУБД було обрано PostgreSQL та створено базу даних, що складається з однієї таблиці Emails, яка містить наступні поля:
- id
- email
- created_at
