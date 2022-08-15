#!/bin/bash
go build -o bench-fss
./bench-fss --numeval $NUMEVAL
rm bench-fss