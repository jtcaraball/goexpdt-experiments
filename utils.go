package main

import (
	"bufio"
	"errors"
	"fmt"
	"goexpdt-experiments/tree"
	"math/rand"
	"os"

	"github.com/jtcaraball/goexpdt/compute"
	"github.com/jtcaraball/goexpdt/query"
)

const (
	outputdir = "io/output"
	solver    = "./kissat"
)

// solveFormula and return ok, const value. ok is false if the formula is
// unsatisfiable.
func solveFormula(
	f compute.Encodable,
	v query.QVar,
	ctx query.QContext,
	solverPath, cnfPath string,
) (bool, query.QConst, error) {
	exitcode, out, err := compute.Step(
		f,
		ctx,
		solverPath,
		cnfPath,
	)
	if err != nil {
		return false, query.QConst{}, err
	}

	if exitcode == 10 {
		outConst, err := compute.GetValueFromBytes(out, v, ctx)
		if err != nil {
			return false, query.QConst{}, err
		}
		return true, outConst, nil
	}

	return false, query.QConst{}, nil
}

// bToC sets the values of c to constant represented in bytes b.
func bToC(b []byte, c query.QConst) error {
	if len(b) != len(c.Val) {
		return fmt.Errorf(
			"Invalid bytes length %d expected %d.",
			len(b),
			len(c.Val),
		)
	}
	for i, ch := range b {
		switch ch {
		case 48: // 0
			c.Val[i] = query.ZERO
		case 49: // 1
			c.Val[i] = query.ONE
		case 50: // 2
			c.Val[i] = query.BOT
		default:
			return fmt.Errorf("Invalid const feature value in index %d", i)
		}
	}
	return nil
}

// sToC sets the values of c to constant represented as s.
func sToC(s string, c query.QConst) error {
	if len(s) != len(c.Val) {
		return fmt.Errorf(
			"Invalid string length %d expected %d.",
			len(s),
			len(c.Val),
		)
	}
	for i, ch := range s {
		switch ch {
		case 48: // 0
			c.Val[i] = query.ZERO
		case 49: // 1
			c.Val[i] = query.ONE
		case 95: // _
			c.Val[i] = query.BOT
		default:
			return fmt.Errorf("Invalid const feature value in index %d", i)
		}
	}
	return nil
}

// randConst the values of c to a random partial instance. If full==true then
// no values will be set to query.BOT.
func randConst(c query.QConst, full bool) {
	limit := 3
	if full {
		limit = 2
	}

	for i := 0; i < len(c.Val); i++ {
		r := rand.Intn(limit)
		switch r {
		case 0:
			c.Val[i] = query.ZERO
		case 1:
			c.Val[i] = query.ONE
		case 2:
			c.Val[i] = query.BOT
		}
	}
}

// randValConst sets the value of c to a random partial instance with a
// classification equal to tVal.
func randValConst(c query.QConst, tVal bool, ctx query.QContext) error {
	match := false
	for !match {
		randConst(c, true)
		val, err := evalConst(c, ctx)
		if err != nil {
			return err
		}
		match = val == tVal
	}
	return nil
}

// evalConst runs the classification model over c and returns its
// classification. Returns a non nil error if the constant c is not full or
// the model in ctx is invalid.
func evalConst(c query.QConst, ctx query.QContext) (bool, error) {
	nodes := ctx.Nodes()
	node := nodes[0]

	for !node.IsLeaf() {
		if node.Feat < 0 || node.Feat >= len(c.Val) {
			return false, errors.New("Node feature out of index.")
		}
		switch c.Val[node.Feat] {
		case query.ONE:
			node = nodes[node.OChild]
		case query.ZERO:
			node = nodes[node.ZChild]
		default:
			return false, errors.New("Constant is not full")
		}
	}

	return node.Value, nil
}

// genContext returns a query.QContext based on the decision tree encoded in
// the file represented by its absolute path treePath.
func genContext(treePath string) (query.QContext, error) {
	t, err := tree.Load(treePath)
	if err != nil {
		return nil, err
	}
	ctx := query.BasicQContext(&t)
	return ctx, nil
}

// parseTIInput returns instances and context represented in the tree-instance
// input file passed by path.
func parseTIInput(inf string) ([]query.QConst, query.QContext, error) {
	treeFP, instStrings, err := scanTIFile(inf)
	if err != nil {
		return nil, nil, err
	}

	ctx, err := genContext(treeFP)
	if err != nil {
		return nil, nil, err
	}

	instances := make([]query.QConst, len(instStrings))
	for i, cb := range instStrings {
		instances[i] = query.AllBotConst(ctx.Dim())
		err := sToC(cb, instances[i])
		if err != nil {
			return nil, nil, err
		}
	}

	return instances, ctx, nil
}

// scanTIFIle scans a tree/instance input file by path. Returns its tree file
// path and a slice of instances represented as strings.
func scanTIFile(path string) (string, []string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return "", nil, errors.New("Empty input file.")
	}
	treeFP := scanner.Text()

	instStrings := []string{}
	for scanner.Scan() {
		instStrings = append(instStrings, scanner.Text())
	}
	if len(instStrings) == 0 {
		return "", nil, errors.New("No instances in input file.")
	}

	if scanner.Err() != nil {
		return "", nil, err
	}

	return treeFP, instStrings, nil
}
