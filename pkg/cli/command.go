package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Command struct {
	Name       string
	Usage      string
	Arguments  []string
	HandleFunc func(group *Group, command *Command, arguments []string) int
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

func (c *Command) PrintHelp(group *Group) int {
	fmt.Println("Usage:")

	usage := "  " + group.Name + " " + c.Name

	for _, argument := range c.Arguments {
		usage += " <" + argument + ">"
	}

	fmt.Println(usage)

	return Success
}

func (c *Command) PrintError(err error) int {
	log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return Failure
}
