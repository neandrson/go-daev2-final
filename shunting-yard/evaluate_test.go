package shuntingYard

import (
	"reflect"
	"testing"
)

func TestEvaluate(t *testing.T) {
	cases := []struct {
		name    string
		tokens  []*RPNToken
		want    float64
		wantErr bool
	}{
		{
			"nil []*RPNToken",
			nil,
			0,
			true,
		},
		{
			"empty []*RPNToken",
			[]*RPNToken{},
			0,
			true,
		},
		{
			"one operand",
			[]*RPNToken{NewRPNOperandToken(1)},
			1,
			false,
		},
		{
			"one operator",
			[]*RPNToken{NewRPNOperatorToken("+")},
			0,
			true,
		},
		{
			"unknown operator",
			[]*RPNToken{NewRPNOperandToken(1), NewRPNOperandToken(2), NewRPNOperatorToken("a")},
			0,
			true,
		},
		{
			"expression without error",
			[]*RPNToken{NewRPNOperandToken(1), NewRPNOperandToken(2), NewRPNOperatorToken("^"), NewRPNOperandToken(2), NewRPNOperandToken(2), NewRPNOperatorToken("+"), NewRPNOperandToken(2), NewRPNOperandToken(4), NewRPNOperatorToken("*"), NewRPNOperatorToken("/"), NewRPNOperatorToken("-")},
			0.5,
			false,
		},
		{
			"expression with error",
			[]*RPNToken{NewRPNOperandToken(1), NewRPNOperatorToken("^"), NewRPNOperandToken(2), NewRPNOperandToken(2), NewRPNOperatorToken("+"), NewRPNOperandToken(2), NewRPNOperandToken(4), NewRPNOperatorToken("*"), NewRPNOperatorToken("/"), NewRPNOperatorToken("-")},
			0,
			true,
		},
	}

	for _, tc := range cases {
		tc := tc
		// запуск отдельного теста
		t.Run(tc.name, func(t *testing.T) {
			// тестируем функцию Evaluate
			got, err := Evaluate(tc.tokens)
			// проверим полученное значение
			if (err != nil && !tc.wantErr) || (err == nil && tc.wantErr) {
				t.Errorf("got err: %v; want err: %v", err != nil, tc.wantErr)
			} else if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v; want %v", got, tc.want)
			}
		})
	}
}
