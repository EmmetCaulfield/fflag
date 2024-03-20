#!/usr/bin/env bash

set -Eeuo pipefail

stringRoundtrip () {
    local -r type="$1"; shift
    local -r scalar="$1"; shift
    local -r vector="$1"; shift
    local -r sep="$1"; shift

    cat<<EOF
func TestStrConv_${type}(t *testing.T) {
    a := ${type}(${scalar})
    var b ${type}
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

    sa := []${type}{${vector}}
    sb := []${type}{}
    s = StrConv(sa, WithSep("${sep}"))
    err := FromStr(&sb, s, WithSep("${sep}"))
    if err != nil {
        t.Errorf("error converting '%s': %v", s, err)
    }
    for i, v := range(sa) {
        if v != sb[i] {
            t.Errorf("failed to roundtrip %v via '%s', got %v", sa, s, sb)
            break
        }
    }

    s = StrConv(&sa, WithSep("${sep}"))
    sb = sb[:0]
    err = FromStr(&sb, s, WithSep("${sep}"))
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

EOF
}

sliceLengthAndItem () {
    local -r  type="$1"; shift
    local -ir index="$1"; shift
    local -r  list="$1"; shift
    local -ir len=$(echo $(echo "${list}" | tr -cd , | wc -c)+1 | bc)
    local -r  item=$(echo "${list}" | cut -d, -f$((${index}+1)))

    cat<<EOF
func TestSliceOps_${type}(t *testing.T) {
    slice := []${type}{${list}}
    length := SliceLen(slice)
    if length != ${len} {
        t.Errorf("wrong length for slice %v, expected ${len}, got %d", slice, length)
    }
    length = SliceLen(&slice)
    if length != ${len} {
        t.Errorf("wrong length for slice %v, expected ${len}, got %d", &slice, length)
    }

    item := ItemAt(slice, ${index}).(${type})
    if item != ${item} {
        t.Errorf(\`wrong value item at ${index} in %v, expected ${item}, got %v\`, slice, item)
    }
    item = ItemAt(&slice, ${index}).(${type})
    if item != ${item} {
        t.Errorf(\`wrong value item at ${index} in %v, expected ${item}, got %v\`, &slice, item)
    }
}

EOF
}



header() {
    cat<<EOF
package types
import(
    "testing"
)

EOF
}

header

sliceLengthAndItem bool     1  true,false,true
sliceLengthAndItem int      0  1,2,3
sliceLengthAndItem int8     1  1,2,3
sliceLengthAndItem int16    5  1,2,3,4,5,6,7,8,9
sliceLengthAndItem int32    1  1,2,3
sliceLengthAndItem int64    1  1,2,3
sliceLengthAndItem uint     1  1,2,3
sliceLengthAndItem uint8    1  1,2,3
sliceLengthAndItem uint16   1  1,2,3
sliceLengthAndItem uint32   1  1,2,3
sliceLengthAndItem uint64   1  1,2,3
sliceLengthAndItem float32  1  1.0,2.5,3.0
sliceLengthAndItem float64  2  1.0,2.5,3.0
sliceLengthAndItem byte     1  1,2,3
sliceLengthAndItem rune     1  1,2,3
sliceLengthAndItem string   1  '"A","B","C"'

stringRoundtrip "bool"   true "true,false,true", "="
stringRoundtrip "int"     -3  "3,2,1",   "/"
stringRoundtrip "int8"     3  "3,+2,1",  "|"
stringRoundtrip "int16"   -9  "-3,2,-1", "%"
stringRoundtrip "int32"    3  "3,-2,1",  "###"
stringRoundtrip "int64"    3  "3,2,1",   "/"
stringRoundtrip "uint"     3  "3,+2,1",  "/"
stringRoundtrip "uint8"    3  "3,2,1",   "/"
stringRoundtrip "uint16"   3  "3,2,1",   "/"
stringRoundtrip "uint32"   3  "3,2,1",   "/"
stringRoundtrip "uint64"   3  "3,2,1",   "/"
stringRoundtrip "float32"  3  "3,2,1",   "/"
stringRoundtrip "float64"  3  "3,2,1",   "/"
stringRoundtrip "byte"     3  "3,2,1",   "/"
stringRoundtrip "rune"     3  "3,2,1",   "/"
stringRoundtrip "string"  '"foo"' '"foo","bar","baz"',   "/"

