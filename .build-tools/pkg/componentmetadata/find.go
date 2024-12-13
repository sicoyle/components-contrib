/*
Copyright 2022 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package componentmetadata

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindComponents finds all component folders and returns their paths.
func FindComponents(folders []string, skip []string) ([]string, error) {
	skipMap := map[string]struct{}{}
	for _, v := range skip {
		// Normalize all slashes
		v = filepath.Clean(v)
		skipMap[v] = struct{}{}
	}

	res := []string{}
	for _, folder := range folders {
		folder = filepath.Clean(folder)
		err := findInDirectory(folder, skipMap, &res)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func findInDirectory(dir string, skip map[string]struct{}, res *[]string) error {
	read, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, e := range read {
		// Ignore anything but directories
		if !e.IsDir() {
			continue
		}

		path := filepath.Join(dir, e.Name())

		// Add the directory if not skipped
		if _, ok := skip[path]; !ok {
			*res = append(*res, path)
		} else {
			fmt.Fprintln(os.Stderr, "Info: skipped folder "+path)
		}

		// Read the directory recursively
		findInDirectory(path, skip, res)
	}
	return nil
}

func FindValidComponents(folders []string, skip []string) ([]string, error) {
	skipMap := make(map[string]struct{}, len(skip))
	for _, v := range skip {
		// Normalize all slashes
		v = filepath.Clean(v)
		skipMap[v] = struct{}{}
	}

	var result []string
	for _, folder := range folders {
		folder = filepath.Clean(folder)
		if err := findValidInDirectory(folder, skipMap, &result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func findValidInDirectory(dir string, skip map[string]struct{}, res *[]string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		path := filepath.Join(dir, entry.Name())

		// Skip the folder if it exists in skipMap
		if _, shouldSkip := skip[path]; shouldSkip {
			fmt.Fprintln(os.Stderr, "Info: skipped folder "+path)
			continue
		}

		// Add the folder to the result list
		*res = append(*res, path)

		// Recursively process the folder
		if err := findInDirectory(path, skip, res); err != nil {
			return err
		}
	}

	return nil
}
