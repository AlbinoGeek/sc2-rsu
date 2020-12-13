package fynemd

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

var (
	// IconSize used with (in display pixels) used with Material Design elements
	IconSize = 24

	// Padding (in display pixels) used with Material Design elements
	Padding = 16
)

// AppBar ...
type AppBar struct {
	widget.BaseWidget

	Dense    bool
	Extended bool
	Title    string

	actions []*widget.Button
	nav     *NavigationDrawer
	objects []fyne.CanvasObject
}

// NewAppBar ...
func NewAppBar(title string) *AppBar {
	bar := &AppBar{
		Title: title,
	}

	bar.ExtendBaseWidget(bar)

	return bar
}

// SetDense ...
func (bar *AppBar) SetDense(dense bool) {
	bar.Dense = dense
	bar.Refresh()
}

// SetExtended ...
func (bar *AppBar) SetExtended(extended bool) {
	bar.Extended = extended
	bar.Refresh()
}

// SetTitle ...
func (bar *AppBar) SetTitle(title string) {
	bar.Title = title
	bar.Refresh()
}

// SetNavigation ...
func (bar *AppBar) SetNavigation(nav *NavigationDrawer) {
	bar.nav = nav
	bar.Refresh()
}

// CreateRenderer returns a new WidgetRenderer for this widget.
// This should not be called by regular code, it is used internally to render a widget.
//
// Implements: fyne.Widget
func (bar *AppBar) CreateRenderer() fyne.WidgetRenderer {
	rend := &appBarRenderer{
		bar: bar,
	}

	rend.Init()

	return rend
}

// --

// appBarRenderer defines the behaviour of a widget's implementation.
// This is returned from a widget's declarative object through the CreateRenderer()
// function and should be exactly one instance per widget in memory.
//
// Implements: fyne.WidgetRenderer
type appBarRenderer struct {
	bar *AppBar

	navIcon *widget.Button
	title   *canvas.Text
}

// BackgroundColor returns the color that should be used to draw the background of this renderer’s widget.
//
// Deprecated: Widgets will no longer have a background to support hover and selection indication in collection widgets.
// If a widget requires a background color or image, this can be achieved by using a canvas.Rect or canvas.Image
// as the first child of a MaxLayout, followed by the rest of the widget components.
func (*appBarRenderer) BackgroundColor() color.Color {
	return theme.PrimaryColor()
}

// Destroy is for internal use.
func (*appBarRenderer) Destroy() {}

func (br *appBarRenderer) Init() {
	br.navIcon = widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
		if br.bar.nav != nil {
			br.bar.nav.Show()
		}
	})
	br.title = NewScaledText(TextSizeHeading5, br.bar.Title)
	// br.title.TextStyle.Bold = true
	br.bar.objects = []fyne.CanvasObject{
		br.navIcon, br.title,
	}
}

// Layout is a hook that is called if the widget needs to be laid out.
// This should never call Refresh.
// ! should respond to theme values
func (br *appBarRenderer) Layout(space fyne.Size) {
	pos := fyne.NewPos(Padding, Padding)

	if br.bar.Dense {
		pos.Y = Padding / 2
	}

	if br.navIcon.Visible() {
		br.navIcon.Move(pos)
		br.navIcon.Resize(fyne.NewSize(IconSize, IconSize))
	}

	pos.Y -= Padding / 4

	if br.bar.nav != nil && !br.bar.nav.Visible() {
		pos.X += br.navIcon.Size().Width + Padding*2
		br.navIcon.Show()
	} else {
		br.navIcon.Hide()
	}

	br.title.Move(pos)
	br.title.Resize(br.title.MinSize())

	if br.bar.Title == "" && br.title.Visible() {
		br.title.Hide()
	}

	// TODO: Layout actions from right

	// TODO: if len(actions) > 3 { actionsMenu = actions[2:] } ...
}

// MinSize returns the minimum size of the widget that is rendered by this renderer.
func (br *appBarRenderer) MinSize() fyne.Size {
	// pad := theme.Padding()

	size := fyne.NewSize(360, 56) // material specs

	if br.bar.Dense {
		size.Height = 40
	} else if br.bar.Extended {
		size.Height = 128
	}

	// TODO: enough space for navigationIcon if visible

	// enough space for the text
	size = size.Max(fyne.MeasureText(br.bar.Title, theme.TextSize(), br.bar.nav.title.TextStyle))

	// TODO: enough space for action buttons when implemented

	return size
}

// Objects returns all objects that should be drawn.
func (br *appBarRenderer) Objects() []fyne.CanvasObject {
	return br.bar.objects
}

// Refresh is a hook that is called if the widget has updated and needs to be redrawn.
// This might trigger a Layout.
func (br *appBarRenderer) Refresh() {
	for _, o := range br.Objects() {
		// ? should eventually look into why these are nil
		if o == nil || !o.Visible() {
			continue
		}

		o.Refresh()
	}
}
