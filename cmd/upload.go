package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"

	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/golog"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload [filter]",
	Args:  cobra.MinimumNArgs(1),
	Short: "(re)Upload a back catalog of replays specified",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not yet implemented")
	},
}

func uploadReplay(apikey, filename string) error {
	var buf bytes.Buffer
	mpw := multipart.NewWriter(&buf)
	w, err := mpw.CreateFormFile("replay_file", filename)
	if err != nil {
		return err
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = io.Copy(w, f); err != nil {
		return err
	}
	if err := mpw.WriteField("upload_method", fmt.Sprintf("sc2-rsu-%s", runtime.GOOS)); err != nil {
		return err
	}
	if err := mpw.Close(); err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/replay", sc2rsAPIBase), &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", apikey)
	req.Header.Set("Content-Type", mpw.FormDataContentType())
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sc2replaystats returned error: %v", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	response := make(map[string]string)
	if err := jsoniter.Unmarshal(b, &response); err != nil {
		return err
	}
	if id, ok := response["replay_queue_id"]; ok {
		golog.Infof("sc2replaystats accepted our replay, queue ID: %v", id)
	}

	return nil
}
