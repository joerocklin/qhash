/*
qhash_test.go - test code for the qhash code
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
	"bytes"
	"encoding/hex"
	"testing"
)

const small_test_file = "./test_data/small_test_file"

var small_sum_sha224, _ = hex.DecodeString("3a167ed2f5013762c019ee4f2ee20a72270c51d8ecd26feea43166c2")
var small_sum_sha256, _ = hex.DecodeString("20fe9d7e663079224fcedbe70e14d18345de592a3a6e84131695c574df1e0cf6")

func TestSha256(t *testing.T) {
	var complete_channel = make(chan Sumlist, 1)
	flag_sha256 = true
	go process_file(small_test_file, complete_channel)

	sumlist := <-complete_channel
	if len(sumlist.sums) != 1 {
		t.Errorf("Expected 1 checksum, received %d", len(sumlist.sums))
	}

	csum := sumlist.sums[0].hashFunc.Sum(nil)
	if !bytes.Equal(csum, small_sum_sha256) {
		t.Errorf("Wrong checksum received:\n  %x <- Expected\n  %x <- Received", small_sum_sha256, csum)
	}
}

func TestMultipleHashes(t *testing.T) {
	var complete_channel = make(chan Sumlist, 1)

	flag_sha224 = true
	flag_sha256 = true

	go process_file(small_test_file, complete_channel)

	sumlist := <-complete_channel
	if len(sumlist.sums) != 2 {
		t.Errorf("Expected 2 checksum, received %d", len(sumlist.sums))
	}

	csum := sumlist.sums[0].hashFunc.Sum(nil)
	if !bytes.Equal(csum, small_sum_sha224) {
		t.Errorf("Wrong checksum received:\n  %x <- Expected\n  %x <- Received", small_sum_sha224, csum)
	}

	csum = sumlist.sums[1].hashFunc.Sum(nil)
	if !bytes.Equal(csum, small_sum_sha256) {
		t.Errorf("Wrong checksum received:\n  %x <- Expected\n  %x <- Received", small_sum_sha256, csum)
	}
}
