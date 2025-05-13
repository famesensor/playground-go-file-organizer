package main

import (
	"flag"
	"log"
	"strings"

	"github.com/famesensor/playground-go-file-organizer/organizer"
)

func main() {
	dir := flag.String("dir", ".", "Directory to organize")
	recursive := flag.Bool("recursive", false, "Include subdirectories")
	dryRun := flag.Bool("dry-run", false, "Preview changes without moving files")
	ignore := flag.String("ignore", "", "Ignore extensions (e.g. txt,pdf)")
	ctExtMap := flag.String("custom-map-ext", "", "Custom extension mapping (e.g. jpg,png=Images;pdf=Documents)")

	flag.Parse()

	ignoreExt := parseIgnoreExtension(*ignore)
	extMapping := parseCustomExtensionMapping(*ctExtMap)
	err := organizer.Organize(*dir, *recursive, *dryRun, ignoreExt, extMapping)
	if err != nil {
		log.Fatal(err)
	}
}

func parseIgnoreExtension(ignoreExt string) map[string]bool {
	ignoreMap := make(map[string]bool)
	for _, ext := range strings.Split(ignoreExt, ",") {
		ignoreMap[ext] = true
	}
	return ignoreMap
}

func parseCustomExtensionMapping(customExtMap string) map[string]string {
	ctExtMap := make(map[string]string)
	for _, pair := range strings.Split(customExtMap, ";") {
		parts := strings.Split(pair, "=")
		if len(parts) != 2 {
			continue
		}
		for _, ext := range strings.Split(parts[0], ",") {
			ctExtMap[ext] = parts[1]
		}
	}
	return ctExtMap
}
