package types

import(
    "fmt"
    "strconv"

    "golang.org/x/exp/constraints"
)

func pow2(n int) uint64 {
	return uint64(1) << n
}

func pow2m1(n int) uint64 {
	u := uint64(1)
	for i := 0; i < n-1; i++ {
		u = (u << 1) | uint64(1)
	}
	return u
}

// Returns the max and min values for a type
func (tp *TypeId) MinAndMax() (int64, uint64) {
	if tp.TstBoolBit() {
		return 0, 1
	}
	n := tp.BitSize()
	if tp.TstIntBit() {
		return int64(-pow2(n - 1)), pow2m1(n - 1)
	}
	if tp.TstUintBit() {
		return 0, pow2m1(n)
	}
	if tp.TstFloatBit() {
		if n == 32 {
			return int64(-pow2(24)), pow2(24)
		}
		return int64(-pow2(53)), pow2(53)
	}
	return 0, 0
}

type Number interface {
	constraints.Integer | constraints.Float
}

// RangeTest returns nil if the value of the test type is
// representable by the value of the reference type:
func RangeTest[T Number, R Number](t T, r R) error {
	tt := Type(t)
	tr := Type(r)
	// If the types are the same, there's no issue:
	if tt == tr {
		return nil
	}
	// If the reference type has more bits than the test type, but
	// they're otherwise the same, there's no issue:
	if (tt & ^NumBits) == (tr & ^NumBits) {
		if tt.BitSize() <= tr.BitSize() {
			return nil
		}
		// The reference type has fewer bits than the test type, so we
		// need to check the value
	}
	rmin, rmax := tr.MinAndMax()
	if uint64(t) <= rmax && int64(t) >= rmin {
		return nil
	}
	return fmt.Errorf("value %v<%T> is not representable in %T", t, t, r)
}

func CoerceScalar(ref interface{}, val interface{}) (interface{}, error) {
    if ref == nil || val == nil {
        return nil, fmt.Errorf("nil argument given")
    }
    switch ref.(type) {
    case bool:
        switch v := val.(type) {
        case bool:
            return v, nil
        case int:
            return v != int(0), nil
        case int8:
            return v != int8(0), nil
        case int16:
            return v != int16(0), nil
        case int32:
            return v != int32(0), nil
        case int64:
            return v != int64(0), nil
        case uint:
            return v != uint(0), nil
        case uint8:
            return v != uint8(0), nil
        case uint16:
            return v != uint16(0), nil
        case uint32:
            return v != uint32(0), nil
        case uint64:
            return v != uint64(0), nil
        case float32:
            return v != float32(0), nil
        case float64:
            return v != float64(0), nil
        case string:
            n, err := strconv.ParseBool(v)
            return bool(n), err
        }
    case int:
        switch v := val.(type) {
        case bool:
            if v {
                return int(1), nil
            }
            return int(0), nil
        case int:
            return v, nil
        case int8:
            return int(v), nil
            // a 8-bit int is always representable in a 32-bit int
        case int16:
            return int(v), nil
            // a 16-bit int is always representable in a 32-bit int
        case int32:
            return int(v), nil
            // a 32-bit int is always representable in a 32-bit int
        case int64:
            return int(v), RangeTest(v, int(0))
            // Value test needed: int has 32 bits, int64 has 64 bits
        case uint:
            return int(v), RangeTest(v, int(0))
            // Value test needed: int has 32 bits, uint has 32 bits
        case uint8:
            return int(v), nil
            // uint8 is always representable in int32
        case uint16:
            return int(v), nil
            // uint16 is always representable in int32
        case uint32:
            return int(v), RangeTest(v, int(0))
            // Value test needed: int has 32 bits, uint32 has 32 bits
        case uint64:
            return int(v), RangeTest(v, int(0))
            // Value test needed: int has 32 bits, uint64 has 64 bits
        case float32:
            return int(v), RangeTest(v, int(0))
            // Value test needed: int has 32 bits, float32 has 24 bits
        case float64:
            return int(v), RangeTest(v, int(0))
            // Value test needed: int has 32 bits, float64 has 53 bits
        case string:
            n, err := strconv.ParseInt(v, 10, strconv.IntSize)
            return int(n), err
        }
    case int8:
        switch v := val.(type) {
        case bool:
            if v {
                return int8(1), nil
            }
            return int8(0), nil
        case int:
            return int8(v), RangeTest(v, int8(0))
            // Value test needed: int8 has 8 bits, int has 32 bits
        case int8:
            return v, nil
        case int16:
            return int8(v), RangeTest(v, int8(0))
            // Value test needed: int8 has 8 bits, int16 has 16 bits
        case int32:
            return int8(v), RangeTest(v, int8(0))
            // Value test needed: int8 has 8 bits, int32 has 32 bits
        case int64:
            return int8(v), RangeTest(v, int8(0))
            // Value test needed: int8 has 8 bits, int64 has 64 bits
        case uint:
            return int8(v), RangeTest(v, int8(0))
            // Value test needed: int8 has 8 bits, uint has 32 bits
        case uint8:
            return int8(v), RangeTest(v, int8(0))
            // Value test needed: int8 has 8 bits, uint8 has 8 bits
        case uint16:
            return int8(v), RangeTest(v, int8(0))
            // Value test needed: int8 has 8 bits, uint16 has 16 bits
        case uint32:
            return int8(v), RangeTest(v, int8(0))
            // Value test needed: int8 has 8 bits, uint32 has 32 bits
        case uint64:
            return int8(v), RangeTest(v, int8(0))
            // Value test needed: int8 has 8 bits, uint64 has 64 bits
        case float32:
            return int8(v), RangeTest(v, int8(0))
            // Value test needed: int8 has 8 bits, float32 has 24 bits
        case float64:
            return int8(v), RangeTest(v, int8(0))
            // Value test needed: int8 has 8 bits, float64 has 53 bits
        case string:
            n, err := strconv.ParseInt(v, 10, 8)
            return int8(n), err
        }
    case int16:
        switch v := val.(type) {
        case bool:
            if v {
                return int16(1), nil
            }
            return int16(0), nil
        case int:
            return int16(v), RangeTest(v, int16(0))
            // Value test needed: int16 has 16 bits, int has 32 bits
        case int8:
            return int16(v), nil
            // a 8-bit int is always representable in a 16-bit int
        case int16:
            return v, nil
        case int32:
            return int16(v), RangeTest(v, int16(0))
            // Value test needed: int16 has 16 bits, int32 has 32 bits
        case int64:
            return int16(v), RangeTest(v, int16(0))
            // Value test needed: int16 has 16 bits, int64 has 64 bits
        case uint:
            return int16(v), RangeTest(v, int16(0))
            // Value test needed: int16 has 16 bits, uint has 32 bits
        case uint8:
            return int16(v), nil
            // uint8 is always representable in int16
        case uint16:
            return int16(v), RangeTest(v, int16(0))
            // Value test needed: int16 has 16 bits, uint16 has 16 bits
        case uint32:
            return int16(v), RangeTest(v, int16(0))
            // Value test needed: int16 has 16 bits, uint32 has 32 bits
        case uint64:
            return int16(v), RangeTest(v, int16(0))
            // Value test needed: int16 has 16 bits, uint64 has 64 bits
        case float32:
            return int16(v), RangeTest(v, int16(0))
            // Value test needed: int16 has 16 bits, float32 has 24 bits
        case float64:
            return int16(v), RangeTest(v, int16(0))
            // Value test needed: int16 has 16 bits, float64 has 53 bits
        case string:
            n, err := strconv.ParseInt(v, 10, 16)
            return int16(n), err
        }
    case int32:
        switch v := val.(type) {
        case bool:
            if v {
                return int32(1), nil
            }
            return int32(0), nil
        case int:
            return int32(v), nil
            // a 32-bit int is always representable in a 32-bit int
        case int8:
            return int32(v), nil
            // a 8-bit int is always representable in a 32-bit int
        case int16:
            return int32(v), nil
            // a 16-bit int is always representable in a 32-bit int
        case int32:
            return v, nil
        case int64:
            return int32(v), RangeTest(v, int32(0))
            // Value test needed: int32 has 32 bits, int64 has 64 bits
        case uint:
            return int32(v), RangeTest(v, int32(0))
            // Value test needed: int32 has 32 bits, uint has 32 bits
        case uint8:
            return int32(v), nil
            // uint8 is always representable in int32
        case uint16:
            return int32(v), nil
            // uint16 is always representable in int32
        case uint32:
            return int32(v), RangeTest(v, int32(0))
            // Value test needed: int32 has 32 bits, uint32 has 32 bits
        case uint64:
            return int32(v), RangeTest(v, int32(0))
            // Value test needed: int32 has 32 bits, uint64 has 64 bits
        case float32:
            return int32(v), RangeTest(v, int32(0))
            // Value test needed: int32 has 32 bits, float32 has 24 bits
        case float64:
            return int32(v), RangeTest(v, int32(0))
            // Value test needed: int32 has 32 bits, float64 has 53 bits
        case string:
            n, err := strconv.ParseInt(v, 10, 32)
            return int32(n), err
        }
    case int64:
        switch v := val.(type) {
        case bool:
            if v {
                return int64(1), nil
            }
            return int64(0), nil
        case int:
            return int64(v), nil
            // a 32-bit int is always representable in a 64-bit int
        case int8:
            return int64(v), nil
            // a 8-bit int is always representable in a 64-bit int
        case int16:
            return int64(v), nil
            // a 16-bit int is always representable in a 64-bit int
        case int32:
            return int64(v), nil
            // a 32-bit int is always representable in a 64-bit int
        case int64:
            return v, nil
        case uint:
            return int64(v), nil
            // uint32 is always representable in int64
        case uint8:
            return int64(v), nil
            // uint8 is always representable in int64
        case uint16:
            return int64(v), nil
            // uint16 is always representable in int64
        case uint32:
            return int64(v), nil
            // uint32 is always representable in int64
        case uint64:
            return int64(v), RangeTest(v, int64(0))
            // Value test needed: int64 has 64 bits, uint64 has 64 bits
        case float32:
            return int64(v), RangeTest(v, int64(0))
            // Value test needed: int64 has 64 bits, float32 has 24 bits
        case float64:
            return int64(v), RangeTest(v, int64(0))
            // Value test needed: int64 has 64 bits, float64 has 53 bits
        case string:
            n, err := strconv.ParseInt(v, 10, 64)
            return int64(n), err
        }
    case uint:
        switch v := val.(type) {
        case bool:
            if v {
                return uint(1), nil
            }
            return uint(0), nil
        case int:
            return uint(v), RangeTest(v, uint(0))
            // Value test needed: uint has 32 bits, int has 32 bits
        case int8:
            return uint(v), RangeTest(v, uint(0))
            // Value test needed: uint has 32 bits, int8 has 8 bits
        case int16:
            return uint(v), RangeTest(v, uint(0))
            // Value test needed: uint has 32 bits, int16 has 16 bits
        case int32:
            return uint(v), RangeTest(v, uint(0))
            // Value test needed: uint has 32 bits, int32 has 32 bits
        case int64:
            return uint(v), RangeTest(v, uint(0))
            // Value test needed: uint has 32 bits, int64 has 64 bits
        case uint:
            return v, nil
        case uint8:
            return uint(v), nil
            // a 8-bit uint is always representable in a 32-bit uint
        case uint16:
            return uint(v), nil
            // a 16-bit uint is always representable in a 32-bit uint
        case uint32:
            return uint(v), nil
            // a 32-bit uint is always representable in a 32-bit uint
        case uint64:
            return uint(v), RangeTest(v, uint(0))
            // Value test needed: uint has 32 bits, uint64 has 64 bits
        case float32:
            return uint(v), RangeTest(v, uint(0))
            // Value test needed: uint has 32 bits, float32 has 24 bits
        case float64:
            return uint(v), RangeTest(v, uint(0))
            // Value test needed: uint has 32 bits, float64 has 53 bits
        case string:
            n, err := strconv.ParseUint(v, 10, strconv.IntSize)
            return uint(n), err
        }
    case uint8:
        switch v := val.(type) {
        case bool:
            if v {
                return uint8(1), nil
            }
            return uint8(0), nil
        case int:
            return uint8(v), RangeTest(v, uint8(0))
            // Value test needed: uint8 has 8 bits, int has 32 bits
        case int8:
            return uint8(v), RangeTest(v, uint8(0))
            // Value test needed: uint8 has 8 bits, int8 has 8 bits
        case int16:
            return uint8(v), RangeTest(v, uint8(0))
            // Value test needed: uint8 has 8 bits, int16 has 16 bits
        case int32:
            return uint8(v), RangeTest(v, uint8(0))
            // Value test needed: uint8 has 8 bits, int32 has 32 bits
        case int64:
            return uint8(v), RangeTest(v, uint8(0))
            // Value test needed: uint8 has 8 bits, int64 has 64 bits
        case uint:
            return uint8(v), RangeTest(v, uint8(0))
            // Value test needed: uint8 has 8 bits, uint has 32 bits
        case uint8:
            return v, nil
        case uint16:
            return uint8(v), RangeTest(v, uint8(0))
            // Value test needed: uint8 has 8 bits, uint16 has 16 bits
        case uint32:
            return uint8(v), RangeTest(v, uint8(0))
            // Value test needed: uint8 has 8 bits, uint32 has 32 bits
        case uint64:
            return uint8(v), RangeTest(v, uint8(0))
            // Value test needed: uint8 has 8 bits, uint64 has 64 bits
        case float32:
            return uint8(v), RangeTest(v, uint8(0))
            // Value test needed: uint8 has 8 bits, float32 has 24 bits
        case float64:
            return uint8(v), RangeTest(v, uint8(0))
            // Value test needed: uint8 has 8 bits, float64 has 53 bits
        case string:
            n, err := strconv.ParseUint(v, 10, 8)
            return uint8(n), err
        }
    case uint16:
        switch v := val.(type) {
        case bool:
            if v {
                return uint16(1), nil
            }
            return uint16(0), nil
        case int:
            return uint16(v), RangeTest(v, uint16(0))
            // Value test needed: uint16 has 16 bits, int has 32 bits
        case int8:
            return uint16(v), RangeTest(v, uint16(0))
            // Value test needed: uint16 has 16 bits, int8 has 8 bits
        case int16:
            return uint16(v), RangeTest(v, uint16(0))
            // Value test needed: uint16 has 16 bits, int16 has 16 bits
        case int32:
            return uint16(v), RangeTest(v, uint16(0))
            // Value test needed: uint16 has 16 bits, int32 has 32 bits
        case int64:
            return uint16(v), RangeTest(v, uint16(0))
            // Value test needed: uint16 has 16 bits, int64 has 64 bits
        case uint:
            return uint16(v), RangeTest(v, uint16(0))
            // Value test needed: uint16 has 16 bits, uint has 32 bits
        case uint8:
            return uint16(v), nil
            // a 8-bit uint is always representable in a 16-bit uint
        case uint16:
            return v, nil
        case uint32:
            return uint16(v), RangeTest(v, uint16(0))
            // Value test needed: uint16 has 16 bits, uint32 has 32 bits
        case uint64:
            return uint16(v), RangeTest(v, uint16(0))
            // Value test needed: uint16 has 16 bits, uint64 has 64 bits
        case float32:
            return uint16(v), RangeTest(v, uint16(0))
            // Value test needed: uint16 has 16 bits, float32 has 24 bits
        case float64:
            return uint16(v), RangeTest(v, uint16(0))
            // Value test needed: uint16 has 16 bits, float64 has 53 bits
        case string:
            n, err := strconv.ParseUint(v, 10, 16)
            return uint16(n), err
        }
    case uint32:
        switch v := val.(type) {
        case bool:
            if v {
                return uint32(1), nil
            }
            return uint32(0), nil
        case int:
            return uint32(v), RangeTest(v, uint32(0))
            // Value test needed: uint32 has 32 bits, int has 32 bits
        case int8:
            return uint32(v), RangeTest(v, uint32(0))
            // Value test needed: uint32 has 32 bits, int8 has 8 bits
        case int16:
            return uint32(v), RangeTest(v, uint32(0))
            // Value test needed: uint32 has 32 bits, int16 has 16 bits
        case int32:
            return uint32(v), RangeTest(v, uint32(0))
            // Value test needed: uint32 has 32 bits, int32 has 32 bits
        case int64:
            return uint32(v), RangeTest(v, uint32(0))
            // Value test needed: uint32 has 32 bits, int64 has 64 bits
        case uint:
            return uint32(v), nil
            // a 32-bit uint is always representable in a 32-bit uint
        case uint8:
            return uint32(v), nil
            // a 8-bit uint is always representable in a 32-bit uint
        case uint16:
            return uint32(v), nil
            // a 16-bit uint is always representable in a 32-bit uint
        case uint32:
            return v, nil
        case uint64:
            return uint32(v), RangeTest(v, uint32(0))
            // Value test needed: uint32 has 32 bits, uint64 has 64 bits
        case float32:
            return uint32(v), RangeTest(v, uint32(0))
            // Value test needed: uint32 has 32 bits, float32 has 24 bits
        case float64:
            return uint32(v), RangeTest(v, uint32(0))
            // Value test needed: uint32 has 32 bits, float64 has 53 bits
        case string:
            n, err := strconv.ParseUint(v, 10, 32)
            return uint32(n), err
        }
    case uint64:
        switch v := val.(type) {
        case bool:
            if v {
                return uint64(1), nil
            }
            return uint64(0), nil
        case int:
            return uint64(v), RangeTest(v, uint64(0))
            // Value test needed: uint64 has 64 bits, int has 32 bits
        case int8:
            return uint64(v), RangeTest(v, uint64(0))
            // Value test needed: uint64 has 64 bits, int8 has 8 bits
        case int16:
            return uint64(v), RangeTest(v, uint64(0))
            // Value test needed: uint64 has 64 bits, int16 has 16 bits
        case int32:
            return uint64(v), RangeTest(v, uint64(0))
            // Value test needed: uint64 has 64 bits, int32 has 32 bits
        case int64:
            return uint64(v), RangeTest(v, uint64(0))
            // Value test needed: uint64 has 64 bits, int64 has 64 bits
        case uint:
            return uint64(v), nil
            // a 32-bit uint is always representable in a 64-bit uint
        case uint8:
            return uint64(v), nil
            // a 8-bit uint is always representable in a 64-bit uint
        case uint16:
            return uint64(v), nil
            // a 16-bit uint is always representable in a 64-bit uint
        case uint32:
            return uint64(v), nil
            // a 32-bit uint is always representable in a 64-bit uint
        case uint64:
            return v, nil
        case float32:
            return uint64(v), RangeTest(v, uint64(0))
            // Value test needed: uint64 has 64 bits, float32 has 24 bits
        case float64:
            return uint64(v), RangeTest(v, uint64(0))
            // Value test needed: uint64 has 64 bits, float64 has 53 bits
        case string:
            n, err := strconv.ParseUint(v, 10, 64)
            return uint64(n), err
        }
    case float32:
        switch v := val.(type) {
        case bool:
            if v {
                return float32(1), nil
            }
            return float32(0), nil
        case int:
            return float32(v), RangeTest(v, float32(0))
            // Value test needed: float32 has 24 bits, int has 32 bits
        case int8:
            return float32(v), nil
            // a float32 can exactly represent any int8
        case int16:
            return float32(v), nil
            // a float32 can exactly represent any int16
        case int32:
            return float32(v), RangeTest(v, float32(0))
            // Value test needed: float32 has 24 bits, int32 has 32 bits
        case int64:
            return float32(v), RangeTest(v, float32(0))
            // Value test needed: float32 has 24 bits, int64 has 64 bits
        case uint:
            return float32(v), RangeTest(v, float32(0))
            // Value test needed: float32 has 24 bits, uint has 32 bits
        case uint8:
            return float32(v), nil
            // a float32 can exactly represent any uint8
        case uint16:
            return float32(v), nil
            // a float32 can exactly represent any uint16
        case uint32:
            return float32(v), RangeTest(v, float32(0))
            // Value test needed: float32 has 24 bits, uint32 has 32 bits
        case uint64:
            return float32(v), RangeTest(v, float32(0))
            // Value test needed: float32 has 24 bits, uint64 has 64 bits
        case float32:
            return v, nil
        case float64:
            return float32(v), RangeTest(v, float32(0))
            // Value test needed: float32 has 24 bits, float64 has 53 bits
        case string:
            n, err := strconv.ParseFloat(v, 32)
            return float32(n), err
        }
    case float64:
        switch v := val.(type) {
        case bool:
            if v {
                return float64(1), nil
            }
            return float64(0), nil
        case int:
            return float64(v), nil
            // a float64 can exactly represent any int
        case int8:
            return float64(v), nil
            // a float64 can exactly represent any int8
        case int16:
            return float64(v), nil
            // a float64 can exactly represent any int16
        case int32:
            return float64(v), nil
            // a float64 can exactly represent any int32
        case int64:
            return float64(v), RangeTest(v, float64(0))
            // Value test needed: float64 has 53 bits, int64 has 64 bits
        case uint:
            return float64(v), nil
            // a float64 can exactly represent any uint
        case uint8:
            return float64(v), nil
            // a float64 can exactly represent any uint8
        case uint16:
            return float64(v), nil
            // a float64 can exactly represent any uint16
        case uint32:
            return float64(v), nil
            // a float64 can exactly represent any uint32
        case uint64:
            return float64(v), RangeTest(v, float64(0))
            // Value test needed: float64 has 53 bits, uint64 has 64 bits
        case float32:
            return float64(v), nil
            // a 24-bit float is always representable in a 53-bit float
        case float64:
            return v, nil
        case string:
            n, err := strconv.ParseFloat(v, 64)
            return float64(n), err
        }
    case string:
        switch v := val.(type) {
        case bool:
            return strconv.FormatBool(v), nil
        case int:
            return strconv.FormatInt(int64(v), 10), nil
        case int8:
            return strconv.FormatInt(int64(v), 10), nil
        case int16:
            return strconv.FormatInt(int64(v), 10), nil
        case int32:
            return strconv.FormatInt(int64(v), 10), nil
        case int64:
            return strconv.FormatInt(int64(v), 10), nil
        case uint:
            return strconv.FormatUint(uint64(v), 10), nil
        case uint8:
            return strconv.FormatUint(uint64(v), 10), nil
        case uint16:
            return strconv.FormatUint(uint64(v), 10), nil
        case uint32:
            return strconv.FormatUint(uint64(v), 10), nil
        case uint64:
            return strconv.FormatUint(uint64(v), 10), nil
        case float32:
            return strconv.FormatFloat(float64(v), 'g', -1, 32), nil
        case float64:
            return strconv.FormatFloat(float64(v), 'g', -1, 64), nil
        case string:
            return v, nil
        }
    }
    return nil, fmt.Errorf("no convertible value")
}
