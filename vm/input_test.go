package vm

import (
	"context"
	"testing"

	"git.defalsify.org/vise.git/cache"
	"git.defalsify.org/vise.git/state"
)

func TestPhoneInput(t *testing.T) {
	v, err := ValidInput([]byte("+12345"))
	if err != nil {
		t.Fatal(err)
	}
	if v != -1 {
		t.Fatalf("expected -1, got %d", v)
	}
}

func TestMenuInputs(t *testing.T) {
	var err error
	_, err = ValidInput([]byte("0"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidInput([]byte("99"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidInput([]byte("foo"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidInput([]byte("foo Bar"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestFalseInput(t *testing.T) {
	_, err := ValidInput([]byte{0x0a})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTargetInput(t *testing.T) {
	var err error
	st := state.NewState(1)
	_, err = CheckTarget([]byte(""), st)
	if err == nil {
		t.Fatal("expected error")
	}
	_, err = CheckTarget([]byte("_"), st)
	if err == nil {
		t.Fatal("expected error")
	}
	st.Down("foo")
	v, err := CheckTarget([]byte("_"), st)
	if err != nil {
		t.Fatal(err)
	}
	if !v {
		t.Fatal("expected true")
	}
	v, err = CheckTarget([]byte("<"), st)
	if err != nil {
		t.Fatal(err)
	}
	if v {
		t.Fatal("expected false")
	}
	v, err = CheckTarget([]byte(">"), st)
	if err != nil {
		t.Fatal(err)
	}
	if !v {
		t.Fatal("expected true")
	}
	v, err = CheckTarget([]byte("%"), st)
	if err == nil {
		t.Fatal("expected error")
	}
	v, err = CheckTarget([]byte("foo"), st)
	if err != nil {
		t.Fatal(err)
	}
	if !v {
		t.Fatal("expected true")
	}
}

func TestApplyTarget(t *testing.T) {
	var err error
	ctx := context.Background()
	st := state.NewState(0)
	st.Down("root")
	st.Down("one")
	st.Down("two")
	ca := cache.NewCache()
	rs := newTestResource(st)
	rs.Lock()
	b := NewLine(nil, INCMP, []string{"^", "0"}, nil, nil)
	vm := NewVm(st, rs, ca, nil)

	st.SetInput([]byte("0"))
	b, err = vm.Run(ctx, b)
	if err != nil {
		t.Fatal(err)
	}

	st.Restart()
	st.Down("foo")
	b = NewLine(nil, INCMP, []string{"_", "0"}, nil, nil)
	vm = NewVm(st, rs, ca, nil)

	st.SetInput([]byte("0"))
	b, err = vm.Run(ctx, b)
	if err != nil {
		t.Fatal(err)
	}

	st.Restart()
	b = NewLine(nil, INCMP, []string{".", "0"}, nil, nil)
	vm = NewVm(st, rs, ca, nil)

	st.SetInput([]byte("0"))
	b, err = vm.Run(ctx, b)
	if err != nil {
		t.Fatal(err)
	}

	st.Restart()
	b = NewLine(nil, INCMP, []string{">", "0"}, nil, nil)
	vm = NewVm(st, rs, ca, nil)

	st.SetInput([]byte("0"))
	b, err = vm.Run(ctx, b)
	if err != nil {
		t.Fatal(err)
	}

	st.Restart()
	b = NewLine(nil, INCMP, []string{"<", "0"}, nil, nil)
	vm = NewVm(st, rs, ca, nil)

	st.SetInput([]byte("0"))
	b, err = vm.Run(ctx, b)
	if err != nil {
		t.Fatal(err)
	}
}

func TestVmCustomInputValid(t *testing.T) {
	var err error
	s := []byte{0x07, 0x6a, 0x6f, 0x6f}
	_, err = ValidInput(s)
	if err == nil {
		t.Fatal("expected error")
	}

	err = RegisterInputValidator(42, "^\x07[a-z]+")
	if err != nil {
		t.Fatal(err)
	}

	v, err := ValidInput(s)
	if err != nil {
		t.Fatal(err)
	}
	if v != 42 {
		t.Fatalf("expected 42, got %d", v)
	}
}
