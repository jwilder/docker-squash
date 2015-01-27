package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/winchman/libsquash"
)

var (
	buildVersion string
)

func main() {
	//var from string
	var input, output, tag string
	var version bool
	flag.StringVar(&input, "i", "", "Read from a tar archive file, instead of STDIN")
	flag.StringVar(&output, "o", "", "Write to a file, instead of STDOUT")
	flag.StringVar(&tag, "t", "", "Repository name and tag for new image")
	//flag.StringVar(&from, "from", "", "Squash from layer ID (default: first FROM layer)") // TODO: should this be reintroduced
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.BoolVar(&version, "v", false, "Print version information and quit")

	flag.Usage = func() {
		fmt.Printf("\nUsage: docker-squash [options]\n\n")
		fmt.Printf("Squashes the layers of a tar archive on STDIN and streams it to STDOUT\n\n")
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if version {
		fmt.Println(buildVersion)
		return
	}

	if tag != "" && strings.Contains(tag, ":") {
		parts := strings.Split(tag, ":")
		if parts[0] == "" || parts[1] == "" {
			fatalf("bad tag format: %s\n", tag)
		}
	}

	var err error

	instream := os.Stdin
	if input != "" {
		instream, err = os.Open(input)
		if err != nil {
			fatal(err)
		}
		defer instream.Close()
	}

	outstream := os.Stdout
	if output != "" {
		outstream, err = os.Create(output)
		if err != nil {
			fatal(err)
		}
		defer outstream.Close()
		debugf("Tarring new image to %s\n", output)
	} else {
		debugf("Tarring new image to STDOUT\n")
	}

	imageIDBuffer := new(bytes.Buffer)

	libsquash.Verbose = verbose
	if err := libsquash.Squash(instream, outstream, imageIDBuffer); err != nil {
		fatal(err)
	}

	//TODO: add this functionality to libsquash
	//if tag != "" {
	//tagPart := "latest"
	//repoPart := tag
	//parts := strings.Split(tag, ":")
	//if len(parts) > 1 {
	//repoPart = parts[0]
	//tagPart = parts[1]
	//}
	//tagInfo := TagInfo{}
	//layer := export.LastChild()

	//tagInfo[tagPart] = layer.LayerConfig.Id
	//export.Repositories[repoPart] = &tagInfo

	//debugf("Tagging %s as %s:%s\n", layer.LayerConfig.Id[0:12], repoPart, tagPart)
	//err := export.WriteRepositoriesJson()
	//if err != nil {
	//fatal(err)
	//}
	//}
}
