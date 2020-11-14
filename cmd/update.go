package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/google/go-github/v32/github"
	"github.com/kataras/golog"
	"github.com/mitchellh/go-wordwrap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	stripmd "github.com/writeas/go-strip-markdown"

	"github.com/AlbinoGeek/sc2-rsu/utils"
)

var (
	ghClient = github.NewClient(nil)
	ghOwner  = "AlbinoGeek"
	ghRepo   = "sc2-rsu"

	minimumUpdatePeriod = time.Minute * 10

	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Checks for and optionally downloads program updates",
		Run: func(cmd *cobra.Command, args []string) {
			golog.Infof("You are currently running version: %v", VERSION)

			rel := checkUpdate()
			if rel == nil {
				golog.Info("No updates found. You are on the latest release version.")
				return
			}

			printUpdateNotice(rel)
			consoleReader := bufio.NewReaderSize(os.Stdin, 1)
		outer:
			for {
				fmt.Print("Would you like to download the update? [y/N]: ")
				input, _, _ := consoleReader.ReadLine()
				switch string(input) {
				case "Y", "y":
					break outer
				case "N", "n":
					golog.Warn("declined automatic update")
					return
				}
			}

			if err := downloadUpdate(rel); err != nil {
				golog.Fatalf("failed to download update: %s", err)
			}
		},
	}
)

func checkUpdateEvery(period time.Duration) {
	func() {
		for {
			tt := time.Now()
			if rel := checkUpdate(); rel != nil {
				printUpdateNotice(rel)

				if viper.GetBool("update.automatic.enabled") {
					if err := downloadUpdate(rel); err != nil {
						golog.Errorf("failed to download update: %s", err)
					}
				}
				break // only notify the user once
			}
			golog.Debugf("update check took: %v", time.Since(tt))
			time.Sleep(period)
		}
	}()
}

func checkUpdate() *github.RepositoryRelease {
	rels, _, err := ghClient.Repositories.ListReleases(context.TODO(), ghOwner, ghRepo, nil)
	if err != nil {
		golog.Errorf("failed update check, could not list releases: %v", err)
		return nil
	}

	for _, rel := range rels {
		if tag := rel.GetTagName(); utils.CompareSemVer(VERSION, tag) > 0 {
			return rel
		}
	}

	return nil
}

func downloadUpdate(rel *github.RepositoryRelease) error {
	golog.Info("Starting update...")

	var asset *github.ReleaseAsset
	for _, a := range rel.Assets {
		if name := a.GetName(); strings.Contains(name, runtime.GOARCH) &&
			strings.Contains(name, runtime.GOOS) {
			if a.GetBrowserDownloadURL() != "" {
				golog.Debugf("candidate release asset: %v", name)
				asset = a
				break
			}
		}
	}
	if asset == nil {
		return fmt.Errorf("found no suitable packages for your OS/ARCH -- please report this bug")
	}

	absName, err := filepath.Abs(os.Args[0])
	if err != nil {
		return fmt.Errorf("unable to determine executing directory: %v", err)
	}

	fpath, fname, ext := utils.SplitFilepath(absName)
	tag := rel.GetTagName()
	if tag[0] == 'v' {
		tag = tag[1:]
	}
	fname = filepath.Join(fpath, fmt.Sprintf("%s-%s%s", fname, tag, ext))

	golog.Debugf("creating file: %v", fname)
	tmp, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer tmp.Close()

	golog.Infof("Downloading %s update: %v", humanize.Bytes(uint64(asset.GetSize())), asset.GetName())
	if err = utils.DownloadWriter(asset.GetBrowserDownloadURL(), tmp); err != nil {
		return err
	}

	if err = tmp.Chmod(0755); err != nil {
		return fmt.Errorf("failed to mark file executable: %v", err)
	}

	golog.Infof("Update complete!\n\nPlease close the program and start the new version: %s", fname)
	return nil
}

func getUpdateDuration() time.Duration {
	dur := viper.GetString("update.check.period")
	period, err := time.ParseDuration(dur)
	if err != nil || period < minimumUpdatePeriod {
		golog.Warnf("update.check.period invalid or too short: %v", err)
		period = time.Duration(minimumUpdatePeriod)
	}

	return period
}

func printUpdateNotice(rel *github.RepositoryRelease) {
	line := strings.Repeat("=", termWidth)
	fmt.Printf(
		"%s\nUpdate Detected! New Release Version: %v\n\n%s\n%s\n",
		line,
		rel.GetTagName(),
		wordwrap.WrapString(stripmd.Strip(rel.GetBody()), uint(termWidth)),
		line)
}
