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
	NodeID      		ID
	SelfContact 		Contact
	hash 						map[ID][]byte
	rt							[]KBucket
	findContactChan		chan findContactCommand
	findLocalValueChan	chan findLocalValueCommand
	updateChan				chan updateCommand
	storeChan					chan storeCommand
	findNodeChan			chan findNodeCommand
	findValueChan			chan findValueCommand
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
	// hostname, port, err := net.SplitHostPort(laddr)
	// if err != nil {
	// 	return nil
	// }
	// s.HandleHTTP(rpc.DefaultRPCPath+hostname+port,
	// 	rpc.DefaultDebugPath+hostname+port)
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
// func (k *Kademlia) DoPing(host string, port string) (*Contact, error) {
	// TODO: Implement
// <<<<<<< HEAD
// 	//log.Println("DoPing Called.")
// 	hostStr := host.String()
// 	portStr := strconv.Itoa(int(port))

// 	// hostStr = "localhost"
// 	//log.Printf("DoPing: rpc.DefaultRPCPath+hostStr+portStr:", rpc.DefaultRPCPath+hostStr+portStr)
// =======

	hostStr, portStr := IpPortToString(host, port)

	// hostStr = "localhost"
	log.Printf(net.JoinHostPort(hostStr, portStr))
	log.Printf(rpc.DefaultRPCPath+hostStr+portStr)
// >>>>>>> 9f3bdb3f413107e183dde86aef1ff6d29697c791

	client, err := rpc.DialHTTPPath("tcp", net.JoinHostPort(hostStr, portStr),
		rpc.DefaultRPCPath+hostStr+portStr)
	if err != nil {
		// log.Printf(rpc.DefaultRPCPath+hostStr+portStr)
		log.Fatal("DialHTTP: ", err)
	}

	//log.Printf("Pinging initial peer\n")

	// This is a sample of what an RPC looks like
	// TODO: Replace this with a call to your completed DoPing!
	ping := new(PingMessage)
	ping.MsgID = NewRandomID()
	ping.Sender = k.SelfContact
	var pong PongMessage

	//log.Println("ping.Sender in DoPing:", ping.Sender)
	err = client.Call("KademliaRPC.Ping", ping, &pong)

	if err != nil {
		log.Fatal("Call: ", err)
		return nil, &CommandFailed {
			"Unable to ping " + fmt.Sprintf("%s:%v", host.String(), port) }
	}
	//log.Printf("ping msgID: %s\n", ping.MsgID.AsString())
	//log.Printf("pong msgID: %s\n\n", pong.MsgID.AsString())
	//log.Println("Pong rcved. Update.")
	pongCmd := updateCommand{ pong.Sender }
	k.updateChan <- pongCmd
	// k.update(pong.Sender)
	return &pong.Sender, nil
}

func (k *Kademlia) DoStore(contact *Contact, key ID, value []byte) error {
	// TODO: Implement
	hostStr, portStr := IpPortToString(contact.Host, contact.Port)

	// hostStr = "localhost"
	// log.Printf(net.JoinHostPort(hostStr, portStr))
	// log.Printf(rpc.DefaultRPCPath+hostStr+portStr)
// >>>>>>> 9f3bdb3f413107e183dde86aef1ff6d29697c791

	client, err := rpc.DialHTTPPath("tcp", net.JoinHostPort(hostStr, portStr),
		rpc.DefaultRPCPath+hostStr+portStr)
	if err != nil {
		// log.Printf(rpc.DefaultRPCPath+hostStr+portStr)
		log.Fatal("DialHTTP: ", err)
	}
	req := new(StoreRequest)
	req.Sender = k.SelfContact
	req.MsgID = NewRandomID()
	req.Key = key
	req.Value = value
	var res StoreResult
	err = client.Call("KademliaRPC.Store", req, &res)
	if err != nil {
		log.Fatal("Call: ", err)
		return &CommandFailed {
			"Unable to store " + fmt.Sprintf("%s:%v", contact.Host.String(), contact.Port) }
	}
	pongCmd := updateCommand{ *contact }
	k.updateChan <- pongCmd
	return nil
	// return &CommandFailed{"Not implemented"}
}

func (k *Kademlia) DoFindNode(contact *Contact, searchKey ID) ([]Contact, error) {
	// TODO: Implement
	return nil, &CommandFailed{"Not implemented"}
}

func (k *Kademlia) DoFindValue(contact *Contact,
	searchKey ID) (value []byte, contacts []Contact, err error) {
	// TODO: Implement
	return nil, nil, &CommandFailed{"Not implemented"}
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
