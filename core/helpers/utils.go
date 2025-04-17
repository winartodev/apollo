package helpers

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core"
	cErrors "github.com/winartodev/apollo/core/errors"
	"github.com/winartodev/apollo/modules/application"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
	"html/template"
	"net/mail"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const (
	otpChars = "1234567890"
	untitled = "untitled"
)

var (
	slugInvalidCharsRegex = regexp.MustCompile(`[^a-z0-9-]+`)
	slugDashesRegex       = regexp.MustCompile(`-+`)
)

func ReadYAMLFile(path string, out interface{}) error {
	if path == "" {
		return errorPathIsEmpty
	}

	completePath, err := GetCompletePath(path)
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

func GetUserIDFromFiberContext(ctx *fiber.Ctx) (id int64, err error) {
	if localID, ok := ctx.Locals(core.CtxUserID).(int64); ok {
		id = localID
	} else {
		return 0, errors.New("no user id")
	}

	return id, nil
}

func GetApplicationAccessFromFiberContext(ctx *fiber.Ctx) (access *application.Access, err error) {
	if access, ok := ctx.Locals(core.CtxApplicationAccess).(*application.Access); ok {
		return access, nil
	}

	return nil, errors.New("no application access")
}

func GetUserIDFromContext(ctx context.Context) (id int64, err error) {
	if localID, ok := ctx.Value(core.CtxUserID).(int64); ok {
		id = localID
	} else {
		return 0, errors.New("no user id")
	}

	return id, nil
}

func GetApplicationAccessFromContext(ctx context.Context) (access *application.Access, err error) {
	if access, ok := ctx.Value(core.CtxApplicationAccess).(*application.Access); ok {
		return access, nil
	}

	return nil, errors.New("no application access")
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

func GenerateOTP(length int) (res *string, err error) {
	if length < 6 {
		length = 6
	}

	buffer := make([]byte, length)
	_, err = rand.Read(buffer)
	if err != nil {
		return nil, err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	otp := string(buffer)

	return &otp, nil
}

func FormatDuration(duration time.Duration) string {
	if duration < time.Second {
		return fmt.Sprintf("%d nanoseconds", duration.Nanoseconds())
	}

	if duration < time.Minute {
		return fmt.Sprintf("%.0f seconds", duration.Seconds())
	}

	if duration < time.Hour {
		return fmt.Sprintf("%.0f minutes", duration.Minutes())
	}

	return fmt.Sprintf("%.0f hours", duration.Hours())
}

func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GetCompletePath(path string) (completePath string, err error) {
	completePath, err = filepath.Abs(path)
	if err != nil {
		return "", errorInvalidPath
	}

	return completePath, nil
}

func CurrentOS(os string) bool {
	return runtime.GOOS == os
}

func ParseHTMLTemplateAndExecute(path string, body *bytes.Buffer, data any) (err error) {
	tmpl := template.Must(template.ParseFiles(path))

	err = tmpl.Execute(body, data)
	if err != nil {
		return err
	}

	return nil
}

func FormatIndonesianPhoneNumber(phone string) (string, error) {
	cleaned := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)

	// Validate length (10-15 digits after country code)
	if len(cleaned) < 10 || len(cleaned) > 15 {
		return "", fmt.Errorf("invalid phone number length")
	}

	switch {
	case strings.HasPrefix(cleaned, "0"):
		return "+62" + cleaned[1:], nil
	case strings.HasPrefix(cleaned, "62"):
		return "+" + cleaned, nil
	case strings.HasPrefix(cleaned, "8"): // Sometimes numbers start with 8 directly
		return "+62" + cleaned, nil
	default:
		return "", fmt.Errorf("invalid Indonesian phone number format")
	}
}

// NormalizePhoneNumber Remove all non-digit characters including '+'
func NormalizePhoneNumber(phone string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)
}

func PrintJSON(v interface{}) {
	marshaled, _ := json.MarshalIndent(v, "", "   ")
	fmt.Println(string(marshaled))
}

func MakeSlug(text string) string {
	slug := strings.ToLower(strings.TrimSpace(text))
	slug = slugInvalidCharsRegex.ReplaceAllString(slug, "-")
	slug = slugDashesRegex.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = untitled
	}

	return slug
}

func IsValidSlug(slug string) cErrors.Errors {
	if slug == "" {
		return cErrors.InvalidSlugErr.WithReason("slug is empty")
	}

	if slug == untitled {
		return nil
	}

	if strings.HasPrefix(slug, "-") {
		return cErrors.InvalidSlugErr.WithReason("slug cannot start with '-'")
	}

	if strings.HasSuffix(slug, "-") {
		return cErrors.InvalidSlugErr.WithReason("slug cannot end with '-'")
	}

	invalidCharsRegex := regexp.MustCompile(`[^a-z0-9-]`)
	if invalidCharsRegex.MatchString(slug) {
		return cErrors.InvalidSlugErr.WithReason("slug can only contain lowercase letters, numbers, and single hyphens")
	}

	consecutiveHyphensRegex := regexp.MustCompile(`--+`)
	if consecutiveHyphensRegex.MatchString(slug) {
		return cErrors.InvalidSlugErr.WithReason("slug cannot contain consecutive hyphens")
	}

	return nil
}
