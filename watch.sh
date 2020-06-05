#!/bin/bash

filewatcher 'src/**/*' 'printf "\ec" && ./test.sh'
#filewatcher '**/*' 'printf "\ec" && cd test && go run ../src/'
