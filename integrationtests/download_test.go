package integrationtests

import (
	"fmt"
	"testing"

	"github.com/kosmikko/go-areena"

	"github.com/kosmikko/go-areena-dl/download"
)

func Test(t *testing.T) {
	cfg, err := areena.NewConfig()
	if err != nil {
		t.Fatal(err)
	}

	apiClient, err := areena.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	pd, err := apiClient.ProgramDetails("1-2371565")
	if err != nil {
		t.Fatal(err)
	}

	dl := download.NewDownloader()
	err = dl.DownloadVideo(pd.HLSURL, fmt.Sprintf("%s.mp4", pd.Slug))
	if err != nil {
		t.Fatal(err)
	}

	err = dl.DownloadToFile(pd.SubtitleURL, fmt.Sprintf("%s.srt", pd.Slug))
	if err != nil {
		t.Fatal(err)
	}
}
