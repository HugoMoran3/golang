package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// FileMove represents a file move operation
type FileMove struct {
	OriginalPath string
	NewPath      string
	FileSize     int64
}

// FileMoveError represents a file move error
type FileMoveError struct {
	Path  string
	Error error
}

// getUniqueFilename ensures the filename is unique in the target directory
// only adding numbers if there's an actual conflict
func getUniqueFilename(targetDir, originalFilename string, dryRun bool) string {
	baseFile := filepath.Base(originalFilename)
	newPath := filepath.Join(targetDir, baseFile)
	
	if dryRun {
		// For dry run, check if we've already "used" this filename in our simulation
		return newPath
	}
	
	// If the file doesn't exist, use the original name
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		return newPath
	}
	
	// If there's a conflict, then start adding numbers
	ext := filepath.Ext(baseFile)
	nameWithoutExt := strings.TrimSuffix(baseFile, ext)
	counter := 1
	
	for {
		newPath = filepath.Join(targetDir, fmt.Sprintf("%s_%d%s", nameWithoutExt, counter, ext))
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
		counter++
	}
}

// MoveLongPaths finds and moves files with paths longer than maxLength
func MoveLongPaths(sourceDir, targetDir string, maxLength int, dryRun bool) ([]FileMove, []FileMoveError) {
	var movedFiles []FileMove
	var errorFiles []FileMoveError
	usedNames := make(map[string]bool) // Track used names in dry run mode

	// Convert to absolute paths
	absSourceDir, err := filepath.Abs(sourceDir)
	if err != nil {
		fmt.Printf("Error resolving source path: %v\n", err)
		return movedFiles, errorFiles
	}

	absTargetDir, err := filepath.Abs(targetDir)
	if err != nil {
		fmt.Printf("Error resolving target path: %v\n", err)
		return movedFiles, errorFiles
	}

	fmt.Printf("\nMode: %s\n", map[bool]string{true: "DRY RUN (no files will be moved)", false: "ACTUAL RUN"}[dryRun])
	fmt.Printf("Processing source directory: %s\n", absSourceDir)
	fmt.Printf("Target directory: %s\n", absTargetDir)
	fmt.Printf("Maximum path length: %d\n\n", maxLength)

	if !dryRun {
		// Create target directory if it doesn't exist
		err = os.MkdirAll(absTargetDir, 0755)
		if err != nil {
			fmt.Printf("Error creating target directory: %v\n", err)
			return movedFiles, errorFiles
		}
	}

	// Walk through all files in source directory
	err = filepath.Walk(absSourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errorFiles = append(errorFiles, FileMoveError{Path: path, Error: err})
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return nil // Continue walking
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check path length
		if len(path) > maxLength {
			baseFile := filepath.Base(path)
			newPath := filepath.Join(absTargetDir, baseFile)
			
			if dryRun {
				// In dry run mode, simulate file naming conflicts
				if usedNames[baseFile] {
					// If name is already used, add a number
					ext := filepath.Ext(baseFile)
					nameWithoutExt := strings.TrimSuffix(baseFile, ext)
					counter := 1
					for {
						testName := fmt.Sprintf("%s_%d%s", nameWithoutExt, counter, ext)
						if !usedNames[testName] {
							newPath = filepath.Join(absTargetDir, testName)
							usedNames[testName] = true
							break
						}
						counter++
					}
				} else {
					usedNames[baseFile] = true
				}
				
				movedFiles = append(movedFiles, FileMove{
					OriginalPath: path,
					NewPath:      newPath,
					FileSize:     info.Size(),
				})
			} else {
				newPath = getUniqueFilename(absTargetDir, path, false)
				
				// Actually move the file
				err = copyFile(path, newPath)
				if err != nil {
					errorFiles = append(errorFiles, FileMoveError{Path: path, Error: err})
					fmt.Printf("Error copying %s: %v\n", path, err)
					return nil
				}

				// Verify sizes match
				srcInfo, err := os.Stat(path)
				if err != nil {
					errorFiles = append(errorFiles, FileMoveError{Path: path, Error: err})
					fmt.Printf("Error verifying source file %s: %v\n", path, err)
					return nil
				}

				dstInfo, err := os.Stat(newPath)
				if err != nil {
					errorFiles = append(errorFiles, FileMoveError{Path: path, Error: err})
					fmt.Printf("Error verifying destination file %s: %v\n", newPath, err)
					return nil
				}

				if srcInfo.Size() == dstInfo.Size() {
					// Delete original file
					err = os.Remove(path)
					if err != nil {
						errorFiles = append(errorFiles, FileMoveError{Path: path, Error: err})
						fmt.Printf("Error removing original file %s: %v\n", path, err)
						return nil
					}

					movedFiles = append(movedFiles, FileMove{
						OriginalPath: path,
						NewPath:      newPath,
						FileSize:     srcInfo.Size(),
					})
					fmt.Printf("Moved: %s\nTo: %s\n\n", path, newPath)
				} else {
					err := fmt.Errorf("size mismatch after copy")
					errorFiles = append(errorFiles, FileMoveError{Path: path, Error: err})
					fmt.Printf("Error: size mismatch after copying %s\n", path)
					// Attempt to clean up failed copy
					os.Remove(newPath)
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}

	return movedFiles, errorFiles
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the contents
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Sync to ensure write is complete
	return destFile.Sync()
}

// formatFileSize returns a human-readable file size
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}

func main() {
	// Define command line flags
	dryRun := flag.Bool("dry-run", false, "Perform a dry run (no files will be moved)")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	// Get source directory
	fmt.Print("Enter source directory path: ")
	sourceDir, _ := reader.ReadString('\n')
	sourceDir = strings.TrimSpace(sourceDir)

	// Get target directory
	fmt.Print("Enter target directory path: ")
	targetDir, _ := reader.ReadString('\n')
	targetDir = strings.TrimSpace(targetDir)

	// Get maximum length
	fmt.Print("Enter maximum path length (press Enter for default 180): ")
	maxLengthStr, _ := reader.ReadString('\n')
	maxLengthStr = strings.TrimSpace(maxLengthStr)

	maxLength := 180
	if maxLengthStr != "" {
		if val, err := strconv.Atoi(maxLengthStr); err == nil {
			maxLength = val
		}
	}

	movedFiles, errorFiles := MoveLongPaths(sourceDir, targetDir, maxLength, *dryRun)

	// Print summary
	fmt.Printf("\nFound %d files with paths longer than %d characters\n", len(movedFiles), maxLength)

	if len(movedFiles) > 0 {
		fmt.Println("\nFiles to be moved:")
		var totalSize int64
		for _, move := range movedFiles {
			totalSize += move.FileSize
			fmt.Printf("\nFrom: %s\nTo: %s\nSize: %s\n", 
				move.OriginalPath, 
				move.NewPath,
				formatFileSize(move.FileSize))
		}
		fmt.Printf("\nTotal size: %s\n", formatFileSize(totalSize))
	}

	if len(errorFiles) > 0 {
		fmt.Printf("\nEncountered %d errors:\n", len(errorFiles))
		for _, err := range errorFiles {
			fmt.Printf("File: %s\nError: %v\n\n", err.Path, err.Error)
		}
	}

	if *dryRun {
		fmt.Println("\nThis was a dry run - no files were actually moved.")
		fmt.Println("Run without --dry-run flag to perform the actual move operation.")
	}
}
