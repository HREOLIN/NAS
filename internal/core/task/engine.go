package task

import (
	"time"

	"github.com/HREOLIN/NAS/internal/core/store"
)

type taskError struct {
	message  string
	rollback bool
}

func (e *taskError) Error() string {
	return e.message
}

type Context struct {
	task *store.Task
}

func (c *Context) Step(name, detail string) {
	c.task.Steps = append(c.task.Steps, store.TaskStep{
		Name:   name,
		Detail: detail,
		At:     time.Now().UTC(),
	})
	if c.task.Progress < 90 {
		c.task.Progress += 20
	}
}

func (c *Context) Fail(message string, rollback bool) error {
	c.task.Rollback = rollback
	c.task.Error = message
	return &taskError{
		message:  message,
		rollback: rollback,
	}
}

type Engine struct {
	store *store.Store
}

func NewEngine(state *store.Store) *Engine {
	return &Engine{store: state}
}

func (e *Engine) Run(module, action string, payload any, run func(*Context) (any, error)) *store.Task {
	task := &store.Task{
		ID:        e.store.NextID("task"),
		Module:    module,
		Action:    action,
		Status:    "queued",
		Progress:  0,
		Payload:   payload,
		Steps:     []store.TaskStep{},
		StartedAt: time.Now().UTC(),
	}

	e.store.AddTask(task)
	task.Status = "running"

	result, err := run(&Context{task: task})
	task.EndedAt = time.Now().UTC()

	if err != nil {
		if task.Rollback {
			task.Status = "rolled_back"
		} else {
			task.Status = "failed"
		}
		task.Result = result
		if task.Error == "" {
			task.Error = err.Error()
		}
		return task
	}

	task.Progress = 100
	task.Result = result
	if task.Rollback {
		task.Status = "rolled_back"
	} else {
		task.Status = "success"
	}
	return task
}
