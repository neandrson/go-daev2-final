package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/neandrson/go-daev2-final/orchestrator/internal/services/auth"
	"github.com/neandrson/go-daev2-final/orchestrator/internal/services/expression"
)

type CalcHandler struct {
	expressionService *expression.ExpressionService
}

func NewCalcHandler(expressionService *expression.ExpressionService) *CalcHandler {
	return &CalcHandler{
		expressionService: expressionService,
	}
}

func (h *CalcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	var request struct {
		Expression string `json:"expression"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}
	user_id := r.Context().Value(auth.ContextKeyUserID).(int)
	// Логика спрятана сюда
	id, err := h.expressionService.ProcessExpression(request.Expression, user_id)

	if err != nil {
		if errors.Is(err, expression.ErrStorage) || errors.Is(err, expression.ErrService) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"id": id})
}
