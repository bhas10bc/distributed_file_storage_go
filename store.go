package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var defaultRootFolderName = "myFolder"

func CASPathtranformFunc (key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashStr)/blocksize

	paths := make ([]string, sliceLen)

	for i := 0; i < sliceLen; i ++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName : strings.Join(paths, "/"),
		FileName: hashStr,
	}
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	PathName string
	FileName string
}

func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}

var  DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

type StoreOpts struct {
	Root string
	PathtransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathtransformFunc == nil {
		opts.PathtransformFunc = DefaultPathTransformFunc
	}

	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string) bool {
	pathKey :=s.PathtransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	_, err := os.Stat(fullPathWithRoot)
	
	 return !errors.Is(err, os.ErrNotExist)
}

func(s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathtransformFunc(key)
	firstPathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FirstPathName())
	err := os.RemoveAll(firstPathNameWithRoot)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Read(key string) (io.Reader, error){
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf,err
}


func(s *Store) readStream(key string) (io.ReadCloser, error) {
	PathKey := s.PathtransformFunc(key)
	PathKeyWithRoot := fmt.Sprintf("%s/%s", s.Root, PathKey.FullPath())
	return os.Open(PathKeyWithRoot)
	
}

func(s *Store) Write(key string, r io.Reader) error {
	return s.writeStream(key,r)
}

func (s *Store) writeStream(key string, r io.Reader) error {
	PathKey := s.PathtransformFunc(key)
	PathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, PathKey.PathName)
	var err error
	err = os.MkdirAll(PathNameWithRoot, os.ModePerm)
	if err != nil {
		return err
	}


	FullPath := PathKey.FullPath()
	FullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, FullPath)
	f, err := os.Create(FullPathWithRoot)
	if err != nil {
		return err
	}
	n, err := io.Copy(f,r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, FullPathWithRoot)

	return nil
}