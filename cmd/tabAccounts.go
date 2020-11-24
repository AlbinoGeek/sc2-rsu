package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/kataras/golog"
	"github.com/spf13/viper"

	"github.com/AlbinoGeek/sc2-rsu/sc2utils"
)

func (main *windowMain) genAccountList() {
	if main.accList != nil {
		objects := main.accList.Objects
		for _, o := range objects {
			main.accList.Remove(o)
		}
	} else {
		main.accList = container.NewVBox()
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
		header := newHeader(acc)
		header.Move(fyne.NewPos(main.UI.Theme.Padding()/2, 1+main.UI.Theme.Padding()/2))
		main.accList.Add(fyne.NewContainerWithoutLayout(header))
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
			toggleBtn.OnTapped = main.toggleUploading(toggleBtn, id)
			main.uploadEnabled[id] = true

			main.accList.Add(
				container.NewBorder(nil, nil,
					toggleBtn,
					newText(sc2utils.RegionsMap[parts[0]], .9, false),
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
