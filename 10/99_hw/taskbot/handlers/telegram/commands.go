package handlers

import (
	"strings"
	"taskbot/db"
)

func (h *TelegramBotHandler) SendTasksCommand(tgMsgManager *TelegramMessageManager) {
	h.SendTasks(tgMsgManager)
}

func (h *TelegramBotHandler) CreateNewTaskCommand(tgMsgManager *TelegramMessageManager) {
	title := strings.TrimSpace(strings.TrimPrefix(tgMsgManager.Update.Message.Text, "/new"))
	newTask := db.Task{
		Title:  title,
		Author: tgMsgManager.Chat.User,
	}
	h.CreateTask(tgMsgManager, &newTask)
}

func (h *TelegramBotHandler) AssignTaskCommand(tgMsgManager *TelegramMessageManager) {
	idStr := strings.TrimSpace(strings.TrimPrefix(tgMsgManager.Update.Message.Text, "/assign_"))
	task, err := h.GetTaskById(idStr)
	if err != nil {
		h.SendError(tgMsgManager, err)
		return
	}
	h.AssignTask(tgMsgManager, task)
}

func (h *TelegramBotHandler) UnassignTaskCommand(tgMsgManager *TelegramMessageManager) {
	idStr := strings.TrimSpace(strings.TrimPrefix(tgMsgManager.Update.Message.Text, "/unassign_"))
	task, err := h.GetTaskById(idStr)
	if err != nil {
		h.SendError(tgMsgManager, err)
		return
	}
	h.UnassignTask(tgMsgManager, task)
}

func (h *TelegramBotHandler) ResolveTaskCommand(tgMsgManager *TelegramMessageManager) {
	idStr := strings.TrimSpace(strings.TrimPrefix(tgMsgManager.Update.Message.Text, "/resolve_"))
	task, err := h.GetTaskById(idStr)
	if err != nil {
		h.SendError(tgMsgManager, err)
		return
	}
	h.ResolveTask(tgMsgManager, task)
}

func (h *TelegramBotHandler) GetMyTaskCommand(tgMsgManager *TelegramMessageManager) {
	h.GetMyTasks(tgMsgManager)
}

func (h *TelegramBotHandler) GetOwnTaskCommand(tgMsgManager *TelegramMessageManager) {
	h.GetOwnTasks(tgMsgManager)
}
