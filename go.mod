module github.com/d0ctr/bldbr-bot-go

go 1.21

require (
	github.com/dotenv-org/godotenvvault v0.6.0
	gopkg.in/telebot.v3 v3.1.3
)

require github.com/joho/godotenv v1.5.1 // indirect

replace gopkg.in/telebot.v3 => github.com/d0ctr/telebot v0.0.0-20230820223142-a42b93b462ff
