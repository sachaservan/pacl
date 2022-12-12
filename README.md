# Private Access Control Lists (PACLs)

Implementation of PACLs for DPFs and VDPFs.

| **Code organization** ||
| :--- | :---|
| Implementation||
| [pacl-pk/](pacl-pk/) | Implementation of the public-key (V)DPF-PACL construction|
| [pacl-sk/](pacl-sk/) | Implementation of the secret-key (V)DPF-PACL construction|
| [sposs/](sposs/) | Implementation of the Schnorr Proof over Secret Shares (SPoSS)|
| [algebra/](algebra/) | Bare-bones implementation of fields and groups|
| [ec/](ec/) | A wrapper for the P256 elliptic curve|
| Evaluation and results||
| [bench-fss/](bench-fss/) | DPF-PACLs and VDPF-PACLs benchmarks|
| [bench-anon/](bench-anon/) | Anonymous communication benchmarks using VDPF-PACLs|
| [bench-auth/](bench-auth/) | Anonymous authentication benchmarks using VDPF-PACLs|
| [bench-pir/](bench-pir/) | PIR benchmarks using VDPF-PACLs|
| [paper_results/](paper_results/) | Raw evaluation data (.json) used in the paper |



## Dependencies 
* GMP 
* Go 1.13 or higher 
* OpenSSL 1.1.1f
* GNU Make
* Cmake

## Getting everything to run (tested on Ubuntu, CentOS, and MacOS)

|Dependency |Install dependencies (Ubuntu): | Install dependencies (CentOS):|
|--------------|--------------|-----------|
|GMP library |```sudo apt-get install libgmp3-dev```| ```sudo yum install gmp-devel```|
|Go |```sudo apt-get install golang-go```| ```sudo yum install golang```|
|OpenSSL |```sudo apt install libssl-dev```|```sudo yum install openssl-devel```|
|Make |```sudo apt-get install build-essential``` |  ```sudo yum groupinstall 'Development Tools'```|
|Cmake |```sudo apt-get install cmake```| ```sudo yum install cmake```|

For optimal performance, you should compile the C code with clang (approximately 10-20 percent faster than the default on some distributions).

|Dependency |Install dependencies (Ubuntu): | Install dependencies (CentOS):|
|--------------|--------------|-----------|
|Clang-11 |```sudo apt install clang```| ```sudo yum install clang```|
|LLVM-AR |```sudo apt install llvm```| ```sudo yum install llvm```|


### 0) initialize the VDPF submodule 
```
git submodule update --init --recursive
``` 

### 1) Compiling the C VDPF/DPF library
```
go mod tidy
cd vdpf/src && make
```

### 2) Running the benchmarks

| **Benchmarks** ||
| :--- | :---|
| FSS benchmarks | ```cd bench-fss && bash run.sh```|
| Spectrum & Express | ```cd bench-anon && bash run.sh```|
| Anonymous authentication | ```cd bench-auth && bash run.sh```|
| Private Information Retrieval | ```cd bench-pir && bash run.sh```|

### 3) Plotting! 

Raw JSON data and plotting scripts are located in [paper_results/](paper_results/).
```
cd paper_results/
python plot_fss.py --file fss.json
python plot_vfss.py --file fss.json
python plot_pir.py --file pir.json
python plot_anon.py --file anon.json
```


## ⚠️ Important Warning
<b>This implementation of is intended for *research purposes only*. The code has NOT been vetted by security experts. 
As such, no portion of the code should be used in any real-world or production setting!</b>

## License
Copyright © 2022 Sacha Servan-Schreiber

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
