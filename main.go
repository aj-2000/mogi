package main

import (
	"log"
	mogiApp "mogi/app"
	"mogi/color"
	"mogi/examples"
	"mogi/math"
	"mogi/ui"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
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
		app.Run(func(app *mogiApp.App) ui.IComponent {

			bgColor := color.Yellow
			r := app.Container().
				SetID("app_container").
				SetBackgroundColor(bgColor).
				AddChildren( // Add all children at once
					// tabs,
					// examples.ChessboardComponent(app),
					examples.BuyNowCardComponent(app),
					// examples.NestedContainersComponent(app),
					// examples.ExampleMarginPaddingBorder(app),
					// examples.ClayDemoComponent(app),
					// examples.BoxesOneComponent(app),
					// examples.BoxesNLevelComponent(app, 3, 2, 100),
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
