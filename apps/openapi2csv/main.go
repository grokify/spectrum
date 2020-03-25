package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/grokify/gotilla/io/ioutilmore"
	"github.com/grokify/gotilla/os/osutil"
	"github.com/grokify/swaggman/csv"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Directory string `short:"d" long:"dir" description:"Source Directory" required:"true"`
	Regexp    string `short:"r" long:"regexp" description:"matching " required:"true"`
	Output    string `short:"o" long:"output" description:"Output CSV File" required:"true"`
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutilmore.DirEntriesReNotEmpty(
		opts.Directory, regexp.MustCompile(opts.Regexp))

	filepaths := osutil.FinfosToFilepaths(opts.Directory, files)

	tbl, err := csv.TableFromSpecFiles(filepaths, true)
	if err != nil {
		log.Fatal(fmt.Sprintf("TableFromSpecFiles [%v]\n", err.Error()))
	}

	err = tbl.WriteCSV(opts.Output)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WROTE [%v]\n", opts.Output)

	fmt.Println("DONE")
}
