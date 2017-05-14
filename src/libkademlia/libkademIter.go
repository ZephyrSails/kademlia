package libkademlia

// For project 2!
func (k *Kademlia) DoIterativeFindNode(id ID) (shortlist []Contact, err error) {
  findNodeCmd := findNodeCommand{ id, make(chan FindNodeResult), alpha }
	k.findNodeChan <- findNodeCmd

  findNodeRes := <- findNodeCmd.ResChan
  alphaNodes := findNodeRes.Nodes

  bufferContacts := make(Contact[])

  stop := false
  for (!stop) {

    shortlistChan

    for i := 0; i < len(alphaNodes); i++ {
      go ParallelDoFindNode(alphaNodes, shortlistChan)
    }

    for i := 0; i < len(alphaNodes); i++ {
      currDoFindNodeRes <- shortlistChan

      if !currDoFindNodeRes.err {
       bufferContacts.append(currDoFindNodeRes.ContactList)

       shortlist.append(alphaNodeContact)
      }
    }
    shortlist.sort()
    bufferContacts.sort()
    bufferContacts = bufferContacts[:kMax]

    if len(shortlist) >= kMax {
      stop = true

    } else if !shortlist.empty() && id.Xor(bufferContacts[0].Id).Less(id.Xor(shortlist[0].Id)) {
      alphaNodes = bufferContacts[:alpha]

    } else {
      shortlistChan
      for i := 0; i < kMax - len(shortlist); i++ {
        go ParallelDoPing(alphaNodes, shortlistChan)
      }
      for i := 0; i < len(alphaNodes); i++ {
        currDoPingRes <- ParallelDoPingChan

        if !currDoPingRes.err {
         shortlist.append(currDoPingRes.Id)
        }
      }
      stop = true
    }
  }

  shortlist = shortlist[:kMax]
  return

	return nil, &CommandFailed{"Not implemented"}
}

func (k *Kademlia) DoIterativeStore(key ID, value []byte) ([]Contact, error) {
	return nil, &CommandFailed{"Not implemented"}
}

func (k *Kademlia) DoIterativeFindValue(key ID) (value []byte, err error) {
	return nil, &CommandFailed{"Not implemented"}
}
