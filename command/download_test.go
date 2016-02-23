package command

import (
	"os"
	"testing"

	"github.com/mitchellh/cli"
)

func newTestDownloadCommand() *DownloadCommand {
	meta := Meta{
		UI: &cli.ColoredUi{
			InfoColor:  cli.UiColorBlue,
			ErrorColor: cli.UiColorRed,
			Ui: &cli.BasicUi{
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
				Reader:      os.Stdin,
			},
		}}
	return &DownloadCommand{Meta: meta}
}

func TestDownloadCommand(t *testing.T) {
	var dl cli.Command = newTestDownloadCommand()
	args := []string{}
	res := dl.Run(args)
	if res != 1 {
		t.Fatal("should fail with missing args")
	}

	args = []string{
		"abc, foo",
	}
	res = dl.Run(args)
	if res != 1 {
		t.Fatal("should fail with invalid program ids")
	}
}
