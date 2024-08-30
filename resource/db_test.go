package resource

import (
	"bytes"
	"context"
	"testing"

	"git.defalsify.org/vise.git/db"
)

func TestDb(t *testing.T) {
	ctx := context.Background()
	store := db.NewMemDb(ctx)
	tg, err := NewDbFuncGetter(store, db.DATATYPE_TEMPLATE)
	if err != nil {
		t.Fatal(err)
	}
	rs := NewMenuResource()
	rs.WithTemplateGetter(tg.GetTemplate)

	s, err := rs.GetTemplate(ctx, "foo")
	if err == nil {
		t.Fatal("expected error")
	}

	store.SetPrefix(db.DATATYPE_TEMPLATE)
	err = store.Put(ctx, []byte("foo"), []byte("bar"))
	if err != nil {
		t.Fatal(err)
	}
	s, err = rs.GetTemplate(ctx, "foo")
	if err != nil {
		t.Fatal(err)
	}
	if s != "bar" {
		t.Fatalf("expected 'bar', got %s", s)
	}

	// test support check
	store.SetPrefix(db.DATATYPE_BIN)
	err = store.Put(ctx, []byte("xyzzy"), []byte("deadbeef"))
	if err != nil {
		t.Fatal(err)
	}

	rs.WithCodeGetter(tg.GetCode)
	b, err := rs.GetCode(ctx, "xyzzy")
	if err == nil {
		t.Fatal("expected error")
	}

	tg, err = NewDbFuncGetter(store, db.DATATYPE_TEMPLATE, db.DATATYPE_BIN)
	if err != nil {
		t.Fatal(err)
	}
	rs.WithTemplateGetter(tg.GetTemplate)

	rs.WithCodeGetter(tg.GetCode)
	b, err = rs.GetCode(ctx, "xyzzy")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, []byte("deadbeef")) {
		t.Fatalf("expected 'deadbeef', got %x", b)
	}

	tg, err = NewDbFuncGetter(store, db.DATATYPE_TEMPLATE, db.DATATYPE_BIN, db.DATATYPE_MENU)
	if err != nil {
		t.Fatal(err)
	}
	store.SetPrefix(db.DATATYPE_MENU)
	err = store.Put(ctx, []byte("inky"), []byte("pinky"))
	if err != nil {
		t.Fatal(err)
	}
	rs.WithMenuGetter(tg.GetMenu)

}
