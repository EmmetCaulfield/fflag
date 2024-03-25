#!/usr/bin/env bash

set -Eeuo pipefail

readonly -a basetypes=(bool int int8 int16 int32 int64 uint uint8 uint16 uint32 uint64 float32 float64 string)
readonly -A parsefunc=([bool]='ParseBool(v)'
                       [int]='ParseInt(v, 10, strconv.IntSize)'
                       [int8]='ParseInt(v, 10, 8)'
                       [int16]='ParseInt(v, 10, 16)'
                       [int32]='ParseInt(v, 10, 32)'
                       [int64]='ParseInt(v, 10, 64)'
                       [uint]='ParseUint(v, 10, strconv.IntSize)'
                       [uint8]='ParseUint(v, 10, 8)'
                       [uint16]='ParseUint(v, 10, 16)'
                       [uint32]='ParseUint(v, 10, 32)'
                       [uint64]='ParseUint(v, 10, 64)'
                       [float32]='ParseFloat(v, 32)'
                       [float64]='ParseFloat(v, 64)')
readonly -A formatfun=([bool]='FormatBool(v)'
                       [int]='FormatInt(int64(v), 10)'
                       [int8]='FormatInt(int64(v), 10)'
                       [int16]='FormatInt(int64(v), 10)'
                       [int32]='FormatInt(int64(v), 10)'
                       [int64]='FormatInt(int64(v), 10)'
                       [uint]='FormatUint(uint64(v), 10)'
                       [uint8]='FormatUint(uint64(v), 10)'
                       [uint16]='FormatUint(uint64(v), 10)'
                       [uint32]='FormatUint(uint64(v), 10)'
                       [uint64]='FormatUint(uint64(v), 10)'
                       [float32]="FormatFloat(float64(v), 'g', -1, 32)"
                       [float64]="FormatFloat(float64(v), 'g', -1, 64)")

base() {
    echo "$1" | sed 's/[1368][246]\?$//'
}

nbits() {
    case "$1" in
        float32 )
            echo 24
            ;;
        float64 )
            echo 53
            ;;
        int | uint )
            # At least...
            echo 32
            ;;
        * )
            echo "$1" | sed 's/[^0-9]\+//'
    esac
}


rangecheck() {
    local -r valt="$1"; shift # Value type, converting FROM
    local -r reft="$1"; shift # Reference type, converting TO
    local -r varn="$1"; shift # Variable name being converted

    local -r refb=$(base $reft)
    local -r valb=$(base $valt)
    local -i refn=$(nbits $reft)
    local -i valn=$(nbits $valt)

    # If the base types are the same and the number of bits in the
    # reference type at least equals the number of bits in the value
    # type, the overflow can't happen
    if [ $refb == $valb ] && [ $refn -ge $valn ]; then
        echo nil
        echo "            // a $valn-bit $valb is always representable in a $refn-bit $refb"
        return
    fi
    # If the reference type is int and the value type is uint,
    # overflow can't happen if there's even one more bit in the
    # reference type
    if [ $refb == 'int' ] && [ $valb == 'uint' ] && [ $refn -gt $valn ]; then
        echo nil
        echo "            // $valb$valn is always representable in $refb$refn"
        return
    fi
    # If the reference type is float and it has at least one more bit
    # than the value type, it can be exactly represented:
    if [ $refb == 'float' ] && [ $refn -gt $valn ]; then
        echo nil
        echo "            // a $reft can exactly represent any $valt"
        return
    fi
    echo "RangeTest($varn, $reft(0))"
    echo "            // Value test needed: $reft has $refn bits, $valt has $valn bits"
}

coerce () {
    for r in ${basetypes[*]}; do
        cat <<EOF
    case ${r}:
        switch v := val.(type) {
EOF
        for v in ${basetypes[*]}; do
            cat<<EOF    
        case ${v}:
EOF
            if [ "$v" = "$r" ]; then
                cat<<EOF
            return v, nil
EOF
            elif [ "$v" = "string" ]; then
                cat<<EOF
            n, err := strconv.${parsefunc[$r]}
            return ${r}(n), err
EOF
            elif [ "$r" = "string" ]; then
                cat<<EOF
            return strconv.${formatfun[$v]}, nil
EOF
            elif [ "$r" = "bool" ]; then
                cat<<EOF
            return v != ${v}(0), nil
EOF
            elif [ "$v" = "bool" ]; then
                cat<<EOF
            if v {
                return ${r}(1), nil
            }
            return ${r}(0), nil
EOF
            else
                cat<<EOF
            return ${r}(v), $(rangecheck $v $r v)
EOF
            fi
        done
    cat <<EOF
        }
EOF
    done
}



header() {
    cat coerce.go-top
    cat<<EOF
func CoerceScalar(ref interface{}, val interface{}) (interface{}, error) {
    if ref == nil || val == nil {
        return nil, fmt.Errorf("nil argument given")
    }
    switch ref.(type) {
EOF
}

trailer() {
    cat<<EOF
    }
    return nil, fmt.Errorf("no convertible value")
}
EOF
}

header
coerce
trailer
