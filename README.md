# YADRO Telecom Challenge

Система обработки событий для игрового подземелья, написанная в рамках тестового задания на стажировку в команду телеком YADRO.

## Архитектура проекта

- State Machine Engine: Вся бизнес-логика реализована в виде конечного автомата (пакет `processor`).
- Data Streaming: Файл с событиями не загружается в память целиком, а читается построчно через `bufio.Scanner`, что позволяет обрабатывать логи любого размера (Robustness).
- Graceful Error Handling: Парсер не паникует от битых строк, а логирует их в `stderr` и продолжает работу.
- Go 1.22, стандартная библиотека

## Быстрый старт

Проект включает `Makefile` для удобного запуска.

1. Клонируйте репозиторий:

   ```bash
   git clone https://github.com/Mkz-Prog/yadro-telecom-test.git
   cd yadro-telecom-test
   ```

2. Запустите симуляцию (используются `config.json` и `events.txt` из корня проекта):

   ```bash
   make run
   ```

## Тестирование

В проекте реализованы тесты (Table-Driven Tests) для проверки логики парсинга.

```bash
make test
```

## Структура проекта

- `/cmd/app` — точка входа;
- `/internal/config` — загрузчик конфигурации;
- `/internal/domain` — модели Player, Event;
- `/internal/parser` — построчный парсер логов;
- `/internal/processor` — ядро бизнес-логики.
