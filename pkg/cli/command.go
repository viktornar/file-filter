package cli

import (
	"file-filter/pkg/logger"
	"flag"
	"fmt"
)

type Command struct {
	Name       string
	Usage      string
	Arguments  []string
	HandleFunc func(name string, command *Command, arguments []string) int
	flagSet    *flag.FlagSet
}

func (c *Command) FlagSet() *flag.FlagSet {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet(c.Name, flag.ExitOnError)
	}

	return c.flagSet
}

func (c *Command) Parse(arguments []string) ([]string, error) {
	err := c.FlagSet().Parse(arguments)
	if err != nil {
		return nil, err
	}

	if len(c.FlagSet().Args()) < len(c.Arguments) {
		return arguments, ErrMissingArguments
	}

	return c.FlagSet().Args(), nil
}

func (c *Command) PrintHelp() int {
	fmt.Println("Usage:")

	usage := "  " + c.Name

	for _, argument := range c.Arguments {
		usage += " <" + argument + ">"
	}

	fmt.Println(usage)

	return Success
}

func (c *Command) PrintError(err error) int {
	logger.Error.Println(err)

	return Failure
}
