package main

import (
	// "encoding/json"
	"encoding/binary"
	"fmt"
	"strconv"
)


type Item struct {
	Id uint16
	Name string `json:"name"`
	ListId uint16
}


func (item *Item) FmtId(num interface{}) (err error) {
	switch v := num.(type) {
		case []uint8:
			item.Id = binary.BigEndian.Uint16(v)
			return
		case uint16: 
			item.Id = v
			return
		case uint64:
			id := uint16(v)
			item.Id = id
			return
		case int:
			item.Id = uint16(v)
		default:
			err = fmt.Errorf("Unaccepted type %T", num)
			return
	}
	return
}


func (item *Item) SetListId(num interface{}) (err error) {
	switch v := num.(type) {
		case string:
			intId, err := strconv.Atoi(v)
			// if err != nil {
			// 	return err
			// }
			item.ListId = uint16(intId) 
			return err
		default:
			err = fmt.Errorf("Unaccepted type")
			return err
	}
	return
}