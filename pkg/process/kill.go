package process

import (
	"errors"
	"syscall"
)

const maxProcsToScan = 2048

func KillAll(name string, signal syscall.Signal) error {

	filters := FiltersInit("", "")
	infoMap := make(ProcInfoMap)
	pids := make(Pidlist, 0, maxProcsToScan)
	procPrev := NewProcSampleList(maxProcsToScan)
	GetPidList(&pids, maxProcsToScan)
	ProcStatsReader(pids, filters, &procPrev, infoMap)

	count := 0
	for _, proc := range infoMap {
		if proc.Friendly == name {
			count++
			err := syscall.Kill(int(proc.Pid), signal)
			if err != nil {
				return err
			}
		}
	}
	if count == 0 {
		return errors.New("no such process")
	}
	return nil

}
