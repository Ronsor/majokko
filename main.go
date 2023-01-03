// Copyright 2023 Ronsor Labs. All rights reserved.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/pborman/getopt/v2"
	"github.com/ronsor/majokko/henshin"
)

var hasError = false

var (
	doHelp = false
	doVersion = false
	doIdentify = false
	doConvert = false
	doListFormats = false
	maxWorkers int = 1
	noOutputFileNames = false

	identifyFormatString = "%wx%h, hash: %H, comment: %c"

	filterArgs FilterArgs
)

// FilterArgs is a set of filtering arguments.
// They are listed in the order they will be applied.
type FilterArgs struct {
	Strip bool
	AddComments []string
	SetComments []string
	Crop string
	Resize string
	CompressionLevel int
}

func init() {
	getopt.FlagLong(&doHelp, "help", 'h', "Show this help information")
	getopt.FlagLong(&doVersion, "version", 'v', "Show version information")
	getopt.FlagLong(&doIdentify, "identify", 0, "Print information about the image").SetGroup("action")
	getopt.FlagLong(&doConvert, "convert", 0, "Convert or process image (default)").SetGroup("action")
	getopt.FlagLong(&doListFormats, "list-formats", 0, "List supported image formats").SetGroup("action")
	getopt.FlagLong(&maxWorkers, "workers", 'W', "Maximum concurrent workers")
	getopt.FlagLong(&noOutputFileNames, "no-names", 'N', "Don't include file names in output messages")

	getopt.FlagLong(&identifyFormatString, "identify-format", 0, "Format string for --identify output")

	getopt.FlagLong(&filterArgs.Strip, "strip", 'S', "Strip metadata from image")
	getopt.FlagLong(&filterArgs.AddComments, "comment", 'C', "Add comment to image metadata")
	getopt.FlagLong(&filterArgs.SetComments, "set-comment", 0, "Set comments for image metadata")
	getopt.FlagLong(&filterArgs.Crop, "crop", 'c', "Crop image")
	getopt.FlagLong(&filterArgs.Resize, "resize", 'r', "Resize image")
	getopt.FlagLong(&filterArgs.CompressionLevel, "compress", 0, "Compression level, if applicable (0-100)")
	filterArgs.CompressionLevel = -1 // Set to default

	getopt.SetParameters("[images ...] [output path]")
}

func processFilterArgs(wand *henshin.Wand, fa *FilterArgs) {
	wand.ForceRGBA()

	if fa.Strip {
		wand.Strip()
	}

	if fa.AddComments != nil {
		for _, c := range fa.AddComments {
			wand.AddComment(c)
		}
	}

	if fa.SetComments != nil {
		wand.SetComments(fa.SetComments)
	}

	// TODO: crop

	if fa.Resize != "" {
		if fa.Resize[0] == '@' {
			area, err := strconv.ParseInt(fa.Resize[1:], 10, 32)
			if err == nil {
				wand.ResizeMaxArea(int(area), henshin.BiLinearStrategy)
			}
		} else {
			var w, h int
			n, err := fmt.Sscanf(fa.Resize, "%dx%d", &w, &h)
			if n == 2 && err == nil {
				wand.Resize(w, h, henshin.BiLinearStrategy)
			}
		}
	}

	if fa.CompressionLevel != -1 {
		wand.SetCompressionLevel(fa.CompressionLevel)
	}
}

func actionIdentify(wand *henshin.Wand, logPrefix string, inFile string) {
	fmt.Printf("%s%s\n", logPrefix, wand.FormatString(identifyFormatString))
}

func actionConvert(wand *henshin.Wand, logPrefix string, maxArg int, args []string, inFile string) {
	outFile := args[maxArg]
	if maxArg > 1 {
		outFile = filepath.Join(outFile, filepath.Base(inFile))
	}

	processFilterArgs(wand, &filterArgs)

	err := wand.WriteImage(outFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sWriteImage (to %s): %v\n", logPrefix, outFile, err)
		hasError = true
	}
}

func actionVersion() {
	fmt.Printf("Majokko %s (C) 2022-2023 Ronsor Labs. Licensed under the MIT license.\n", VERSION)
	fmt.Printf("Supported formats:")
	for _, c := range henshin.Codecs() {
		fmt.Printf(" %s", c.Name())
	}
	fmt.Printf("\nFor more information, use the --list-formats option.\n")
}

func main() {
	getopt.Parse()
	args := getopt.Args()

	if doHelp {
		getopt.Usage()
		return
	}

	if doVersion {
		actionVersion()
		if !doListFormats { return }
	}

	if doListFormats {
		fmt.Println("+---- Can decode?")
		fmt.Println("|+--- Can encode?")
		fmt.Println("||+-- Can handle metadata?")
		fmt.Println("|||   === Format name ===")
		for _, c := range henshin.Codecs() {
			_, canDecode := c.(henshin.Decoder)
			_, canEncode := c.(henshin.Encoder)
			name := c.Name()

			fmt.Printf("%s%s%s   %s\n",
				map[bool]string{true: "D", false: "-"}[canDecode],
				map[bool]string{true: "E", false: "-"}[canEncode],
				"?",
				name)
		}
		return
	}

	if !doConvert {
		doConvert = !doIdentify
	}

	maxArg := len(args)
	if doConvert { maxArg = maxArg - 1 }

	var wg sync.WaitGroup

	n := 0
	for i := 0; i < maxArg; i++ {
		wg.Add(1)
		n++

		go func(i int) {
			defer wg.Done()

			inFile := args[i]
			logPrefix := ""
			if !noOutputFileNames {
				logPrefix = inFile + ": "
			}

			wand := henshin.NewWand()
			err := wand.ReadImage(args[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%sReadImage: %v\n", logPrefix, err)
				hasError = true
				return
			}

			switch {
				case doIdentify: actionIdentify(wand, logPrefix, inFile)
				case doConvert: actionConvert(wand, logPrefix, maxArg, args, inFile)
			}
		} (i)

		if n == maxWorkers { wg.Wait(); n = 0 }
	}

	wg.Wait()
	if hasError {
		os.Exit(1)
	}
}
