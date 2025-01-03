#!/bin/sh

go run . -pk /data/ec/w.json -t /data/ec/templates.txt  -threads 10  -keep_order -min_len 1 -max_len 20 -v 1  -re 10  
#go run src/ethcracker.go -pk ~/test/ethwallet-q.json -t ~/test/pattern.txt -presale -threads 4  -min_len 1 -v 1 -start_from 0
# -dump ~/test/v.txt 