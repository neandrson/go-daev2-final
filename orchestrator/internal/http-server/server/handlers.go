package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	model "github.com/neandrson/go-daev2-final/orchestrator/internal/domain/models"
	expressionparser "github.com/neandrson/go-daev2-final/orchestrator/internal/lib/expressionParser"
)

type myRequest struct {
	Expression string `json:"expression"`
	Id         string `json:"-"`
}

func (s *Server) MainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, use curl :)")
}

func (s *Server) GetExpressionById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	expression, err := s.httpService.GetExpressionById(context.Background(), id, r.Context().Value("uid").(int))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.Write([]byte("expression not found!"))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(
		[]byte(
			fmt.Sprintf("expression: %s\nstatus: %s\nresult: %f\n\n", expression.InfinixExpression, expression.Status, expression.Result), //created_at: %v\nsolved_at: %v\n
		), //, expression.CreatedAt, expression.SolvedAt
	)

	//fmt.Fprintln(w, "Expression doesn't exist")
}

func (s *Server) EvaluateExpression(w http.ResponseWriter, r *http.Request) {
	expression := r.Context().Value("expression").(model.Expression)

	// // Проверяем выражение на наличие результата в базе данных (в ином случае отправляем агенту на вычисление)
	// // if el, ok := cache.Get(expression.IdExpression); ok {
	// // 	fmt.Fprintln(w, el.Result)
	// // 	return
	// // }

	// expression, err := orchestrator.SolveExpression(&expression)
	// if err != nil {
	// 	log.Println(err.Error())
	// }
	// // записываем в бд

	expression.IdExpression = fmt.Sprintf("%d", time.Now().UnixNano())

	id, err := s.httpService.EvaluateExpression(context.Background(), &expression, r.Context().Value("uid").(int)) //result, err
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Result is: %v", id))) //result
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := s.httpService.Login(context.Background(), request.Login, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Your authorization token: " + token))
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := s.httpService.Register(context.Background(), request.Login, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("Your uid: %d", id)))
}

func (s *Server) GetAgentStates(w http.ResponseWriter, r *http.Request) {
	agents, err := s.httpService.GetAgentStates(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, a := range agents {
		w.Write([]byte(fmt.Sprintf("id: %d\nlast_heartbeat: %v\nstatus: %s\n\n", a.ID, a.LastHeartbeat, a.Status)))
	}
}

func (s *Server) GetExpressionsForUser(w http.ResponseWriter, r *http.Request) {
	expressions, err := s.httpService.GetExpressionsForUser(context.Background(), r.Context().Value("uid").(int))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.Write([]byte("No expressions for you"))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(expressions) == 0 {
		w.Write([]byte("No expressions for you"))
		return
	}
	for _, e := range expressions {
		w.Write(
			[]byte(
				fmt.Sprintf("id: %s\nstatus: %s\nresult: %f\n\n", e.IdExpression, e.Status, e.Result), //expression: %s\ncreated_at: %v\nsolved_at: %v\n
			), //, e.InfinixExpression, e.CreatedAt, e.SolvedAt
		)
	}
}

// вообще по идее оно должно создаваться на фронтэнде, но т.к пока нет фронта - создаём на бэкенде (по запросу с фронта)
func (s *Server) GetImpodenceKeyHandler(w http.ResponseWriter, r *http.Request) {
	// Парсим JSON
	var request myRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Error while parsing JSON", http.StatusBadRequest)
		return
	}

	request.Expression = strings.ReplaceAll(request.Expression, " ", "")

	// и получаем ключ
	key := expressionparser.CreateImpodenceKey(request.Expression)

	w.Write([]byte(key))
}
