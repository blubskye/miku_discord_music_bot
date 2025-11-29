/*
 * Miku Discord Music Bot
 * Copyright (C) 2025 blubskye (https://github.com/blubskye)
 * Discord: blubaustin
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 * Source: https://github.com/blubskye/miku_discord_music_bot
 */

package music

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type LocalFile struct {
	Name     string
	Path     string
	Folder   string
	Duration int // Duration in seconds (can be extracted from metadata later)
}

type Library struct {
	rootPath string
	files    map[string][]*LocalFile // folder -> files
	mu       sync.RWMutex
}

var supportedExtensions = map[string]bool{
	".mp3":  true,
	".flac": true,
	".wav":  true,
	".ogg":  true,
	".m4a":  true,
	".opus": true,
	".aac":  true,
	".wma":  true,
}

func NewLibrary(rootPath string) (*Library, error) {
	if rootPath == "" {
		return nil, fmt.Errorf("music folder path is not configured")
	}

	// Check if path exists
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("music folder does not exist: %s", rootPath)
	}

	lib := &Library{
		rootPath: rootPath,
		files:    make(map[string][]*LocalFile),
	}

	// Scan the directory
	if err := lib.Scan(); err != nil {
		return nil, err
	}

	return lib, nil
}

func (l *Library) Scan() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Clear existing files
	l.files = make(map[string][]*LocalFile)

	// Walk through the directory
	err := filepath.WalkDir(l.rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Check if file has supported extension
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if !supportedExtensions[ext] {
			return nil
		}

		// Get relative path from root
		relPath, err := filepath.Rel(l.rootPath, path)
		if err != nil {
			return err
		}

		// Get folder name (use "root" if file is in root directory)
		folder := filepath.Dir(relPath)
		if folder == "." {
			folder = "root"
		}

		// Create local file entry
		file := &LocalFile{
			Name:   d.Name(),
			Path:   path,
			Folder: folder,
		}

		// Add to files map
		l.files[folder] = append(l.files[folder], file)

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to scan music folder: %w", err)
	}

	return nil
}

func (l *Library) GetFolders() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	folders := make([]string, 0, len(l.files))
	for folder := range l.files {
		folders = append(folders, folder)
	}

	return folders
}

func (l *Library) GetFiles(folder string) []*LocalFile {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if files, exists := l.files[folder]; exists {
		// Return a copy to prevent race conditions
		filesCopy := make([]*LocalFile, len(files))
		copy(filesCopy, files)
		return filesCopy
	}

	return []*LocalFile{}
}

func (l *Library) GetAllFiles() []*LocalFile {
	l.mu.RLock()
	defer l.mu.RUnlock()

	allFiles := make([]*LocalFile, 0)
	for _, files := range l.files {
		allFiles = append(allFiles, files...)
	}

	return allFiles
}

func (l *Library) SearchByName(query string) []*LocalFile {
	l.mu.RLock()
	defer l.mu.RUnlock()

	query = strings.ToLower(query)
	results := make([]*LocalFile, 0)

	for _, files := range l.files {
		for _, file := range files {
			if strings.Contains(strings.ToLower(file.Name), query) {
				results = append(results, file)
			}
		}
	}

	return results
}

func (l *Library) GetFileByFolderAndName(folder, name string) (*LocalFile, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	files, exists := l.files[folder]
	if !exists {
		return nil, fmt.Errorf("folder not found: %s", folder)
	}

	// Try exact match first
	for _, file := range files {
		if file.Name == name {
			return file, nil
		}
	}

	// Try case-insensitive match
	nameLower := strings.ToLower(name)
	for _, file := range files {
		if strings.ToLower(file.Name) == nameLower {
			return file, nil
		}
	}

	// Try partial match
	for _, file := range files {
		if strings.Contains(strings.ToLower(file.Name), nameLower) {
			return file, nil
		}
	}

	return nil, fmt.Errorf("file not found: %s in folder %s", name, folder)
}

func (l *Library) GetTotalFiles() int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	count := 0
	for _, files := range l.files {
		count += len(files)
	}

	return count
}
