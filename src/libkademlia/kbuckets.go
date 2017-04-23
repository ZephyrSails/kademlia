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

type command struct {

}

func (id1 *ID) dist(id2 ID) (ret ID) {
  // 20 * 8

}

func (bucket *KBucket) updateBucket(contact Contact) {
  for ele := bucket.contacts.Front(); ele != nil; ele = ele.Next() {
    if (ele.NodeID == contact.NodeID) {
      KBucket.contacts.MoveToBack(ele)
      return
    }
  }

  if (KBucket.contacts.Len() < kMax) {
    KBucket.contacts.PushBack(contact)
  } else {
    _, err = k.DoPing(bucket.contacts.Front().Host, uint16(bucket.contacts.Front().Port))
    if err != nil {
      bucket.contacts.Front() = contact
      KBucket.contacts.MoveToBack(bucket.contacts.Front())
    }
  }
  return
}

func (k *Kademlia) update(contact Contact) {
  bucket := k.getBucket(k.SelfContact.NodeID.dist(contact.NodeID))

  bucket.updateBucket(contact)
}

func (k *Kademlia) getBucket(dis ID) (ret *KBucket) {

}

// func (rt *RoutingTable) cloestKContacts(key ID, k int) (ret []Contact) {
//
// }

// rt.getBucket(distance()).update()
