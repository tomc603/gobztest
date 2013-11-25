package main

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"
)

//const READ_BUF_SIZE = 10485760 // Set 10MB read buffer
var cpuprofile *string = flag.String("cpuprofile", "", "Write cpu profile to file")
var compfile *string = flag.String("file", "", "Compressed file to decompress")

func main() {
	flag.Parse()

	// Check arguments for sanity and requirements
	switch {
	case *compfile == "":
		log.Fatal("You must provide a compressed file!\n")
	case *cpuprofile != "":
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	log.Printf("Decompressing %s\n", *compfile)
	err := fileOpen(*compfile)
	if err != nil {
		log.Fatalf(" * ERROR: %v\n", err)
	}
	log.Print("Done\n")
}

func fileOpen(path string) error {
	// Open the file and defer its close
	fh, err := os.OpenFile(path, 0, 0)
	if err != nil {
		return err
	}
	defer fh.Close()

	br := bufio.NewReader(fh)

	switch filepath.Ext(path) {
	case ".bz2":
		// Read bzip2 compressed log file
		cr := bzip2.NewReader(br)
		//err := logRead(cr)
		err := dataParse(cr)
		if err != nil {
			return err
		}
	case ".gz":
		// Read gzip compressed log file
		cr, _ := gzip.NewReader(br)
		defer cr.Close()
		err := dataParse(cr)
		if err != nil {
			return err
		}
	default:
		// Read uncompressed log file
		err := dataParse(br)
		if err != nil {
			return err
		}
	}
	return nil
}

func dataParse(r io.Reader) error {
	br := bufio.NewReader(r)

	st := time.Now()
	defer log.Printf(" * Decompress time: %0.2f sec", time.Since(st).Seconds())
	_, e := io.Copy(ioutil.Discard, br)
	if e != nil {
		return e
	}

	return nil
}
