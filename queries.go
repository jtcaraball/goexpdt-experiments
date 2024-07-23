package main

import (
	"errors"

	"github.com/jtcaraball/goexpdt/compute"
	"github.com/jtcaraball/goexpdt/query"
)

type (
	// closeOptimQueryGenFactory returns a property and strict order generator
	// based on the query.QContext passed.
	closeOptimQueryGenFactory func(ctx query.QContext) (
		compute.SVFormula,
		compute.VCOrder,
		error,
	)
	// closeOptimQueryGenFactory returns a property and strict order generator
	// based on the query.QContext and constants cs passed.
	openOptimQueryGenFactory func(ctx query.QContext, cs ...query.QConst) (
		compute.SVFormula,
		compute.VCOrder,
		error,
	)
)

// =========================== //
//        CLOSE QUERIES        //
// =========================== //

func SR_LL_C(ctx query.QContext) (
	compute.SVFormula,
	compute.VCOrder,
	error,
) {
	c := query.AllBotConst(ctx.Dim())
	err := randValConst(c, true, ctx)
	if err != nil {
		return nil, nil, err
	}
	return srFGF(c), llOGF(), nil
}

func SR_SS_C(ctx query.QContext) (
	compute.SVFormula,
	compute.VCOrder,
	error,
) {
	c := query.AllBotConst(ctx.Dim())
	err := randValConst(c, true, ctx)
	if err != nil {
		return nil, nil, err
	}
	return srFGF(c), ssOGF(), nil
}

func DFS_LL_C(ctx query.QContext) (
	compute.SVFormula,
	compute.VCOrder,
	error,
) {
	return dfsFGF(), llOGF(), nil
}

func CR_LH_C(ctx query.QContext) (
	compute.SVFormula,
	compute.VCOrder,
	error,
) {
	c := query.AllBotConst(ctx.Dim())
	err := randValConst(c, true, ctx)
	if err != nil {
		return nil, nil, err
	}
	return crFGF(c), lhOGF(c), nil
}

func CA_GH_C(ctx query.QContext) (
	compute.SVFormula,
	compute.VCOrder,
	error,
) {
	c := query.AllBotConst(ctx.Dim())
	err := randValConst(c, true, ctx)
	if err != nil {
		return nil, nil, err
	}
	return caFGF(c), ghOGF(c), nil
}

// =========================== //
//         OPEN QUERIES        //
// =========================== //

func DFS_LL_O(ctx query.QContext, cs ...query.QConst) (
	compute.SVFormula,
	compute.VCOrder,
	error,
) {
	return dfsFGF(), llOGF(), nil
}

func CR_LH_O(ctx query.QContext, cs ...query.QConst) (
	compute.SVFormula,
	compute.VCOrder,
	error,
) {
	if len(cs) == 0 {
		return nil, nil, errors.New("Missing constant in query factory.")
	}
	return crFGF(cs[0]), lhOGF(cs[0]), nil
}

func CA_GH_O(ctx query.QContext, cs ...query.QConst) (
	compute.SVFormula,
	compute.VCOrder,
	error,
) {
	if len(cs) == 0 {
		return nil, nil, errors.New("Missing constant in query factory.")
	}
	return caFGF(cs[0]), ghOGF(cs[0]), nil
}

func SR_LL_O(ctx query.QContext, cs ...query.QConst) (
	compute.SVFormula,
	compute.VCOrder,
	error,
) {
	return srFGF(cs[0]), llOGF(), nil
}

func SR_SS_O(ctx query.QContext, cs ...query.QConst) (
	compute.SVFormula,
	compute.VCOrder,
	error,
) {
	return srFGF(cs[0]), ssOGF(), nil
}
