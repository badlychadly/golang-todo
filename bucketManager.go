package main 


import (
	"github.com/boltdb/bolt"
	"fmt"
)


type Bm struct {
	Bucket *bolt.Bucket
	Id interface{}
	err error
}

func (bm *Bm) NextSequence(errc chan error) {
	if bm.err != nil {return}
	if id, err := bm.Bucket.NextSequence(); err != nil {
		bm.err = fmt.Errorf("Problem getting next sequence in Bucket %v\n", err)
		errc <- bm.err
		return 
	} else {
		bm.Id = uint16(id)
		// errc <- nil
	}
	return
}


func (bm *Bm) CreateBucket(indexKey []byte, errc chan error ) {
	if bm.err != nil {return}
	// fmt.Println("In CreateBucket")
	lb, err := bm.Bucket.CreateBucketIfNotExists(indexKey)
	if err != nil {
		bm.err = fmt.Errorf("Could not create bucket: %v\n", err)
		fmt.Printf("CB err %v\n", bm.err)
		errc <- bm.err
		return 
	}
	bm.Bucket = lb

	return

}

func (bm *Bm) Put(obj []byte, errc chan error) {
	if bm.err != nil {return}
	fmt.Printf("inside Put\n")
	err := bm.Bucket.Put(itob(bm.Id), obj)
	if err != nil {
		errc <- fmt.Errorf("failed to add Data %v\n", err)
		return 
	}
	// fmt.Printf("type bm.Id: %T\n", bm.Id)
	return
}