package handlers

import (
	"context"
	"fmt"
	"mentor-training/pkg/logger"
	"mentor-training/pkg/queue"
	"net/http"
	"strings"
	"time"
)

// Handler представляет структуру для обработки запросов
type Handler struct {
	log logger.Logger
	qm  queue.QueueManager
}

// NewHandler
func NewHandler(log logger.Logger, qm queue.QueueManager) *Handler {
	return &Handler{
		log: log,
		qm:  qm,
	}
}

// processQueueName извлекает имя очереди из пути запроса
func processQueueName(r *http.Request) (string, error) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 1 || pathParts[0] == "" {
		return "", fmt.Errorf("имя очереди обязательно")
	}
	return pathParts[0], nil
}

// handlePut обрабатывает PUT запросы
func (h *Handler) handlePut(queueName string, w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("v")
	if message == "" {
		h.log.Error("Отсутствует параметр v")
		http.Error(w, "параметр v обязателен", http.StatusBadRequest)
		return
	}

	h.qm.Put(queueName, message)
	h.log.Infof("Сообщение добавлено в очередь %s: %s", queueName, message)
	w.WriteHeader(http.StatusOK)
}

// handleGet обрабатывает GET запросы
func (h *Handler) handleGet(ctx context.Context, queueName string, w http.ResponseWriter, r *http.Request) {
	timeout := time.Second * 0
	timeoutStr := r.URL.Query().Get("timeout")
	if timeoutStr != "" {
		parsedTimeout, err := time.ParseDuration(timeoutStr + "s")
		if err != nil {
			h.log.Errorf("Некорректный параметр timeout: %s", timeoutStr)
			http.Error(w, "некорректный параметр timeout", http.StatusBadRequest)
			return
		}
		timeout = parsedTimeout
	}

	message, exists := h.qm.Get(ctx, queueName, timeout)
	if !exists {
		h.log.Errorf("Сообщение для очереди %s не найдено в течение %s", queueName, timeout)
		http.Error(w, "сообщение не найдено", http.StatusNotFound)
		return
	}

	h.log.Infof("Сообщение из очереди %s: %s", queueName, message)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

// ServeHTTP обрабатывает запросы сервиса
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	queueName, err := processQueueName(r)
	if err != nil {
		h.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.log.Infof("Получен запрос %s %s", r.Method, r.URL.Path)

	switch r.Method {
	case http.MethodPut:
		h.handlePut(queueName, w, r)
	case http.MethodGet:
		h.handleGet(ctx, queueName, w, r)
	default:
		h.log.Errorf("Получен неподдерживаемый метод: %s", r.Method)
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
	}
}
