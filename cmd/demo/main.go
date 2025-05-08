package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"

	mogiApp "github.com/aj-2000/mogi/app"
	"github.com/aj-2000/mogi/cmd/examples"
	"github.com/aj-2000/mogi/color"
	"github.com/aj-2000/mogi/internal/ui"
	"github.com/aj-2000/mogi/math"
)

func main() {

	// --- pprof setup ---
	fCPU, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatalf("could not create CPU profile: %v", err)
	}
	defer fCPU.Close()
	if err := pprof.StartCPUProfile(fCPU); err != nil {
		log.Fatalf("could not start CPU profile: %v", err)
	}
	defer pprof.StopCPUProfile()

	fMem, err := os.Create("mem.pprof")
	if err != nil {
		log.Fatalf("could not create memory profile: %v", err)
	}
	defer func() {
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(fMem); err != nil {
			log.Fatalf("could not write memory profile: %v", err)
		}
		fMem.Close()
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		app := mogiApp.NewApp(800, 800, "Mogi")
		if app == nil {
			log.Fatalln("Failed to create app")
		}
		app.SetVSync(false)
		// TODO: is it needed?
		defer app.Destroy()

		app.LoadFont("JetBrainsMonoNL-Regular.ttf", 24.0)
		// TODO: fix app is not a type issue
		flag := false

		app.Run(func(app *mogiApp.App) ui.IComponent {
			children := []ui.IComponent{
				examples.ChessboardComponent(app),
				examples.BuyNowCardComponent(app),
				examples.BoxesOneComponent(app),
				// examples.BoxesNLevelComponent(app, 3, 3, 100),
				examples.NestedContainersComponent(app),
				// examples.ClayDemoComponent(app),
				examples.ExampleMarginPaddingBorder(app),
			}

			setTabIndex := func(index int) {
				for i, comp := range children {
					display := ui.DisplayNone
					if i == index {
						display = ui.DisplayBlock
					}
					switch c := comp.(type) {
					case *ui.Button:
						c.SetDisplay(display)
					case *ui.Text:
						c.SetDisplay(display)
					case *ui.Image:
						c.SetDisplay(display)
					case *ui.Container:
						c.SetDisplay(display)
					default:
						log.Printf("Unknown component type: %T", c)

					}

				}
			}
			if !flag {
				setTabIndex(0)
				flag = true
			}

			tabs := app.Container().
				SetID("tabs").
				SetBackgroundColor(color.Magenta).
				SetBorderRadius(5).
				SetPadding(math.Vec2f32{X: 4, Y: 4}).
				AddChildren(
					app.Container().
						SetID("tab_bar").
						SetZIndex(999).
						SetBackgroundColor(color.Gray).
						SetBorderRadius(5).
						SetBorder(math.Vec2f32{X: 1, Y: 1}).
						SetBorderColor(color.Black).
						SetPadding(math.Vec2f32{X: 4, Y: 4}).
						SetGap(math.Vec2f32{X: 3, Y: 3}).
						SetMargin(math.Vec2f32{X: 3, Y: 3}).
						SetDisplay(ui.DisplayBlock).
						AddChildren(
							app.Button("Chessboard").
								SetDisplay(ui.DisplayInline).
								SetOnClick(func(self *ui.Button) {
									setTabIndex(0)
								}),
							app.Button("Buy Now").
								SetDisplay(ui.DisplayInline).
								SetOnClick(func(self *ui.Button) {
									setTabIndex(1)
									log.Println("Buy Now clicked")
								}),
							app.Button("Boxes One").
								SetDisplay(ui.DisplayInline).
								SetOnClick(func(self *ui.Button) {
									setTabIndex(2)
								}),
							// app.Button("Boxes N Level").
							// 	SetDisplay(ui.DisplayInline).
							// 	SetOnClick(func(self *ui.Button) {
							// 		setTabIndex(3)
							// 	}),
							app.Button("Nested Containers").
								SetDisplay(ui.DisplayInline).
								SetOnClick(func(self *ui.Button) {
									setTabIndex(3)
								}),
							// app.Button("Clay Demo").
							// 	SetDisplay(ui.DisplayInline).
							// 	SetOnClick(func(self *ui.Button) {
							// 		setTabIndex(5)
							// 	}),
							app.Button("Margin Padding Border").
								SetDisplay(ui.DisplayInline).
								SetOnClick(func(self *ui.Button) {
									setTabIndex(4)
								}),
						),
					app.Container().
						SetID("tabs_container").
						SetBackgroundColor(color.Pink).
						AddChildren(children...),
				)

			bgColor := color.Transparent
			r := app.Container().
				SetID("app_container").
				SetBackgroundColor(bgColor).
				AddChildren( // Add all children at once
					tabs,
					examples.FPSCounterComponent(app),
				).
				SetMargin(math.Vec2f32{X: 3, Y: 3}).
				SetPadding(math.Vec2f32{X: 4, Y: 4})

			return r
		})
	}()
	// should i rename Vec2f32 to Vec2 and Vec2 to Vec2f64?
	wg.Wait()
}

// TakeFullWidth
// TakeFullHeight
// LinuxVM
