package libkademlia

import (
  "sort"
)

type ParallelDoFindNodeRes struct {
  Target    *Contact
  Nodes     []Contact
  Err       bool
}

func (k *Kademlia) ParallelDoFindNode (target *Contact, key ID, resChan chan ParallelDoFindNodeRes) {
  contacts, err := k.DoFindNode(target, key)
  if err != nil {
    resChan <- ParallelDoFindNodeRes{ target, contacts, true }
  } else {
    resChan <- ParallelDoFindNodeRes{ target, contacts, false }
  }
}

type ParallelDoPingRes struct {
  Target    *Contact
  Err       bool
}

func (k *Kademlia) ParallelDoPing (target *Contact, resChan chan ParallelDoPingRes) {
  _, err := k.DoPing(target.Host, target.Port)
  if err != nil {
    resChan <- ParallelDoPingRes{ target, true }
  } else {
    resChan <- ParallelDoPingRes{ target, false }
  }
}

type ContactTargetSlice struct {
  Contacts  []Contact
  Target    ID
}

func (self ContactTargetSlice) Less(i, j int) bool {
  return self.Contacts[i].NodeID.Xor(self.Target).Less(self.Contacts[j].NodeID.Xor(self.Target))
}

func (self ContactTargetSlice) Len() int {
  return len(self.Contacts)
}

func (self ContactTargetSlice) Swap(i, j int) {
  self.Contacts[i], self.Contacts[j] = self.Contacts[j], self.Contacts[i]
}

func toContactTargetSlice(l []Contact, target ID) (contacts ContactTargetSlice) {
  contacts.Contacts = l
  contacts.Target = target
  return
}



// For project 2!
func (k *Kademlia) DoIterativeFindNode(id ID) (shortlist []Contact, err error) {
  findNodeCmd := findNodeCommand{ id, make(chan FindNodeResult), alpha }
	k.findNodeChan <- findNodeCmd

  findNodeRes := <- findNodeCmd.ResChan
  alphaNodes := findNodeRes.Nodes

  var bufferContacts []Contact

  stop := false
  for (!stop) {

    ParallelDoFindNodeResChan := make(chan ParallelDoFindNodeRes)

    for i := 0; i < len(alphaNodes); i++ {
      go k.ParallelDoFindNode(&alphaNodes[i], id, ParallelDoFindNodeResChan)
    }

    for i := 0; i < len(alphaNodes); i++ {
      currDoFindNodeRes := <- ParallelDoFindNodeResChan

      if !currDoFindNodeRes.Err {
        bufferContacts = append(bufferContacts, currDoFindNodeRes.Nodes...)

        shortlist = append(shortlist, *currDoFindNodeRes.Target)
      }
    }
    temp := toContactTargetSlice(shortlist, id)
    sort.Sort(temp)
    shortlist = temp.Contacts

    temp = toContactTargetSlice(bufferContacts, id)
    sort.Sort(temp)
    bufferContacts = temp.Contacts

    shortlist = removeDup(shortlist)
    bufferContacts = removeDup(bufferContacts)
    bufferContacts = minus(bufferContacts, shortlist)

    bufferContacts = firstKEle(bufferContacts, kMax - len(shortlist))

    if len(shortlist) >= kMax {
      stop = true

    } else if len(shortlist) > 0 && id.Xor(bufferContacts[0].NodeID).Less(id.Xor(shortlist[0].NodeID)) {

      alphaNodes = firstKEle(bufferContacts, alpha)

    } else {
      ParallelDoPingResChan := make(chan ParallelDoPingRes)

      for i := 0; i < len(bufferContacts); i++ {
        go k.ParallelDoPing(&bufferContacts[i], ParallelDoPingResChan)
      }

      for i := 0; i < len(bufferContacts); i++ {
        currDoPingRes := <- ParallelDoPingResChan

        if !currDoPingRes.Err {
          shortlist = append(shortlist, *currDoPingRes.Target)
        }
      }
      stop = true
    }
  }

  shortlist = firstKEle(shortlist, kMax)
  return
}

func (k *Kademlia) DoIterativeStore(key ID, value []byte) (res []Contact, err error) {
  contacts, err := k.DoIterativeFindNode(key)

  if err == nil {
    for i := 0; i < len(contacts); i++ {
      go func () {
        err = k.DoStore(&contacts[i], key, value)
        if err == nil {
          res = append(res, contacts[i])
        }
      } ()
    }
    return res, nil
  }
  return nil, &CommandFailed{ "DoIterativeStore Failed" }
}



func (k *Kademlia) DoIterativeFindValue(key ID) (value []byte, err error) {
  findNodeCmd := findNodeCommand{ id, make(chan FindNodeResult), alpha }
  k.findNodeChan <- findNodeCmd

  findNodeRes := <- findNodeCmd.ResChan
  alphaNodes := findNodeRes.Nodes

  var bufferContacts []Contact

  stop := false
  for (!stop) {

    ParallelDoFindNodeResChan := make(chan ParallelDoFindNodeRes)

    for i := 0; i < len(alphaNodes); i++ {
      go k.ParallelDoFindNode(&alphaNodes[i], id, ParallelDoFindNodeResChan)
    }

    for i := 0; i < len(alphaNodes); i++ {
      currDoFindNodeRes := <- ParallelDoFindNodeResChan

      if !currDoFindNodeRes.Err {
        bufferContacts = append(bufferContacts, currDoFindNodeRes.Nodes...)

        shortlist = append(shortlist, *currDoFindNodeRes.Target)
      }
    }
    temp := toContactTargetSlice(shortlist, id)
    sort.Sort(temp)
    shortlist = temp.Contacts

    temp = toContactTargetSlice(bufferContacts, id)
    sort.Sort(temp)
    bufferContacts = temp.Contacts

    shortlist = removeDup(shortlist)
    bufferContacts = removeDup(bufferContacts)
    bufferContacts = minus(bufferContacts, shortlist)

    bufferContacts = firstKEle(bufferContacts, kMax - len(shortlist))

    if len(shortlist) >= kMax {
      stop = true

    } else if len(shortlist) > 0 && id.Xor(bufferContacts[0].NodeID).Less(id.Xor(shortlist[0].NodeID)) {

      alphaNodes = firstKEle(bufferContacts, alpha)

    } else {
      ParallelDoPingResChan := make(chan ParallelDoPingRes)

      for i := 0; i < len(bufferContacts); i++ {
        go k.ParallelDoPing(&bufferContacts[i], ParallelDoPingResChan)
      }

      for i := 0; i < len(bufferContacts); i++ {
        currDoPingRes := <- ParallelDoPingResChan

        if !currDoPingRes.Err {
          shortlist = append(shortlist, *currDoPingRes.Target)
        }
      }
      stop = true
    }
  }

  shortlist = firstKEle(shortlist, kMax)
  return
}
