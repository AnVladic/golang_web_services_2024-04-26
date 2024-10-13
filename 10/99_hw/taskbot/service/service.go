package service

import (
	"sort"
	"sync"
	"taskbot/db"
)

type TaskManager struct {
	Tasks  map[int]*db.Task
	lastId int
	mu     *sync.Mutex
}

func InitTaskManager() *TaskManager {
	return &TaskManager{
		Tasks:  make(map[int]*db.Task),
		lastId: 0,
		mu:     &sync.Mutex{},
	}
}

func (t *TaskManager) AddTask(task *db.Task) int {
	t.mu.Lock()
	t.lastId++
	t.Tasks[t.lastId] = task
	task.Id = t.lastId
	t.mu.Unlock()
	return t.lastId
}

func (t *TaskManager) RemoveTask(task *db.Task) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.Tasks, task.Id)
}

func (t *TaskManager) GetTasks() map[int]*db.Task {
	return t.Tasks
}

func (t *TaskManager) GetSortedKeys() []int {
	keys := make([]int, 0, len(t.Tasks))
	for k := range t.Tasks {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}
