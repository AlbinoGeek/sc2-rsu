package cmd

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/kataras/golog"
	"github.com/spf13/viper"

	"github.com/AlbinoGeek/sc2-rsu/cmd/gui"
	"github.com/AlbinoGeek/sc2-rsu/cmd/gui/fynewidget"
	"github.com/AlbinoGeek/sc2-rsu/sc2utils"
)

type tabAccounts struct {
	*gui.TabBase

	container *fyne.Container
}

func makeTabAccounts(w gui.Window) gui.Tab {
	tab := &tabAccounts{
		TabBase: gui.NewTabWithIcon("", accIcon, w),
	}

	tab.Init()
	tab.Refresh()
	return tab
}

func (t *tabAccounts) Init() {
	t.container = container.NewVBox()
	t.SetContent(container.NewVScroll(t.container))
}

func (t *tabAccounts) Refresh() {
	// Clear container if it has objects
	objects := t.container.Objects
	for _, o := range objects {
		t.container.Remove(o)
	}

	main := t.GetWindow().(*windowMain)

	players, err := sc2api.GetAccountPlayers()
	if err != nil {
		golog.Errorf("GetAccountPlayers: %v", err)
		return
	}

	accounts, err := sc2utils.EnumerateAccounts(viper.GetString("replaysRoot"))
	if err != nil {
		accounts = []string{"No Accounts Found/"}
	}

	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(main.UI.Theme.Padding(), main.UI.Theme.Padding()))

	for acc, list := range toonList(accounts) {
		header := fynewidget.NewHeader(acc)
		header.Move(fyne.NewPos(main.UI.Theme.Padding()/2, 0))
		t.container.Add(fyne.NewContainerWithoutLayout(header))
		for _, toon := range list {
			parts := strings.Split(toon, "-")

			aLabel := fynewidget.NewText("Unknown Character", 1, false)
			for _, p := range players {
				if parts[len(parts)-1] == strconv.Itoa(int(p.Player.CharacterID)) {
					aLabel.Text = p.Player.Name
				}
			}

			toggleBtn := widget.NewButtonWithIcon("", theme.MediaPauseIcon(), nil)
			toggleBtn.Importance = widget.HighImportance

			id := fmt.Sprintf("%s/%s", acc, toon)
			toggleBtn.OnTapped = main.toggleUploading(toggleBtn, id)
			main.uploadEnabled[id] = true

			t.container.Add(
				container.NewBorder(nil, nil,
					toggleBtn,
					container.NewHBox(fynewidget.NewText(sc2utils.RegionsMap[parts[0]], .9, false), spacer),
					aLabel,
				),
			)
		}
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
