package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kataras/golog"
	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/crypto/ssh/terminal"

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
			key := viper.GetString("apikey")
			if key == "" {
				return errors.New("no API key in configuration, please use the login command")
			}

			if !sc2replaystats.ValidAPIKey(key) {
				return errors.New("invalid API key in configuration, please replace it or use the login command")
			}

			golog.Info("Starting Automatic Replay Uploader...")
			return automaticUpload(key)
		},
	}

	termWidth = 80
)

func automaticUpload(apikey string) error {
	tt := time.Now()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		golog.Fatalf("failed to setup fswatcher: %v", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					golog.Debugf("uploading replay: %v", event.Name)
					if id, err := sc2replaystats.UploadReplay(apikey, event.Name); err != nil {
						golog.Errorf("failed to upload replay: %v: %v", event.Name, err)
					} else {
						golog.Infof("sc2replaystats accepted our replay: [%v] %s", id, filepath.Base(event.Name))
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				golog.Warnf("fswatcher error: %v", err)
			}
		}
	}()

	watchPaths := make([]string, 0)
	replaysRoot := viper.GetString("replaysRoot")
	if f, err := os.Stat(replaysRoot); err != nil || !f.IsDir() {
		golog.Warn("Replay Root not configured correctly, searching for replays directory...")
		if replaysRoot, err = findReplaysRoot(); err != nil {
			golog.Fatalf("unable to automatically determine the path to your replays directory: %v", err)
		}

		viper.Set("replaysRoot", replaysRoot)
		if err := saveConfig(); err != nil {
			return err
		}
		golog.Infof("Using replays directory: %v", replaysRoot)
	}

	accs, err := findAccounts(replaysRoot)
	if err != nil {
		return err
	}

	for _, a := range accs {
		path := filepath.Join(replaysRoot, a, "Replays", "Multiplayer")
		if f, err := os.Stat(path); err == nil && f.IsDir() {
			watchPaths = append(watchPaths, path)
		}
	}

	for _, p := range watchPaths {
		golog.Debugf("Watching replays directory: %v", p)
		if err = watcher.Add(p); err != nil {
			golog.Fatalf("failed to watch replay directory: %v: %v", p, err)
		}
	}

	// Setup Interrupt (Ctrl+C) Handler
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-c
		fmt.Println()
		golog.Warnf("Received signal:%v, Quitting.", sig)
		close(done)
	}()

	golog.Debugf("Startup took: %v", time.Since(tt))
	golog.Info("Ready!")
	<-done
	return nil
}

func findAccounts(root string) (ids []string, err error) {
	golog.Debug("Searching for accounts in replay directory...")

	paths, err := utils.FindDirectoriesBySuffix(viper.GetString("replaysRoot"), "ultiplayer", true)
	if err != nil {
		return nil, fmt.Errorf("FindDirectory error: %v", err)
	}

	i := 0
	uniq := make(map[string]struct{})
	for _, p := range paths {
		p = utils.StripPathParts(p, 2)
		if _, duplicate := uniq[p]; !duplicate {
			uniq[p] = struct{}{}
			paths[i] = utils.StripPathParts(p, -2)
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
	if runtime.GOOS == "linux" {
		scanRoot = "/home"
	} else if runtime.GOOS == "windows" {
		scanRoot = "/Users"
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

	if len(paths) > 1 {
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
	} else if len(paths) == 1 {
		return paths[0], nil
	}

	return
}

// Execute executes the root command once Initialize has been called
func Execute() error {
	defaultCfgFile = fmt.Sprintf("%s.%s", PROGRAM, cfgExt)

	rootCmd.Use = PROGRAM
	rootCmd.Version = fmt.Sprintf("%s, version %s-(%s-%s)", PROGRAM, VERSION, runtime.GOARCH, runtime.GOOS)

	// ? should not require RAW mode just go get the dimensions...
	if w, _, err := terminal.GetSize(0); err == nil {
		termWidth = w
	}

	// Load Configuration on Initialize
	cobra.OnInitialize(loadConfig)

	// Attach Flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is $HOME/%s)", defaultCfgFile))
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable debug logging for troubleshooting sake")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	// Add Commands
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(versionCmd)

	return rootCmd.Execute()
}
