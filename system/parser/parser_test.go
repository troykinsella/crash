package parser

import (
	"testing"
	"bufio"
	"strings"
	"reflect"
	"github.com/troykinsella/crash/system/ast"
	"fmt"
	"github.com/troykinsella/crash/system/token"
)

func TestParser_Parse(t *testing.T) {
	var tests = []struct {
		str string
		n *ast.Statement
		err string
	} {
		{
			"foo op",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "foo",
								},
							},
						},
					},
					false,
					"op",
					nil,
				},
				nil,
			},
			"",
		}, {
			"foo[bar] op",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Index: &ast.IndexExpr{
										Operand: &ast.PrimaryExpr{
											Ident: "foo",
										},
										Index: &ast.PrimaryExpr{
											Ident: "bar",
										},
									},
								},
							},
						},
					},
					false,
					"op",
					nil,
				},
				nil,
			},
			"",
		}, {
			"foo.bar op",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Selector: &ast.SelectorExpr{
										Operand: &ast.PrimaryExpr{
											Ident: "foo",
										},
										Ident: "bar",
									},
								},
							},
						},
					},
					false,
					"op",
					nil,
				},
				nil,
			},
			"",
		}, {
			"foo.bar[baz] op",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Index: &ast.IndexExpr{
										Operand: &ast.PrimaryExpr{
											Selector: &ast.SelectorExpr{
												Operand: &ast.PrimaryExpr{
													Ident: "foo",
												},
												Ident: "bar",
											},
										},
										Index: &ast.PrimaryExpr{
											Ident: "baz",
										},
									},
								},
							},
						},
					},
					false,
					"op",
					nil,
				},
				nil,
			},
			"",
		}, {
			"foo[bar].baz op",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Selector: &ast.SelectorExpr{
										Operand: &ast.PrimaryExpr{
											Index: &ast.IndexExpr{
												Operand: &ast.PrimaryExpr{
													Ident: "foo",
												},
												Index: &ast.PrimaryExpr{
													Ident: "bar",
												},
											},
										},
										Ident: "baz",
									},
								},
							},
						},
					},
					false,
					"op",
					nil,
				},
				nil,
			},
			"",
		}, {
			"foo.bar[baz.biz] op",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Index: &ast.IndexExpr{
										Operand: &ast.PrimaryExpr{
											Selector: &ast.SelectorExpr{
												Operand: &ast.PrimaryExpr{
													Ident: "foo",
												},
												Ident: "bar",
											},
										},
										Index: &ast.PrimaryExpr{
											Selector: &ast.SelectorExpr{
												Operand: &ast.PrimaryExpr{
													Ident: "baz",
												},
												Ident: "biz",
											},
										},
									},
								},
							},
						},
					},
					false,
					"op",
					nil,
				},
				nil,
			},
			"",
		}, {
			"foo.bar[baz[biz]] op",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Index: &ast.IndexExpr{
										Operand: &ast.PrimaryExpr{
											Selector: &ast.SelectorExpr{
												Operand: &ast.PrimaryExpr{
													Ident: "foo",
												},
												Ident: "bar",
											},
										},
										Index: &ast.PrimaryExpr{
											Index: &ast.IndexExpr{
												Operand: &ast.PrimaryExpr{
													Ident: "baz",
												},
												Index: &ast.PrimaryExpr{
													Ident: "biz",
												},
											},
										},
									},
								},
							},
						},
					},
					false,
					"op",
					nil,
				},
				nil,
			},
			"",
		}, {
			"foo[bar].baz op",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Selector: &ast.SelectorExpr{
										Operand: &ast.PrimaryExpr{
											Index: &ast.IndexExpr{
												Operand: &ast.PrimaryExpr{
													Ident: "foo",
												},
												Index: &ast.PrimaryExpr{
													Ident: "bar",
												},
											},
										},
										Ident: "baz",
									},
								},
							},
						},
					},
					false,
					"op",
					nil,
				},
				nil,
			},
			"",
		}, {
			"foo.bar.baz op",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Selector: &ast.SelectorExpr{
										Operand: &ast.PrimaryExpr{
											Selector: &ast.SelectorExpr{
												Operand: &ast.PrimaryExpr{
													Ident: "foo",
												},
												Ident: "bar",
											},
										},
										Ident: "baz",
									},
								},
							},
						},
					},
					false,
					"op",
					nil,
				},
				nil,
			},
			"",
		}, {
			"foo op bar",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "foo",
								},
							},
						},
					},
					false,
					"op",
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "bar",
								},
							},
						},
					},
				},
				nil,
			},
			"",
		}, {
			"foo op 123",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "foo",
								},
							},
						},
					},
					false,
					"op",
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Literal: &ast.Literal{
										Type: token.NUMBER,
										Int: 123,
									},
								},
							},
						},
					},
				},
				nil,
			},
			"",
		}, {
			"foo op 'bar'",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "foo",
								},
							},
						},
					},
					false,
					"op",
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Literal: &ast.Literal{
										Type: token.STRING,
										Str: "bar",
									},
								},
							},
						},
					},
				},
				nil,
			},
			"",
		}, {
			"foo op bar[baz]",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "foo",
								},
							},
						},
					},
					false,
					"op",
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Index: &ast.IndexExpr{
										Operand: &ast.PrimaryExpr{
											Ident: "bar",
										},
										Index: &ast.PrimaryExpr{
											Ident: "baz",
										},
									},
								},
							},
						},
					},
				},
				nil,
			},
			"",
		}, {
			"foo op bar.baz",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "foo",
								},
							},
						},
					},
					false,
					"op",
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Selector: &ast.SelectorExpr{
										Operand: &ast.PrimaryExpr{
											Ident: "bar",
										},
										Ident: "baz",
									},
								},
							},
						},
					},
				},
				nil,
			},
			"",
		}, {
			"foo not op bar",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "foo",
								},
							},
						},
					},
					true,
					"op",
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "bar",
								},
							},
						},
					},
				},
				nil,
			},
			"",
		}, {
			"foo op bar, 'baz'",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "foo",
								},
							},
						},
					},
					false,
					"op",
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "bar",
								},
							},
							&ast.Expression{
								&ast.PrimaryExpr{
									Literal: &ast.Literal{
										Type: token.STRING,
										Str: "baz",
									},
								},
							},
						},
					},
				},
				nil,
			},
			"",
		}, {
			"foo, bar op baz",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "foo",
								},
							},
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "bar",
								},
							},
						},
					},
					false,
					"op",
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "baz",
								},
							},
						},
					},
				},
				nil,
			},
			"",
		},
		{
			"foo op // bar",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "foo",
								},
							},
						},
					},
					false,
					"op",
					nil,
				},
				&ast.IString{
					Str: "bar",
				},
			},
			"",
		}, {
			"foo, bar op baz, 'biz'",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "foo",
								},
							},
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "bar",
								},
							},
						},
					},
					false,
					"op",
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "baz",
								},
							},
							&ast.Expression{
								&ast.PrimaryExpr{
									Literal: &ast.Literal{
										Type: token.STRING,
										Str: "biz",
									},
								},
							},
						},
					},
				},
				nil,
			},
			"",
		}, {
			"foo, bar not op baz, 'biz', 123",
			&ast.Statement{
				&ast.Operation{
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "foo",
								},
							},
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "bar",
								},
							},
						},
					},
					true,
					"op",
					&ast.ExpressionList{
						[]*ast.Expression{
							&ast.Expression{
								&ast.PrimaryExpr{
									Ident: "baz",
								},
							},
							&ast.Expression{
								&ast.PrimaryExpr{
									Literal: &ast.Literal{
										Type: token.STRING,
										Str: "biz",
									},
								},
							},
							&ast.Expression{
								&ast.PrimaryExpr{
									Literal: &ast.Literal{
										Type: token.NUMBER,
										Int: 123,
									},
								},
							},
						},
					},
				},
				nil,
			},
			"",
		},
	}

	for i, test := range tests {
		p := New(bufio.NewReader(strings.NewReader(test.str)))
		n, err := p.Statement()

		if test.err == "" {
			if err != nil {
				t.Errorf("%d. unexpected error: %s\n", i, err.Error())
	 		} else if !reflect.DeepEqual(test.n, n) {
				fmt.Println("Expected:")
				test.n.Dump(0)
				fmt.Println("Actual:")
				n.Dump(0)
				t.Errorf("%d. \"%s\" unexpected ast:\nexpected=%#v,\nactual=%#v\n", i, test.str, test.n, n)
			}
		} else {
			if err == nil {
				t.Errorf("%d. \"%s\" expected error:\nexpected=%s,\nactual=nil\n", i, test.str, test.err)
			} else if test.err != err.Error() {
				t.Errorf("%d. \"%s\" unexpected error:\nexpected=%s,\nactual=%s\n", i, test.str, test.err, err.Error())
			}
		}
	}
}
