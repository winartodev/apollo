package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/winartodev/apollo/core/helpers"
	authEntity "github.com/winartodev/apollo/modules/auth/entities"
	"time"
)

const (
	emailOTPPrefix = "otp_email"
	phoneOTPPrefix = "otp_phone"
	resendAttempt  = "resend_attempt"
)

type VerificationRepositoryItf interface {
	SetEmailOTPRedis(ctx context.Context, email string, data authEntity.OTPData, ttl *time.Duration) (err error)
	SetPhoneOTPRedis(ctx context.Context, phoneNumber string, data authEntity.OTPData, ttl *time.Duration) (err error)
	GetEmailOTPRedis(ctx context.Context, email string) (res *authEntity.OTPData, err error)
	GetPhoneOTPRedis(ctx context.Context, phoneNumber string) (res *authEntity.OTPData, err error)
	SetResendAttemptRedis(ctx context.Context, value string, ttl *time.Duration) (count int64, err error)
	DeletePhoneOTPRedis(ctx context.Context, phoneNumber string) (err error)
	DeleteEmailOTPRedis(ctx context.Context, email string) (err error)
}

type VerificationRepository struct {
	Redis *redis.Client
}

func NewVerificationRepository(repository VerificationRepository) VerificationRepositoryItf {
	return &VerificationRepository{
		Redis: repository.Redis,
	}
}

func (vr *VerificationRepository) SetEmailOTPRedis(ctx context.Context, email string, data authEntity.OTPData, ttl *time.Duration) (err error) {
	key := vr.GenerateRedisKey(emailOTPPrefix, email)
	return vr.setOTPRedis(ctx, key, data, ttl)
}

func (vr *VerificationRepository) SetPhoneOTPRedis(ctx context.Context, phoneNumber string, data authEntity.OTPData, ttl *time.Duration) (err error) {
	phoneNumberStr := helpers.NormalizePhoneNumber(phoneNumber)
	key := vr.GenerateRedisKey(phoneOTPPrefix, phoneNumberStr)
	return vr.setOTPRedis(ctx, key, data, ttl)
}

func (vr *VerificationRepository) GetEmailOTPRedis(ctx context.Context, email string) (res *authEntity.OTPData, err error) {
	key := vr.GenerateRedisKey(emailOTPPrefix, email)

	var data authEntity.OTPData
	err = vr.getRedisKey(ctx, key, &data)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if err == redis.Nil {
		return nil, nil
	}

	return &data, nil
}

func (vr *VerificationRepository) GetPhoneOTPRedis(ctx context.Context, phoneNumber string) (res *authEntity.OTPData, err error) {
	phoneNumberStr := helpers.NormalizePhoneNumber(phoneNumber)
	key := vr.GenerateRedisKey(phoneOTPPrefix, phoneNumberStr)

	var data authEntity.OTPData
	err = vr.getRedisKey(ctx, key, &data)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if err == redis.Nil {
		return nil, nil
	}

	return &data, nil
}

func (vr *VerificationRepository) SetResendAttemptRedis(ctx context.Context, value string, ttl *time.Duration) (count int64, err error) {
	key := vr.GenerateRedisKey(resendAttempt, value)
	count, err = vr.Redis.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if ttl != nil {
		err = vr.Redis.Expire(ctx, key, *ttl).Err()
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}

func (vr *VerificationRepository) DeletePhoneOTPRedis(ctx context.Context, phoneNumber string) (err error) {
	phoneNumberStr := helpers.NormalizePhoneNumber(phoneNumber)
	key := vr.GenerateRedisKey(phoneOTPPrefix, phoneNumberStr)
	return vr.deleteRedisKey(ctx, key)
}

func (vr *VerificationRepository) DeleteEmailOTPRedis(ctx context.Context, email string) (err error) {
	key := vr.GenerateRedisKey(emailOTPPrefix, email)
	return vr.deleteRedisKey(ctx, key)
}

func (vr *VerificationRepository) GenerateRedisKey(prefix string, value string) (key string) {
	return fmt.Sprintf("%s:%s", prefix, value)
}

func (vr *VerificationRepository) deleteRedisKey(ctx context.Context, key string) (err error) {
	err = vr.Redis.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

func (vr *VerificationRepository) getRedisKey(ctx context.Context, key string, data *authEntity.OTPData) (err error) {
	dataStr, err := vr.Redis.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		return err
	}

	return nil
}

func (vr *VerificationRepository) setOTPRedis(ctx context.Context, key string, data authEntity.OTPData, ttl *time.Duration) (err error) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = vr.Redis.SetEX(ctx, key, dataByte, *ttl).Err()
	if err != nil {
		return err
	}

	return nil
}
