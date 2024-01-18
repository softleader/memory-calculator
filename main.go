package main

import (
	"fmt"
	"github.com/softleader/memory-calculator/calc"
	"github.com/spf13/cobra"
	"os"
)

const (
	desc = `This command calculate the JVM memory for applications to run smoothly and stay within the memory limits of the container.
During the computation process, numerous parameters are required, which must be obtained in a specific order and logic.
The sequence and explanations of these parameters are as follows:

  1. Percentage of reserved space allocated by Memory Calculation tool:
     - First, determine if '--head-room' is passed through args.
     - If not, check the OS environment variable $BPL_JVM_HEAD_ROOM.
     - If neither is available, the default value is 0.

  2. Number of classes loaded at runtime:
     - First, determine if '--loaded-class-count' is passed through args.
     - If not, check the OS environment variable $BPL_JVM_LOADED_CLASS_COUNT.
     - If neither is available, dynamically calculate 35% of the total number of classes in the App directory.

  3. Number of user threads at runtime:
     - First, determine if '--thread-count' is passed through args.
     - If not, check the OS environment variable $BPL_JVM_THREAD_COUNT.
     - If neither is available, the default value is 200.

  4. App directory:
     - First, determine if '--app-path' is passed through args.
     - If not, the default directory is /app.

  5. Java startup parameters:
     - First, determine if '--jvm-options' is passed through args.
     - If not, check the OS environment variable $JAVA_OPTS.

  6. Java home:
     - Only check the OS environment variable $JAVA_HOME.

Examples:
  # Use ZGC and output to /tmp/.env
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
