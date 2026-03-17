package todo

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"connectrpc.com/connect"
	todov1 "github.com/example/grpc-learning/gen/go/todo/v1"
)

type Service struct {
	mu    sync.Mutex
	next  int
	todos []*todov1.Todo
}

func NewService() *Service {
	return &Service{next: 1}
}

func (s *Service) AddTodo(
	_ context.Context,
	req *connect.Request[todov1.AddTodoRequest],
) (*connect.Response[todov1.AddTodoResponse], error) {
	title := strings.TrimSpace(req.Msg.Title)
	if title == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("title is required"))
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	todo := &todov1.Todo{
		Id:            fmt.Sprintf("%d", s.next),
		Title:         title,
		CreatedAtUnix: time.Now().Unix(),
	}
	s.next++
	s.todos = append(s.todos, todo)

	return connect.NewResponse(&todov1.AddTodoResponse{Todo: todo}), nil
}

func (s *Service) ListTodos(
	_ context.Context,
	_ *connect.Request[todov1.ListTodosRequest],
) (*connect.Response[todov1.ListTodosResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	copied := make([]*todov1.Todo, len(s.todos))
	copy(copied, s.todos)

	return connect.NewResponse(&todov1.ListTodosResponse{Todos: copied}), nil
}
