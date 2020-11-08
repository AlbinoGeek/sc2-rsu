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

	"github.com/fsnotify/fsnotify"
	"github.com/kataras/golog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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

			if !utils.ValidAPIKey(key) {
				return errors.New("invalid API key in configuration, please replace it or use the login command")
			}

			golog.Infof("Starting Automatic Replay Uploader...")
			return automaticUpload(key)
		},
	}
)

func automaticUpload(apikey string) error {
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
					golog.Infof("Detected new replay file: %v", event.Name)
					if err := uploadReplay(apikey, event.Name); err != nil {
						golog.Errorf("failed to upload replay: %v: %v", event.Name, err)
					}
				}
				// if event.Op&fsnotify.Write == fsnotify.Write {
				// 	log.Println("modified file:", event.Name)
				// }
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
		golog.Warnf("Replay Root not configured correctly, searching for replays directory...")
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

	golog.Info("Ready!")
	<-done
	return nil
}

func findAccounts(root string) (ids []string, err error) {
	golog.Debugf("Searching for accounts in replay directory...")

	paths := make([]string, 0)
	uniq := make(map[string]struct{})
	err = filepath.Walk(root, func(path string, info os.FileInfo, ferr error) error {
		if ferr != nil {
			// silently skip access / permission errors
			if strings.Contains(ferr.Error(), "access denied") ||
				strings.Contains(ferr.Error(), "permission denied") {
				return nil
			}

			return ferr
		}

		if info.IsDir() && strings.HasSuffix(path, "ultiplayer") {
			if parts := strings.Split(path, string(filepath.Separator)); len(parts) > 2 {
				path = strings.Join(parts[len(parts)-4:len(parts)-2], string(filepath.Separator))
				if _, is := uniq[path]; !is {
					uniq[path] = struct{}{}
					paths = append(paths, path)
					golog.Debugf("found candidate account: %v", path)
				}
			}
		}

		return nil
	})
	golog.Debugf("finished scanning for candidates. Found: %v", len(paths))
	return paths, err

	// 	if len(paths) > 1 {
	// 		fmt.Print(`============================================================
	// More than one possible replay directory was located while we
	// scanned for your StarCraft II installation's Accounts folder

	// Please select which directory we should be watching below:`)
	// 		for i, p := range paths {
	// 			fmt.Printf("\n  %d: %s", 1+i, p)
	// 		}
	// 		fmt.Println("\n============================================================")
	// 		consoleReader := bufio.NewReaderSize(os.Stdin, 1)
	// 		for {
	// 			fmt.Printf("Your Choice [1-%d]: ", len(paths))
	// 			input, _, _ := consoleReader.ReadLine()
	// 			choice, err := strconv.Atoi(string(input))
	// 			if err == nil && choice-1 > 0 && choice-1 < len(paths) {
	// 				return paths[choice-1], nil
	// 			}
	// 		}
	// 	} else if len(paths) == 1 {
	// 		return paths[0], nil
	// 	}
}

func findReplaysRoot() (root string, err error) {
	golog.Infof("Determining replays directory... (this could take a few minutes)...")

	paths := make([]string, 0)
	uniq := make(map[string]struct{})
	err = filepath.Walk("/", func(path string, info os.FileInfo, ferr error) error {
		if ferr != nil {
			// silently skip access / permission errors
			if strings.Contains(ferr.Error(), "access denied") ||
				strings.Contains(ferr.Error(), "permission denied") {
				return nil
			}

			return ferr
		}

		if info.IsDir() &&
			strings.Contains(path, "ccounts") &&
			strings.Contains(path, "eplays") &&
			strings.HasSuffix(path, "ultiplayer") {

			if parts := strings.Split(path, string(filepath.Separator)); len(parts) > 4 {
				path = strings.Join(parts[:len(parts)-4], string(filepath.Separator))
				if _, is := uniq[path]; !is {
					uniq[path] = struct{}{}
					paths = append(paths, path)
					golog.Debugf("found candidate replay directory: %v", path)
				}
			}
		}

		return nil
	})
	golog.Debugf("finished scanning for candidates. Found: %v", len(paths))

	if len(paths) > 1 {
		fmt.Print(`============================================================
More than one possible replay directory was located while we
scanned for your StarCraft II installation's Accounts folder

Please select which directory we should be watching below:`)
		for i, p := range paths {
			fmt.Printf("\n  %d: %s", 1+i, p)
		}
		fmt.Println("\n============================================================")
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
