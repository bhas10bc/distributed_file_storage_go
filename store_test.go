package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathtransformFunc:  CASPathtranformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpg bytes"))

	err := s.writeStream("myImages", data)
	if err != nil {
		t.Error(err)
	}
}

func TestPathTransformFunc(t *testing.T){
	key := "mypersonalImage"
	pathname := CASPathtranformFunc(key)

	fmt.Println(pathname)
}