/*
Author: Aggelos Kolaitis
Description: Photos and videos organizer.
Version: 0.1
Usage: ./gorganizer -source INPUT -dest OUTPUT
*/

package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ParsedFile contains the md5 sum and date information for the file
type ParsedFile struct {
	md5sum   string
	filename string
	fileinfo os.FileInfo
}

type truthTable = map[string]bool

// getHash() calculates the md5 sum of a file
func getHash(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(h.Sum(nil))
}

// NewParsedFile parses a file from its name
func NewParsedFile(filename string) *ParsedFile {
	stat, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}

	return &ParsedFile{
		filename: filename,
		md5sum:   getHash(filename),
		fileinfo: stat,
	}
}

// GetName returns the file destination name
func (f ParsedFile) GetName(id int, prefix string) string {
	fix := strings.NewReplacer("/", "-", " ", "-").Replace(f.filename)
	modTime := f.fileinfo.ModTime()
	return fmt.Sprintf("%s/%d/%02d-%s/%05d-%s",
		prefix, modTime.Year(), modTime.Month(), modTime.Month(), id, fix)
}

// CopyTo copies file to destination directory
func (f ParsedFile) CopyTo(id int, destination string) error {
	// open file for reading
	src, err := os.Open(f.filename)
	if err != nil {
		log.Printf("[ERR] Could not open %s for reading!\n", f.filename)
		return err
	}
	defer src.Close()

	// create directory if needed
	destName := f.GetName(id, destination)
	destDir := destName[:strings.LastIndex(destName, "/")]
	err = os.MkdirAll(destDir, 0775)
	if err != nil && !os.IsExist(err) {
		log.Printf("[ERR] Could not create destination directory %s: %s!\n", destDir, err)
	}

	// open file for writing
	dest, err := os.Create(destName)
	if err != nil {
		log.Printf("[ERR] Could not open %s for writing!\n", destName)
		return err
	}
	defer dest.Close()

	// copy contents
	_, err = io.Copy(dest, src)
	if err != nil {
		return err
	}

	// preserve modification time
	err = os.Chtimes(destName, time.Now(), f.fileinfo.ModTime())
	if err != nil {
		log.Printf("[ERR] Could not update times for %s!\n", destName)
	}
	return err
}

func main() {
	// parse command-line arguments
	source := flag.String("source", "", "Specify source directory")
	dest := flag.String("dest", "", "Specify destination directory")
	flag.Parse()
	if *source == "" {
		log.Fatal("[ERR] Missing required flag -source")
	}
	if *dest == "" {
		log.Fatal("[ERR] Missing required flag -dest")
	}

	hashes := make(truthTable)
	var found []ParsedFile

	log.Println("[INF] Parsing existing files")
	filepath.Walk(*dest, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		hashes[NewParsedFile(path).md5sum] = true
		return nil
	})

	log.Println("[INF] Parsing source files")
	filepath.Walk(*source, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		log.Printf("[INF] Found: %s\n", path)
		f := NewParsedFile(path)
		_, exists := hashes[f.md5sum]
		if !exists {
			hashes[f.md5sum] = true
			found = append(found, *f)
		}
		return nil
	})

	sort.Slice(found, func(h1 int, h2 int) bool {
		return found[h1].fileinfo.ModTime().Before(found[h2].fileinfo.ModTime())
	})

	count := make(map[int]int)
	for _, file := range found {
		index := 13*file.fileinfo.ModTime().Year() + int(file.fileinfo.ModTime().Month())
		count[index]++
		file.CopyTo(count[index], *dest)
	}
}
