package types

import (
	"strconv"
	"testing"
)

type PSet struct {
	foo int
}

type VSet struct {
	foo int
}

func (s *PSet) Set(string) error {
	return nil
}

func (s VSet) Set(string) error {
	return nil
}

func TestSetTstClr(t *testing.T) {
	var tp TypeId
	testCases := []struct {
		set func()
		tst func() bool
		clr func()
	}{
		{tp.SetBoolBit, tp.TstBoolBit, tp.ClrBoolBit},
		{tp.SetIntBit, tp.TstIntBit, tp.ClrIntBit},
		{tp.SetUintBit, tp.TstUintBit, tp.ClrUintBit},
		{tp.SetFloatBit, tp.TstFloatBit, tp.ClrFloatBit},
		{tp.SetStringBit, tp.TstStringBit, tp.ClrStringBit},
		{tp.SetSliceBit, tp.TstSliceBit, tp.ClrSliceBit},
		{tp.SetPointerBit, tp.TstPointerBit, tp.ClrPointerBit},
		{tp.SetSetterBit, tp.TstSetterBit, tp.ClrSetterBit},
		{tp.SetOtherBit, tp.TstOtherBit, tp.ClrOtherBit},
		{tp.SetIntBit, tp.TstAnyNumBit, tp.ClrIntBit},
		{tp.SetUintBit, tp.TstAnyNumBit, tp.ClrUintBit},
		{tp.SetFloatBit, tp.TstAnyNumBit, tp.ClrFloatBit},
	}
	for _, fn := range testCases {
		fn.set()
		test := fn.tst()
		if !test {
			t.Errorf("Set/Tst failed %016b, result: %t", tp, test)
		}
		fn.clr()
		test = fn.tst()
		if test {
			t.Errorf("Clr/Tst failed %016b, result: %t", tp, test)
		}
	}
}

func boolSlicesAreEqual(a, b []bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func allBits(ix interface{}) []bool {
	ti := Type(ix)
	return []bool{
		ti.TstBoolBit(),
		ti.TstIntBit(),
		ti.TstUintBit(),
		ti.TstFloatBit(),
		ti.TstStringBit(),
		ti.TstSliceBit(),
		ti.TstPointerBit(),
		ti.TstSetterBit(),
		ti.TstOtherBit(),
		ti.TstAnyNumBit(),
	}
}

func TestAllBits(t *testing.T) {
	testCases := []struct {
		ix     interface{}
		result []bool
		nbits  int
		slilen int
	}{
		//                bool    int   uint  float    str  slice    ptr    set    oth    any   nb  ns
		{true, []bool{true, false, false, false, false, false, false, false, false, false}, 0, -1},
		{new(bool), []bool{true, false, false, false, false, false, true, false, false, false}, 0, -1},
		{[]bool{}, []bool{true, false, false, false, false, true, false, false, false, false}, 0, 0},
		{&([]bool{}), []bool{true, false, false, false, false, true, true, false, false, false}, 0, 0},
		//                        bool    int   uint  float    str  slice    ptr    oth    any
		{int(1), []bool{false, true, false, false, false, false, false, false, false, true}, strconv.IntSize, -1},
		{new(int), []bool{false, true, false, false, false, false, true, false, false, true}, strconv.IntSize, -1},
		{[]int{}, []bool{false, true, false, false, false, true, false, false, false, true}, strconv.IntSize, 0},
		{&([]int{5}), []bool{false, true, false, false, false, true, true, false, false, true}, strconv.IntSize, 1},
		{int8(1), []bool{false, true, false, false, false, false, false, false, false, true}, 8, -1},
		{new(int8), []bool{false, true, false, false, false, false, true, false, false, true}, 8, -1},
		{[]int8{}, []bool{false, true, false, false, false, true, false, false, false, true}, 8, 0},
		{&([]int8{}), []bool{false, true, false, false, false, true, true, false, false, true}, 8, 0},
		{int16(1), []bool{false, true, false, false, false, false, false, false, false, true}, 16, -1},
		{new(int16), []bool{false, true, false, false, false, false, true, false, false, true}, 16, -1},
		{[]int16{}, []bool{false, true, false, false, false, true, false, false, false, true}, 16, 0},
		{&([]int16{}), []bool{false, true, false, false, false, true, true, false, false, true}, 16, 0},
		{int32(1), []bool{false, true, false, false, false, false, false, false, false, true}, 32, -1},
		{new(int32), []bool{false, true, false, false, false, false, true, false, false, true}, 32, -1},
		{[]int32{}, []bool{false, true, false, false, false, true, false, false, false, true}, 32, 0},
		{&([]int32{}), []bool{false, true, false, false, false, true, true, false, false, true}, 32, 0},
		{int64(1), []bool{false, true, false, false, false, false, false, false, false, true}, 64, -1},
		{new(int64), []bool{false, true, false, false, false, false, true, false, false, true}, 64, -1},
		{[]int64{}, []bool{false, true, false, false, false, true, false, false, false, true}, 64, 0},
		{&([]int64{}), []bool{false, true, false, false, false, true, true, false, false, true}, 64, 0},
		//                      bool    int   uint  float    str  slice    ptr    oth    any
		{uint(1), []bool{false, false, true, false, false, false, false, false, false, true}, strconv.IntSize, -1},
		{new(uint), []bool{false, false, true, false, false, false, true, false, false, true}, strconv.IntSize, -1},
		{[]uint{}, []bool{false, false, true, false, false, true, false, false, false, true}, strconv.IntSize, 0},
		{&([]uint{}), []bool{false, false, true, false, false, true, true, false, false, true}, strconv.IntSize, 0},
		{uint8(1), []bool{false, false, true, false, false, false, false, false, false, true}, 8, -1},
		{new(uint8), []bool{false, false, true, false, false, false, true, false, false, true}, 8, -1},
		{[]uint8{}, []bool{false, false, true, false, false, true, false, false, false, true}, 8, 0},
		{&([]uint8{}), []bool{false, false, true, false, false, true, true, false, false, true}, 8, 0},
		{uint16(1), []bool{false, false, true, false, false, false, false, false, false, true}, 16, -1},
		{new(uint16), []bool{false, false, true, false, false, false, true, false, false, true}, 16, -1},
		{[]uint16{}, []bool{false, false, true, false, false, true, false, false, false, true}, 16, 0},
		{&([]uint16{}), []bool{false, false, true, false, false, true, true, false, false, true}, 16, 0},
		{uint32(1), []bool{false, false, true, false, false, false, false, false, false, true}, 32, -1},
		{new(uint32), []bool{false, false, true, false, false, false, true, false, false, true}, 32, -1},
		{[]uint32{}, []bool{false, false, true, false, false, true, false, false, false, true}, 32, 0},
		{&([]uint32{}), []bool{false, false, true, false, false, true, true, false, false, true}, 32, 0},
		{uint64(1), []bool{false, false, true, false, false, false, false, false, false, true}, 64, -1},
		{new(uint64), []bool{false, false, true, false, false, false, true, false, false, true}, 64, -1},
		{[]uint64{}, []bool{false, false, true, false, false, true, false, false, false, true}, 64, 0},
		{&([]uint64{}), []bool{false, false, true, false, false, true, true, false, false, true}, 64, 0},
		//                        bool    int   uint  float    str  slice    ptr    oth    any
		{float32(1), []bool{false, false, false, true, false, false, false, false, false, true}, 32, -1},
		{new(float32), []bool{false, false, false, true, false, false, true, false, false, true}, 32, -1},
		{[]float32{}, []bool{false, false, false, true, false, true, false, false, false, true}, 32, 0},
		{&([]float32{}), []bool{false, false, false, true, false, true, true, false, false, true}, 32, 0},
		{float64(1), []bool{false, false, false, true, false, false, false, false, false, true}, 64, -1},
		{new(float64), []bool{false, false, false, true, false, false, true, false, false, true}, 64, -1},
		{[]float64{}, []bool{false, false, false, true, false, true, false, false, false, true}, 64, 0},
		{&([]float64{}), []bool{false, false, false, true, false, true, true, false, false, true}, 64, 0},
		//                       bool    int   uint  float    str  slice    ptr    oth    any
		{string(""), []bool{false, false, false, false, true, false, false, false, false, false}, 0, -1},
		{new(string), []bool{false, false, false, false, true, false, true, false, false, false}, 0, -1},
		{[]string{}, []bool{false, false, false, false, true, true, false, false, false, false}, 0, 0},
		{&([]string{}), []bool{false, false, false, false, true, true, true, false, false, false}, 0, 0},
		//                        bool    int   uint  float    str  slice    ptr    oth    any
		{struct{}{}, []bool{false, false, false, false, false, false, false, false, true, false}, 0, -1},
		{&(struct{}{}), []bool{false, false, false, false, false, false, true, false, true, false}, 0, -1},
		{[]struct{}{}, []bool{false, false, false, false, false, true, false, false, true, false}, 0, -1},
		{&([]struct{}{}), []bool{false, false, false, false, false, true, true, false, true, false}, 0, -1},
		// Setter                bool    int   uint  float    str  slice    ptr    set   oth    any
		{PSet{}, []bool{false, false, false, false, false, false, false, false, true, false}, 0, -1},
		{&PSet{}, []bool{false, false, false, false, false, false, true, true, false, false}, 0, -1},
		{[]PSet{}, []bool{false, false, false, false, false, true, false, false, true, false}, 0, -1},
		{[]*PSet{{}}, []bool{false, false, false, false, false, true, false, false, true, false}, 0, -1},
		{VSet{}, []bool{false, false, false, false, false, false, false, true, false, false}, 0, -1},
		{&VSet{}, []bool{false, false, false, false, false, false, true, true, false, false}, 0, -1},
		{[]VSet{}, []bool{false, false, false, false, false, true, false, false, true, false}, 0, -1},
		{[]*VSet{{}}, []bool{false, false, false, false, false, true, false, false, true, false}, 0, -1},
	}
	for i, row := range testCases {
		ans := allBits(row.ix)
		if !boolSlicesAreEqual(ans, row.result) {
			t.Errorf("row %d T<%16T> failed bools (%016b)", i, row.ix, Type(row.ix))
		}
		if BitSize(row.ix) != row.nbits {
			t.Errorf("row %d T<%16T> failed nbits (%016b)", i, row.ix, Type(row.ix))
		}
		if SliceLen(row.ix) != row.slilen {
			t.Errorf("row %d T<%16T> failed slice length (%016b), %d != %d", i, row.ix,
				Type(row.ix), SliceLen(row.ix), row.slilen)
		}
	}
}

func TestIsPointerTo(t *testing.T) {
	var a, b bool
	var sa, sb []bool
	var i, j int

	if IsPointerTo(a, b) {
		t.Errorf("got success, expected failure")
	}
	if IsPointerTo(a, i) {
		t.Errorf("got success, expected failure")
	}
	if IsPointerTo(i, sa) {
		t.Errorf("got success, expected failure")
	}
	if IsPointerTo(&i, &j) {
		t.Errorf("got success, expected failure")
	}
	if !IsPointerTo(&a, b) {
		t.Errorf("got failure, expected success")
	}
	if !IsPointerTo(&sa, sb) {
		t.Errorf("got failure, expected success")
	}
	if !IsPointerTo(&i, j) {
		t.Errorf("got failure, expected success")
	}
}
