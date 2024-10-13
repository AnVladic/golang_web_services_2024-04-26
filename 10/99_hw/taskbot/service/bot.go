package service

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"taskbot/db"
)

type BotController struct {
	TaskManager *TaskManager
	Users       map[int]*db.User
	Mu          *sync.Mutex
}

type MessageManager interface {
	SendText(text string) error
	GetChat() *db.Chat
	SetChat(chat *db.Chat)
}

func (b *BotController) trySendText(msgManager MessageManager, text string) {
	err := msgManager.SendText(text)
	if err != nil {
		_ = fmt.Errorf("send error: %s", err)
	}
}

func (b *BotController) SendError(msgManager MessageManager, err error) {
	b.trySendText(msgManager, fmt.Sprintf("Ошибка: %s", err))
}

func (b *BotController) GenerateTaskInfo(chat db.Chat, task db.Task, showMeAssignee bool) string {
	assignee := ""
	commands := ""
	if task.Performer == nil {
		commands = fmt.Sprintf("/assign_%d", task.Id)
	} else {
		if chat.User.Id == task.Performer.Id {
			if showMeAssignee {
				assignee = "assignee: я\n"
			}
			commands = fmt.Sprintf("/unassign_%d /resolve_%d", task.Id, task.Id)
		} else {
			assignee = fmt.Sprintf("assignee: @%s", task.Performer.Name)
		}
	}
	return fmt.Sprintf(`%d. %s by @%s
%s%s`, task.Id, task.Title, task.Author.Name, assignee, commands)
}

func (b *BotController) SendTasks(msgManager MessageManager) {
	tasks := b.TaskManager.GetTasks()
	if len(tasks) == 0 {
		b.trySendText(msgManager, "Нет задач")
		return
	}
	chat := msgManager.GetChat()
	tasksTexts := make([]string, 0, len(tasks))
	sortedTasksIds := b.TaskManager.GetSortedKeys()
	for _, id := range sortedTasksIds {
		tasksTexts = append(tasksTexts, b.GenerateTaskInfo(*chat, *tasks[id], true))
	}

	b.trySendText(msgManager, strings.Join(tasksTexts, "\n\n"))
}

func (b *BotController) CreateTask(msgManager MessageManager, task *db.Task) {
	taskId := b.TaskManager.AddTask(task)
	b.trySendText(msgManager, fmt.Sprintf("Задача \"%s\" создана, id=%d", task.Title, taskId))
}

func (b *BotController) AssignTask(msgManager MessageManager, task *db.Task) {
	prevPerformer := task.Performer
	if prevPerformer == nil {
		prevPerformer = task.Author
	}
	task.Performer = msgManager.GetChat().User
	b.trySendText(
		msgManager, fmt.Sprintf("Задача \"%s\" назначена на вас", task.Title))

	if prevPerformer.Id != task.Performer.Id {
		msgManager.SetChat(prevPerformer.Chat)
		b.trySendText(
			msgManager, fmt.Sprintf("Задача \"%s\" назначена на @%s", task.Title, task.Performer.Name))
	}
}

func (b *BotController) UnassignTask(msgManager MessageManager, task *db.Task) {
	user := msgManager.GetChat().User
	if task.Performer != nil && user.Id == task.Performer.Id {
		task.Performer = nil
		b.trySendText(msgManager, "Принято")
		msgManager.SetChat(task.Author.Chat)
		b.trySendText(msgManager, fmt.Sprintf("Задача \"%s\" осталась без исполнителя", task.Title))
	} else {
		b.trySendText(msgManager, "Задача не на вас")
	}
}

func (b *BotController) ResolveTask(msgManager MessageManager, task *db.Task) {
	user := msgManager.GetChat().User
	if task.Performer != nil && user.Id == task.Performer.Id {
		b.TaskManager.RemoveTask(task)
		b.trySendText(msgManager, fmt.Sprintf("Задача \"%s\" выполнена", task.Title))
		msgManager.SetChat(task.Author.Chat)
		b.trySendText(
			msgManager, fmt.Sprintf("Задача \"%s\" выполнена @%s", task.Title, user.Name))
	} else {
		b.trySendText(msgManager, "Задача не на вас")
	}
}

func (b *BotController) GetTaskBy(
	msgManager MessageManager, showMeAssignee bool, comparison func(task *db.Task, chat *db.Chat) bool) {

	chat := msgManager.GetChat()
	allTasks := b.TaskManager.GetTasks()
	var tasks []*db.Task
	for _, task := range allTasks {
		if comparison(task, chat) {
			tasks = append(tasks, task)
		}
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Id < tasks[j].Id
	})
	tasksTexts := make([]string, 0, len(tasks))
	for _, task := range tasks {
		tasksTexts = append(tasksTexts, b.GenerateTaskInfo(*chat, *task, showMeAssignee))
	}
	b.trySendText(msgManager, strings.Join(tasksTexts, "\n\n"))
}

func (b *BotController) GetMyTasks(msgManager MessageManager) {
	b.GetTaskBy(msgManager, false, func(task *db.Task, chat *db.Chat) bool {
		return task.Performer != nil && task.Performer.Id == chat.User.Id
	})
}

func (b *BotController) GetOwnTasks(msgManager MessageManager) {
	b.GetTaskBy(msgManager, true, func(task *db.Task, chat *db.Chat) bool {
		return task.Author != nil && task.Author.Id == chat.User.Id
	})
}
