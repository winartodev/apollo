package controllers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/winartodev/apollo/core/configs"
	"github.com/winartodev/apollo/core/helpers"
	authEnum "github.com/winartodev/apollo/modules/auth/emums"
	authEntity "github.com/winartodev/apollo/modules/auth/entities"
	authRepo "github.com/winartodev/apollo/modules/auth/repositories"
	"time"
)

const (
	defaultOTPLength = 6
	maxRefreshOTP    = 3

	otpEmailExpiration = 60 * time.Second
	otpPhoneExpiration = 15 * time.Minute
	defaultTTL         = 30 * time.Minute

	otpMailHtmlTemplate = "modules/auth/files/otp-mail-template.html"
	phoneMessageFormat  = "[%s] Your verification code is %s, valid for (%s)"
)

var (
	ErrorOTPAlreadyVerified      = errors.New("otp already verified")
	errorInvalidEmail            = errors.New("invalid email")
	errorInvalidVerificationType = errors.New("invalid verification type")
	errorOTPAlreadyExists        = errors.New("otp already exists")
	ErrorOTPDataEmpty            = errors.New("OTP not found")
	errorOTPDataExpired          = errors.New("OTP is expired")
	errorOTPNotMatch             = errors.New("otp code not match")
	errorOTPMaxAttempts          = errors.New("OTP max attempts exceeded")
)

type OTPMailTemplate struct {
	RecipientName string
	OTPCode       string
	Duration      string
}

type VerificationControllerItf interface {
	GetOTP(ctx context.Context, verificationType int, value string) (data *authEntity.OTPData, err error)
	CreateOTP(ctx context.Context, verificationType int, value string) (err error)
	VerifyOTP(ctx context.Context, verificationType int, value string, code string) (err error)
	ResendOTP(ctx context.Context, verificationType int, value string) (err error)
	DeleteOTP(ctx context.Context, verificationType int, value string) (err error)
}

type VerificationController struct {
	SmtpClient             *configs.SMTPClient
	TwilioClient           *configs.TwilioClient
	VerificationRepository authRepo.VerificationRepositoryItf
}

func NewVerificationController(controller VerificationController) VerificationControllerItf {
	return &VerificationController{
		SmtpClient:             controller.SmtpClient,
		TwilioClient:           controller.TwilioClient,
		VerificationRepository: controller.VerificationRepository,
	}
}

func (vc *VerificationController) GenerateAndStoreOTP(ctx context.Context, verificationType int, value string) (err error) {
	otp, err := helpers.GenerateOTP(defaultOTPLength)
	if err != nil {
		return err
	}

	data := authEntity.OTPData{
		Value:      value,
		OTP:        *otp,
		IsVerified: false,
	}

	switch verificationType {
	case authEnum.VerificationEmail:
		return vc.handleSendEmailOTP(ctx, data)
	case authEnum.VerificationPhone:
		return vc.handleSendPhoneOTP(ctx, data)
	default:
		return nil
	}
}

func (vc *VerificationController) handleSendEmailOTP(ctx context.Context, data authEntity.OTPData) (err error) {
	if !helpers.IsEmailValid(data.Value) {
		return errorInvalidEmail
	}

	data.Expire = time.Now().Add(otpEmailExpiration).Unix()

	ttl := defaultTTL

	err = vc.VerificationRepository.SetEmailOTPRedis(ctx, data.Value, data, &ttl)
	if err != nil {
		return err
	}

	otpMailTemplate := OTPMailTemplate{
		RecipientName: data.Value,
		OTPCode:       data.OTP,
		Duration:      helpers.FormatDuration(otpEmailExpiration),
	}

	go func(templateData OTPMailTemplate) {
		err = vc.SendOTPToEmail(templateData)
		if err != nil {
			log.Errorf("SendOTPToEmail err: %v", err)
		}
	}(otpMailTemplate)

	return nil
}

func (vc *VerificationController) handleSendPhoneOTP(ctx context.Context, data authEntity.OTPData) (err error) {
	ttl := defaultTTL

	data.Expire = time.Now().Add(otpPhoneExpiration).Unix()

	err = vc.VerificationRepository.SetPhoneOTPRedis(ctx, data.Value, data, &ttl)
	if err != nil {
		return err
	}

	expireStr := helpers.FormatDuration(otpPhoneExpiration)
	message := fmt.Sprintf(phoneMessageFormat, "APOLLO", data.OTP, expireStr)

	go func() {
		err = vc.TwilioClient.SendSMS(data.Value, message)
		if err != nil {
			log.Errorf("SendSMS err: %v", err)
		}
	}()

	return nil
}

func (vc *VerificationController) GetOTP(ctx context.Context, verificationType int, value string) (data *authEntity.OTPData, err error) {
	switch verificationType {
	case authEnum.VerificationEmail:
		return vc.VerificationRepository.GetEmailOTPRedis(ctx, value)
	case authEnum.VerificationPhone:
		return vc.VerificationRepository.GetPhoneOTPRedis(ctx, value)
	default:
		return nil, errorInvalidVerificationType
	}
}

func (vc *VerificationController) CreateOTP(ctx context.Context, verificationType int, value string) (err error) {
	data, err := vc.GetOTP(ctx, verificationType, value)
	if err != nil {
		return err
	}

	if data != nil {
		return errorOTPAlreadyExists
	}

	return vc.GenerateAndStoreOTP(ctx, verificationType, value)
}

func (vc *VerificationController) VerifyOTP(ctx context.Context, verificationType int, value string, code string) (err error) {
	data, err := vc.GetOTP(ctx, verificationType, value)
	if err != nil {
		return err
	}

	if data == nil {
		return ErrorOTPDataEmpty
	}

	if data.IsVerified {
		return ErrorOTPAlreadyVerified
	}

	expirationTime := time.Unix(data.Expire, 0)
	if time.Now().After(expirationTime) {
		return errorOTPDataExpired
	}

	if data.OTP != code {
		return errorOTPNotMatch
	}

	data.IsVerified = true

	ttl := defaultTTL

	switch verificationType {
	case authEnum.VerificationEmail:
		return vc.VerificationRepository.SetEmailOTPRedis(ctx, value, *data, &ttl)
	case authEnum.VerificationPhone:
		return vc.VerificationRepository.SetPhoneOTPRedis(ctx, value, *data, &ttl)
	}

	return nil
}

func (vc *VerificationController) ResendOTP(ctx context.Context, verificationType int, value string) (err error) {
	data, err := vc.GetOTP(ctx, verificationType, value)
	if err != nil {
		return err
	}

	if data == nil {
		return ErrorOTPDataEmpty
	}

	if data.IsVerified {
		return ErrorOTPAlreadyVerified
	}

	ttl := defaultTTL

	attempt, err := vc.VerificationRepository.SetResendAttemptRedis(ctx, value, &ttl)
	if err != nil {
		return err
	}

	if attempt > maxRefreshOTP {
		return errorOTPMaxAttempts
	}

	return vc.GenerateAndStoreOTP(ctx, verificationType, value)
}

func (vc *VerificationController) SendOTPToEmail(data OTPMailTemplate) error {
	completePath, err := helpers.GetCompletePath(otpMailHtmlTemplate)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	err = helpers.ParseHTMLTemplateAndExecute(completePath, &body, data)
	if err != nil {
		return err
	}

	emailData := &configs.Email{
		To:      data.RecipientName,
		Subject: "Email OTP Verification Code",
		Body:    body.String(),
		HTML:    true,
	}

	err = vc.SmtpClient.Send(emailData)
	if err != nil {
		return err
	}

	return nil
}

func (vc *VerificationController) DeleteOTP(ctx context.Context, verificationType int, value string) (err error) {
	data, err := vc.GetOTP(ctx, verificationType, value)
	if err != nil {
		return err
	}

	if data == nil {
		return ErrorOTPDataEmpty
	}

	switch verificationType {
	case authEnum.VerificationEmail:
		return vc.VerificationRepository.DeleteEmailOTPRedis(ctx, value)
	case authEnum.VerificationPhone:
		return vc.VerificationRepository.DeletePhoneOTPRedis(ctx, value)
	default:
		return errorInvalidVerificationType
	}
}
