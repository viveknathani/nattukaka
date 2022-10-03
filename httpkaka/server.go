package httpkaka

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/viveknathani/nattukaka/service"
	"github.com/viveknathani/nattukaka/shared"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Server holds together all the configuration needed to run this web service.
type Server struct {
	Service *service.Service
	Router  *mux.Router
}

// ServeHTTP is implemented so that Server can be used for listening to requests.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	requestID := uuid.New().String()
	request := r.Clone(shared.WithRequestID(context.Background(), requestID))
	showRequestMetaData(s.Service.Logger, request)
	s.Router.ServeHTTP(w, request)
}

func zapReqID(r *http.Request) zapcore.Field {

	return zapcore.Field{
		Key:    "requestID",
		String: shared.ExtractRequestID(r.Context()),
		Type:   zapcore.StringType,
	}
}

func showRequestMetaData(l *zap.Logger, r *http.Request) {

	reqMethod := zapcore.Field{
		Key:    "method",
		String: r.Method,
		Type:   zapcore.StringType,
	}

	reqPath := zapcore.Field{
		Key:    "path",
		String: r.URL.String(),
		Type:   zapcore.StringType,
	}

	ip := zapcore.Field{
		Key:    "ip",
		String: getIP(r),
		Type:   zapcore.StringType,
	}

	l.Info("incoming request", zapReqID(r), reqMethod, reqPath, ip)
}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func showRequestEnd(l *zap.Logger, r *http.Request) {

	l.Info("completed processing", zapReqID(r))
}
