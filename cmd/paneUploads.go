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
	"github.com/AlbinoGeek/sc2-rsu/fynex"
	"github.com/AlbinoGeek/sc2-rsu/sc2replaystats"
)

type paneUploads struct {
	fynex.Pane

	table *widget.Table
}

func makePaneUploads(w gui.Window) fynex.Pane {
	p := &paneUploads{
		Pane: fynex.NewPaneWithIcon("Uploads", uploadIcon, w),
	}

	p.Init()

	return p
}

func (t *paneUploads) Init() {
	main := t.GetWindow().(*windowMain)

	t.table = widget.NewTable(
		func() (int, int) { return len(main.uploadStatus), 3 },
		func() fyne.CanvasObject {
			return fynex.NewScaledText(fynex.TextSizeBody1, "@@@@@@@@")
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

	// TODO needs to be in a Layout call, in an overriden widget -_-
	t.table.SetColumnWidth(0, 230)
	t.table.SetColumnWidth(1, 86)
	t.table.SetColumnWidth(2, 90)

	pad := theme.Padding()

	tblName := fynex.NewScaledText(fynex.TextSizeSubtitle1, "Map Name")
	tblName.TextStyle.Bold = true
	tblName.Move(fyne.NewPos(pad*2, 0))

	tblID := fynex.NewScaledText(fynex.TextSizeSubtitle1, "ID")
	tblID.TextStyle.Bold = true
	tblID.Move(fyne.NewPos(228+pad*5, 0))

	tblStatus := fynex.NewScaledText(fynex.TextSizeSubtitle1, "Status")
	tblStatus.TextStyle.Bold = true
	tblStatus.Move(fyne.NewPos(312+pad*7, 0))

	t.SetContent(container.NewBorder(
		fyne.NewContainerWithoutLayout(
			tblName, tblID, tblStatus,
		), nil, nil, nil, t.table,
	))
}

func (t *paneUploads) Refresh() {
	t.table.Refresh()
}
