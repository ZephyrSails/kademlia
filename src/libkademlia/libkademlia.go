package libkademlia

// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
)

const (
	alpha = 3
	b     = 8 * IDBytes
	kMax  = 20
)

// Kademlia type. You can put whatever state you need in this.
type Kademlia struct {
	NodeID      				ID
	SelfContact 				Contact
	hash 								map[ID][]byte
	rt									[]KBucket
	findContactChan			chan findContactCommand
	findLocalValueChan	chan findLocalValueCommand
	updateChan					chan updateCommand
	storeChan						chan storeCommand
	findNodeChan				chan findNodeCommand
	findValueChan				chan findValueCommand
}

func NewKademliaWithId(laddr string, nodeID ID) *Kademlia {
	k := new(Kademlia)
	k.initChans()
	k.NodeID = nodeID
	k.rt = make([]KBucket, IDBits)
	k.hash = make(map[ID][]byte)

	// TODO: Initialize other state here as you add functionality.


	// Set up RPC server
	// NOTE: KademliaRPC is just a wrapper around Kademlia. This type includes
	// the RPC functions.

	s := rpc.NewServer()

	kRPC := KademliaRPC{k}

	// s.Register(&KademliaRPC{k})
	s.Register(&kRPC)

	h, p, _ := StringToIpPort(laddr)
	hostname, port := IpPortToString(h, p)

	s.HandleHTTP(rpc.DefaultRPCPath+hostname+port,
		rpc.DefaultDebugPath+hostname+port)

	l, err := net.Listen("tcp", laddr)
	if err != nil {
		log.Fatal("Listen: ", err)
	}

	// Run RPC server forever.
	go http.Serve(l, nil)

	// Add self contact
	hostname, port, _ = net.SplitHostPort(l.Addr().String())
	port_int, _ := strconv.Atoi(port)

	ipAddrStrings, err := net.LookupHost(hostname)
	var host net.IP
	for i := 0; i < len(ipAddrStrings); i++ {
		host = net.ParseIP(ipAddrStrings[i])
		if host.To4() != nil {
			break
		}
	}
	k.SelfContact = Contact{k.NodeID, host, uint16(port_int)}

	return k
}

func NewKademlia(laddr string) *Kademlia {
	return NewKademliaWithId(laddr, NewRandomID())
}

type ContactNotFoundError struct {
	id  ID
	msg string
}

func (e *ContactNotFoundError) Error() string {
	return fmt.Sprintf("%x %s", e.id, e.msg)
}

func (k *Kademlia) FindContact(nodeId ID) (*Contact, error) {
	// TODO: Search through contacts, find specified ID
	// Find contact with provided ID
	if nodeId == k.SelfContact.NodeID {
		//log.Println("FindContact find itself.")
		return &k.SelfContact, nil
	}

	//TODO: Give this variable a better name
	//Note: use new will cause problem because it generate a pointer
	//cmd := new(findContactCommand)
	//cmd.NodeID = nodeId
	//cmd.ContactChan = make(chan Contact)

	cmd := findContactCommand{nodeId, make(chan findContactResponse)}
	k.findContactChan <- cmd

	//TODO: Give this variable a better name
	result := <- cmd.ContactChan
	//log.Println("result: ", result.result, "err: ", result.err)
	if result.Err == nil {
		//log.Println("ID found.")
		return &result.Result, nil
	} else {
		//log.Println("Not found.")
		return nil, &ContactNotFoundError{nodeId, "Not found"}
	}
}

type CommandFailed struct {
	msg string
}

func (e *CommandFailed) Error() string {
	return fmt.Sprintf("%s", e.msg)
}

func (k *Kademlia) DoPing(host net.IP, port uint16) (*Contact, error) {

	req := PingMessage{ k.SelfContact, NewRandomID()}
	var res PongMessage

	client := getClient(host, port)
	err := client.Call("KademliaRPC.Ping", req, &res)
	if err != nil {
		log.Fatal("Call: ", err)
		return nil, &CommandFailed {
			"Unable to ping " + fmt.Sprintf("%s:%v", host.String(), port) }
	}

	updateCmd := updateCommand{ res.Sender }
	k.updateChan <- updateCmd
	return &res.Sender, err
}

func (k *Kademlia) DoStore(contact *Contact, key ID, value []byte) (error) {

	req := StoreRequest{ k.SelfContact, NewRandomID(), key, value }
	var res StoreResult

	client := getClient(contact.Host, contact.Port)
	err := client.Call("KademliaRPC.Store", req, &res)
	if err != nil {
		log.Fatal("Call: ", err)
		return &CommandFailed {
			"Unable to store " + fmt.Sprintf("%s:%v", contact.Host.String(), contact.Port) }
	}

	updateCmd := updateCommand{ *contact }
	k.updateChan <- updateCmd
	return err
}

func (k *Kademlia) DoFindNode(contact *Contact, searchKey ID) ([]Contact, error) {

	req := FindNodeRequest{ k.SelfContact, NewRandomID(), searchKey }
	var res FindNodeResult

	client := getClient(contact.Host, contact.Port)
	err := client.Call("KademliaRPC.FindNode", req, &res)
	if err != nil {
		log.Fatal("Call: ", err)
		return nil, &CommandFailed {
			"Unable to find node " + fmt.Sprintf("%s:%v", contact.Host.String(), contact.Port) }
	}
	updateCmd := updateCommand{ *contact }
	k.updateChan <- updateCmd

	return res.Nodes, &CommandFailed{"Not sure which error yet"}
}

func (k *Kademlia) DoFindValue(contact *Contact,
	searchKey ID) ([]byte, []Contact, error) {

	req := FindValueRequest{ k.SelfContact, NewRandomID(), searchKey }
	var res FindValueResult

	client := getClient(contact.Host, contact.Port)
	err := client.Call("KademliaRPC.FindValue", req, &res)
	if err != nil {
		log.Fatal("Call: ", err)
		return nil, nil, &CommandFailed {
			"Unable to find value " + fmt.Sprintf("%s:%v", contact.Host.String(), contact.Port) }
	}
	updateCmd := updateCommand{ *contact }
	k.updateChan <- updateCmd

	return res.Value, res.Nodes, &CommandFailed{"Not sure which error yet"}
	// return nil, nil, &CommandFailed{"Not implemented"}
}

// func (k *Kademlia) FindContact(nodeId ID) (*Contact, error) {
// 	if nodeId == k.SelfContact.NodeID {
// 		return &k.SelfContact, nil
// 	}
// 	cmd := findContactCommand{nodeId, make(chan findContactResponse)}
// 	k.findContactChan <- cmd
// 	//TODO: Give this variable a better name
// 	result := <- cmd.ContactChan
// 	if result.Err == nil {
// 		return &result.Result, nil
// 	} else {
// 		return nil, &ContactNotFoundError{nodeId, "Not found"}
// 	}
// }

type LocalValueNotFoundError struct {
	searchKey  ID
	msg string
}

func (e *LocalValueNotFoundError) Error() string {
	return fmt.Sprintf("%x %s", e.searchKey, e.msg)
}

func (k *Kademlia) LocalFindValue(searchKey ID) ([]byte, error) {
	// TODO: Implement
	cmd := findLocalValueCommand{searchKey, make(chan findLocalValueResponse)}
	k.findLocalValueChan <- cmd
	result := <- cmd.LocalValueChan
	if result.Err == nil {
		return result.Result, nil
	} else {
		return nil, &LocalValueNotFoundError{searchKey, "Not found"}
	}
	// return []byte(""), &CommandFailed{"Not implemented"}
}

// For project 2!
func (k *Kademlia) DoIterativeFindNode(id ID) ([]Contact, error) {
	return nil, &CommandFailed{"Not implemented"}
}
func (k *Kademlia) DoIterativeStore(key ID, value []byte) ([]Contact, error) {
	return nil, &CommandFailed{"Not implemented"}
}
func (k *Kademlia) DoIterativeFindValue(key ID) (value []byte, err error) {
	return nil, &CommandFailed{"Not implemented"}
}

// For project 3!
func (k *Kademlia) Vanish(data []byte, numberKeys byte,
	threshold byte, timeoutSeconds int) (vdo VanashingDataObject) {
	return
}

func (k *Kademlia) Unvanish(searchKey ID) (data []byte) {
	return nil
}
