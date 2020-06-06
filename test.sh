#!/bin/bash

printf "\ec"
#./build.sh && cd test && time ../taggenator open <<< "-some other tag, yolo baggins-"
#./build.sh && cd test && time ../taggenator open -sort unseen
#./build.sh && cd test && time ../taggenator open
#cp ./test/data.db.old ./test/data.db

#./build.sh && cd test && time ../taggenator <<< ""
#./build.sh && cd test && time ../taggenator dump_tags 
cp ./test/data.db.old ./test/data.db
#./build.sh && cd test && time ../taggenator open <<< "test"
./build.sh && cd test && time ../taggenator apply_tags --tag --threads 20 silly <<< "test"
#./build.sh && cd test && time ../taggenator apply_tags --tag silly <<< "test"
#./build.sh && cd test && time ../taggenator apply_tags --tag silly <<< "test"
#./build.sh && cd test && time ../taggenator dump
#./build.sh && cd test && time ../taggenator dump search 7005 <<< ""

#./build.sh && cd test && time ../taggenator dump_tags
#./build.sh && cd test && time ../taggenator open search oopsie <<< ""

cd ..
cp ./taggenator ~/Downloads/Shared/
#./build.sh && cd test && time ../taggenator open_read_only
#./build.sh && cd test && time ../taggenator open search "temp" <<< ""
#./build.sh && cd test && time ../taggenator move -destination "/tmp" -sort search 004 <<< "ls"
#./build.sh && cd test && time ../taggenator open <<< "pre"
#./build.sh && cd test && time ../taggenator help
#./build.sh && cd test && time ../taggenator open <<< "y"
#./build.sh && cd test && time ../taggenator open <<< "y"
#./build.sh && cd test && time ../taggenator open
#./build.sh && cd test && time ../taggenator open <<< ""
#./build.sh && cd test && time ../taggenator open <<< "some tag, some other tag, some third tag, some fourth tag, some fifth tag, some sixth tax"
#./build.sh && cd test && time ../taggenator open <<< "< 5"
#./build.sh && cd test && time ../taggenator open <<< "rename:applesauce.txt"
#./build.sh && cd test && time ../taggenator open <<< "rename:AppleSauce.txt"
#./build.sh && cd test && time ../taggenator open -sort unseen -sort reverse <<< "rename:AppleSauce.txt"
#./build.sh && cd test && time ../taggenator open -sort unseen <<< "a"
#./build.sh && cd test && time ../taggenator dump_tags
#./build.sh && cd test && time ../taggenator help
#./build.sh && cd test && time ../taggenator open -sort search tagA tagB

# Multiple process test (failed)
#go build -o taggenator ./src/ && cd test
#../taggenator A &
#../taggenator B &
