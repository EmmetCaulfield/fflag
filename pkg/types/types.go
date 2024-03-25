package types

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type SetValue interface {
	Set(string) error
}

type TypeId uint16

const (
	NumBits  TypeId = 0b0000000000000111
	Bits8    TypeId = 0x0001
	Bits16   TypeId = 0x0002
	Bits32   TypeId = 0x0003
	Bits64   TypeId = 0x0004
	BoolT    TypeId = 0b0000000000001000
	IntT     TypeId = 0b0000000000010000
	UintT    TypeId = 0b0000000000100000
	FloatT   TypeId = 0b0000000001000000
	StringT  TypeId = 0b0000000010000000
	SliceT   TypeId = 0b0000000100000000
	PointerT TypeId = 0b0000001000000000
	SetterT  TypeId = 0b0100000000000000
	OtherT   TypeId = 0b1000000000000000
)

func (tp *TypeId) SetBoolBit()    { *tp = *tp | BoolT }
func (tp *TypeId) SetIntBit()     { *tp = *tp | IntT }
func (tp *TypeId) SetUintBit()    { *tp = *tp | UintT }
func (tp *TypeId) SetFloatBit()   { *tp = *tp | FloatT }
func (tp *TypeId) SetStringBit()  { *tp = *tp | StringT }
func (tp *TypeId) SetSliceBit()   { *tp = *tp | SliceT }
func (tp *TypeId) SetPointerBit() { *tp = *tp | PointerT }
func (tp *TypeId) SetSetterBit()  { *tp = *tp | SetterT }
func (tp *TypeId) SetOtherBit()   { *tp = *tp | OtherT }

func (tp *TypeId) ClrBoolBit()    { *tp = *tp & ^BoolT }
func (tp *TypeId) ClrIntBit()     { *tp = *tp & ^IntT }
func (tp *TypeId) ClrUintBit()    { *tp = *tp & ^UintT }
func (tp *TypeId) ClrFloatBit()   { *tp = *tp & ^FloatT }
func (tp *TypeId) ClrStringBit()  { *tp = *tp & ^StringT }
func (tp *TypeId) ClrSliceBit()   { *tp = *tp & ^SliceT }
func (tp *TypeId) ClrPointerBit() { *tp = *tp & ^PointerT }
func (tp *TypeId) ClrSetterBit()  { *tp = *tp & ^SetterT }
func (tp *TypeId) ClrOtherBit()   { *tp = *tp & ^OtherT }

func (tp *TypeId) TstBoolBit() bool    { return *tp&BoolT != 0 }
func (tp *TypeId) TstIntBit() bool     { return *tp&IntT != 0 }
func (tp *TypeId) TstUintBit() bool    { return *tp&UintT != 0 }
func (tp *TypeId) TstFloatBit() bool   { return *tp&FloatT != 0 }
func (tp *TypeId) TstStringBit() bool  { return *tp&StringT != 0 }
func (tp *TypeId) TstSliceBit() bool   { return *tp&SliceT != 0 }
func (tp *TypeId) TstPointerBit() bool { return *tp&PointerT != 0 }
func (tp *TypeId) TstSetterBit() bool  { return *tp&SetterT != 0 }
func (tp *TypeId) TstOtherBit() bool   { return *tp&OtherT != 0 }
func (tp *TypeId) TstAnyNumBit() bool  { return *tp&IntT != 0 || *tp&UintT != 0 || *tp&FloatT != 0 }

// Returns true if two types have the same underlying basic type
func SameBaseType(a, b TypeId) bool {
	mask := ^(PointerT | SliceT)
	return a&mask == b&mask
}

// Returns the number of bits or zero if not applicable
func (tp *TypeId) BitSize() int {
	n := *tp & NumBits
	if n == 0 {
		return 0
	}
	return 8 * (1 << (n - 1))
}

// Returns the NumBits setting for the default integer size
func IntBits() TypeId {
	nbytes := strconv.IntSize >> 3
	var n uint8
	for n = 0; nbytes != 0; nbytes = (nbytes >> 1) {
		n++
	}
	return TypeId(n)
}

// Returns a TypeId corresponding to the type within the given interface
func Type(ix interface{}) TypeId {
	if ix == nil {
		return TypeId(0)
	}
	switch ix.(type) {
	// Boolean
	case bool:
		return BoolT
	case *bool:
		return PointerT | BoolT
	case []bool:
		return SliceT | BoolT
	case *[]bool:
		return PointerT | SliceT | BoolT

	// Unsigned integers
	case uint:
		return IntBits() | UintT
	case *uint:
		return PointerT | IntBits() | UintT
	case []uint:
		return SliceT | IntBits() | UintT
	case *[]uint:
		return PointerT | SliceT | IntBits() | UintT

	case uint8: // also `byte`
		return Bits8 | UintT
	case *uint8:
		return PointerT | Bits8 | UintT
	case []uint8:
		return SliceT | Bits8 | UintT
	case *[]uint8:
		return PointerT | SliceT | Bits8 | UintT

	case uint16:
		return Bits16 | UintT
	case *uint16:
		return PointerT | Bits16 | UintT
	case []uint16:
		return SliceT | Bits16 | UintT
	case *[]uint16:
		return PointerT | SliceT | Bits16 | UintT

	case uint32:
		return Bits32 | UintT
	case *uint32:
		return PointerT | Bits32 | UintT
	case []uint32:
		return SliceT | Bits32 | UintT
	case *[]uint32:
		return PointerT | SliceT | Bits32 | UintT

	case uint64:
		return Bits64 | UintT
	case *uint64:
		return PointerT | Bits64 | UintT
	case []uint64:
		return SliceT | Bits64 | UintT
	case *[]uint64:
		return PointerT | SliceT | Bits64 | UintT

	// Signed Integers
	case int:
		return IntBits() | IntT
	case *int:
		return PointerT | IntBits() | IntT
	case []int:
		return SliceT | IntBits() | IntT
	case *[]int:
		return PointerT | SliceT | IntBits() | IntT

	case int8:
		return Bits8 | IntT
	case *int8:
		return PointerT | Bits8 | IntT
	case []int8:
		return SliceT | Bits8 | IntT
	case *[]int8:
		return PointerT | SliceT | Bits8 | IntT

	case int16:
		return Bits16 | IntT
	case *int16:
		return PointerT | Bits16 | IntT
	case []int16:
		return SliceT | Bits16 | IntT
	case *[]int16:
		return PointerT | SliceT | Bits16 | IntT

	case int32: // also `rune`
		return Bits32 | IntT
	case *int32:
		return PointerT | Bits32 | IntT
	case []int32:
		return SliceT | Bits32 | IntT
	case *[]int32:
		return PointerT | SliceT | Bits32 | IntT

	case int64:
		return Bits64 | IntT
	case *int64:
		return PointerT | Bits64 | IntT
	case []int64:
		return SliceT | Bits64 | IntT
	case *[]int64:
		return PointerT | SliceT | Bits64 | IntT

	// Floating-point types
	case float32:
		return Bits32 | FloatT
	case *float32:
		return PointerT | Bits32 | FloatT
	case []float32:
		return SliceT | Bits32 | FloatT
	case *[]float32:
		return PointerT | SliceT | Bits32 | FloatT

	case float64:
		return Bits64 | FloatT
	case *float64:
		return PointerT | Bits64 | FloatT
	case []float64:
		return SliceT | Bits64 | FloatT
	case *[]float64:
		return PointerT | SliceT | Bits64 | FloatT

	case string:
		return StringT
	case *string:
		return PointerT | StringT
	case []string:
		return SliceT | StringT
	case *[]string:
		return PointerT | SliceT | StringT

	}

	// The only useful thing we can do is tell whether the thing
	// behind the interface `ix` implements the SetValue interface. We
	// don't get to determine how it's implemented, so whether it's a
	// pointer or a slice or whatever is useless.
	//
	// If we're going to start looking into lists of pointers or
	// pointers to lists of pointers and try to handle them, we'd have
	// to do it for everything and it's not necessary to do what we
	// want to do here: we don't have to broaden the interface to
	// admit absolutely everthing.
	var typeId TypeId
	if _, ok := ix.(SetValue); ok {
		typeId.SetSetterBit()
	} else {
		typeId.SetOtherBit()
	}
	td := fmt.Sprintf("%T", ix)
	// if len(td) > 4 && td[0:4] == "*[]*" {
	//		return typeId | PointerT | SliceT
	//	}
	if len(td) > 2 && td[0:3] == "*[]" {
		return typeId | PointerT | SliceT
	}
	//	if len(td) > 2 && td[0:3] == "[]*" {
	//		return typeId | PointerT | SliceT
	//	}
	if len(td) > 1 && td[0:2] == "[]" {
		return typeId | SliceT
	}
	if len(td) > 0 && td[0] == '*' {
		return typeId | PointerT
	}
	return typeId
}

func IsNum(ix interface{}) bool {
	typeId := Type(ix)
	return typeId.TstAnyNumBit()
}

func IsInt(ix interface{}) bool {
	typeId := Type(ix)
	return typeId.TstIntBit()
}

func IsUint(ix interface{}) bool {
	typeId := Type(ix)
	return typeId.TstUintBit()
}

func IsFloat(ix interface{}) bool {
	typeId := Type(ix)
	return typeId.TstFloatBit()
}

func IsString(ix interface{}) bool {
	typeId := Type(ix)
	return typeId.TstStringBit()
}

func IsOtherT(ix interface{}) bool {
	typeId := Type(ix)
	return typeId.TstOtherBit()
}

func IsSetter(ix interface{}) bool {
	typeId := Type(ix)
	return typeId.TstSetterBit()
}

func IsPointer(ix interface{}) bool {
	typeId := Type(ix)
	return typeId.TstPointerBit()
}

func IsSlice(ix interface{}) bool {
	typeId := Type(ix)
	return typeId.TstSliceBit()
}

// Returns the number of bits or zero if not applicable
func BitSize(ix interface{}) int {
	typeId := Type(ix)
	return typeId.BitSize()
}

// Returns the length of the underlying slice or -1 if not applicable
func SliceLen(ix interface{}) int {
	if ix == nil {
		return -1
	}
	// It seems that there's no way of saying:
	//
	//     if v, ok := ix.([]interface{}); ok { ... }
	//     case []interface{}:
	//     if v, ok := ix.([]any); ok { ... }
	//     case []any:
	switch v := ix.(type) {
	case []bool:
		return len(v)
	case *[]bool:
		return len(*v)
	case []uint:
		return len(v)
	case *[]uint:
		return len(*v)
	case []uint8: // also `byte`
		return len(v)
	case *[]uint8:
		return len(*v)
	case []uint16:
		return len(v)
	case *[]uint16:
		return len(*v)
	case []uint32: // also `rune`
		return len(v)
	case *[]uint32:
		return len(*v)
	case []uint64:
		return len(v)
	case *[]uint64:
		return len(*v)
	case []int:
		return len(v)
	case *[]int:
		return len(*v)
	case []int8:
		return len(v)
	case *[]int8:
		return len(*v)
	case []int16:
		return len(v)
	case *[]int16:
		return len(*v)
	case []int32: // also `rune`
		return len(v)
	case *[]int32:
		return len(*v)
	case []int64:
		return len(v)
	case *[]int64:
		return len(*v)
	case []float32:
		return len(v)
	case *[]float32:
		return len(*v)
	case []float64:
		return len(v)
	case *[]float64:
		return len(*v)
	case []string:
		return len(v)
	case *[]string:
		return len(*v)
	}
	return -1
}

// Returns the element at index `i` of the underlying slice or nil if
// not applicable
func ItemAt(ix interface{}, i int) interface{} {
	if ix == nil {
		return nil
	}
	if i < 0 {
		return nil
	}
	// It seems that there's no way of saying:
	//
	//     if v, ok := ix.([]interface{}); ok { ... }
	//     case []interface{}:
	//     if v, ok := ix.([]any); ok { ... }
	//     case []any:
	switch v := ix.(type) {
	case []bool:
		if i < len(v) {
			return v[i]
		}
	case *[]bool:
		if i < len(*v) {
			return (*v)[i]
		}

	case []uint:
		if i < len(v) {
			return v[i]
		}
	case *[]uint:
		if i < len(*v) {
			return (*v)[i]
		}

	case []uint8: // also `byte`
		if i < len(v) {
			return v[i]
		}
	case *[]uint8:
		if i < len(*v) {
			return (*v)[i]
		}

	case []uint16:
		if i < len(v) {
			return v[i]
		}
	case *[]uint16:
		if i < len(*v) {
			return (*v)[i]
		}

	case []uint32: // also `rune`
		if i < len(v) {
			return v[i]
		}
	case *[]uint32:
		if i < len(*v) {
			return (*v)[i]
		}

	case []uint64:
		if i < len(v) {
			return v[i]
		}
	case *[]uint64:
		if i < len(*v) {
			return (*v)[i]
		}

	case []int:
		if i < len(v) {
			return v[i]
		}
	case *[]int:
		if i < len(*v) {
			return (*v)[i]
		}

	case []int8:
		if i < len(v) {
			return v[i]
		}
	case *[]int8:
		if i < len(*v) {
			return (*v)[i]
		}

	case []int16:
		if i < len(v) {
			return v[i]
		}
	case *[]int16:
		if i < len(*v) {
			return (*v)[i]
		}

	case []int32: // also `rune`
		if i < len(v) {
			return v[i]
		}
	case *[]int32:
		if i < len(*v) {
			return (*v)[i]
		}

	case []int64:
		if i < len(v) {
			return v[i]
		}
	case *[]int64:
		if i < len(*v) {
			return (*v)[i]
		}

	case []float32:
		if i < len(v) {
			return v[i]
		}
	case *[]float32:
		if i < len(*v) {
			return (*v)[i]
		}

	case []float64:
		if i < len(v) {
			return v[i]
		}
	case *[]float64:
		if i < len(*v) {
			return (*v)[i]
		}

	case []string:
		if i < len(v) {
			return v[i]
		}
	case *[]string:
		if i < len(*v) {
			return (*v)[i]
		}
	}
	return nil
}

// Returns true if ixa is a pointer to the same type as ixb
func IsPointerTo(ixa, ixb interface{}) bool {
	ta := Type(ixa)
	tb := Type(ixb)
	if !ta.TstPointerBit() || tb.TstPointerBit() {
		// ixa is not a pointer or ixb is a pointer
		return false
	}
	if tb|PointerT == ta {
		return true
	}
	return false
}

// See the various strconv.Format<Type> functions
type StrConvParams struct {
	base int    // for strconv.FormatInt() and .FormatUint()
	fmt  byte   // for strconv.FormatFloat()
	prec int    // for strconv.FormatFloat()
	sep  string // separator for slice elements in returned string
}

const baseDefault int = 10
const fmtDefault byte = byte('g')
const precDefault int = -1
const sepDefault string = ", "

type StrConvOption = func(f *StrConvParams)

func WithBase(base int) StrConvOption {
	return func(p *StrConvParams) {
		p.base = base
	}
}
func WithFmt(fmt byte) StrConvOption {
	return func(p *StrConvParams) {
		p.fmt = fmt
	}
}
func WithPrec(prec int) StrConvOption {
	return func(p *StrConvParams) {
		p.prec = prec
	}
}
func WithSep(sep string) StrConvOption {
	return func(p *StrConvParams) {
		p.sep = sep
	}
}

func StrConv(ix interface{}, opts ...StrConvOption) string {
	param := &StrConvParams{
		base: baseDefault,
		fmt:  fmtDefault,
		prec: precDefault,
		sep:  sepDefault,
	}
	for _, opt := range opts {
		opt(param)
	}

	// Return the empty string for any zero-length slice or slice pointer:
	if SliceLen(ix) == 0 {
		return ""
	}

	buf := bytes.Buffer{}
	switch v := ix.(type) {
	// Boolean
	case bool:
		return strconv.FormatBool(v)
	case *bool:
		return strconv.FormatBool(*v)
	case []bool:
		buf.WriteString(strconv.FormatBool(v[0]))
		for _, b := range v[1:] {
			buf.WriteString(param.sep + strconv.FormatBool(b))
		}
	case *[]bool:
		buf.WriteString(strconv.FormatBool((*v)[0]))
		for _, b := range (*v)[1:] {
			buf.WriteString(param.sep + strconv.FormatBool(b))
		}

	// Unsigned integers
	case uint:
		return strconv.FormatUint(uint64(v), param.base)
	case *uint:
		return strconv.FormatUint(uint64(*v), param.base)
	case []uint:
		buf.WriteString(strconv.FormatUint(uint64(v[0]), param.base))
		for _, u := range v[1:] {
			buf.WriteString(param.sep + strconv.FormatUint(uint64(u), param.base))
		}
	case *[]uint:
		buf.WriteString(strconv.FormatUint(uint64((*v)[0]), param.base))
		for _, u := range (*v)[1:] {
			buf.WriteString(param.sep + strconv.FormatUint(uint64(u), param.base))
		}

	case uint8: // also `byte`
		return strconv.FormatUint(uint64(v), param.base)
	case *uint8:
		return strconv.FormatUint(uint64(*v), param.base)
	case []uint8:
		buf.WriteString(strconv.FormatUint(uint64(v[0]), param.base))
		for _, u := range v[1:] {
			buf.WriteString(param.sep + strconv.FormatUint(uint64(u), param.base))
		}
	case *[]uint8:
		buf.WriteString(strconv.FormatUint(uint64((*v)[0]), param.base))
		for _, u := range (*v)[1:] {
			buf.WriteString(param.sep + strconv.FormatUint(uint64(u), param.base))
		}

	case uint16:
		return strconv.FormatUint(uint64(v), param.base)
	case *uint16:
		return strconv.FormatUint(uint64(*v), param.base)
	case []uint16:
		buf.WriteString(strconv.FormatUint(uint64(v[0]), param.base))
		for _, u := range v[1:] {
			buf.WriteString(param.sep + strconv.FormatUint(uint64(u), param.base))
		}
	case *[]uint16:
		buf.WriteString(strconv.FormatUint(uint64((*v)[0]), param.base))
		for _, u := range (*v)[1:] {
			buf.WriteString(param.sep + strconv.FormatUint(uint64(u), param.base))
		}

	case uint32:
		return strconv.FormatUint(uint64(v), param.base)
	case *uint32:
		return strconv.FormatUint(uint64(*v), param.base)
	case []uint32:
		buf.WriteString(strconv.FormatUint(uint64(v[0]), param.base))
		for _, u := range v[1:] {
			buf.WriteString(param.sep + strconv.FormatUint(uint64(u), param.base))
		}
	case *[]uint32:
		buf.WriteString(strconv.FormatUint(uint64((*v)[0]), param.base))
		for _, u := range (*v)[1:] {
			buf.WriteString(param.sep + strconv.FormatUint(uint64(u), param.base))
		}

	case uint64:
		return strconv.FormatUint(uint64(v), param.base)
	case *uint64:
		return strconv.FormatUint(uint64(*v), param.base)
	case []uint64:
		buf.WriteString(strconv.FormatUint(uint64(v[0]), param.base))
		for _, u := range v[1:] {
			buf.WriteString(param.sep + strconv.FormatUint(uint64(u), param.base))
		}
	case *[]uint64:
		buf.WriteString(strconv.FormatUint(uint64((*v)[0]), param.base))
		for _, u := range (*v)[1:] {
			buf.WriteString(param.sep + strconv.FormatUint(uint64(u), param.base))
		}

	// Signed integers
	case int:
		return strconv.FormatInt(int64(v), param.base)
	case *int:
		return strconv.FormatInt(int64(*v), param.base)
	case []int:
		buf.WriteString(strconv.FormatInt(int64(v[0]), param.base))
		for _, u := range v[1:] {
			buf.WriteString(param.sep + strconv.FormatInt(int64(u), param.base))
		}
	case *[]int:
		buf.WriteString(strconv.FormatInt(int64((*v)[0]), param.base))
		for _, u := range (*v)[1:] {
			buf.WriteString(param.sep + strconv.FormatInt(int64(u), param.base))
		}

	case int8:
		return strconv.FormatInt(int64(v), param.base)
	case *int8:
		return strconv.FormatInt(int64(*v), param.base)
	case []int8:
		buf.WriteString(strconv.FormatInt(int64(v[0]), param.base))
		for _, u := range v[1:] {
			buf.WriteString(param.sep + strconv.FormatInt(int64(u), param.base))
		}
	case *[]int8:
		buf.WriteString(strconv.FormatInt(int64((*v)[0]), param.base))
		for _, u := range (*v)[1:] {
			buf.WriteString(param.sep + strconv.FormatInt(int64(u), param.base))
		}

	case int16:
		return strconv.FormatInt(int64(v), param.base)
	case *int16:
		return strconv.FormatInt(int64(*v), param.base)
	case []int16:
		buf.WriteString(strconv.FormatInt(int64(v[0]), param.base))
		for _, u := range v[1:] {
			buf.WriteString(param.sep + strconv.FormatInt(int64(u), param.base))
		}
	case *[]int16:
		buf.WriteString(strconv.FormatInt(int64((*v)[0]), param.base))
		for _, u := range (*v)[1:] {
			buf.WriteString(param.sep + strconv.FormatInt(int64(u), param.base))
		}

	case int32: // also `rune`
		return strconv.FormatInt(int64(v), param.base)
	case *int32:
		return strconv.FormatInt(int64(*v), param.base)
	case []int32:
		buf.WriteString(strconv.FormatInt(int64(v[0]), param.base))
		for _, u := range v[1:] {
			buf.WriteString(param.sep + strconv.FormatInt(int64(u), param.base))
		}
	case *[]int32:
		buf.WriteString(strconv.FormatInt(int64((*v)[0]), param.base))
		for _, u := range (*v)[1:] {
			buf.WriteString(param.sep + strconv.FormatInt(int64(u), param.base))
		}

	case int64:
		return strconv.FormatInt(int64(v), param.base)
	case *int64:
		return strconv.FormatInt(int64(*v), param.base)
	case []int64:
		buf.WriteString(strconv.FormatInt(int64(v[0]), param.base))
		for _, u := range v[1:] {
			buf.WriteString(param.sep + strconv.FormatInt(int64(u), param.base))
		}
	case *[]int64:
		buf.WriteString(strconv.FormatInt(int64((*v)[0]), param.base))
		for _, u := range (*v)[1:] {
			buf.WriteString(param.sep + strconv.FormatInt(int64(u), param.base))
		}

	// Floating-point types
	case float32:
		return strconv.FormatFloat(float64(v), param.fmt, param.prec, 64)
	case *float32:
		return strconv.FormatFloat(float64(*v), param.fmt, param.prec, 64)
	case []float32:
		buf.WriteString(strconv.FormatFloat(float64(v[0]), param.fmt, param.prec, 64))
		for _, u := range v[1:] {
			buf.WriteString(param.sep + strconv.FormatFloat(float64(u), param.fmt, param.prec, 64))
		}
	case *[]float32:
		buf.WriteString(strconv.FormatFloat(float64((*v)[0]), param.fmt, param.prec, 64))
		for _, u := range (*v)[1:] {
			buf.WriteString(param.sep + strconv.FormatFloat(float64(u), param.fmt, param.prec, 64))
		}

	case float64:
		return strconv.FormatFloat(float64(v), param.fmt, param.prec, 64)
	case *float64:
		return strconv.FormatFloat(float64(*v), param.fmt, param.prec, 64)
	case []float64:
		buf.WriteString(strconv.FormatFloat(float64(v[0]), param.fmt, param.prec, 64))
		for _, u := range v[1:] {
			buf.WriteString(param.sep + strconv.FormatFloat(float64(u), param.fmt, param.prec, 64))
		}
	case *[]float64:
		buf.WriteString(strconv.FormatFloat(float64((*v)[0]), param.fmt, param.prec, 64))
		for _, u := range (*v)[1:] {
			buf.WriteString(param.sep + strconv.FormatFloat(float64(u), param.fmt, param.prec, 64))
		}

	// A little bit silly, but for completeness:
	case string:
		return v
	case *string:
		return *v
	case []string:
		return strings.Join(v, param.sep)
	case *[]string:
		return strings.Join(*v, param.sep)
	}

	return buf.String()
}

func FromStr(ix interface{}, str string, opts ...StrConvOption) error {
	// Prefer the SetValue interface
	if settee, ok := ix.(SetValue); ok {
		return settee.Set(str)
	}

	param := &StrConvParams{
		base: baseDefault,
		sep:  ",",
		// `fmt` and `prec` are ignored
	}
	for _, opt := range opts {
		opt(param)
	}

	typeId := Type(ix)
	if typeId.TstOtherBit() {
		return fmt.Errorf("interface (%v) does not represent a supported type (%T)", ix, ix)
	}
	if !typeId.TstPointerBit() {
		return fmt.Errorf("interface (%v) does not represent a pointer (%T)", ix, ix)
	}

	switch v := ix.(type) {
	// Booleans
	case *bool:
		b, err := strconv.ParseBool(str)
		if err != nil {
			return err
		}
		*v = b
		return nil
	case *[]bool:
		for _, item := range strings.Split(str, param.sep) {
			trimmed := strings.TrimSpace(item)
			b, err := strconv.ParseBool(trimmed)
			if err != nil {
				return err
			}
			*v = append(*v, b)
		}
		return nil

	// Unsigned integers
	case *uint:
		u64, err := strconv.ParseUint(str, param.base, strconv.IntSize)
		if err != nil {
			return err
		}
		*v = uint(u64)
		return nil
	case *[]uint:
		for _, item := range strings.Split(str, param.sep) {
			trimmed := strings.TrimSpace(item)
			u64, err := strconv.ParseUint(trimmed, param.base, strconv.IntSize)
			if err != nil {
				return err
			}
			*v = append(*v, uint(u64))
		}
		return nil

	case *uint8: // also `byte`:
		u64, err := strconv.ParseUint(str, param.base, 8)
		if err != nil {
			return err
		}
		*v = uint8(u64)
		return nil
	case *[]uint8:
		for _, item := range strings.Split(str, param.sep) {
			trimmed := strings.TrimSpace(item)
			u64, err := strconv.ParseUint(trimmed, param.base, 8)
			if err != nil {
				return err
			}
			*v = append(*v, uint8(u64))
		}
		return nil

	case *uint16:
		u64, err := strconv.ParseUint(str, param.base, 16)
		if err != nil {
			return err
		}
		*v = uint16(u64)
		return nil
	case *[]uint16:
		for _, item := range strings.Split(str, param.sep) {
			trimmed := strings.TrimSpace(item)
			u64, err := strconv.ParseUint(trimmed, param.base, 16)
			if err != nil {
				return err
			}
			*v = append(*v, uint16(u64))
		}
		return nil

	case *uint32:
		u64, err := strconv.ParseUint(str, param.base, 32)
		if err != nil {
			return err
		}
		*v = uint32(u64)
		return nil
	case *[]uint32:
		for _, item := range strings.Split(str, param.sep) {
			trimmed := strings.TrimSpace(item)
			u64, err := strconv.ParseUint(trimmed, param.base, 32)
			if err != nil {
				return err
			}
			*v = append(*v, uint32(u64))
		}
		return nil

	case *uint64:
		u64, err := strconv.ParseUint(str, param.base, 64)
		if err != nil {
			return err
		}
		*v = u64
		return nil
	case *[]uint64:
		for _, item := range strings.Split(str, param.sep) {
			trimmed := strings.TrimSpace(item)
			u64, err := strconv.ParseUint(trimmed, param.base, 64)
			if err != nil {
				return err
			}
			*v = append(*v, u64)
		}
		return nil

	// Signed integers
	case *int:
		i64, err := strconv.ParseInt(str, param.base, strconv.IntSize)
		if err != nil {
			return err
		}
		*v = int(i64)
		return nil
	case *[]int:
		for _, item := range strings.Split(str, param.sep) {
			trimmed := strings.TrimSpace(item)
			i64, err := strconv.ParseInt(trimmed, param.base, strconv.IntSize)
			if err != nil {
				return err
			}
			*v = append(*v, int(i64))
		}
		return nil

	case *int8:
		i64, err := strconv.ParseInt(str, param.base, 8)
		if err != nil {
			return err
		}
		*v = int8(i64)
		return nil
	case *[]int8:
		for _, item := range strings.Split(str, param.sep) {
			trimmed := strings.TrimSpace(item)
			i64, err := strconv.ParseInt(trimmed, param.base, 8)
			if err != nil {
				return err
			}
			*v = append(*v, int8(i64))
		}
		return nil

	case *int16:
		i64, err := strconv.ParseInt(str, param.base, 16)
		if err != nil {
			return err
		}
		*v = int16(i64)
		return nil
	case *[]int16:
		for _, item := range strings.Split(str, param.sep) {
			trimmed := strings.TrimSpace(item)
			i64, err := strconv.ParseInt(trimmed, param.base, 16)
			if err != nil {
				return err
			}
			*v = append(*v, int16(i64))
		}
		return nil

	case *int32: // also `rune`
		i64, err := strconv.ParseInt(str, param.base, 32)
		if err != nil {
			return err
		}
		*v = int32(i64)
		return nil
	case *[]int32:
		for _, item := range strings.Split(str, param.sep) {
			trimmed := strings.TrimSpace(item)
			i64, err := strconv.ParseInt(trimmed, param.base, 32)
			if err != nil {
				return err
			}
			*v = append(*v, int32(i64))
		}
		return nil

	case *int64:
		i64, err := strconv.ParseInt(str, param.base, 64)
		if err != nil {
			return err
		}
		*v = i64
		return nil
	case *[]int64:
		for _, item := range strings.Split(str, param.sep) {
			trimmed := strings.TrimSpace(item)
			i64, err := strconv.ParseInt(trimmed, param.base, 64)
			if err != nil {
				return err
			}
			*v = append(*v, i64)
		}
		return nil

	// Floating-point number
	case *float32:
		f64, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return err
		}
		*v = float32(f64)
		return nil
	case *[]float32:
		for _, item := range strings.Split(str, param.sep) {
			trimmed := strings.TrimSpace(item)
			f64, err := strconv.ParseFloat(trimmed, 32)
			if err != nil {
				return err
			}
			*v = append(*v, float32(f64))
		}
		return nil

	case *float64:
		f64, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		*v = f64
		return nil
	case *[]float64:
		for _, item := range strings.Split(str, param.sep) {
			trimmed := strings.TrimSpace(item)
			f64, err := strconv.ParseFloat(trimmed, 64)
			if err != nil {
				return err
			}
			*v = append(*v, f64)
		}
		return nil

	// Strings - a little silly, but for completeness
	case *string:
		*v = str
		return nil
	case *[]string:
		for _, item := range strings.Split(str, param.sep) {
			// trimmed := strings.TrimSpace(item)
			*v = append(*v, item)
		}
		return nil

	}
	return nil
}
