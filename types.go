package main

import (
	"fmt"
	"math/big"
)

type TermKind string

const (
	KindInt      TermKind = "Int"
	KindStr      TermKind = "Str"
	KindBool     TermKind = "Bool"
	KindBinary   TermKind = "Binary"
	KindCall     TermKind = "Call"
	KindFunction TermKind = "Function"
	KindLet      TermKind = "Let"
	KindIf       TermKind = "If"
	KindPrint    TermKind = "Print"
	KindFirst    TermKind = "First"
	KindSecond   TermKind = "Second"
	KindTuple    TermKind = "Tuple"
	KindVar      TermKind = "Var"
)

type BinaryOp string

const (
	Add BinaryOp = "Add"
	Sub BinaryOp = "Sub"
	Mul BinaryOp = "Mul"
	Div BinaryOp = "Div"
	Rem BinaryOp = "Rem"
	Eq  BinaryOp = "Eq"
	Neq BinaryOp = "Neq"
	Lt  BinaryOp = "Lt"
	Gt  BinaryOp = "Gt"
	Lte BinaryOp = "Lte"
	Gte BinaryOp = "Gte"
	And BinaryOp = "And"
	Or  BinaryOp = "Or"
)

type Term interface{}

type Int struct {
	Kind  TermKind `json:"kind"`
	Value int64    `json:"value"`
}

type Str struct {
	Kind  TermKind `json:"kind"`
	Value string   `json:"value"`
}

type Bool struct {
	Kind  TermKind `json:"kind"`
	Value bool     `json:"value"`
}

type Binary struct {
	Kind TermKind `json:"kind"`
	LHS  Term     `json:"lhs"`
	Op   BinaryOp `json:"op"`
	RHS  Term     `json:"rhs"`
}

type Print struct {
	Kind  TermKind `json:"kind"`
	Value Term     `json:"value"`
}

type First struct {
	Kind  TermKind `json:"kind"`
	Value Term     `json:"value"`
}

type Second struct {
	Kind  TermKind `json:"kind"`
	Value Term     `json:"value"`
}

type If struct {
	Kind      TermKind `json:"kind"`
	Condition Term     `json:"condition"`
	Then      Term     `json:"then"`
	Otherwise Term     `json:"otherwise"`
}

type Tuple struct {
	Kind   TermKind `json:"kind"`
	First  Term     `json:"first"`
	Second Term     `json:"second"`
}

type Parameter struct {
	Text string `json:"text"`
}

type Call struct {
	Kind      TermKind `json:"kind"`
	Callee    Term     `json:"callee"`
	Arguments []Term   `json:"arguments"`
}

type Let struct {
	Kind  TermKind  `json:"kind"`
	Name  Parameter `json:"name"`
	Value Term      `json:"value"`
	Next  Term      `json:"next"`
}

type Var struct {
	Kind TermKind `json:"kind"`
	Text string   `json:"text"`
}

type Function struct {
	Kind       TermKind    `json:"kind"`
	Parameters []Parameter `json:"parameters"`
	Value      Term        `json:"value"`
}

type File struct {
	Name       string `json:"name"`
	Expression Term   `json:"expression"`
}

func convertToInt(lhs interface{}, rhs interface{}) (*big.Int, *big.Int) {
	var lhsInt int64
	var rhsInt int64
	var okLhs bool = false
	var okRhs bool = false

	if _, ok := lhs.(int64); ok {
		lhsInt = lhs.(int64)
		okLhs = true
	}

	if _, ok := rhs.(int64); ok {
		rhsInt = rhs.(int64)
		okRhs = true
	}

	if _, ok := lhs.(*big.Int); ok {
		lhsInt = lhs.(*big.Int).Int64()
		okLhs = true
	}

	if _, ok := rhs.(*big.Int); ok {
		rhsInt = rhs.(*big.Int).Int64()
		okRhs = true
	}

	if !okLhs || !okRhs {
		fmt.Println("Error converting to int")
		return nil, nil
	}

	return big.NewInt(lhsInt), big.NewInt(rhsInt)
}

func convertToBool(lhs interface{}, rhs interface{}) (bool, bool) {
	var okLhs bool = false
	var okRhs bool = false

	if _, ok := lhs.(int64); ok {
		if lhs != 0 {
			okLhs = true
		}
	}

	if _, ok := rhs.(int64); ok {
		if rhs != 0 {
			okRhs = true
		}
	}

	if _, ok := lhs.(string); ok {
		if lhs != "" {
			okLhs = true
		}
	}

	if _, ok := rhs.(string); ok {
		if rhs != "" {
			okRhs = true
		}
	}

	if _, ok := lhs.(bool); ok {
		okLhs = lhs.(bool)
	}

	if _, ok := rhs.(bool); ok {
		okRhs = rhs.(bool)
	}

	if lhs == nil {
		okLhs = false
	}

	if rhs == nil {
		okRhs = false
	}

	return okLhs, okRhs
}
