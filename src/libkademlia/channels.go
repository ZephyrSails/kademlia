package libkademlia

import (
  "errors"
)

type findContactResponse struct {
	Result   Contact
	Err      error
}

type findLocalValueResponse struct {
	Result   []byte
	Err  	 error
}

type findLocalValueCommand struct {
	SearchKey          ID
	LocalValueChan     chan findLocalValueResponse
}

type findContactCommand struct {
	NodeID       ID
	ContactChan  chan findContactResponse
}

type updateCommand struct {
	Sender     Contact
}

type storeCommand struct {
	Key        ID
	Value      []byte
	ErrChan    chan error
}

type findNodeCommand struct {
	NodeID 	   ID
	ResChan    chan FindNodeResult
  N          int
}

// type findValueCommand struct {
// 	Key    	   ID
//   ResChan    chan FindValueResult
// }

func (k *Kademlia) routingHandler() {
	//log.Println("Handler online")
	for {
		select {
    case findContactCmd := <- k.findContactChan:
			//log.Println("findContactCmd received")
			findContactCmd.ContactChan <- k.getContact(findContactCmd.NodeID)

		case updateCmd := <- k.updateChan:
			//log.Println("updateCmd received: ", updateCmd.Sender)
			k.update(updateCmd.Sender)

    case findNodeCmd := <- k.findNodeChan:
      findNodeCmd.ResChan <- k.getNContacts(findNodeCmd.NodeID, findNodeCmd.N)
		}
	}
}

func (k *Kademlia) storageHandler() {
	for {
		select {
    case storeCmd := <- k.storeChan:
      if _, ok := k.hash[storeCmd.Key]; ok {
        // storeCmd.ErrChan <- errors.New("Hash conflict")
        storeCmd.ErrChan <- errors.New("")
      } else {
        k.hash[storeCmd.Key] = storeCmd.Value
        storeCmd.ErrChan <- nil
      }

    case findLocalValueCmd := <- k.findLocalValueChan:
      findLocalValueCmd.LocalValueChan <- k.getLocalValue(findLocalValueCmd.SearchKey)
		}
	}
}

func (k *Kademlia) initChans() {
  k.findContactChan = make(chan findContactCommand)
  k.updateChan = make(chan updateCommand)
  k.storeChan = make(chan storeCommand)
  k.findNodeChan = make(chan findNodeCommand)
  k.findLocalValueChan = make(chan findLocalValueCommand)

  go k.routingHandler()
  go k.storageHandler()
}
