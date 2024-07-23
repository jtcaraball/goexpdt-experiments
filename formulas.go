package main

import (
	"github.com/jtcaraball/goexpdt/compute"
	"github.com/jtcaraball/goexpdt/query"
	"github.com/jtcaraball/goexpdt/query/extensions/allcomp"
	"github.com/jtcaraball/goexpdt/query/extensions/dfs"
	"github.com/jtcaraball/goexpdt/query/extensions/full"
	"github.com/jtcaraball/goexpdt/query/logop"
	"github.com/jtcaraball/goexpdt/query/predicates/subsumption"
)

// dfsFGF returns a query generator for the formula Determinant Feature Set.
func dfsFGF() compute.SVFormula {
	return func(v query.QVar) compute.Encodable {
		return dfs.Var{I: v}
	}
}

// srFGF returns a query generator for the formula Sufficient Reason of
// constant c.
func srFGF(c query.QConst) compute.SVFormula {
	return func(v query.QVar) compute.Encodable {
		return logop.WithVar{
			I: v,
			Q: logop.And{
				Q1: subsumption.VarConst{I1: v, I2: c},
				Q2: logop.And{
					Q1: logop.Or{
						Q1: logop.Not{
							Q: allcomp.Const{I: c, LeafValue: true},
						},
						Q2: allcomp.Var{
							I:               v,
							LeafValue:       true,
							ReachNodeVarGen: varGenNodeReach,
						},
					},
					Q2: logop.Or{
						Q1: logop.Not{
							Q: allcomp.Const{I: c, LeafValue: false},
						},
						Q2: allcomp.Var{
							I:               v,
							LeafValue:       false,
							ReachNodeVarGen: varGenNodeReach,
						},
					},
				},
			},
		}
	}
}

// crFGF returns a query generator for the formula Change Required of the
// constant c.
func crFGF(c query.QConst) compute.SVFormula {
	return func(v query.QVar) compute.Encodable {
		return logop.WithVar{
			I: v,
			Q: logop.And{
				Q1: full.Var{I: v},
				Q2: logop.And{
					Q1: full.Const{I: c},
					Q2: logop.Or{
						Q1: logop.And{
							Q1: allcomp.Var{
								I:               v,
								LeafValue:       true,
								ReachNodeVarGen: varGenNodeReach,
							},
							Q2: logop.Not{
								Q: allcomp.Const{I: c, LeafValue: true},
							},
						},
						Q2: logop.And{
							Q1: allcomp.Const{I: c, LeafValue: true},
							Q2: logop.Not{
								Q: allcomp.Var{
									I:               v,
									LeafValue:       true,
									ReachNodeVarGen: varGenNodeReach,
								},
							},
						},
					},
				},
			},
		}
	}
}

// caFGF returns a query generator for the formula Change Allowed of the
// constant c.
func caFGF(c query.QConst) compute.SVFormula {
	return func(v query.QVar) compute.Encodable {
		return logop.WithVar{
			I: v,
			Q: logop.And{
				Q1: full.Var{I: v},
				Q2: logop.And{
					Q1: full.Const{I: c},
					Q2: logop.And{
						Q1: logop.Or{
							Q1: logop.Not{
								Q: allcomp.Var{
									I:               v,
									LeafValue:       true,
									ReachNodeVarGen: varGenNodeReach,
								},
							},
							Q2: allcomp.Const{I: c, LeafValue: true},
						},
						Q2: logop.Or{
							Q1: logop.Not{
								Q: allcomp.Const{I: c, LeafValue: true},
							},
							Q2: allcomp.Var{
								I:               v,
								LeafValue:       true,
								ReachNodeVarGen: varGenNodeReach,
							},
						},
					},
				},
			},
		}
	}
}
