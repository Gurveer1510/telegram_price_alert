package telegram

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Gurveer1510/telegram_price_tracker/internal/types"
	"github.com/Gurveer1510/telegram_price_tracker/internal/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	Bot       *tgbotapi.BotAPI
	TGUsecase *usecase.TelegramUsecase
}

func Newbot(token string, tgUsecase *usecase.TelegramUsecase) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true

	return &TelegramBot{Bot: bot, TGUsecase: tgUsecase}, nil
}

func (tb *TelegramBot) GetUpdates() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := tb.Bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		
		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /setalert <instrument> <price> <above|below>"
		case "setalert":
			args := strings.Fields(update.Message.CommandArguments())

			if len(args) != 3 {
				msg.Text = "Usage: /setalert <instrument> <price> <above|below>"
				break
			}
			instrument := strings.ToUpper(args[0])

			price, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				msg.Text = "Invalid price. Example: /setalert IDFCFIRSTB 393.21 above"
				break
			}
			condition := strings.ToLower(args[2])
			if condition != "above" && condition != "below" {
				msg.Text = "Condition must be 'above' or 'below'"
				break
			}
			newAlert := &types.Alert{
				Instrument_name: instrument,
				ChatId:          update.Message.Chat.ID,
				Exchange:        "NSE",
				Trigger_price:   price,
				Condition:       condition,
			}
			_, err = tb.TGUsecase.CreateAlert(context.Background(), newAlert)
			if err != nil {
				log.Println(err.Error())
				msg.Text = "Something went wrong :("
				break
			}
			msg.Text = fmt.Sprintf("Alert set! %s @ %.2f (%s)", instrument, price, condition)
		case "setaccesstoken":
			if update.Message.From.UserName != "gurveer1510" {
				msg.Text = "Not authorized to use this command"
			} else {
				args := strings.Fields(update.Message.CommandArguments())
				if len(args) != 1 {
					msg.Text = "Usage: /setaccesstoken <token>"
				} else {
					token := args[0]
					err := tb.TGUsecase.StoreAccessToken(context.Background(), token)
					if err != nil{
						msg.Text = err.Error()
						break
					}
					msg.Text = "Access token set successfully"
				}
			}

		default:
			msg.Text = "I don't know that command"
		}
		if _, err := tb.Bot.Send(msg); err != nil {
			log.Println(err.Error())
		}
	}
}

func (tb *TelegramBot) SetAlert() {

}
