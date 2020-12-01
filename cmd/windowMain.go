package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/fsnotify/fsnotify"
	"github.com/kataras/golog"
	"github.com/spf13/viper"

	"github.com/AlbinoGeek/sc2-rsu/cmd/gui"
	"github.com/AlbinoGeek/sc2-rsu/cmd/gui/fynewidget"
	"github.com/AlbinoGeek/sc2-rsu/sc2replaystats"
	"github.com/AlbinoGeek/sc2-rsu/sc2utils"
	"github.com/AlbinoGeek/sc2-rsu/utils"
)

type windowMain struct {
	*gui.WindowBase
	gettingStarted uint
	modal          *widget.PopUp
	uploadEnabled  map[string]bool
	uploadStatus   []*uploadRecord
	watcher        *fsnotify.Watcher

	nav *fynewidget.NavigationDrawer

	// Tabs
	accounts *paneAccounts
	uploads  *paneUploads
	settings *paneSettings
}

type uploadRecord struct {
	Filename string
	Filesize string
	MapName  string
	QueueID  string
	ReplayID string
	Status   string
}

func (main *windowMain) Init() {
	w := main.App.NewWindow("SC2ReplayStats Uploader")
	main.SetWindow(w)

	main.uploadEnabled = make(map[string]bool)
	main.uploadStatus = make([]*uploadRecord, 0)

	// closing the main window should quit the application
	w.SetCloseIntercept(func() {
		if main.settings.unsaved {
			main.settings.onClose()
			return
		}

		if main.watcher != nil {
			main.watcher.Close()
		}

		w.Close()
		main.App.Quit()
	})

	if sc2api == nil {
		sc2api = sc2replaystats.New(viper.GetString("apikey"))
	}

	main.accounts = makePaneAccounts(main).(*paneAccounts)
	main.uploads = makePaneUploads(main).(*paneUploads)
	main.settings = makePaneSettings(main).(*paneSettings)

	main.nav = fynewidget.NewNavigationDrawer(
		PROGRAM,
		"",
		// NewNavigationLabelWithIcon("Overview", theme.HomeIcon(),
		// 	widget.NewLabel("Overview Contents"),
		// ),
		fynewidget.NewNavigationLabelWithIcon("Accounts", accIcon,
			main.accounts.GetContent(),
		),
		fynewidget.NewNavigationLabelWithIcon("Uploads", uploadIcon,
			main.uploads.GetContent(),
		),
		fynewidget.NewNavigationSeparator(),
		fynewidget.NewNavigationLabelWithIcon("Settings", theme.SettingsIcon(),
			main.settings.GetContent(),
		),
		fynewidget.NewNavigationLabelWithIcon("Help & Feedback", feedbackIcon,
			makePaneAbout(main).GetContent(),
		),
	)
	main.nav.SetImage(theme.InfoIcon())

	content := container.NewMax(layout.NewSpacer()) // ? what's a better way ?
	main.nav.OnSelect = func(ni fynewidget.NavigationItem) {
		content.Objects = []fyne.CanvasObject{ni.GetContent()}
	}
	main.nav.Select(0)

	main.GetWindow().SetContent(container.NewBorder(nil, nil,
		main.nav, nil,
		content))

	w.Resize(fyne.NewSize(600, 480))
	w.CenterOnScreen()
	w.Show()

	main.setupUploader()

	if viper.GetString("version") == "" || viper.GetString("apikey") == "" {
		main.openGettingStarted1()
	}
}

func (main *windowMain) WizardModal(skipText, nextText string, skipFn, nextFn func(), contents ...fyne.CanvasObject) {
	if skipFn == nil {
		skipFn = func() { main.modal.Hide() }
	}

	if nextFn == nil {
		nextFn = func() { main.modal.Hide() }
	}

	buttons := make([]fyne.CanvasObject, 0)

	if skipText != "" {
		btn := widget.NewButtonWithIcon(skipText, theme.CancelIcon(), skipFn)
		btn.Importance = widget.LowImportance
		buttons = append(buttons, btn)
	}

	if nextText != "" {
		btn := widget.NewButtonWithIcon(nextText, theme.NavigateNextIcon(), nextFn)
		btn.Importance = widget.HighImportance
		buttons = append(buttons, btn)
	}

	// Re-use existing Modal
	if main.modal != nil {
		box := main.modal.Content.(*widget.Box)
		box.Children[0].(*widget.Card).Content.(*fyne.Container).Objects = contents
		box.Children[len(box.Children)-1] = fyne.NewContainerWithLayout(
			layout.NewGridLayout(len(buttons)),
			buttons...,
		)

		main.modal.Show()
		main.modal.Refresh()

		return
	}

	// Create Fresh Modal
	container := fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		contents...,
	)

	box := widget.NewVBox(
		widget.NewCard("Welcome!", "First-Time Setup",
			container,
		),
		layout.NewSpacer(),
		widget.NewSeparator(),
		fyne.NewContainerWithLayout(
			layout.NewGridLayout(len(buttons)),
			buttons...,
		),
	)

	main.modal = widget.NewModalPopUp(box, main.GetWindow().Canvas())
	main.modal.Show()

	size := fyne.NewSize(360, 240)
	main.modal.Resize(size)
	box.Resize(size)
}

func (main *windowMain) handleReplay(replayFilename string) {
	_, mapName, _ := utils.SplitFilepath(replayFilename)
	entry := &uploadRecord{
		Filename: replayFilename,
		MapName:  mapName,
		Status:   "pending",
	}

	golog.Debugf("uploading replay: %v", replayFilename)

	main.uploadStatus = append(main.uploadStatus, entry)
	main.uploads.Refresh()

	tries := 0
	wait := time.Second * 3

	// naive retry logic
	for {
		tries++
		if tries > 3 {
			return
		}

		// wait for the replay to have finished being written (large enough filesize)
		var lastSize int64

		for {
			time.Sleep(time.Millisecond * 250)

			// ! smallest replay I've seen is 27418 bytes (-3 second long)
			if s, err := os.Stat(replayFilename); err == nil && s.Size() > validReplaySize {
				// check that the replay has stopped growing
				if s.Size() > lastSize {
					lastSize = s.Size()
				} else {
					break
				}
			}
		}

		entry.Status = "uploading"
		main.uploads.Refresh()

		rqid, err := sc2api.UploadReplay(replayFilename)
		entry.QueueID = rqid

		if err != nil {
			entry.Status = "u failed"
			main.uploads.Refresh()

			dialog.NewError(fmt.Errorf("replay upload failed:%v\n%v", mapName, err), main.GetWindow())
		} else {
			entry.Status = "processing"
			main.uploads.Refresh()

			if err := main.watchReplayStatus(entry); err == nil {
				return
			}
		}

		time.Sleep(wait)
		wait *= 2
	}
}

// OpenGitHub launches the user's browser to a given GitHub URL relative to
// this project's repository root
func (main *windowMain) OpenGitHub(slug string) func() {
	u, _ := url.Parse(ghLink(slug))

	return func() {
		if err := main.UI.App.OpenURL(u); err != nil {
			dialog.ShowError(err, main.GetWindow())
		}
	}
}

func (main *windowMain) openGettingStarted1() {
	main.gettingStarted = 1
	main.WizardModal("Skip", "Next", nil, func() {
		if viper.GetString("replaysroot") == "" {
			main.nav.Select(3) // ! ID BASED IS ERROR PRONE
		} else {
			main.gettingStarted = 0
		}
		main.modal.Hide()
	},
		labelWithWrapping("You are only two steps away from having your replays automatically uploaded!"),
		labelWithWrapping("1) We will find your Replays Directory"),
		labelWithWrapping("2) We will find your sc2replaystats API Key"),
	)
}

// func (main *windowMain) openGettingStarted2() {
// 	main.gettingStarted = 2

// 	btnSettings := widget.NewButtonWithIcon("Open Settings", theme.SettingsIcon(), func() {
// 		main.tabs.Select(3) // ! ID BASED IS ERROR PRONE
// 	})
// 	btnSettings.Importance = widget.HighImportance

// 	// TODO: Refactor this to actually have the settings UI, not just direct the user to settings
// 	main.WizardModal("", "", nil, nil,
// 		labelWithWrapping("First thing's first. Please use the button below to open the Settings dialog, and under the StarCraft II section, add your Replays Directory."),
// 		btnSettings,
// 		labelWithWrapping("Once you have found your replays directory and saved the settings, this setup wizard will automatically advance to the next step."),
// 	)
// }

// func (main *windowMain) openGettingStarted3() {
// 	main.gettingStarted = 3

// 	btnSettings := widget.NewButtonWithIcon("Open Settings", theme.SettingsIcon(), func() {
// 		main.tabs.Select(3) // ! ID BASED IS ERROR PRONE
// 	})
// 	btnSettings.Importance = widget.HighImportance

// 	// TODO: Refactor this to actually have the settings UI, not just direct the user to settings
// 	main.WizardModal("", "", nil, nil,
// 		labelWithWrapping("Lastly, please set your sc2replaystats API key. If you do not know how to find this, use the \"Login and find it for me\" button to have us login to your account and generate one on your behalf."),
// 		btnSettings,
// 	)
// }

// func (main *windowMain) openGettingStarted4() {
// 	main.gettingStarted = 0

// 	main.WizardModal("Close", "", func() {
// 		main.gettingStarted = 0
// 		main.modal.Hide()
// 	}, nil,
// 		labelWithWrapping("Contratulations! You have finished first-time setup. You can change these settings at any time by going to File -> Settings."),
// 	)
// }

func (main *windowMain) setupUploader() {
	w := main.GetWindow()
	replaysRoot := viper.GetString("replaysRoot")

	if replaysRoot == "" {
		return
	}

	accs, err := sc2utils.EnumerateAccounts(replaysRoot)
	if err != nil {
		dialog.NewError(err, w)
	}

	paths := make([]string, len(accs))
	for i, a := range accs {
		paths[i] = filepath.Join(replaysRoot, a, "Replays", "Multiplayer")
	}

	// TODO : should just clear watch paths instead of making a new watcher
	// in case we were setup again (replaysRoot changed)
	if main.watcher != nil {
		main.watcher.Close()
	}

	watch, err := newWatcher(paths)
	if err != nil {
		dialog.NewError(fmt.Errorf("Failed to start uploader:\n%v", err), w)
		return
	}

	main.watcher = watch

	go func() {
		for {
			select {
			case event, ok := <-watch.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					// bug: SC2 sometime writes out ".SC2Replay.writeCacheBackup" files
					if strings.HasSuffix(event.Name, "eplay") {
						go main.handleReplay(event.Name)
					}
				}
			case err, ok := <-watch.Errors:
				if !ok {
					return
				}

				golog.Warnf("fswatcher error: %v", err)
			}
		}
	}()
}

func (main *windowMain) toggleUploading(btn *widget.Button, id string) func() {
	return func() {
		w := main.GetWindow()
		replaysRoot := viper.GetString("replaysRoot")

		main.uploadEnabled[id] = !main.uploadEnabled[id]
		if main.uploadEnabled[id] {
			if err := main.watcher.Remove(filepath.Join(replaysRoot, id, "Replays", "Multiplayer")); err != nil {
				dialog.NewError(err, w)
				return
			}

			btn.Importance = widget.HighImportance
			btn.Icon = theme.MediaPauseIcon()
		} else {
			if err := main.watcher.Add(filepath.Join(replaysRoot, id, "Replays", "Multiplayer")); err != nil {
				dialog.NewError(err, w)
				return
			}

			btn.Importance = widget.MediumImportance
			btn.Icon = theme.MediaPlayIcon()
		}
	}
}

func (main *windowMain) watchReplayStatus(entry *uploadRecord) error {
	defer main.uploads.Refresh()
	for {
		time.Sleep(time.Second)

		rid, err := sc2api.GetReplayStatus(entry.QueueID)

		if err != nil {
			golog.Errorf("error checking reply status: %v: %v", entry.QueueID, err)
			entry.Status = "p failed"
			return err // could not check status
		}

		if rid != "" {
			entry.Status = "success"
			entry.ReplayID = rid
			return nil // replay parsed!
		}

		golog.Debugf("sc2replaystats process..: [%v] %s", entry.QueueID, rid)
	}
}
