package fs

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"git.defalsify.org/vise.git/db"
)

// holds string (filepath) versions of LookupKey
type fsLookupKey struct {
	Default     string
	Translation string
}

// pure filesystem backend implementation if the Db interface.
type fsDb struct {
	*db.DbBase
	dir         string
	elements    []os.DirEntry
	matchPrefix []byte
	binary      bool
}

// NewFsDb creates a filesystem backed Db implementation.
func NewFsDb() *fsDb {
	db := &fsDb{
		DbBase: db.NewDbBase(),
	}
	return db
}

func (fdb *fsDb) WithBinary() *fsDb {
	fdb.binary = true
	return fdb
}

// String implements the string interface.
func (fdb *fsDb) String() string {
	return "fsdb: " + fdb.dir
}

// Connect implements the Db interface.
func (fdb *fsDb) Connect(ctx context.Context, connStr string) error {
	if fdb.dir != "" {
		logg.WarnCtxf(ctx, "already connected", "conn", fdb.dir)
		return nil
	}
	err := os.MkdirAll(connStr, 0700)
	if err != nil {
		return err
	}
	fdb.DbBase.Connect(ctx, connStr)
	fdb.dir = fdb.Connection()
	return nil
}

// ToKey overrides the BaseDb implementation, creating a base64 string
// if binary keys have been enabled
func (fdb *fsDb) ToKey(ctx context.Context, key []byte) (db.LookupKey, error) {
	if fdb.binary {
		s := base64.StdEncoding.EncodeToString(key)
		key = []byte(s)
	}
	return fdb.DbBase.ToKey(ctx, key)
}

func (fdb *fsDb) DecodeKey(ctx context.Context, key []byte) ([]byte, error) {
	key, err := fdb.DbBase.DecodeKey(ctx, key)
	if err != nil {
		return nil, err
	}
	if !fdb.binary {
		return key, nil
	}
	oldKey := key
	key, err = base64.StdEncoding.DecodeString(string(key))
	if err != nil {
		return []byte{}, fmt.Errorf("base64 decode error '%s': %v", oldKey, err)
	}
	logg.TraceCtxf(ctx, "decoding base64 key", "base64", oldKey, "bin", key)
	return key, nil
}

// Get implements the Db interface.
func (fdb *fsDb) Get(ctx context.Context, key []byte) ([]byte, error) {
	var f *os.File
	lk, err := fdb.ToKey(ctx, key)
	if err != nil {
		return nil, err
	}
	flk, err := fdb.pathFor(ctx, &lk)
	if err != nil {
		return nil, err
	}
	flka, err := fdb.altPathFor(ctx, &lk)
	if err != nil {
		return nil, err
	}
	for i, fp := range []string{flk.Translation, flka.Translation, flk.Default, flka.Default} {
		if fp == "" {
			logg.TraceCtxf(ctx, "fs get skip missing", "i", i)
			continue
		}
		logg.TraceCtxf(ctx, "trying fs get", "i", i, "key", key, "path", fp)
		f, err = os.Open(fp)
		if err == nil {
			break
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	}
	if f == nil {
		return nil, db.NewErrNotFound(key)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Put implements the Db interface.
func (fdb *fsDb) Put(ctx context.Context, key []byte, val []byte) error {
	if !fdb.CheckPut() {
		return errors.New("unsafe put and safety set")
	}
	lk, err := fdb.ToKey(ctx, key)
	if err != nil {
		return err
	}
	flk, err := fdb.pathFor(ctx, &lk)
	if err != nil {
		return err
	}
	logg.TraceCtxf(ctx, "fs put", "key", key, "lk", lk, "flk", flk, "val", val)
	if flk.Translation != "" {
		err = ioutil.WriteFile(flk.Translation, val, 0600)
		if err != nil {
			return err
		}
		return nil
	}
	return ioutil.WriteFile(flk.Default, val, 0600)
}

// Close implements the Db interface.
func (fdb *fsDb) Close(ctx context.Context) error {
	return nil
}

// create a key safe for the filesystem.
func (fdb *fsDb) pathFor(ctx context.Context, lk *db.LookupKey) (fsLookupKey, error) {
	var flk fsLookupKey
	lk.Default[0] += 0x30
	flk.Default = path.Join(fdb.dir, string(lk.Default))
	if lk.Translation != nil {
		lk.Translation[0] += 0x30
		flk.Translation = path.Join(fdb.dir, string(lk.Translation))
	}
	return flk, nil
}

// create a key safe for the filesystem, matching legacy resource.FsResource name.
func (fdb *fsDb) altPathFor(ctx context.Context, lk *db.LookupKey) (fsLookupKey, error) {
	var flk fsLookupKey
	fb := string(lk.Default[1:])
	if fdb.Prefix() == db.DATATYPE_BIN {
		fb += ".bin"
	}
	flk.Default = path.Join(fdb.dir, fb)

	if lk.Translation != nil {
		fb = string(lk.Translation[1:])
		if fdb.Prefix() == db.DATATYPE_BIN {
			fb += ".bin"
		}
		flk.Translation = path.Join(fdb.dir, fb)
	}

	return flk, nil
}
