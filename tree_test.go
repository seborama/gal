package gal_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/seborama/gal"
)

func TestTree_Eval_Expressions(t *testing.T) {
	tt := map[string]struct {
		tree gal.Tree
		want gal.Value
	}{
		"starts with *": {
			tree: gal.Tree{
				gal.Multiply,
				gal.NewNumber(-4),
			},
			want: gal.NewUndefinedWithReasonf("syntax error: expression starts with '*'"),
		},
		"starts with + -4": {
			tree: gal.Tree{
				gal.Plus,
				gal.NewNumber(-4),
			},
			want: gal.NewNumber(-4),
		},
		"starts with - -4": {
			tree: gal.Tree{
				gal.Minus,
				gal.NewNumber(-4),
			},
			want: gal.NewNumber(4),
		},
		"chained * and /": {
			tree: gal.Tree{
				// 3 * 4 / 2 / 3 * 4
				gal.NewNumber(3),
				gal.Multiply,
				gal.NewNumber(4),
				gal.Divide,
				gal.NewNumber(2),
				gal.Divide,
				gal.NewNumber(3),
				gal.Multiply,
				gal.NewNumber(4),
			},
			want: gal.NewNumber(8),
		},
		"chained and tree'ed * and /": {
			tree: gal.Tree{
				// (((3)) * (4)) / (2) / (3) * (4)
				gal.Tree{
					gal.Tree{
						gal.Tree{
							gal.NewNumber(3),
						},
					},
					gal.Multiply,
					gal.Tree{
						gal.NewNumber(4),
					},
				},
				gal.Divide,
				gal.Tree{
					gal.NewNumber(2),
				},
				gal.Divide,
				gal.Tree{
					gal.NewNumber(3),
				},
				gal.Multiply,
				gal.Tree{
					gal.NewNumber(4),
				},
			},
			want: gal.NewNumber(8),
		},
		"rich tree": {
			tree: gal.Tree{
				// 3 - 4 * (-2) - 5 => 3 - ( 4 * (-2) ) - 5
				gal.NewNumber(3),
				gal.Minus,
				gal.NewNumber(4),
				gal.Multiply,
				gal.Tree{
					gal.Minus,
					gal.NewNumber(2),
				},
				gal.Minus,
				gal.NewNumber(5),
			},
			want: gal.NewNumber(6),
		},
		"multiple levels of decreasing operator precedence": {
			tree: gal.Tree{
				// 10 ^ 2 * 4 + 3 => 10 ^ 2 * 4 + 3
				gal.NewNumber(10),
				gal.Power,
				gal.NewNumber(2),
				gal.Multiply,
				gal.NewNumber(4),
				gal.Plus,
				gal.NewNumber(3),
			},
			want: gal.NewNumber(403),
		},
		"multiple levels of operator precedence": {
			tree: gal.Tree{
				// 10 + 5 * 4 ^ 3 * 2 + 6 * 7 => 10 + ( 5 * ( 4 ^ 3 ) * 2 ) + ( 6 * 7 )
				gal.NewNumber(10),
				gal.Plus,
				gal.NewNumber(5),
				gal.Multiply,
				gal.NewNumber(4),
				gal.Power,
				gal.NewNumber(3),
				gal.Multiply,
				gal.NewNumber(2),
				gal.Plus,
				gal.NewNumber(6),
				gal.Multiply,
				gal.NewNumber(7),
			},
			want: gal.NewNumber(692),
		},
		"rich sub-trees": {
			tree: gal.Tree{
				gal.NewNumber(10),
				gal.Plus,
				gal.Tree{
					gal.NewNumber(5),
					gal.Multiply,
					gal.Tree{
						gal.Minus,
						gal.NewNumber(4),
						gal.Modulus,
						gal.Tree{
							gal.Minus,
							gal.NewNumber(3),
						},
					},
				},
			},
			want: gal.NewNumber(5),
		},
		"function": {
			tree: gal.Tree{
				gal.NewNumber(10),
				gal.Plus,
				gal.NewFunction(
					"trunc",
					gal.Trunc,
					gal.Tree{
						gal.NewFunction(
							"sqrt",
							gal.Sqrt,
							gal.Tree{
								gal.NewNumber(10),
							},
						),
					},
					gal.Tree{
						gal.NewNumber(6),
					},
				),
			},
			want: gal.NewNumberFromFloat(13.162277),
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			val := tc.tree.Eval()

			if !cmp.Equal(tc.want, val) {
				if _, ok := val.(gal.Undefined); ok {
					t.Log("Value:", val.String())
				}
				t.Errorf(cmp.Diff(tc.want, val))
			}
		})
	}
}
