# 🎮 dotaward-bot

## О проекте

Telegram-бот для отслеживания матчей с друзьями. Достаточно привязать свой аккаунт и можно следить за матчами прямо в чате. Под капотом - Go, SQLite и OpenDota API.

## Структура

```
.
├── cmd/bot/          # точка входа
├── config/           # загрузка конфига из .env
└── internal/
    ├── handlers/     # обработчики команд и колбэков
    ├── models/       # модели данных
    ├── opendota/     # клиент OpenDota API, хелперы
    └── repository/   # SQLite репозиторий
```

## Команды

| Команда | Описание |
|--------|----------|
| `/register <dota_account_id>` | Привязать Dota аккаунт |
| `/profile` | Статистика аккаунта |
| `/lastmatch [@username]` | Статистика последнего матча |
| `/streak` | Текущая серия (последние 20 матчей) |
| `/maxstreak` | Максимальная серия за всё время |
| `/why` | AI анализ матча (new) |
| `/help` | Список команд |

> `/lastmatch` - без аргументов твой матч, с `@username` - чужой, ответом на сообщение - того на кого ответил.

## В планах

- `/top` — топ игроков группы за неделю
- `/profile` — расширить: любимый герой, средний KDA
- `/compare @user1 @user2` — сравнить двух игроков
- `/rank` — текущий ранг и MMR
- Автоматические уведомления о новых матчах

## Быстрый старт

1. Клонируй репозиторий

```bash
git clone https://github.com/nxwex/dotaward-bot
cd dotaward-bot
```

2. Создай `.env` файл

```bash
cp .env.example .env
```

Открой `.env` и вставь токен от [@BotFather](https://t.me/BotFather):

```env
BOT_TOKEN=твой_токен
DB_PATH=users.db
```

3. Запусти

```bash
go run ./cmd/bot
```

> Account ID можно найти на [opendota.com](https://opendota.com) после входа через Steam, либо скопировать из самой игры.

## Стек

- [Go 1.26+](https://go.dev)
- [telebot v3](https://gopkg.in/telebot.v3)
- [OpenDota API](https://docs.opendota.com)