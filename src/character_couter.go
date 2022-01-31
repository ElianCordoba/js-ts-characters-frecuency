package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func main() {
	root_directory := "./data"

	// Taken from ScriptElementKindModifier in typescript/src/services/type.js
	// .d.ts are counted as .ts, same with mts and cts
	file_extensions_to_include := []string{".ts", ".tsx", ".js", ".jsx", ".mts", ".mjs", ".cts", ".cjs"}

	found_files, file_extensions_found := scan_for_files(root_directory, file_extensions_to_include)

	fmt.Println(file_extensions_found)

	result := count_characters(found_files, root_directory)

	jsonRaw, err := json.Marshal(result)

	if err != nil {
		return
	}

	write_to_file(jsonRaw)
}

func scan_for_files(root_directory string, file_extensions_to_include []string) ([]string, map[string]int) {
	var found_files []string
	file_extensions_found := make(map[string]int)

	accept_file_extension := func(file_extension string) bool {
		return string_in_slice(file_extensions_to_include, file_extension)
	}

	fileSystem := os.DirFS(root_directory)

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			log.Fatal(err)
		}

		if d.IsDir() {
			return nil
		}

		file_extension := filepath.Ext(path)

		// Skip files with unwanted extensions
		if !accept_file_extension(file_extension) {
			return nil
		}

		file_extensions_found[file_extension] = file_extensions_found[file_extension] + 1

		found_files = append(found_files, path)

		return nil
	})

	return found_files, file_extensions_found
}

func count_characters(files []string, root_directory string) map[byte]int {
	// Key: Ascii code
	// Value: Number of occurrences
	stats := make(map[byte]int)

	for _, file_path := range files {
		file, err := os.Open(root_directory + "/" + file_path)

		if err != nil {
			fmt.Println(err)
		}

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanRunes)

		for scanner.Scan() {
			ascii_code := scanner.Text()[0]
			stats[ascii_code] = stats[ascii_code] + 1
		}

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}

		file.Close()
	}

	return stats
}

func write_to_file(content []byte) {
	err := os.WriteFile("./output.json", content, 0644)

	if err != nil {
		return
	}
}

// Util
func string_in_slice(slice []string, string_to_find string) bool {
	for _, b := range slice {
		if b == string_to_find {
			return true
		}
	}
	return false
}
