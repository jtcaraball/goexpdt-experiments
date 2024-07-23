package main

import (
	"github.com/jtcaraball/goexpdt/compute"
	"github.com/jtcaraball/goexpdt/query"
	"github.com/jtcaraball/goexpdt/query/extensions/leh"
	"github.com/jtcaraball/goexpdt/query/logop"
	"github.com/jtcaraball/goexpdt/query/predicates/lel"
	"github.com/jtcaraball/goexpdt/query/predicates/subsumption"
)

// llOGF returns a query generator for the strict partial order Less Level.
func llOGF() compute.VCOrder {
	return func(v query.QVar, c query.QConst) compute.Encodable {
		return logop.And{
			Q1: lel.VarConst{
				I1:          v,
				I2:          c,
				CountVarGen: varGenBotCount,
			},
			Q2: logop.Not{
				Q: lel.ConstVar{
					I1:          c,
					I2:          v,
					CountVarGen: varGenBotCount,
				},
			},
		}
	}
}

// ssOGF returns a query generator for the strict partial order Strict
// Subsumption.
func ssOGF() compute.VCOrder {
	return func(v query.QVar, c query.QConst) compute.Encodable {
		return logop.And{
			Q1: subsumption.VarConst{I1: v, I2: c},
			Q2: logop.Not{Q: subsumption.ConstVar{I1: c, I2: v}},
		}
	}
}

// lhOGF returns a query generator for the strict partial order Less Hamming
// Distance.
func lhOGF(cp query.QConst) compute.VCOrder {
	return func(v query.QVar, c query.QConst) compute.Encodable {
		return logop.And{
			Q1: leh.ConstVarConst{
				I1:                    cp,
				I2:                    v,
				I3:                    c,
				HammingDistanceVarGen: varGenHammingDistance,
			},
			Q2: logop.Not{
				Q: leh.ConstConstVar{
					I1:                    cp,
					I2:                    c,
					I3:                    v,
					HammingDistanceVarGen: varGenHammingDistance,
				},
			},
		}
	}
}

// ghOGF returns a query generator for the strict partial order Greater Hamming
// Distance.
func ghOGF(cp query.QConst) compute.VCOrder {
	return func(v query.QVar, c query.QConst) compute.Encodable {
		return logop.And{
			Q1: leh.ConstConstVar{
				I1:                    cp,
				I2:                    c,
				I3:                    v,
				HammingDistanceVarGen: varGenHammingDistance,
			},
			Q2: logop.Not{
				Q: leh.ConstVarConst{
					I1:                    cp,
					I2:                    v,
					I3:                    c,
					HammingDistanceVarGen: varGenHammingDistance,
				},
			},
		}
	}
}
