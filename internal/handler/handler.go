package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"onlineusersub/internal/domain"
	"onlineusersub/internal/service"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type SubHandler struct {
	service service.SubService
}

func NewSubHandler(service service.SubService) *SubHandler {
	return &SubHandler{service: service}
}

func (s *SubHandler) Routes() chi.Router {

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/subscription", func(r chi.Router) {
		r.Post("/publish", s.SubPublish)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", s.SubGetByID)
			r.Patch("/", s.UpdateSubByID)
			r.Delete("/", s.DeleteByID)
		})
	})

	r.Route("/subscriptions", func(r chi.Router) {
		r.Get("/", s.SubsGetAll)
		r.Get("/sum", s.GetTotalSum)
	})

	return r
}

// Добавление записи в БД
func (s *SubHandler) SubPublish(w http.ResponseWriter, r *http.Request) {

	var req struct {
		UserID         uuid.UUID `json:"user_id"`
		ServiceName    string    `json:"service_name"`
		Price          int       `json:"price"`
		DateCreated    string    `json:"start_date"`
		DateConclusion string    `json:"conclusion_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	layout := "01-2006"

	startDate, err := time.Parse(layout, req.DateCreated)
	if err != nil {
		writeError(w, "invalid JSON date format", http.StatusBadRequest)
		return
	}

	var endDate time.Time
	if req.DateConclusion != "" {
		endDate, err = time.Parse(layout, req.DateConclusion)
		if err != nil {
			writeError(w, "invalid JSON date format", http.StatusBadRequest)
			return
		}
	}

	subscription := domain.Subscriptions{
		UserID:         req.UserID,
		ServiceName:    req.ServiceName,
		Price:          req.Price,
		DateCreated:    startDate,
		DateConclusion: endDate,
	}

	sub, err := s.service.CreateSub(r.Context(), subscription)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, sub, http.StatusCreated)

}

// Возвращение подписки из БД
func (s *SubHandler) SubGetByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	if id == "" {
		writeError(w, "UID подписки не указан", http.StatusBadRequest)
		return
	}

	sub, err := s.service.GetSub(r.Context(), id)
	if err != nil {

		if errors.Is(err, service.ErrSubNotFound) {
			writeError(w, "Подписка не найдена", http.StatusNotFound)
			return
		}

		slog.Info("(handler) ошибка GetSub", "id", "err", id, err)
		writeError(w, "Внутреняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	writeJSON(w, sub, http.StatusOK)

}

// Возвращение всех подписок
func (s *SubHandler) SubsGetAll(w http.ResponseWriter, r *http.Request) {

	subs, err := s.service.GetAllSub(r.Context())
	if err != nil {
		slog.Info("(handler) ошибка SubsGetAll", "err", err)
		writeError(w, "Ошибка получения подписок", http.StatusInternalServerError)
		return
	}

	if subs == nil {
		subs = []domain.Subscriptions{}
	}

	writeJSON(w, subs, http.StatusOK)
}

// Обновление данных
func (s *SubHandler) UpdateSubByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	uid, err := uuid.Parse(id)
	if err != nil {
		writeError(w, "Невалидный ID", http.StatusBadRequest)
		return
	}

	var req struct {
		UserID         uuid.UUID `json:"user_id"`
		ServiceName    string    `json:"service_name"`
		Price          int       `json:"price"`
		DateCreated    string    `json:"start_date"`
		DateConclusion string    `json:"conclusion_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	layout := "01-2006"

	startDate, err := time.Parse(layout, req.DateCreated)
	if err != nil {
		writeError(w, "invalid JSON date format", http.StatusBadRequest)
		return
	}

	var endDate time.Time
	if req.DateConclusion != "" {
		endDate, err = time.Parse(layout, req.DateConclusion)
		if err != nil {
			writeError(w, "invalid JSON date format", http.StatusBadRequest)
		}
		return
	}

	subscription := domain.Subscriptions{
		ID:             uid,
		UserID:         req.UserID,
		ServiceName:    req.ServiceName,
		Price:          req.Price,
		DateCreated:    startDate,
		DateConclusion: endDate,
	}

	if err := s.service.UpdateSub(r.Context(), subscription); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, subscription, http.StatusOK)

}

// Удаление записи по ID
func (s *SubHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	err := s.service.DeleteSub(r.Context(), id)
	if err != nil {

		if errors.Is(err, service.ErrSubNotFound) {
			writeError(w, "Подписка не найдена", http.StatusNotFound)
			return
		}

		slog.Info("(handler) ошибка DeleteByID", "id", id, "err", err)
		writeError(w, "Внутреняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Возвращение суммы всех подписок пользователя
func (s *SubHandler) GetTotalSum(w http.ResponseWriter, r *http.Request) {

	uID := r.URL.Query().Get("user_id")
	srv := r.URL.Query().Get("service")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	layout := "01-2006"

	startDate, err := time.Parse(layout, from)
	if err != nil {
		writeError(w, "invalid JSON date format", http.StatusBadRequest)
		return
	}

	var endDate time.Time
	if to != "" {
		endDate, err = time.Parse(layout, to)
		if err != nil {
			writeError(w, "invalid JSON date format", http.StatusBadRequest)
			return
		}

	} else {
		endDate = time.Now()
	}

	sum, err := s.service.GetSumSub(r.Context(), uID, srv, startDate, endDate)
	if err != nil {
		writeError(w, "Invalid URL Query", http.StatusBadRequest)
		return
	}

	writeJSON(w, sum, http.StatusOK)
}

// Вспомогатльеные фукнции

// Сериализация данных в JSON и отправка
func writeJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Отправка JSON ответа с ошибкой
func writeError(w http.ResponseWriter, message string, status int) {
	writeJSON(w, map[string]string{"error": message}, status)
}
