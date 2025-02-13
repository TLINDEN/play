package main

import (
	"image/color"
	"log"

	"golang.org/x/image/font/gofont/goregular"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
)

type game struct {
	ui *ebitenui.UI
}

func newList(entries []interface{}, res *uiResources, widgetOpts ...widget.WidgetOpt) *widget.List {
	return widget.NewList(
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widgetOpts...)),
		widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(res.list.image)),
		widget.ListOpts.SliderOpts(
			widget.SliderOpts.Images(res.list.track, res.list.handle),
			widget.SliderOpts.MinHandleSize(res.list.handleSize),
			widget.SliderOpts.TrackPadding(res.list.trackPadding),
		),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(string)
		}),
		widget.ListOpts.EntryFontFace(res.list.face),
		widget.ListOpts.EntryColor(res.list.entry),
		widget.ListOpts.EntryTextPadding(res.list.entryPadding),
	)
}

func main() {
	ebiten.SetWindowSize(900, 800)
	ebiten.SetWindowTitle("Ebiten UI Hello World")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	res, err := newUIResources()
	if err != nil {
		log.Fatal(err)
	}

	// This creates the root container for this UI.
	// All other UI elements must be added to this container.
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			// It is using a GridLayout with a single column
			widget.GridLayoutOpts.Columns(1),
			// It uses the Stretch parameter to define how the rows will be layed out.
			// - a fixed sized header
			// - a content row that stretches to fill all remaining space
			// - a fixed sized footer
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			// Padding defines how much space to put around the outside of the grid.
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    50,
				Bottom: 50,
			}),
			// Spacing defines how much space to put between each column and row
			widget.GridLayoutOpts.Spacing(0, 20))),
	)

	listsContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: false,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Stretch([]bool{false, false, true}, []bool{true}),
			widget.GridLayoutOpts.Spacing(10, 0))))

	entries1 := []interface{}{"One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten"}
	list1 := newList(entries1, res, widget.WidgetOpts.LayoutData(widget.GridLayoutData{
		MaxHeight: 400,
	}))

	list2 := newList(entries1, res, widget.WidgetOpts.LayoutData(widget.GridLayoutData{
		MaxHeight: 400,
	}))
	listsContainer.AddChild(list1)
	listsContainer.AddChild(list2)

	// This loads a font and creates a font face.
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal("Error Parsing Font", err)
	}
	fontFace := truetype.NewFace(ttfFont, &truetype.Options{
		Size: 32,
	})

	// This creates a text widget that says "Hello World!"
	helloWorldLabel := widget.NewText(
		widget.TextOpts.Text("Hello World!", fontFace, color.White),
	)

	// To display the text widget, we have to add it to the root container.
	rootContainer.AddChild(helloWorldLabel)

	rootContainer.AddChild(listsContainer)

	// This adds the root container to the UI, so that it will be rendered.
	eui := &ebitenui.UI{
		Container: rootContainer,
	}

	game := game{
		ui: eui,
	}

	err = ebiten.RunGame(&game)
	if err != nil {
		log.Print(err)
	}
}

func (g *game) Update() error {
	// ui.Update() must be called in ebiten Update function, to handle user input and other things
	g.ui.Update()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	// ui.Draw() should be called in the ebiten Draw function, to draw the UI onto the screen.
	// It should also be called after all other rendering for your game so that it shows up on top of your game world.
	g.ui.Draw(screen)
}

func (g *game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
