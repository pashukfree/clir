package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Define lipgloss styles for various output elements
var (
	titleStyle         = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("57"))    // Deep Pink/Magenta
	headerStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))    // Bright Purple
	directoryNameStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("211"))   // Bright Pink
	pathStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))              // Lighter Grey
	infoStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))              // Very Light Grey
	successStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("84"))               // Teal/Green
	warningStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))              // Gold/Yellow
	errorStyle         = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("197"))   // Hot Pink/Red
	promptStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("227"))              // Light Yellow
	sizeStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("213"))              // Soft Pink
	skippedItemStyle   = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("242")) // Grey Italic (unchanged)
	separatorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))              // Darker Grey
	emojiInfoStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("117"))              // Light Blue
	emojiSuccessStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("47"))               // Bright Green
	emojiWarningStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))              // Orange (unchanged)
	emojiErrorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))              // Red (unchanged)
)

// Style for the main title box
var titleBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("211")).
	PaddingLeft(5).PaddingRight(5)

var bannerBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("63")).
	PaddingLeft(2).PaddingRight(2)

// dirSize calculates the total size of files within a directory and counts skipped items due to permissions.
func dirSize(path string) (int64, int, error) {
	var size int64
	var skippedCount int
	var pathError error // To store error specific to the initial path stat

	// Walk through the directory tree
	err := filepath.Walk(path, func(currentPath string, info fs.FileInfo, err error) error {
		// Handle errors during walk
		if err != nil {
			// If permission denied, increment skipped count and skip directory if it's a directory
			if os.IsPermission(err) {
				skippedCount++
				if info != nil && info.IsDir() {
					return filepath.SkipDir // Skip the entire directory if permission is denied
				}
				return nil // Skip this file/item if permission is denied
			}
			// Store the error if it's for the initial path
			if currentPath == path && pathError == nil {
				pathError = err
			}
			// Return other errors unless it's a SkipDir error
			if err != filepath.SkipDir {
				return err
			}
			return nil
		}
		// If it's not a directory, add its size to the total
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	// Return accumulated size, skipped count, and any error (prioritizing pathError)
	if err != nil && err != filepath.SkipDir {
		return size, skippedCount, err
	}
	if pathError != nil {
		return size, skippedCount, pathError
	}
	return size, skippedCount, nil
}

// formatSize converts bytes to a human-readable string (B, KB, MB, GB).
func formatSize(size int64) string {
	if size < 0 {
		return errorStyle.Render("N/A") // Handle error state size
	}
	// Convert and format based on size magnitude
	if size < 1024 {
		return sizeStyle.Render(fmt.Sprintf("%d B", size))
	}
	sizeKB := float64(size) / 1024
	if sizeKB < 1024 {
		return sizeStyle.Render(fmt.Sprintf("%.1f KB", sizeKB))
	}
	sizeMB := sizeKB / 1024
	if sizeMB < 1024 {
		return sizeStyle.Render(fmt.Sprintf("%.1f MB", sizeMB))
	}
	sizeGB := sizeMB / 1024
	return sizeStyle.Render(fmt.Sprintf("%.1f GB", sizeGB))
}

// confirm prompts the user with a yes/no question and returns true if 'y' or 'yes' is entered.
func confirm(promptMsg string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(promptStyle.Render(fmt.Sprintf("%s [y/N]: ", promptMsg)))
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))
		// Check user input
		if input == "y" || input == "yes" {
			return true
		} else if input == "n" || input == "no" || input == "" {
			return false
		}
		// Invalid input message
		fmt.Println(warningStyle.Render("   Invalid input. Please enter 'y' or 'n'."))
	}
}

// clearDirectoryContents attempts to remove all items within a specified directory.
func clearDirectoryContents(path string) (int, int, error) {
	var deletedCount int
	var skippedCount int

	// Read directory entries
	dirs, err := os.ReadDir(path)
	if err != nil {
		// Handle permission denied error for reading directory
		if os.IsPermission(err) {
			return 0, 0, fmt.Errorf("list %s: %s", pathStyle.Render(path), errorStyle.Render("permission denied"))
		}
		// Handle other read directory errors
		return 0, 0, fmt.Errorf("read dir %s: %w", pathStyle.Render(path), err)
	}

	// Check if directory is empty
	if len(dirs) == 0 {
		fmt.Println(infoStyle.Render(fmt.Sprintf("   %s Empty or all inaccessible.", emojiInfoStyle.Render("💨"))))
		return 0, 0, nil
	}

	// Print clearing message
	fmt.Println(infoStyle.Render(fmt.Sprintf("   %s Clearing %s...", emojiInfoStyle.Render("✨"), pathStyle.Render(path))))
	// Iterate and remove each item
	for _, d := range dirs {
		itemPath := filepath.Join(path, d.Name())
		err := os.RemoveAll(itemPath)
		// Handle removal errors
		if err != nil {
			// If permission denied, skip and continue
			if os.IsPermission(err) {
				skippedCount++
				continue
			} else {
				// Print warning for other errors and skip
				fmt.Println(warningStyle.Render(fmt.Sprintf("      %s Del %s: %v", emojiWarningStyle.Render("⚠️"), pathStyle.Render(itemPath), err)))
				skippedCount++
			}
		} else {
			// Print success message for deleted item
			fmt.Println(infoStyle.Render(fmt.Sprintf("      %s %s", emojiSuccessStyle.Render("🗑️"), pathStyle.Render(itemPath)))) // Changed emoji
			deletedCount++
		}
	}

	// Print summary of deleted and skipped items
	if deletedCount > 0 {
		fmt.Println(successStyle.Render(fmt.Sprintf("   %s Deleted %d items.", emojiSuccessStyle.Render("✅"), deletedCount)))
	}
	if skippedCount > 0 {
		fmt.Println(warningStyle.Render(fmt.Sprintf("   %s Skipped %d items.", emojiWarningStyle.Render("⚠️"), skippedCount)))
	}
	// Message if no items were touched
	if deletedCount == 0 && skippedCount == 0 && len(dirs) > 0 {
		fmt.Println(infoStyle.Render(fmt.Sprintf("   %s No items touched (all might be inaccessible).", emojiInfoStyle.Render("🤔"))))
	}

	return deletedCount, skippedCount, nil
}

// DirToScan struct holds information about a directory to be scanned and potentially cleared.
type DirToScan struct {
	Path             string // Absolute path to the directory
	DisplayName      string // User-friendly name for the directory
	Category         string // Category of the directory (e.g., Cache, Logs)
	Warning          string // Optional warning message for the user
	ScannedSize      int64  // Size of accessible files in the directory
	InitialSkipCount int    // Number of items skipped during the initial scan
	HasCriticalError bool   // Flag indicating a critical error during scan
	ScanErrorMsg     string // Error message if a critical error occurred
}

func main() {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("%s Error getting home directory: %v", emojiErrorStyle.Render("❌"), err)))
		return
	}

	// Define the list of directories to scan and potentially clear
	dirsToScan := []DirToScan{
		{Path: filepath.Join(homeDir, "Library", "Caches"), DisplayName: "System Caches", Category: "Cache", Warning: "Some cache subfolders might be system-protected & skipped."},
		{Path: filepath.Join(homeDir, "Library", "Logs"), DisplayName: "System Logs", Category: "Logs", Warning: ""},
	}

	var totalInitialScannedSize int64

	// First pass: Scan directories to get initial sizes and identify errors
	for i := range dirsToScan {
		dirInfo := &dirsToScan[i]

		// Check if the directory exists
		if _, statErr := os.Stat(dirInfo.Path); os.IsNotExist(statErr) {
			msg := fmt.Sprintf("Not found: %s. Skip.", pathStyle.Render(dirInfo.Path))
			dirInfo.ScannedSize = -1 // Indicate not found
			dirInfo.HasCriticalError = true
			dirInfo.ScanErrorMsg = msg
			continue // Skip to the next directory
		}

		// Calculate directory size and count skipped items
		dirPathSize, skippedItems, sizeErr := dirSize(dirInfo.Path)
		dirInfo.InitialSkipCount = skippedItems

		// Handle errors during size calculation
		if sizeErr != nil {
			msg := fmt.Sprintf("Scan %s: %v", directoryNameStyle.Render(dirInfo.DisplayName), sizeErr)
			dirInfo.ScannedSize = 0 // Indicate error, size unknown
			dirInfo.HasCriticalError = true
			dirInfo.ScanErrorMsg = msg
		} else {
			dirInfo.ScannedSize = dirPathSize
			totalInitialScannedSize += dirPathSize
		}
	}

	// Print the application title
	fmt.Println(titleBoxStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left,
		titleStyle.Render(fmt.Sprintf("%s Clir:", titleStyle.Render("🌸"))),
		infoStyle.Render(fmt.Sprintf(" %s", formatSize(totalInitialScannedSize))),
	)))

	// Second pass: Display info and prompt user to clear each directory
	for i := range dirsToScan {
		dirInfo := &dirsToScan[i]

		// Print directory display name and path
		fmt.Println(directoryNameStyle.Render(fmt.Sprintf("%s Dir: %s", emojiInfoStyle.Render("📁"), dirInfo.DisplayName)))
		fmt.Println(infoStyle.Render(fmt.Sprintf("   Path: %s", pathStyle.Render(dirInfo.Path))))

		// If a critical error occurred during scan, report it and skip clearing
		if dirInfo.HasCriticalError {
			fmt.Println(errorStyle.Render(fmt.Sprintf("   Status: %s", dirInfo.ScanErrorMsg)))
			fmt.Println(separatorStyle.Render(strings.Repeat("-", 30)))
			continue
		}

		// Print the scanned size
		fmt.Println(infoStyle.Render(fmt.Sprintf("   Size: %s", formatSize(dirInfo.ScannedSize))))

		// Handle cases where the directory is empty or all contents are inaccessible
		if dirInfo.ScannedSize == 0 && dirInfo.InitialSkipCount == 0 {
			fmt.Println(infoStyle.Render(fmt.Sprintf("   %s Empty.", emojiInfoStyle.Render("💨"))))
			fmt.Println(separatorStyle.Render(strings.Repeat("-", 30)))
			continue
		} else if dirInfo.ScannedSize == 0 && dirInfo.InitialSkipCount > 0 {
			fmt.Println(infoStyle.Render(fmt.Sprintf("   %s All contents inaccessible.", emojiInfoStyle.Render("💨"))))
			fmt.Println(separatorStyle.Render(strings.Repeat("-", 30)))
			continue
		}

		// Prompt user to confirm clearing
		if confirm(fmt.Sprintf("   %s Clear %s?", emojiWarningStyle.Render("🗑️"), directoryNameStyle.Render(dirInfo.DisplayName))) { // Shortened
			// Attempt to clear directory contents
			_, _, clearErr := clearDirectoryContents(dirInfo.Path)
			if clearErr != nil {
				fmt.Println(errorStyle.Render(fmt.Sprintf("   %s Clear %s: %v", emojiErrorStyle.Render("❌"), directoryNameStyle.Render(dirInfo.DisplayName), clearErr))) // Shortened
			}

			// Recalculate size after clearing attempt
			newSize, newSkipped, newSizeErr := dirSize(dirInfo.Path)
			dirInfo.InitialSkipCount = newSkipped
			if newSizeErr != nil {
				fmt.Println(warningStyle.Render(fmt.Sprintf("   %s Size recalc %s: %v", emojiWarningStyle.Render("⚠️"), directoryNameStyle.Render(dirInfo.DisplayName), newSizeErr))) // Shortened
				dirInfo.ScannedSize = -2                                                                                                                                              // Indicate error during recalculation
			} else {
				dirInfo.ScannedSize = newSize
			}
			// Print new size and inaccessible item count
			fmt.Println(infoStyle.Render(fmt.Sprintf("   %s New Size %s: %s", emojiInfoStyle.Render("📊"), directoryNameStyle.Render(dirInfo.DisplayName), formatSize(dirInfo.ScannedSize))))
			if dirInfo.InitialSkipCount > 0 {
				inaccessibleMsgPart := skippedItemStyle.Render(fmt.Sprintf("(%d inaccessible)", dirInfo.InitialSkipCount))
				fmt.Println(infoStyle.Render(fmt.Sprintf("      %s", inaccessibleMsgPart)))
			} else {
				fmt.Println(skippedItemStyle.Render(fmt.Sprintf("      (%d still inaccessible)", dirInfo.InitialSkipCount)))
			}
		} else {
			// Message if clearing was skipped by the user
			fmt.Println(successStyle.Render(fmt.Sprintf("   %s Skipped clearing.", emojiSuccessStyle.Render("👍"))))
		}
		// Print separator line
		fmt.Println(separatorStyle.Render(strings.Repeat("-", 30))) // Shortened separator
	}

	// Final Summary section
	fmt.Println(headerStyle.Render(fmt.Sprintf("\n%s Final Summary:", emojiSuccessStyle.Render("🎉"))))
	var finalTotalAccessibleSize int64
	// Iterate through directories for final status report
	for _, dirInfo := range dirsToScan {
		// Report critical errors or post-clear errors
		if dirInfo.HasCriticalError {
			fmt.Println(errorStyle.Render(fmt.Sprintf("   %s: %s", directoryNameStyle.Render(dirInfo.DisplayName), dirInfo.ScanErrorMsg)))
		} else if dirInfo.ScannedSize == -2 {
			fmt.Println(errorStyle.Render(fmt.Sprintf("   %s: Error post-clear", directoryNameStyle.Render(dirInfo.DisplayName))))
		} else {
			// Report final size and inaccessible count
			baseMsg := fmt.Sprintf("   %s: %s", directoryNameStyle.Render(dirInfo.DisplayName), formatSize(dirInfo.ScannedSize))
			if dirInfo.InitialSkipCount > 0 {
				skippedMsg := skippedItemStyle.Render(fmt.Sprintf("(%d inaccessible)", dirInfo.InitialSkipCount))
				fmt.Println(infoStyle.Render(fmt.Sprintf("%s %s", baseMsg, skippedMsg)))
			} else {
				fmt.Println(infoStyle.Render(baseMsg))
			}
			// Add accessible size to final total
			if dirInfo.ScannedSize > 0 {
				finalTotalAccessibleSize += dirInfo.ScannedSize
			}
		}
	}

	// Print banner and reboot message
	fmt.Println()
	fmt.Println(bannerBoxStyle.Render(titleStyle.Render(fmt.Sprintf("%s, %s %s \n%s", warningStyle.Render("Stand with Ukrainian defenders"), infoStyle.Render("donate to"), directoryNameStyle.Render("Come Back Alive Fund"), headerStyle.Render("https://savelife.in.ua/en/donate-en/#donate-army-card-once")))))
	fmt.Println(directoryNameStyle.Render("\nTo trigger MacOS to recalculate your disk space, reboot your Mac"))
	fmt.Println()
}
