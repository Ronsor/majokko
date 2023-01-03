// Copyright 2023 Ronsor Labs. All rights reserved.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	groupParamOpen, groupParamClose *int
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

	groupParamOpen = getopt.CounterLong("group", '(', "Open filter parameter group")
	groupParamClose = getopt.CounterLong("end-group", ')', "Close filter parameter group")

	initFilterArgs(&filterArgs, getopt.CommandLine)

	getopt.SetParameters("[images ...] [output path]")
}

func initFilterArgs(filterArgs *FilterArgs, optSet *getopt.Set) {
	optSet.FlagLong(&filterArgs.Strip, "strip", 'S', "Strip metadata from image")
	optSet.FlagLong(&filterArgs.AddComments, "comment", 'C', "Add comment to image metadata")
	optSet.FlagLong(&filterArgs.SetComments, "set-comment", 0, "Set comments for image metadata")
	optSet.FlagLong(&filterArgs.Crop, "crop", 'c', "Crop image")
	optSet.FlagLong(&filterArgs.Resize, "resize", 'r', "Resize image")
	optSet.FlagLong(&filterArgs.CompressionLevel, "compress", 0, "Compression level, if applicable (0-100)")
	filterArgs.CompressionLevel = -1 // Set to default
}

func parseFilterArgs(args []string) (ret []*FilterArgs) {
	groupStartIdx := -1
	groupEndIdx := -1
	groupCount := 1

	for i, arg := range args {
		if arg == "--group" || arg == "-(" {
			groupStartIdx = i
		} else if arg == "--end-group" || arg == "-)" {
			groupEndIdx = i
		}
		if groupStartIdx != -1 && groupEndIdx != -1 {
			if (groupEndIdx - groupStartIdx) == -1 {
				groupStartIdx = -1
				groupEndIdx = -1
				groupCount++
				continue
			}

			optSet := getopt.New()
			filterArgs := &FilterArgs{}
			initFilterArgs(filterArgs, optSet)

			section := args[groupStartIdx:groupEndIdx]
			err := optSet.Getopt(section, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Filter group %d: %v\n", groupCount, err)
				optSet.SetProgram("majokko -(")
				optSet.SetParameters("-)")
				optSet.PrintUsage(os.Stderr)
				os.Exit(1)
			}

			ret = append(ret, filterArgs)

			groupStartIdx = -1
			groupEndIdx = -1
			groupCount++
		}
	}

	return
}

func processFilterArgs(wand *henshin.Wand, fa *FilterArgs) {
	if fa.Strip {
		wand.Strip()
	}

	if fa.AddComments != nil {
		for _, c := range fa.AddComments {
			wand.AddComment(wand.FormatString(c))
		}
	}

	if fa.SetComments != nil {
		wand.SetComments(fa.SetComments)
	}

	if fa.Crop != "" {
		var (
			w = -1
			h = -1
			xoff = 0
			yoff = 0
		)
		n, _ := fmt.Sscanf(fa.Crop, "%dx%d+%d+%d", &w, &h, &xoff, &yoff)
		if n > 0 {
			wand.Crop(w, h, xoff, yoff)
		}
	}

	if fa.Resize != "" {
		resizeOpt := fa.Resize
		shrinkLarger := resizeOpt[len(resizeOpt)-1] == '>'
		enlargeSmaller := resizeOpt[len(resizeOpt)-1] == '<'
		if (shrinkLarger || enlargeSmaller) && len(resizeOpt) > 1 {
			resizeOpt = resizeOpt[:len(resizeOpt)-1]
		}

		percentScaling := strings.ContainsRune(resizeOpt, '%')
		if percentScaling {
			resizeOpt = strings.ReplaceAll(resizeOpt, "%", "")
			// Special case
			if strings.ContainsRune("0123456789", rune(resizeOpt[0])) && !strings.ContainsRune(resizeOpt, 'x') {
				// Ensure things like `-resize 50%` work
				resizeOpt = resizeOpt + "x"
			}
		}

		if resizeOpt[0] == '@' {
			var area int
			n, err := fmt.Sscanf(resizeOpt, "@%d", &area)
			if n == 1 && err == nil {
				currentArea := wand.Width() * wand.Height()
				if currentArea < area && shrinkLarger {
					goto SKIP
				} else if currentArea > area && enlargeSmaller {
					goto SKIP
				}

				wand.ResizeArea(int(area), henshin.BiLinearStrategy)
			}
		} else if resizeOpt[0] == 'x' {
			var h int
			n, err := fmt.Sscanf(resizeOpt, "x%d", &h)

			if n == 1 && err == nil {
				if percentScaling && h != 0 {
					h = int((float64(h) / 100) * float64(wand.Height()))
				}

				if wand.Height() < h && shrinkLarger {
					goto SKIP
				} else if wand.Height() > h && enlargeSmaller {
					goto SKIP
				}

				wand.Resize(-1, h, henshin.BiLinearStrategy)
			}
		} else if resizeOpt[len(resizeOpt)-1] == 'x' {
			var w int
			n, err := fmt.Sscanf(resizeOpt, "%dx", &w)

			if n == 1 && err == nil {
				if percentScaling && w != 0 {
					w = int((float64(w) / 100) * float64(wand.Width()))
				}

				if wand.Width() < w && shrinkLarger {
					goto SKIP
				} else if wand.Width() > w && enlargeSmaller {
					goto SKIP
				}

				wand.Resize(w, -1, henshin.BiLinearStrategy)
			}			
		} else {
			var w, h int
			n, err := fmt.Sscanf(resizeOpt, "%dx%d", &w, &h)
			if n == 2 && err == nil {
				if percentScaling && w != 0 {
					w = int((float64(w) / 100) * float64(wand.Width()))
				}
				if percentScaling && h != 0 {
					h = int((float64(h) / 100) * float64(wand.Height()))
				}

				area := w * h
				currentArea := wand.Width() * wand.Height()
				if currentArea < area && shrinkLarger {
					goto SKIP
				} else if currentArea > area && enlargeSmaller {
					goto SKIP
				}

				wand.Resize(w, h, henshin.BiLinearStrategy)
			}
		}

		SKIP:
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

	wand.ForceRGBA()

	faGroups := parseFilterArgs(os.Args)
	if faGroups == nil {
		processFilterArgs(wand, &filterArgs)
	} else {
		for _, fa := range faGroups {
			processFilterArgs(wand, fa)
		}
	}

	err := wand.WriteImage(outFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sWriteImage (to %s): %v\n", logPrefix, outFile, err)
		hasError = true
	}
}

func actionVersion(full bool) {
	fmt.Printf("Majokko %s (C) 2022-2023 Ronsor Labs. Licensed under the MIT license.\n", VERSION)
	fmt.Printf("Supported formats:")
	for _, c := range henshin.Codecs() {
		fmt.Printf(" %s", c.Name())
	}

	if !full { fmt.Println(""); return }
	fmt.Printf("\nFor more information, use the --list-formats option.\n")
}

func main() {
	getopt.Parse()
	args := getopt.Args()

	if doHelp {
		actionVersion(false)
		fmt.Println("")
		getopt.Usage()
		return
	}

	if doVersion {
		actionVersion(true)
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

	if *groupParamOpen != *groupParamClose {
		fmt.Fprintln(os.Stderr, "Unbalanced filter groups.")
		getopt.Usage()
		os.Exit(1)
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
