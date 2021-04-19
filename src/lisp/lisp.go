package lisp

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Expression interface {
	Evaluate(context EvaluationContext) EvaluationResult
	GetType() string
	Print() string
}


type BaseTypeExpression interface {
	 Expression
}

type Block struct {
	Expression
	SubExpressions []Expression
}

func (b Block) GetType() string {
	return "block"
}

func (b Block) Print() string {
	result := "block"
	return result
}

type Int struct {
	BaseTypeExpression
	Value int
}

func (i Int) GetType() string {
	return "int"
}

func (i Int) Print() string {
	return fmt.Sprintf("%d", i.Value)
}

type String struct {
	BaseTypeExpression
	Value string
}

func (s String) GetType() string {
	return "string"
}

func (s String) Print() string {
	return s.Value
}

type Boolean struct {
	BaseTypeExpression
	Value bool
}

func (b Boolean) GetType() string {
	return "boolean"
}

func (b Boolean) Print() string {
	if b.Value {
		return "T"
	} else {
		return "NIL"
	}
}

type List struct {
	BaseTypeExpression
	left Expression
	right Expression
	//Value []Expression
}

func (l List) GetType() string {
	return "list"
}

func (l List) Print() string {
	result := "("
	elementsToString := []string{}
	elementsToString = append(elementsToString, l.left.Print())
	if ! (l.right.GetType() == "boolean" && ! l.right.(Boolean).Value) {
		rightValue := l.right.Print()
		if l.right.GetType() != "list" {
			rightValue = ". " + rightValue
		} else {
			rightValue = rightValue[1:len(rightValue)-1]
		}
		elementsToString = append(elementsToString, rightValue)
	}
	result += strings.Join(elementsToString, " ")
	result += ")"
	return result
}

type FunctionCall struct {
	Expression
	functionName string
	arguments []Expression
}

func (fc FunctionCall) GetType() string {
	return "functionCall"
}

func (fc FunctionCall) Print() string {
	result := fc.functionName + "("
	argumentsAsStrings := []string{}
	for _, argument := range fc.arguments {
		argumentsAsStrings = append(argumentsAsStrings, argument.Print())
	}
	result += strings.Join(argumentsAsStrings, ",")
	result += ")"
	return result
}

type Variable struct {
	Expression
	Name string
}

func (v Variable) GetType() string {
	return "variable"
}

func (v Variable) Print() string {
	return fmt.Sprintf("@%s", v.Name)
}

type FunctionDeclaration struct {
	Expression
	functionName          string
	functionDocumentation string
	arguments             []Variable
	body                  Block
}

func (fd FunctionDeclaration) GetType() string {
	return "functionDeclaration"
}

func (fd FunctionDeclaration)  Print() string {
	return fmt.Sprintf("#%s", fd.functionName)
}

type ParseResult interface {
	IsSucccessful() bool
}

type SuccessfulParseResult struct {
	ParseResult
	Expression Expression
}

type UnsuccessfulParseResult struct {
	ParseResult
	Message string
}

func (r SuccessfulParseResult) IsSucccessful() bool {
	return true
}

func (r UnsuccessfulParseResult) IsSucccessful() bool {
	return false
}

func ParseString(expression string) ParseResult {
	var is_surrounded_by_quotes = false
	if expression[:1] == "\"" && expression[len(expression)-1:] == "\"" {
		is_surrounded_by_quotes = true
	}

	if is_surrounded_by_quotes {
		payload := expression[1:len(expression)-1]

		if is_surrounded_by_quotes {
			return SuccessfulParseResult{
				Expression: String {
					Value: payload,
				},
			}
		}
	}

	errorMsg := fmt.Sprintf("Cannot parse String from \"%s\"", expression)
	return UnsuccessfulParseResult{
		Message: errorMsg,
	}
}

func ParseVariable(expression string) ParseResult {

	match, _ := regexp.MatchString("[_a-zA-Z][_a-zA-Z1-9]*", expression)
	if match {
		return SuccessfulParseResult{
			Expression: Variable {
				Name: expression,
			},
		}
	}

	errorMsg := fmt.Sprintf("Cannot parse Variable from \"%s\"", expression)
	return UnsuccessfulParseResult{
		Message: errorMsg,
	}
}

func ParseInt(expression string) ParseResult {
	int_value, err := strconv.Atoi(expression)

	if err == nil {
		return SuccessfulParseResult{
			Expression: Int{
				Value: int_value,
			},
		}
	}

	errorMsg := fmt.Sprintf("Cannot parse Int from \"%s\"", expression)
	return UnsuccessfulParseResult{
		Message: errorMsg,
	}
}

func ParseBoolean(expression string) ParseResult {

	if expression == "T" {
		return SuccessfulParseResult{
			Expression: Boolean{
				Value: true,
			},
		}
	}

	if expression == "NIL" {
		return SuccessfulParseResult{
			Expression: Boolean{
				Value: false,
			},
		}
	}

	errorMsg := fmt.Sprintf("Cannot parse Int from \"%s\"", expression)
	return UnsuccessfulParseResult{
		Message: errorMsg,
	}
}

func ParseSeveralExpressionsString(expression string) []string {
	arguments := []string{}

	payload := expression //lisp[1:len(lisp)-1]

	// combine parts of a same string
	currentWord := ""

	_ = arguments
	_ = currentWord

	isProcessingASplittedWord := false
	currentParenthesisLevel := 0
	currentSubExpressionStart := 0

	for i := 0; i < len(payload); i++ {
		if payload[i:i+1] == "\"" && i > 0 && payload[i-1:i] == "\\" {
			isProcessingASplittedWord = !isProcessingASplittedWord
		}

		if payload[i:i+1] == "(" {
			if currentParenthesisLevel == 0 {
				if i > 0 && payload[i-1:i] == "'" {
					currentSubExpressionStart = i-1
				} else {
					currentSubExpressionStart = i
				}
			}
			currentParenthesisLevel += 1
		}

		if payload[i:i+1] == ")" {
			currentParenthesisLevel -= 1
			if currentParenthesisLevel == 0 {
				subExpression := payload[currentSubExpressionStart:i+1]
				arguments = append(arguments, subExpression)
				currentWord = ""
				continue
			}
		}

		if currentParenthesisLevel == 0 {
			if !isProcessingASplittedWord && (payload[i:i+1] == " " || payload[i:i+1] == "\n") {
				if len(currentWord) > 0 {
					arguments = append(arguments, currentWord)
					currentWord = ""
				}
				continue
			} else {
				currentWord += payload[i:i+1]
			}
		}
	}

	if currentWord != "" && currentWord != "\n" {
		arguments = append(arguments, currentWord)
	}

	return arguments
}

func countSubExpressions(expression string) int {
	count := 0
	currentParenthesisLevel := 0
	isProcessingASplittedWord := false

	for i := 0; i < len(expression); i++ {
		if expression[i:i+1] == "\"" && i > 0 && expression[i-1:i] == "\\" {
			isProcessingASplittedWord = !isProcessingASplittedWord
		}

		if expression[i:i+1] == "(" {
			currentParenthesisLevel += 1
		}

		if expression[i:i+1] == ")" {
			currentParenthesisLevel -= 1
			if currentParenthesisLevel == 0 {
				count += 1
			}
		}
	}

	return count
}

func ParseBlock(expression string) ParseResult {

	isProcessingASplittedWord := false
	currentParenthesisLevel := 0
	currentSubExpressionStart := 0

	expressions := []Expression{}
	arguments := []string{}

	_ = expressions

	currentWord := ""

	for i := 0; i < len(expression); i++ {
		if expression[i:i+1] == "\"" && i > 0 && expression[i-1:i] == "\\" {
			isProcessingASplittedWord = !isProcessingASplittedWord
		}

		if expression[i:i+1] == "(" {
			if currentParenthesisLevel == 0 {
				if i > 0 && expression[i-1:i] == "'" {
					currentSubExpressionStart = i-1
				} else {
					currentSubExpressionStart = i
				}
			}
			currentParenthesisLevel += 1
		}

		if expression[i:i+1] == ")" {
			currentParenthesisLevel -= 1
			if currentParenthesisLevel == 0 {
				subExpression := expression[currentSubExpressionStart:i+1]
				arguments = append(arguments, subExpression)
				currentWord = ""
				continue
			}
		}

		if currentParenthesisLevel == 0 {
			if !isProcessingASplittedWord && (expression[i:i+1] == " " || expression[i:i+1] == "\n") {
				if len(currentWord) > 0 {
					arguments = append(arguments, currentWord)
					currentWord = ""
				}
				continue
			} else {
				currentWord += expression[i:i+1]
			}
		}
	}

	if currentParenthesisLevel > 0 || isProcessingASplittedWord {
		errorMessage := "Cannot extract any expression from " + expression
		return UnsuccessfulParseResult{
			Message: errorMessage,
		}
	}

	for _, argument := range arguments {
		parseResult := Parse(argument)
		if parseResult.IsSucccessful() {
			expressions = append(expressions, parseResult.(SuccessfulParseResult).Expression)
		} else {
			return parseResult
		}
	}

	return SuccessfulParseResult{
		Expression: Block{
			SubExpressions: expressions,
		},
	}
}

func ParseFunctionCall(expression string) ParseResult {
	if len(expression) < 2 {
		errorMsg := fmt.Sprintf("Cannot parse function call from \"%s\"", expression)
		return UnsuccessfulParseResult{
			Message: errorMsg,
		}
	}

	arguments := ParseSeveralExpressionsString(expression[1: len(expression)-1])

	if len(arguments) == 0 {
		errorMsg := fmt.Sprintf("Cannot parse Function call from \"%s\"", expression)
		return UnsuccessfulParseResult{
			Message: errorMsg,
		}
	}

	functionName := arguments[0]

	functionArgumentsExpressions := []Expression{}

	for _, functionArgument := range arguments[1:] {
		functionArgumentParseResult := Parse(functionArgument)

		if functionArgumentParseResult.IsSucccessful() {
			r := functionArgumentParseResult.(SuccessfulParseResult)
			functionArgumentsExpressions = append(functionArgumentsExpressions, r.Expression)
		} else {
			errorMsg := fmt.Sprintf("Cannot parse \"%s\"", functionArgument)
			return UnsuccessfulParseResult{
				Message: errorMsg,
			}
		}
	}

	return SuccessfulParseResult{
		Expression: FunctionCall {
			functionName: functionName,
			arguments: functionArgumentsExpressions,
		},
	}
}

func ParseList(expression string) ParseResult {
	if len(expression) < 3 {
		errorMsg := fmt.Sprintf("Cannot parse List from \"%s\"", expression)
		return UnsuccessfulParseResult{
			Message: errorMsg,
		}
	}

	arguments := ParseSeveralExpressionsString(expression[2: len(expression)-1])

	if len(arguments) == 0 {
		return SuccessfulParseResult{
			Expression: Boolean{
				Value: false,
			},
		}
	}

	var firstElement List
	var previousElement List
	var currentElement List

	for i := (len(arguments) -1) ; i >= 0; i-- {
		functionArgument := arguments[i]
		functionArgumentParseResult := Parse(functionArgument)
		if functionArgumentParseResult.IsSucccessful() {
			r := functionArgumentParseResult.(SuccessfulParseResult)
			if i == (len(arguments) - 1) {
				currentElement = List{
					left:  r.Expression,
					right: Boolean{Value: false},
				}
			} else {
				currentElement = List{
					left:  r.Expression,
					right: previousElement,
				}
			}
			previousElement = currentElement
		}
	}
	firstElement = currentElement

	return SuccessfulParseResult{
		Expression: firstElement,
	}
}

func Parse(expression string) ParseResult {

	var isSurroundedByParentheses = false
	var isDefiningAList = false

	if len(expression) >= 2 && expression[:1] == "(" && expression[len(expression)-1:] == ")" {
		isSurroundedByParentheses = true
	}

	if len(expression) >= 3 && expression[:2] == "'(" && expression[len(expression)-1:] == ")" {
		isDefiningAList = true
	}

	payload := expression

	if countSubExpressions(expression) > 1 {
		blockParseResult := ParseBlock(payload)
		if blockParseResult.IsSucccessful() {
			return blockParseResult
		}
	}

	if isDefiningAList {
		listParseResult := ParseList(payload)
		if listParseResult.IsSucccessful() {
			return listParseResult
		}

	} else if isSurroundedByParentheses {
		functionCallParseResult := ParseFunctionCall(payload)
		if functionCallParseResult.IsSucccessful() {
			return functionCallParseResult
		}
	} else {

		booleanParseResult := ParseBoolean(payload)
		if booleanParseResult.IsSucccessful() {
			return booleanParseResult
		}

		stringParseResult := ParseString(payload)
		if stringParseResult.IsSucccessful() {
			return stringParseResult
		}

		intParseResult := ParseInt(payload)
		if intParseResult.IsSucccessful() {
			return intParseResult
		}

		variableParseResult := ParseVariable(payload)
		if variableParseResult.IsSucccessful() {
			return variableParseResult
		}
	}

	errorMsg := fmt.Sprintf("Cannot parse \"%s\"", expression)
	return UnsuccessfulParseResult{
		Message: errorMsg,
	}
}

type EvaluationContext struct {
	Parent *EvaluationContext
	variables map[string]Expression
	functions map[string]FunctionDeclaration
}


type EvaluationResult interface {
	IsSuccessful() bool
}

type SuccessfulEvaluationResult struct {
	EvaluationResult
	Expression Expression
}

type UnsuccessfulEvaluationResult struct {
	EvaluationResult
	message string
}

func (er SuccessfulEvaluationResult) IsSuccessful() bool {
	return true
}

func (er UnsuccessfulEvaluationResult) IsSuccessful() bool {
	return false
}

func (block Block) Evaluate(context EvaluationContext) EvaluationResult {

	results := []Expression{}

	for _, expression := range block.SubExpressions {
		evaluationResult := expression.Evaluate(context)

		if evaluationResult.IsSuccessful() {
			results = append(results, evaluationResult.(SuccessfulEvaluationResult).Expression)
		} else {
			return UnsuccessfulEvaluationResult{}
		}
	}

	result := SuccessfulEvaluationResult {
	}
	if len(results) > 0 {
		result = SuccessfulEvaluationResult {
			Expression: results[len(results)-1],
		}
	}

	return result
}

func (re Boolean) Evaluate(context EvaluationContext) EvaluationResult {
	return SuccessfulEvaluationResult {
		Expression: re,
	}
}

func (re Int) Evaluate(context EvaluationContext) EvaluationResult {
	return SuccessfulEvaluationResult {
		Expression: re,
	}
}

func (re String) Evaluate(context EvaluationContext) EvaluationResult {
	return SuccessfulEvaluationResult {
		Expression: re,
	}
}

func (re List) Evaluate(context EvaluationContext) EvaluationResult {
	return SuccessfulEvaluationResult {
		Expression: re,
	}
}

func (v Variable) Evaluate(context EvaluationContext) EvaluationResult {

	if variableValue, ok := context.variables[v.Name]; ok {
		return SuccessfulEvaluationResult {
			Expression: variableValue,
		}
	}

	return UnsuccessfulEvaluationResult {
	}
}

func compare(evaluatedArguments []EvaluationResult, side int) EvaluationResult {

	if len(evaluatedArguments) != 2 {
		return UnsuccessfulEvaluationResult{}
	}

	if ! (evaluatedArguments[0].IsSuccessful() && evaluatedArguments[1].IsSuccessful()) {
		return UnsuccessfulEvaluationResult{}
	}

	arg1 := evaluatedArguments[0].(SuccessfulEvaluationResult).Expression
	arg2 := evaluatedArguments[1].(SuccessfulEvaluationResult).Expression

	if ! (arg1.GetType() == "int" && arg2.GetType() == "int") {
		return UnsuccessfulEvaluationResult{}
	}

	arg1AsInt := arg1.(Int)
	arg2AsInt := arg2.(Int)

	result := false

	if side == -1 {
		result = arg1AsInt.Value > arg2AsInt.Value
	} else if side == 0 {
		result = arg1AsInt.Value == arg2AsInt.Value
	} else if side == 1 {
		result = arg1AsInt.Value < arg2AsInt.Value
	} else if side == 2 {
		result = arg1AsInt.Value != arg2AsInt.Value
	} else {
		return UnsuccessfulEvaluationResult{}
	}


	return SuccessfulEvaluationResult{
		Expression: Boolean{
			Value: result,
		},
	}
}

func plus(evaluatedArguments []EvaluationResult) EvaluationResult {
	if len(evaluatedArguments) != 2 {
		return UnsuccessfulEvaluationResult{}
	}

	if ! (evaluatedArguments[0].IsSuccessful() && evaluatedArguments[1].IsSuccessful()) {
		return UnsuccessfulEvaluationResult{}
	}

	arg1 := evaluatedArguments[0].(SuccessfulEvaluationResult).Expression
	arg2 := evaluatedArguments[1].(SuccessfulEvaluationResult).Expression

	if ! (arg1.GetType() == "int" && arg2.GetType() == "int") {
		return UnsuccessfulEvaluationResult{}
	}

	arg1AsInt := arg1.(Int)
	arg2AsInt := arg2.(Int)

	return SuccessfulEvaluationResult{
		Expression: Int{
			Value: arg1AsInt.Value + arg2AsInt.Value,
		},
	}
}

func minus(evaluatedArguments []EvaluationResult) EvaluationResult {
	if len(evaluatedArguments) != 2 {
		return UnsuccessfulEvaluationResult{}
	}

	if ! (evaluatedArguments[0].IsSuccessful() && evaluatedArguments[1].IsSuccessful()) {
		return UnsuccessfulEvaluationResult{}
	}

	arg1 := evaluatedArguments[0].(SuccessfulEvaluationResult).Expression
	arg2 := evaluatedArguments[1].(SuccessfulEvaluationResult).Expression

	if ! (arg1.GetType() == "int" && arg2.GetType() == "int") {
		return UnsuccessfulEvaluationResult{}
	}

	arg1AsInt := arg1.(Int)
	arg2AsInt := arg2.(Int)

	return SuccessfulEvaluationResult{
		Expression: Int{
			Value: arg1AsInt.Value - arg2AsInt.Value,
		},
	}
}

func mult(evaluatedArguments []EvaluationResult) EvaluationResult {
	if len(evaluatedArguments) != 2 {
		return UnsuccessfulEvaluationResult{}
	}

	if ! (evaluatedArguments[0].IsSuccessful() && evaluatedArguments[1].IsSuccessful()) {
		return UnsuccessfulEvaluationResult{}
	}

	arg1 := evaluatedArguments[0].(SuccessfulEvaluationResult).Expression
	arg2 := evaluatedArguments[1].(SuccessfulEvaluationResult).Expression

	if ! (arg1.GetType() == "int" && arg2.GetType() == "int") {
		return UnsuccessfulEvaluationResult{}
	}

	arg1AsInt := arg1.(Int)
	arg2AsInt := arg2.(Int)

	return SuccessfulEvaluationResult{
		Expression: Int{
			Value: arg1AsInt.Value * arg2AsInt.Value,
		},
	}
}

func divide(evaluatedArguments []EvaluationResult) EvaluationResult {
	if len(evaluatedArguments) != 2 {
		return UnsuccessfulEvaluationResult{}
	}

	if ! (evaluatedArguments[0].IsSuccessful() && evaluatedArguments[1].IsSuccessful()) {
		return UnsuccessfulEvaluationResult{}
	}

	arg1 := evaluatedArguments[0].(SuccessfulEvaluationResult).Expression
	arg2 := evaluatedArguments[1].(SuccessfulEvaluationResult).Expression

	if ! (arg1.GetType() == "int" && arg2.GetType() == "int") {
		return UnsuccessfulEvaluationResult{}
	}

	arg1AsInt := arg1.(Int)
	arg2AsInt := arg2.(Int)

	return SuccessfulEvaluationResult{
		Expression: Int{
			Value: arg1AsInt.Value / arg2AsInt.Value,
		},
	}
}

func ifFunction(re FunctionCall, context EvaluationContext) EvaluationResult {

	if len(re.arguments) != 3 {
		return UnsuccessfulEvaluationResult{}
	}

	arg1EvaluationResult := re.arguments[0].Evaluate(context)

	if ! (arg1EvaluationResult.IsSuccessful()) {
		return UnsuccessfulEvaluationResult{}
	}

	arg1 := arg1EvaluationResult.(SuccessfulEvaluationResult).Expression

	//if ! (arg1.GetType() == "boolean") {
	//	return UnsuccessfulEvaluationResult{}
	//}

	//arg1AsBoolean := arg1.(Boolean)

	var expressionToExecute Expression

	if ! (arg1.GetType() == "boolean" && ! arg1.(Boolean).Value) {
		expressionToExecute = re.arguments[1]
	} else {
		expressionToExecute = re.arguments[2]
	}

	evaluationResult := expressionToExecute.Evaluate(context)

	if ! evaluationResult.IsSuccessful() {
		return UnsuccessfulEvaluationResult{}
	}

	resultExpression := evaluationResult.(SuccessfulEvaluationResult).Expression

	return SuccessfulEvaluationResult{
		Expression: resultExpression,
	}
}

func makeList(evaluatedArguments []EvaluationResult) EvaluationResult {

	if len(evaluatedArguments) != 2 {
		return UnsuccessfulEvaluationResult{}
	}

	if ! (evaluatedArguments[0].IsSuccessful() && evaluatedArguments[1].IsSuccessful()) {
		return UnsuccessfulEvaluationResult{}
	}

	arg1 := evaluatedArguments[0].(SuccessfulEvaluationResult).Expression
	arg2 := evaluatedArguments[1].(SuccessfulEvaluationResult).Expression

	return SuccessfulEvaluationResult{
		Expression: List{
			left: arg1,
			right: arg2,
		},
	}
}

func defineFunction(expressions []Expression, context EvaluationContext) EvaluationResult {

	functionName := "lambda"
	arguments := []Variable{}
	functionDocumentation := ""

	bodyExpressions := []Expression{}

	for idx, expression := range expressions {
		variableType := expression.GetType()

		if idx == 0 && variableType == "variable" {
			functionName = expression.(Variable).Name
			continue
		}

		if idx == 1 && variableType == "functionCall" {
			expressionAsFunctionCall := expression.(FunctionCall)
			arguments = append(arguments, Variable{Name: expressionAsFunctionCall.functionName})

			for _, functionArgument := range expressionAsFunctionCall.arguments {
				if functionArgument.GetType() == "variable" {
					arguments = append(arguments, functionArgument.(Variable))
				} else {
					return UnsuccessfulEvaluationResult{}
				}
			}
			continue
		}

		if idx == 2 && variableType == "string" {
			functionDocumentation = expression.(String).Value
			continue
		}

		if idx >= 2 {
			bodyExpressions = append(bodyExpressions, expression)
			continue
		}
	}

	functionDeclaration := FunctionDeclaration{
		functionName:          functionName,
		functionDocumentation: functionDocumentation,
		arguments: arguments,
		body: Block{SubExpressions: bodyExpressions},
	}

	context.functions[functionName] = functionDeclaration

	return SuccessfulEvaluationResult{
		Expression: functionDeclaration,
	}
}

func makeAssignment(re FunctionCall, context EvaluationContext) EvaluationResult {
	if len(re.arguments) != 2 {
		return UnsuccessfulEvaluationResult{}
	}

	if re.arguments[0].GetType() != "variable" {
		return UnsuccessfulEvaluationResult{}
	}

	variableName := re.arguments[0].(Variable).Name
	variableValue := re.arguments[1]

	context.variables[variableName] = variableValue

	return SuccessfulEvaluationResult{
		Expression: variableValue,
	}
}

func headList(evaluatedArguments []EvaluationResult) EvaluationResult {
	if len(evaluatedArguments) != 1 {
		return UnsuccessfulEvaluationResult{}
	}

	if ! evaluatedArguments[0].IsSuccessful() {
		return UnsuccessfulEvaluationResult{}
	}

	arg1 := evaluatedArguments[0].(SuccessfulEvaluationResult).Expression

	if arg1.GetType() == "boolean" {
		if !arg1.(Boolean).Value {
			return SuccessfulEvaluationResult{
				Expression: arg1,
			}
		} else {
			return UnsuccessfulEvaluationResult{}
		}
	}

	if arg1.GetType() != "list" {
		return UnsuccessfulEvaluationResult{}
	}

	result := arg1.(List).left

	return SuccessfulEvaluationResult{
		Expression: result,
	}
}

func restList(evaluatedArguments []EvaluationResult) EvaluationResult {
	if len(evaluatedArguments) != 1 {
		return UnsuccessfulEvaluationResult{}
	}

	if ! evaluatedArguments[0].IsSuccessful() {
		return UnsuccessfulEvaluationResult{}
	}

	arg1 := evaluatedArguments[0].(SuccessfulEvaluationResult).Expression

	if arg1.GetType() == "boolean" {
		if !arg1.(Boolean).Value {
			return SuccessfulEvaluationResult{
				Expression: arg1,
			}
		} else {
			return UnsuccessfulEvaluationResult{}
		}
	}

	if arg1.GetType() != "list" {
		return UnsuccessfulEvaluationResult{}
	}

	result := arg1.(List).right

	return SuccessfulEvaluationResult{
		Expression: result,
	}
}


func (re FunctionCall) Evaluate(context EvaluationContext) EvaluationResult {

	newContextvariables := make(map[string]Expression)
	newContextFunctions := make(map[string]FunctionDeclaration)

	// Copy variables from pred
	for key, value := range context.variables {
		newContextvariables[key] = value
	}

	// Copy from the original map to the target map
	for key, value := range context.functions {
		newContextFunctions[key] = value
	}

	functionContext := EvaluationContext{
		Parent: &context,
		variables: newContextvariables,
		functions: newContextFunctions,
	}

	evaluatedArguments := []EvaluationResult{}

	for i := 0; i < len(re.arguments); i++ {
		//variableName := re.arguments[i].(Variable).Name
		//_= variableName
		//variableValue := re.arguments[i].Evaluate(functionContext)
		//if variableValue.IsSuccessful() {
		//	functionContext.variables[variableName] = variableValue.(SuccessfulEvaluationResult).Expression
		//} else {
		//	return UnsuccessfulEvaluationResult{}
		//}
	}

	if re.functionName == "defun" {
		return defineFunction(re.arguments, context)
	}

	if re.functionName == "setq" {
		return makeAssignment(re, functionContext)
	}

	if re.functionName == "if" {
		return ifFunction(re, functionContext)
	}

	for _, functionCallArgument := range re.arguments {
		evaluatedArguments = append(evaluatedArguments, functionCallArgument.Evaluate(context))
	}

	if re.functionName == "car" {
		return headList(evaluatedArguments)
	}

	if re.functionName == "cdr" {
		return restList(evaluatedArguments)
	}

	if re.functionName == ">" {
		return compare(evaluatedArguments, -1)
	}

	if re.functionName == "<" {
		return compare(evaluatedArguments, 0)
	}

	if re.functionName == "=" {
		return compare(evaluatedArguments, 1)
	}

	if re.functionName == "/=" {
		return compare(evaluatedArguments, 2)
	}

	if re.functionName == "+" {
		return plus(evaluatedArguments)
	}

	if re.functionName == "-" {
		return minus(evaluatedArguments)
	}

	if re.functionName == "/" {
		return divide(evaluatedArguments)
	}

	if re.functionName == "*" {
		return mult(evaluatedArguments)
	}

	if re.functionName == "cons" {
		return makeList(evaluatedArguments)
	}

	if functionDeclaration, ok := context.functions[re.functionName]; ok {
		if ! (len(re.arguments) == len(functionDeclaration.arguments)) {
			return  UnsuccessfulEvaluationResult{}
		}

		for i := 0; i < len(re.arguments); i++ {
			variableName := functionDeclaration.arguments[i].Name
			variableValue := re.arguments[i].Evaluate(functionContext)
			if variableValue.IsSuccessful() {
				functionContext.variables[variableName] = variableValue.(SuccessfulEvaluationResult).Expression
			} else {
				return UnsuccessfulEvaluationResult{}
			}
		}

		result := functionDeclaration.body.Evaluate(functionContext)
		return result
	}

	return UnsuccessfulEvaluationResult{}
}

func Evaluate(expression Expression) EvaluationResult {
	context := EvaluationContext{
		variables: make(map[string]Expression),
		functions: make(map[string]FunctionDeclaration),
	}
	return expression.Evaluate(context)
}
