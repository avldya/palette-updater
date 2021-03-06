package service_control

import (
	"golang.org/x/sys/windows/svc"
)

type ServiceControl struct {
}

func (sc *ServiceControl) Install(svcName, svcDisplayName, svcDescription string) error {
	return installService(svcName, svcDisplayName, svcDescription)
}

func (sc *ServiceControl) Remove(svcName string) error {
	return removeService(svcName)
}

func (sc *ServiceControl) Start(svcName string) error {
	return startService(svcName)
}

func (sc *ServiceControl) Stop(svcName string) error {
	return controlService(svcName, svc.Stop, svc.Stopped)
}

func (sc *ServiceControl) Query(svcName string) (svc.Status, error) {
	return queryService(svcName)
}

// The following functions are not necessary, but they were
// already implemented on Windows.
func (sc *ServiceControl) Pause(svcName string) error {
	return controlService(svcName, svc.Pause, svc.Paused)
}

func (sc *ServiceControl) Continue(svcName string) error {
	return controlService(svcName, svc.Continue, svc.Running)
}
