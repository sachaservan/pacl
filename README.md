# Verifiable Distributed Point Functions

Implementation of VDPFs in C with a Go wrapper. See the [paper](https://eprint.iacr.org/2021/580.pdf) for details.

## Dependencies 
* Go 1.13 or higher 
* OpenSSL 1.1.1f
* GNU Make
* Cmake

## Getting everything to run (tested on Ubuntu, CentOS, and MacOS)

|Install dependencies (Ubuntu): | Install dependencies (CentOS):|
|--------------|-----------|
|```sudo apt-get install build-essential``` |  ```sudo yum groupinstall 'Development Tools'```|
|```sudo apt-get install cmake```| ```sudo yum install cmake```|
|```sudo apt install libssl-dev```|```sudo yum install openssl-devel```|
|```sudo apt-get install golang-go```| ```sudo yum install golang```|


For optimal performance, you should compile the C code with clang (approximately 10-20 percent faster than the default on some distributions).
- Clang-11: On Ubuntu run ```sudo apt install clang```.  On CentOS, ```sudo yum install clang```.
  - You'll also need llvm if you use clang. 
- LLVM-AR: On Ubuntu run ```sudo apt install llvm```. On CentOS, ```sudo yum install llvm```.

### 1) Compiling the C VDPF/DPF library
```
go mod tidy
cd src && make
```

### 2a) Running tests (in C) 
```
./test
```
See also [```src/test.c```](src/test.c)

### 2b) Running tests (in Go)
```
go test
```
See also [```dpf_test.go```](dpf_test.go)

## ⚠️ Important Warning
<b>This implementation of is intended for *research purposes only*. The code has NOT been vetted by security experts. 
As such, no portion of the code should be used in any real-world or production setting!</b>

## License
Copyright © 2022 Sacha Servan-Schreiber and Simon Langowski

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
# pacl
