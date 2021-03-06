package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"encoding/binary"
	// "bytes"

	"github.com/boltdb/bolt"
)

type LDB struct {
	Lists *bolt.DB
	filepath string
}

var DB *LDB

func (db *LDB) Initialize(filepath string)  {
	db.filepath = filepath
	db.Open()
	defer db.Close()

	err := db.Lists.Update(func(txn *bolt.Tx) error {
		_, err := txn.CreateBucketIfNotExists([]byte("LISTS"))
		if err != nil {
			return err
		}
		return nil
	})
	if err == nil {
		fmt.Println("DB Initialized")

	}
	return
}

func (db *LDB) Open() {
	if db.filepath == "" {
		log.Fatal("Filepath required for Database")
	}
	bdb, err := bolt.Open(db.filepath, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	db.Lists = bdb
}

func (db *LDB) Close() {
	db.Lists.Close()
}





func (db *LDB) CreateList(list *List) (err error) {
	db.Open()
	defer db.Close()

	errc := make(chan error, 1)
	err = db.Lists.Update(func(txn *bolt.Tx) error {
		go func() {
			defer close(errc)
			var bm Bm
			bm.Bucket = txn.Bucket([]byte("LISTS"))
			bm.NextSequence(errc)
	
			list.FmtId(bm.Id, errc)
				
			bm.CreateBucket(itob(list.Id), errc)
			bm.CreateChildBucket("ITEMS", errc)
			bm.Put(list.ToBytes(), errc)	
		}()
			err = ListenToChan(errc)

		return err
	})
	if err != nil {
		fmt.Printf("New Error %v\n", err)
		return err	
	}
	fmt.Println("New list added")
	return
}

func (db *LDB) ViewLists() (listSlice []List) {
	db.Open()
	defer db.Close()

	errc := make(chan error, 1)
	err := db.Lists.View(func(tx *bolt.Tx) error {
		go func(){
			defer close(errc)
			var bm Bm
			bm.Bucket = tx.Bucket([]byte("LISTS"))
			// lsb := tx.Bucket([]byte("LISTS"))
			c := bm.Bucket.Cursor()
			var ok bool
	
			listSlice, ok = GetAll(c, bm.Bucket, listSlice).([]List) 
			if !ok {
				// err := fmt.Errorf("Did not Work %v", listSlice)
				return
			}

		}()
		err := ListenToChan(errc)

		return err
	})
	// fmt.Printf("the list: %v", listSlice)
	if err != nil {
		log.Fatal(err)
	}
	return
}


func GetAll(c *bolt.Cursor, lsb *bolt.Bucket, objSlice interface{}) interface{} {
 
	switch mainSlice := objSlice.(type) {
		case []List:
			for k, v := c.First(); k != nil; k, v = c.Next() {
				lb := lsb.Bucket(k)
				ld := lb.Get(k)
				list := List{}
				err := json.Unmarshal(ld, &list)
				if err != nil {
					fmt.Errorf("error is: %v\n", err)
				}
				fmt.Printf("list data: %v, other value: %v\n", list, v)
				
				mainSlice = append(mainSlice, list)
				// fmt.Printf("key=%T, value=%s\n", k, list.Name)
			}
			return mainSlice
		case []Item:
			for k, v := c.First(); k != nil; k, v = c.Next() {
				item := Item{}
				item.Name = string(v)
				if err := item.FmtId(k); err != nil {
					return mainSlice
				}
				mainSlice = append(mainSlice, item)
				fmt.Printf("key=%T, value=%s\n", k, v)
			}
			return mainSlice
		default:
			return objSlice
	}
}


func (db *LDB) ViewList(id string) (list List, err error) {
	var items []Item
	db.Open()
	defer db.Close()
	err = db.Lists.View(func(tx *bolt.Tx) error {
		intId, err := strconv.Atoi(id)
		if err != nil {
			err = fmt.Errorf("Error converting to Integer %v", err)
			return err
		}
		lb := tx.Bucket([]byte("LISTS")).Bucket(itob(uint16(intId)))
		

		if lv := lb.Get(itob(uint16(intId))); lv == nil {
			err = fmt.Errorf("No List with id %v\n", intId)
			return err
		} else {
			if err = json.Unmarshal(lv, &list); err != nil {
				err = fmt.Errorf("error: %v\n", err)
				return err
			}
	
			if items, err = GetListItems(lb); err == nil {
				list.Items = items
			} 
		}

		return nil
	})

	return
}


func GetListItems(lb *bolt.Bucket) (itemSlice []Item, err error) {
	ib := lb.Bucket([]byte("ITEMS"))
	c := ib.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		item := Item{}
		err = json.Unmarshal(v, &item)
		if err != nil {
			err = fmt.Errorf("error is: %v\n", err)
			return
		}
		
		itemSlice = append(itemSlice, item)
	}
	return
}



func (db *LDB) DeleteList(id string) (err error) {
	db.Open()
	defer db.Close()
	err = db.Lists.Update(func(tx *bolt.Tx) error {
		lb := tx.Bucket([]byte("LISTS"))
		intId, _ := strconv.Atoi(id)
		err = lb.Delete(itob(uint16(intId)))
		if err != nil {
			return err
		}
		return nil
	})
	fmt.Printf("Successfully Deleted list with id %v\n", id)
	return
}

func (db *LDB) CreateItem(item *Item, listId string) (err error){
	lId, _ := strconv.Atoi(listId)
	db.Open()
	defer db.Close()
	err = db.Lists.Update(func(tx *bolt.Tx) error {
		lb := tx.Bucket([]byte("LISTS")).Bucket(itob(uint16(lId)))
		ib, err := lb.CreateBucketIfNotExists([]byte("ITEMS"))
		if err != nil {
			// err = fmt.Errorf("Bucket ITEMS could not be created %v\n", err)
			return err
		}
		itemId, err := ib.NextSequence()
		if err != nil {
			return err
		}

		item.SetListId(listId)
		item.FmtId(itemId)
		itemBytes, err := json.Marshal(item)
		if err != nil {
			return err
		}
		err = ib.Put(itob(item.Id), itemBytes)
		if err != nil {
			return err
		}
		return nil
	})
	fmt.Println("Item added db.go")
	return 
}




func itob(num interface{}) []byte {
	b := make([]byte, 8)
	switch v := num.(type) {
	case uint64:
		binary.BigEndian.PutUint16(b, uint16(v))
		return b
	case uint16:
		binary.BigEndian.PutUint16(b, uint16(v))
	}
    return b
}


func ListenToChan(ch chan error) (err error) {
	for err = range ch {
		fmt.Printf("Cl error: %v\n", err)
		if err != nil {
			return
		}
	}
	return 
}