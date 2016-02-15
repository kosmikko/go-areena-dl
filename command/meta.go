package command

import (
	"bufio"
	"flag"
	"io"

	"github.com/mitchellh/cli"
)

// Meta contain the meta-option that nearly all subcommand inherited.
type Meta struct {
	UI cli.Ui
}

func (m *Meta) FlagSet(n string) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)

	errR, errW := io.Pipe()
	errScanner := bufio.NewScanner(errR)
	go func() {
		for errScanner.Scan() {
			m.UI.Error(errScanner.Text())
		}
	}()
	f.SetOutput(errW)

	return f
}
