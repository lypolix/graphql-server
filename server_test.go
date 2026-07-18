package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"graphql-server/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/vektah/gqlparser/v2/ast"
)

func testServer() http.Handler {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Store: graph.SeedStore()}}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})
	return srv
}

func TestProjectsQuery(t *testing.T) {
	srv := testServer()

	body := `{"query":"query { projects { id name } }"}`
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "projects") {
		t.Fatalf("response does not contain projects: %s", rec.Body.String())
	}
}

func TestCreateProjectMutation(t *testing.T) {
	srv := testServer()

	body := `{"query":"mutation { createProject(input: {name: \"Demo\", description: \"Test\"}) { id name } }"}`
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "Demo") {
		t.Fatalf("project was not created: %s", rec.Body.String())
	}
}

func TestInvalidFieldQuery(t *testing.T) {
	srv := testServer()

	body := `{"query":"query { projects { unknownField } }"}`
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "errors") {
		t.Fatalf("expected validation error: %s", rec.Body.String())
	}
}

func TestAddTaskToUnknownProject(t *testing.T) {
	srv := testServer()

	body := `{"query":"mutation { addTask(input: {projectId: \"missing\", title: \"Task\", status: NEW}) { id title } }"}`
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)

	if !strings.Contains(rec.Body.String(), "errors") {
		t.Fatalf("expected domain error: %s", rec.Body.String())
	}
}
