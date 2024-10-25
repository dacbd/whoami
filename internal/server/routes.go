package server

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	//"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

func parseIP(logger *slog.Logger, raw string) string {
	host, _, err := net.SplitHostPort(raw)
	if err != nil {
		logger.Warn("error spliting host and port", "error", err)
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip.String()
	}
	logger.Warn("Returning blank")
	return ""
}

func parseForwardFor(logger *slog.Logger, raw string) string {
	vals := strings.Split(raw, ",")
	logger.Debug("parseForwardFor", "len", len(vals))
	for _, v := range vals {
		parsedIP := parseIP(logger, strings.Trim(v, " "))
		if parsedIP != "" {
			return parsedIP
		}
	}
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
	rawIP := parseIP(logger, r.RemoteAddr)
	logger.Info("WhoAmI Request", "rawIP", rawIP)
	headers := r.Header
	logger.Info("IP Headers",
		"X-Forwarded-For", headers.Get("X-Forwarded-For"),
		"Forwarded", headers.Get("Forwarded"),
		"X-Real-Ip", headers.Get("X-Real-Ip"),
		"Via", headers.Get("Via"),
	)

	forwardedIP := parseForwardFor(logger, headers.Get("X-Forwarded-For"))
	var realIP string
	if forwardedIP != "" {
		realIP = forwardedIP
	} else {
		realIP = rawIP
	}

	if _, err := w.Write([]byte(fmt.Sprintln(realIP))); err != nil {
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
