package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/fsnotify/fsnotify"
	"github.com/google/go-github/v32/github"
	"github.com/kataras/golog"
	"github.com/spf13/viper"

	"github.com/AlbinoGeek/sc2-rsu/cmd/gui"
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

	accList    *fyne.Container
	tabs       *container.AppTabs
	uploadList *widget.Table

	settings *tabSettings
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

	main.Refresh()

	w.Resize(fyne.NewSize(700, 600))
	w.CenterOnScreen()
	w.Show()

	main.setupUploader()

	if viper.GetString("version") == "" || viper.GetString("apikey") == "" {
		main.openGettingStarted1()
	}
}

func (main *windowMain) Refresh() {
	main.genAccountList()
	main.genUploadList()

	tblName := newText("Map Name", 1, true)
	tblName.Move(fyne.NewPos(8, 3))

	tblID := newText("ID", 1, true)
	tblID.Move(fyne.NewPos(248, 3))

	tblStatus := newText("Status", 1, true)
	tblStatus.Move(fyne.NewPos(334, 3))

	sourceURL, _ := url.Parse(ghLink(""))

	if main.settings == nil {
		main.settings = &tabSettings{Window: main}
	}

	main.tabs = container.NewAppTabs(
		container.NewTabItemWithIcon("", theme.HomeIcon(),
			container.NewVScroll(
				layout.NewSpacer(),
			),
		),

		container.NewTabItemWithIcon("", accIcon,
			container.NewVScroll(main.accList),
		),

		container.NewTabItemWithIcon("", uploadIcon,
			container.NewBorder(
				fyne.NewContainerWithoutLayout(
					tblName, tblID, tblStatus,
				), nil, nil, nil, main.uploadList,
			),
		),

		container.NewTabItemWithIcon("", theme.SettingsIcon(),
			container.NewVScroll(main.settings.Init()),
		),

		container.NewTabItemWithIcon("", theme.InfoIcon(),
			container.NewCenter(widget.NewVBox(
				widget.NewHBox(
					layout.NewSpacer(),
					widget.NewVBox(
						newHeader(PROGRAM),
						widget.NewHyperlink("Browse Source", sourceURL),
					),
					layout.NewSpacer(),
					widget.NewForm(
						widget.NewFormItem("Author", widget.NewLabel(ghOwner)),
						widget.NewFormItem("Version", widget.NewLabel(VERSION)),
					),
					layout.NewSpacer(),
				),
				widget.NewVBox(
					widget.NewButtonWithIcon("Check for Updates", theme.ViewRefreshIcon(), func() { go main.checkUpdate() }),
					widget.NewButtonWithIcon("Request Feedback", feedbackIcon, main.OpenGitHub("issues/new?assignees=AlbinoGeek&labels=enhancement&template=feature-request.md&title=%5BFEATURE+REQUEST%5D")),
					widget.NewButtonWithIcon("Report A Bug", reportBugIcon, main.OpenGitHub("issues/new?assignees=AlbinoGeek&labels=bug&template=bug-report.md&title=%5BBUG%5D")),
				),
			)),
		),
	)

	main.tabs.SetTabLocation(widget.TabLocationTrailing)
	main.GetWindow().SetContent(main.tabs)
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

func (main *windowMain) checkUpdate() {
	w := main.GetWindow()
	dlg := dialog.NewProgressInfinite("Check for Updates", "Checking for new releases...", w)
	dlg.Show()

	rel := checkUpdate()

	dlg.Hide()

	if rel == nil {
		dialog.ShowInformation("Check for Updates", "No updates are available at this time.", w)
		return
	}

	dialog.ShowConfirm("Update Available!",
		fmt.Sprintf("You are running version %s.\nAn update is available: %s\nWould you like us to download it now?", VERSION, rel.GetTagName()),
		main.doUpdate(rel), main.GetWindow())
}

func (main *windowMain) doUpdate(rel *github.RepositoryRelease) func(bool) {
	return func(ok bool) {
		if !ok {
			return
		}

		w := main.GetWindow()

		// otherwise we might block the fyne event queue...
		go func() {
			// TODO: display download progress, filename and size
			dlg := dialog.NewProgressInfinite("Downloading Update",
				fmt.Sprintf("Downloading version %s nomain...", rel.GetTagName()), w)
			dlg.Show()

			err := downloadUpdate(rel)

			dlg.Hide()

			if err != nil {
				dialog.ShowError(err, w)
			} else {
				dialog.ShowInformation("Update Complete!", "Please close the program and start the new binary.", w)
			}
		}()
	}
}

func (main *windowMain) genUploadList() {
	main.uploadList = widget.NewTable(
		func() (int, int) { return len(main.uploadStatus), 3 },
		func() fyne.CanvasObject {
			return newText("@@@@@@@@", 1, false)
		},
		func(tci widget.TableCellID, f fyne.CanvasObject) {
			l := f.(*canvas.Text)
			switch atom := main.uploadStatus[tci.Row]; tci.Col {
			case 0:
				l.Text = atom.MapName
			case 1:
				l.Text = atom.ReplayID
			case 2:
				l.Text = atom.Status
			}
			l.Refresh()
		},
	)
	main.uploadList.OnSelected = func(id widget.TableCellID) {
		if id.Row > len(main.uploadStatus)-1 {
			return // selected row that does not exist
		}

		if rid := main.uploadStatus[id.Row].ReplayID; rid != "" {
			u, _ := url.Parse(fmt.Sprintf("%s/replay/%s", sc2replaystats.WebRoot, rid))
			main.App.OpenURL(u)
		}
	}
	main.uploadList.SetColumnWidth(0, 230)
	main.uploadList.SetColumnWidth(1, 76)
	main.uploadList.SetColumnWidth(2, 84)
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

	tries := 0
	wait := time.Second * 3

	// naive retry logic
	for {
		tries++
		if tries > 3 {
			return
		}

		main.uploadList.Refresh()

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

		main.uploadList.Refresh()

		rqid, err := sc2api.UploadReplay(replayFilename)
		entry.QueueID = rqid

		if err != nil {
			entry.Status = "u failed"

			main.uploadList.Refresh()

			dialog.NewError(fmt.Errorf("replay upload failed:%v\n%v", mapName, err), main.GetWindow())
		} else {
			entry.Status = "processing"
			main.uploadList.Refresh()
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
			main.openGettingStarted2()
		} else {
			main.gettingStarted = 0
			main.modal.Hide()
		}
	},
		labelWithWrapping("You are only two steps away from having your replays automatically uploaded!"),
		labelWithWrapping("1) We will find your Replays Directory"),
		labelWithWrapping("2) We will find your sc2replaystats API Key"),
	)
}

func (main *windowMain) openGettingStarted2() {
	main.gettingStarted = 2

	btnSettings := widget.NewButtonWithIcon("Open Settings", theme.SettingsIcon(), func() {
		main.UI.OpenWindow(WindowSettings)
	})
	btnSettings.Importance = widget.HighImportance

	// TODO: Refactor this to actually have the settings UI, not just direct the user to settings
	main.WizardModal("", "", nil, nil,
		labelWithWrapping("First thing's first. Please use the button below to open the Settings dialog, and under the StarCraft II section, add your Replays Directory."),
		btnSettings,
		labelWithWrapping("Once you have found your replays directory and saved the settings, this setup wizard will automatically advance to the next step."),
	)
}

func (main *windowMain) openGettingStarted3() {
	main.gettingStarted = 3

	btnSettings := widget.NewButtonWithIcon("Open Settings", theme.SettingsIcon(), func() {
		main.UI.OpenWindow(WindowSettings)
	})
	btnSettings.Importance = widget.HighImportance

	// TODO: Refactor this to actually have the settings UI, not just direct the user to settings
	main.WizardModal("", "", nil, nil,
		labelWithWrapping("Lastly, please set your sc2replaystats API key. If you do not know how to find this, use the \"Login and find it for me\" button to have us login to your account and generate one on your behalf."),
		btnSettings,
	)
}

func (main *windowMain) openGettingStarted4() {
	main.gettingStarted = 0

	main.WizardModal("Close", "", func() {
		main.gettingStarted = 0
		main.modal.Hide()
	}, nil,
		labelWithWrapping("Contratulations! You have finished first-time setup. You can change these settings at any time by going to File -> Settings."),
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
	for {
		time.Sleep(time.Second)

		rid, err := sc2api.GetReplayStatus(entry.QueueID)

		if err != nil {
			golog.Errorf("error checking reply status: %v: %v", entry.QueueID, err)
			entry.Status = "p failed"

			main.uploadList.Refresh()

			return err // could not check status
		}

		if rid != "" {
			entry.Status = "success"
			entry.ReplayID = rid

			main.uploadList.Refresh()

			return nil // replay parsed!
		}

		golog.Debugf("sc2replaystats process..: [%v] %s", entry.QueueID, rid)
	}
}
