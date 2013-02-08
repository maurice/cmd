// Copyright 2013 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// command empties finds empty directories
package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

var ignored = []string{}

var dirNames = []string{}
var fileCounts = map[string]int{}
var ignoreHits = map[string][]string{}

func visitDir(path string, info os.FileInfo, err error) error {
	dir := filepath.Dir(path)

	var ignore bool = false
	for _, name := range ignored {
		if base := filepath.Base(path); base == name {
			ignoreHits[dir] = append(ignoreHits[dir], name)
			ignore = true
			break
		}
	}

	if info.IsDir() {
		if ignore {
			// don't descend into ignored dirs
			return filepath.SkipDir
		}
		dirNames = append(dirNames, path)
		fileCounts[path] = 0
	}

	// don't count ignored files
	if !ignore {
		fileCounts[dir] = fileCounts[dir] + 1
	}

	return nil
}

func main() {
	log.SetPrefix("")
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <dir> [ignored]\n", os.Args[0])
	}

	if len(os.Args) == 3 {
		ignored = strings.Split(os.Args[2], ",")
	}

	dir := os.Args[1]
	if len(ignored) > 0 {
		log.Printf("Scanning '%s', ignoring %v...\n", dir, strings.Join(ignored, ", "))
	} else {
		log.Printf("Scanning '%s'...\n", dir)
	}

	err := filepath.Walk(dir, visitDir)
	if err != nil {
		log.Fatalf("Failed to stat dir: %v\n", err)
	}

	var numEmpty = 0
	for _, name := range dirNames {
		if count, ok := fileCounts[name]; ok && count == 0 {
			numEmpty += 1
			if ignored := ignoreHits[name]; ignored != nil {
				log.Printf("%s (only %s)", name, strings.Join(ignored, ", "))
			} else {
				log.Printf("%s\n", name)
			}
		}
	}
	log.Printf("%d total, %d empty\n", len(dirNames), numEmpty)
}
