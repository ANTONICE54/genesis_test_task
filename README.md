## Функціонал
Відповідно до умови було реалізовано 3 запити:
- /rate - повертає поточний курс долара (USD) у гривні (UAH);
- /subscibe - підписує певний емейл на розсилку інформації про курс
- /sendEmails - відправляє на всі підписані емейли листи з актуальним курсом

Також раз на день відбувається автоматична розсилка листів.

Важливо зазначити, що для відправки кожного повідомлення використовується окрема горутина(goroutine), а це в свою чергу значно підвищує продуктивність системи.

## Опис компонентів




### ***Web Framework***
Для створення API був використаний [Gin web framework](https://github.com/gin-gonic/gin).



### ***DB***
В якості СУБД було обрано PostgreSQL та створено базу даних, що складається з однієї таблиці Emails, яка містить наступні поля:
- id
- email
- created_at

При створенні таблиці, для зручності подальшого тестування, до неї також додаються 3 емейли:

example111@gmail.com; 

example222@gmail.com; 

example333@gmail.com.





### ***SMTP server***
Для тестування відправки електронних листів був обраний [Mailhog](https://mailtrap.io/blog/mailhog-explained/).




### ***Third-party API***
Для того, щоб отримувати інформацію про курс обміну я використав [ExchangeRate-API](https://www.exchangerate-api.com/), проте в його використанні є свої недоліки - використовуючи безплатний план дані про курс оновлюються раз на день та кількість запитів на місяць обмежена( максимально можна відправити 1500 запитів).





## How to use
Для того, щоб запустити сервер необхідно в терміналі перейти в корневу папку з проектом та виконати команду:
```
docker compose up
```
Після запуску контейнерів можна використовувати весь доступний функціонал.

Також можна налаштувати час, в який будуть автоматично розсилатись емейли, відредагувавши значення DAILY_EMAILS_TIME в app.env файлі, перед збіркою контейнерів. 

**Час має бути введений в форматі ГОДИНИ:ХВИЛИНИ:СЕКУНДИ.**

Варто відмітити, що для використання ExchangeRate-API необхідно зареєструватись, після чого можна буде отримати ключ доступу, що використовується при відправці запитів. На даний момент використовується мій ключ доступу, проте поміняти його можна шляхом редагування змінної RATE_API_KEY в app.env файлі.


Для того, щоб переглянути відправлені емейли необхідно в браузері перейти за посиланням localhost:8025.

Для відправки запитів використовуйте localhost:8080


## Короткий опис основних компонентів коду
В функції main ініціалізується структура Server та запускається HTTP сервер. Структура Server містить в собі всі поля необхідні для функціонування серверу:
- router - відповідає за маршрутизацію;
- config - структура, що містить в собі змінні необхідні для роботи сервера;
- InfoLog, ErrorLog - логери;
- store - обгортка навколо [sql.DB](https://pkg.go.dev/database/sql#DB), яка дає змогу визначити свої DAO методи;
- wait - об'єкт типу sync.WaitGroup, який допомагає контролювати закінчення роботи всіх горутин перед виходом програми;
- mailer - об'єкт за допомогою якого відбувається відправлення листів.
- cronOperator - за допомогою цього об'єкту виконується розсилка листів раз на день.

**Міграція бази даних відбувається в коді, використовуючи функцію під назвою runDBMifration().**
