#!/bin/bash

# download the code and all its dependencies
_go_get() {
    # Note: keep "github.com/ugorji/go/codec" in quotes, as script pushing to github will replace it appropriately
    ${go[@]} get -u \
       "github.com/ugorji/go/codec" \
       github.com/tinylib/msgp/msgp github.com/tinylib/msgp \
       github.com/pquerna/ffjson/ffjson github.com/pquerna/ffjson \
       github.com/Sereal/Sereal/Go/sereal \
       bitbucket.org/bodhisnarkva/cbor/go \
       github.com/fxamacker/cbor/v2 \
       github.com/davecgh/go-xdr/xdr2 \
       gopkg.in/mgo.v2/bson \
       gopkg.in/vmihailenco/msgpack.v2 \
       github.com/json-iterator/go \
       go.mongodb.org/mongo-driver/bson \
       github.com/globalsign/mgo/bson \
       github.com/goccy/go-json \
       github.com/vmihailenco/msgpack/v4 \
       github.com/mailru/easyjson/...
}

# add generated tag to the top of each file
_prependbt() {
    local tag="${3:-generated}"
    cat > ${2} <<EOF
// +build ${tag}

EOF
    cat ${1} >> ${2}
    rm -f ${1}
}

_sed_in_file() {
    case "$(uname -s)" in
        Darwin*) sed -i '' "$@";;
        *) sed -i "$@";;
    esac
}

_check_shared_files() {
    local d=${1:-../../go/codec}
    for i in values_test.go init_test.go codec_init_test.go bench_test.go codec_bench_test.go; do
        echo $i
        diff -s $i $d/$i
        echo '........................'
    done
}

# To run the full suite of benchmarks, including executing against the external frameworks
# listed above, you MUST first run code generation for the frameworks that support it.
#
# If you want to run the benchmarks against code generated values.
# Then first generate the code generated values from values_test.go named typed.
# we cannot normally read a _test.go file, so temporarily copy it into a readable file.
_gen() {
    local zsfx="_generated_test.go"
    # local z=`pwd`
    # z=${z%%/src/*}
    # NOTE: MARKER: ffjson has been generating bad uncompilable code, so ignore it for now
    cp values_test.go v.go &&
        echo "msgp ... " &&
        msgp -unexported -tests=false -o=m9.go -file=v.go &&
        _prependbt m9.go values_msgp${zsfx} &&
        echo "easyjson ... " &&
        easyjson -all -no_std_marshalers -omit_empty -output_filename e9.go v.go &&
        _prependbt e9.go values_easyjson${zsfx} &&
        echo "ffjson ... " && 
        ffjson -force-regenerate -reset-fields -w f9.go v.go &&
        _prependbt f9.go values_ffjson${zsfx} ignore &&
        _sed_in_file -e 's+ MarshalJSON(+ _MarshalJSON(+g' values_ffjson${zsfx} &&
        _sed_in_file -e 's+ UnmarshalJSON(+ _UnmarshalJSON(+g' values_ffjson${zsfx} &&
        rm -f easyjson-bootstrap*.go ffjson-inception* &&
        rm -f v.go &&
        echo "... DONE"
}

# run the full suite of tests
#
# Basically, its a sequence of
# ${go[@]} test -tags "alltests x codec.safe" -bench "CodecSuite or AllSuite or XSuite" -benchmem
#

_suite_tests() {
    if [[ "${do_x}" = "1" ]]; then
        printf "\n==== X Baseline ====\n"
        ${go[@]} test "${zargs[@]}" -tags x -v "$@"
    else
        printf "\n==== Baseline ====\n"
        ${go[@]} test "${zargs[@]}" -v "$@"
    fi
    if [[ "${do_x}" = "1" ]]; then
        printf "\n==== X Generated ====\n"
        ${go[@]} test "${zargs[@]}" -tags "x generated" -v "$@"
    else
        printf "\n==== Generated ====\n"
        ${go[@]} test "${zargs[@]}" -tags "generated" -v "$@"
    fi
}

_suite_tests_strip_file_line() {
    # sed -e 's/^\([^a-zA-Z0-9]\+\)[a-zA-Z0-9_]\+\.go:[0-9]\+:/\1/'
    sed -e 's/[a-zA-Z0-9_]*.go:[0-9]*://g'
}

_suite_any() {
    local x="$1"
    local g="$2"
    local b="$3"
    shift 3
    local a=( "" "codec.safe"  "codec.notfastpath" "codec.notfastpath codec.safe" )
    if [[ "$g" = "g" ]]; then a=( "generated" "generated codec.safe"); fi
    for i in "${a[@]}"; do
        echo ">>>> bench TAGS: 'alltests $x $i' SUITE: $b"
        ${go[@]} test "${zargs[@]}" -tags "alltests $x $i" -bench "$b" -benchmem "$@"
    done
}

_suite_very_quick_json_non_suite() {
    local t="${1:-x}"
    shift
    echo ">>>> very quick json bench"
    local tags=("${t}" "${t} generated" "${t} codec.safe" "${t} generated codec.safe" "${t} codec.notfastpath")
    local js=( En De )
    for t in "${tags[@]}"; do
        echo "---- tags: ${t} ----"
        local b="Json"
        if [[ "${t}" =~ x && ! "${t}" =~ safe && ! "${t}" =~ notfastpath ]]; then
            b="Json|Std_Json|JsonIter|GoccyJson"
            if [[ "${t}" =~ generated ]]; then b="Json|Easyjson"; fi
        fi            
        for j in "${js[@]}"; do
            ${go[@]} test "${zargs[@]}" -tags "${t}" -bench "__(${b})__.*${j}" -benchmem "$@"
            [[ "${b}" != Json ]] && echo # echo if more than 1 line is printed
        done
    done
}

_suite_very_quick_benchmark() {
    local a="${1:-Json}" # Json Cbor Msgpack Simple
    # case "$1" in
    #     Json|Cbor|Msgpack|Simple|Binc) a="${1}"; shift ;;
    # esac
    local b="${2}"
    local c="${3}"
    local t="${4:-4s}" # 4s 1x
    shift 4
    ${go[@]} test "${zargs[@]}" -tags "alltests ${c}" -bench "__${a}__.*${b}" -benchmem -benchtime "${t}" "$@"
}

_suite_trim_output() {
    grep -v -E "^(goos:|goarch:|pkg:|cpu:|PASS|ok|=== RUN|--- PASS)"
}

_bench_dot_out_dot_txt() {
  printf "**** STATS ****\n\n"
  ./bench.sh -tx
  printf "**** SUITE **** (without libs doing code generation)\n\n"
  ./bench.sh -sx
  printf "**** SUITE **** (with libs doing code generation) ****\n\n"
  ./bench.sh -sgx
}

_suite_debugging() {
    _suite_very_quick_benchmark "$@"
}

_usage() {
    printf "usage: bench.sh -[dcbsgjqpz] for \n"
    printf "\t-d download\n"
    printf "\t-c code-generate\n"
    printf "\t-e compare files shared with go/codec\n"
    printf "\t-tx tests (show stats for each format and whether encoded == decoded); if x, do external also\n"
    printf "\t-sgx run test suite for codec; if g, use generated files; if x, do external also\n"
    printf "\t-jq run test suite for [json, json-quick]\n"
    printf "\t-p run test suite with profiles: defaults to json: [format/prefix] [suffix] [tags] [benchtime]\n"
    printf "\t-z run tests for bench.out.txt\n"
    printf "\t-f [pprof file] run pprof\n"
    printf "\t-y run debugging suite (during development only): [format/prefix] [suffix] [tags] [benchtime]\n"
}

_main() {
    if [[ "$1" == "" || "$1" == "-h" || "$1" == "-?" ]]; then
        _usage
        return 1
    fi
    
    # export GODEBUG=asyncpreemptoff=1 # TODO remove
    
    local go=( "${MYGOCMD:-go}" )
    local zargs=("-count" "1" "-tr")
    local args=()
    local do_x="0"
    local do_g="0"
    while getopts "dcbsjqptxklgzfye" flag
    do
        case "$flag" in
            d|c|b|s|j|q|p|t|x|k|l|g|z|f|y|e) args+=( "$flag" ) ;;
            *) _usage; return 1 ;;
        esac
    done
    shift "$((OPTIND-1))"
    
    [[ " ${args[*]} " == *"x"* ]] && do_x="1"
    [[ " ${args[*]} " == *"g"* ]] && do_g="1"
    [[ " ${args[*]} " == *"k"* ]] && zargs+=("-gcflags" "all=-B")
    [[ " ${args[*]} " == *"l"* ]] && zargs+=("-gcflags" "all=-l=4")
    [[ " ${args[*]} " == *"d"* ]] && _go_get "$@"
    [[ " ${args[*]} " == *"c"* ]] && _gen "$@"
    
    [[ " ${args[*]} " == *"s"* && "${do_x}" == 0 && "${do_g}" == 0 ]] && _suite_any - - BenchmarkCodecSuite "$@" | _suite_trim_output
    [[ " ${args[*]} " == *"s"* && "${do_x}" == 0 && "${do_g}" == 1 ]] && _suite_any - g BenchmarkCodecSuite "$@" | _suite_trim_output
    [[ " ${args[*]} " == *"s"* && "${do_x}" == 1 && "${do_g}" == 0 ]] && _suite_any x - BenchmarkCodecXSuite "$@" | _suite_trim_output
    [[ " ${args[*]} " == *"s"* && "${do_x}" == 1 && "${do_g}" == 1 ]] && _suite_any x g BenchmarkCodecXGenSuite "$@" | _suite_trim_output
    
    [[ " ${args[*]} " == *"j"* ]] && _suite_any x - BenchmarkCodecQuickAllJsonSuite "$@" | _suite_trim_output

    # These are some very specific suites
    # [[ " ${args[*]} " == *"q"* ]] && _suite_very_quick_json_via_suite "$@" | _suite_trim_output
    [[ " ${args[*]} " == *"q"* ]] && _suite_very_quick_json_non_suite "$@" | _suite_trim_output

    # These are just helpers (not really running benchmark suites)
    [[ " ${args[*]} " == *"t"* ]] && _suite_tests "$@" | _suite_trim_output | _suite_tests_strip_file_line
    [[ " ${args[*]} " == *"p"* ]] && zargs+=("-cpuprofile" "cpu.out" "-memprofile" "mem.out" "-memprofilerate" "1") && _suite_very_quick_benchmark "$@" | _suite_trim_output
    [[ " ${args[*]} " == *"f"* ]] && ${go[@]} tool pprof bench.test ${1:-mem.out}
    [[ " ${args[*]} " == *"z"* ]] && _bench_dot_out_dot_txt
    [[ " ${args[*]} " == *"y"* ]] && _suite_debugging "$@" | _suite_trim_output
    [[ " ${args[*]} " == *"e"* ]] && _check_shared_files "$@"
    
    true
    # shift $((OPTIND-1))
}

if [ "." = `dirname $0` ]
then
    _main "$@"
else
    echo "bench.sh must be run from the directory it resides in"
    _usage
fi 
