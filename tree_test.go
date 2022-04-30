package gal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTree_PrioritiseOperators(t *testing.T) {
	inTree := Tree{
		NewString(`"-3 + -4"`),
		plus,
		minus,
		NewNumber(3),
		multiply,
		NewNumber(4),
		divide,
		Tree{
			NewNumber(11),
			plus,
			NewNumber(12),
			plus,
			NewNumber(13),
			plus,
			NewNumber(14),
		},
		modulus,
		NewNumber(6),
		power,
		NewNumber(2),
		plus,
		NewFunction(
			"log",
			Tree{
				NewNumber(10),
			},
		),
	}
	// from: "-3 + -4" + -3 * 4 / (1 + 2 + 3 + 4) % 6 ^ 2 + log( 10 )
	//   to: "-3 + -4" + - ( 3 * 4 / (1 + 2 + 3 + 4) % ( 6 ^ 2 ) ) + log( 10 )

	outTree := inTree.prioritiseOperators()

	expectedTree := Tree{
		NewString(`"-3 + -4"`),
		plus,
		minus,
		Tree{
			NewNumber(3),
			multiply,
			NewNumber(4),
			divide,
			Tree{
				NewNumber(11),
				plus,
				NewNumber(12),
				plus,
				NewNumber(13),
				plus,
				NewNumber(14),
			},
			modulus,
			Tree{
				NewNumber(6),
				power,
				NewNumber(2),
			},
		},
		plus,
		NewFunction(
			"log",
			Tree{
				NewNumber(10),
			},
		),
	}

	if !cmp.Equal(expectedTree, outTree) {
		t.Error(cmp.Diff(expectedTree, outTree))
	}
}

func TestTree_Eval_Expressions(t *testing.T) {
	tt := map[string]struct {
		tree Tree
		want Value
	}{
		// "starts with *": {
		// 	tree: Tree{
		// 		multiply,
		// 		NewNumber(-4),
		// 	},
		// 	want: Undefined{reason: "syntax error: expression starts with '*'"},
		// },
		// "starts with + -4": {
		// 	tree: Tree{
		// 		NewNumber(-4),
		// 	},
		// 	want: NewNumber(-4),
		// },
		// "starts with - -4": {
		// 	tree: Tree{
		// 		minus,
		// 		NewNumber(-4),
		// 	},
		// 	want: NewNumber(4),
		// },
		// "rich tree": {
		// 	tree: Tree{
		// 		NewNumber(3),
		// 		minus,
		// 		NewNumber(4),
		// 		multiply,
		// 		Tree{
		// 			minus,
		// 			NewNumber(2),
		// 		},
		// 		minus,
		// 		NewNumber(5),
		// 	},
		// 	want: NewNumber(6),
		// },
		"multiple levels of operator precedence": {
			tree: Tree{
				NewNumber(10),
				plus,
				NewNumber(5),
				multiply,
				NewNumber(4),
				power,
				NewNumber(3),
			},
			want: NewNumber(330),
		},
		// "rich sub-trees": {
		// 	tree: Tree{
		// 		NewNumber(10),
		// 		plus,
		// 		Tree{
		// 			NewNumber(5),
		// 			multiply,
		// 			Tree{
		// 				minus,
		// 				NewNumber(4),
		// 				modulus,
		// 				Tree{
		// 					minus,
		// 					NewNumber(3),
		// 				},
		// 			},
		// 		},
		// 	},
		// 	want: NewNumber(5),
		// },
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			val := tc.tree.Eval()

			if !cmp.Equal(tc.want, val) {
				t.Error(cmp.Diff(tc.want, val))
			}
		})
	}
}
