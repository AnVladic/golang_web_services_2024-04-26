package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
	"taskbot/db"
)

func (h *TelegramBotHandler) Route(update tgbotapi.Update) {
	if update.Message == nil || update.Message.Text == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error: Unknown command.")
		_, _ = h.Bot.Send(msg)
	}
	text := update.Message.Text
	chat := &db.Chat{
		Id: update.Message.Chat.ID,
	}
	userId := update.Message.From.ID
	h.Mu.Lock()
	user := h.Users[userId]
	if user == nil {
		user = &db.User{
			Id:   userId,
			Name: update.Message.From.UserName,
		}
		h.Users[userId] = user
	}
	user.Chat = chat
	chat.User = user
	h.Mu.Unlock()

	tgMsgManager := &TelegramMessageManager{
		Bot:    h.Bot,
		Update: &update,
		Chat:   chat,
	}

	if strings.HasPrefix(text, "/tasks") {
		h.SendTasksCommand(tgMsgManager)
	} else if strings.HasPrefix(text, "/new") {
		h.CreateNewTaskCommand(tgMsgManager)
	} else if strings.HasPrefix(text, "/assign_") {
		h.AssignTaskCommand(tgMsgManager)
	} else if strings.HasPrefix(text, "/unassign_") {
		h.UnassignTaskCommand(tgMsgManager)
	} else if strings.HasPrefix(text, "/resolve_") {
		h.ResolveTaskCommand(tgMsgManager)
	} else if strings.HasPrefix(text, "/my") {
		h.GetMyTasks(tgMsgManager)
	} else if strings.HasPrefix(text, "/own") {
		h.GetOwnTasks(tgMsgManager)
	}
}
