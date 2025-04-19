#!/bin/bash
set -e

GODSLSRC="../go-dsl/app"

cleanup() { # $1 = example name
    rm -rf example-$1 
    rm -rf $1/dsl_* 
    rm -rf $1/template_* 
}

build() { # $1 = example name
    go build -o ./example-$1 ./$1/ || {
        printRed "Failed to build $1"
        exit 1
    }
}

build-go-dsl() { # $1 = source dir
    go build -C $1 -o $(pwd)/go-dsl . || {
        printRed "Failed to build GoDSL"
        exit 1
    }
}

genDSL() { # $1 = example name
    ./go-dsl $1 "$1 script" "An example implementation of $1 script" "1.0.0" "$1" ./$1/
    go mod tidy # to ensure we have all dependencies
}

run() { # $1 = example name
    ./example-$@ || {
        printRed "Failed to execute example"
        exit 1
    }
}

printBlue() { # $1 = message
    echo -e "\033[1;34m$1\033[0m"
}

printGreen() { # $1 = message
    echo -e "\033[1;32m$1\033[0m"
}   

printRed() { # $1 = message
    echo -e "\033[1;31m$1\033[0m"
}   

#########################################
printBlue "Preparing build..."
#########################################
go mod tidy
cleanup "basic"
cleanup "calculator"
cleanup "ohms-law"
cleanup "machine-intel" 

#########################################
printBlue "Building GoDSL..."
#########################################
build-go-dsl $GODSLSRC -o ./go-dsl 

#########################################
printBlue "Executing Basic example..."
#########################################
genDSL "basic" 
build "basic"
run "basic" 'l:+(test-function-1(1 2 "hi, this will be printed") *(50 2)) +(l +(1 gx))'

#########################################
printBlue "Executing Calculator example..."
#########################################
genDSL "calculator" 
build "calculator"
run "calculator"

#########################################
printBlue "Executing Ohm's law example..."
#########################################
genDSL "ohms-law" 
build "ohms-law"
run "ohms-law" 

#########################################
printBlue "Executing Machine Intel example..."
#########################################
genDSL "machine-intel" 
build "machine-intel"
run "machine-intel" 

#########################################
printBlue "Executing Image Filter example..."
#########################################
genDSL "image-filter" 
build "image-filter"
run "image-filter" 

#########################################
printBlue "Cleaning up..."
#########################################
cleanup "basic"
cleanup "calculator"
cleanup "ohms-law"
cleanup "machine-intel" 
cleanup "image-filter" 
rm -rf go-dsl
rm -rf doc.*

printGreen "All examples executed successfully"
exit 0
