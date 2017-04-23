package libkademlia

import (
  "container/list"
)

// while node.next {
//   do sth
//
//   node = node.next
// }

type KBucket struct {
  contacts  list.List
}

// type RoutingTable []KBucket
// type RoutingTable struct {
//   buckets
// }

// type command struct {
//
// }


// func (k *Kademlia) updateBucket(contact Contact) {
//
// }

func (k *Kademlia) update(contact Contact) {
  bucket := k.getBucket(k.SelfContact.NodeID.Xor(contact.NodeID))

  for ele := bucket.contacts.Front(); ele != nil; ele = ele.Next() {
    if (ele.Value.(Contact).NodeID == contact.NodeID) {
      bucket.contacts.MoveToBack(ele)
      return
    }
  }

  if (bucket.contacts.Len() < kMax) {
    bucket.contacts.PushBack(contact)
  } else {
    _, err := k.DoPing(bucket.contacts.Front().Value.(Contact).Host, uint16(bucket.contacts.Front().Value.(Contact).Port))
    if err != nil {
      bucket.contacts.Remove(bucket.contacts.Front())
      bucket.contacts.PushBack(contact)
    }
  }
}

func (k *Kademlia) getBucket(dis ID) (ret *KBucket) {
  ret = &k.rt[dis.PrefixLen()]
  return
}
