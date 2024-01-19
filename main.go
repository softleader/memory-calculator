package main

import (
	"fmt"
	"github.com/softleader/memory-calculator/calc"
	"github.com/spf13/cobra"
	"os"
)

const (
	desc = `This command calculate the JVM memory for applications to run smoothly and stay within the memory limits of the container.

In the calculation process, most parameters have default values, with the loaded class count being the most crucial one:

  - First, it checks if '--loaded-class-count' has been passed as an argument.
  - If not, it will examine the environment variable $BPL_JVM_LOADED_CLASS_COUNT.
  - If neither option is available, it calculates the sum of the class counts in the App directory and the JVM as the loaded class count.

  Additionally, the App directory will first consider whether '--app-path' has been passed as an argument; if not, it defaults to using /app.
  The JVM class count will initially reference whether '--jvm-class-count' has been passed as an argument;
  if not, it will automatically calculate the number of classes under JAVA_HOME.

Examples:
  # Minimum example of input parameters
  memory-calculator --loaded-class-count 10000

  # Use ZGC and output to /tmp/.env and auto detect the loaded class count
  memory-calculator --jvm-options '-XX:+UseZGC' -o '/tmp/.env'

  # Print the version and exit
  memory-calculator --version
`
)

var version = "<unknown>"

type config struct {
	output  string
	version bool
	calc    calc.Calculator
}

func main() {
	c := config{
		version: false,
		output:  "",
		calc:    calc.NewCalculator(),
	}
	cmd := &cobra.Command{
		Use:          "memory-calculator",
		Short:        "JVM Memory Calculator",
		Long:         desc,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(c)
		},
	}
	f := cmd.Flags()
	f.Var(c.calc.JVMOptions, calc.FlagJVMOptions, calc.UsageJVMOptions)
	f.Var(c.calc.HeadRoom, calc.FlagHeadRoom, calc.UsageHeadRoom)
	f.Var(c.calc.ThreadCount, calc.FlagThreadCount, calc.UsageThreadCount)
	f.Var(c.calc.LoadedClassCount, calc.FlagLoadedClassCount, calc.UsageLoadedClassCount)
	f.Var(c.calc.JVMClassCount, calc.FlagJVMClassCount, calc.UsageJVMClassCount)
	f.Var(c.calc.JVMClassAdj, calc.FlagJVMClassAdj, calc.UsageJVMClassAdj)
	f.Var(c.calc.JVMCacerts, calc.FlagJVMCacerts, calc.UsageJVMCacerts)
	f.Var(c.calc.AppPath, calc.FlagAppPath, calc.UsageAppPath)
	f.Var(c.calc.EnableNmt, calc.FlagEnableNmt, calc.UsageEnableNmt)
	f.Var(c.calc.EnableJfr, calc.FlagEnableJfr, calc.UsageEnableJfr)
	f.Var(c.calc.EnableJmx, calc.FlagEnableJmx, calc.UsageEnableJmx)
	f.Var(c.calc.EnableJdwp, calc.FlagEnableJdwp, calc.UsageEnableJdwp)
	f.VarP(c.calc.Verbose, calc.FlagVerbose, calc.FlagShortVerbose, calc.UsageVerbose)
	f.StringVarP(&c.output, "output", "o", c.output, "write to a file, instead of STDOUT")
	f.BoolVar(&c.version, "version", c.version, "print version and exit")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(c config) error {
	if c.version {
		fmt.Println(version)
		return nil
	}
	j, err := c.calc.Execute()
	if err != nil {
		return err
	}
	return c.out(j)
}

func (c *config) out(j *calc.JavaToolOptions) error {
	if c.output == "" {
		j.Print()
		return nil
	}
	return j.WriteFile(c.output)
}
