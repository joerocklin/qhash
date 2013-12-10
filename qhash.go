/*
qhash.go - main body of the quick hashing tools
Copyright (C) 2013 Joe Rocklin <joe.rocklin@gmail.com>

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License along
with this program; if not, write to the Free Software Foundation, Inc.,
51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
*/

package main

import (
	"bufio"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"hash"
	"hash/crc32"
	"hash/crc64"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Checksum struct {
	name     string
	hashFunc hash.Hash
}

type Sumlist struct {
	filename string
	sums     []Checksum
}

var flag_help bool
var flag_all bool
var flag_crc32 bool
var flag_crc64 bool
var flag_sha224 bool
var flag_sha256 bool
var flag_sha384 bool
var flag_sha512 bool

var flag_num_procs int

var files []string

func init() {
	flag.BoolVar(&flag_help, "h", false, "Display this help message")
	flag.BoolVar(&flag_crc32, "crc32", false, "Compute the CRC32 hash")
	flag.BoolVar(&flag_crc64, "crc64", false, "Compute the CRC64 hash")
	flag.BoolVar(&flag_sha224, "sha224", false, "Compute the SHA224 hash")
	flag.BoolVar(&flag_sha256, "sha256", false, "Compute the SHA256 hash")
	flag.BoolVar(&flag_sha384, "sha384", false, "Compute the SHA384 hash")
	flag.BoolVar(&flag_sha512, "sha512", false, "Compute the SHA512 hash")
	flag.BoolVar(&flag_all, "all", false, "Compute all supported hashes")

	flag.IntVar(&flag_num_procs, "n", 1, "The number of prcesses to execute in parallel")

	flag.Parse()

	args := flag.Args()

	if len(args) > 0 {
		fileargs := args[0:]

		// Check the strings for any file globs in case the shell doesn't do this for us
		for _, file := range fileargs {
			glob_files, err := filepath.Glob(file)
			if err != nil {
				log.Fatal(err)
			}
			files = append(files, glob_files...)
		}
	}

	if flag_all {
		flag_crc32 = true
		flag_crc64 = true
		flag_sha224 = true
		flag_sha256 = true
		flag_sha384 = true
		flag_sha512 = true
	}

	if !(flag_crc32 || flag_crc64 ||
		flag_sha224 || flag_sha256 ||
		flag_sha384 || flag_sha512) {
		//log.Printf("Defaulting to SHA256\n")
		flag_sha256 = true
	}

	if len(files) == 0 {
		flag_help = true
		//log.Printf("ERROR: No files listed to hash\n")
	}

}

func main() {
	if flag_help {
		fmt.Println("")
		flag.Usage()
		return
	}

	// Set the number of processes to execute in parallel
	runtime.GOMAXPROCS(flag_num_procs)

	var complete_channel = make(chan Sumlist, len(files))

	for _, filename := range files {
		go process_file(filename, complete_channel)
	}

	for i := 1; i <= len(files); i++ {
		// Read a value
		sumlist := <-complete_channel

		// Print the sums
		for _, csum := range sumlist.sums {
			fmt.Printf("%x %s %s\n", csum.hashFunc.Sum(nil), csum.name, sumlist.filename)
		}
	}
}

func process_file(filename string, complete chan Sumlist) {
	sumlist := Sumlist{}
	sumlist.filename = filename

	// Open the file and bail if we fail
	infile, err := os.Open(filename)
	if err != nil {
		log.Printf("Unable to open %s: %s", filename, err)
		complete <- sumlist
		return
	}
	defer infile.Close()

	// Create the checksum objects
	if flag_crc32 {
		sumlist.sums = append(sumlist.sums, Checksum{"CRC32", crc32.New(crc32.IEEETable)})
	}
	if flag_crc64 {
		sumlist.sums = append(sumlist.sums, Checksum{"CRC64", crc64.New(crc64.MakeTable(crc64.ISO))})
	}
	if flag_sha224 {
		sumlist.sums = append(sumlist.sums, Checksum{"SHA224", sha256.New224()})
	}
	if flag_sha256 {
		sumlist.sums = append(sumlist.sums, Checksum{"SHA256", sha256.New()})
	}
	if flag_sha384 {
		sumlist.sums = append(sumlist.sums, Checksum{"SHA384", sha512.New384()})
	}
	if flag_sha512 {
		sumlist.sums = append(sumlist.sums, Checksum{"SHA512", sha512.New()})
	}

	// Create our file reader
	reader := bufio.NewReader(infile)

	// Start a buffer and loop to read the entire file
	buf := make([]byte, 4096)
	for {
		read_count, err := reader.Read(buf)
		// If we get an error that is not EOF, then we have a problem
		if err != nil && err != io.EOF {
			log.Printf("Unable to open %s: %s", filename, err)
			complete <- sumlist
			return
		}
		// If the returned size is zero, we're at the end of the file
		if read_count == 0 {
			break
		}

		// Add the buffer contents to the checksum calculation
		for _, sum := range sumlist.sums {
			sum.hashFunc.Write(buf[:read_count])
		}

	}

	complete <- sumlist
}
