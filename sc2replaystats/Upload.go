package sc2replaystats

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	jsoniter "github.com/json-iterator/go"
)

func prepareMultipartUpload(filename string) (buf bytes.Buffer, contentType string, err error) {
	mpw := multipart.NewWriter(&buf)
	defer mpw.Close()

	w, err := mpw.CreateFormFile("replay_file", filename)
	if err != nil {
		return buf, "", err
	}

	f, err := os.Open(filename)
	if err != nil {
		return buf, "", err
	}
	defer f.Close()

	if _, err = io.Copy(w, f); err != nil {
		return
	}

	if err = mpw.WriteField("upload_method", ClientIdentifier); err != nil {
		return
	}

	if err = mpw.Close(); err != nil {
		return
	}

	return buf, mpw.FormDataContentType(), nil
}

// UploadReplay sends the specified replay to sc2replaystats queue for processing
func UploadReplay(apikey, filename string) (string, error) {
	buf, contentType, err := prepareMultipartUpload(filename)
	if err != nil {
		return "", fmt.Errorf("failed to prepare formdata: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/replay", APIRoot), &buf)
	if err != nil {
		return "", fmt.Errorf("failed to prepare request: %v", err)
	}

	req.Header.Set("Authorization", apikey)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("sc2replaystats API returned error: %v", resp.Status)
	}

	response := make(map[string]string)
	if err = jsoniter.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	// return the replay_queue_id if present
	if id, ok := response["replay_queue_id"]; ok {
		return id, nil
	}

	return "-1", nil
}
