package todo

import (
	"context"
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	t.Parallel()

	get := get{
		repo: mockRepo{
			tasks: []task{
				{
					Title: "a",
				},
				{
					Title: "b",
				},
				{
					Title: "c",
				},
			},
		},
	}

	okTests := []struct {
		ctx      context.Context
		expected []task
	}{
		{
			ctx: context.Background(),
			expected: []task{
				{
					Title: "a",
				},
				{
					Title: "b",
				},
				{
					Title: "c",
				},
			},
		},
	}

	for _, test := range okTests {
		test := test

		t.Run("ok", func(t *testing.T) {
			t.Parallel()

			actual, err := get.do(test.ctx)
			if err != nil {
				t.Errorf("failed to get: %v", err)
			}
			for i, actual := range actual {
				if actual.Title != test.expected[i].Title {
					t.Errorf("unexpected title: got %v, but expected %v", actual.Title, test.expected[i].Title)
				}
			}
		})
	}

	failedTests := []struct {
		ctx    context.Context
		reason string
	}{
		{
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 0)
				defer cancel()
				return ctx
			}(),
			reason: "timed out",
		},
		{
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				return ctx
			}(),
			reason: "canceled",
		},
	}

	for _, test := range failedTests {
		test := test

		t.Run("failed", func(t *testing.T) {
			t.Parallel()

			if _, err := get.do(test.ctx); err == nil {
				t.Errorf("should fail to get: %v", test.reason)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()

	create := create{
		repo: mockRepo{},
	}

	okTests := []struct {
		ctx      context.Context
		title    string
		expected task
	}{
		{
			ctx:   context.Background(),
			title: "Buy some milk",
			expected: task{
				Title: "Buy some milk",
			},
		},
	}

	for _, test := range okTests {
		test := test

		t.Run(fmt.Sprintf("ok:%v", test.title), func(t *testing.T) {
			t.Parallel()

			actual, err := create.do(test.ctx, test.title)
			if err != nil {
				t.Errorf("failed to create: %v", err)
			}
			if err := actual.validate(); err != nil {
				t.Errorf("task should be valid: %v", err)
			}
			if actual.Title != test.expected.Title {
				t.Errorf("unexpected title: got %v, but expected %v", actual.Title, test.expected.Title)
			}
		})
	}

	failedTests := []struct {
		ctx    context.Context
		title  string
		reason string
	}{
		{
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 0)
				defer cancel()
				return ctx
			}(),
			reason: "timed out",
		},
		{
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				return ctx
			}(),
			reason: "canceled",
		},
		{
			ctx:    context.Background(),
			reason: "empty title",
		},
	}

	for _, test := range failedTests {
		test := test

		t.Run(fmt.Sprintf("failed:%v", test.title), func(t *testing.T) {
			t.Parallel()

			if _, err := create.do(test.ctx, test.title); err == nil {
				t.Errorf("should fail to create: %v", test.reason)
			}
		})
	}
}

func TestTask(t *testing.T) {
	t.Parallel()

	okTests := []task{
		{
			ID:    "bd9c3caa-1ec9-4023-b43f-aba41ec5f570",
			Title: "Buy some milk",
		},
	}

	for _, test := range okTests {
		test := test

		t.Run(fmt.Sprintf("ok:%v", test), func(t *testing.T) {
			t.Parallel()

			if err := test.validate(); err != nil {
				t.Errorf("task should be valid: %v", err)
			}
		})
	}

	failedTests := []struct {
		task   task
		reason string
	}{
		{
			task:   task{},
			reason: "zero value",
		},
		{
			task: task{
				Title: "Call him",
			},
			reason: "empty id",
		},
		{
			task: task{
				ID: "bd9c3caa-1ec9-4023-b43f-aba41ec5f570",
			},
			reason: "empty title",
		},
	}

	for _, test := range failedTests {
		test := test

		t.Run(fmt.Sprintf("failed:%v", test.task), func(t *testing.T) {
			t.Parallel()

			if err := test.task.validate(); err == nil {
				t.Errorf("task should not be valid: %v", test.reason)
			}
		})
	}
}

type mockRepo struct {
	tasks []task
}

func (r mockRepo) get(ctx context.Context) ([]task, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return r.tasks, nil
	}
}

func (r mockRepo) create(ctx context.Context, title taskTitle) (task, error) {
	select {
	case <-ctx.Done():
		return task{}, ctx.Err()
	default:
		if err := validate(title); err != nil {
			return task{}, err
		}

		return task{
			ID:    "bd9c3caa-1ec9-4023-b43f-aba41ec5f570",
			Title: title,
		}, nil
	}
}
