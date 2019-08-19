package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/sony/sonyflake"
)

func readMachineIDFromLocalFile() uint16 {
	// read from file, assume value = 1
	return 1
}

func generateMachineID() (uint16, error) {
	// use random utils to generate machine ID
	var data uint16
	data = uint16(rand.Intn(10))
	return data, nil
}

func saddMachineIDToRedisSet() (uint16, error) {
	// mock func: add machine ID to redis
	return 0, nil
}

func saveMachineIDToLocalFile(machineID uint16) error {
	// mock func: save machine ID to redis
	return nil
}

func getMachineID() (uint16, error) {
	var machineID uint16
	var err error
	machineID = readMachineIDFromLocalFile()
	if machineID == 0 {
		machineID, err = generateMachineID()
		if err != nil {
			return 0, err
		}
	}

	return machineID, nil
}

func checkMachineID(machineID uint16) bool {
	saddResult, err := saddMachineIDToRedisSet()
	if err != nil || saddResult == 0 {
		return true
	}

	err = saveMachineIDToLocalFile(machineID)
	if err != nil {
		return true
	}

	return false
}

func main() {
	t, _ := time.Parse("2006-01-02", "2018-01-01")
	settings := sonyflake.Settings{
		StartTime:      t,
		MachineID:      getMachineID,
		CheckMachineID: checkMachineID,
	}

	sf := sonyflake.NewSonyflake(settings)
	id, err := sf.NextID()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(id)
}
