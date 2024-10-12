package server

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	//"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

func parseIP(raw string) string {
	logger := slog.Default().WithGroup("parseIp").With("input", raw)
	host, _, err := net.SplitHostPort(raw)
	if err != nil {
		slog.Warn("error spliting host and port", "error", err)
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip.String()
	}
	logger.Warn("Returning blank")
	return ""
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(s.httpLogger))
	//r.Use(middleware.Logger)

	r.Get("/", s.WhoAmIHandler)
	r.Get("/health", s.HealthCheckHandler)

	return r
}

func (s *Server) WhoAmIHandler(w http.ResponseWriter, r *http.Request) {
	logger := httplog.LogEntry(r.Context())
	rawIP := parseIP(r.RemoteAddr)
	logger.Info("WhoAmI Request", "rawIP", rawIP)
	headers := r.Header
	logger.Info("IP Headers",
		"X-Forwarded-For", headers.Get("X-Forwarded-For"),
		"Forwarded", headers.Get("Forwarded"),
		"X-Real-Ip", headers.Get("X-Real-Ip"),
		"Via", headers.Get("Via"),
	)

	if _, err := w.Write([]byte(fmt.Sprintln(rawIP))); err != nil {
		logger.Error("ResponseWriter.Write failed", "err", err)
	}
}

func (s *Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	logger := httplog.LogEntry(r.Context())
	fixedResponse := []byte("OK\n")

	_, err := w.Write(fixedResponse)
	if err != nil {
		logger.Error("ResponseWriter.Write failed", "err", err)
	}
}
