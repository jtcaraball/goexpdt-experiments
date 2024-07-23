package main

import (
	"io"
	"os"
	"path"
	"strings"
	"time"
)

// driver for running the optimization algorithm over a set of inputs.
type driver interface {
	Run(out io.Writer, solver string, args ...string) error
}

// experiment corresponds to a particular instance of a query, determined by
// the driver d, that can be executed over a set of inputs. It additionally
// contains information about the experiment.
type experiment struct {
	Name        string
	Description string
	d           driver
}

// Run the experiment over the set of inputs and using the options contained in
// args.
func (e experiment) Run(args ...string) error {
	ts := strings.Join(strings.Split(time.Now().String(), " ")[:2], "_")
	outp := path.Join(outputdir, e.Name+"_"+ts+".csv")

	of, err := os.Create(outp)
	if err != nil {
		return err
	}
	defer of.Close()

	if err = e.d.Run(of, solver, args...); err != nil {
		os.Remove(of.Name())
		return err
	}

	return nil
}

var experiments = []experiment{
	{
		"optim:rand:stats:dfs-ll",
		"Optimum (Stats, Random Instances) - DFS under Lesser Level Order.\n" +
			"Arguments:\n" +
			"  - n (instances per input\n" +
			"  - List of <tree_file_inputs>",
		randStatsDriver{DFS_LL_C},
	},
	{
		"optim:rand:stats:sr-ll",
		"Optimum (Stats, Random Instances) - SR under Lesser Level Order.\n" +
			"Arguments:\n" +
			"  - n (instances per input\n" +
			"  - List of <tree_file_inputs>",
		randStatsDriver{SR_LL_C},
	},
	{
		"optim:rand:stats:sr-ss",
		"Optimum (Stats, Random Instances) - SR under Strict Subsumption" +
			" Order.\n" +
			"Arguments:\n" +
			"  - n (instances per input\n" +
			"  - List of <tree_file_inputs>",
		randStatsDriver{SR_SS_C},
	},
	{
		"optim:rand:stats:cr-lh",
		"Optimum (Stats, Random Instances) - CR under Lesser Hamming" +
			" Distance Order.\n" +
			"Arguments:\n" +
			"  - n (instances per input\n" +
			"  - List of <tree_file_inputs>",
		randStatsDriver{CR_LH_C},
	},
	{
		"optim:rand:stats:ca-gh",
		"Optimum (Stats, Random Instances) - CA under Greater Hamming" +
			" Distance Order.\n" +
			"Arguments:\n" +
			"  - n (instances per input\n" +
			"  - List of <tree_file_inputs>",
		randStatsDriver{CA_GH_C},
	},
	{
		"optim:rand:val:dfs-ll",
		"Optimum (Value, Random Instances) - DFS under Lesser Level Order.\n" +
			"Arguments:\n" +
			"  - n (instances per input\n" +
			"  - List of <tree_file_inputs>",
		randCompValDriver{DFS_LL_C},
	},
	{
		"optim:val:dfs-ll",
		"Optimum (Value) - DFS under Lesser Level Order.\n" +
			"Arguments:\n" +
			"  - n (instances per input\n" +
			"  - List of <tree_file_inputs>",
		compValDriver{DFS_LL_O},
	},
	{
		"optim:val:sr-ll",
		"Optimum (Value) - SR under Lesser Level Order.\n" +
			"Arguments:\n" +
			"  - List of <optim_file_input>",
		compValDriver{SR_LL_O},
	},
	{
		"optim:val:sr-ss",
		"Optimum (Value) - SR under Strict Subsumption Order.\n" +
			"Arguments:\n" +
			"  - List of <optim_file_input>",
		compValDriver{SR_SS_O},
	},
	{
		"optim:val:cr-lh",
		"Optimum (Value) - CR under Less Hamming Distance Order.\n" +
			"Arguments:\n" +
			"  - List of <optim_file_input>",
		compValDriver{CR_LH_O},
	},
	{
		"optim:val:ca-gh",
		"Optimum (Value) - CA under Greater Hamming Distance Order.\n" +
			"Arguments:\n" +
			"  - List of <optim_file_input>",
		compValDriver{CA_GH_O},
	},
}

// expMap returns map of implemented experiments with their name as key.
func expMap() map[string]experiment {
	exps := make(map[string]experiment)
	for _, exp := range experiments {
		exps[exp.Name] = exp
	}
	return exps
}
