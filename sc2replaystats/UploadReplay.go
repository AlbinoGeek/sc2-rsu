package sc2replaystats

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func prepareMultipartUpload(filename string) (buf bytes.Buffer, contentType string, err error) {
	mpw := multipart.NewWriter(&buf)
	defer mpw.Close()

	w, err := mpw.CreateFormFile("replay_file", filename)
	if err != nil {
		return buf, "", fmt.Errorf("create form file: %v", err)
	}

	f, err := os.Open(filename)
	if err != nil {
		return buf, "", fmt.Errorf("open form file: %v", err)
	}
	defer f.Close()

	if _, err = io.Copy(w, f); err != nil {
		return
	}

	if err = mpw.WriteField("upload_method", ClientIdentifier); err != nil {
		return
	}

	return buf, mpw.FormDataContentType(), nil
}

// UploadReplay sends the specified replay to sc2replaystats queue for processing
func (client *Client) UploadReplay(filename string) (replayQueueID string, err error) {
	buf, contentType, err := prepareMultipartUpload(filename)

	if err != nil {
		return "", fmt.Errorf("failed to prepare formdata: %v", err)
	}

	res, err := client.requestMap(http.MethodPost, "replay", contentType, &buf)

	// return the replay_queue_id if present
	if rqid, ok := res["replay_queue_id"]; ok {
		replayQueueID = rqid
	}

	return
}
