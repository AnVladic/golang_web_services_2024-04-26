package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"sync"
	"taskbot/db"
	"taskbot/service"
)

type TelegramBotHandler struct {
	service.BotController
	Bot *tgbotapi.BotAPI
}

type TelegramMessageManager struct {
	Bot    *tgbotapi.BotAPI
	Update *tgbotapi.Update
	Chat   *db.Chat
}

func (h *TelegramBotHandler) GetTaskById(idStr string) (*db.Task, error) {
	taskId, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	task := h.TaskManager.Tasks[taskId]
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}
	return task, nil
}

func (h *TelegramBotHandler) SetDefault() {
	h.TaskManager = service.InitTaskManager()
	h.Users = map[int]*db.User{}
	h.Mu = &sync.Mutex{}
}

func (h *TelegramMessageManager) SendText(text string) error {
	msg := tgbotapi.NewMessage(h.Chat.Id, text)
	_, err := h.Bot.Send(msg)
	return err
}

func (h *TelegramMessageManager) GetChat() *db.Chat {
	return h.Chat
}

func (h *TelegramMessageManager) SetChat(chat *db.Chat) {
	h.Chat = chat
}
