package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
	"os"
	"strconv"
	"tg-dice-roller/dice"
	"time"
)

var token string
var Bot *tgbotapi.BotAPI
var cache = make(map[int][]string, 0)
var randomSource = rand.NewSource(time.Now().UnixNano())
var randGenerator = rand.New(randomSource)

func init() {
	token = os.Getenv("BOT_TOKEN")
}

func main() {
	Start()
}

func Start() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	Bot = bot

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		query := update.InlineQuery
		if query != nil {

			fmt.Printf("query: %v, %v\n", query.ID, query.Query)
			var responses []interface{}

			if query.Query == "" { //todo optimize duplicated code
				for _, str := range cache[query.From.ID] {
					res, explanation, err := dice.Roll(str)
					if err != nil {
						continue
					}
					message := fmt.Sprintf("<u>%v:</u> %v = <b>%v</b>", str, explanation, res)

					responses = append(responses, tgbotapi.InlineQueryResultArticle{
						Type:  "article",
						ID:    strconv.FormatInt(randGenerator.Int63(), 10),
						Title: str,

						InputMessageContent: map[string]string{
							"message_text": message,
							"parse_mode":   "HTML",
						},
					})
				}
				inlineConfig := tgbotapi.InlineConfig{
					InlineQueryID: query.ID,
					Results:       responses,
				}
				_, err = Bot.AnswerInlineQuery(inlineConfig)

				continue
			}

			//todo add timeout
			res, explanation, err := dice.Roll(query.Query)
			if err != nil {
				continue
			}
			updateCache(query.From.ID, query.Query) //todo use choseninlineresult to rotate the cache

			message := fmt.Sprintf("<u>%v:</u> %v = <b>%v</b>", query.Query, explanation, res)
			responses = append(responses, tgbotapi.InlineQueryResultArticle{
				Type:  "article",
				ID:    strconv.FormatInt(time.Now().UnixNano(), 10),
				Title: "Roll it!",

				InputMessageContent: map[string]string{
					"message_text": message,
					"parse_mode":   "HTML",
				},
			})

			inlineConfig := tgbotapi.InlineConfig{
				InlineQueryID: query.ID,
				Results:       responses,
			}
			_, err = Bot.AnswerInlineQuery(inlineConfig)
		}
	}
}

func updateCache(userId int, query string) {
	if len(cache[userId]) < 5 {
		cache[userId] = append(cache[userId], query) //todo check existence
	} else {
		for i := 0; i < len(cache[userId])-1; i++ {
			cache[userId][i] = cache[userId][i+1]
		}
		cache[userId][4] = query
	}
}
