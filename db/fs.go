package db

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type FsDb struct {
	BaseDb
	dir string
}

func(fdb *FsDb) Connect(ctx context.Context, connStr string) error {
	fi, err := os.Stat(connStr)
	if err != nil {
		return err
	}
	if !fi.IsDir()  {
		return fmt.Errorf("fs db %s is not a directory", connStr)
	}
	fdb.dir = connStr
	return nil
}

func(fdb *FsDb) Get(ctx context.Context, key []byte) ([]byte, error) {
	fp, err := fdb.pathFor(key)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func(fdb *FsDb) Put(ctx context.Context, key []byte, val []byte) error {
	fp, err := fdb.pathFor(key)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fp, val, 0600)
}

func(fdb *FsDb) Close() error {
	return nil
}	
 
func(fdb *FsDb) pathFor(key []byte) (string, error) {
	kb, err := fdb.ToKey(key)
	if err != nil {
		return "", err
	}
	kb[0] += 30
	return path.Join(fdb.dir, string(kb)), nil
}
