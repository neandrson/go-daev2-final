package shuntingYard

import (
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name       string
		expression interface{}
		want       []*RPNToken
		wantErr    bool
	}{
		{
			name:       "empty string",
			expression: "",
			want:       []*RPNToken{},
			wantErr:    false,
		},
		{
			name:       "correct expression",
			expression: "1+2*(2-3)^3",
			want:       []*RPNToken{NewRPNOperandToken(1), NewRPNOperandToken(2), NewRPNOperandToken(2), NewRPNOperandToken(3), NewRPNOperatorToken("-"), NewRPNOperandToken(3), NewRPNOperatorToken("^"), NewRPNOperatorToken("*"), NewRPNOperatorToken("+")},
			wantErr:    false,
		},
		{
			name:       "incorrect brackets (opening)",
			expression: "1+(2",
			want:       []*RPNToken{},
			wantErr:    true,
		},
		{
			name:       "incorrect bracket (closing)",
			expression: "1+2)",
			want:       []*RPNToken{},
			wantErr:    true,
		},
		{
			name:       "unknown operator",
			expression: []string{"hello world"},
			want:       []*RPNToken{},
			wantErr:    true,
		},
	}

	for _, tc := range cases {
		tc := tc
		// запуск отдельного теста
		t.Run(tc.name, func(t *testing.T) {
			// тестируем функцию Evaluate
			var tokens []string
			switch tc.expression.(type) {
			case string:
				tokens, _ = Scan(tc.expression.(string))
			case []string:
				tokens = tc.expression.([]string)
			}
			got, err := Parse(tokens)
			// проверим полученное значение
			if (err != nil && !tc.wantErr) || (err == nil && tc.wantErr) {
				t.Errorf("got err: %v; want err: %v", err != nil, tc.wantErr)
			} else if !Equal(got, tc.want) {
				t.Errorf("got %v; want %v", got, tc.want)
			}
		})
	}
}

func Equal(slice1, slice2 []*RPNToken) (isEqual bool) {
	isEqual = true
	if len(slice1) != len(slice2) {
		isEqual = false
		return isEqual
	}
	for i := 0; i < len(slice1); i++ {
		if (slice1[i].Value != slice2[i].Value) && (slice1[i].Type != slice2[i].Type) {
			isEqual = false
			return isEqual
		}
	}
	return isEqual
}
