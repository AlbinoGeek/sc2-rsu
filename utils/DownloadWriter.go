package utils

import (
	"fmt"
	"io"
	"net/http"
)

// DownloadWriter retrieves a file from a given URL, and writes the returned
// contents into the given writer, making sure a status 200 OK was returned.
func DownloadWriter(URL string, w io.Writer) error {
	resp, err := http.DefaultClient.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status != 200")
	}

	_, err = io.Copy(w, resp.Body)

	return err
}
