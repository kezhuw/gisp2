package gisp

import (
	"testing"

	p "github.com/Dwarfartisan/goparsec2"
)

func TestFloatParser0(t *testing.T) {
	data := "0.012"
	st := p.BasicStateFromText(data)
	o, err := FloatParser(&st)
	if err != nil {
		t.Fatalf("expect a Float but error %v", err)
	}
	if f, ok := o.(Float); ok {
		if f != Float(0.012) {
			t.Fatalf("expect a Float 0.012 but %v", f)
		}
	} else {
		t.Fatalf("expect Float but %v", o)
	}
}

func TestFloatParser1(t *testing.T) {
	data := "3.1415926"
	st := p.BasicStateFromText(data)
	o, err := FloatParser(&st)
	if err != nil {
		t.Fatalf("expect a Float but error %v", err)
	}
	if f, ok := o.(Float); ok {
		if f != Float(3.1415926) {
			t.Fatalf("expect a Float 3.1415926 but %v", f)
		}
	} else {
		t.Fatalf("expect Float but %v", o)
	}
}

func TestFloatParser2(t *testing.T) {
	data := "234.0"
	st := p.BasicStateFromText(data)
	o, err := FloatParser(&st)
	if err != nil {
		t.Fatalf("expect a Float but error %v", err)
	}
	if f, ok := o.(Float); ok {
		if f != Float(234) {
			t.Fatalf("expect a Float 234.0 but %v", f)
		}
	} else {
		t.Fatalf("expect Float but %v", o)
	}
}

func TestFloatParser3(t *testing.T) {
	data := ".5"
	st := p.BasicStateFromText(data)
	o, err := FloatParser(&st)
	if err != nil {
		t.Fatalf("expect a Float but error %v", err)
	}
	if f, ok := o.(Float); ok {
		if f != Float(0.5) {
			t.Fatalf("expect a Float 0.5 but %v", f)
		}
	} else {
		t.Fatalf("expect Float but %v", o)
	}
}

func TestFloatParser4(t *testing.T) {
	data := "f234.0"
	st := p.BasicStateFromText(data)
	o, err := FloatParser(&st)
	if err == nil {
		t.Fatalf("expect a Float parse error but got %v", o)
	}
}

func TestFloatParser5(t *testing.T) {
	data := "234"
	st := p.BasicStateFromText(data)
	o, err := FloatParser(&st)
	if err == nil {
		t.Fatalf("expect a Float parse error but got %v", o)
	}
}
