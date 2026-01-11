//go:build windows

package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

const serviceName = "mma"
const serviceDisplayName = "Modbus Memory Appliance"

func handleCLI() {
	// 1️⃣ If launched by SCM → run as service
	isService, err := svc.IsWindowsService()
	if err == nil && isService {
		_ = svc.Run(serviceName, &mmaService{})
		os.Exit(0) // 🔑 MUST exit
	}

	// 2️⃣ If NOT a service command → return to normal flow
	if len(os.Args) < 2 || os.Args[1] != "service" {
		return
	}

	// 3️⃣ Handle service subcommands
	if len(os.Args) < 3 {
		fmt.Println("Usage: mma service [install|uninstall|start|stop]")
		os.Exit(1)
	}

	cmd := os.Args[2]

	m, err := mgr.Connect()
	if err != nil {
		fmt.Println("Failed to connect to service manager:", err)
		os.Exit(1)
	}
	defer m.Disconnect()

	switch cmd {
	case "install":
		exe, _ := os.Executable()
		s, err := m.CreateService(
			serviceName,
			exe,
			mgr.Config{
				DisplayName: serviceDisplayName,
			},
		)
		if err != nil {
			fmt.Println("Install failed:", err)
			os.Exit(1)
		}
		s.Close()
		fmt.Println("Service installed")

	case "uninstall":
		s, err := m.OpenService(serviceName)
		if err != nil {
			fmt.Println("Service not found")
			os.Exit(1)
		}
		s.Delete()
		s.Close()
		fmt.Println("Service uninstalled")

	case "start":
		s, err := m.OpenService(serviceName)
		if err != nil {
			fmt.Println("Service not found")
			os.Exit(1)
		}
		s.Start()
		s.Close()
		fmt.Println("Service started")

	case "stop":
		s, err := m.OpenService(serviceName)
		if err != nil {
			fmt.Println("Service not found")
			os.Exit(1)
		}
		s.Control(svc.Stop)
		s.Close()
		fmt.Println("Service stopped")

	default:
		fmt.Println("Unknown service command:", cmd)
		os.Exit(1)
	}

	// 🔑 THIS IS THE CRITICAL LINE
	os.Exit(0)
}
