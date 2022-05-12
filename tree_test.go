package gal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTree_Eval_Expressions(t *testing.T) {
	tt := map[string]struct {
		tree Tree
		want Value
	}{
		"starts with *": {
			tree: Tree{
				multiply,
				NewNumber(-4),
			},
			want: Undefined{reason: "syntax error: expression starts with '*'"},
		},
		"starts with + -4": {
			tree: Tree{
				plus,
				NewNumber(-4),
			},
			want: NewNumber(-4),
		},
		"starts with - -4": {
			tree: Tree{
				minus,
				NewNumber(-4),
			},
			want: NewNumber(4),
		},
		"chained * and /": {
			tree: Tree{
				// 3 * 4 / 2 / 3 * 4
				NewNumber(3),
				multiply,
				NewNumber(4),
				divide,
				NewNumber(2),
				divide,
				NewNumber(3),
				multiply,
				NewNumber(4),
			},
			want: NewNumber(8),
		},
		"chained and tree'ed * and /": {
			tree: Tree{
				// (((3)) * (4)) / (2) / (3) * (4)
				Tree{
					Tree{
						Tree{
							NewNumber(3),
						},
					},
					multiply,
					Tree{
						NewNumber(4),
					},
				},
				divide,
				Tree{
					NewNumber(2),
				},
				divide,
				Tree{
					NewNumber(3),
				},
				multiply,
				Tree{
					NewNumber(4),
				},
			},
			want: NewNumber(8),
		},
		"rich tree": {
			tree: Tree{
				// 3 - 4 * (-2) - 5 => 3 - ( 4 * (-2) ) - 5
				NewNumber(3),
				minus,
				NewNumber(4),
				multiply,
				Tree{
					minus,
					NewNumber(2),
				},
				minus,
				NewNumber(5),
			},
			want: NewNumber(6),
		},
		"multiple levels of decreasing operator precedence": {
			tree: Tree{
				// 10 ^ 2 * 4 + 3 => 10 ^ 2 * 4 + 3
				NewNumber(10),
				power,
				NewNumber(2),
				multiply,
				NewNumber(4),
				plus,
				NewNumber(3),
			},
			want: NewNumber(403),
		},
		"multiple levels of operator precedence": {
			tree: Tree{
				// 10 + 5 * 4 ^ 3 * 2 + 6 * 7 => 10 + ( 5 * ( 4 ^ 3 ) * 2 ) + ( 6 * 7 )
				NewNumber(10),
				plus,
				NewNumber(5),
				multiply,
				NewNumber(4),
				power,
				NewNumber(3),
				multiply,
				NewNumber(2),
				plus,
				NewNumber(6),
				multiply,
				NewNumber(7),
			},
			want: NewNumber(692),
		},
		"rich sub-trees": {
			tree: Tree{
				NewNumber(10),
				plus,
				Tree{
					NewNumber(5),
					multiply,
					Tree{
						minus,
						NewNumber(4),
						modulus,
						Tree{
							minus,
							NewNumber(3),
						},
					},
				},
			},
			want: NewNumber(5),
		},
		"function": {
			tree: Tree{
				NewNumber(10),
				plus,
				NewFunction(
					"trunc",
					Tree{
						NewNumber(6),
					},
					Tree{
						NewFunction(
							"sqrt",
							Tree{
								NewNumber(10),
							},
						),
					},
				),
			},
			want: NewNumberFromFloat(13.162277),
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			val := tc.tree.Eval()

			if !cmp.Equal(tc.want, val) {
				if _, ok := val.(Undefined); ok {
					t.Log("Value:", val.String())
				}
				t.Errorf(cmp.Diff(tc.want, val))
			}
		})
	}
}
