package utils

import (
	"errors"
	"time"

	"github.com/sony/sonyflake"
)

// cluster node
var flake *sonyflake.Sonyflake

func GenID() (id uint64, err error) {
	if flake == nil {
		return id, errors.New("sony flake not inited")
	}
	id, err = flake.NextID() //该节点进入集群等待队列中等待生成ID
	return
}

// below codes are the initilization of flake node
func getMachineID() (uint16, error) {
	return 10000, nil
}

func Init(machineID uint16) (err error) {
	t, _ := time.Parse("2006-01-02", "2024-05-01")
	settings := sonyflake.Settings{
		StartTime: t,
		MachineID: getMachineID,
	}
	flake = sonyflake.NewSonyflake(settings)
	return
}
