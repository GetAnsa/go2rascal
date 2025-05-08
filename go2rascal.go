package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

var rascalizer *strings.Replacer = strings.NewReplacer("<", "\\<", ">", "\\>", "\n", "\\n", "\t", "\\t", "\r", "\\r", "\\", "\\\\", "\"", "\\\"", "'", "\\'")
var filePath string = ""

// Map applies a function to each element of a slice and returns a new slice with the results.
func Map[T, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

// parses file and makes sure there is no error in file input.
func processFile(addLocs bool) string {
	fset := token.NewFileSet()
	if file, err := parser.ParseFile(fset, filePath, nil, 0); err != nil {
		return fmt.Sprintf("errorFile(Could not process file %s, %s)", filePath, err.Error())
	} else {
		return visitFile(file, fset, addLocs)
	}
}

// RascalNode denotes starting and end position of ident(type) portion and
// bracketed((|[|{) expression, barrier value of -1 if not applicable to node
type RascalNode struct {
	identStartPos   int
	identEndPos     int
	bracketStartPos int
	bracketEndPos   int
	children        []*RascalNode
}

func (r RascalNode) inIdent() bool {
	return r.identStartPos != -1 && r.identEndPos == -1
}

func (r RascalNode) inBracket() bool {
	return r.bracketStartPos != -1 && r.bracketEndPos == -1
}

func (r RascalNode) noIdent() bool {
	return r.identStartPos == -1
}

func (r RascalNode) blank() bool {
	return r.identStartPos == -1 && r.identEndPos == -1 && r.bracketStartPos == -1 && r.bracketEndPos == -1
}

// todo proper unicode support
func isWordChar(ch uint8) bool {
	return unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_'
}

func isBracket(ch uint8) bool {
	_, ok := bracketMap[ch]
	return ok
}

// todo go:generate reverse with sed 's/.*uint8{//' | sed 's/\}$//' | tr ',' '\n' | sed 's/ *\(.*\): *\(.*\)/\2: \1/' | sed "s/: '/:\'/g" | tr '\n' ',' | sed 's/\(.*\)/\tvar reverseBracketMap = map[uint8]uint8{\1}\n/' | sed 's/,,}$/}/'
var bracketMap = map[uint8]uint8{'(': ')', '{': '}', '[': ']', '<': '>', ':': ':', '"': '"', '\'': '\''}

var reverseBracketMap = map[uint8]uint8{')': '(', '}': '{', ']': '[', '>': '<', ':': ':', '"': '"', '\'': '\''}

type bracketCounter struct {
	counter    map[uint8]uint
	totalCount uint
}

func newBracketCounter() bracketCounter {
	return bracketCounter{
		counter: map[uint8]uint{},
	}
}

func (b bracketCounter) count(ch uint8) uint {
	if _, ok := bracketMap[ch]; ok {
		b.counter[ch]++
		b.totalCount++
	} else if _, ok := reverseBracketMap[ch]; ok {
		b.counter[ch]--
		b.totalCount--
	}

	return b.totalCount
}

var nodeMap map[int]*RascalNode = make(map[int]*RascalNode)

func visitRascalAst(ast string, headPtr uint, tailPtr uint) {
	var currentNode = RascalNode{
		identStartPos:   -1,
		identEndPos:     -1,
		bracketStartPos: -1,
		bracketEndPos:   -1,
	}

	// so the headPtr should either be at the start of an ident or a beginning of a bracket
	// though after the ident there may be other stuff like = in the case of addLocs
	for isWordChar(ast[headPtr]) {
		// not necessary yet but I have a hunch
		if currentNode.noIdent() {
			currentNode.identStartPos = int(headPtr)
		}
		headPtr++
		if !isWordChar(ast[headPtr]) {
			currentNode.identEndPos = int(headPtr - 1)
		}
	}

	if isBracket(ast[headPtr]) {
		// same hunch
		if !currentNode.inBracket() {
			currentNode.bracketStartPos = int(headPtr)
		}

		// this should be the most common (and only??) case where the tailPtr
		// passed in is already at the closing bracket
		if ast[tailPtr] == bracketMap[ast[headPtr]] {
			currentNode.bracketEndPos = int(tailPtr)
		} else {
			// probably makes sense to get the commas in this loop and just
			// opportunistically double check the closing bracket
			tempSeekPtr := headPtr + 1
			pda := newBracketCounter()
			for tempSeekPtr < tailPtr {
				if isBracket(ast[tempSeekPtr]) {
					if pda.count(ast[tempSeekPtr]) == 0 && ast[tempSeekPtr] == reverseBracketMap[ast[headPtr]] {
						currentNode.bracketEndPos = int(tempSeekPtr)
					}
				}

				tempSeekPtr++
			}
			if tempSeekPtr == tailPtr {
				panic("mismatched bracket char")
			}
		}

	}

}

//func visitRascalAst(ast string, headPtr uint, tailPtr uint) {
//
//	//var headPtr, tailPtr = 0, len(ast) - 1
//
//	var rootNode = RascalNode{
//		identStartPos:   -1,
//		identEndPos:     -1,
//		bracketStartPos: -1,
//		bracketEndPos:   -1,
//	}
//
//	var currentNode = rootNode
//
//	// todo forgot about commas, 2 pointer solution is busted
//	// just keep count of open brackets and iterate forward from headPtr
//	for headPtr < tailPtr {
//		if (ast[headPtr] >= 'a' && ast[headPtr] <= 'z') ||
//			(ast[headPtr] >= 'A' && ast[headPtr] <= 'Z') ||
//			ast[headPtr] >= '0' && ast[headPtr] <= '9' {
//			if currentNode.blank() {
//				currentNode.identStartPos = int(headPtr)
//			} else if currentNode.inIdent() {
//
//			}
//		} else if _, ok := bracketMap[ast[headPtr]]; ok {
//			if currentNode.inIdent() {
//				currentNode.identEndPos = int(headPtr - 1)
//			}
//			if !currentNode.inBracket() {
//				currentNode.bracketStartPos = int(headPtr)
//			}
//			closingPtr := headPtr + 1
//			openCount := 0
//			for closingPtr < tailPtr {
//				if openCount == 0 && ast[closingPtr] == bracketMap[ast[headPtr]] {
//					break
//				} else if ast[closingPtr] == ast[headPtr] {
//					openCount += 1
//				} else if ast[closingPtr] == bracketMap[ast[headPtr]] {
//					openCount -= 1
//				}
//				closingPtr++
//			}
//			if closingPtr == tailPtr {
//				panic(fmt.Sprintf("Mismatched %cs at %d", ast[headPtr], headPtr))
//			}
//
//			if ast[closingPtr] == bracketMap[ast[headPtr]] {
//				currentNode.bracketEndPos = int(closingPtr)
//			} else {
//				panic(fmt.Sprintf("Mismatched %cs at %d", ast[headPtr], headPtr))
//			}
//			// todo inside the bracket will presumably be some stuff?
//			// conceivably several stuffs because commas
//			// also need to handle the locations case
//		}
//	}
//}

func visitFile(node *ast.File, fset *token.FileSet, addLocs bool) string {
	decls := Map(node.Decls, func(t ast.Decl) string {
		return visitDeclaration(&t, fset, addLocs)
	})
	declString := strings.Join(decls, ",")

	packageName := node.Name.Name

	if addLocs {
		locationString := computeLocation(fset, node.FileStart, node.FileEnd)
		return fmt.Sprintf("file(\"%s\", [%s], at=%s)", packageName, declString, locationString)
	} else {
		return fmt.Sprintf("file(\"%s\", [%s])", packageName, declString)
	}

}

func visitDeclaration(node *ast.Decl, fset *token.FileSet, addLocs bool) string {
	switch d := (*node).(type) {
	case *ast.GenDecl:
		return visitGeneralDeclaration(d, fset, addLocs)
	case *ast.FuncDecl:
		return visitFunctionDeclaration(d, fset, addLocs)
	default:
		return "unknownDeclaration()" // This is an error, should panic here
	}
}

func optionalNameToRascal(node *ast.Ident) string {
	if node != nil {
		return fmt.Sprintf("someName(\"%s\")", node.Name)
	} else {
		return "noName()"
	}
}

func visitSpec(node *ast.Spec, fset *token.FileSet, addLocs bool) string {
	switch d := (*node).(type) {
	case *ast.ImportSpec:
		return visitImportSpec(d, fset, addLocs)
	case *ast.ValueSpec:
		return visitValueSpec(d, fset, addLocs)
	case *ast.TypeSpec:
		return visitTypeSpec(d, fset, addLocs)
	default:
		return "unknownSpec()" // This is an error, should panic here
	}
}

func visitImportSpec(node *ast.ImportSpec, fset *token.FileSet, addLocs bool) string {
	specName := optionalNameToRascal(node.Name)
	specPath := literalToRascal(node.Path, fset, addLocs)

	return outputRascalString(node, fset, "importSpec", []string{specName, specPath}, addLocs)
}

func outputRascalString(node ast.Node, fset *token.FileSet, typeName string, values []string, addLocs bool) string {
	builder := strings.Builder{}
	builder.WriteString(typeName)
	builder.WriteString("(")

	// prime the pump to avoid trailing ,
	if len(values) > 0 {
		builder.WriteString(values[0])
	}

	for _, value := range values[1:] {
		builder.WriteString(",")
		builder.WriteString(value)
	}
	if addLocs {
		builder.WriteString(",at=")
		builder.WriteString(computeLocation(fset, node.Pos(), node.End()))
	}

	builder.WriteString(")")

	return builder.String()
}

func visitValueSpec(node *ast.ValueSpec, fset *token.FileSet, addLocs bool) string {
	names := Map(node.Names, func(t *ast.Ident) string {
		return fmt.Sprintf("\"%s\"", t.Name)
	})
	namesStr := fmt.Sprintf("[%s]", strings.Join(names, ","))
	typeStr := visitOptionExpr(&node.Type, fset, addLocs)
	values := visitExprList(node.Values, fset, addLocs)

	return outputRascalString(node, fset, "valueSpec", []string{namesStr, typeStr, values}, addLocs)

}

func visitTypeSpec(node *ast.TypeSpec, fset *token.FileSet, addLocs bool) string {
	typeParams := visitFieldList(node.TypeParams, fset, addLocs)
	typeStr := visitExpr(&node.Type, fset, addLocs)

	return outputRascalString(
		node,
		fset,
		"typeSpec",
		[]string{"\"" + node.Name.Name + "\"", typeParams, typeStr},
		addLocs,
	)
}

func visitSpecList(nodes []ast.Spec, fset *token.FileSet, addLocs bool) string {
	specs := Map(nodes, func(t ast.Spec) string {
		return visitSpec(&t, fset, addLocs)
	})
	return fmt.Sprintf("[%s]", strings.Join(specs, ","))
}

func visitGeneralDeclaration(node *ast.GenDecl, fset *token.FileSet, addLocs bool) string {
	declType := declTypeToRascal(node.Tok)
	specList := visitSpecList(node.Specs, fset, addLocs)

	return outputRascalString(node, fset, "genDecl", []string{declType, specList}, addLocs)
}

func visitFunctionDeclaration(node *ast.FuncDecl, fset *token.FileSet, addLocs bool) string {
	receivers := visitFieldList(node.Recv, fset, addLocs)
	signature := visitFuncType(node.Type, fset, addLocs)
	body := visitOptionalBlockStmt(node.Body, fset, addLocs)

	// todo was printing the full node.Name ident before, but it looks like the syntax in go-analysis doesn't use it?
	return outputRascalString(
		node,
		fset,
		"funDecl",
		[]string{"\"" + node.Name.Name + "\"", receivers, signature, body},
		addLocs,
	)
}

func visitStmt(node *ast.Stmt, fset *token.FileSet, addLocs bool) string {
	switch t := (*node).(type) {
	case *ast.DeclStmt:
		return visitDeclStmt(t, fset, addLocs)
	case *ast.EmptyStmt:
		return visitEmptyStmt(t, fset, addLocs)
	case *ast.LabeledStmt:
		return visitLabeledStmt(t, fset, addLocs)
	case *ast.ExprStmt:
		return visitExprStmt(t, fset, addLocs)
	case *ast.SendStmt:
		return visitSendStmt(t, fset, addLocs)
	case *ast.IncDecStmt:
		return visitIncDecStmt(t, fset, addLocs)
	case *ast.AssignStmt:
		return visitAssignStmt(t, fset, addLocs)
	case *ast.GoStmt:
		return visitGoStmt(t, fset, addLocs)
	case *ast.DeferStmt:
		return visitDeferStmt(t, fset, addLocs)
	case *ast.ReturnStmt:
		return visitReturnStmt(t, fset, addLocs)
	case *ast.BranchStmt:
		return visitBranchStmt(t, fset, addLocs)
	case *ast.BlockStmt:
		return visitBlockStmt(t, fset, addLocs)
	case *ast.IfStmt:
		return visitIfStmt(t, fset, addLocs)
	case *ast.CaseClause:
		return visitCaseClause(t, fset, addLocs)
	case *ast.SwitchStmt:
		return visitSwitchStmt(t, fset, addLocs)
	case *ast.TypeSwitchStmt:
		return visitTypeSwitchStmt(t, fset, addLocs)
	case *ast.CommClause:
		return visitCommClause(t, fset, addLocs)
	case *ast.SelectStmt:
		return visitSelectStmt(t, fset, addLocs)
	case *ast.ForStmt:
		return visitForStmt(t, fset, addLocs)
	case *ast.RangeStmt:
		return visitRangeStmt(t, fset, addLocs)
	default:
		return fmt.Sprintf("unknownStmt(\"%s\")", reflect.TypeOf(node).Name())
	}
}

// todo option wrapper?
func visitOptionStmt(node *ast.Stmt, fset *token.FileSet, addLocs bool) string {
	if *node == nil {
		return "noStmt()"
	} else {
		return fmt.Sprintf("someStmt(%s)", visitStmt(node, fset, addLocs))
	}
}

func visitOptionExpr(node *ast.Expr, fset *token.FileSet, addLocs bool) string {
	if *node == nil {
		return "noExpr()"
	} else {
		return fmt.Sprintf("someExpr(%s)", visitExpr(node, fset, addLocs))
	}
}

//func rascalOptionExprToGo(rascalExpr string) *ast.Expr {
//	if rascalExpr == "noExpr()" {
//		return nil
//	} else {
//		goExpr := rascalExprToGo(
//			strings.TrimSuffix(
//				strings.TrimPrefix(rascalExpr, "someExpr("),
//				")",
//			),
//		)
//		return goExpr
//	}
//}

func visitDeclStmt(node *ast.DeclStmt, fset *token.FileSet, addLocs bool) string {
	declStr := visitDeclaration(&node.Decl, fset, addLocs)

	return outputRascalString(
		node,
		fset,
		"declStmt",
		[]string{declStr},
		addLocs,
	)
}

func visitEmptyStmt(node *ast.EmptyStmt, fset *token.FileSet, addLocs bool) string {
	return outputRascalString(node, fset, "emptyStmt", []string{}, addLocs)
}

func visitLabeledStmt(node *ast.LabeledStmt, fset *token.FileSet, addLocs bool) string {
	labelStr := labelToRascal(node.Label)
	stmtStr := visitStmt(&node.Stmt, fset, addLocs)

	return outputRascalString(node, fset, "labeledStmt", []string{labelStr, stmtStr}, addLocs)
}

func visitExprStmt(node *ast.ExprStmt, fset *token.FileSet, addLocs bool) string {
	exprStr := visitExpr(&node.X, fset, addLocs)

	return outputRascalString(node, fset, "exprStmt", []string{exprStr}, addLocs)
}

func visitSendStmt(node *ast.SendStmt, fset *token.FileSet, addLocs bool) string {
	chanStr := visitExpr(&node.Chan, fset, addLocs)
	valStr := visitExpr(&node.Value, fset, addLocs)

	return outputRascalString(node, fset, "sendStmt", []string{chanStr, valStr}, addLocs)
}

func visitIncDecStmt(node *ast.IncDecStmt, fset *token.FileSet, addLocs bool) string {
	opStr := opToRascal(node.Tok)
	exprStr := visitExpr(&node.X, fset, addLocs)

	return outputRascalString(node, fset, "incDecStmt", []string{opStr, exprStr}, addLocs)
}

func visitAssignStmt(node *ast.AssignStmt, fset *token.FileSet, addLocs bool) string {
	assignOp := assignmentOpToRascal(node.Tok)
	left := visitExprList(node.Lhs, fset, addLocs)
	right := visitExprList(node.Rhs, fset, addLocs)

	return outputRascalString(node, fset, "assignStmt", []string{left, right, assignOp}, addLocs)
}

func visitGoStmt(node *ast.GoStmt, fset *token.FileSet, addLocs bool) string {
	exprStr := visitCallExpr(node.Call, fset, addLocs)

	return outputRascalString(node, fset, "goStmt", []string{exprStr}, addLocs)
}

func visitDeferStmt(node *ast.DeferStmt, fset *token.FileSet, addLocs bool) string {
	exprStr := visitCallExpr(node.Call, fset, addLocs)

	return outputRascalString(node, fset, "deferStmt", []string{exprStr}, addLocs)
}

func visitReturnStmt(node *ast.ReturnStmt, fset *token.FileSet, addLocs bool) string {
	results := visitExprList(node.Results, fset, addLocs)

	return outputRascalString(node, fset, "returnStmt", []string{results}, addLocs)
}

func visitBranchStmt(node *ast.BranchStmt, fset *token.FileSet, addLocs bool) string {
	typeStr := branchTypeToRascal(node.Tok)
	labelStr := labelToRascal(node.Label)

	return outputRascalString(node, fset, "branchStmt", []string{typeStr, labelStr}, addLocs)
}

// todo option wrapper sounding better by the minute
func visitOptionalBlockStmt(node *ast.BlockStmt, fset *token.FileSet, addLocs bool) string {
	if node == nil {
		return "noStmt()"
	} else {
		return fmt.Sprintf("someStmt(%s)", visitBlockStmt(node, fset, addLocs))
	}
}

func visitBlockStmt(node *ast.BlockStmt, fset *token.FileSet, addLocs bool) string {
	stmts := visitStmtList(node.List, fset, addLocs)

	return outputRascalString(node, fset, "blockStmt", []string{stmts}, addLocs)
}

func visitIfStmt(node *ast.IfStmt, fset *token.FileSet, addLocs bool) string {
	initStmt := visitOptionStmt(&node.Init, fset, addLocs)
	condExpr := visitExpr(&node.Cond, fset, addLocs)
	body := visitBlockStmt(node.Body, fset, addLocs)
	elseStmt := visitOptionStmt(&node.Else, fset, addLocs)

	return outputRascalString(node, fset, "ifStmt", []string{initStmt, condExpr, body, elseStmt}, addLocs)
}

// todo maybe optionwrapper could be generalized to include this?
func caseToRascal(nodes []ast.Expr, fset *token.FileSet, addLocs bool) string {
	if nodes != nil {
		return fmt.Sprintf("regularCase(%s)", visitExprList(nodes, fset, addLocs))
	} else {
		return "defaultCase()"
	}
}

func visitCaseClause(node *ast.CaseClause, fset *token.FileSet, addLocs bool) string {
	caseStr := caseToRascal(node.List, fset, addLocs)
	stmtsString := visitStmtList(node.Body, fset, addLocs)

	return outputRascalString(node, fset, "caseClause", []string{caseStr, stmtsString}, addLocs)
}

func visitCaseClauseList(node *ast.BlockStmt, fset *token.FileSet, addLocs bool) string {
	var cases []string
	for i := 0; node != nil && node.List != nil && i < len(node.List); i++ {
		nodeAsCase, ok := node.List[i].(*ast.CaseClause)
		if ok {
			cases = append(cases, visitCaseClause(nodeAsCase, fset, addLocs))
		} else {
			cases = append(cases, "invalidCaseClause()")
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(cases, ","))
}

func visitSwitchStmt(node *ast.SwitchStmt, fset *token.FileSet, addLocs bool) string {
	init := visitOptionStmt(&node.Init, fset, addLocs)
	tag := visitOptionExpr(&node.Tag, fset, addLocs)
	block := visitCaseClauseList(node.Body, fset, addLocs)

	return outputRascalString(node, fset, "switchStmt", []string{init, tag, block}, addLocs)
}

func visitTypeSwitchStmt(node *ast.TypeSwitchStmt, fset *token.FileSet, addLocs bool) string {
	init := visitOptionStmt(&node.Init, fset, addLocs)
	assign := visitStmt(&node.Assign, fset, addLocs)
	block := visitCaseClauseList(node.Body, fset, addLocs)

	return outputRascalString(node, fset, "typeSwitchStmt", []string{init, assign, block}, addLocs)
}

func clauseToRascal(node *ast.Stmt, fset *token.FileSet, addLocs bool) string {
	if *node == nil {
		return "defaultComm()"
	} else {
		return fmt.Sprintf("regularComm(%s)", visitStmt(node, fset, addLocs))
	}
}

func visitCommClause(node *ast.CommClause, fset *token.FileSet, addLocs bool) string {
	clauseStr := clauseToRascal(&node.Comm, fset, addLocs)
	stmtsString := visitStmtList(node.Body, fset, addLocs)

	return outputRascalString(node, fset, "commClause", []string{clauseStr, stmtsString}, addLocs)
}

func visitCommClauseList(node *ast.BlockStmt, fset *token.FileSet, addLocs bool) string {
	var clauses []string
	for i := 0; node != nil && node.List != nil && i < len(node.List); i++ {
		nodeAsClause, ok := node.List[i].(*ast.CommClause)
		if ok {
			clauses = append(clauses, visitCommClause(nodeAsClause, fset, addLocs))
		} else {
			clauses = append(clauses, "invalidCommClause()")
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(clauses, ","))
}

func visitSelectStmt(node *ast.SelectStmt, fset *token.FileSet, addLocs bool) string {
	clauses := visitCommClauseList(node.Body, fset, addLocs)

	return outputRascalString(node, fset, "selectStmt", []string{clauses}, addLocs)
}

func visitForStmt(node *ast.ForStmt, fset *token.FileSet, addLocs bool) string {
	init := visitOptionStmt(&node.Init, fset, addLocs)
	post := visitOptionStmt(&node.Post, fset, addLocs)
	cond := visitOptionExpr(&node.Cond, fset, addLocs)
	block := visitBlockStmt(node.Body, fset, addLocs)

	return outputRascalString(node, fset, "forStmt", []string{init, cond, post, block}, addLocs)
}

func visitRangeStmt(node *ast.RangeStmt, fset *token.FileSet, addLocs bool) string {
	key := visitOptionExpr(&node.Key, fset, addLocs)
	value := visitOptionExpr(&node.Value, fset, addLocs)
	assignOp := assignmentOpToRascal(node.Tok)
	rangeExpr := visitExpr(&node.X, fset, addLocs)
	block := visitBlockStmt(node.Body, fset, addLocs)

	return outputRascalString(node, fset, "rangeStmt", []string{key, value, assignOp, rangeExpr, block}, addLocs)
}

func visitExpr(node *ast.Expr, fset *token.FileSet, addLocs bool) string {
	switch t := (*node).(type) {

	case *ast.Ident:
		return visitIdent(t, fset, addLocs)

	case *ast.Ellipsis:
		return visitEllipsis(t, fset, addLocs)

	case *ast.BasicLit:
		return visitBasicLit(t, fset, addLocs)

	case *ast.FuncLit:
		return visitFuncLit(t, fset, addLocs)

	case *ast.CompositeLit:
		return visitCompositeLit(t, fset, addLocs)

	case *ast.ParenExpr:
		return visitParenExpr(t, fset, addLocs)

	case *ast.SelectorExpr:
		return visitSelectorExpr(t, fset, addLocs)

	case *ast.IndexExpr:
		return visitIndexExpr(t, fset, addLocs)

	case *ast.IndexListExpr:
		return visitIndexListExpr(t, fset, addLocs)

	case *ast.SliceExpr:
		return visitSliceExpr(t, fset, addLocs)

	case *ast.TypeAssertExpr:
		return visitTypeAssertExpr(t, fset, addLocs)

	case *ast.CallExpr:
		return visitCallExpr(t, fset, addLocs)

	case *ast.StarExpr:
		return visitStarExpr(t, fset, addLocs)

	case *ast.UnaryExpr:
		return visitUnaryExpr(t, fset, addLocs)

	case *ast.BinaryExpr:
		return visitBinaryExpr(t, fset, addLocs)

	case *ast.KeyValueExpr:
		return visitKeyValueExpr(t, fset, addLocs)

	case *ast.ArrayType:
		return visitArrayType(t, fset, addLocs)

	case *ast.StructType:
		return visitStructType(t, fset, addLocs)

	case *ast.FuncType:
		return visitFuncType(t, fset, addLocs)

	case *ast.InterfaceType:
		return visitInterfaceType(t, fset, addLocs)

	case *ast.MapType:
		return visitMapType(t, fset, addLocs)

	case *ast.ChanType:
		return visitChanType(t, fset, addLocs)

	default:
		return fmt.Sprintf("unknownExpr(\"%s\")", reflect.TypeOf(node).Name())
	}
}

func visitIdent(node *ast.Ident, fset *token.FileSet, addLocs bool) string {
	name := node.Name

	return outputRascalString(node, fset, "ident", []string{name}, addLocs)
}

func visitEllipsis(node *ast.Ellipsis, fset *token.FileSet, addLocs bool) string {
	elt := visitOptionExpr(&node.Elt, fset, addLocs)

	return outputRascalString(node, fset, "ellipsis", []string{elt}, addLocs)
}

func literalToRascal(node *ast.BasicLit, fset *token.FileSet, addLocs bool) string {
	switch node.Kind {
	case token.INT:
		if parsed, err := strconv.ParseInt(node.Value, 0, 64); err == nil {
			return outputRascalString(node, fset, "literalInt", []string{strconv.FormatInt(parsed, 10)}, addLocs)
		} else if parsed, err := strconv.ParseUint(node.Value, 0, 64); err == nil {
			return outputRascalString(node, fset, "literalInt", []string{strconv.FormatUint(parsed, 10)}, addLocs)
		} else if parsed, err := strconv.ParseFloat(node.Value, 64); err == nil {
			// NOTE: This is here because we can include an INT literal as an argument
			// when setting up a floating-point number that is outside the possible
			// bounds of an int64, e.g.,340282346638528860000000000000000000000. This
			// is coerced into a float in that case, but cannot be represented as either
			// an int32 or an int64.
			return outputRascalString(
				node,
				fset,
				"literalFloat",
				[]string{strconv.FormatFloat(parsed, 'f', -1, 64)},
				addLocs,
			)
		} else {
			return outputRascalString(node, fset, "unknownLiteral", []string{"\"" + node.Value + "\""}, addLocs)
		}
	case token.FLOAT:
		if parsed, err := strconv.ParseFloat(node.Value, 64); err == nil {
			return outputRascalString(
				node,
				fset,
				"literalFloat",
				[]string{strconv.FormatFloat(parsed, 'f', -1, 64)},
				addLocs,
			)
		} else {
			return outputRascalString(node, fset, "unknownLiteral", []string{"\"" + node.Value + "\""}, addLocs)
		}
	case token.CHAR:
		return outputRascalString(node, fset, "literalChar", []string{rascalizeChar(node.Value)}, addLocs)
	case token.STRING:
		return outputRascalString(node, fset, "literalString", []string{rascalizeString(node.Value)}, addLocs)
	case token.IMAG:
		if ic, err := strconv.ParseComplex(node.Value, 128); err == nil {
			return outputRascalString(
				node,
				fset,
				"literalImaginary",
				[]string{
					strconv.FormatFloat(real(ic), 'f', -1, 64),
					strconv.FormatFloat(imag(ic), 'f', -1, 64),
				},
				addLocs,
			)
		} else {
			return outputRascalString(node, fset, "unknownLiteral", []string{"\"" + node.Value + "\""}, addLocs)
		}
	default:
		return outputRascalString(node, fset, "unknownLiteral", []string{"\"" + node.Value + "\""}, addLocs)
	}
}

func visitBasicLit(node *ast.BasicLit, fset *token.FileSet, addLocs bool) string {
	value := literalToRascal(node, fset, addLocs)

	return outputRascalString(node, fset, "basicLit", []string{value}, addLocs)
}

func visitFuncLit(node *ast.FuncLit, fset *token.FileSet, addLocs bool) string {
	typeStr := visitFuncType(node.Type, fset, addLocs)
	body := visitBlockStmt(node.Body, fset, addLocs)

	return outputRascalString(node, fset, "funcLit", []string{typeStr, body}, addLocs)
}

func boolToRascal(val bool) string {
	if val {
		return "true"
	} else {
		return "false"
	}
}

func visitCompositeLit(node *ast.CompositeLit, fset *token.FileSet, addLocs bool) string {
	typeStr := visitOptionExpr(&node.Type, fset, addLocs)
	elts := visitExprList(node.Elts, fset, addLocs)

	return outputRascalString(
		node,
		fset,
		"compositeLit",
		[]string{typeStr, elts, boolToRascal(node.Incomplete)},
		addLocs,
	)
}

// todo this is definitely going to need to turn into a real part of the grammar if
// this thing is going to be properly bidirectional
func visitParenExpr(node *ast.ParenExpr, fset *token.FileSet, addLocs bool) string {
	// We are building a tree, we do not need to keep explicit parens
	return visitExpr(&node.X, fset, addLocs)
}

func visitSelectorExpr(node *ast.SelectorExpr, fset *token.FileSet, addLocs bool) string {
	exprStr := visitExpr(&node.X, fset, addLocs)

	return outputRascalString(node, fset, "selectorExpr", []string{exprStr, node.Sel.Name}, addLocs)
}

func visitIndexExpr(node *ast.IndexExpr, fset *token.FileSet, addLocs bool) string {
	x := visitExpr(&node.X, fset, addLocs)
	index := visitExpr(&node.Index, fset, addLocs)

	return outputRascalString(node, fset, "indexExpr", []string{x, index}, addLocs)
}

func visitIndexListExpr(node *ast.IndexListExpr, fset *token.FileSet, addLocs bool) string {
	x := visitExpr(&node.X, fset, addLocs)
	indices := visitExprList(node.Indices, fset, addLocs)

	return outputRascalString(node, fset, "indexListExpr", []string{x, indices}, addLocs)
}

func visitSliceExpr(node *ast.SliceExpr, fset *token.FileSet, addLocs bool) string {
	x := visitExpr(&node.X, fset, addLocs)
	low := visitOptionExpr(&node.Low, fset, addLocs)
	high := visitOptionExpr(&node.High, fset, addLocs)
	maxI := visitOptionExpr(&node.Max, fset, addLocs)

	return outputRascalString(
		node,
		fset,
		"sliceExpr",
		[]string{x, low, high, maxI, boolToRascal(node.Slice3)},
		addLocs,
	)
}

func visitTypeAssertExpr(node *ast.TypeAssertExpr, fset *token.FileSet, addLocs bool) string {
	x := visitExpr(&node.X, fset, addLocs)
	types := visitOptionExpr(&node.Type, fset, addLocs)

	return outputRascalString(node, fset, "typeAssertExpr", []string{x, types}, addLocs)
}

func visitCallExpr(node *ast.CallExpr, fset *token.FileSet, addLocs bool) string {
	fun := visitExpr(&node.Fun, fset, addLocs)
	args := visitExprList(node.Args, fset, addLocs)
	hasEllipses := boolToRascal(node.Ellipsis != token.NoPos)

	return outputRascalString(node, fset, "callExpr", []string{fun, args, hasEllipses}, addLocs)
}

func visitStarExpr(node *ast.StarExpr, fset *token.FileSet, addLocs bool) string {
	x := visitExpr(&node.X, fset, addLocs)

	return outputRascalString(node, fset, "starExpr", []string{x}, addLocs)
}

func visitUnaryExpr(node *ast.UnaryExpr, fset *token.FileSet, addLocs bool) string {
	x := visitExpr(&node.X, fset, addLocs)
	tok := opToRascal(node.Op)

	return outputRascalString(node, fset, "unaryExpr", []string{x, tok}, addLocs)
}

func visitBinaryExpr(node *ast.BinaryExpr, fset *token.FileSet, addLocs bool) string {
	x := visitExpr(&node.X, fset, addLocs)
	tok := opToRascal(node.Op)
	y := visitExpr(&node.Y, fset, addLocs)

	return outputRascalString(node, fset, "binaryExpr", []string{x, y, tok}, addLocs)
}

func visitKeyValueExpr(node *ast.KeyValueExpr, fset *token.FileSet, addLocs bool) string {
	key := visitExpr(&node.Key, fset, addLocs)
	value := visitExpr(&node.Value, fset, addLocs)

	return outputRascalString(node, fset, "keyValueExpr", []string{key, value}, addLocs)
}

func visitArrayType(node *ast.ArrayType, fset *token.FileSet, addLocs bool) string {
	lens := visitOptionExpr(&node.Len, fset, addLocs)
	elt := visitExpr(&node.Elt, fset, addLocs)

	return outputRascalString(node, fset, "arrayType", []string{lens, elt}, addLocs)
}

// todo that optional wrapper sounds real good
func visitOptionBasicLiteral(literal *ast.BasicLit, fset *token.FileSet, addLocs bool) string {
	if literal == nil {
		return "noLiteral()"
	} else {
		return fmt.Sprintf("someLiteral(%s)", literalToRascal(literal, fset, addLocs))
	}
}

func fieldToRascal(field *ast.Field, fset *token.FileSet, addLocs bool) string {
	names := Map(field.Names, func(t *ast.Ident) string {
		return fmt.Sprintf("\"%s\"", t.Name)
	})
	namesStr := fmt.Sprintf("[%s]", strings.Join(names, ","))
	fieldType := visitOptionExpr(&field.Type, fset, addLocs)
	fieldTag := visitOptionBasicLiteral(field.Tag, fset, addLocs)

	return outputRascalString(field, fset, "field", []string{namesStr, fieldType, fieldTag}, addLocs)
}

func visitFieldList(fieldList *ast.FieldList, fset *token.FileSet, addLocs bool) string {
	if fieldList == nil {
		return "[]"
	}

	fields := Map(fieldList.List, func(t *ast.Field) string {
		return fieldToRascal(t, fset, addLocs)
	})

	return fmt.Sprintf("[%s]", strings.Join(fields, ","))
}

func visitStructType(node *ast.StructType, fset *token.FileSet, addLocs bool) string {
	fieldStr := visitFieldList(node.Fields, fset, addLocs)

	return outputRascalString(node, fset, "structType", []string{fieldStr}, addLocs)
}

func visitFuncType(node *ast.FuncType, fset *token.FileSet, addLocs bool) string {
	typeParams := visitFieldList(node.TypeParams, fset, addLocs)
	params := visitFieldList(node.Params, fset, addLocs)
	returns := visitFieldList(node.Results, fset, addLocs)

	return outputRascalString(node, fset, "funcType", []string{typeParams, params, returns}, addLocs)
}

func visitInterfaceType(node *ast.InterfaceType, fset *token.FileSet, addLocs bool) string {
	methods := visitFieldList(node.Methods, fset, addLocs)

	return outputRascalString(node, fset, "interfaceType", []string{methods}, addLocs)
}

func visitMapType(node *ast.MapType, fset *token.FileSet, addLocs bool) string {
	key := visitExpr(&node.Key, fset, addLocs)
	value := visitExpr(&node.Value, fset, addLocs)

	return outputRascalString(node, fset, "mapType", []string{key, value}, addLocs)
}

func channelDirToRascal(dir ast.ChanDir) string {
	if (dir&ast.SEND == ast.SEND) && (dir&ast.RECV == ast.RECV) {
		return "bidirectional()"
	} else if dir&ast.SEND == ast.SEND {
		return "send()"
	} else if dir&ast.RECV == ast.RECV {
		return "receive()"
	} else {
		return "unknown()"
	}
}

func rascalChannelDirToGo(dir string) ast.ChanDir {
	switch dir {
	case "bidirectional()":
		return ast.SEND | ast.RECV
	case "send()":
		return ast.SEND
	case "receive()":
		return ast.RECV
	default:
		panic("unknown channel dir")
	}
}

func visitChanType(node *ast.ChanType, fset *token.FileSet, addLocs bool) string {
	value := visitExpr(&node.Value, fset, addLocs)
	chanSend := channelDirToRascal(node.Dir)

	return outputRascalString(node, fset, "chanType", []string{value, chanSend}, addLocs)
}

func visitExprList(nodes []ast.Expr, fset *token.FileSet, addLocs bool) string {
	exprs := Map(nodes, func(t ast.Expr) string {
		return visitExpr(&t, fset, addLocs)
	})
	return fmt.Sprintf("[%s]", strings.Join(exprs, ","))
}

func visitStmtList(nodes []ast.Stmt, fset *token.FileSet, addLocs bool) string {
	stmts := Map(nodes, func(t ast.Stmt) string {
		return visitStmt(&t, fset, addLocs)
	})

	return fmt.Sprintf("[%s]", strings.Join(stmts, ","))
}

// computes location of the ast nodes
func computeLocation(fset *token.FileSet, start token.Pos, end token.Pos) string {
	if start.IsValid() && end.IsValid() {
		startPos := fset.Position(start)
		//endPos := fset.Position(end)
		return fmt.Sprintf("|file://%s|(%d,%d)",
			filePath, startPos.Offset, end-start)
		//startPos.Line, startPos.Column-1, endPos.Line, endPos.Column-1)
	} else {
		return fmt.Sprintf("|file://%s|", fset.Position(start).Filename)
	}
}

// todo option wrapper is looking pretty non-optional
func labelToRascal(node *ast.Ident) string {
	if node != nil {
		return fmt.Sprintf("someLabel(\"%s\")", node.Name)
	} else {
		return "noLabel()"
	}
}

func rascalLabelToGo(rascalIdent string) *ast.Ident {
	if rascalIdent == "noLabel()" {
		return nil
	} else {
		return &ast.Ident{
			Name: strings.TrimSuffix(
				strings.TrimPrefix(rascalIdent, "someLabel("),
				")",
			),
		}
	}
}

// todo I could be crazy but I think this is a no-op?
func rascalizeString(input string) string {
	//r := strings.NewReplacer("<", "\\<", ">", "\\>", "\n", "\\n", "\t", "\\t", "\r", "\\r", "\\", "\\\\", "\"", "\\\"", "'", "\\'")
	s1, _ := strings.CutPrefix(input, "\"")
	s2, _ := strings.CutSuffix(s1, "\"")
	return fmt.Sprintf("\"%s\"", rascalizer.Replace(s2))
}

func rascalizeChar(input string) string {
	//r := strings.NewReplacer("<", "\\<", ">", "\\>", "\n", "\\n", "\t", "\\t", "\r", "\\r", "\\", "\\\\", "\"", "\\\"", "'", "\\'")
	s1, _ := strings.CutPrefix(input, "'")
	s2, _ := strings.CutSuffix(s1, "'")
	return fmt.Sprintf("\"%s\"", rascalizer.Replace(s2))
}

func main() {
	flag.StringVar(&filePath, "filePath", "", "The file to be processed")

	var addLocations bool
	flag.BoolVar(&addLocations, "addLocs", true, "Include location annotations")

	flag.Parse()

	if filePath != "" {
		//fmt.Printf("Processing file %s\n", filePath)
		fmt.Println(processFile(addLocations))
	} else {
		fmt.Println("No file given")
	}
}
