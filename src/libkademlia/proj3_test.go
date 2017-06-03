package libkademlia

import (
	"log"
	"testing"
)

func TestVanish(t *testing.T) {
	data := []byte("Hello World")

	instance1 := NewKademlia("localhost:9404")
	instance2 := NewKademlia("localhost:9405")

	host2, port2, _ := StringToIpPort("localhost:9305")
	instance1.DoPing(host2, port2)

	_, vdoid := instance1.Vanish(data, byte(20), byte(10), 0)

	vdo_data := instance2.Unvanish(instance1.NodeID, vdoid)

	log.Println(vdo_data)

	t.Error("test error")
}