# BLE Device Scanner

This is a simple Bluetooth Low Energy (BLE) device scanner built using Go and the Gio UI framework. The application allows users to scan for nearby BLE devices and display their details in a graphical interface.

## Prerequisites

Before you can build and run this project, ensure you have the following installed:

1. **Go** (Golang) installed on your system. Download and install it from [golang.org](https://go.dev/dl/).
2. **Git** installed to clone repositories and manage dependencies.
3. Required Go modules and dependencies.

## Installation

### 1. Install Go
If you haven't installed Go yet, download it from [Go's official website](https://go.dev/dl/) and follow the installation instructions for your operating system.

After installation, verify Go is installed correctly by running:
```sh
 go version
```

### 2. Clone the Repository

### 3. Install Dependencies
Run the following command to install the required Go modules:
```sh
go mod tidy
```

Ensure the required dependencies are installed:
```sh
go get gioui.org/app gioui.org/layout gioui.org/op gioui.org/text gioui.org/widget gioui.org/widget/material tinygo.org/x/bluetooth
```

## Building and Running the Application

### Run the Application
To run the application locally, execute:
```sh
go run src/main.go
```

### Build the Executable
To create an executable file for Windows:
```sh
go build -o build/FINDBLE.exe src/main.go
```
For Linux/macOS:
```sh
go build -o build/FINDBLE src/main.go
```

### Running the Executable
On Windows:
```sh
./build/FINDBLE.exe
```
On Linux/macOS:
```sh
./build/FINDBLE
```

## Notes
- Ensure your system has Bluetooth enabled before running the application.
- If you're on Linux, you may need to run the application with `sudo` to access Bluetooth features.

## License
This project is licensed under the MIT License.

## Author
Created by [Your Name]. Feel free to contribute and improve this project!
