package cmd

import (
	"fmt"
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/AlbinoGeek/sc2-rsu/cmd/gui"
	"github.com/AlbinoGeek/sc2-rsu/cmd/gui/fynewidget"
	"github.com/AlbinoGeek/sc2-rsu/sc2replaystats"
)

type paneUploads struct {
	fynewidget.Pane

	table *widget.Table
}

func makePaneUploads(w gui.Window) fynewidget.Pane {
	p := &paneUploads{
		Pane: fynewidget.NewPaneWithIcon("", uploadIcon, w),
	}

	p.Init()
	return p
}

func (t *paneUploads) Init() {
	main := t.GetWindow().(*windowMain)

	t.table = widget.NewTable(
		func() (int, int) { return len(main.uploadStatus), 3 },
		func() fyne.CanvasObject {
			return fynewidget.NewText("@@@@@@@@", 1, false)
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
	t.table.OnSelected = func(id widget.TableCellID) {
		if id.Row > len(main.uploadStatus)-1 {
			return // selected row that does not exist
		}

		if rid := main.uploadStatus[id.Row].ReplayID; rid != "" {
			u, _ := url.Parse(fmt.Sprintf("%s/replay/%s", sc2replaystats.WebRoot, rid))
			main.App.OpenURL(u)
		}
	}
	t.table.SetColumnWidth(0, 230)
	t.table.SetColumnWidth(1, 86)
	t.table.SetColumnWidth(2, 90)
	pad := theme.Padding()

	tblName := fynewidget.NewText("Map Name", 1, true)
	tblName.Move(fyne.NewPos(pad*2, 3))

	tblID := fynewidget.NewText("ID", 1, true)
	tblID.Move(fyne.NewPos(228+pad*5, 3))

	tblStatus := fynewidget.NewText("Status", 1, true)
	tblStatus.Move(fyne.NewPos(312+pad*7, 3))

	t.SetContent(container.NewBorder(
		fyne.NewContainerWithoutLayout(
			tblName, tblID, tblStatus,
		), nil, nil, nil, t.table,
	))
}

func (t *paneUploads) Refresh() {
	t.table.Refresh()
}
