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
	Sender 	   Contact
	NodeID 	   ID
	ResChan    chan []Contact
}

type findValueCommand struct {
	Sender 	   Contact
	Key    	   ID
	ValChan	   chan []byte
	NodeChan 	 chan []Contact
}

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
		}
	}
}

func (k *Kademlia) storageHandler() {
	for {
		select {
    case storeCmd := <- k.storeChan:
      if _, ok := k.hash[storeCmd.Key]; ok {
        storeCmd.ErrChan <- errors.New("Contact Not found")
      } else {
        k.hash[storeCmd.Key] = storeCmd.Value
        storeCmd.ErrChan <- nil
        

      }

      case findLocalValuecmd := <- k.findLocalValueChan:
      		findLocalValuecmd.LocalValueChan <- k.getLocalValue(findLocalValuecmd.SearchKey)

		}
	}
}

func (k *Kademlia) initChans() {
  k.findContactChan = make(chan findContactCommand)
  k.updateChan = make(chan updateCommand)
  k.storeChan = make(chan storeCommand)
  k.findNodeChan = make(chan findNodeCommand)
  k.findValueChan = make(chan findValueCommand)
  k.findLocalValueChan = make(chan findLocalValueCommand)

  go k.routingHandler()
  go k.storageHandler()
}
