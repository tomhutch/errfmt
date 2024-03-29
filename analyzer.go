package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"reflect"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type Config struct {
	// ErrWrapSigs allows you to specify additional functions that errfmt will validate.
	//
	// For example, an errWrapSig of `[]string{"stacktrace.Propagate("}` will validate errors
	// messages passed to the stacktrace package's Propagate function:
	ErrWrapSigs []string `mapstructure:"ignoreSigs" yaml:"ignoreSigs"`
}

func NewDefaultConfig() Config {
	return Config{
		ErrWrapSigs: []string{
			".Errorf(",
			"errors.New(",
			"errors.Unwrap(",
			".Wrap(",
			".Wrapf(",
			".WithMessage(",
			".WithMessagef(",
			".WithStack(",
			"multierr.Append(",
		},
	}
}

func NewAnalyzer(cfg Config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "errfmt",
		Doc:  "Checks message format of wrapped errors",
		Run:  run(cfg),
	}
}

func run(cfg Config) func(*analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		for _, file := range pass.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				ret, ok := n.(*ast.ReturnStmt)
				if !ok {
					return true
				}

				if len(ret.Results) < 1 {
					return true
				}

				// Iterate over the values to be returned looking for errors
				for _, expr := range ret.Results {
					// Check if the return expression is a function call, if it is, we need
					// to handle it by checking the return params of the function.
					retFn, ok := expr.(*ast.CallExpr)
					if ok {
						// If the return type of the function is a single error. This will not
						// match an error within multiple return values, for that, the below
						// tuple check is required.

						if isError(pass.TypesInfo.TypeOf(expr)) {
							checkCallExpr(pass, retFn, retFn.Pos(), cfg, file)
							return true
						}

						// Check if one of the return values from the function is an error
						tup, ok := pass.TypesInfo.TypeOf(expr).(*types.Tuple)
						if !ok {
							return true
						}

						// Iterate over the return values of the function looking for error
						// types
						for i := 0; i < tup.Len(); i++ {
							v := tup.At(i)
							if v == nil {
								return true
							}
							if isError(v.Type()) {
								checkCallExpr(pass, retFn, expr.Pos(), cfg, file)
								return true
							}
						}
					}

					if !isError(pass.TypesInfo.TypeOf(expr)) {
						continue
					}

					ident, ok := expr.(*ast.Ident)
					if !ok {
						return true
					}

					var call *ast.CallExpr

					// Attempt to find the most recent short assign
					if shortAss := prevErrAssign(pass, file, ident); shortAss != nil {
						call, ok = shortAss.Rhs[0].(*ast.CallExpr)
						if !ok {
							return true
						}
					} else if isUnresolved(file, ident) {
						// TODO Check if the identifier is unresolved, and try to resolve it in
						// another file.
						return true
					} else {
						// Check for ValueSpec nodes in order to locate a possible var
						// declaration.
						if ident.Obj == nil {
							return true
						}

						vSpec, ok := ident.Obj.Decl.(*ast.ValueSpec)
						if !ok {
							// We couldn't find a short or var assign for this error return.
							// This is an error. Where did this identifier come from? Possibly a
							// function param.
							//
							// TODO decide how to handle this case, whether to follow function
							// param back, or assert wrapping at call site.

							return true
						}

						if len(vSpec.Values) < 1 {
							return true
						}

						call, ok = vSpec.Values[0].(*ast.CallExpr)
						if !ok {
							return true
						}
					}

					// Make sure there is a call identified as producing the error being
					// returned, otherwise just bail
					if call == nil {
						return true
					}

					checkCallExpr(pass, call, ident.NamePos, cfg, file)
				}

				return true
			})
		}

		return nil, nil
	}
}

func checkCallExpr(pass *analysis.Pass, call *ast.CallExpr, tokenPos token.Pos, cfg Config, file *ast.File) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	fnSig := pass.TypesInfo.ObjectOf(sel.Sel).String()
	if contains(cfg.ErrWrapSigs, fnSig) {
		if len(call.Args) > 0 {
			// Find upstream function call that assigned the error.
			upstreamFnCall := findUpstreamFnCall(pass, call, file)
			if upstreamFnCall == nil {
				return
			}

			prefixParts := crawlExprSelectorChain(pass, upstreamFnCall.Fun)

			arg, ok := call.Args[0].(*ast.BasicLit)
			if ok {
				validateErrorFormat(pass, prefixParts, arg, tokenPos)
			}
		}
	}
}

func validateErrorFormat(pass *analysis.Pass, prefixParts []string, arg *ast.BasicLit, tokenPos token.Pos) {
	expected := fmt.Sprintf("\"%s", strings.Join(prefixParts, "."))

	// Inspect message passed to error wrapping function (e.g. fmt.Sprintf("error message: %w", err)).
	if strings.HasPrefix(arg.Value, expected) {
		return
	}

	msg := fmt.Sprintf("error message not prefixed in expected format (value = [%s], expected = \"[%s\"])", arg.Value, expected)
	suggestedFixText := fmt.Sprintf("%s: %%w\"", expected)
	suggestedFix := analysis.SuggestedFix{
		Message: "try this",
		TextEdits: []analysis.TextEdit{
			{
				Pos:     arg.Pos(),
				End:     arg.End(),
				NewText: []byte(suggestedFixText),
			},
		},
	}
	pass.Report(analysis.Diagnostic{
		Pos:            arg.Pos(),
		End:            0,
		Category:       "",
		Message:        msg,
		SuggestedFixes: []analysis.SuggestedFix{suggestedFix},
		Related:        nil,
	})
}

func findUpstreamFnCall(pass *analysis.Pass, call *ast.CallExpr, file *ast.File) *ast.CallExpr {
	var upstreamFnCall *ast.CallExpr
	for _, arg := range call.Args {
		if isError(pass.TypesInfo.TypeOf(arg)) {
			if ident, ok := arg.(*ast.Ident); ok {
				assignmentStmt := prevErrAssign(pass, file, ident)
				if assignmentStmt != nil {
					if callExpr, ok := assignmentStmt.Rhs[0].(*ast.CallExpr); ok {
						upstreamFnCall = callExpr
					}
				}
			}
		}
	}
	return upstreamFnCall
}

func findSelectorExprInCallExprChain(pass *analysis.Pass, fn *ast.CallExpr) *ast.SelectorExpr {
	fnSelExpr, ok := fn.Fun.(*ast.SelectorExpr)
	if ok {
		return fnSelExpr
	}

	rv := reflect.ValueOf(fn.Fun)
	if rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		fnCallExpr, ok := fn.Fun.(*ast.CallExpr)
		if ok {
			return findSelectorExprInCallExprChain(pass, fnCallExpr)
		}
	}

	return nil
}

func crawlExprSelectorChain(pass *analysis.Pass, sel ast.Expr) []string {
	callExpr, ok := sel.(*ast.CallExpr)
	if ok {
		upstreamFnCallSel := findSelectorExprInCallExprChain(pass, callExpr)
		return append(crawlExprSelectorChain(pass, upstreamFnCallSel.X), upstreamFnCallSel.Sel.Name)
	}

	selExpr, ok := sel.(*ast.SelectorExpr)
	if ok {
		return append(crawlExprSelectorChain(pass, selExpr.X), selExpr.Sel.Name)
	}

	ident, ok := sel.(*ast.Ident)
	if ok {
		return []string{ident.Name}
	}

	basicLit, ok := sel.(*ast.BasicLit)
	if ok {
		return []string{basicLit.Value}
	}

	indexExpr, ok := sel.(*ast.IndexExpr)
	if ok {
		return formatIndexExpr(
			crawlExprSelectorChain(pass, indexExpr.X),
			crawlExprSelectorChain(pass, indexExpr.Index),
		)
	}

	return nil
}

// formatIndexExpr appends a formatted index expr index (i.e. a BasicLit map key) to the end of the index expression.
func formatIndexExpr(indexExprX []string, indexExprIndexes []string) []string {
	var index string
	if len(indexExprIndexes) == 1 {
		index = indexExprIndexes[0]
		if strings.HasPrefix(index, "\"") {
			index = strings.TrimPrefix(index, "\"")
		}
		if strings.HasSuffix(index, "\"") {
			index = strings.TrimSuffix(index, "\"")
		}
		index = fmt.Sprintf("[%s]", index)
	}
	indexExprX[len(indexExprX)-1] = indexExprX[len(indexExprX)-1] + index
	return indexExprX
}

// prevErrAssign traverses the AST of a file looking for the most recent
// assignment to an error declaration which is specified by the returnIdent
// identifier.
//
// This only returns short form assignments and reassignments, e.g. `:=` and
// `=`. This does not include `var` statements. This function will return nil if
// the only declaration is a `var` (aka ValueSpec) declaration.
func prevErrAssign(pass *analysis.Pass, file *ast.File, returnIdent *ast.Ident) *ast.AssignStmt {
	// A slice containing all the assignments which contain an identifier
	// referring to the source declaration of the error. This is to catch
	// cases where err is defined once, and then reassigned multiple times
	// within the same block. In these cases, we should check the method of
	// the most recent call.
	var assigns []*ast.AssignStmt

	// Find all assignments which have the same declaration
	ast.Inspect(file, func(n ast.Node) bool {
		if ass, ok := n.(*ast.AssignStmt); ok {
			for _, expr := range ass.Lhs {
				if !isError(pass.TypesInfo.TypeOf(expr)) {
					continue
				}

				if assIdent, ok := expr.(*ast.Ident); ok {
					if assIdent.Obj == nil || returnIdent.Obj == nil {
						// If we can't find the Obj for one of the identifiers, just skip
						// it.
						return true
					} else if assIdent.Obj.Decl == returnIdent.Obj.Decl {
						assigns = append(assigns, ass)
					}
				}
			}
		}

		return true
	})

	// Iterate through the assignments, comparing the token positions to
	// find the assignment that directly precedes the return position
	var mostRecentAssign *ast.AssignStmt

	for _, ass := range assigns {
		if ass.Pos() > returnIdent.Pos() {
			break
		}

		mostRecentAssign = ass
	}

	return mostRecentAssign
}

func contains(slice []string, el string) bool {
	for _, s := range slice {
		if strings.Contains(el, s) {
			return true
		}
	}

	return false
}

// isError returns whether or not the provided type interface is an error
func isError(typ types.Type) bool {
	if typ == nil {
		return false
	}

	return typ.String() == "error"
}

func isUnresolved(file *ast.File, ident *ast.Ident) bool {
	for _, unresolvedIdent := range file.Unresolved {
		if unresolvedIdent.Pos() == ident.Pos() {
			return true
		}
	}

	return false
}
