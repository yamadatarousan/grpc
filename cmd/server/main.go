package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	todov1connect "github.com/example/grpc-learning/gen/go/todo/v1/todov1connect"
	"github.com/example/grpc-learning/internal/todo"
)

func main() {
	mux := http.NewServeMux()

	todoService := todo.NewService()
	path, handler := todov1connect.NewTodoServiceHandler(todoService)
	mux.Handle(path, handler)

	webDist := filepath.Join("web", "dist")
	if _, err := os.Stat(webDist); err == nil {
		mux.Handle("/", http.FileServer(http.Dir(webDist)))
	} else if errors.Is(err, os.ErrNotExist) {
		mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = w.Write([]byte("web/dist が見つかりません。`cd web && npm install && npm run build` を先に実行してください。\n"))
		})
	} else {
		log.Fatalf("failed to read web/dist: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("server listening on http://localhost%s", addr)
	log.Printf("grpc-web endpoint: http://localhost%s%s", addr, todov1connect.TodoServiceAddTodoProcedure)

	if err := http.ListenAndServe(addr, h2c.NewHandler(mux, &http2.Server{})); err != nil {
		log.Fatal(err)
	}
}
