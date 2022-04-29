package service

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SplitExpression(str string) ([]string, error) {

	for strings.Contains(str, "sqrt(") {
		sqrFront := strings.Index(str, "sqrt(") + len("sqrt")
		sqrBack := getMatchBracket(str, sqrFront)
		if sqrBack == -1 {
			return nil, fiber.ErrBadRequest
		}
		innerSqrt := str[sqrFront+1 : sqrBack]
		sqrtWhole := str[sqrFront-len("sqrt") : sqrBack+1]

		innerSplit, err := SplitExpression(innerSqrt)
		if err != nil {
			return nil, err
		}

		sqrtRes, err := Evaluate(innerSplit)
		if err != nil {
			return nil, err
		}

		res, err := strconv.ParseFloat(sqrtRes, 64)
		if err != nil {
			return nil, err
		}
		res = math.Sqrt(res)
		sqrtRes = fmt.Sprintf("%f", res)

		str = strings.Replace(str, sqrtWhole, sqrtRes, 1)
	}

	splitted := strings.FieldsFunc(str, split)
	if splitted[0] == "-" {
		temp := make([]string, len(splitted)+1)
		temp[0] = "0"
		copy(temp[1:], splitted)
		splitted = temp
	}

	for i := 1; i < len(splitted); i++ {
		if splitted[i] == "-" && splitted[i-1] == "(" {
			splitted[i] = ""
			splitted[i+1] = "-" + splitted[i+1]
		}
	}

	return splitted, nil
}

func Evaluate(tokens []string) (string, error) {

	var val Stack
	var ops Stack

	containsVal := false
	for i := 0; i < len(tokens); i++ {
		if !isOperator(tokens[i]) {
			containsVal = true
			break
		}
	}
	if !containsVal { // Preventing CPU Spiking when input contains only operator and has no operands
		return "", fiber.ErrBadRequest
	}

	for i := 0; i < len(tokens); i++ {
		if tokens[i] == "" || tokens[i] == " " {
			continue
		} else if !isOperator(tokens[i]) {
			_, err := strconv.ParseFloat(tokens[i], 64)
			if err == nil {
				val.Push(tokens[i])
			} else if err != nil {
				return "", err
			}
		} else if tokens[i] == "(" {
			ops.Push("(")
		} else if tokens[i] == ")" {
			for len(ops) != 0 && ops.Top() != "(" {
				b, err := strconv.ParseFloat(val.Pop(), 32)
				if err != nil {
					return "", err
				}

				a, err := strconv.ParseFloat(val.Pop(), 32)
				if err != nil {
					return "", err
				}

				op := ops.Pop()

				result := applyOperator(a, b, op)
				if result == "Divided By Zero" {
					return result, nil
				}

				val.Push(result)
			}

			ops.Pop()

		} else {
			for len(ops) != 0 && precedence(ops.Top()) >= precedence(tokens[i]) {
				b, err := strconv.ParseFloat(val.Pop(), 64)
				if err != nil {
					return "", err
				}

				a, err := strconv.ParseFloat(val.Pop(), 64)
				if err != nil {
					return "", err
				}

				op := ops.Pop()

				result := applyOperator(a, b, op)
				if result == "Divided By Zero" {
					return result, nil
				}

				val.Push(result)
			}

			ops.Push(tokens[i])
		}
	}

	for len(ops) != 0 {
		if val.Top() == "" || val.Top() == " " {
			continue
		}
		b, err := strconv.ParseFloat(val.Pop(), 64)
		if err != nil {
			return "", err
		}

		a, err := strconv.ParseFloat(val.Pop(), 64)
		if err != nil {
			return "", err
		}

		op := ops.Pop()

		result := applyOperator(a, b, op)
		if result == "Divided By Zero" {
			return result, nil
		}

		val.Push(result)
	}

	result := val.Top()

	if len(val) > 1 {
		temp := 1.0
		for i := 0; i < len(val); i++ {
			a, err := strconv.ParseFloat(val[i], 64)
			if err != nil {
				return "", err
			}

			temp = temp * a
		}
		result = fmt.Sprintf("%f", temp)
	}

	if strings.Contains(result, ".") {
		result = strings.TrimRight(result, "0")
		result = strings.TrimRight(result, ".")
	}

	return result, nil
}

func applyOperator(a float64, b float64, op string) string {
	if op == "+" {
		return fmt.Sprintf("%f", a+b)
	}
	if op == "-" {
		return fmt.Sprintf("%f", a-b)
	}
	if op == "x" {
		return fmt.Sprintf("%f", a*b)
	}
	if op == "/" {
		if b == 0 {
			return "Divided By Zero"
		}
		return fmt.Sprintf("%f", a/b)
	}
	if op == "^" {
		return fmt.Sprintf("%f", math.Pow(a, b))
	}

	return "0"
}

func precedence(op string) int {
	if op == "+" || op == "-" {
		return 1
	}
	if op == "x" || op == "/" {
		return 2
	}
	if op == "^" {
		return 3
	}
	return 0
}

func isOperator(token string) bool {

	if token == "+" || token == "-" || token == "x" || token == "/" || token == "^" || token == "(" || token == ")" {
		return true
	}
	return false
}

func split(r rune) bool {
	return r == ' '
}

func getMatchBracket(expression string, index int) int {
	if expression[index:index+1] != "(" {
		return -1
	}
	var st Stack

	for i := index; i < len(expression); i++ {

		if expression[i:i+1] == "(" {
			st.Push(expression[index : index+1])
		} else if expression[i:i+1] == ")" {
			st.Pop()
			if st.IsEmpty() {
				return i
			}
		}
	}
	return -1
}
