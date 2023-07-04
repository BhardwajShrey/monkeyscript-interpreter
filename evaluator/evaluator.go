package evaluator

import (
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
            return evalStatements(node.Statements)
        case *ast.ExpressionStatement:
            return Eval(node.Expression)
        // EXPRESSIONS
        case *ast.IntegerLiteral:
            return &object.Integer{Value: node.Value}
        case *ast.Boolean:
            return boolToBoolean(node.Value)
    }

    return nil
}

func boolToBoolean(val bool) *object.Boolean {
    if val {
        return TRUE
    }
    return FALSE
}

func evalStatements(stmts [] ast.Statement) object.Object {
    var result object.Object

    for _, statement := range stmts {
        result = Eval(statement)
    }

    return result
}
