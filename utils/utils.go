package utils

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Log(variable any) {
	log.Debug().
		Interface("", variable).Msg("LOGGER ")
}

func IsEnabled(key bool) func(c *fiber.Ctx) bool {
	if key {
		return nil
	}

	return func(c *fiber.Ctx) bool { return true }
}

func InlineCondition(condition bool, ifResult any, elseResult any) any {
	if condition {
		return ifResult
	}
	return elseResult
}

func GetIntInParams(c *fiber.Ctx, key string) (uint64, error) {
	return strconv.ParseUint(c.Params(key), 10, 64)
}

func GetIntInQueries(c *fiber.Ctx, key string) (uint64, error) {
	return strconv.ParseUint(c.Query(key), 10, 64)
}

func GetAuthenticatedUser(c *fiber.Ctx) (schema.User, error) {
	user, ok := c.Locals("user").(schema.User)
	if ok {
		return user, nil
	}
	return schema.User{}, errors.New("user doesn't exist")
}

func ValidateMobileNumber(number string) error {
	// Define the regular expression pattern for a mobile number
	pattern := `9(1[0-9]|3[1-9]|2[1-9])-?[0-9]{3}-?[0-9]{4}`

	// Compile the regex pattern
	regex := regexp.MustCompile(pattern)

	// Match the number against the regex pattern
	valid := regex.MatchString(number)
	if !valid {
		return errors.New("تلفن همراه معتبر نمی باشد")
	}

	return nil
}

func Random(min int, max int) uint64 {
	return uint64(rand.Intn(max-min) + min) //nolint:gosec
}

func RandomFromArray(arr []uint64) uint64 {
	randomIndex := rand.Intn(len(arr))
	return arr[randomIndex]
}

func GetFakeTableIDs(db *gorm.DB, table any) (ids []uint64, err error) {
	err = db.Model(&table).Select("id").Find(&ids).Error
	if err != nil {
		log.Err(err)
	}
	return
}

func RandomDateTime() time.Time {
	year := rand.Intn(2023) //nolint:gosec
	month := time.Month(rand.Intn(12) + 1)
	day := rand.Intn(31) + 1 //nolint:gosec
	hour := rand.Intn(24)    //nolint:gosec
	min := rand.Intn(60)     //nolint:gosec
	sec := rand.Intn(60)     //nolint:gosec

	return time.Date(year, month, day, hour, min, sec, 0, time.UTC)
}

var imageOptimizedWidths = []int{600, 300}

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

func resizeAndSaveImage(folderPath string, imageName string, newWidth int) (size int64, err error) {
	tempFile, err := ioutil.TempFile("", "resize-*.jpg")
	if err != nil {
		return 0, err
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name()) // Remove the temporary file after processing

	inputPath := filepath.Join(folderPath, imageName)
	img, err := imaging.Open(inputPath)
	if err != nil {
		return 0, err
	}

	if newWidth > img.Bounds().Size().X {
		return 0, nil
	}

	resizedImg := imaging.Resize(img, newWidth, 0, imaging.NearestNeighbor)

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

	// Remove the extension from the file name
	baseName := strings.TrimSuffix(imageName, ext)

	// Create the new file name
	newFileName := fmt.Sprintf("%s-%dx%d%s", baseName, newWidth, newWidth, ext)

	outputPath := filepath.Join(folderPath, newFileName)
	// Move the temporary file to the desired output path
	err = os.Rename(tempFile.Name(), outputPath)
	if err != nil {
		return 0, err
	}

	// Get the file size of the resized image
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		return 0, err
	}

	fileSize := fileInfo.Size()

	return fileSize, nil
}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	ext := filepath.Ext(filePath)
	if slices.Contains([]string{".jpg", ".jpeg", ".png"}, ext) {
		baseName := strings.TrimSuffix(filePath, ext)
		for _, width := range imageOptimizedWidths {
			filePath := fmt.Sprintf("%s-%dx%d%s", baseName, width, width, ext)
			os.Remove(filePath)
		}
	}

	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func DeleteFoldersInDirectory(dir string) error {
	// Read the contents of the directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// Iterate through the entries
	for _, entry := range entries {
		// Skip non-directories
		if !entry.IsDir() {
			continue
		}

		// Build the path to the subdirectory
		subdir := filepath.Join(dir, entry.Name())

		// Remove the subdirectory and its contents (files and subdirectories)
		err = deleteAllContents(subdir)
		if err != nil {
			return fmt.Errorf("failed to delete subdirectory '%s': %w", subdir, err)
		}
	}

	return nil
}

func deleteAllContents(path string) error {
	// Read the contents of the directory
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// Iterate through the entries and remove them
	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			// Recursively delete the subdirectory and its contents
			err = deleteAllContents(entryPath)
			if err != nil {
				return fmt.Errorf("failed to delete subdirectory '%s': %w", entryPath, err)
			}
		} else {
			// Remove the file
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
