package main

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type Scope map[string]Term

func eval(scope Scope, termData Term) Term {
	kind := termData.(map[string]interface{})["kind"].(string)

	switch TermKind(kind) {
	case KindInt:
		var intValue Int
		decode(termData, &intValue)

		return big.NewInt(intValue.Value)
	case KindStr:
		var stringValue Str
		decode(termData, &stringValue)

		return stringValue.Value
	case KindBool:
		var boolValue Bool
		decode(termData, &boolValue)

		return boolValue.Value
	case KindPrint:
		var printValue Print
		decode(termData, &printValue)

		fmt.Println(eval(scope, printValue.Value))
	case KindBinary:
		var binaryValue Binary
		decode(termData, &binaryValue)

		lhs := eval(scope, binaryValue.LHS)
		rhs := eval(scope, binaryValue.RHS)
		op := BinaryOp(binaryValue.Op)

		switch op {
		case Add:
			if intLhs, ok := lhs.(*big.Int); ok {
				if intRhs, ok := rhs.(*big.Int); ok {
					return new(big.Int).Add(intRhs, intLhs)
				} else if strRhs, ok := rhs.(string); ok {
					return fmt.Sprintf("%d%s", intLhs, strRhs)
				}
			} else if strLhs, ok := lhs.(string); ok {
				if intRhs, ok := rhs.(*big.Int); ok {
					return fmt.Sprintf("%s%d", strLhs, intRhs)
				} else if strRhs, ok := rhs.(string); ok {
					return fmt.Sprintf("%s%s", strLhs, strRhs)
				}
			}
		case Sub:
			intLhs, intRhs := convertToInt(lhs, rhs)
			return new(big.Int).Sub(intLhs, intRhs)
		case Mul:
			intLhs, intRhs := convertToInt(lhs, rhs)
			return new(big.Int).Mul(intLhs, intRhs)
		case Div:
			intLhs, intRhs := convertToInt(lhs, rhs)
			return new(big.Int).Div(intLhs, intRhs)
		case Rem:
			intLhs, intRhs := convertToInt(lhs, rhs)
			return new(big.Int).Rem(intLhs, intRhs)
		case Eq:
			return fmt.Sprintf("%v", lhs) == fmt.Sprintf("%v", rhs)
		case Neq:
			return fmt.Sprintf("%v", lhs) != fmt.Sprintf("%v", rhs)
		case Lt:
			intLhs, intRhs := convertToInt(lhs, rhs)
			comp := intLhs.Cmp(intRhs)
			return comp < 0
		case Gt:
			intLhs, intRhs := convertToInt(lhs, rhs)
			comp := intLhs.Cmp(intRhs)
			return comp > 0
		case Lte:
			intLhs, intRhs := convertToInt(lhs, rhs)
			comp := intLhs.Cmp(intRhs)
			return comp <= 0
		case Gte:
			intLhs, intRhs := convertToInt(lhs, rhs)
			comp := intLhs.Cmp(intRhs)
			return comp >= 0
		case And:
			boolLhs, boolRhs := convertToBool(lhs, rhs)
			return boolLhs && boolRhs
		case Or:
			boolLhs, boolRhs := convertToBool(lhs, rhs)
			return boolLhs || boolRhs
		}
	case KindIf:
		var ifValue If
		decode(termData, &ifValue)

		condition := eval(scope, ifValue.Condition)
		if bool(condition.(bool)) {
			return eval(scope, ifValue.Then)
		} else {
			return eval(scope, ifValue.Otherwise)
		}
	case KindTuple:
		var tupleValue Tuple
		decode(termData, &tupleValue)

		firstEl := eval(scope, tupleValue.First)
		secondEl := eval(scope, tupleValue.Second)

		return fmt.Sprintf("(%v, %v)", firstEl, secondEl)
	case KindFirst:
		var firstValue First
		var tuple Tuple
		decode(termData, &firstValue)
		decode(firstValue.Value, &tuple)

		if tuple.Kind != KindTuple {
			fmt.Println("Error evaluating first element of tuple")
			return nil
		}

		firstEl := tuple.First
		firstElValue := eval(scope, firstEl)
		return firstElValue
	case KindSecond:
		var secondValue First
		var tuple Tuple
		decode(termData, &secondValue)
		decode(secondValue.Value, &tuple)

		if tuple.Kind != KindTuple {
			fmt.Println("Error evaluating second element of tuple")
			return nil
		}

		secondEl := tuple.Second
		secondElValue := eval(scope, secondEl)
		return secondElValue
	case KindCall:
		var callValue Call
		decode(termData, &callValue)

		function := reflect.ValueOf(eval(scope, callValue.Callee))

		var args []Term
		for _, arg := range callValue.Arguments {
			args = append(args, eval(scope, arg))
		}

		return function.Call([]reflect.Value{reflect.ValueOf(args), reflect.ValueOf(scope)})[0].Interface().(Term)
	case KindLet:
		var letValue Let
		decode(termData, &letValue)

		scope[letValue.Name.Text] = eval(scope, letValue.Value)
		return eval(scope, letValue.Next)
	case KindVar:
		var varValue Var
		decode(termData, &varValue)

		var value Term
		var ok bool

		if value, ok = scope[varValue.Text]; !ok {
			fmt.Println("Error: variable not found")
			return nil
		}

		return value
	case KindFunction:
		var function Function
		decode(termData, &function)

		return func(args []Term, scope Scope) Term {
			if len(args) != len(function.Parameters) {
				fmt.Println("Error: incorrect number of arguments")
				return nil
			}
			iScope := Scope{}
			for key, value := range scope {
				iScope[key] = value
			}
			for index, value := range function.Parameters {
				iScope[value.Text] = args[index]
			}

			return eval(iScope, function.Value)
		}
	}

	return nil
}

func decode(term Term, value Term) Term {
	err := mapstructure.Decode(term, &value)

	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	return value
}
