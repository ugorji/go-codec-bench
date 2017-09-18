#!/bin/bash
_bench() {
    local t=alltests
    local s=CodecSuite
    if [[ "x$1" != "x" ]]; then
        s="$1"; shift;
        if [[ $s == *X* || $s == *All* ]]; then t="$t x"; fi
    fi
    local a=( "default" "safe"  "notfastpath" "codecgen" "codecgen safe" "notfastpath safe" )
    for i in "${a[@]}"
    do
        echo ">>>> TAGS: '$t $i' SUITE: '$s'"
        go test "-tags=$t $i" -bench "$s" -benchmem "$@"
        # -cpuprofile=cpu.out -memprofile=mem.out -test.memprofilerate=1 
    done
}
_bench "$@"
