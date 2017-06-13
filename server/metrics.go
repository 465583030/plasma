package server

import (
	"net/http"

	"github.com/openfresh/plasma/config"
	"github.com/openfresh/plasma/metrics"
	"go.uber.org/zap"
)

type metricsHandler struct {
	accessLogger *zap.Logger
	errorLogger  *zap.Logger
	config       config.Config
	mux          *http.ServeMux
}

func (h metricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func NewMetricsHandler(opt Option) metricsHandler {

	h := metricsHandler{
		accessLogger: opt.AccessLogger,
		errorLogger:  opt.ErrorLogger,
		config:       opt.Config,
		mux:          http.NewServeMux(),
	}

	h.mux.HandleFunc("/go", h.metrics)
	return h
}

func (h *metricsHandler) metrics(w http.ResponseWriter, r *http.Request) {
	metrics.HTTPHandler(w, r)
}
