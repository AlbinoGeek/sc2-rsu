package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/kataras/golog"
	"github.com/mitchellh/go-wordwrap"
	stripmd "github.com/writeas/go-strip-markdown"
)

var (
	ghClient = github.NewClient(nil)
	ghOwner  = "AlbinoGeek"
	ghRepo   = "sc2-rsu"
)

func isNewer(old, new string) bool {
	if len(new) < 1 {
		return false
	}

	nparts := strings.Split(strings.Split(new[1:], "-")[0], ".")
	oparts := strings.Split(old, ".")
	for i := range nparts {
		if len(oparts)-1 < i {
			break
		}

		n, err1 := strconv.Atoi(nparts[i])
		o, err2 := strconv.Atoi(oparts[i])
		if err1 == nil && err2 == nil && n > o {
			return true
		}
	}
	return false
}

func updateCheckEvery(period time.Duration) {
	func() {
		for {
			tt := time.Now()
			if rel := updateCheck(); rel != nil {
				line := strings.Repeat("=", termWidth)
				fmt.Printf(
					"%s\nUpdate Detected! New Release Version: %v\n\n%s\n%s\n",
					line,
					rel.GetTagName(),
					wordwrap.WrapString(stripmd.Strip(rel.GetBody()), uint(termWidth)),
					line)
				break // only notify the user once
			}
			golog.Debugf("update check took: %v", time.Since(tt))
			time.Sleep(period)
		}
	}()
}

func updateCheck() *github.RepositoryRelease {
	rels, _, err := ghClient.Repositories.ListReleases(context.TODO(), ghOwner, ghRepo, nil)
	if err != nil {
		golog.Errorf("failed update check, could not list releases: %v", err)
		return nil
	}

	for _, rel := range rels {
		if tag := rel.GetTagName(); isNewer(VERSION, tag) {
			return rel
		}
	}

	return nil
}
