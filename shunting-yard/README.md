# go-shunting-yard

A simple implementation of shunting-yard algorithm in Go(Golang).

## Installation
```sh
go get github.com/neandrson/go-shunting-yard
```

## Supported Operators

| Name | Description | Precedence | Associativity |
|------|-------------|------------|---------------|
| ^  | Exponentiation | 3 | Left-associative |
| * /  | Multiplication, Division | 2 | Left-associative |
| + -  | Addition, Subtraction    | 1 | Left-associative |


## Example
```go
import (
	"fmt"

	"github.com/neandrson/go-shunting-yard"
)

func main() {
	// calculate 12 + 4 * 3 / 5 - 2
	input := "12 + 4 * 3 / 5 - 2"

	// parse input expression to infix notation
	infixTokens, err := shuntingYard.Scan(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("Infix Tokens:")
	fmt.Println(infixTokens)

	// convert infix notation to postfix notation(RPN)
	postfixTokens, err := shuntingYard.Parse(infixTokens)
	if err != nil {
		panic(err)
	}
	fmt.Println("Postfix(RPN) Tokens:")
	for _, t := range postfixTokens {
		fmt.Printf("%v ", t.GetDescription())
	}
	fmt.Println()

	// evaluate RPN tokens
	result, err := shuntingYard.Evaluate(postfixTokens)
	if err != nil {
		panic(err)
	}

	// output the result
	fmt.Printf("Result: %v", result)
}
```

Output:
```
Infix Tokens:
[12 + 4 * 3 / 5 - 2]
Postfix(RPN) Tokens:
(1)12 (1)4 (1)3 (2)* (1)5 (2)/ (2)+ (1)2 (2)-
Result: 12
```
