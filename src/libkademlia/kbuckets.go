package libkademlia

import (
  "container/list"
  "log"
  "errors"
)

// while node.next {
//   do sth
//
//   node = node.next
// }

type KBucket struct {
  contacts  list.List
}

func (k *Kademlia) update(contact Contact) {
  //log.Println("update contact called: ", contact)
  bucket := k.getBucket(k.SelfContact.NodeID.Xor(contact.NodeID))

  for ele := bucket.contacts.Front(); ele != nil; ele = ele.Next() {
    if (ele.Value.(Contact).NodeID == contact.NodeID) {
      bucket.contacts.MoveToBack(ele)
      return
    }
  }

  if (bucket.contacts.Len() < kMax) {
    //log.Println("PushBack new contact.")
    bucket.contacts.PushBack(contact)
    //k.printBucket(bucket)
  } else {
    _, err := k.DoPing(bucket.contacts.Front().Value.(Contact).Host, uint16(bucket.contacts.Front().Value.(Contact).Port))
    if err != nil {
      bucket.contacts.Remove(bucket.contacts.Front())
      bucket.contacts.PushBack(contact)
    }
  }
}

func (k *Kademlia) getContact(NodeID ID) (ret findContactResponse) {
  //log.Println("getContact called with NodeID: ", NodeID)

  bucket := k.getBucket(k.SelfContact.NodeID.Xor(NodeID))
  //k.printBucket(bucket)

  for ele := bucket.contacts.Front(); ele != nil; ele = ele.Next() {
    if (ele.Value.(Contact).NodeID == NodeID) {
      //log.Println("NodeID: ", NodeID, " found.")
      ret =  findContactResponse{ele.Value.(Contact), nil}
      return
    }
  }
  //log.Println("NodeID not found.")
  //TODO: Find the better implementation
  con := Contact{}
  ret = findContactResponse{con, errors.New("Contact Not found")}
  return
}

func (k *Kademlia) getLocalValue(searchKey ID) (ret findLocalValueResponse) {
  //log.Println("getContact called with NodeID: ", NodeID)
  if _, ok := k.hash[searchKey]; ok {
      
      ret = findLocalValueResponse{k.hash[searchKey], nil}
     } else {
       ret = findLocalValueResponse{make([]byte, 0), errors.New("Local Value Not found")}
     }
  return
}
 

func (k *Kademlia) getBucket(dis ID) (ret *KBucket) {
  ret = &k.rt[dis.PrefixLen()]
  return
}

func (k *Kademlia) printBucket(bucket *KBucket) {
  log.Println("Printing bucket.")

  for ele := bucket.contacts.Front(); ele != nil; ele = ele.Next() {
    log.Println("NodeID: ", ele.Value.(Contact).NodeID)
  }

  log.Println("Finish printing.")
  return
}
