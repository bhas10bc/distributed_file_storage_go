package main

import (
	"io"
	"log"
	"os"
)

func CASPathtranformFunc (key string) string {
	return ""
}

type PathTransformFunc func(string) string

var  DefaultPathTransformFunc = func(key string) string {
	return key
}

type StoreOpts struct {
	PathtransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) writeStream(key string, r io.Reader) error {

	var err error
	pathName := key

	err = os.Mkdir(pathName, os.ModePerm)
	if err != nil {
		return err
	}

	filename := "somefilename"

	pathAndFileName := pathName + "/" + filename

	f, err := os.Create(pathAndFileName)
	if err != nil {
		return err
	}
	n, err := io.Copy(f,r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, pathAndFileName)

	return nil
}