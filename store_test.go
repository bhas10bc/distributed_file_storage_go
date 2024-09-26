package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathtransformFunc:  CASPathtranformFunc,
	}
	s := NewStore(opts)
	key := "myImages2"
	storeData := []byte("some more 2 jpg bytes")

	err := s.writeStream(key, bytes.NewReader(storeData))
	if err != nil {
		t.Error(err)
	}

	ok := s.Has(key)
	if !ok {
		t.Errorf("expected to have key %s", key)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	retrievedData , err := ioutil.ReadAll(r)

	if string(storeData) != string(retrievedData){
		t.Errorf("want %s but we have %s", storeData, retrievedData)
	}

	


}

func TestPathTransformFunc(t *testing.T){
	key := "mypersonalImage"
	pathname := CASPathtranformFunc(key)

	fmt.Println(pathname)
}

func TestStoreDeleteKey(t *testing.T){
	opts := StoreOpts{
		PathtransformFunc:  CASPathtranformFunc,
	}
	s := NewStore(opts)
	key := "myImages2"
	// storeData := []byte("deletable_data")

	// err := s.writeStream(key, bytes.NewReader(storeData))
	// if err != nil {
	// 	t.Error(err)
	// }


	err := s.Delete(key)
	if err != nil {
		t.Error(err)
	}

	ok := s.Has(key)
	if ok {
		t.Errorf("expected to have key %s", key)
	}


}