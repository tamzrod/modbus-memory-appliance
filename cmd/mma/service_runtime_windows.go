//go:build windows

package main

import "golang.org/x/sys/windows/svc"

type mmaService struct{}

func (m *mmaService) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (bool, uint32) {
	s <- svc.Status{State: svc.StartPending}

	go appMain() // 🔑 start existing runtime (unchanged)

	s <- svc.Status{
		State:   svc.Running,
		Accepts: svc.AcceptStop | svc.AcceptShutdown,
	}

	for c := range r {
		switch c.Cmd {
		case svc.Stop, svc.Shutdown:
			// hard exit is acceptable; no runtime refactor
			s <- svc.Status{State: svc.StopPending}
			return false, 0
		}
	}

	return false, 0
}
