package gal_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

	"github.com/seborama/gal/v10"
)

func TestTree_FullLen(t *testing.T) {
	tt := map[string]struct {
		tree    gal.Tree
		wantLen int
	}{
		"empty tree": {
			tree:    gal.Tree{},
			wantLen: 0,
		},
		"simple tree": {
			tree: gal.Tree{
				gal.Plus,
				gal.NewNumberFromInt(-4),
			},
			wantLen: 2,
		},
		"semi-complex tree": {
			tree: gal.Tree{
				gal.NewNumberFromInt(-4),
				gal.Plus,
				gal.Tree{},
			},
			wantLen: 2,
		},
		"complex tree": {
			tree: gal.Tree{
				gal.NewNumberFromInt(-4),
				gal.Plus,
				gal.Tree{
					gal.NewNumberFromInt(-4),
					gal.Plus,
					gal.Tree{
						gal.NewNumberFromInt(-4),
						gal.Plus,
						gal.Tree{
							gal.NewNumberFromInt(-4),
							gal.Plus,
							gal.Tree{
								gal.NewNumberFromInt(-4),
								gal.Plus,
								gal.Tree{
									gal.NewNumberFromInt(-4),
								},
							},
						},
					},
				},
			},
			wantLen: 11,
		},
	}

	for name, tc := range tt {
		name := name
		tc := tc

		t.Run(name, func(t *testing.T) {
			got := tc.tree.FullLen()
			assert.Equal(t, tc.wantLen, got)
		})
	}
}

func TestTree_Eval_Expressions(t *testing.T) {
	tt := map[string]struct {
		tree gal.Tree
		want gal.Value
	}{
		"starts with *": {
			tree: gal.Tree{
				gal.Multiply,
				gal.NewNumberFromInt(-4),
			},
			want: gal.NewUndefinedWithReasonf("syntax error: missing left hand side value for operator '*'"),
		},
		"starts with + -4": {
			tree: gal.Tree{
				gal.Plus,
				gal.NewNumberFromInt(-4),
			},
			want: gal.NewNumberFromInt(-4),
		},
		"starts with - -4": {
			tree: gal.Tree{
				gal.Minus,
				gal.NewNumberFromInt(-4),
			},
			want: gal.NewNumberFromInt(4),
		},
		"chained * and /": {
			tree: gal.Tree{
				gal.NewNumberFromInt(3),
				gal.Multiply,
				gal.NewNumberFromInt(4),
				gal.Divide,
				gal.NewNumberFromInt(2),
				gal.Divide,
				gal.NewNumberFromInt(3),
				gal.Multiply,
				gal.NewNumberFromInt(4),
			},
			want: gal.NewNumberFromInt(8),
		},
		"chained and tree'ed * and /": {
			tree: gal.Tree{
				gal.Tree{
					gal.Tree{
						gal.Tree{
							gal.NewNumberFromInt(3),
						},
					},
					gal.Multiply,
					gal.Tree{
						gal.NewNumberFromInt(4),
					},
				},
				gal.Divide,
				gal.Tree{
					gal.NewNumberFromInt(2),
				},
				gal.Divide,
				gal.Tree{
					gal.NewNumberFromInt(3),
				},
				gal.Multiply,
				gal.Tree{
					gal.NewNumberFromInt(4),
				},
			},
			want: gal.NewNumberFromInt(8),
		},
		"rich tree": {
			tree: gal.Tree{
				gal.NewNumberFromInt(3),
				gal.Minus,
				gal.NewNumberFromInt(4),
				gal.Multiply,
				gal.Tree{
					gal.Minus,
					gal.NewNumberFromInt(2),
				},
				gal.Minus,
				gal.NewNumberFromInt(5),
			},
			want: gal.NewNumberFromInt(6),
		},
		"multiple levels of decreasing operator precedence": {
			tree: gal.Tree{
				gal.NewNumberFromInt(10),
				gal.Power,
				gal.NewNumberFromInt(2),
				gal.Multiply,
				gal.NewNumberFromInt(4),
				gal.Plus,
				gal.NewNumberFromInt(3),
			},
			want: gal.NewNumberFromInt(403),
		},
		"multiple levels of operator precedence": {
			tree: gal.Tree{
				gal.NewNumberFromInt(10),
				gal.Plus,
				gal.NewNumberFromInt(5),
				gal.Multiply,
				gal.NewNumberFromInt(4),
				gal.Power,
				gal.NewNumberFromInt(3),
				gal.Multiply,
				gal.NewNumberFromInt(2),
				gal.Plus,
				gal.NewNumberFromInt(6),
				gal.Multiply,
				gal.NewNumberFromInt(7),
			},
			want: gal.NewNumberFromInt(692),
		},
		"rich sub-trees": {
			tree: gal.Tree{
				gal.NewNumberFromInt(10),
				gal.Plus,
				gal.Tree{
					gal.NewNumberFromInt(5),
					gal.Multiply,
					gal.Tree{
						gal.Minus,
						gal.NewNumberFromInt(4),
						gal.Modulus,
						gal.Tree{
							gal.Minus,
							gal.NewNumberFromInt(3),
						},
					},
				},
			},
			want: gal.NewNumberFromInt(5),
		},
		"function": {
			tree: gal.Tree{
				gal.NewNumberFromInt(10),
				gal.Plus,
				gal.NewFunction(
					"trunc",
					gal.Trunc,
					gal.Tree{
						gal.NewFunction(
							"sqrt",
							gal.Sqrt,
							gal.Tree{
								gal.NewNumberFromInt(10),
							},
						),
					},
					gal.Tree{
						gal.NewNumberFromInt(6),
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
				t.Error(cmp.Diff(tc.want, val))
			}
		})
	}
}
