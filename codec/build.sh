#!/bin/bash

# Run all the different permutations of all the tests and other things
# This helps ensure that nothing gets broken.

# # is a generation needed?
# _ng() {
#     local a="$1"
#     if [[ ! -e "$a" ]]; then echo 1; return; fi 
#     for i in `ls -1 *.go.tmpl gen.go gen_mono.go values_test.go`
#     do
#         if [[ "$a" -ot "$i" ]]; then echo 1; return; fi 
#     done
# }

_build_proceed() {
    # return success (0) if we should, and 1 (fail) if not
    if [[ "${zforce}" ]]; then return 0; fi
    for a in "fast-path.generated.go" "json.mono.generated.go"; do
        if [[ ! -e "$a" ]]; then return 0; fi
        for b in `ls -1 *.go.tmpl gen.go gen_mono.go values_test.go`; do
            if [[ "$a" -ot "$b" ]]; then return 0; fi
        done
    done
    return 1
}

# _build generates fast-path.go
_build() {
    # if ! [[ "${zforce}" || $(_ng "fast-path.generated.go") || $(_ng "json.mono.generated.go") ]]; then return 0; fi
    _build_proceed
    if [ $? -eq 1 ]; then return 0; fi
    if [ "${zbak}" ]; then
        _zts=`date '+%m%d%Y_%H%M%S'`
        _gg=".generated.go"
        [ -e "fast-path${_gg}" ] && mv fast-path${_gg} fast-path${_gg}__${_zts}.bak
        [ -e "gen${_gg}" ] && mv gen${_gg} gen${_gg}__${_zts}.bak
    fi 
    rm -f fast-path.generated.go *_generated_test.go gen-from-tmpl*.generated.go

    cat > gen-from-tmpl.codec.generated.go <<EOF
package codec
func GenTmplRun2Go(in, out string) { genTmplRun2Go(in, out) }
func GenMonoAll() { genMonoAll() }
EOF

    # explicitly return 0 if this passes, else return 1
    local btags="codec.gen codec.generics codec.safe codec.notfastpath"
    rm -f fast-path.generated.go mammoth_generated_test.go
    
    cat > gen-from-tmpl.generated.go <<EOF
//go:build ignore
package main
import "${zpkg}"
func main() {
codec.GenTmplRun2Go("fast-path.go.tmpl", "fast-path.generated.go")
codec.GenTmplRun2Go("mammoth-test.go.tmpl", "mammoth_generated_test.go")
}
EOF

    ${gocmd} run -tags "$btags" gen-from-tmpl.generated.go || return 1
    rm -f gen-from-tmpl.generated.go

    cat > gen-from-tmpl.generated.go <<EOF
//go:build ignore
package main
import "${zpkg}"
func main() {
codec.GenMonoAll()
}
EOF
    # btags="codec.safe codec.gen codec.generics"
    ${gocmd} run -tags "$btags" gen-from-tmpl.generated.go || return 1
    rm -f gen-from-tmpl*.generated.go
    return 0
}

_prebuild() {
    local d="$PWD"
    local zfin="test_values.generated.go"
    local zfin2="test_values_flex.generated.go"
    local zpkg="github.com/ugorji/go/codec"
    local returncode=1

    # zpkg=${d##*/src/}
    # zgobase=${d%%/src/*}
    # rm -f *_generated_test.go 
    true &&
        _build &&
        cp $d/values_test.go $d/$zfin &&
        cp $d/values_flex_test.go $d/$zfin2 &&
        if [[ "$(type -t _codegenerators_external )" = "function" ]]; then _codegenerators_external ; fi &&
        if [[ $zforce ]]; then ${gocmd} install ${zargs[*]} .; fi &&
        returncode=0 &&
        echo "prebuild done successfully"
    rm -f $d/$zfin $d/$zfin2
    return $returncode
    # unset zfin zfin2 zpkg
}

# _make() {
#     local makeforce=${zforce}
#     zforce=1
#     _prebuild && ${gocmd} install ${zargs[*]} .
#     zforce=${makeforce}
# }

_make() {
    _prebuild && ${gocmd} install ${zargs[*]} .
}

_clean() {
    rm -f \
       gen-from-tmpl.*generated.go \
       test_values.generated.go test_values_flex.generated.go
}

_tests() {
    local vet="" # TODO: make it off
    local gover=$( ${gocmd} version | cut -f 3 -d ' ' )
    [[ $( ${gocmd} version ) == *"gccgo"* ]] && zcover=0
    [[ $( ${gocmd} version ) == *"gollvm"* ]] && zcover=0
    case $gover in
        go1.2[0-9]*|go2.*|devel*) true ;;
        *) return 1
    esac
    # we test the following permutations wnich all execute different code paths as below.
    echo "TestCodecSuite: (fastpath/unsafe), (!fastpath/unsafe), (fastpath/!unsafe), (!fastpath/!unsafe)"
    local echo=1
    local nc=2 # count
    local cpus="1,$(nproc)"
    # if using the race detector, then set nc to
    if [[ " ${zargs[@]} " =~ "-race" ]]; then
        cpus="$(nproc)"
    fi
    local a=( "" "codec.safe" "codec.generics" "codec.generics codec.safe" "codec.generics codec.notfastpath" )
    local b=()
    local c=()
    for i in "${a[@]}"
    do
        local i2=${i:-default}
        [[ "$zwait" == "1" ]] && echo ">>>> TAGS: 'alltests $i'; RUN: 'TestCodecSuite'"
        [[ "$zcover" == "1" ]] && c=( -coverprofile "${i2// /-}.cov.out" )
        true &&
            ${gocmd} vet -printfuncs "errorf" "$@" &&
            if [[ "$echo" == 1 ]]; then set -o xtrace; fi &&
            ${gocmd} test ${zargs[*]} ${ztestargs[*]} -vet "$vet" -tags "alltests $i" -count $nc -cpu $cpus -run "TestCodecSuite" "${c[@]}" "$@" &
        if [[ "$echo" == 1 ]]; then set +o xtrace; fi
        b+=("${i2// /-}.cov.out")
        [[ "$zwait" == "1" ]] && wait
            
        # if [[ "$?" != 0 ]]; then return 1; fi
    done
    if [[ "$zextra" == "1" ]]; then
        [[ "$zwait" == "1" ]] && echo ">>>> TAGS: 'codec.generics codec.notfastpath x'; RUN: 'Test.*X$'"
        [[ "$zcover" == "1" ]] && c=( -coverprofile "x.cov.out" )
        ${gocmd} test ${zargs[*]} ${ztestargs[*]} -vet "$vet" -tags "codec.generics codec.notfastpath x" -count $nc -run 'Test.*X$' "${c[@]}" &
        b+=("x.cov.out")
        [[ "$zwait" == "1" ]] && wait
    fi
    wait
    # go tool cover is not supported for gccgo, gollvm, other non-standard go compilers
    [[ "$zcover" == "1" ]] &&
        command -v gocovmerge &&
        gocovmerge "${b[@]}" > __merge.cov.out &&
        ${gocmd} tool cover -html=__merge.cov.out
}

_release() {
    local reply
    read -p "Pre-release validation takes a few minutes and MUST be run from within GOPATH/src. Confirm y/n? " -n 1 -r reply
    echo
    if [[ ! $reply =~ ^[Yy]$ ]]; then return 1; fi

    # expects GOROOT, GOROOT_BOOTSTRAP to have been set.
    if [[ -z "${GOROOT// }" || -z "${GOROOT_BOOTSTRAP// }" ]]; then return 1; fi
    # (cd $GOROOT && git checkout -f master && git pull && git reset --hard)
    (cd $GOROOT && git pull)
    local f=`pwd`/make.release.out
    cat > $f <<EOF
========== `date` ===========
EOF
    # go 1.6 and below kept giving memory errors on Mac OS X during SDK build or go run execution,
    # that is fine, as we only explicitly test the last 3 releases and tip (2 years).
    local makeforce=${zforce}
    zforce=1
    for i in 1.10 1.11 1.12 master
    do
        echo "*********** $i ***********" >>$f
        if [[ "$i" != "master" ]]; then i="release-branch.go$i"; fi
        (false ||
             (echo "===== BUILDING GO SDK for branch: $i ... =====" &&
                  cd $GOROOT &&
                  git checkout -f $i && git reset --hard && git clean -f . &&
                  cd src && ./make.bash >>$f 2>&1 && sleep 1 ) ) &&
            echo "===== GO SDK BUILD DONE =====" &&
            _prebuild &&
            echo "===== PREBUILD DONE with exit: $? =====" &&
            _tests "$@"
        if [[ "$?" != 0 ]]; then return 1; fi
    done
    zforce=${makeforce}
    echo "++++++++ RELEASE TEST SUITES ALL PASSED ++++++++"
}

_usage() {
    # hidden args:
    # -pf [p=prebuild (f=force)]
    
    cat <<EOF
primary usage: $0 
    -t[esow]   -> t=tests [e=extra, s=short, o=cover, w=wait]
    -[md]      -> [m=make, d=race detector]
    -[n l i]   -> [n=inlining diagnostics, l=mid-stack inlining, i=check inlining for path (path)]
    -v         -> v=verbose
EOF
    if [[ "$(type -t _usage_run)" = "function" ]]; then _usage_run ; fi
}

_main() {
    if [[ -z "$1" ]]; then _usage; return 1; fi
    local x # determines the main action to run in this build
    local zforce # force
    local zcover # generate cover profile and show in browser when done
    local zwait # run tests in sequence, not parallel ie wait for one to finish before starting another
    local zextra # means run extra (python based tests, etc) during testing
    
    local ztestargs=()
    local zargs=()
    local zverbose=()
    local zbenchflags=""

    local gocmd=${MYGOCMD:-go}
    
    OPTIND=1
    while getopts ":cetmnrgpfvldsowkxyzi" flag
    do
        case "x$flag" in
            'xo') zcover=1 ;;
            'xe') zextra=1 ;;
            'xw') zwait=1 ;;
            'xf') zforce=1 ;;
            'xs') ztestargs+=("-short") ;;
            'xv') zverbose+=(1) ;;
            'xl') zargs+=("-gcflags"); zargs+=("-l=4") ;;
            'xn') zargs+=("-gcflags"); zargs+=("-m=2") ;;
            'xd') zargs+=("-race") ;;
            # 'xi') x='i'; zbenchflags=${OPTARG} ;;
            x\?) _usage; return 1 ;;
            *) x=$flag ;;
        esac
    done
    shift $((OPTIND-1))
    # echo ">>>> _main: extra args: $@"
    case "x$x" in
        'xt') _tests "$@" ;;
        'xm') _make "$@" ;;
        'xr') _release "$@" ;;
        'xg') _go ;;
        'xp') _prebuild "$@" ;;
        'xc') _clean "$@" ;;
        'xi') _check_inlining_one "$@" ;;
        'xk') _go_compiler_validation_suite ;;
        'xx') _analyze_checks "$@" ;;
        'xy') _analyze_debug_types "$@" ;;
        'xz') _analyze_do_inlining_and_more "$@" ;;
    esac
    # unset zforce zargs zbenchflags
}

[ "." = `dirname $0` ] && _main "$@"

