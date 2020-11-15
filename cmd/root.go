package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kataras/golog"
	"github.com/mitchellh/go-wordwrap"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/AlbinoGeek/sc2-rsu/sc2replaystats"
	"github.com/AlbinoGeek/sc2-rsu/utils"
)

var (
	rootCmd = &cobra.Command{
		Short: "SC2ReplayStats Uploader",
		Long:  `Unofficial SC2ReplayStats Uploader by AlbinoGeek`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetBool("update.check.enabled") {
				go checkUpdateEvery(getUpdateDuration())
			}

			key := viper.GetString("apikey")
			if key == "" {
				return errors.New("no API key in configuration, please use the login command")
			}

			if !sc2replaystats.ValidAPIKey(key) {
				return errors.New("invalid API key in configuration, please replace it or use the login command")
			}

			paths, err := getWatchPaths()
			if err != nil {
				return err
			}

			golog.Info("Starting Automatic Replay Uploader...")
			sc2api = sc2replaystats.New(key)

			done := make(chan struct{})
			w, err := automaticUpload(paths)
			if err != nil {
				return err
			}
			defer w.Close()

			// Setup Interrupt (Ctrl+C) Handler
			c := make(chan os.Signal)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			go func() {
				sig := <-c
				fmt.Println()
				golog.Warnf("Received signal:%v, Quitting.", sig)
				close(done)
			}()

			golog.Debugf("Startup took: %v", time.Since(startTime))
			golog.Info("Ready!")
			<-done
			return nil
		},
	}
	sc2api    *sc2replaystats.Client
	startTime = time.Now()
	termWidth = 80
)

func automaticUpload(paths []string) (w *fsnotify.Watcher, err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to setup fswatcher: %v", err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					for {
						time.Sleep(time.Millisecond * 100)
						if s, err := os.Stat(event.Name); err == nil && s.Size() > 256 {
							break
						}
					}
					go handleReplay(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				golog.Warnf("fswatcher error: %v", err)
			}
		}
	}()

	for _, p := range paths {
		golog.Debugf("Watching replays directory: %v", p)
		if err = watcher.Add(p); err != nil {
			golog.Fatalf("failed to watch replay directory: %v: %v", p, err)
		}
	}

	return watcher, nil
}

func findAccounts(root string) (ids []string, err error) {
	golog.Debug("Searching for accounts in replay directory...")

	paths, err := utils.FindDirectoriesBySuffix(root, "ultiplayer", true)
	if err != nil {
		return nil, fmt.Errorf("FindDirectory error: %v", err)
	}

	i := 0
	uniq := make(map[string]struct{})
	for _, p := range paths {
		p = utils.StripPathParts(p, 2)
		if _, duplicate := uniq[p]; !duplicate {
			uniq[p] = struct{}{}
			paths[i] = p[len(root):]
			golog.Debugf("found candidate account: %v", paths[i])
			i++
		}
	}
	paths = paths[:i] // truncate duplicates
	golog.Debugf("finished scanning for candidates. Found: %v", len(paths))

	return paths, err
}

func findReplaysRoot() (root string, err error) {
	golog.Info("Determining replays directory... (this could take a few minutes)...")

	scanRoot := "/"
	if home, err := os.UserHomeDir(); err == nil {
		scanRoot = home
	}

	paths, err := utils.FindDirectoriesBySuffix(scanRoot, "ultiplayer", true)
	if err != nil {
		return "", fmt.Errorf("FindDirectory error: %v", err)
	}

	i := 0
	uniq := make(map[string]struct{})
	for _, p := range paths {
		p = utils.StripPathParts(p, 4)
		if _, duplicate := uniq[p]; !duplicate && p != "/" {
			uniq[p] = struct{}{}
			paths[i] = p
			golog.Debugf("found candidate replay directory: %v", p)
			i++
		}
	}
	paths = paths[:i] // truncate duplicates
	golog.Debugf("finished scanning for candidates. Found: %v", len(paths))

	if len(paths) == 1 {
		return paths[0], nil
	}

	if len(paths) == 0 {
		return "", fmt.Errorf("no root found")
	}

	line := strings.Repeat("=", termWidth/2)
	fmt.Printf("\n%s\n%s\n", line, wordwrap.WrapString("More than one possible replay directory was located while we scanned for your StarCraft II installation's Accounts folder.\n\nPlease select which directory we should be watching below:", uint(termWidth/2)))
	for i, p := range paths {
		fmt.Printf("\n  %d: %s", 1+i, p)
	}
	fmt.Printf("\n%s\n", line)
	consoleReader := bufio.NewReaderSize(os.Stdin, 1)
	for {
		fmt.Printf("Your Choice [1-%d]: ", len(paths))
		input, _, _ := consoleReader.ReadLine()
		choice, err := strconv.Atoi(string(input))
		if err == nil && choice-1 > 0 && choice-1 < len(paths) {
			return paths[choice-1], nil
		}
	}
}

func getWatchPaths() ([]string, error) {
	replaysRoot := viper.GetString("replaysRoot")
	if f, err := os.Stat(replaysRoot); err != nil || !f.IsDir() {
		golog.Warn("Replay Root not configured correctly, searching for replays directory...")
		if replaysRoot, err = findReplaysRoot(); err != nil {
			golog.Fatalf("unable to automatically determine the path to your replays directory: %v", err)
		}

		viper.Set("replaysRoot", replaysRoot)
		if err := saveConfig(); err != nil {
			return nil, err
		}
		golog.Infof("Using replays directory: %v", replaysRoot)
	}

	accs, err := findAccounts(replaysRoot)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0)
	for _, a := range accs {
		p := filepath.Join(replaysRoot, a, "Replays", "Multiplayer")
		if f, err := os.Stat(p); err == nil && f.IsDir() {
			paths = append(paths, p)
		}
	}

	return paths, nil
}

func handleReplay(replayFilename string) {
	golog.Debugf("uploading replay: %v", replayFilename)
	_, mapName, _ := utils.SplitFilepath(replayFilename)

	rqid, err := sc2api.UploadReplay(replayFilename)
	if err != nil {
		golog.Errorf("failed to upload replay: %v: %v", mapName, err)
		return
	}
	golog.Infof("sc2replaystats accepted : [%v] %s", rqid, mapName)
	go watchReplayStatus(rqid)
}

func watchReplayStatus(rqid string) {
	for {
		time.Sleep(time.Second)
		rid, err := sc2api.GetReplayStatus(rqid)
		if err != nil {
			golog.Errorf("error checking reply status: %v: %v", rqid, err)
			return // could not check status
		}

		if rid != "" {
			golog.Infof("sc2replaystats processed: [%v] %s", rqid, rid)
			return // replay parsed!
		}

		golog.Debugf("sc2replaystats process..: [%v] %s", rqid, rid)
	}
}
