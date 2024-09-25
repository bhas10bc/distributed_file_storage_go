package main

import (
	"bytes"
	"testing"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathtransformFunc: DefaultPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpg bytes"))

	err := s.writeStream("myImage", data)
	if err != nil {
		t.Error(err)
	}
}

func TestPathTransformFunc(t *testing.T){
	
}