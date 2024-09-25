package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

func CASPathtranformFunc (key string) string {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashStr)/blocksize

	paths := make ([]string, sliceLen)

	for i := 0; i < sliceLen; i ++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}

	return strings.Join(paths, "/")
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
	pathName := s.PathtransformFunc(key)
	var err error
	err = os.MkdirAll(pathName, os.ModePerm)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, r)

	fileNameBytes := md5.Sum(buf.Bytes())
	filename := hex.EncodeToString(fileNameBytes[:])
	pathAndFileName := pathName + "/" + filename

	f, err := os.Create(pathAndFileName)
	if err != nil {
		return err
	}
	n, err := io.Copy(f,buf)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, pathAndFileName)

	return nil
}