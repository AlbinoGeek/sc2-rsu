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
	"github.com/AlbinoGeek/sc2-rsu/fynex"
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

	nav    *fynex.NavDrawer
	topbar *fynex.AppBar

	// Panes
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
	main.uploadEnabled = make(map[string]bool)
	main.uploadStatus = make([]*uploadRecord, 0)

	w := main.App.NewWindow("SC2ReplayStats Uploader")
	w.SetPadded(false)
	w.Resize(fyne.NewSize(640, 560))
	w.CenterOnScreen()
	main.SetWindow(w)

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

	main.topbar = fynex.NewAppBar(PROGRAM)
	main.nav = fynex.NewNavDrawer(
		PROGRAM,
		"",
		main.accounts,
		main.uploads,
		fynex.NewNavSeparator(),
		main.settings,
		makePaneAbout(main),
		// fynex.NewNavSeparator(),
		// makePaneDeveloper(main),
	)
	// main.nav.SetImage(theme.InfoIcon())
	main.topbar.SetNav(main.nav)

	mobile := fyne.CurrentDevice().IsMobile()
	content := container.NewPadded(layout.NewSpacer()) // ? what's a better way ?
	main.nav.OnSelect = func(ni fynex.NavItem) {
		content.Objects = []fyne.CanvasObject{ni.GetContent()}

		if mobile {
			main.topbar.SetTitle(ni.GetTitle())
			main.topbar.SetNavClosed(true)
		}
	}

	if mobile {
		main.GetWindow().SetContent(
			container.NewBorder(
				nil, nil, main.nav, nil,
				container.NewBorder(
					main.topbar,
					nil,
					nil,
					nil,
					content,
				),
			),
		)

		main.topbar.SetNavClosed(true)
	} else {
		main.GetWindow().SetContent(
			container.NewBorder(
				main.topbar,
				nil,
				main.nav,
				nil,
				content,
			),
		)

		main.nav.SetTitle("")
	}

	w.Show()

	main.nav.Select(0) // Cannot select before window is shown!
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
	labelWithWrapping := func(text string) *widget.Label {
		label := widget.NewLabel(text)
		label.Wrapping = fyne.TextWrapWord

		return label
	}

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

	paths := make([]string, 0)

	for _, a := range accs {
		if getToonEnabled(a) {
			paths = append(paths, filepath.Join(replaysRoot, a, "Replays", "Multiplayer"))
		}
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
			if err := main.watcher.Add(filepath.Join(replaysRoot, id, "Replays", "Multiplayer")); err != nil {
				dialog.NewError(err, w)

				return
			}

			btn.Importance = widget.HighImportance
			btn.Icon = theme.MediaPauseIcon()
		} else {
			if err := main.watcher.Remove(filepath.Join(replaysRoot, id, "Replays", "Multiplayer")); err != nil {
				dialog.NewError(err, w)

				return
			}

			btn.Importance = widget.MediumImportance
			btn.Icon = theme.MediaPlayIcon()
		}

		main.updateEnabledToons()
	}
}

func (main *windowMain) updateEnabledToons() {
	enabledToons := make([]string, 0)

	for key, enabled := range main.uploadEnabled {
		if enabled {
			enabledToons = append(enabledToons, key)
		}
	}

	setToons(enabledToons)
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
