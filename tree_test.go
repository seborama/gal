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
		times,
		NewNumber(4),
		dividedBy,
		Tree{
			NewNumber(1),
			plus,
			NewNumber(2),
			plus,
			NewNumber(3),
			plus,
			NewNumber(4),
		},
		modulus,
		NewNumber(6),
		plus,
		NewFunction(
			"log",
			Tree{
				NewNumber(10),
			},
		),
	}
	// from: "-3 + -4" + -3 * 4 / (1 + 2 + 3 + 4) % 6 + log( 10 )
	//   to: "-3 + -4" + ( -3 * 4 / (1 + 2 + 3 + 4) % 6 ) + log( 10 )

	outTree := inTree.PrioritiseOperators()

	expectedTree := Tree{
		NewString(`"-3 + -4"`),
		plus,
		minus,
		Tree{
			NewNumber(3),
			times,
			NewNumber(4),
			dividedBy,
			Tree{
				NewNumber(1),
				plus,
				NewNumber(2),
				plus,
				NewNumber(3),
				plus,
				NewNumber(4),
			},
			modulus,
			NewNumber(6),
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
