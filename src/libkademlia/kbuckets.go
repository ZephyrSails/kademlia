type KBucket struct {
  contacts  []Contact
}


type RoutingTable struct {
  buckets   []KBucket
}



func distance(id1 ID, id2 ID) (ret ID) {
  // 20 * 8
}

func (kb *KBucket) update(nodeId ID) {

}

func (rt *RoutingTable) getBucket(dis ID) (ret *KBucket) {

}

func (rt *RoutingTable) cloestKContacts(key ID, k int) (ret []Contact) {

}

// rt.getBucket(distance()).update()
