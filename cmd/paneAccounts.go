package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/kataras/golog"
	"github.com/spf13/viper"

	"github.com/AlbinoGeek/sc2-rsu/cmd/gui"
	"github.com/AlbinoGeek/sc2-rsu/fynemd"
	"github.com/AlbinoGeek/sc2-rsu/sc2utils"
)

type paneAccounts struct {
	fynemd.Pane

	container *fyne.Container
}

func makePaneAccounts(w gui.Window) fynemd.Pane {
	p := &paneAccounts{
		Pane: fynemd.NewPaneWithIcon("Accounts", accIcon, w),
	}

	p.container = container.NewVBox()
	p.SetContent(container.NewVScroll(p.container))

	go p.Init() // takes hundreds of ms

	return p
}

func (t *paneAccounts) Init() {
	t.Update()
	t.container.Refresh()
}

func (t *paneAccounts) Refresh() {
	t.container.Refresh()
}

func (t *paneAccounts) Update() {
	players, err := sc2api.GetAccountPlayers()

	if err != nil {
		golog.Errorf("GetAccountPlayers: %v", err)
		return
	}

	accounts, err := sc2utils.EnumerateAccounts(viper.GetString("replaysRoot"))
	if err != nil {
		accounts = []string{"No Accounts Found/"}
	}

	// Clear container if it has objects
	objects := t.container.Objects

	for _, o := range objects {
		t.container.Remove(o)
	}

	objects = nil

	main := t.GetWindow().(*windowMain)

	for acc, list := range toonList(accounts) {
		for _, toon := range list {
			name := ""

			// find toon name via sc2replaystats account players
			parts := strings.Split(toon, "-")
			for _, p := range players {
				if parts[len(parts)-1] == strconv.Itoa(int(p.Player.CharacterID)) {
					name = p.Player.Name
				}
			}

			card := widget.NewCard(name, sc2utils.RegionsMap[parts[0]], nil)
			id := fmt.Sprintf("%s/%s", acc, toon)

			btnToggle := widget.NewButtonWithIcon("", theme.MediaPauseIcon(), nil)
			btnToggle.Importance = widget.HighImportance
			btnToggle.OnTapped = main.toggleUploading(btnToggle, id)

			// todo: reverse this map ( disableUpload )
			main.uploadEnabled[id] = true

			if !getToonEnabled(id) {
				if err := main.watcher.Remove(filepath.Join(viper.GetString("replaysRoot"), id, "Replays", "Multiplayer")); err != nil {
					dialog.NewError(err, t.GetWindow().GetWindow())

					return
				}

				btnToggle.Importance = widget.MediumImportance
				btnToggle.Icon = theme.MediaPlayIcon()
			}

			t.container.Add(container.NewBorder(nil, nil, btnToggle, nil, card))
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
