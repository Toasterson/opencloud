package installd

import (
	"fmt"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/toasterson/mozaik/logger"
)

func HTTPDownload(url string, location string) (err error) {
	if location == "" {
		location = "/tmp/"
	}
	req, _ := grab.NewRequest(location, url)
	return doDownload(req).Err()
}

func HTTPDownloadTo(url string, location string) (file string, err error){
	if location == "" {
		location = "/tmp/"
	}
	req, _ := grab.NewRequest(location, url)
	resp := doDownload(req)
	return resp.Filename, resp.Err()
}

func doDownload(request *grab.Request) (resp *grab.Response){
	client := grab.NewClient()
	// start download
	logger.Info(fmt.Sprintf("Downloading %v...", request.URL()))
	resp = client.Do(request)

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

ProgressLoop:
	for {
		select {
		case <-t.C:
			logger.Info(fmt.Sprintf("  transferred %v / %v bytes (%.2f%%)",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress()))

		case <-resp.Done:
			// download is complete
			break ProgressLoop
		}
	}

	logger.Info(fmt.Sprintf("Download saved to %v", resp.Filename))
	return
}
