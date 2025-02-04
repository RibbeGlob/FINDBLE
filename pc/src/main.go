package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"tinygo.org/x/bluetooth"
)

var (
	adapter = bluetooth.DefaultAdapter

	deviceList []string

	scanButton widget.Clickable
	stopButton widget.Clickable

	scanning   bool
	stopScanCh chan struct{}

	listWidget = &widget.List{
		List: layout.List{Axis: layout.Vertical},
	}
	scrollToTop bool
)

func main() {
	go func() {
		w := new(app.Window)
		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func scanBLEDevices(win *app.Window) {
	if err := adapter.Enable(); err != nil {
		log.Fatalf("Error during startup %v", err)
	}

	deviceList = nil
	scanning = true

	stopScanCh = make(chan struct{})

	go func() {
		err := adapter.Scan(func(a *bluetooth.Adapter, d bluetooth.ScanResult) {
			deviceInfo := fmt.Sprintf("Device: %s [%s] RSSI: %d",
				d.LocalName(), d.Address.String(), d.RSSI)

			deviceList = append([]string{deviceInfo}, deviceList...)

			fmt.Println(deviceInfo)

			scrollToTop = true
			win.Invalidate()
		})
		if err != nil {
			log.Printf("Error during scanning %v", err)
		}
	}()

	select {
	case <-stopScanCh:
		fmt.Println("Użytkownik przerwał skanowanie.")
	}

	adapter.StopScan()
	scanning = false

	win.Invalidate()
}

func run(window *app.Window) error {
	theme := material.NewTheme()
	var ops op.Ops

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					title := material.H1(theme, "BLE Device Scanner")
					title.Color = color.NRGBA{R: 127, G: 0, B: 0, A: 255}
					title.Alignment = text.Middle
					return title.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:    layout.Horizontal,
						Spacing: layout.SpaceEvenly,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if scanButton.Clicked(gtx) && !scanning {
								go scanBLEDevices(window)
							}
							return material.Button(theme, &scanButton, "Scan BLE Devices").Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if stopButton.Clicked(gtx) && scanning {
								close(stopScanCh)
							}
							return material.Button(theme, &stopButton, "Stop Scanning").Layout(gtx)
						}),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Spacer{Height: 20}.Layout(gtx)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					if scrollToTop && len(deviceList) > 0 {
						listWidget.Position.First = 0
						listWidget.Position.Offset = 0
						scrollToTop = false
					}
					return material.List(theme, listWidget).Layout(gtx, len(deviceList),
						func(gtx layout.Context, i int) layout.Dimensions {
							return material.Body1(theme, deviceList[i]).Layout(gtx)
						},
					)
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}
