package main
import (
	"fmt"
	"go/token"
)

func branchTypeToRascal(node token.Token) string {
	switch node {
	case token.GOTO:
		return "gotoBranch()"
	case token.FALLTHROUGH:
		return "fallthroughBranch()"
	case token.BREAK:
		return "breakBranch()"
	case token.CONTINUE:
		return "continueBranch()"
	default:
		panic(fmt.Sprintf("unknownBranchType(\"%s\"), node.String()))
	}
}

func rascalBranchTypeToGo(rascalBranchType string) token.Token {
	switch rascalBranchType{
	case "gotoBranch()":
		return token.GOTO
	case "fallthroughBranch()":
		return token.FALLTHROUGH
	case "breakBranch()":
		return token.BREAK
	case "continueBranch()":
		return token.CONTINUE
	default:
		panic(fmt.Sprintf("unknownBranchType(\"%s\"), rascalBranchType))
	}
}

func declTypeToRascal(node token.Token) string {
	switch node {
	case token.TYPE:
		return "typeDecl()"
	case token.VAR:
		return "varDecl()"
	case token.IMPORT:
		return "importDecl()"
	case token.CONST:
		return "constDecl()"
	default:
		panic(fmt.Sprintf("unknownDeclType(\"%s\"), node.String()))
	}
}

func rascalDeclTypeToGo(rascalDeclType string) token.Token {
	switch rascalDeclType{
	case "typeDecl()":
		return token.TYPE
	case "varDecl()":
		return token.VAR
	case "importDecl()":
		return token.IMPORT
	case "constDecl()":
		return token.CONST
	default:
		panic(fmt.Sprintf("unknownDeclType(\"%s\"), rascalDeclType))
	}
}

func assignmentOpToRascal(node token.Token) string {
	switch node {
	case token.MUL_ASSIGN:
		return "mulAssign()"
	case token.QUO_ASSIGN:
		return "quoAssign()"
	case token.REM_ASSIGN:
		return "remAssign()"
	case token.AND_ASSIGN:
		return "andAssign()"
	case token.OR_ASSIGN:
		return "orAssign()"
	case token.XOR_ASSIGN:
		return "xorAssign()"
	case token.SHL_ASSIGN:
		return "shiftLeftAssign()"
	case token.SHR_ASSIGN:
		return "shiftRightAssign()"
	case token.AND_NOT_ASSIGN:
		return "andNotAssign()"
	case token.DEFINE:
		return "defineAssign()"
	case token.ASSIGN:
		return "assign()"
	case token.ILLEGAL:
		return "noKey()"
	case token.ADD_ASSIGN:
		return "addAssign()"
	case token.SUB_ASSIGN:
		return "subAssign()"
	default:
		panic(fmt.Sprintf("unknownAssignmentOp(\"%s\"), node.String()))
	}
}

func rascalAssignmentOpToGo(rascalAssignmentOp string) token.Token {
	switch rascalAssignmentOp{
	case "mulAssign()":
		return token.MUL_ASSIGN
	case "quoAssign()":
		return token.QUO_ASSIGN
	case "remAssign()":
		return token.REM_ASSIGN
	case "andAssign()":
		return token.AND_ASSIGN
	case "orAssign()":
		return token.OR_ASSIGN
	case "xorAssign()":
		return token.XOR_ASSIGN
	case "shiftLeftAssign()":
		return token.SHL_ASSIGN
	case "shiftRightAssign()":
		return token.SHR_ASSIGN
	case "andNotAssign()":
		return token.AND_NOT_ASSIGN
	case "defineAssign()":
		return token.DEFINE
	case "assign()":
		return token.ASSIGN
	case "noKey()":
		return token.ILLEGAL
	case "addAssign()":
		return token.ADD_ASSIGN
	case "subAssign()":
		return token.SUB_ASSIGN
	default:
		panic(fmt.Sprintf("unknownAssignmentOp(\"%s\"), rascalAssignmentOp))
	}
}

func opToRascal(node token.Token) string {
	switch node {
	case token.GEQ:
		return "greaterThanEq()"
	case token.MUL:
		return "mul()"
	case token.TILDE:
		return "tilde()"
	case token.QUO:
		return "quo()"
	case token.REM:
		return "rem()"
	case token.AND:
		return "and()"
	case token.OR:
		return "or()"
	case token.XOR:
		return "xor()"
	case token.SHL:
		return "shiftLeft()"
	case token.SHR:
		return "shiftRight()"
	case token.AND_NOT:
		return "andNot()"
	case token.LAND:
		return "logicalAnd()"
	case token.LOR:
		return "logicalOr()"
	case token.ARROW:
		return "arrow()"
	case token.INC:
		return "inc()"
	case token.DEC:
		return "dec()"
	case token.EQL:
		return "equal()"
	case token.LSS:
		return "lessThan()"
	case token.GTR:
		return "greaterThan()"
	case token.NOT:
		return "not()"
	case token.NEQ:
		return "notEqual()"
	case token.ADD:
		return "add()"
	case token.LEQ:
		return "lessThanEq()"
	case token.SUB:
		return "sub()"
	default:
		panic(fmt.Sprintf("unknownOp(\"%s\"), node.String()))
	}
}

func rascalOpToGo(rascalOp string) token.Token {
	switch rascalOp{
	case "greaterThanEq()":
		return token.GEQ
	case "mul()":
		return token.MUL
	case "tilde()":
		return token.TILDE
	case "quo()":
		return token.QUO
	case "rem()":
		return token.REM
	case "and()":
		return token.AND
	case "or()":
		return token.OR
	case "xor()":
		return token.XOR
	case "shiftLeft()":
		return token.SHL
	case "shiftRight()":
		return token.SHR
	case "andNot()":
		return token.AND_NOT
	case "logicalAnd()":
		return token.LAND
	case "logicalOr()":
		return token.LOR
	case "arrow()":
		return token.ARROW
	case "inc()":
		return token.INC
	case "dec()":
		return token.DEC
	case "equal()":
		return token.EQL
	case "lessThan()":
		return token.LSS
	case "greaterThan()":
		return token.GTR
	case "not()":
		return token.NOT
	case "notEqual()":
		return token.NEQ
	case "add()":
		return token.ADD
	case "lessThanEq()":
		return token.LEQ
	case "sub()":
		return token.SUB
	default:
		panic(fmt.Sprintf("unknownOp(\"%s\"), rascalOp))
	}
}
