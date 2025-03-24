package helpers

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ReadYAMLFile(path string, out interface{}) error {
	if path == "" {
		return errorPathIsEmpty
	}

	completePath, err := filepath.Abs(path)
	if err != nil {
		return errorInvalidPath
	}

	yamlFile, err := os.Open(completePath)
	if err != nil {
		return errors.New(fmt.Sprintf(errorReadYamlFile, err.Error()))
	}

	defer yamlFile.Close()

	if yamlFile == nil {
		return errorYamlFileIsEmpty
	}

	decoder := yaml.NewDecoder(yamlFile)
	err = decoder.Decode(out)
	if err != nil {
		return errors.New(fmt.Sprintf(errorDecodeYamlFile, err.Error()))
	}

	return nil
}

func GetFormValue(ctx *fiber.Ctx, key string, required bool) (value string, error error) {
	if key == "" {
		return "", errors.New("key is empty")
	}

	value = ctx.FormValue(key)

	if value == "" && required {
		formattedKey := cases.Title(language.Und).
			String(strings.ReplaceAll(key, "_", " "))

		return "", fmt.Errorf("%s is required", formattedKey)
	}

	return value, nil
}

func GetUserIDFromContext(ctx *fiber.Ctx) (id int64, err error) {
	if idStr, ok := ctx.Locals("id").(float64); ok {
		id = int64(idStr)
	} else {
		return 0, errors.New("no user id")
	}

	return id, nil
}

func FormatUnixTime(unixTime int64) *time.Time {
	if unixTime == 0 {
		return nil
	}

	t := time.Unix(unixTime, 0)
	return &t
}

func HashPassword(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
