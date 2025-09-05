package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Fedasov/Effective-Mobile/internal/model"
	"github.com/Fedasov/Effective-Mobile/internal/service"

	"github.com/gorilla/mux"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
}

func NewSubscriptionHandler(service service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// Create обрабатывает запрос на создание подписки
// @Summary Создать новую подписку
// @Description Создает новую запись о подписке пользователя
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body model.SubscriptionCreateRequest true "Данные подписки"
// @Success 201 {object} model.Subscription
// @Failure 400 {object} map[string]string "Неверный формат данных"
// @Router /api/v1/subscriptions [post]
func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.SubscriptionCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subscription, err := h.service.Create(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subscription)
}

// GetByID обрабатывает запрос на получение подписки по ID
// @Summary Получить подписку по ID
// @Description Возвращает информацию о подписке по её идентификатору
// @Tags subscriptions
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Подписка не найдена"
// @Router /api/v1/subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	subscription, err := h.service.GetByID(uint32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)
}

// Update обрабатывает запрос на обновление подписки
// @Summary Обновить подписку
// @Description Обновляет информацию о существующей подписке
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "ID подписки"
// @Param input body model.SubscriptionCreateRequest true "Новые данные подписки"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Подписка не найдена"
// @Router /api/v1/subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req model.SubscriptionCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subscription, err := h.service.Update(uint32(id), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)
}

// Delete обрабатывает запрос на удаление подписки
// @Summary Удалить подписку
// @Description Удаляет запись о подписке по её идентификатору
// @Tags subscriptions
// @Param id path int true "ID подписки"
// @Success 204 "Подписка успешно удалена"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Подписка не найдена"
// @Router /api/v1/subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(uint32(id)); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List обрабатывает запрос на получение списка подписок
// @Summary Получить список подписок
// @Description Возвращает список подписок с поддержкой пагинации
// @Tags subscriptions
// @Produce json
// @Param limit query int false "Лимит записей (по умолчанию 10)"
// @Param offset query int false "Смещение (по умолчанию 0)"
// @Success 200 {array} model.Subscription
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/v1/subscriptions [get]
func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	subscriptions, err := h.service.List(int32(limit), int32(offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscriptions)
}

// GetTotalCost обрабатывает запрос на расчет общей стоимости
// @Summary Рассчитать общую стоимость
// @Description Вычисляет общую стоимость подписок за указанный период с возможностью фильтрации
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body model.TotalCostRequest true "Параметры расчета"
// @Success 200 {object} map[string]int "Общая стоимость"
// @Failure 400 {object} map[string]string "Неверные параметры"
// @Router /api/v1/subscriptions/total-cost [post]
func (h *SubscriptionHandler) GetTotalCost(w http.ResponseWriter, r *http.Request) {
	var req model.TotalCostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	total, err := h.service.CalculateTotalCost(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int32{"total_cost": total})
}
