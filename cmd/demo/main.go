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

			_ = app.Container().
				SetID("tabs").
				SetBackgroundColor(color.Magenta).
				SetBorderRadius(5).
				SetPadding(math.Vec2f32{X: 4, Y: 4}).
				AddChildren(
					app.Container().
						SetID("tab_bar").
						// TODO: fix percent precendence
						SetWidthPercent(100).
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
			// cursorSize := float32(30)
			// cursorSize1 := float32(25)

			// mousePos := app.GetMousePos()
			rows := []ui.Row{
				{
					Cells: []string{"India", "1,450,000", "4,614", "4,730", "2,182", "5,500", "2,820", "42,000", "6,464", "India maintains a large and diverse military force, with significant investments in aircraft and armored vehicles. It continues to modernize its capabilities."},
				},
				{
					Cells: []string{"Russia", "1,014,000", "13,132", "4,182", "1,531", "4,173", "1,511", "13,000", "4,173", "Russia possesses extensive military hardware, particularly in aircraft and tanks, reflecting its historical emphasis on armored and air superiority."},
				},
				{
					Cells: []string{"China", "2,035,000", "5,250", "3,285", "1,250", "3,285", "1,250", "5,250", "3,285", "China's military is rapidly expanding and modernizing, with a strong focus on air power and technological advancements in defense."},
				},
				{
					Cells: []string{"USA", "1,390,000", "13,300", "6,287", "5,550", "6,287", "5,550", "13,300", "6,287", "The USA maintains one of the most technologically advanced and powerful militaries globally, with extensive air and armored capabilities."},
				},
				{
					Cells: []string{"Pakistan", "654,000", "1,500", "1,200", "1,200", "1,200", "1,200", "2,000", "1,200", "Pakistan's military is focused on regional defense and has a significant number of aircraft and armored vehicles."},
				},
			}

			table := app.Table().
				SetID("table").
				SetHeader([]string{
					"Country",
					"Active Personnel",
					"Total Aircraft",
					"Fighter Aircraft",
					"Attack Aircraft",
					"Helicopters",
					"Attack Helicopters",
					"Combat Tanks",
					"Self-Propelled Artillery",
					"Summary",
				}).
				SetBackgroundColor(color.Red).
				AddRows(rows)
			r := app.Container().
				SetID("app_container").
				SetBackgroundColor(bgColor).
				AddChildren( // Add all children at once
					// examples.FPSCounterComponent(app),
					// app.Container().
					// 	SetID("mouse_pos").
					// 	SetZIndex(10001).
					// 	SetBackgroundColor(color.Orange).
					// 	SetBorderRadius(cursorSize/2).
					// 	SetSize(math.Vec2f32{X: cursorSize, Y: cursorSize}).
					// 	SetPosition(ui.Position{X: mousePos.X - (cursorSize / 2), Y: mousePos.Y - (cursorSize / 2), Type: ui.PositionTypeAbsolute}),
					// app.Container().
					// 	SetID("mouse_pos1").
					// 	SetZIndex(10001).
					// 	SetBackgroundColor(color.Cyan).
					// 	SetBorderRadius(cursorSize1/2).
					// 	SetSize(math.Vec2f32{X: cursorSize1, Y: cursorSize1}).
					// 	SetPosition(ui.Position{X: mousePos.X - (cursorSize1 / 2), Y: mousePos.Y - (cursorSize1 / 2), Type: ui.PositionTypeAbsolute}),
					// tabs,
					table,
					// app.Text("Lorem Ipsum is simply dummy text of the printing and typesetting industry.").
					// 	SetColor(color.Red).
					// 	SetTextWrapped(true),

				// table,
				).
				SetMargin(math.Vec2f32{X: 3, Y: 3}).
				SetPadding(math.Vec2f32{X: 4, Y: 4})
				// SetSize(math.Vec2f32{X: 120, Y: 200})

			return r
		})
	}()
	// should i rename Vec2f32 to Vec2 and Vec2 to Vec2f64?
	wg.Wait()
}

// TakeFullWidth
// TakeFullHeight
// LinuxVM
