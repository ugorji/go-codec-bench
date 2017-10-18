#!/bin/bash
_bench() {
    if [[ "$zquick" == "1" ]]; then
        # echo "RUNNING IN QUICK MODE"
        for k in En De; do
            for j in 1 2; do
                for y in "" "-ti=1024"; do
                    for i in $(seq 1 1); do
                        echo ">>>>>> go test -bd $j -run Melody -benchmem -benchtime ${j}s -tags x  -bench _(Json).*$k $y"
                        go test -bd $j -run Melody -benchmem -benchtime "${j}s" -tags x  -bench "_(Json).*$k" "$y"
                    done
                done
            done
        done
        return
    fi
    
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
