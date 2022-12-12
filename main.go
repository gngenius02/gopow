package main

import (
	// "crypto/sha256"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/minio/sha256-simd"
)

var charSet = []byte("abcdef0123456789")

// var charSet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

var hashes = []int{}

func Comma(v int64) string {
	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1
	for v > 999 {
		parts[j] = strconv.FormatInt(v%1000, 10)
		switch len(parts[j]) {
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		}
		v = v / 1000
		j--
	}
	parts[j] = strconv.Itoa(int(v))
	return strings.Join(parts[j:], ",")
}

func main() {
	start := time.Now()
	prefix := []byte("abc")
	solutionChan := make(chan []byte, 1)

	for i := 0; i < runtime.NumCPU(); i++ {
		POWOnCores(prefix, solutionChan)
	}
	solution := <-solutionChan
	fmt.Println(string(solution))
	totalHashesProcessed := 0

	for _, value := range hashes {
		totalHashesProcessed += value
	}

	end := time.Since(start)

	fmt.Printf("time: %s\n", end)
	fmt.Printf("Processed %s\n", Comma(int64(totalHashesProcessed)))
	fmt.Printf("Processed/sec %s\n", Comma(int64(float64(totalHashesProcessed)/end.Seconds())))
}

func POWOnCores(prefix []byte, solutionChannel chan []byte) {
	blockSize := 1024

	unprocessIndex := make(chan int, 2)
	processIndex := make(chan int, 2)
	offset := len(prefix)

	blocks := make([][][]byte, 2)

	for idx := range blocks {
		unprocessIndex <- idx
		blocks[idx] = make([][]byte, blockSize)
		for i := 0; i < blockSize; i++ {
			blocks[idx][i] = make([]byte, 20)
			blocks[idx][i] = append(prefix, blocks[idx][i]...)
		}
	}
	go func() {
		seed := uint64(time.Now().Local().UnixNano())
		for {
			blockIndex := <-unprocessIndex
			for idx := range blocks[blockIndex] {
				seed = RandStr(blocks[blockIndex][idx], offset, seed)
			}
			processIndex <- blockIndex
		}
	}()

	idx := len(hashes)
	hashes = append(hashes, 0)
	go func(hashIndex int) {
		var hash bool
		for {
			blockIndex := <-processIndex
			hashes[hashIndex] += blockSize
			for _, random := range blocks[blockIndex] {
				hash = Hash(random, 31)
				if hash {
					solutionChannel <- random
					break
				}
			}
			unprocessIndex <- blockIndex
		}
	}(idx)
}

func RandNum(seed uint64) uint64 {
	seed ^= seed << 21
	seed ^= seed >> 35
	seed ^= seed << 4
	return seed
}

func RandStr(str []byte, offset int, seed uint64) uint64 {
	for i := offset; i < len(str); i++ {
		seed = RandNum(seed)
		str[i] = charSet[seed%16]
	}
	return seed
}

func Hash(data []byte, bits int) bool {
	bs := sha256.Sum256(data)
	nbytes := bits / 8
	nbits := bits % 8
	idx := 0
	for idx < nbytes {
		if bs[idx] > 0 {
			return false
		}
		idx++
	}
	return (bs[idx] >> (8 - nbits)) == 0
}
