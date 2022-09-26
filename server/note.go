package server

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/viveknathani/nattukaka/entity"
	"github.com/viveknathani/nattukaka/shared"
)

type notesPageVariables struct {
	NotesList *[]entity.Note
}

func (s *Server) handleNotesIndex(w http.ResponseWriter, r *http.Request) {

	indexFilePath := "static/pages/note.html"
	t, err := template.ParseFiles(indexFilePath)
	if err != nil {
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}

	list, err := s.Service.GetAllNotes(r.Context(), shared.ExtractUserID(r.Context()))
	if err != nil {
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}

	err = t.Execute(w, notesPageVariables{
		NotesList: list,
	})

	if err != nil {
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}
}

func (s *Server) handleNoteContent(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	id := params["id"][0]
	list, err := s.Service.GetNote(r.Context(), id, shared.ExtractUserID(r.Context()))
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

func (s *Server) handleNoteCreate(w http.ResponseWriter, r *http.Request) {

	var n noteCreateRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&n)
	if err != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
		if ok := sendClientError(w, err.Error()); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}

	err = s.Service.CreateNote(r.Context(), &entity.Note{
		UserId:  shared.ExtractUserID(r.Context()),
		Id:      "",
		Content: "",
		Title:   n.Title,
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

func (s *Server) handleNoteUpdate(w http.ResponseWriter, r *http.Request) {

	var n noteUpdateRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&n)
	if err != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
		if ok := sendClientError(w, err.Error()); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}

	err = s.Service.UpdateNote(r.Context(), &entity.Note{
		UserId:  shared.ExtractUserID(r.Context()),
		Id:      n.Id,
		Content: n.Content,
		Title:   n.Title,
	})

	if err != nil {
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}

	if ok := sendUpdated(w); ok != nil {
		s.Service.Logger.Error(ok.Error(), zapReqID(r))
	}
	showRequestEnd(s.Service.Logger, r)
}
