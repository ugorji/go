#!/bin/bash

# download the code and all its dependencies 
_go_get() {
    go get -u \
       "github.com/ugorji/go/codec" "github.com/ugorji/go/codec"/codecgen \
       github.com/tinylib/msgp/msgp github.com/tinylib/msgp \
       github.com/pquerna/ffjson/ffjson github.com/pquerna/ffjson \
       github.com/Sereal/Sereal/Go/sereal \
       bitbucket.org/bodhisnarkva/cbor/go \
       github.com/davecgh/go-xdr/xdr2 \
       gopkg.in/mgo.v2/bson \
       gopkg.in/vmihailenco/msgpack.v2 \
       github.com/json-iterator/go \
       github.com/mailru/easyjson/...
}

# add generated tag to the top of each file
_prependbt() {
    cat > ${2} <<EOF
// +build generated

EOF
    cat ${1} >> ${2}
    rm -f ${1}
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
    # Note: ensure you run the codecgen for this codebase
    cp values_test.go v.go &&
        echo "codecgen ..." &&
        codecgen -nx -rt codecgen -t 'codecgen generated' -o values_codecgen${zsfx} -d 19780 v.go &&
        echo "msgp ... " &&
        msgp -unexported -tests=false -o=m9.go -file=v.go &&
        _prependbt m9.go values_msgp${zsfx} &&
        echo "easyjson ... " &&
        easyjson -all -no_std_marshalers -omit_empty -output_filename e9.go v.go &&
        _prependbt e9.go values_easyjson${zsfx} &&
        echo "ffjson ... " && 
        ffjson -force-regenerate -reset-fields -w f9.go v.go &&
        _prependbt f9.go values_ffjson${zsfx} &&
        sed -i '' -e 's+ MarshalJSON(+ _MarshalJSON(+g' values_ffjson${zsfx} &&
        sed -i '' -e 's+ UnmarshalJSON(+ _UnmarshalJSON(+g' values_ffjson${zsfx} &&
        rm -f easyjson-bootstrap*.go ffjson-inception* &&
        rm -f v.go &&
        echo "... DONE"
}

# run the full suite of tests
#
# Basically, its a sequence of
# go test -tags "alltests x safe codecgen generated" -bench "CodecSuite or AllSuite or XSuite" -benchmem
# 
_suite() {
    local t="alltests x"
    local a=( "" "safe"  "notfastpath" "notfastpath safe" "codecgen" "codecgen safe")
    local b=( "generated" "generated safe")
    for i in "${a[@]}"
    do
        echo ">>>> bench TAGS: '$t $i' SUITE: BenchmarkCodecXSuite"
        go test -run Nothing -tags "$t $i" -bench BenchmarkCodecXSuite -benchmem "$@"
    done
    for i in "${b[@]}"
    do
        echo ">>>> bench TAGS: '$t $i' SUITE: BenchmarkCodecXGenSuite"
        go test -run Nothing -tags "$t $i" -bench BenchmarkCodecXGenSuite -benchmem "$@"
    done
}

_usage() {
    echo "usage: bench.sh -[dcs] for [download, code-generate and suite-of-tests] respectively"
}

_main() {
    if [[ "$1" == "" ]]
    then
        _usage
        return 1
    fi
    while getopts "dcs" flag
    do
        case "x$flag" in
            'xd') shift; _go_get "$@" ;;
            'xc') shift; _gen "$@" ;;
            'xs') shift; _suite "$@" ;;
            *) shift; _usage; return 1 ;;
        esac
    done
    # shift $((OPTIND-1))
}

if [ "." = `dirname $0` ]
then
    _main "$@"
else
    echo "bench.sh must be run from the directory it resides in"
    _usage
fi 
