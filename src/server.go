package main

import (
	"./lisp"
	"fmt"
	"log"
	"net/http"
)

func viewHandler(w http.ResponseWriter, r *http.Request) {

	status := 0
	message := ""
	requestResult := ""

	expression := r.FormValue("expression")

	if expression == "" {
		status = -1
		message = "Please provide an 'expression' POST form value"
	} else {
		parseResult := lisp.Parse(expression)
		if parseResult.IsSucccessful() {
			successfulParseResult := parseResult.(lisp.SuccessfulParseResult)
			evaluationResult := lisp.Evaluate(successfulParseResult.Expression)
			if evaluationResult.IsSuccessful() {
				successfulEvaluationResult := evaluationResult.(lisp.SuccessfulEvaluationResult)
				strResult := successfulEvaluationResult.Expression.Print()
				print("Result: ")
				println(strResult)
				status = 0
				message = "Compilation and evaluation were successful"
				requestResult = strResult
			} else {
				print("Error: ")
				println(parseResult)
				status = 1
				message = "Compilation was successful but evaluation was not successful"
			}
		} else {
			print("Compilation error: ")
			println(parseResult.(lisp.UnsuccessfulParseResult).Message)
			status = 2
			message = "Compilation was not successful"
		}
	}

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	fmt.Fprintf(w, "{\"status:\": %d, \"result\": \"%s\", \"msg\": \"%s\"}", status, requestResult, message)
}

func main() {
	http.HandleFunc("/lispgo", viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
