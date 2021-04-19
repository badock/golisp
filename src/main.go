package main

import (
	"./lisp"
)

func main() {

	testingLispExpressions := []string{
		"(+ 1 2)",
		"(+ 1 2)(+ 1 2)",
		"T",
		"NIL",
		"(> 1 0)",
		"(if (> 1 0) 1 0)",
		"(setq toto (+ 1 2))",
		"(+ (+ 3 (+ 44 29)) (* 1 2))",
		"(cons 1 (cons 2 (cons 3 4)))",
		"(cons 1 (cons 2 (cons 3 NIL)))",
		"(cons (cons (cons 1 1) 0) (cons 1 (cons 3 (cons 4 5))))",
		"'( 1 2 3 4 5 '( 1 2 3))",
		"(car '(1 2 3 4))",
		"(cdr '(1 2 3 4))",
		"(car '())",
		"(cdr '())",
		"(defun multiply_by_seven (number) \"Multiply NUMBER by seven.\" (* 7 number))",
		"(defun multiply_by_seven (number) (* 7 number)(* 7 number)) (multiply_by_seven 8)",
		"(defun add (a b) (if (> a 0) (add (- a 1) (+ b 1)) b))(add 10 5)",
		"(defun list_length_ (list n) (if list (list_length_ (cdr list) (+ n 1) ) n))(defun list_length (list) (list_length_ list 0))(list_length '(1 2 3 4))",
	}

	for _, testingLispExpression := range testingLispExpressions {
		parseResult := lisp.Parse(testingLispExpression)

		if parseResult.IsSucccessful() {
			successfulParseResult := parseResult.(lisp.SuccessfulParseResult)
			evaluationResult := lisp.Evaluate(successfulParseResult.Expression)
			if evaluationResult.IsSuccessful() {
				successfulEvaluationResult := evaluationResult.(lisp.SuccessfulEvaluationResult)
				strResult := successfulEvaluationResult.Expression.Print()
				print("Result: ")
				println(strResult)
			} else {
				print("Error: ")
				println(parseResult)
			}
		} else {
			print("Compilation error: ")
			println(parseResult.(lisp.UnsuccessfulParseResult).Message)
		}
	}
}
