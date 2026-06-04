# 🎮 dotaward-bot

Telegram бот для отслеживания статистики Dota 2. Привяжи свой аккаунт и следи за матчами прямо в чате.

## Команды

| Команда | Описание |
|--------|----------|
| `/register <dota_account_id>` | Привязать Dota аккаунт |
| `/lastmatch` | Статистика последнего матча |
| `/help` | Список команд |

## В планах

- `/lastmatch @username` — статистика последнего матча конкретного юзера
- `/top` — топ игроков группы за неделю
- `/profile` — винрейт, любимый герой, средний KDA
- `/compare @user1 @user2` — сравнить двух игроков
- `/rank` — текущий ранг и MMR
- Автоматические уведомления о новых матчах

## Запуск

```bash
BOT_TOKEN=токен ./dotaward-bot
```

> Account ID можно найти на [opendota.com](https://opendota.com) после входа через Steam, либо же скопировать из самой игры.

## Стек

- [Go 1.26+](https://go.dev)
- [telebot v3](https://gopkg.in/telebot.v3)
- [OpenDota API](https://docs.opendota.com)