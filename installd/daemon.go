package installd

import (
	"fmt"
	"time"
	"github.com/cavaliercoder/grab"
	"github.com/toasterson/mozaik/logger"
)

func HTTPDownload(url string, location string) (err error) {
	if location == ""{
		location = "/tmp/"
	}
	client := grab.NewClient()
	req, _ := grab.NewRequest(location, url)

	// start download
	logger.Info(fmt.Sprintf("Downloading %v...\n", req.URL()))
	resp := client.Do(req)

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

ProgressLoop:
	for {
		select {
		case <-t.C:
			logger.Info(fmt.Sprintf("  transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress()))

		case <-resp.Done:
			// download is complete
			break ProgressLoop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		logger.Critical(fmt.Sprintf("Download failed: %v\n", err))
		return
	}

	logger.Info(fmt.Sprintf("Download saved to %v \n", resp.Filename))
	return
}