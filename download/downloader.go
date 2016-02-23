package download

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"

	"github.com/grafov/m3u8"
)

var log = getLogger()

// Downloader downloads HLS streams
type Downloader struct {
	client *http.Client
	output *os.File
}

// NewDownloader constructor
func NewDownloader() *Downloader {
	// need cookies for Akamai acl
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	return &Downloader{
		client: client,
	}
}

// DownloadVideo downloads HLS video to file
func (d *Downloader) DownloadVideo(url string, outputFilename string) (err error) {
	log.Infof("Downloading video %s to file %s", url, outputFilename)
	out, err := os.Create(outputFilename)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer out.Close()
	d.output = out
	err = d.fetchPlaylist(url)
	if err != nil {
		log.Infof("Video downloaded successfully")
	}
	return
}

// DownloadToFile download e.g. subtitles to file
func (d *Downloader) DownloadToFile(url string, outputFilename string) (err error) {
	log.Infof("Downloading %s to file %s", url, outputFilename)
	out, err := os.Create(outputFilename)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer out.Close()
	resp, err := d.request(url)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to fetch %s, status %v", url, resp.StatusCode)
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()
	return
}

// FetchPlaylist GET playlist content & fetch files
func (d *Downloader) fetchPlaylist(purl string) (err error) {
	purl = stripInvalidChars(purl)
	if err != nil {
		return
	}

	resp, err := d.request(purl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body := bufio.NewReader(resp.Body)
	playlist, listType, err := m3u8.DecodeFrom(body, true)
	if err != nil {
		return
	}

	switch listType {
	case m3u8.MEDIA:
		mediapl := playlist.(*m3u8.MediaPlaylist)
		log.Infof("downloading playlist, %d segments",
			mediapl.Count(),
		)
		d.downloadSegments(mediapl)

	case m3u8.MASTER:
		masterpl := playlist.(*m3u8.MasterPlaylist)

		err = d.fetchPlaylist(urlForBestQualityVariant(masterpl))
		if err != nil {
			return err
		}
	default:
		log.Error("playlist type not supported")
	}
	return
}

func (d *Downloader) downloadSegments(mpl *m3u8.MediaPlaylist) error {
	for _, segment := range mpl.Segments {
		if segment != nil && segment.URI != "" {
			log.Debugf("downloading segment [%s]",
				segment.URI,
			)
			d.downloadSegmentToFile(segment.URI)
		}
	}
	return nil
}

func (d *Downloader) downloadSegmentToFile(uri string) {
	resp, err := d.request(uri)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if resp.StatusCode != 200 {
		log.Errorf("Received HTTP %d for %s", resp.StatusCode, uri)
		return
	}
	_, err = io.Copy(d.output, resp.Body)
	if err != nil {
		log.Errorf("Failed to write file: %s", err.Error())
		return
	}
	resp.Body.Close()
}

func (d *Downloader) request(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	resp, err = d.client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode >= 400 {
		contents, parseErr := ioutil.ReadAll(resp.Body)
		if parseErr != nil {
			return nil, parseErr
		}
		err = errors.New(string(contents))
		resp.Body.Close()
	}
	return
}

func stripInvalidChars(str string) string {
	b := make([]byte, len(str))
	var bl int
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 32 && c < 127 {
			b[bl] = c
			bl++
		}
	}
	return string(b[:bl])
}

func urlForBestQualityVariant(masterpl *m3u8.MasterPlaylist) (url string) {
	var bestVariant *m3u8.Variant
	for _, variant := range masterpl.Variants {
		log.Debugf("available variant bandwith: %d, codec: %s, alternatives: %+v", variant.Bandwidth, variant.Codecs, variant.Alternatives)
		if bestVariant == nil || variant.VariantParams.Bandwidth > bestVariant.VariantParams.Bandwidth {
			bestVariant = variant
		}
	}
	log.Infof("selected variant %d", bestVariant.Bandwidth)
	return bestVariant.URI
}
