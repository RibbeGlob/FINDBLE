package main

import (
	"fmt"
	"log"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"tinygo.org/x/bluetooth"
)

var (
	adapter    = bluetooth.DefaultAdapter
	deviceList []bluetooth.ScanResult
	scanning   bool
	stopScanCh chan struct{}
	mutex      sync.Mutex
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("BLE Scanner")
	myWindow.Resize(fyne.NewSize(400, 400))

	list := widget.NewList(
		func() int { return len(deviceList) },
		func() fyne.CanvasObject { return widget.NewButton("", nil) },
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			btn := obj.(*widget.Button)
			device := deviceList[i]
			btn.SetText(fmt.Sprintf("%s [%s] RSSI: %d", device.LocalName(), device.Address.String(), device.RSSI))
			btn.OnTapped = func() {
				connectToDevice(device)
			}
		},
	)

	scanButton := widget.NewButton("Scan for BLE Devices", func() {
		if !scanning {
			go scanBLEDevices(list)
		} else {
			log.Println("Scanning is already in progress...")
		}
	})

	stopButton := widget.NewButton("Stop Scanning", func() {
		if scanning {
			close(stopScanCh)
		} else {
			log.Println("No active scan to stop.")
		}
	})

	exitButton := widget.NewButton("Exit", func() {
		myApp.Quit()
	})

	buttonContainer := container.NewVBox(scanButton, stopButton, exitButton)
	mainContainer := container.NewBorder(nil, buttonContainer, nil, nil, list)

	myWindow.SetContent(mainContainer)
	myWindow.ShowAndRun()
}

func scanBLEDevices(list *widget.List) {
	if err := adapter.Enable(); err != nil {
		log.Fatalf("Error enabling adapter: %v", err)
	}

	mutex.Lock()
	deviceList = nil
	mutex.Unlock()

	scanning = true
	stopScanCh = make(chan struct{})

	log.Println("Scanning for BLE devices...")

	go func() {
		err := adapter.Scan(func(a *bluetooth.Adapter, d bluetooth.ScanResult) {
			mutex.Lock()
			deviceList = append(deviceList, d)
			mutex.Unlock()
			list.Refresh()
		})
		if err != nil {
			log.Printf("Error during scanning: %v", err)
		}
	}()

	<-stopScanCh
	log.Println("Scanning stopped.")

	adapter.StopScan()
	scanning = false
}

func connectToDevice(device bluetooth.ScanResult) {
	log.Printf("Connecting to device %s [%s]...", device.LocalName(), device.Address.String())
	peer, err := adapter.Connect(device.Address, bluetooth.ConnectionParams{})
	if err != nil {
		log.Printf("Failed to connect to %s: %v", device.LocalName(), err)
		return
	}
	log.Printf("Connected to %s [%s]", device.LocalName(), device.Address.String())
	defer peer.Disconnect()
}
