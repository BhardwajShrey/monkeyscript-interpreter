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

func Eval(node ast.Node) object.Object {
    switch node := node.(type) {
        // STATEMENTS
        case *ast.Program:
            return evalProgram(node)
        case *ast.ExpressionStatement:
            return Eval(node.Expression)
        case *ast.BlockStatement:
            return evalBlockStatement(node)
        case *ast.ReturnStatement:
            val := Eval(node.ReturnValue)
            if isError(val) {
                return val
            }
            return &object.ReturnValue{Value: val}
        // EXPRESSIONS
        case *ast.IntegerLiteral:
            return &object.Integer{Value: node.Value}
        case *ast.PrefixExpression:
            right := Eval(node.Right)
            if isError(right) {
                return right
            }
            return evalPrefixExpression(node.Operator, right)
        case *ast.InfixExpression:
            left := Eval(node.Left)
            if isError(left) {
                return left
            }

            right := Eval(node.Right)
            if isError(right) {
                return right
            }

            return evalInfixExpression(node.Operator, left, right)
        case *ast.Boolean:
            return boolToBoolean(node.Value)
        case *ast.IfExpression:
            return evalIfExpression(node)
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

func evalProgram(program *ast.Program) object.Object {
    var result object.Object

    for _, statement := range program.Statements {
        result = Eval(statement)

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
func evalIfExpression(ie *ast.IfExpression) object.Object {
    condition := Eval(ie.Condition)
    if isError(condition) {
        return condition
    }

    if isTruthy(condition) {
        return Eval(ie.Consequence)
    } else if (ie.Alternative != nil) {
        return Eval(ie.Alternative)
    } else {
        return NULL
    }
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
    var result object.Object

    for _, statement := range block.Statements {
        result = Eval(statement)

        if result != nil && (result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ) {
            return result
        }
    }

    return result
}
