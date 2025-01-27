# TODO List

TODO List - это веб-приложение на Go для управления задачами с поддержкой повторяющихся событий.


## Установка

### Склонируйте репозиторий
```
git clone https://github.com/Epicpt/simple_todo.git
cd simple_todo
go mod tidy
```
## Запуск
```
go run main.go
```

Сервер доступен по адресу http://localhost:7540

## Стэк технологий

* golang
* SQLite
* html, css, js

<details>
<summary>Структура</summary>
<ul>
<li>main.go: Главный файл приложения, точка входа сервера.</li>
<li>database/: Пакет для инициализации базы данных, взаимодействия с базой данных.</li>
<li>handlers/: Пакет с обработчиками API запросов.</li>
<li>model/: Пакет с моделями данных.</li>
<li>api/: Пакет с логикой обработки задач и вычисления следующей даты выполнения.</li>
<li>web/: Директория для статических файлов (HTML, CSS, JS).</li>
</ul>
</details>
