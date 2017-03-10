package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"regexp"
	"sync"
	"syscall"
)

var (
	buildVersion string
	signals      chan os.Signal
	wg           sync.WaitGroup
)

func shutdown(tempdir string) {
	defer wg.Done()
	<-signals
	debugf("Removing tempdir %s\n", tempdir)
	err := os.RemoveAll(tempdir)
	if err != nil {
		fatal(err)
	}

}

func main() {
	var from, input, output, tempdir, tag string
	var keepTemp, version, last bool
	flag.StringVar(&input, "i", "", "Read from a tar archive file, instead of STDIN")
	flag.StringVar(&output, "o", "", "Write to a file, instead of STDOUT")
	flag.StringVar(&tag, "t", "", "Repository name and tag for new image")
	flag.StringVar(&from, "from", "", "Squash from layer ID (default: first FROM layer)")
	flag.BoolVar(&last, "last", false, "Squash from last found layer ID (Inverts order for automatic root-layer selection")
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

	var tagRegExp = regexp.MustCompile("^(.*):(.*)$")
	if tag != "" && strings.Contains(tag, ":") {
		parts := tagRegExp.FindStringSubmatch(tag)
		if parts[1] == "" || parts[2] == "" {
			fatalf("bad tag format: %s\n", tag)
		}
	}

	signals = make(chan os.Signal, 1)

	if !keepTemp {
		wg.Add(1)
		signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGTERM)
		go shutdown(tempdir)
	}

	export, err := LoadExport(input, tempdir)
	if err != nil {
		fatal(err)
	}

	// Export may have multiple branches with the same parent.
	// We can't handle that currently so abort.
	for _, v := range export.Repositories {
		commits := map[string]string{}
		for tag, commit := range *v {
			commits[commit] = tag
		}
		if len(commits) > 1 {
			fatal("This image is a full repository export w/ multiple images in it.  " +
				"You need to generate the export from a specific image ID or tag.")
		}

	}

	var start *ExportedImage
	if last {
		start = export.LastSquash()
		// Can't find a previously squashed layer, use last FROM
		if start == nil {
			start = export.LastFrom()
		}
	} else {
		start = export.FirstSquash()
		// Can't find a previously squashed layer, use first FROM
		if start == nil {
			start = export.FirstFrom()
		}
	}
	// Can't find a FROM, default to root
	if start == nil {
		start = export.Root()
	}

	if from != "" {

		if from == "root" {
			start = export.Root()
		} else {
			start, err = export.GetById(from)
			if err != nil {
				fatal(err)
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
			cmd := strings.Join(e.LayerConfig.ContainerConfig().Cmd, " ")
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
	err = newEntry.TarLayer()
	if err != nil {
		fatal(err)
	}

	debugf("Removing extracted layers\n")
	// remove our expanded "layer" dirs
	err = export.RemoveExtractedLayers()
	if err != nil {
		fatal(err)
	}

	if tag != "" {
		tagPart := "latest"
		repoPart := tag
		parts := tagRegExp.FindStringSubmatch(tag)
		if len(parts) > 2 {
			repoPart = parts[1]
			tagPart = parts[2]
		}
		tagInfo := TagInfo{}
		layer := export.LastChild()

		tagInfo[tagPart] = layer.LayerConfig.Id
		export.Repositories[repoPart] = &tagInfo

		debugf("Tagging %s as %s:%s\n", layer.LayerConfig.Id[0:12], repoPart, tagPart)
		err := export.WriteRepositoriesJson()
		if err != nil {
			fatal(err)
		}
	}

	ow := os.Stdout
	if output != "" {
		var err error
		ow, err = os.Create(output)
		if err != nil {
			fatal(err)
		}
		debugf("Tarring new image to %s\n", output)
	} else {
		debugf("Tarring new image to STDOUT\n")
	}
	// bundle up the new image
	err = export.TarLayers(ow)
	if err != nil {
		fatal(err)
	}

	debug("Done. New image created.")
	// print our new history
	export.PrintHistory()

	signals <- os.Interrupt
	wg.Wait()
}
