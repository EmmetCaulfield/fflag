package types
import(
    "testing"
)

func TestSliceOps_bool(t *testing.T) {
    slice := []bool{true,false,true}
    length := SliceLen(slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", slice, length)
    }
    length = SliceLen(&slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", &slice, length)
    }

    item := ItemAt(slice, 1).(bool)
    if item != false {
        t.Errorf(`wrong value item at 1 in %v, expected false, got %v`, slice, item)
    }
    item = ItemAt(&slice, 1).(bool)
    if item != false {
        t.Errorf(`wrong value item at 1 in %v, expected false, got %v`, &slice, item)
    }
}

func TestSliceOps_int(t *testing.T) {
    slice := []int{1,2,3}
    length := SliceLen(slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", slice, length)
    }
    length = SliceLen(&slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", &slice, length)
    }

    item := ItemAt(slice, 0).(int)
    if item != 1 {
        t.Errorf(`wrong value item at 0 in %v, expected 1, got %v`, slice, item)
    }
    item = ItemAt(&slice, 0).(int)
    if item != 1 {
        t.Errorf(`wrong value item at 0 in %v, expected 1, got %v`, &slice, item)
    }
}

func TestSliceOps_int8(t *testing.T) {
    slice := []int8{1,2,3}
    length := SliceLen(slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", slice, length)
    }
    length = SliceLen(&slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", &slice, length)
    }

    item := ItemAt(slice, 1).(int8)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, slice, item)
    }
    item = ItemAt(&slice, 1).(int8)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, &slice, item)
    }
}

func TestSliceOps_int16(t *testing.T) {
    slice := []int16{1,2,3,4,5,6,7,8,9}
    length := SliceLen(slice)
    if length != 9 {
        t.Errorf("wrong length for slice %v, expected 9, got %d", slice, length)
    }
    length = SliceLen(&slice)
    if length != 9 {
        t.Errorf("wrong length for slice %v, expected 9, got %d", &slice, length)
    }

    item := ItemAt(slice, 5).(int16)
    if item != 6 {
        t.Errorf(`wrong value item at 5 in %v, expected 6, got %v`, slice, item)
    }
    item = ItemAt(&slice, 5).(int16)
    if item != 6 {
        t.Errorf(`wrong value item at 5 in %v, expected 6, got %v`, &slice, item)
    }
}

func TestSliceOps_int32(t *testing.T) {
    slice := []int32{1,2,3}
    length := SliceLen(slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", slice, length)
    }
    length = SliceLen(&slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", &slice, length)
    }

    item := ItemAt(slice, 1).(int32)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, slice, item)
    }
    item = ItemAt(&slice, 1).(int32)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, &slice, item)
    }
}

func TestSliceOps_int64(t *testing.T) {
    slice := []int64{1,2,3}
    length := SliceLen(slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", slice, length)
    }
    length = SliceLen(&slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", &slice, length)
    }

    item := ItemAt(slice, 1).(int64)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, slice, item)
    }
    item = ItemAt(&slice, 1).(int64)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, &slice, item)
    }
}

func TestSliceOps_uint(t *testing.T) {
    slice := []uint{1,2,3}
    length := SliceLen(slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", slice, length)
    }
    length = SliceLen(&slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", &slice, length)
    }

    item := ItemAt(slice, 1).(uint)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, slice, item)
    }
    item = ItemAt(&slice, 1).(uint)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, &slice, item)
    }
}

func TestSliceOps_uint8(t *testing.T) {
    slice := []uint8{1,2,3}
    length := SliceLen(slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", slice, length)
    }
    length = SliceLen(&slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", &slice, length)
    }

    item := ItemAt(slice, 1).(uint8)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, slice, item)
    }
    item = ItemAt(&slice, 1).(uint8)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, &slice, item)
    }
}

func TestSliceOps_uint16(t *testing.T) {
    slice := []uint16{1,2,3}
    length := SliceLen(slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", slice, length)
    }
    length = SliceLen(&slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", &slice, length)
    }

    item := ItemAt(slice, 1).(uint16)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, slice, item)
    }
    item = ItemAt(&slice, 1).(uint16)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, &slice, item)
    }
}

func TestSliceOps_uint32(t *testing.T) {
    slice := []uint32{1,2,3}
    length := SliceLen(slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", slice, length)
    }
    length = SliceLen(&slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", &slice, length)
    }

    item := ItemAt(slice, 1).(uint32)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, slice, item)
    }
    item = ItemAt(&slice, 1).(uint32)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, &slice, item)
    }
}

func TestSliceOps_uint64(t *testing.T) {
    slice := []uint64{1,2,3}
    length := SliceLen(slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", slice, length)
    }
    length = SliceLen(&slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", &slice, length)
    }

    item := ItemAt(slice, 1).(uint64)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, slice, item)
    }
    item = ItemAt(&slice, 1).(uint64)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, &slice, item)
    }
}

func TestSliceOps_float32(t *testing.T) {
    slice := []float32{1.0,2.5,3.0}
    length := SliceLen(slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", slice, length)
    }
    length = SliceLen(&slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", &slice, length)
    }

    item := ItemAt(slice, 1).(float32)
    if item != 2.5 {
        t.Errorf(`wrong value item at 1 in %v, expected 2.5, got %v`, slice, item)
    }
    item = ItemAt(&slice, 1).(float32)
    if item != 2.5 {
        t.Errorf(`wrong value item at 1 in %v, expected 2.5, got %v`, &slice, item)
    }
}

func TestSliceOps_float64(t *testing.T) {
    slice := []float64{1.0,2.5,3.0}
    length := SliceLen(slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", slice, length)
    }
    length = SliceLen(&slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", &slice, length)
    }

    item := ItemAt(slice, 2).(float64)
    if item != 3.0 {
        t.Errorf(`wrong value item at 2 in %v, expected 3.0, got %v`, slice, item)
    }
    item = ItemAt(&slice, 2).(float64)
    if item != 3.0 {
        t.Errorf(`wrong value item at 2 in %v, expected 3.0, got %v`, &slice, item)
    }
}

func TestSliceOps_byte(t *testing.T) {
    slice := []byte{1,2,3}
    length := SliceLen(slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", slice, length)
    }
    length = SliceLen(&slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", &slice, length)
    }

    item := ItemAt(slice, 1).(byte)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, slice, item)
    }
    item = ItemAt(&slice, 1).(byte)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, &slice, item)
    }
}

func TestSliceOps_rune(t *testing.T) {
    slice := []rune{1,2,3}
    length := SliceLen(slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", slice, length)
    }
    length = SliceLen(&slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", &slice, length)
    }

    item := ItemAt(slice, 1).(rune)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, slice, item)
    }
    item = ItemAt(&slice, 1).(rune)
    if item != 2 {
        t.Errorf(`wrong value item at 1 in %v, expected 2, got %v`, &slice, item)
    }
}

func TestSliceOps_string(t *testing.T) {
    slice := []string{"A","B","C"}
    length := SliceLen(slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", slice, length)
    }
    length = SliceLen(&slice)
    if length != 3 {
        t.Errorf("wrong length for slice %v, expected 3, got %d", &slice, length)
    }

    item := ItemAt(slice, 1).(string)
    if item != "B" {
        t.Errorf(`wrong value item at 1 in %v, expected "B", got %v`, slice, item)
    }
    item = ItemAt(&slice, 1).(string)
    if item != "B" {
        t.Errorf(`wrong value item at 1 in %v, expected "B", got %v`, &slice, item)
    }
}

func TestStrConv_bool(t *testing.T) {
    a := bool(true)
    var b bool
    s := StrConv(a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    s = StrConv(&a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    sa := []bool{true,false,true,}
    sb := []bool{}
    s = StrConv(sa, WithSep("="))
    err := FromStr(&sb, s, WithSep("="))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", sa, s, sb)
            break
        }
    }

    s = StrConv(&sa, WithSep("="))
    err = FromStr(&sb, s, WithSep("="))
    if err == nil {
        t.Errorf("unexpected success rountripping %v via '%s' with nonempty slice %v", &sa, s, sb)
    }
    sb = sb[:0]
    err = FromStr(&sb, s, WithSep("="))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", &sa, s, sb)
            break
        }
    }
}

func TestStrConv_int(t *testing.T) {
    a := int(-3)
    var b int
    s := StrConv(a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    s = StrConv(&a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    sa := []int{3,2,1,}
    sb := []int{}
    s = StrConv(sa, WithSep("/"))
    err := FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", sa, s, sb)
            break
        }
    }

    s = StrConv(&sa, WithSep("/"))
    err = FromStr(&sb, s, WithSep("/"))
    if err == nil {
        t.Errorf("unexpected success rountripping %v via '%s' with nonempty slice %v", &sa, s, sb)
    }
    sb = sb[:0]
    err = FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", &sa, s, sb)
            break
        }
    }
}

func TestStrConv_int8(t *testing.T) {
    a := int8(3)
    var b int8
    s := StrConv(a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    s = StrConv(&a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    sa := []int8{3,+2,1,}
    sb := []int8{}
    s = StrConv(sa, WithSep("|"))
    err := FromStr(&sb, s, WithSep("|"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", sa, s, sb)
            break
        }
    }

    s = StrConv(&sa, WithSep("|"))
    err = FromStr(&sb, s, WithSep("|"))
    if err == nil {
        t.Errorf("unexpected success rountripping %v via '%s' with nonempty slice %v", &sa, s, sb)
    }
    sb = sb[:0]
    err = FromStr(&sb, s, WithSep("|"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", &sa, s, sb)
            break
        }
    }
}

func TestStrConv_int16(t *testing.T) {
    a := int16(-9)
    var b int16
    s := StrConv(a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    s = StrConv(&a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    sa := []int16{-3,2,-1,}
    sb := []int16{}
    s = StrConv(sa, WithSep("%"))
    err := FromStr(&sb, s, WithSep("%"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", sa, s, sb)
            break
        }
    }

    s = StrConv(&sa, WithSep("%"))
    err = FromStr(&sb, s, WithSep("%"))
    if err == nil {
        t.Errorf("unexpected success rountripping %v via '%s' with nonempty slice %v", &sa, s, sb)
    }
    sb = sb[:0]
    err = FromStr(&sb, s, WithSep("%"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", &sa, s, sb)
            break
        }
    }
}

func TestStrConv_int32(t *testing.T) {
    a := int32(3)
    var b int32
    s := StrConv(a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    s = StrConv(&a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    sa := []int32{3,-2,1,}
    sb := []int32{}
    s = StrConv(sa, WithSep("###"))
    err := FromStr(&sb, s, WithSep("###"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", sa, s, sb)
            break
        }
    }

    s = StrConv(&sa, WithSep("###"))
    err = FromStr(&sb, s, WithSep("###"))
    if err == nil {
        t.Errorf("unexpected success rountripping %v via '%s' with nonempty slice %v", &sa, s, sb)
    }
    sb = sb[:0]
    err = FromStr(&sb, s, WithSep("###"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", &sa, s, sb)
            break
        }
    }
}

func TestStrConv_int64(t *testing.T) {
    a := int64(3)
    var b int64
    s := StrConv(a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    s = StrConv(&a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    sa := []int64{3,2,1,}
    sb := []int64{}
    s = StrConv(sa, WithSep("/"))
    err := FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", sa, s, sb)
            break
        }
    }

    s = StrConv(&sa, WithSep("/"))
    err = FromStr(&sb, s, WithSep("/"))
    if err == nil {
        t.Errorf("unexpected success rountripping %v via '%s' with nonempty slice %v", &sa, s, sb)
    }
    sb = sb[:0]
    err = FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", &sa, s, sb)
            break
        }
    }
}

func TestStrConv_uint(t *testing.T) {
    a := uint(3)
    var b uint
    s := StrConv(a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    s = StrConv(&a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    sa := []uint{3,+2,1,}
    sb := []uint{}
    s = StrConv(sa, WithSep("/"))
    err := FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", sa, s, sb)
            break
        }
    }

    s = StrConv(&sa, WithSep("/"))
    err = FromStr(&sb, s, WithSep("/"))
    if err == nil {
        t.Errorf("unexpected success rountripping %v via '%s' with nonempty slice %v", &sa, s, sb)
    }
    sb = sb[:0]
    err = FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", &sa, s, sb)
            break
        }
    }
}

func TestStrConv_uint8(t *testing.T) {
    a := uint8(3)
    var b uint8
    s := StrConv(a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    s = StrConv(&a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    sa := []uint8{3,2,1,}
    sb := []uint8{}
    s = StrConv(sa, WithSep("/"))
    err := FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", sa, s, sb)
            break
        }
    }

    s = StrConv(&sa, WithSep("/"))
    err = FromStr(&sb, s, WithSep("/"))
    if err == nil {
        t.Errorf("unexpected success rountripping %v via '%s' with nonempty slice %v", &sa, s, sb)
    }
    sb = sb[:0]
    err = FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", &sa, s, sb)
            break
        }
    }
}

func TestStrConv_uint16(t *testing.T) {
    a := uint16(3)
    var b uint16
    s := StrConv(a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    s = StrConv(&a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    sa := []uint16{3,2,1,}
    sb := []uint16{}
    s = StrConv(sa, WithSep("/"))
    err := FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", sa, s, sb)
            break
        }
    }

    s = StrConv(&sa, WithSep("/"))
    err = FromStr(&sb, s, WithSep("/"))
    if err == nil {
        t.Errorf("unexpected success rountripping %v via '%s' with nonempty slice %v", &sa, s, sb)
    }
    sb = sb[:0]
    err = FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", &sa, s, sb)
            break
        }
    }
}

func TestStrConv_uint32(t *testing.T) {
    a := uint32(3)
    var b uint32
    s := StrConv(a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    s = StrConv(&a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    sa := []uint32{3,2,1,}
    sb := []uint32{}
    s = StrConv(sa, WithSep("/"))
    err := FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", sa, s, sb)
            break
        }
    }

    s = StrConv(&sa, WithSep("/"))
    err = FromStr(&sb, s, WithSep("/"))
    if err == nil {
        t.Errorf("unexpected success rountripping %v via '%s' with nonempty slice %v", &sa, s, sb)
    }
    sb = sb[:0]
    err = FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", &sa, s, sb)
            break
        }
    }
}

func TestStrConv_uint64(t *testing.T) {
    a := uint64(3)
    var b uint64
    s := StrConv(a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    s = StrConv(&a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    sa := []uint64{3,2,1,}
    sb := []uint64{}
    s = StrConv(sa, WithSep("/"))
    err := FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", sa, s, sb)
            break
        }
    }

    s = StrConv(&sa, WithSep("/"))
    err = FromStr(&sb, s, WithSep("/"))
    if err == nil {
        t.Errorf("unexpected success rountripping %v via '%s' with nonempty slice %v", &sa, s, sb)
    }
    sb = sb[:0]
    err = FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", &sa, s, sb)
            break
        }
    }
}

func TestStrConv_float32(t *testing.T) {
    a := float32(3)
    var b float32
    s := StrConv(a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    s = StrConv(&a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    sa := []float32{3,2,1,}
    sb := []float32{}
    s = StrConv(sa, WithSep("/"))
    err := FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", sa, s, sb)
            break
        }
    }

    s = StrConv(&sa, WithSep("/"))
    err = FromStr(&sb, s, WithSep("/"))
    if err == nil {
        t.Errorf("unexpected success rountripping %v via '%s' with nonempty slice %v", &sa, s, sb)
    }
    sb = sb[:0]
    err = FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", &sa, s, sb)
            break
        }
    }
}

func TestStrConv_float64(t *testing.T) {
    a := float64(3)
    var b float64
    s := StrConv(a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    s = StrConv(&a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    sa := []float64{3,2,1,}
    sb := []float64{}
    s = StrConv(sa, WithSep("/"))
    err := FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", sa, s, sb)
            break
        }
    }

    s = StrConv(&sa, WithSep("/"))
    err = FromStr(&sb, s, WithSep("/"))
    if err == nil {
        t.Errorf("unexpected success rountripping %v via '%s' with nonempty slice %v", &sa, s, sb)
    }
    sb = sb[:0]
    err = FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", &sa, s, sb)
            break
        }
    }
}

func TestStrConv_byte(t *testing.T) {
    a := byte(3)
    var b byte
    s := StrConv(a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    s = StrConv(&a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    sa := []byte{3,2,1,}
    sb := []byte{}
    s = StrConv(sa, WithSep("/"))
    err := FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", sa, s, sb)
            break
        }
    }

    s = StrConv(&sa, WithSep("/"))
    err = FromStr(&sb, s, WithSep("/"))
    if err == nil {
        t.Errorf("unexpected success rountripping %v via '%s' with nonempty slice %v", &sa, s, sb)
    }
    sb = sb[:0]
    err = FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", &sa, s, sb)
            break
        }
    }
}

func TestStrConv_rune(t *testing.T) {
    a := rune(3)
    var b rune
    s := StrConv(a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    s = StrConv(&a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    sa := []rune{3,2,1,}
    sb := []rune{}
    s = StrConv(sa, WithSep("/"))
    err := FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", sa, s, sb)
            break
        }
    }

    s = StrConv(&sa, WithSep("/"))
    err = FromStr(&sb, s, WithSep("/"))
    if err == nil {
        t.Errorf("unexpected success rountripping %v via '%s' with nonempty slice %v", &sa, s, sb)
    }
    sb = sb[:0]
    err = FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", &sa, s, sb)
            break
        }
    }
}

func TestStrConv_string(t *testing.T) {
    a := string("foo")
    var b string
    s := StrConv(a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    s = StrConv(&a)
    FromStr(&b, s)
    if a != b {
        t.Errorf("failed to roundtrip %v via '%s', got %v", a, s, b)
    }

    sa := []string{"foo","bar","baz",}
    sb := []string{}
    s = StrConv(sa, WithSep("/"))
    err := FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", sa, s, sb)
            break
        }
    }

    s = StrConv(&sa, WithSep("/"))
    err = FromStr(&sb, s, WithSep("/"))
    if err == nil {
        t.Errorf("unexpected success rountripping %v via '%s' with nonempty slice %v", &sa, s, sb)
    }
    sb = sb[:0]
    err = FromStr(&sb, s, WithSep("/"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", &sa, s, sb)
            break
        }
    }
}

