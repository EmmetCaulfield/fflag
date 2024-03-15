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
    err = FromStr(&sb, s, WithSep("${sep}"))
    if err == nil {
        t.Errorf("unexpected success rountripping %v via '%s' with nonempty slice %v", &sa, s, sb)
    }
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

header() {
    cat<<EOF
package types
import(
    "testing"
)

EOF
}

header
stringRoundtrip "bool"   true "true,false,true", "="
stringRoundtrip "int"    -3 "3,2,1",   "/"
stringRoundtrip "int8"    3 "3,+2,1",  "|"
stringRoundtrip "int16"  -9 "-3,2,-1", "%"
stringRoundtrip "int32"   3 "3,-2,1",  "###"
stringRoundtrip "int64"   3 "3,2,1",   "/"
stringRoundtrip "uint"    3 "3,+2,1",  "/"
stringRoundtrip "uint8"   3 "3,2,1",   "/"
stringRoundtrip "uint16"  3 "3,2,1",   "/"
stringRoundtrip "uint32"  3 "3,2,1",   "/"
stringRoundtrip "uint64"  3 "3,2,1",   "/"
stringRoundtrip "byte"    3 "3,2,1",   "/"
stringRoundtrip "rune"    3 "3,2,1",   "/"
stringRoundtrip "string"  '"foo"' '"foo","bar","baz"',   "/"
