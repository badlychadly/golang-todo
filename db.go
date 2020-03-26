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
		// _, err = lb.CreateBucketIfNotExists([]byte("ITEM"))
		// if err != nil {
		// 	return err
		// }
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

	err = db.Lists.Update(func(txn *bolt.Tx) error {
		lsb := txn.Bucket([]byte("LISTS"))
		id, err := lsb.NextSequence()
		if err != nil {
			return err
		}
		if err = list.FmtId(id); err != nil {
			return err
		}
		lb, err := lsb.CreateBucketIfNotExists(itob(list.Id))
		if err != nil {
			return err
		}

		
		
		// list.Id = strconv.Itoa(int(id))
		listBytes, err := json.Marshal(list)
		if err != nil {
			return err
		}
		// err = lb.Put([]byte("Id"), itob(list.Id))
		// err = lb.Put([]byte("Name"), []byte(list.Name))
		err = lb.Put(itob(list.Id), listBytes)
		if err != nil {
			return err
		}
		return nil
	})
	fmt.Println("New list added")
	return
}

func (db *LDB) ViewLists() (listSlice []List) {
	db.Open()
	defer db.Close()
	// var btoi uint64
	// var listSlice []List

	err := db.Lists.View(func(tx *bolt.Tx) error {
		lsb := tx.Bucket([]byte("LISTS"))
		c := lsb.Cursor()
		var ok bool

		listSlice, ok = GetAll(c, lsb, listSlice).([]List) 
		if !ok {
			err := fmt.Errorf("Did not Work %v", listSlice)
			return err
		}
		// listSlice = GetAll(c)
		// if listSlice == nil {
		// 	err := fmt.Errorf("Empty set returned")
		// 	return err
		// }
		return nil
	})
	// fmt.Printf("the list: %v", listSlice)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// func GetAll(c *bolt.Cursor) (listSlice []List){
// 	for k, v := c.First(); k != nil; k, v = c.Next() {
// 		// btoi := binary.BigEndian.Uint16(k)
// 		list := List{Name: string(v)}
// 		if err := list.FmtId(k); err != nil {
// 			return listSlice
// 		}
// 		listSlice = append(listSlice, list)
// 		fmt.Printf("key=%T, value=%s\n", k, v)
// 	}
// 	return
// }

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
		lb := tx.Bucket([]byte("LIST"))
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




func itob(v uint16) []byte {
    b := make([]byte, 8)
	binary.BigEndian.PutUint16(b, uint16(v))
    return b
}