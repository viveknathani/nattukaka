package httpkaka

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/viveknathani/nattukaka/entity"
	"github.com/viveknathani/nattukaka/shared"
)

func (s *Server) handleTodoPending(w http.ResponseWriter, r *http.Request) {

	list, err := s.Service.GetAllPendingTodos(r.Context(), shared.ExtractUserID(r.Context()))
	if err != nil {
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(err.Error(), zapReqID(r))
		}
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, ok := w.Write(data); ok != nil {
		s.Service.Logger.Error(ok.Error(), zapReqID(r))
	}
	showRequestEnd(s.Service.Logger, r)
}

func (s *Server) handleTodoCreate(w http.ResponseWriter, r *http.Request) {

	var t todoCreateRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)
	if err != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
		if ok := sendClientError(w, err.Error()); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}

	err = s.Service.CreateTodo(r.Context(), &entity.Todo{
		UserId:      shared.ExtractUserID(r.Context()),
		Task:        t.Task,
		Status:      "pending",
		Deadline:    &t.Deadline,
		CompletedAt: nil,
	})

	if err != nil {
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}

	if ok := sendCreated(w); ok != nil {
		s.Service.Logger.Error(ok.Error(), zapReqID(r))
	}
	showRequestEnd(s.Service.Logger, r)
}

func (s *Server) handleTodoUpdate(w http.ResponseWriter, r *http.Request) {

	var t todoUpdateRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)
	if err != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
		if ok := sendClientError(w, err.Error()); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}

	fmt.Println("completedAt is: ", t.CompletedAt)
	err = s.Service.UpdateTodo(r.Context(), &entity.Todo{
		UserId:      shared.ExtractUserID(r.Context()),
		Id:          t.Id,
		Task:        t.Task,
		Status:      t.Status,
		Deadline:    &t.Deadline,
		CompletedAt: &t.CompletedAt,
	})

	if err != nil {
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}

	if ok := sendCreated(w); ok != nil {
		s.Service.Logger.Error(ok.Error(), zapReqID(r))
	}
	showRequestEnd(s.Service.Logger, r)
}
func (s *Server) handleTodoDelete(w http.ResponseWriter, r *http.Request) {

	var t todoDeleteRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)
	if err != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
		if ok := sendClientError(w, err.Error()); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}

	err = s.Service.DeleteTodo(r.Context(), t.Id)

	if err != nil {
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(err.Error(), zapReqID(r))
		}
		return
	}

	if ok := sendUpdated(w); ok != nil {
		s.Service.Logger.Error(ok.Error(), zapReqID(r))
	}
	showRequestEnd(s.Service.Logger, r)
}

func (s *Server) serveTodo(w http.ResponseWriter, r *http.Request) {

	p, err := template.ParseFiles("static/pages/todo.html")
	if err != nil {
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}
	err = p.Execute(w, nil)
	if err != nil {
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}
}
