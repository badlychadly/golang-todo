package main

import (
	// "encoding/json"
	"encoding/binary"
	"fmt"
	// "strconv"
)

type List struct {
	Id uint16
	Name string `json:"name"`
	Items []Item `json:"items"`
}

func (list *List) FmtId(num interface{}) (err error) {
	switch v := num.(type) {
		case []uint8:
			list.Id = binary.BigEndian.Uint16(v)
			return
		case uint16: 
			list.Id = v
			return
		case uint64:
			id := uint16(v)
			list.Id = id
			return
		case int:
			list.Id = uint16(v)
		default:
			err = fmt.Errorf("Unaccepted type")
			return
	}
	return
}




func (list *List) Empty() (empty bool) {
	fmt.Printf("list %v, empty %v\n", list.Id, empty)
	return
}