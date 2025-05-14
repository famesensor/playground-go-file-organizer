package organizer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Organize(root string, recursive, dryRun bool, ignoreExt map[string]bool, extMapping map[string]string) error {
	var scanned, moved, skipped int
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		scanned += 1
		if err != nil {
			return err
		}

		// skip the root directory itself
		if path == root {
			skipped += 1
			return nil
		}

		// skip subdirectories if not in recursive mode
		if d.IsDir() {
			skipped += 1
			if !recursive {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == "" {
			ext = "others"
		} else {
			ext = ext[1:] // remove the dot
		}

		// ignore extension
		if ignoreExt[ext] {
			skipped += 1
			fmt.Printf("ignore path: %s\n", path)
			return nil
		}

		// custom folder extension mapping
		if mapping, ok := extMapping[ext]; ok {
			ext = mapping
		}

		destDir := filepath.Join(root, ext)
		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			skipped += 1
			return err
		}

		destPath := filepath.Join(destDir, filepath.Base(path))
		newPath, err := resolveConflict(destPath)
		if err != nil {
			skipped += 1
			return err
		}

		fmt.Printf("-> %s -> %s\n", path, newPath)
		if !dryRun {
			err = os.Rename(path, newPath)
			if err != nil {
				skipped += 1
				return err
			}
		}
		moved += 1

		return nil
	})
	if err != nil {
		return err
	}

	fmt.Println("\n✔ Summary:")
	fmt.Printf("  Total files scanned: %d\n", scanned)
	fmt.Printf("  Moved: %d\n", moved)
	fmt.Printf("  Skipped: %d\n", skipped)
	fmt.Printf("  (Dry run: %v)\n", dryRun)

	return nil
}

func resolveConflict(destPath string) (string, error) {
	dir := filepath.Dir(destPath)
	name := filepath.Base(destPath)
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	i := 1
	for {
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			return destPath, nil // no conflict
		}

		// Generate new file name with (1), (2), etc.
		newName := fmt.Sprintf("%s(%d)%s", base, i, ext)
		destPath = filepath.Join(dir, newName)
		i++
		if i > 9999 {
			return "", fmt.Errorf("too many conflicting files for %s", name)
		}
	}
}
