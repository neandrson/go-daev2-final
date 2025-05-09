package agent

import (
	"log"
	"time"

	"github.com/neandrson/go-daev2-final/agent/internal/domain/models"
	shuntingYard "github.com/neandrson/go-daev2-final/shunting-yard"
)

type Calculator struct {
	taskChan   chan *models.ExpressionPart
	expression *models.ExpressionPart
	isBusy     bool
	id         int
}

func NewCalculator(i int) *Calculator {
	c := &Calculator{
		taskChan: make(chan *models.ExpressionPart),
		isBusy:   false,
		id:       i,
	}

	c.Start()

	return c
}

func (c *Calculator) AddTask(task *models.ExpressionPart) bool {
	if c.isBusy {
		return false
	}

	c.taskChan <- task
	return true
}

func (c *Calculator) Start() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
				c.Start()
			}
		}()

		for {
			task := <-c.taskChan
			c.expression = task
			c.isBusy = true

			log.Printf("Calculator[%v]: got task to solve", c.id)

			c.SolveExpression(task)

			log.Printf("Calculator[%v]: task solved", c.id)

			c.expression = nil
			c.isBusy = false
		}
	}()
}

func (c *Calculator) SolveExpression(expr *models.ExpressionPart) {
	time.Sleep(expr.Duration)

	if result, err := shuntingYard.Evaluate([]*shuntingYard.RPNToken{expr.FirstOperand, expr.SecondOperand, expr.Operation}); err == nil {
		tokenizedResult := shuntingYard.NewRPNOperandToken(result)
		expr.Result <- tokenizedResult
		return
	}
	expr.Result <- shuntingYard.NewRPNOperandToken(0)
}
