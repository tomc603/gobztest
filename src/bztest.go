package main

import (
	"bufio"
	"compress/bzip2"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
)

//const READ_BUF_SIZE = 10485760 // Set 10MB read buffer
var cpuprofile *string = flag.String("cpuprofile", "", "Write cpu profile to file")
var logdir *string = flag.String("logdir", "", "Log directory")

func main() {
	flag.Parse()

	// Check arguments for sanity and requirements
	switch {
	case *logdir == "":
		log.Fatal("You must provide a log source")
	case *cpuprofile != "":
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	fileCollect(logdir)
}

func scan(path string, f os.FileInfo, err error) error {
	if f != nil {
		if !f.IsDir() {
			log.Printf("Processing %s\n", path)
			e := fileOpen(path)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

// Recurse through path p and collect files of interest.
func fileCollect(p *string) error {
	// Walk logdir and call scan function for each file
	err := filepath.Walk(*p, scan)
	if err != nil {
		log.Printf("ERROR: fileCollect: %v\n", err)
		return err
	}
	return nil
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
	// case ".gz":
	// 	// Read bzip2 compressed log file
	// 	cr, _ := gzip.NewReader(br)
	// 	defer cr.Close()
	// 	err := dataParse(cr)
	// 	if err != nil {
	// 		return err
	// 	}
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

	// s := bufio.NewScanner(br)
	// for s.Scan() {
	//  _ = s.Text()
	// }

	// err := s.Err()
	// if err != nil {
	//  return err
	// }

	_, e := io.Copy(ioutil.Discard, br)
	if e != nil {
		return e
	}

	return nil
}
