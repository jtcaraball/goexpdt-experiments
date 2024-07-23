package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/jtcaraball/goexpdt/compute"
	"github.com/jtcaraball/goexpdt/query"
)

// randCompValDriver corresponds to the driver for experiments that use random
// positively classified instances to compute an optimal value based on the
// property and order generated by queryGF.
type randCompValDriver struct {
	queryGF closeOptimQueryGenFactory
}

// Run executes the experiment over the inputs passed in args and writes the
// results to out.
func (d randCompValDriver) Run(
	out io.Writer,
	solver string,
	args ...string,
) error {
	if len(args) < 2 {
		return errors.New("Missing arguments")
	}

	w := csv.NewWriter(out)

	m, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid multiplier '%s'", args[0])
	}

	if err = w.Write(
		[]string{
			"file_name",
			"tree_dim",
			"tree_nodes",
			"iter",
			"time",
			"value",
		},
	); err != nil {
		return err
	}

	for _, tp := range args[1:] {
		ctx, err := genContext(tp)
		if err != nil {
			return err
		}

		if err = d.eval(tp, solver, m, ctx, w); err != nil {
			return err
		}
	}

	return nil
}

// eval runs the experiment on a single input m amount of times and writes the
// output to w.
func (d randCompValDriver) eval(
	id, solver string,
	m int,
	ctx query.QContext,
	w *csv.Writer,
) error {
	v := query.QVar("x")
	dim := strconv.Itoa(ctx.Dim())
	nc := strconv.Itoa(len(ctx.Nodes()))

	for i := 0; i < m; i++ {
		fg, og, err := d.queryGF(ctx)
		if err != nil {
			return err
		}

		t := time.Now()

		out, err := compute.ComputeOptim(fg, og, v, ctx, solver)
		if err != nil {
			return fmt.Errorf("Compute error: %s", err.Error())
		}

		val := "-"
		if out.Found {
			val = out.Value.AsString()
		}
		ts := strconv.Itoa(int(time.Since(t)))

		if err = w.Write(
			[]string{id, dim, nc, strconv.Itoa(i), val, ts},
		); err != nil {
			return err
		}

		w.Flush() // Experiments are long. Save outputs often.
		ctx.Reset()
	}

	return nil
}

// randStatsDriver corresponds to the driver for experiments that use random
// positively classified instances to calculate stats for computing optimal
// values based on the property and order generated by queryGF.
type randStatsDriver struct {
	queryGF closeOptimQueryGenFactory
}

// Run executes the experiment over the inputs passed in args and writes the
// results to out.
func (d randStatsDriver) Run(
	out io.Writer,
	solver string,
	args ...string,
) error {
	if len(args) < 2 {
		return errors.New("Missing arguments")
	}

	w := csv.NewWriter(out)

	m, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid multiplier '%s'", args[0])
	}

	if err = w.Write(
		[]string{
			"file_name",
			"tree_dim",
			"tree_nodes",
			"iter",
			"#bots",
			"#calls",
			"time",
		},
	); err != nil {
		return err
	}

	for _, tp := range args[1:] {
		ctx, err := genContext(tp)
		if err != nil {
			return err
		}

		if err = d.eval(tp, solver, m, ctx, w); err != nil {
			return err
		}
	}

	return nil
}

// eval runs the experiment on a single input m amount of times and writes the
// output to w.
func (d randStatsDriver) eval(
	id, solver string,
	m int,
	ctx query.QContext,
	w *csv.Writer,
) error {
	v := query.QVar("x")
	dim := strconv.Itoa(ctx.Dim())
	nc := strconv.Itoa(len(ctx.Nodes()))

	for i := 0; i < m; i++ {
		fg, og, err := d.queryGF(ctx)
		if err != nil {
			return err
		}

		t := time.Now()

		out, err := compute.ComputeOptim(fg, og, v, ctx, solver)
		if err != nil {
			return fmt.Errorf("Compute error: %s", err.Error())
		}

		if err = w.Write(
			[]string{
				id,
				dim,
				nc,
				strconv.Itoa(i),
				strconv.Itoa(out.Value.BotCount()),
				strconv.Itoa(out.Calls),
				strconv.Itoa(int(time.Since(t))),
			},
		); err != nil {
			return err
		}

		w.Flush() // Experiments are long. Save outputs often.
		ctx.Reset()
	}

	return nil
}

// compValDriver corresponds to the driver for experiments that compute an
// optimal value based on the property and order generated by queryGF for a
// specific set of partial instances passed as input.
type compValDriver struct {
	queryGF openOptimQueryGenFactory
}

// Run executes the experiment over the inputs passed in args and writes the
// results to out.
func (d compValDriver) Run(out io.Writer, solver string, args ...string) error {
	if len(args) == 0 {
		return errors.New("Missing arguments")
	}

	w := csv.NewWriter(out)

	if err := w.Write(
		[]string{"file_name", "tree_dim", "tree_nodes", "time", "value"},
	); err != nil {
		return err
	}

	for _, tp := range args {
		if err := d.eval(tp, solver, w); err != nil {
			return err
		}
	}

	return nil
}

// eval runs the experiment on a single input  writes the outputs to w.
func (d compValDriver) eval(ip, solver string, w *csv.Writer) error {
	inst, ctx, err := parseTIInput(ip)
	if err != nil {
		return err
	}

	v := query.QVar("x")
	dim := strconv.Itoa(ctx.Dim())
	nc := strconv.Itoa(len(ctx.Nodes()))

	for _, c := range inst {
		fg, og, err := d.queryGF(ctx, c)
		if err != nil {
			return err
		}

		t := time.Now()

		out, err := compute.ComputeOptim(fg, og, v, ctx, solver)
		if err != nil {
			return fmt.Errorf("Compute error: %s", err.Error())
		}

		val := "-"
		if out.Found {
			val = out.Value.AsString()
		}

		if err = w.Write(
			[]string{ip, dim, nc, strconv.Itoa(int(time.Since(t))), val},
		); err != nil {
			return err
		}

		w.Flush() // Experiments are long. Save outputs often.
		ctx.Reset()
	}

	return nil
}
