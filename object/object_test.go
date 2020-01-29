package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := NewStringObject("Hello World")
	hello2 := NewStringObject("Hello World")
	diff1 := NewStringObject("My name is johnny")
	diff2 := NewStringObject("My name is johnny")

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	true1 := &Boolean{Value: true}
	true2 := &Boolean{Value: true}
	false1 := &Boolean{Value: false}
	false2 := &Boolean{Value: false}

	if true1.HashKey() != true2.HashKey() {
		t.Errorf("trues do not have same hash key")
	}

	if false1.HashKey() != false2.HashKey() {
		t.Errorf("falses do not have same hash key")
	}

	if true1.HashKey() == false1.HashKey() {
		t.Errorf("true has same hash key as false")
	}
}

func TestIntegerHashKey(t *testing.T) {
	one1 := &Integer{Value: 1}
	one2 := &Integer{Value: 1}
	two1 := &Integer{Value: 2}
	two2 := &Integer{Value: 2}

	if one1.HashKey() != one2.HashKey() {
		t.Errorf("integers with same content have twoerent hash keys")
	}

	if two1.HashKey() != two2.HashKey() {
		t.Errorf("integers with same content have twoerent hash keys")
	}

	if one1.HashKey() == two1.HashKey() {
		t.Errorf("integers with twoerent content have same hash keys")
	}
}

func TestHashPairEqual(t *testing.T) {
	hp1 := &HashPair{NewStringObject("abc"), NewStringObject("xyz")}
	hp2 := &HashPair{NewStringObject("abc"), NewStringObject("xyz")}
	hp3 := &HashPair{NewStringObject("foo"), NewStringObject("bar")}

	if hp1.Key.(*String).Compare(hp2.Key.(*String)) != 0 {
		t.Errorf("the same string compare got unequal.")
	}

	if hp1.Key.(*String).Compare(hp3.Key.(*String)) == 0 {
		t.Errorf("the different string compare got equal.")
	}
}
