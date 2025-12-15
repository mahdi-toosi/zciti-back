// Package utils provides common utility functions used throughout the application.
// It includes helpers for HTTP request handling, file operations, image processing,
// random generation, and date/time manipulation.
package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-fiber-starter/app/database/schema"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

// IsChildProcess returns true if the current process is a child process.
// Used with Fiber's prefork mode to detect worker processes.
func IsChildProcess() bool {
	return fiber.IsChild()
}

// Log outputs the given variable to debug log with zerolog.
// Includes a warning prefix when running in a child process (prefork mode).
func Log(variable any) {
	if fiber.IsChild() {
		log.Warn().Msg("IN CHILD PROCESS ")
	}

	log.Debug().Interface("", variable).Msg("LOGGER ")
}

// IsEnabled returns a Fiber skip function based on the enabled state.
// Returns nil (middleware enabled) if key is true, or a skip-all function if false.
// Used for conditionally enabling/disabling middleware.
func IsEnabled(key bool) func(c *fiber.Ctx) bool {
	if key {
		return nil
	}

	return func(c *fiber.Ctx) bool { return true }
}

// InlineCondition provides ternary-like conditional logic for Go.
// Returns ifResult when condition is true, otherwise returns elseResult.
func InlineCondition(condition bool, ifResult any, elseResult any) any {
	if condition {
		return ifResult
	}
	return elseResult
}

// IsForUser checks if the current request context has "forUser" flag set.
// Returns false if the flag is not set or not a boolean.
func IsForUser(c *fiber.Ctx) bool {
	val, ok := c.Locals("forUser").(bool)
	if !ok {
		return false
	}
	return val
}

// GetIntInParams parses a uint64 from the route parameters by key.
func GetIntInParams(c *fiber.Ctx, key string) (uint64, error) {
	return strconv.ParseUint(c.Params(key), 10, 64)
}

// GetUintInQueries parses a uint64 from the query string by key.
func GetUintInQueries(c *fiber.Ctx, key string) (uint64, error) {
	return strconv.ParseUint(c.Query(key), 10, 64)
}

// GetIntInQueries parses an int64 from the query string by key.
func GetIntInQueries(c *fiber.Ctx, key string) (int64, error) {
	return strconv.ParseInt(c.Query(key), 10, 64)
}

// GetDateInQueries parses a date from the query string by key using Asia/Tehran timezone.
// Returns a zero time if the query parameter is empty or parsing fails.
// Expected format: YYYY-MM-DD (time.DateOnly).
func GetDateInQueries(c *fiber.Ctx, key string) *time.Time {
	if c.Query(key) == "" {
		return &time.Time{}
	}
	loc, _ := time.LoadLocation("Asia/Tehran")
	date, err := time.ParseInLocation(time.DateOnly, c.Query(key), loc)
	if err != nil {
		return &time.Time{}
	}
	return &date
}

// GetAuthenticatedUser retrieves the authenticated user from the request context.
// The user is expected to be set by the authentication middleware in c.Locals("user").
// Returns an error with Persian message if the user is not authenticated.
func GetAuthenticatedUser(c *fiber.Ctx) (schema.User, error) {
	user, ok := c.Locals("user").(schema.User)
	if ok && c.Locals("user") != nil {

		return user, nil
	}
	return schema.User{}, errors.New("ابتدا وارد حساب کاربری خود شوید و دوباره تلاش کنید")
}

// ValidateMobileNumber validates an Iranian mobile phone number format.
// Valid format: starts with 9, followed by exactly 9 digits (e.g., "9123456789").
// Returns a fiber error with Persian message if validation fails.
func ValidateMobileNumber(number string) error {
	// Iranian mobile numbers: 9XXXXXXXXX (10 digits starting with 9)
	pattern := `^9\d{9}$`
	regex := regexp.MustCompile(pattern)

	valid := regex.MatchString(number)
	if !valid {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "تلفن همراه معتبر نمی باشد",
		}
	}

	return nil
}

// Random generates a random uint64 between min (inclusive) and max (exclusive).
func Random(min int, max int) uint64 {
	return uint64(rand.Intn(max-min) + min) //nolint:gosec
}

// RandomFromArray returns a random element from the given uint64 slice.
// Panics if the array is empty.
func RandomFromArray(arr []uint64) uint64 {
	randomIndex := rand.Intn(len(arr))
	return arr[randomIndex]
}

// GetFakeTableIDs retrieves all IDs from a GORM model table.
// Used for generating fake/test data with valid foreign key references.
func GetFakeTableIDs(db *gorm.DB, table any) (ids []uint64, err error) {
	err = db.Model(&table).Select("id").Find(&ids).Error
	if err != nil {
		log.Err(err)
	}
	return
}

// GetFakeTableIDsWithConditions retrieves IDs from a GORM model table with WHERE conditions.
// Conditions map format: {"column = ?": []any{value}} or {"column IN (?)": []any{values...}}.
// Used for generating fake/test data with filtered foreign key references.
func GetFakeTableIDsWithConditions(db *gorm.DB, table any, conditions map[string][]any) (ids []uint64, err error) {
	query := db.Model(&table).Select("id")

	for key, value := range conditions {
		if len(value) > 0 {
			query = query.Where(key, value)
		}
	}

	err = query.Find(&ids).Error
	if err != nil {
		log.Err(err)
	}
	return
}

// RandomDateTime generates a random date and time.
// Note: Year range is 0-2022, which may produce invalid dates (e.g., Feb 31).
// Primarily used for test data generation.
func RandomDateTime() time.Time {
	year := rand.Intn(2023) //nolint:gosec
	month := time.Month(rand.Intn(12) + 1)
	day := rand.Intn(31) + 1 //nolint:gosec
	hour := rand.Intn(24)    //nolint:gosec
	min := rand.Intn(60)     //nolint:gosec
	sec := rand.Intn(60)     //nolint:gosec

	return time.Date(year, month, day, hour, min, sec, 0, time.UTC)
}

// imageOptimizedWidths defines the widths (in pixels) for generating optimized image versions.
// Used for responsive images to reduce bandwidth on different screen sizes.
var imageOptimizedWidths = []int{600, 300}

// StoreImageOptimizedVersions creates resized versions of an image at predefined widths.
// For each width in imageOptimizedWidths, creates a new file with format: {name}-{width}x{width}.{ext}
// Returns the total size of all created optimized images in bytes.
func StoreImageOptimizedVersions(path string, fileName string) (resultSizes int64, err error) {
	for _, size := range imageOptimizedWidths {
		resultSize, err := resizeAndSaveImage(path, fileName, size)
		if err != nil {
			return 0, err
		}
		resultSizes += resultSize
	}
	return resultSizes, nil
}

// resizeAndSaveImage resizes an image to the specified width while maintaining aspect ratio.
// Output filename format: {originalName}-{width}x{width}.{ext}
// Supports JPEG, JPG, and PNG formats. Skips resizing if image is smaller than target width.
// Uses temporary file for atomic write operation.
func resizeAndSaveImage(folderPath string, imageName string, newWidth int) (size int64, err error) {
	tempFile, err := ioutil.TempFile("", "resize-*.jpg")
	if err != nil {
		return 0, err
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	inputPath := filepath.Join(folderPath, imageName)
	img, err := imaging.Open(inputPath)
	if err != nil {
		return 0, err
	}

	// Skip if image is already smaller than target width (no upscaling)
	if newWidth > img.Bounds().Size().X {
		return 0, nil
	}

	// Resize with NearestNeighbor for performance (height=0 maintains aspect ratio)
	resizedImg := imaging.Resize(img, newWidth, 0, imaging.NearestNeighbor)

	// Encode based on original file extension
	ext := filepath.Ext(imageName)
	switch ext {
	case ".jpeg":
		err = imaging.Encode(tempFile, resizedImg, imaging.JPEG)
	case ".jpg":
		err = imaging.Encode(tempFile, resizedImg, imaging.JPEG)
	case ".png":
		err = imaging.Encode(tempFile, resizedImg, imaging.PNG)
	}
	if err != nil {
		return 0, err
	}

	// Build output filename: {baseName}-{width}x{width}.{ext}
	baseName := strings.TrimSuffix(imageName, ext)
	newFileName := fmt.Sprintf("%s-%dx%d%s", baseName, newWidth, newWidth, ext)

	outputPath := filepath.Join(folderPath, newFileName)
	// Atomic move from temp to final destination
	err = os.Rename(tempFile.Name(), outputPath)
	if err != nil {
		return 0, err
	}

	// Return file size for storage tracking
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		return 0, err
	}

	return fileInfo.Size(), nil
}

// DeleteFile removes a file and its optimized versions (if it's an image).
// For image files (.jpg, .jpeg, .png), also deletes the resized variants
// created by StoreImageOptimizedVersions.
func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	// Clean up optimized image versions if this is an image file
	ext := filepath.Ext(filePath)
	if slices.Contains([]string{".jpg", ".jpeg", ".png"}, ext) {
		baseName := strings.TrimSuffix(filePath, ext)
		for _, width := range imageOptimizedWidths {
			optimizedPath := fmt.Sprintf("%s-%dx%d%s", baseName, width, width, ext)
			os.Remove(optimizedPath) // Ignore errors for optimized versions
		}
	}

	return nil
}

// letterBytes defines the character set for RandomStringBytes (letters only).
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandomStringBytes generates a random string of length n using only alphabetic characters.
func RandomStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// DeleteFoldersInDirectory removes all subdirectories (and their contents) within a directory.
// Files directly in the root directory are preserved.
func DeleteFoldersInDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		// Skip files, only process directories
		if !entry.IsDir() {
			continue
		}

		subdir := filepath.Join(dir, entry.Name())
		err = deleteAllContents(subdir)
		if err != nil {
			return fmt.Errorf("failed to delete subdirectory '%s': %w", subdir, err)
		}
	}

	return nil
}

// deleteAllContents recursively removes all files and subdirectories within a path,
// then removes the directory itself. Used as a helper for DeleteFoldersInDirectory.
func deleteAllContents(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			// Recursively delete subdirectory contents first
			err = deleteAllContents(entryPath)
			if err != nil {
				return fmt.Errorf("failed to delete subdirectory '%s': %w", entryPath, err)
			}
		} else {
			err = os.Remove(entryPath)
			if err != nil {
				return fmt.Errorf("failed to delete file '%s': %w", entryPath, err)
			}
		}
	}

	// Remove the now-empty directory
	err = os.Remove(path)
	if err != nil {
		return fmt.Errorf("failed to remove empty directory '%s': %w", path, err)
	}

	return nil
}

// PrettyJSON converts any value to a formatted JSON string with indentation.
// Useful for debugging and logging structured data.
func PrettyJSON(v interface{}) string {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Err(err).Msg("Error marshaling to JSON with utils.PrettyJSON")
	}

	return string(jsonData)
}

// BoolPtr returns a pointer to the given boolean value.
// Useful when you need *bool for optional fields or API requests.
func BoolPtr(val bool) *bool {
	return &val
}

// GenerateRandomString generates a random lowercase alphanumeric string of specified length.
// Character set: a-z and 0-9.
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// StartOfDay returns a new time.Time set to 00:00:00.000000000 of the given day.
// Preserves the original timezone.
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// StartOfDayString returns the start of day as an RFC3339 formatted string.
func StartOfDayString(t time.Time) string {
	return StartOfDay(t).Format(time.RFC3339)
}

// EndOfDay returns a new time.Time set to 23:59:59.999999999 of the given day.
// Preserves the original timezone.
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// EndOfDayString returns the end of day as an RFC3339 formatted string.
func EndOfDayString(t time.Time) string {
	return EndOfDay(t).Format(time.RFC3339)
}
