package httpkaka

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/viveknathani/nattukaka/entity"
	"github.com/viveknathani/nattukaka/service"
	"github.com/viveknathani/nattukaka/shared"
)

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var u userLoginRequest
	err := decoder.Decode(&u)

	if err != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
		if ok := sendClientError(w, err.Error()); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}

	token, err := s.Service.Login(r.Context(), &entity.User{
		Email:    u.Email,
		Password: []byte(u.Password),
	})
	if err != nil {

		s.Service.Logger.Error(err.Error(), zapReqID(r))
		switch {
		case err == service.ErrInvalidEmailPassword:
			{
				if ok := sendClientError(w, err.Error()); ok != nil {
					s.Service.Logger.Error(ok.Error(), zapReqID(r))
				}
				return
			}
		default:
			{
				if ok := sendServerError(w); ok != nil {
					s.Service.Logger.Error(ok.Error(), zapReqID(r))
				}
				return
			}
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		MaxAge:   int(time.Hour.Seconds() * 24 * 3),
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   (os.Getenv("MODE") == "prod"),
	})

	if ok := sendResponse(w, "ok", http.StatusOK); ok != nil {
		s.Service.Logger.Error(ok.Error(), zapReqID(r))
	}
	showRequestEnd(s.Service.Logger, r)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("token")
	if err != nil {
		s.Service.Logger.Error(err.Error(), zapReqID(r))
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}

	err = s.Service.Logout(r.Context(), cookie.Value)
	if err != nil {

		s.Service.Logger.Error(err.Error())
		if ok := sendServerError(w); ok != nil {
			s.Service.Logger.Error(ok.Error(), zapReqID(r))
		}
		return
	}

	if ok := sendResponse(w, "ok", http.StatusOK); ok != nil {
		s.Service.Logger.Error(ok.Error(), zapReqID(r))
	}
	showRequestEnd(s.Service.Logger, r)
}

func (s *Server) middlewareTokenVerification(handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			s.Service.Logger.Error(err.Error(), zapReqID(r))
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		id, err := s.Service.VerifyAndDecodeToken(r.Context(), cookie.Value)
		if err != nil {

			s.Service.Logger.Error(err.Error(), zapReqID(r))
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		request := r.Clone(shared.WithUserID(r.Context(), id))
		handler.ServeHTTP(w, request)
	}
}

func (s *Server) serveLogin(w http.ResponseWriter, r *http.Request) {

	p, err := template.ParseFiles("static/pages/login.html")
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
