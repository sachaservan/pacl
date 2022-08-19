#!/bin/bash

go build -o bench-auth
./bench-auth
rm bench-auth