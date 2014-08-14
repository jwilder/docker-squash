package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var buildVersion string

func main() {
	var from, input, output, tempdir, tag string
	var keepTemp, version bool
	flag.StringVar(&input, "i", "", "Read from a tar archive file, instead of STDIN")
	flag.StringVar(&output, "o", "", "Write to a file, instead of STDOUT")
	flag.StringVar(&tag, "t", "", "Repository name and tag for new image")
	flag.StringVar(&from, "from", "", "Squash from layer ID (default is root)")
	flag.BoolVar(&keepTemp, "keepTemp", false, "Keep temp dir when done. (Useful for debugging)")
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

	var err error
	tempdir, err = ioutil.TempDir("", "docker-squash")
	if err != nil {
		fatal(err)
	}

	if tag != "" {
		if !strings.Contains(tag, ":") {
			fatalf("bad tag format: %s\n", tag)
			return
		}
		parts := strings.Split(tag, ":")
		if parts[0] == "" || parts[1] == "" {
			fatalf("bad tag format: %s\n", tag)
			return
		}
	}

	if !keepTemp {
		defer func() {
			debugf("Removing tempdir %s\n", tempdir)
			os.RemoveAll(tempdir)
		}()
	}

	export, err := LoadExport(input, tempdir)
	if err != nil {
		fatal(err)
		return
	}

	if len(export.Repositories) > 0 {
		fatal("This image is a full repository export w/ multiple images in it.  " +
			"You need to generate the export from a specific image ID or tag.")
		return
	}

	start := export.FirstFrom()
	if from != "" {

		if from == "root" {
			start = export.Root()
		} else {
			start, err = export.GetById(from)
			if err != nil {
				fatal(err)
				return
			}
		}
	}

	if start == nil {
		fatalf("no layer matching %s\n", from)
		return
	}

	// extract each "layer.tar" to "layer" dir
	err = export.ExtractLayers()
	if err != nil {
		fatal(err)
		return
	}

	// insert a new layer after our squash point
	newEntry, err := export.InsertLayer(start.LayerConfig.Id)
	if err != nil {
		fatal(err)
		return
	}

	debugf("Inserted new layer %s after %s\n", newEntry.LayerConfig.Id[0:12],
		newEntry.LayerConfig.Parent[0:12])

	if verbose {
		e := export.Root()
		for {
			if e == nil {
				break
			}
			cmd := strings.Join(e.LayerConfig.ContainerConfig.Cmd, " ")
			if len(cmd) > 60 {
				cmd = cmd[:60]
			}

			if e.LayerConfig.Id == newEntry.LayerConfig.Id {
				debugf("  -> %s %s\n", e.LayerConfig.Id[0:12], cmd)
			} else {
				debugf("  -  %s %s\n", e.LayerConfig.Id[0:12], cmd)
			}
			e = export.ChildOf(e.LayerConfig.Id)
		}
	}

	// squash all later layers into our new layer
	err = export.SquashLayers(newEntry, newEntry)
	if err != nil {
		fatal(err)
		return
	}

	debugf("Tarring up squashed layer %s\n", newEntry.LayerConfig.Id[:12])
	// create a layer.tar from our squashed layer
	newEntry.TarLayer()

	debugf("Removing extracted layers\n")
	// remove our expanded "layer" dirs
	export.RemoveExtractedLayers()

	if tag != "" {
		parts := strings.Split(tag, ":")
		tagInfo := TagInfo{}
		layer := export.LastChild()
		tagInfo[parts[1]] = layer.LayerConfig.Id
		export.Repositories[parts[0]] = &tagInfo

		debugf("Tagging %s as %s\n", layer.LayerConfig.Id[0:12], tag)
		err := export.WriteRepositoriesJson()
		if err != nil {
			fatal(err)
			return
		}
	}

	ow := os.Stdout
	if output != "" {
		var err error
		ow, err = os.Create(output)
		if err != nil {
			fatal(err)
			return
		}
		debugf("Tarring new image to %s\n", output)
	} else {
		debugf("Tarring new image to STDOUT\n")
	}
	// bundle up the new image
	export.TarLayers(ow)

	debug("Done. New image created.")

	// print our new history
	export.PrintHistory()
}
