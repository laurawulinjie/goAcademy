package todo

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestActor_ConcurrentUpdateSameTodo(t *testing.T) {
	ctx, ctxDone := context.WithCancel(context.Background())
	defer ctxDone()

	StartTodoActor(ctx)

	createResChan := make(chan any)
	RequestQueue <- Request{
		Ctx:      ctx,
		Action:   "create",
		Payload:  "Original Task",
		Response: createResChan,
	}

	createRes := (<-createResChan).(struct {
		Todo Todo
		Err  error
	})

	if createRes.Err != nil {
		t.Fatalf("Failed to create todo: %v", createRes.Err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		updateResChan := make(chan any)
		RequestQueue <- Request{
			Ctx:    ctx,
			Action: "update",
			Payload: struct {
				ID     int
				Task   string
				Status string
			}{
				ID:     1,
				Task:   "Task-update1",
				Status: Started,
			},
			Response: updateResChan,
		}
		<-updateResChan
	}()

	go func() {
		defer wg.Done()
		updateResChan := make(chan any)
		RequestQueue <- Request{
			Ctx:    ctx,
			Action: "update",
			Payload: struct {
				ID     int
				Task   string
				Status string
			}{
				ID:     1,
				Task:   "Task-update2",
				Status: Completed,
			},
			Response: updateResChan,
		}
		<-updateResChan
	}()

	wg.Wait()

	getResChan := make(chan any)
	RequestQueue <- Request{
		Ctx:      ctx,
		Action:   "getAll",
		Payload:  nil,
		Response: getResChan,
	}

	todos := (<-getResChan).(map[int]Todo)
	fmt.Println("updated todos: ", todos)

	if !(todos[1].Task == "Task-update1" || todos[1].Task == "Task-update2") {
		t.Fatalf("Unexpected final task value: %v", todos[1].Task)
	}
}

func TestActor_ConcurrentDeleteSameTodo(t *testing.T) {
	ctx, ctxDone := context.WithCancel(context.Background())
	defer ctxDone()

	StartTodoActor(ctx)

	createResChan := make(chan any)
	RequestQueue <- Request{
		Ctx:      ctx,
		Action:   "create",
		Payload:  "Original Task",
		Response: createResChan,
	}

	createRes := (<-createResChan).(struct {
		Todo Todo
		Err  error
	})

	if createRes.Err != nil {
		t.Fatalf("Failed to create todo: %v", createRes.Err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	errors := make(chan error, 2)

	go func() {
		defer wg.Done()
		deleteResChan := make(chan any)
		RequestQueue <- Request{
			Ctx:      ctx,
			Action:   "delete",
			Payload:  1,
			Response: deleteResChan,
		}
		res := (<-deleteResChan).(struct {
			Err error
		})
		errors <- res.Err
	}()

	go func() {
		defer wg.Done()
		deleteResChan := make(chan any)
		RequestQueue <- Request{
			Ctx:      ctx,
			Action:   "delete",
			Payload:  1,
			Response: deleteResChan,
		}
		res := (<-deleteResChan).(struct {
			Err error
		})
		errors <- res.Err
	}()

	wg.Wait()
	close(errors)

	successCount := 0
	errorCount := 0

	for err := range errors {
		if err == nil {
			successCount++
		} else if err.Error() == "todo not found" {
			errorCount++
		} else {
			t.Fatalf("Unexpected error: %v", err)
		}
	}

	if successCount != 1 || errorCount != 1 {
		t.Fatalf("Expected 1 success and 1 error, got %d success and %d error", successCount, errorCount)
	}
}
