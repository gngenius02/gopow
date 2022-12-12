# POW

- implementing a pow algo using go.
- Based on live coding video by Tanmay Bakshi
    - [Live Coding: Optimizing a Proof of Work algorithm in Golang & C! (feat. Omer Kamal)](https://youtu.be/3zK_ogtwQw8)


## NOTE: I believe these improvements are only AMD specific. 
- The simd sha256 package does have ARM stuff in it but I am unable to test since I dont have an ARM cpu.
- The intel cpu's with AVX512 get slower with this package.


### I changed the following things and noticed some improvements.
- CPU: AMD Ryzen 9 5900X 12-core, 24-Thread
- Difficult: 31 bits
- added simd package
    - improved output from 70m/s to 160m/s
- removed loop at the end of POWOnCores function
    - this seemed to give a speed up from 160m/s to over 200m/s


```shell
~/work/pow > go run .
    abcc628cac49b36d7562fc6
    time: 4.95525255s
    Processed 1,074,202,624
    Processed/sec 216,780,600

~/work/pow > go run .
    abcab3ef3ebe7b8512b3d9a
    time: 7.500131134s
    Processed 1,604,866,048
    Processed/sec 213,978,398

~/work/pow > go run .
    abc35ef30da2a988533f712
    time: 2.379043285s
    Processed 511,780,864
    Processed/sec 215,120,450

```


## ERRORS?

- I noticed that somtimes (about 20-30% of the time it gives an incorrect solution)

```shell
~/work/pow > go run .
    abcaaa8c05f2fe483ccfe44
    time: 226.145284ms
    Processed 49,139,712
    Processed/sec 217,292,667

~/work/pow > echo -n abcaaa8c05f2fe483ccfe44 | sha256sum
94a10fcdc7d226b74f50b1cc6bc92536601a4099ac21d3220ecaf02fa8c9a4a6  -
```