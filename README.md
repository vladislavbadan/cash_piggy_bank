# Cash Piggy Bank

Проект для достижения финансовых целей, написанный на Go.

## Описание

Cash Piggy Bank - это телеграм бот для постановки и достижения целей, который умеет:
- Запоминать финансовые цели пользователя и следить за прогрессом
- Считать дневной лимит денег
- Отправлять мотивирующие напоминания
- Показывать курсы валют

## Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/vladislavbadan/cash_piggy_bank.git
cd cash_piggy_bank
```

2. Установите зависимости:
```bash
go mod download
```

3. Определите переменную окружения TOKEN

4. Запустите приложение:
```bash
go run cmd/main.go
```

## Структура проекта

```
.
├── cmd/        # Точка входа в приложение
├── internal/   # Внутренние пакеты
└── pkg/        # Публичные пакеты
```

## Лицензия

MIT 
