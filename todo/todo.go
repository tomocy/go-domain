package todo

import (
	"context"
	"fmt"
)

const (
	createURL = "http://localhost"
)

type get struct {
	repo repo
}

func (g get) do(ctx context.Context) ([]task, error) {
	return g.repo.get(ctx)
}

type create struct {
	repo repo
}

func (c create) do(ctx context.Context, rawTitle string) (task, error) {
	title := taskTitle(rawTitle)
	if err := validate(title); err != nil {
		return task{}, err
	}

	return c.repo.create(ctx, title)
}

type repo interface {
	get(ctx context.Context) ([]task, error)
	create(ctx context.Context, title taskTitle) (task, error)
}

type task struct {
	ID    taskID    `json:"id"`
	Title taskTitle `json:"title"`
}

func (t task) validate() error {
	if err := validate(t.ID, t.Title); err != nil {
		return err
	}

	return nil
}

type taskID string

func (i taskID) validate() error {
	if i == "" {
		return fmt.Errorf("task id should be set")
	}
	return nil
}

type taskTitle string

func (t taskTitle) validate() error {
	if t == "" {
		return fmt.Errorf("task title should be set")
	}
	return nil
}

type errorResp struct {
	Message string `json:"error"`
}

func (r errorResp) Error() string {
	return r.Message
}

func validate(vs ...validator) error {
	for _, v := range vs {
		if err := v.validate(); err != nil {
			return err
		}
	}
	return nil
}

type validator interface {
	validate() error
}
