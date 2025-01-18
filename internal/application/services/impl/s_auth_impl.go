package impl

import (
	"arfdev/chat/config"
	"arfdev/chat/internal/application/dtos/requests"
	"arfdev/chat/internal/application/dtos/responses"
	"arfdev/chat/internal/application/services"
	"arfdev/chat/internal/domain/entities"
	"arfdev/chat/internal/domain/repositories"
	"arfdev/chat/pkg/helpers"
	"bytes"
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

type authServiceImpl struct {
	repo repositories.UserRepositoryInterface
	rdb  *redis.Client
}

func NewAuthServiceImpl(repo repositories.UserRepositoryInterface, rdb *redis.Client) services.AuthServiceInterface {
	return &authServiceImpl{
		repo: repo,
		rdb:  rdb,
	}
}

func (s *authServiceImpl) CheckEmail(ctx context.Context, payload requests.CheckEmailRequestPayload) responses.APIBaseResponse {
	if ctx.Err() != nil {
		return responses.NewResponse(http.StatusBadRequest, false, "Bad Request", nil)
	}

	user, err := s.repo.FindByEmail(ctx, payload.Email)
	if !payload.IsLogin {
		if user.Email == "" {
			return responses.NewResponse(http.StatusOK, true, "Email can be used for register", nil)
		}
		if err != nil {
			return responses.NewResponse(http.StatusInternalServerError, false, err.Error(), nil)
		}
		return responses.NewResponse(http.StatusBadRequest, false, "Email already exists", nil)
	}

	if user.Email != "" {
		response := responses.CheckPhoneIsLoginResponsePayload{
			DeviceID: user.DeviceID,
		}
		return responses.NewResponse(http.StatusOK, true, "User exists", response)
	}

	return responses.NewResponse(http.StatusBadRequest, false, "User does not exist", nil)
}

func (s *authServiceImpl) SendOtp(ctx context.Context, payload requests.OTPRequestPayload) responses.APIBaseResponse {
	if ctx.Err() != nil {
		return responses.NewResponse(http.StatusBadRequest, false, "Bad Request", nil)
	}

	signatureID, err := helpers.GenerateUniqueID()
	otp, err := helpers.GenerateOTP()
	if err != nil {
		return responses.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
	}

	data := &entities.Mst_otp{
		OTPCode:     otp,
		DeviceID:    payload.DeviceID,
		SignatureID: signatureID,
		ExpiredAt:   time.Now().Add(3 * time.Minute),
	}

	err = s.repo.Insert(ctx, *data)
	if err != nil {
		return responses.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err = helpers.SendEmailOTP(ctx, otp, payload.Email)
	if err != nil {
		switch {
		case ctx.Err() == context.DeadlineExceeded:
			return responses.NewResponse(http.StatusGatewayTimeout, false, "Request timeout: OTP email sending failed", nil)
		case strings.Contains(err.Error(), "authentication failed"):
			log.Printf("SMTP authentication error: %v", err)
			return responses.NewResponse(http.StatusInternalServerError, false, "Failed to send OTP email", nil)
		case strings.Contains(err.Error(), "connection failed"):
			log.Printf("SMTP connection error: %v", err)
			return responses.NewResponse(http.StatusServiceUnavailable, false, "Email service temporarily unavailable", nil)
		default:
			log.Printf("Email sending error: %v", err)
			return responses.NewResponse(http.StatusInternalServerError, false, "Failed to send OTP email", nil)
		}
	}

	response := responses.OTPResponsePayload{
		SignatureID: signatureID,
	}

	return responses.NewResponse(http.StatusOK, true, "Success request OTP", response)
}

func (s *authServiceImpl) VerifyOTP(ctx context.Context, payload requests.ValidateOTPRequestPayload) responses.APIBaseResponse {
	if ctx.Err() != nil {
		return responses.NewResponse(http.StatusBadRequest, false, "Bad Request", nil)
	}

	find, err := s.repo.FindOTP(ctx, payload.OTPCode.(string), payload.SignatureID)
	if err != nil {
		if err.Error() == "Invalid Signature ID or OTP" {
			return responses.NewResponse(http.StatusBadRequest, false, err.Error(), nil)
		}
		return responses.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
	}

	now := time.Now()
	if now.After(find.ExpiredAt) {
		return responses.NewResponse(http.StatusBadRequest, false, "OTP Expired", nil)
	}
	return responses.NewResponse(http.StatusOK, true, "Verification success", nil)
}

func (s *authServiceImpl) Register(ctx context.Context, payload requests.RegisterRequestPayload) responses.APIBaseResponse {
	bucketname := os.Getenv("MINIO_BUCKETNAME")
	minioClient := config.NewMinioClient()
	userId := uuid.New().String()
	if ctx.Err() != nil {
		return responses.NewResponse(http.StatusBadRequest, false, "Bad Request", nil)
	}

	pubKeyFile, err := payload.PublicKey.Open()
	if err != nil {
		return responses.NewResponse(http.StatusBadRequest, false, "Failed to process pub key", nil)
	}

	pubKeyBytes, err := io.ReadAll(pubKeyFile)
	if err != nil {
		return responses.NewResponse(http.StatusBadRequest, false, "Failed to read pub key", nil)
	}

	block, _ := pem.Decode(pubKeyBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return responses.NewResponse(http.StatusBadRequest, false, "Invalid Pub Key Format.", nil)
	}

	_, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return responses.NewResponse(http.StatusBadRequest, false, "Invalid Pub Key", nil)
	}

	profilePicPath := "default-profile-icon-24.jpg"
	if payload.ImgFile != nil {
		profileImg, err := payload.ImgFile.Open()
		if err != nil {
			return responses.NewResponse(http.StatusBadRequest, false, "Failed to process profile picture", nil)
		}

		contentType := payload.ImgFile.Header.Get("Content-Type")
		if !helpers.ValidFileType(contentType) {
			return responses.NewResponse(http.StatusBadRequest, false, "Invalid Image Format", nil)
		}

		profilePicPath = fmt.Sprintf("profilePic/%s%s", userId, filepath.Ext(payload.ImgFile.Filename))
		_, err = minioClient.PutObject(ctx, bucketname, profilePicPath, profileImg, payload.ImgFile.Size, minio.PutObjectOptions{
			ContentType: payload.ImgFile.Header.Get("Content-Type"),
		})
		if err != nil {
			return responses.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
		}
		defer profileImg.Close()
	}

	pubKeyPath := fmt.Sprintf("pubKey/%s.pem", userId)
	_, err = minioClient.PutObject(ctx, bucketname, pubKeyPath, bytes.NewReader(pubKeyBytes), payload.PublicKey.Size, minio.PutObjectOptions{
		ContentType: payload.PublicKey.Header.Get("Content-Type"),
	})
	if err != nil {
		return responses.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
	}

	usersData := entities.Mst_users{
		ID:            userId,
		Email:         payload.Email,
		DeviceID:      payload.DeviceID,
		PublicKeyPath: pubKeyPath,
		IsOnline:      false,
		IsDeleted:     false,
	}

	usersDetail := entities.Mst_users_detail{
		UserID:   usersData.ID,
		Username: payload.Username,
		FullName: payload.Fullname,
		ImgPath:  profilePicPath,
		Gender:   payload.Gender,
		Age:      int8(payload.Age),
	}

	err = s.repo.InsertUser(ctx, usersData)
	if err != nil {
		return responses.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
	}
	err = s.repo.InsertUserDetail(ctx, usersDetail)
	if err != nil {
		return responses.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
	}

	return responses.NewResponse(http.StatusCreated, true, "Register Success", nil)
}

func (s *authServiceImpl) Login(ctx context.Context, payload requests.LoginRequestPayload) responses.APIBaseResponse {
	secret := []byte(os.Getenv("JWT_SECRET"))
	accessExp, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRES_IN"))
	refreshExp, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRES_IN"))
	accessExpDuration := time.Duration(accessExp) * time.Hour
	refreshExpDuration := time.Duration(refreshExp) * time.Hour

	jwtAccessHelper := helpers.NewJWTHelper(secret, accessExpDuration)
	jwtRefreshHelper := helpers.NewJWTHelper(secret, refreshExpDuration)

	if ctx.Err() != nil {
		return responses.NewResponse(http.StatusBadRequest, false, "Bad Request", nil)
	}

	find, err := s.repo.FindByEmail(ctx, payload.Email)
	if err == nil && "" == find.Email {
		return responses.NewResponse(http.StatusBadRequest, false, "User has not been registered yet", nil)
	}

	if err != nil {
		return responses.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
	}

	accessToken, errToken := jwtAccessHelper.GenerateToken(find.ID)
	refreshToken, errToken := jwtRefreshHelper.GenerateToken(find.ID)
	if errToken != nil {
		return responses.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
	}

	accessRdbKey := fmt.Sprintf("access-%s", find.ID)
	refreshRdbKey := fmt.Sprintf("refresh-%s", find.ID)

	errRdb := s.rdb.Set(ctx, accessRdbKey, accessToken, accessExpDuration).Err()
	errRdb = s.rdb.Set(ctx, refreshRdbKey, refreshToken, refreshExpDuration).Err()
	if errRdb != nil {
		return responses.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
	}

	response := &responses.LoginResponsePayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return responses.NewResponse(http.StatusOK, true, "Login success", response)
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, payload requests.RefreshTokenRequestPayload) responses.APIBaseResponse {
	secret := []byte(os.Getenv("JWT_SECRET"))
	accessExp, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRES_IN"))
	refreshExp, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRES_IN"))
	accessExpDuration := time.Duration(accessExp) * time.Hour
	refreshExpDuration := time.Duration(refreshExp) * time.Hour

	jwtAccessHelper := helpers.NewJWTHelper(secret, accessExpDuration)
	jwtRefreshHelper := helpers.NewJWTHelper(secret, refreshExpDuration)

	if ctx.Err() != nil {
		return responses.NewResponse(http.StatusBadRequest, false, "Bad Request", nil)
	}

	claims, err := jwtRefreshHelper.ValidateToken(payload.RefreshToken)
	if err != nil {
		return responses.NewResponse(http.StatusUnauthorized, false, "Unautorized", nil)
	}

	find, err := s.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return responses.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
	}

	accessToken, errToken := jwtAccessHelper.GenerateToken(find.ID)
	if errToken != nil {
		return responses.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
	}

	response := responses.RefreshTokenResponsePayload{
		AccessToken: accessToken,
	}

	errRdb := s.rdb.Set(ctx, fmt.Sprintf("access-%s", find.ID), accessToken, accessExpDuration).Err()
	if errRdb != nil {
		return responses.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
	}

	return responses.NewResponse(http.StatusOK, true, "Success get new access token", response)
}

func (s *authServiceImpl) Logout(ctx context.Context, userID string) responses.APIBaseResponse {
	if ctx.Err() != nil {
		return responses.NewResponse(http.StatusBadRequest, false, "Bad Request", nil)
	}

	accessRdbKey := fmt.Sprintf("access-%s", userID)
	refreshRdbKey := fmt.Sprintf("refresh-%s", userID)

	errDel := s.rdb.Del(ctx, accessRdbKey).Err()
	errDel = s.rdb.Del(ctx, refreshRdbKey).Err()
	if errDel != nil {
		return responses.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
	}

	return responses.NewResponse(http.StatusOK, true, "Logout Success", nil)
}
