package evaluator

import (
    "fmt"
    "monkey/ast"
    "monkey/object"
)

// no need to create new instances of true and false every time if we can reference them
var (
    NULL = &object.Null{}
    TRUE = &object.Boolean{Value: true}
    FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
    switch node := node.(type) {
        // STATEMENTS
        case *ast.Program:
            return evalProgram(node, env)
        case *ast.ExpressionStatement:
            return Eval(node.Expression, env)
        case *ast.BlockStatement:
            return evalBlockStatement(node, env)
        case *ast.LetStatement:
            val := Eval(node.Value, env)
            if isError(val) {
                return val
            }
            env.Set(node.Name.Value, val)
        case *ast.ReturnStatement:
            val := Eval(node.ReturnValue, env)
            if isError(val) {
                return val
            }
            return &object.ReturnValue{Value: val}

        // EXPRESSIONS
        case *ast.IntegerLiteral:
            return &object.Integer{Value: node.Value}
        case *ast.PrefixExpression:
            right := Eval(node.Right, env)
            if isError(right) {
                return right
            }
            return evalPrefixExpression(node.Operator, right)
        case *ast.InfixExpression:
            left := Eval(node.Left, env)
            if isError(left) {
                return left
            }

            right := Eval(node.Right, env)
            if isError(right) {
                return right
            }

            return evalInfixExpression(node.Operator, left, right)
        case *ast.Boolean:
            return boolToBoolean(node.Value)
        case *ast.IfExpression:
            return evalIfExpression(node, env)
        case *ast.Identifier:
            return evalIdentifier(node, env)
        case *ast.FunctionLiteral:
            return &object.Function{
                // reuse Parameters and Body of the AST node
                Parameters: node.Parameters,
                Body: node.Body,
                Env: env,
            }
        case *ast.CallExpression:
            function := Eval(node.Function, env)
            if isError(function) {
                return function
            }

            args := evalExpressions(node.Arguments, env)
            if len(args) == 1 && isError(args[0]) {
                return args[0]
            }

            return applyFunction(function, args)
        case *ast.IndexExpression:
            left := Eval(node.Left, env)
            if isError(left) {
                return left
            }
            index := Eval(node.Index, env)
            if isError(index) {
                return index
            }
            return evalIndexExpression(left, index)
        case *ast.StringLiteral:
            return &object.String{
                Value: node.Value,
            }
        case *ast.ArrayLiteral:
            elems := evalExpressions(node.Elements, env)
            if len(elems) == 1 && isError(elems[0]) {
                return elems[0]
            }
            return &object.Array{Elements: elems}
        case *ast.HashLiteral:
            return evalHashLiteral(node, env)
    }

    return nil
}

func newError(format string, a ...interface{}) *object.Error {
    return &object.Error{
        Message: fmt.Sprintf(format, a...),
    }
}

func isError(obj object.Object) bool {
    if obj != nil {
        return obj.Type() == object.ERROR_OBJ
    }
    return false
}

func boolToBoolean(val bool) *object.Boolean {
    if val {
        return TRUE
    }
    return FALSE
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
    var result object.Object

    for _, statement := range program.Statements {
        result = Eval(statement, env)

        switch result := result.(type) {
        case *object.ReturnValue:
            return result.Value
        case *object.Error:
            return result
        }
    }

    return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
    switch operator {
        case "!":
            return evalBangOperatorExpression(right)
        case "-":
            return evalMinusPrefixOperatorExpression(right)
        default:
            return newError("unknown operator: %s%s", operator, right.Type())
    }
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
    switch {
    case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
        return evalIntegerInfixExpression(operator, left, right)
    case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
        return evalStringInfixExpression(operator, left, right)
    // since we're handling pointers in the next few statements, check for all other operands
    // will have to happen before them. Otherwise 5 == 5 would evaluate to false since new pointers
    // are initialized whenever an Integer object is created. Comparison between booleans is therefore
    // faster as the interpreter has to unwrap Integer objects before a comparison is made
    case operator == "==":
        return boolToBoolean(left == right)
    case operator == "!=":
        return boolToBoolean(left != right)
    case left.Type() != right.Type():
        return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
    default:
        return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
    }
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
    // supporting only concat for now. Can support == and != later
    if operator != "+" {
        return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
    }

    leftVal := left.(*object.String).Value
    rightVal := right.(*object.String).Value
    return &object.String{Value: leftVal + rightVal}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
    // this is where unwrapping of the value happens
    leftVal := left.(*object.Integer).Value
    rightVal := right.(*object.Integer).Value

    switch operator {
    case "+":
        return &object.Integer{Value: leftVal + rightVal}
    case "-":
        return &object.Integer{Value: leftVal - rightVal}
    case "*":
        return &object.Integer{Value: leftVal * rightVal}
    case "/":
        return &object.Integer{Value: leftVal / rightVal}
    case "<":
        return boolToBoolean(leftVal < rightVal)
    case ">":
        return boolToBoolean(leftVal > rightVal)
    case "==":
        return boolToBoolean(leftVal == rightVal)
    case "!=":
        return boolToBoolean(leftVal != rightVal)
    default:
        return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
    }
}

func evalBangOperatorExpression(right object.Object) object.Object {
    switch right {
        case TRUE:
            return FALSE
        case FALSE:
            return TRUE
        case NULL:
            return TRUE
        default:
            return FALSE
    }
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
    if right.Type() != object.INTEGER_OBJ {
        return newError("unknown operator: -%s", right.Type())
    }

    value := right.(*object.Integer).Value
    return &object.Integer{Value: -value}
}

func isTruthy(obj object.Object) bool {
    switch obj {
    case NULL:
        return false
    case TRUE:
        return true
    case FALSE:
        return false
    default:
        return true
    }
}

// consequence is evaluated when condition is truthy i.e. not null and not false
// can design this to have consequence evaluated when condition is strictly true as well
func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
    condition := Eval(ie.Condition, env)
    if isError(condition) {
        return condition
    }

    if isTruthy(condition) {
        return Eval(ie.Consequence, env)
    } else if (ie.Alternative != nil) {
        return Eval(ie.Alternative, env)
    } else {
        return NULL
    }
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
    var result object.Object

    for _, statement := range block.Statements {
        result = Eval(statement, env)

        if result != nil && (result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ) {
            return result
        }
    }

    return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
    if val, ok := env.Get(node.Value); ok {
        return val
    }

    if builtin, ok := builtins[node.Value]; ok {
        return builtin
    }

    return newError("identifier not found: " + node.Value)
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
    pairs := make(map[object.HashKey]object.HashPair)

    for keyNode, valueNode := range node.Pairs {
        key := Eval(keyNode, env)
        if isError(key) {
            return key
        }

        hashKey, ok := key.(object.Hashable)
        if !ok {
            return newError("unusable as hash key: %s", key.Type())
        }

        value := Eval(valueNode, env)
        if isError(value) {
            return value
        }

        hashed := hashKey.HashKey()
        pairs[hashed] = object.HashPair{Key: key, Value: value}
    }

    return &object.Hash{Pairs: pairs}
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
    var result []object.Object

    for _, exp := range exps {
        evaluated := Eval(exp, env)
        if isError(evaluated) {
            return []object.Object{evaluated}
        }
        result = append(result, evaluated)
    }

    return result
}

func evalIndexExpression(left, index object.Object) object.Object {
    switch {
    case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
        return evalArrayIndexExpression(left, index)
    case left.Type() == object.HASH_OBJ:
        return evalHashIndexExpression(left, index)
    default:
        return newError("index operator not supported: %s", left.Type())
    }
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
    arrayObject := array.(*object.Array)
    idx := index.(*object.Integer).Value
    max := int64(len(arrayObject.Elements) - 1)

    // handles index out of bounds
    if idx < 0 || idx > max {
        return NULL
    }

    return arrayObject.Elements[idx]
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
    hashObject := hash.(*object.Hash)

    key, ok := index.(object.Hashable)
    if !ok {
        return newError("unusable as hash key: %s", index.Type())
    }

    pair, ok := hashObject.Pairs[key.HashKey()]
    if !ok {
        return NULL
    }

    return pair.Value
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
    switch function := fn.(type) {
    case *object.Function:
        extendedEnv := extendFunctionEnv(function, args)
        evaluated := Eval(function.Body, extendedEnv)
        return unwrapReturnValue(evaluated)

    case *object.BuiltIn:
        return function.Fn(args...)

        default:return newError("not a function: %s", fn.Type())
    }
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
    env := object.NewEnclosedEnvironment(fn.Env)

    for paramIdx, param := range fn.Parameters {
        env.Set(param.Value, args[paramIdx])
    }

    return env
}

func unwrapReturnValue(obj object.Object) object.Object {
    if returnValue, ok := obj.(*object.ReturnValue); ok {
        return returnValue.Value
    }

    return obj
}
