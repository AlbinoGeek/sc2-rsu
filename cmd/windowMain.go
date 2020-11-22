package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
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

	"github.com/AlbinoGeek/sc2-rsu/sc2replaystats"
	"github.com/AlbinoGeek/sc2-rsu/sc2utils"
	"github.com/AlbinoGeek/sc2-rsu/utils"
)

type windowMain struct {
	*windowBase
	gettingStarted uint
	modal          *widget.PopUp
	uploadEnabled  map[string]bool
	uploadStatus   []*uploadRecord
	watcher        *fsnotify.Watcher

	accList    *fyne.Container
	uploadList *widget.Table
}

type uploadRecord struct {
	Filename string
	Filesize string
	MapName  string
	QueueID  string
	ReplayID string
	Status   string
}

func (w *windowMain) Init() {
	w.windowBase.Window = w.windowBase.app.NewWindow("SC2ReplayStats Uploader")

	w.uploadEnabled = make(map[string]bool)
	w.uploadStatus = make([]*uploadRecord, 0)

	w.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Check for Updates", func() { go w.checkUpdate() }),
			fyne.NewMenuItem("Settings", func() { w.ui.OpenWindow(WindowSettings) }),
		),
		fyne.NewMenu("Help",
			fyne.NewMenuItem("Report Bug", w.ui.OpenGitHub("issues/new?assignees=AlbinoGeek&labels=bug&template=bug-report.md&title=%5BBUG%5D")),
			fyne.NewMenuItem("Request Feature", w.ui.OpenGitHub("issues/new?assignees=AlbinoGeek&labels=enhancement&template=feature-request.md&title=%5BFEATURE+REQUEST%5D")),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("About", func() { w.ui.OpenWindow(WindowAbout) }),
		),
	))

	// closing the main window should quit the application
	w.SetCloseIntercept(func() {
		// Close "About" if it's open
		if win := w.ui.windows[WindowAbout].GetWindow(); win != nil {
			win.Close()
		}

		win := w.ui.windows[WindowSettings]
		if win.GetWindow() != nil {
			settings := win.(*windowSettings)
			if settings.unsaved {
				settings.onClose()
				return
			}
			settings.Close()
		}

		if w.watcher != nil {
			w.watcher.Close()
		}

		w.Close()
		w.app.Quit()
	})

	if sc2api == nil {
		sc2api = sc2replaystats.New(viper.GetString("apikey"))
	}
	w.Refresh()

	w.Resize(fyne.NewSize(420, 360))
	w.CenterOnScreen()
	w.Show()

	w.setupUploader()

	if viper.GetString("version") == "" || viper.GetString("apikey") == "" {
		w.openGettingStarted1()
	}
}

func (w *windowMain) Refresh() {
	w.genAccountList()
	w.genUploadList()

	tblName := newText("Map Name", 1, true)
	tblName.Move(fyne.NewPos(8, 3))

	tblID := newText("ID", 1, true)
	tblID.Move(fyne.NewPos(248, 3))

	tblStatus := newText("Status", 1, true)
	tblStatus.Move(fyne.NewPos(334, 3))

	w.SetContent(container.NewAppTabs(
		container.NewTabItem("Accounts",
			container.NewVScroll(w.accList),
		),
		container.NewTabItem("Uploads",
			container.NewBorder(
				fyne.NewContainerWithoutLayout(
					tblName, tblID, tblStatus,
				), nil, nil, nil, w.uploadList,
			),
		),
	))
}

func (w *windowMain) WizardModal(skipText, nextText string, skipFn, nextFn func(), contents ...fyne.CanvasObject) {
	if skipFn == nil {
		skipFn = func() { w.modal.Hide() }
	}
	if nextFn == nil {
		nextFn = func() { w.modal.Hide() }
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
	if w.modal != nil {
		box := w.modal.Content.(*widget.Box)
		box.Children[0].(*widget.Card).Content.(*fyne.Container).Objects = contents
		box.Children[len(box.Children)-1] = fyne.NewContainerWithLayout(
			layout.NewGridLayout(len(buttons)),
			buttons...,
		)

		w.modal.Show()
		w.modal.Refresh()
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

	w.modal = widget.NewModalPopUp(box, w.Canvas())
	w.modal.Show()

	size := fyne.NewSize(360, 240)
	w.modal.Resize(size)
	box.Resize(size)
}

func (w *windowMain) checkUpdate() {
	dlg := dialog.NewProgressInfinite("Check for Updates", "Checking for new releases...", w)
	dlg.Show()
	rel := checkUpdate()
	dlg.Hide()

	if rel == nil {
		dialog.ShowInformation("Check for Updates",
			fmt.Sprintf("You are running version %s.\nNo updates are available at this time.", VERSION), w)
		return
	}

	dialog.ShowConfirm("Update Available!",
		fmt.Sprintf("You are running version %s.\nAn update is available: %s\nWould you like us to download it now?", VERSION, rel.GetTagName()),
		w.doUpdate(rel), w)
}

func (w *windowMain) doUpdate(rel *github.RepositoryRelease) func(bool) {
	return func(ok bool) {
		if !ok {
			return
		}

		// otherwise we might block the fyne event queue...
		go func() {
			// TODO: display download progress, filename and size
			dlg := dialog.NewProgressInfinite("Downloading Update",
				fmt.Sprintf("Downloading version %s now...", rel.GetTagName()), w)
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

func (w *windowMain) genAccountList() {
	if w.accList != nil {
		objects := w.accList.Objects
		for _, o := range objects {
			w.accList.Remove(o)
		}
	} else {
		w.accList = container.NewVBox()
	}

	players, err := sc2api.GetAccountPlayers()
	if err != nil {
		golog.Errorf("GetAccountPlayers: %v", err)
		return
	}

	accounts, err := sc2utils.EnumerateAccounts(viper.GetString("replaysRoot"))
	if err != nil {
		accounts = []string{"No Accounts Found/"}
	}

	for acc, list := range toonList(accounts) {
		header := newText(acc, 1.6, true)
		header.Move(fyne.NewPos(w.ui.theme.Padding()/2, 1+w.ui.theme.Padding()/2))
		w.accList.Add(fyne.NewContainerWithoutLayout(header))
		for _, toon := range list {
			parts := strings.Split(toon, "-")

			aLabel := newText("Unknown Character", 1, false)
			for _, p := range players {
				if parts[len(parts)-1] == strconv.Itoa(int(p.Player.CharacterID)) {
					aLabel.Text = p.Player.Name
				}
			}

			toggleBtn := widget.NewButtonWithIcon("", theme.MediaPauseIcon(), nil)
			toggleBtn.Importance = widget.HighImportance

			id := fmt.Sprintf("%s/%s", acc, toon)
			toggleBtn.OnTapped = w.toggleUploading(toggleBtn, id)
			w.uploadEnabled[id] = true

			w.accList.Add(
				container.NewBorder(nil, nil,
					toggleBtn,
					newText(sc2utils.RegionsMap[parts[0]], .9, false),
					aLabel,
				),
			)
		}
	}
}

func (w *windowMain) genUploadList() {
	w.uploadList = widget.NewTable(
		func() (int, int) { return len(w.uploadStatus), 3 },
		func() fyne.CanvasObject {
			return newText("@@@@@@@@", 1, false)
		},
		func(tci widget.TableCellID, f fyne.CanvasObject) {
			l := f.(*canvas.Text)
			switch atom := w.uploadStatus[tci.Row]; tci.Col {
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
	w.uploadList.OnSelected = func(id widget.TableCellID) {
		if id.Row > len(w.uploadStatus)-1 {
			return // selected row that does not exist
		}

		if rid := w.uploadStatus[id.Row].ReplayID; rid != "" {
			u, _ := url.Parse(fmt.Sprintf("%s/replay/%s", sc2replaystats.WebRoot, rid))
			w.app.OpenURL(u)
		}
	}
	w.uploadList.SetColumnWidth(0, 230)
	w.uploadList.SetColumnWidth(1, 76)
	w.uploadList.SetColumnWidth(2, 84)
}

func (w *windowMain) handleReplay(replayFilename string) {
	_, mapName, _ := utils.SplitFilepath(replayFilename)
	entry := &uploadRecord{
		Filename: replayFilename,
		MapName:  mapName,
		Status:   "pending",
	}

	golog.Debugf("uploading replay: %v", replayFilename)
	w.uploadStatus = append(w.uploadStatus, entry)
	w.uploadList.Refresh()

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
	w.uploadList.Refresh()

	rqid, err := sc2api.UploadReplay(replayFilename)
	entry.QueueID = rqid

	if err != nil {
		entry.Status = "u failed"
		w.uploadList.Refresh()

		dialog.NewError(fmt.Errorf("replay upload failed:%v\n%v", mapName, err), w)
		return
	}

	entry.Status = "processing"
	w.uploadList.Refresh()

	go w.watchReplayStatus(entry)
}

func (w *windowMain) openGettingStarted1() {
	w.gettingStarted = 1
	w.WizardModal("Skip", "Next", nil, func() {
		if viper.GetString("replaysroot") == "" {
			w.openGettingStarted2()
		} else {
			w.gettingStarted = 0
			w.modal.Hide()
		}
	},
		labelWithWrapping("You are only two steps away from having your replays automatically uploaded!"),
		labelWithWrapping("1) We will find your Replays Directory"),
		labelWithWrapping("2) We will find your sc2replaystats API Key"),
	)
}

func (w *windowMain) openGettingStarted2() {
	w.gettingStarted = 2

	btnSettings := widget.NewButtonWithIcon("Open Settings", theme.SettingsIcon(), func() {
		w.ui.OpenWindow(WindowSettings)
	})
	btnSettings.Importance = widget.HighImportance

	// TODO: Refactor this to actually have the settings UI, not just direct the user to settings
	w.WizardModal("", "", nil, nil,
		labelWithWrapping("First thing's first. Please use the button below to open the Settings dialog, and under the StarCraft II section, add your Replays Directory."),
		btnSettings,
		labelWithWrapping("Once you have found your replays directory and saved the settings, this setup wizard will automatically advance to the next step."),
	)
}

func (w *windowMain) openGettingStarted3() {
	w.gettingStarted = 3

	btnSettings := widget.NewButtonWithIcon("Open Settings", theme.SettingsIcon(), func() {
		w.ui.OpenWindow(WindowSettings)
	})
	btnSettings.Importance = widget.HighImportance

	// TODO: Refactor this to actually have the settings UI, not just direct the user to settings
	w.WizardModal("", "", nil, nil,
		labelWithWrapping("Lastly, please set your sc2replaystats API key. If you do not know how to find this, use the \"Login and find it for me\" button to have us login to your account and generate one on your behalf."),
		btnSettings,
	)
}

func (w *windowMain) openGettingStarted4() {
	w.gettingStarted = 0

	w.WizardModal("Close", "", func() {
		w.gettingStarted = 0
		w.modal.Hide()
	}, nil,
		labelWithWrapping("Contratulations! You have finished first-time setup. You can change these settings at any time by going to File -> Settings."),
	)
}

func (w *windowMain) setupUploader() {
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
	if w.watcher != nil {
		w.watcher.Close()
	}

	watch, err := newWatcher(paths)
	if err != nil {
		dialog.NewError(fmt.Errorf("Failed to start uploader:\n%v", err), w)
		return
	}
	w.watcher = watch

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
						go w.handleReplay(event.Name)
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

func (w *windowMain) toggleUploading(btn *widget.Button, id string) func() {
	return func() {
		replaysRoot := viper.GetString("replaysRoot")

		w.uploadEnabled[id] = !w.uploadEnabled[id]
		if w.uploadEnabled[id] {
			if err := w.watcher.Remove(filepath.Join(replaysRoot, id, "Replays", "Multiplayer")); err != nil {
				dialog.NewError(err, w)
				return
			}

			btn.Importance = widget.HighImportance
			btn.Icon = theme.MediaPauseIcon()
		} else {
			if err := w.watcher.Add(filepath.Join(replaysRoot, id, "Replays", "Multiplayer")); err != nil {
				dialog.NewError(err, w)
				return
			}

			btn.Importance = widget.MediumImportance
			btn.Icon = theme.MediaPlayIcon()
		}
	}
}

func (w *windowMain) watchReplayStatus(entry *uploadRecord) {
	for {
		time.Sleep(time.Second)
		rid, err := sc2api.GetReplayStatus(entry.QueueID)
		if err != nil {
			golog.Errorf("error checking reply status: %v: %v", entry.QueueID, err)
			entry.Status = "p failed"
			w.uploadList.Refresh()
			return // could not check status
		}

		if rid != "" {
			entry.Status = "success"
			entry.ReplayID = rid
			w.uploadList.Refresh()
			return // replay parsed!
		}

		golog.Debugf("sc2replaystats process..: [%v] %s", entry.QueueID, rid)
	}
}

func toonList(accounts []string) (toons map[string][]string) {
	toons = make(map[string][]string)
	for _, acc := range accounts {
		parts := strings.Split(acc[1:], string(filepath.Separator))
		toonList, ok := toons[parts[0]]
		if !ok {
			toons[parts[0]] = []string{parts[1]}
		} else {
			toons[parts[0]] = append(toonList, parts[1])
		}
	}

	return toons
}
