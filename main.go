package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"tg-dice-roller/dice"
	"time"
)

var randomSource = rand.NewSource(time.Now().UnixNano())
var randGenerator = rand.New(randomSource)

var token string
var Bot *tgbotapi.BotAPI

var cache = make(map[int][]string, 0)

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
		chosenResult := update.ChosenInlineResult
		query := update.InlineQuery

		if chosenResult != nil && chosenResult.Query != "" { // change it (!= "") if the cache rotation should be implemented
			trimmedQuery := strings.ReplaceAll(chosenResult.Query, " ", "")
			updateCache(chosenResult.From.ID, trimmedQuery)
			continue
		}

		if query != nil {
			var responses []interface{}
			trimmedQuery := strings.ReplaceAll(query.Query, " ", "")
			if trimmedQuery == "" {
				for _, str := range cache[query.From.ID] {
					res, err := prepareInlineResult(str, str)
					if err != nil {
						continue
					}
					responses = append(responses, res)
				}
			} else {
				res, err := prepareInlineResult(trimmedQuery, "Roll it!")
				if err != nil {
					continue
				}
				responses = append(responses, res)
			}
			answerInline(responses, query.ID)
		}
	}
}

func updateCache(userId int, query string) {
	if contains(cache[userId], query) { //todo rotate cache to move to top the often-used roll
		return
	}
	if len(cache[userId]) < 5 {
		cache[userId] = append(cache[userId], query)
	} else {
		for i := 0; i < len(cache[userId])-1; i++ {
			cache[userId][i] = cache[userId][i+1]
		}
		cache[userId][4] = query
	}
}

func contains(slice []string, elem string) bool {
	for _, s := range slice {
		if s == elem {
			return true
		}
	}
	return false
}

func prepareInlineResult(rollStr, title string) (tgbotapi.InlineQueryResultArticle, error) {
	res, explanation, err := dice.Roll(rollStr)
	if err != nil {
		fmt.Printf("Can't roll the dice: %v, error: %v", rollStr, err)
		return tgbotapi.InlineQueryResultArticle{}, err
	}
	var message string
	if explanation != "" {
		message = fmt.Sprintf("<u>%v:</u> %v = <b>%v</b>", rollStr, explanation, res)
	} else {
		message = fmt.Sprintf("<u>%v:</u> <b>%v</b>", rollStr, res)
	}

	return tgbotapi.InlineQueryResultArticle{
		Type:  "article",
		ID:    strconv.FormatInt(randGenerator.Int63(), 10),
		Title: title,
		InputMessageContent: map[string]string{
			"message_text": message,
			"parse_mode":   "HTML",
		},
	}, nil
}

func answerInline(responses []interface{}, queryId string) {
	inlineConfig := tgbotapi.InlineConfig{
		InlineQueryID: queryId,
		Results:       responses,
		CacheTime:     0,
	}
	_, _ = Bot.AnswerInlineQuery(inlineConfig)
}
