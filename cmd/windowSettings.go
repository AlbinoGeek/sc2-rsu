package cmd

import (
	"errors"
	"fmt"
	"image/color"
	"os"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/AlbinoGeek/sc2-rsu/sc2replaystats"
	"github.com/AlbinoGeek/sc2-rsu/sc2utils"
	"github.com/spf13/viper"
)

type windowSettings struct {
	*windowBase

	// do we have unsaved changes in the form?
	unsaved bool

	// widgets
	replaysRoot  *widget.Entry
	updatePeriod *widget.Entry
}

// TODO: candidate for refactor
func (w *windowSettings) Init() {
	w.windowBase.Window = w.windowBase.app.NewWindow("Settings")

	// widgets we need to access from other funcs
	w.replaysRoot = widget.NewEntry()
	w.replaysRoot.SetText(viper.GetString("replaysRoot"))

	w.updatePeriod = widget.NewEntry()
	w.updatePeriod.SetText(getUpdateDuration().String())
	w.updatePeriod.Validator = func(period string) (err error) {
		_, err = time.ParseDuration(period)
		return
	}

	// widgets we don't need to access
	autoDownload := widget.NewCheck("Automatically Download Updates?", func(checked bool) {
		w.unsaved = true
		viper.Set("update.automatic.enabled", checked)
	})
	autoDownload.SetChecked(viper.GetBool("update.automatic.enabled"))

	checkUpdates := widget.NewCheck("Check for Updates Periodically?", func(checked bool) {
		w.unsaved = true
		if checked {
			autoDownload.Enable()
			w.updatePeriod.Enable()
		} else {
			autoDownload.Disable()
			w.updatePeriod.Disable()
		}
		viper.Set("update.check.enabled", checked)
	})
	checkUpdates.SetChecked(viper.GetBool("update.check.enabled"))
	w.unsaved = false // otherwise set by the above line

	if !checkUpdates.Checked {
		autoDownload.Disable()
		w.updatePeriod.Disable()
	}

	apiKey := widget.NewEntry()
	apiKey.SetText(viper.GetString("apiKey"))
	apiKey.Validator = func(key string) (err error) {
		if !sc2replaystats.ValidAPIKey(key) {
			err = errors.New("invalid API key format")
		}
		return
	}

	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(5, 5))
	w.SetContent(widget.NewVBox(
		widget.NewCard(fmt.Sprintf("%s Settings", PROGRAM), "", widget.NewVBox(
			checkUpdates,
			autoDownload,
			fyne.NewContainerWithLayout(
				layout.NewFormLayout(),
				widget.NewLabel("Check Every"),
				w.updatePeriod,
			),
		)),
		spacer,
		widget.NewCard("sc2ReplayStats Account", "", widget.NewVBox(
			fyne.NewContainerWithLayout(
				layout.NewFormLayout(),
				widget.NewLabel("API Key"),
				widget.NewHScrollContainer(apiKey),
			),
			widget.NewButtonWithIcon("Login and Generate it for me...", theme.ComputerIcon(), func() {
				// ! IMPLEMENT LOGIN FORM
			}),
		)),
		spacer,
		widget.NewCard("StarCraft II", "", widget.NewVBox(
			fyne.NewContainerWithLayout(
				layout.NewFormLayout(),
				widget.NewLabel("Replays Root"),
				widget.NewHScrollContainer(w.replaysRoot),
			),
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(2),
				widget.NewButtonWithIcon("Find it for me...", theme.SearchIcon(), func() { go w.findReplaysRoot() }),
				widget.NewButtonWithIcon("Browse...", theme.FolderOpenIcon(), func() {
					dlg := dialog.NewFolderOpen(w.browseReplaysRoot, w)
					dlg.Resize(w.Canvas().Size().Subtract(fyne.NewSize(20, 20))) // ! can't be larger than the settings window
					dlg.Show()
				}),
			),
		)),
		spacer,
		layout.NewSpacer(),
		widget.NewSeparator(),
		fyne.NewContainerWithLayout(
			layout.NewGridLayout(2),
			widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
				w.unsaved = false
				loadConfig()
				w.onClose()
			}),
			widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
				if err := saveConfig(); err != nil {
					dialog.ShowError(err, w)
				} else {
					dialog.ShowInformation("Saved!", "Your settings have been saved.", w.ui.windows[WindowMain].GetWindow())
					w.unsaved = false
					w.Close()
				}
			}),
		),
	))

	w.SetCloseIntercept(w.onClose)
	w.SetOnClosed(func() {
		w.SetWindow(nil)
	})

	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			w.onClose()
		}
	})

	w.SetPadded(false)
	w.SetFixedSize(true)

	w.Resize(fyne.NewSize(600, 600))
	w.CenterOnScreen()
	w.Show()
}

// TODO: candidate for refactor
func (w *windowSettings) browseReplaysRoot(uri fyne.ListableURI, err error) {
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	if uri == nil {
		return // cancelled
	}

	root := uri.String()
	if strings.HasPrefix(root, "file://") {
		root = root[7:] // ? is this reliable
	}

	// TODO: record the newly found accounts if confirmed
	if accs, err := sc2utils.EnumerateAccounts(root); err != nil || len(accs) == 0 {
		dialog.ShowConfirm("Invalid Directory!",
			fmt.Sprintf("We could not find any accounts in that directory.\nAre you sure you want to use it anyways?\n\n%s", root),
			func(ok bool) {
				if ok {
					w.unsaved = true
					viper.Set("replaysRoot", root)
					w.replaysRoot.SetText(root)
				}
			}, w)
		return
	}

	w.unsaved = true
	viper.Set("replaysRoot", root)
	w.replaysRoot.SetText(root)
}

// TODO: candidate for refactor
func (w *windowSettings) findReplaysRoot() {
	scanRoot := "/"
	if home, err := os.UserHomeDir(); err == nil {
		scanRoot = home
	}

	dlg := dialog.NewProgressInfinite("Searching for Replays Root...",
		"Please wait while we search for a valid Replays folder.\nThis could take several minutes.", w)
	dlg.Show()
	roots, err := sc2utils.FindReplaysRoot(scanRoot)
	dlg.Hide()

	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	if len(roots) == 0 {
		dialog.ShowError(errors.New("no replay directories found"), w)
		return
	}

	if len(roots) == 1 {
		accs, err := sc2utils.EnumerateAccounts(roots[0])
		if err != nil {
			dialog.ShowError(fmt.Errorf("error scanning for accounts: %v", err), w)
			return
		}

		dialog.ShowInformation("Replays Root Found!",
			fmt.Sprintf("We found your replays directory!\nIt contains %d account/toons.\n%s", len(accs), roots[0]), w)
		w.replaysRoot.SetText(roots[0])
		return
	}

	selected := -1
	longest := ""
	for _, s := range roots {
		if l := len(s); l > len(longest) {
			longest = s
		}
	}

	listWidget := widget.NewList(func() int {
		return len(roots)
	}, func() fyne.CanvasObject {
		return widget.NewLabel(longest)
	}, func(id int, obj fyne.CanvasObject) {
		obj.(*widget.Label).SetText(roots[id])
	})
	listWidget.OnSelected = func(id int) {
		selected = id
	}
	dlg2 := dialog.NewCustomConfirm("Multiple Possible Roots Found",
		"Select", "Cancel", widget.NewHScrollContainer(listWidget), func(ok bool) {
			if !ok {
				return
			}
			viper.Set("replaysRoot", roots[selected])
			w.replaysRoot.SetText(roots[selected])
		}, w)

	size := fyne.MeasureText(longest, theme.TextSize(), fyne.TextStyle{})
	size.Height *= len(roots)

	dlg2.Resize(fyne.NewSize(60, 144).Add(size))
	dlg2.Show()
}

func (w *windowSettings) onClose() {
	if !w.unsaved {
		w.Close()
		return
	}

	dialog.ShowConfirm("Unsaved Changes",
		"You have not saved your settings.\nAre you sure you want to discard amy changes?",
		func(ok bool) {
			if ok {
				w.Close()
			}
		}, w)
}
