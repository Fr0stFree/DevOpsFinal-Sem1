# Финальный проект 1 семестра

REST API сервис для загрузки и выгрузки данных о ценах.

## Требования к системе

- go 1.23.3
- postgres 13.3+ (предполагается, что база данных уже создана и доступна по адресу `localhost:5432`)

## Установка и запуск

Для компиляции приложения и создания таблицы в БД необходимо выполнить команду:

```bash
chmod +x scripts/prepare.sh
./scripts/prepare.sh
```

Для запуска приложения необходимо выполнить команду:

```bash
chmod +x scripts/run.sh
./scripts/run.sh
```

## Тестирование

Директория `sample_data` - это пример директории, которая является разархивированной версией файла `sample_data.zip`

для заупуска тестов необходимо выполнить команду:

```bash
chmod +x scripts/test.sh
./scripts/test.sh 3
```

## Контакты

me@frostfree.ru
или 
https://t.me/fr0stfree

https://github.com/Fr0stFree/DevOpsFinal-Sem1
