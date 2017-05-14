package libkademlia

import (
	//"net"
	"strconv"
	"testing"
	//"time"
	"log"
	//"fmt"
)

func TestIterativeFindNode(t *testing.T) {
	instance1 := NewKademlia("localhost:9104")
	instance2 := NewKademlia("localhost:9105")
	host2, port2, _ := StringToIpPort("localhost:9105")
	instance1.DoPing(host2, port2)


	tree_node := make([]*Kademlia, 20)
	for i := 0; i < 10; i++ {
		address := "localhost:" + strconv.Itoa(9106+i)
		tree_node[i] = NewKademlia(address)
		host_number, port_number, _ := StringToIpPort(address)
		instance2.DoPing(host_number, port_number)
	}


	contacts, err := tree_node[0].DoIterativeFindNode(tree_node[3].NodeID)

	if err == nil {
		for i := 0; i < len(contacts); i++ {
			log.Println(contacts[i].NodeID.AsString())
		}
	}
	log.Println(tree_node[3].NodeID.AsString())

	t.Error("test error")

}
