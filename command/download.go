package command

import (
	"fmt"
	"log"
	"strings"

	"github.com/kosmikko/go-areena"
	"github.com/kosmikko/go-areena-dl/download"
)

// DownloadCommand command to download video
type DownloadCommand struct {
	Meta
}

// Run the command
func (c *DownloadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("download")
	flags.Usage = func() { c.Meta.UI.Output(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}
	parsedArgs := flags.Args()
	if len(parsedArgs) != 1 {
		msg := fmt.Sprintf("Invalid arguments: usage download PROGRAM_ID")
		c.UI.Error(msg)
		return 1
	}
	programID := parsedArgs[0]
	log.Printf("Downloading %s...", programID)

	cfg, err := areena.NewConfig()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	apiClient, err := areena.NewClient(cfg)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	dl := download.NewDownloader()
	pd, err := apiClient.ProgramDetails(programID)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	err = dl.DownloadVideo(pd.HLSURL, fmt.Sprintf("%s.mp4", pd.Slug))
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if pd.SubtitleURL != "" {
		err = dl.DownloadToFile(pd.SubtitleURL, fmt.Sprintf("%s.srt", pd.Slug))
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
	}

	return 0
}

func (c *DownloadCommand) Synopsis() string {
	return "Downloads video & subtitles"
}

func (c *DownloadCommand) Help() string {
	helpText := `
Usage:

goareena download [options] PROGRAM_ID

Options:

-output, -O  Change output file name. By default filename is generated from program name.

`
	return strings.TrimSpace(helpText)
}
