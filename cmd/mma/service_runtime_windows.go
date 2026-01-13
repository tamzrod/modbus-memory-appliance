//go:build windows

package main

import "golang.org/x/sys/windows/svc"

// mmaService implements the Windows Service interface.
// It acts only as a lifecycle bridge between SCM and appMain().
type mmaService struct{}

// Execute is called by the Windows Service Control Manager (SCM).
// This function MUST block until a stop or shutdown signal is received.
func (m *mmaService) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (bool, uint32) {
	// Notify SCM that service startup is in progress
	s <- svc.Status{State: svc.StartPending}

	// Start the normal application runtime (same as console/Docker)
	go appMain()

	// Notify SCM that service is now running
	s <- svc.Status{
		State:   svc.Running,
		Accepts: svc.AcceptStop | svc.AcceptShutdown,
	}

	// BLOCK here until SCM sends a stop or shutdown command.
	// This keeps the process alive exactly like NSSM would.
	for {
		c := <-r
		switch c.Cmd {
		case svc.Stop, svc.Shutdown:
			s <- svc.Status{State: svc.StopPending}
			return false, 0
		}
	}
}
