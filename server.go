package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	titlePattern = "/title/"
)

var CacheMissError = errors.New("Cache Miss Error")

type APIHandler func(w http.ResponseWriter, r *http.Request) error

type APIError struct {
	Error string `json:"error"`
}

func NewAPIError(err string) *APIError {
	return &APIError{
		Error: err,
	}
}

func (a APIHandler) decorate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := a(w, r); err != nil {
			err := writeJSON(w, http.StatusBadRequest, *NewAPIError(err.Error()))
			check(err)
		}
	}
}

type Server struct {
	addr        string
	cacheClient CacheClient
	dbClient    DbClient
}

func NewServer(addr string, cacheClient CacheClient, dbClient DbClient) *Server {
	return &Server{
		addr:        addr,
		cacheClient: cacheClient,
		dbClient:    dbClient,
	}
}

func (s *Server) start() {
	http.HandleFunc(titlePattern, APIHandler(s.handleRequest).decorate())

	fmt.Println(http.ListenAndServe(os.Getenv("SERVER_PORT"), nil))
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" || strings.ToLower(contentType) != "application/json" {
		return fmt.Errorf("Expected %q to be %q. But received %q", "Content-Type", "application/json", contentType)
	}

	switch r.Method {
	case "GET":
		return s.handleGetTitle(w, r)
	default:
		return fmt.Errorf("Unsupported method: %q", r.Method)
	}
}

func (s *Server) handleGetTitle(w http.ResponseWriter, r *http.Request) error {
	url := r.URL.Path
	id, err := parseID(url)
	if err != nil {
		return err
	}

	title, err := s.cacheClient.Get(id)
	if err == nil {
		fmt.Println("Cache Hit")
		return writeJSON(w, http.StatusOK, *title)
	} else if err != CacheMissError {
		return err
	}
	fmt.Println("Cache Miss")

	title, err = s.dbClient.Get(id)
	if err != nil {
		return err
	}

	err = s.cacheClient.Set(id, title)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, *title)
}

func writeJSON[V TitleBasics | APIError](w http.ResponseWriter, statusCode int, data V) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	return err
}

func parseID(url string) (string, error) {
	id := strings.TrimSpace(strings.Replace(url, titlePattern, "", 1))
	if isEmpty(id) {
		return id, fmt.Errorf("id cannot be empty")
	}
	return id, nil
}

func isEmpty(s string) bool {
	return s == ""
}
