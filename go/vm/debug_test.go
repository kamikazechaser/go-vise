package vm

import (
	"testing"
)


func TestToString(t *testing.T) {
	var b []byte
	var r string
	var expect string
	var err error

	b = NewLine(nil, CATCH, []string{"xyzzy"}, []byte{0x0d}, []uint8{1})
	r, err = ToString(b)
	if err != nil {
		t.Fatal(err)
	}
	expect = "CATCH xyzzy 13 1 # invertmatch=true\n"
	if r != expect {
		t.Fatalf("expected:\n\t%v\ngot:\n\t%v", expect, r)
	}

	b = NewLine(nil, CROAK, nil, []byte{0x0d}, []uint8{1})
	r, err = ToString(b)
	if err != nil {
		t.Fatal(err)
	}
	expect = "CROAK 13 1 # invertmatch=true\n"
	if r != expect {
		t.Fatalf("expected:\n\t%v\ngot:\n\t%v", expect, r)
	}

	b = NewLine(nil, LOAD, []string{"foo"}, []byte{0x0a}, nil)
	r, err = ToString(b)
	if err != nil {
		t.Fatal(err)
	}
	expect = "LOAD foo 10\n"
	if r != expect {
		t.Fatalf("expected:\n\t%v\ngot:\n\t%v", expect, r)
	}

	b = NewLine(nil, RELOAD, []string{"bar"}, nil, nil)
	r, err = ToString(b)
	if err != nil {
		t.Fatal(err)
	}
	expect = "RELOAD bar\n"
	if r != expect {
		t.Fatalf("expected:\n\t%v\ngot:\n\t%v", expect, r)
	}

	b = NewLine(nil, MAP, []string{"inky_pinky"}, nil, nil) 
	r, err = ToString(b)
	if err != nil {
		t.Fatal(err)
	}
	expect = "MAP inky_pinky\n"
	if r != expect {
		t.Fatalf("expected:\n\t%v\ngot:\n\t%v", expect, r)
	}

	b = NewLine(nil, MOVE, []string{"blinky_clyde"}, nil, nil) 
	r, err = ToString(b)
	if err != nil {
		t.Fatal(err)
	}
	expect = "MOVE blinky_clyde\n"
	if r != expect {
		t.Fatalf("expected:\n\t%v\ngot:\n\t%v", expect, r)
	}

	b = NewLine(nil, HALT, nil, nil, nil) 
	r, err = ToString(b)
	if err != nil {
		t.Fatal(err)
	}
	expect = "HALT\n"
	if r != expect {
		t.Fatalf("expected:\n\t%v\ngot:\n\t%v", expect, r)
	}

	b = NewLine(nil, INCMP, []string{"13", "baz"}, nil, nil) 
	r, err = ToString(b)
	if err != nil {
		t.Fatal(err)
	}
	expect = "INCMP 13 baz\n"
	if r != expect {
		t.Fatalf("expected:\n\t%v\ngot:\n\t%v", expect, r)
	}
}

func TestToStringMultiple(t *testing.T) {
	b := NewLine(nil, INCMP, []string{"1", "foo"}, nil, nil)
	b = NewLine(b, INCMP, []string{"2", "bar"}, nil, nil)
	b = NewLine(b, CATCH, []string{"aiee"}, []byte{0x02, 0x9a}, []uint8{0})
	b = NewLine(b, LOAD, []string{"inky"}, []byte{0x2a}, nil)
	b = NewLine(b, HALT, nil, nil, nil)
	r, err := ToString(b)
	if err != nil {
		t.Fatal(err)
	}
	expect := `INCMP 1 foo
INCMP 2 bar
CATCH aiee 666 0 # invertmatch=false
LOAD inky 42
HALT
`
	if r != expect {
		t.Fatalf("expected:\n\t%v\ngot:\n\t%v", expect, r)
	}
}
