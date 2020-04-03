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

func (list *List) FmtId(num interface{}, errc chan error) {
	// fmt.Printf("in FmtId %T\n", num)
	switch v := num.(type) {
		case []uint8:
			list.Id = binary.BigEndian.Uint16(v)
		case uint16: 
			list.Id = v
		case uint64:
			id := uint16(v)
			list.Id = id
		case int:
			list.Id = uint16(v)
		default:
			fmt.Printf("num val %v\n", v)
			
			errc <- fmt.Errorf("Unaccepted type %v\n", v)
			return
	}
	return
}




func (list *List) Empty() (empty bool) {
	fmt.Printf("list %v, empty %v\n", list.Id, empty)
	return
}