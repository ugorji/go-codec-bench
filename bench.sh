#!/bin/bash
_bench() {
    local t="alltests"
    local s="CodecSuite"
    if [[ "x$1" != "x" ]]; then
        s="$1"; shift;
        if [[ $s == *X* ]]; then t="$t x"; fi
    fi
    local a=( "default" "safe"  "notfastpath" "notfastpath safe" "codecgen" "codecgen safe" "generated" "generated safe")
    for i in "${a[@]}"
    do
        echo ">>>> bench TAGS: '$t $i' SUITE: '$s'"
        go test -run Nothing -tags "$t $i" -bench "$s" -benchmem "$@"
        # -cpuprofile=cpu.out -memprofile=mem.out -test.memprofilerate=1 
    done
}
_bench "$@"
