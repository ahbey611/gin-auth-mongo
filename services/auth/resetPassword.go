package auth

import (
	"encoding/base64"
	"errors"
	"os"
	"strings"
	"time"

	"gin-auth-mongo/databases"
	"gin-auth-mongo/models/requests"
	"gin-auth-mongo/repositories"
	"gin-auth-mongo/utils/consts"
	"gin-auth-mongo/utils/datetime"

	"gin-auth-mongo/utils/crypto"
	"gin-auth-mongo/utils/mail"
)

func GenerateResetPasswordFlowID(email string) (string, string, error) {

	expiredAt := time.Now().Add(time.Duration(consts.VERIFY_EMAIL_RESET_PWD_LINK_EXPIRY) * time.Minute).Format(consts.DATETIME_FORMAT)
	code := email + "$" + expiredAt

	// base64 encode
	flowID := base64.RawURLEncoding.EncodeToString([]byte(code))
	// log.Println("resetCode(after encode): ", encodedCode)

	return flowID, expiredAt, nil
}

func DecodeResetPasswordFlowID(flowID string) (string, string, error) {

	decodedCode, err := base64.RawURLEncoding.DecodeString(flowID)
	if err != nil {
		return "", "", err
	}

	splitted := strings.Split(string(decodedCode), "$")
	if len(splitted) != 2 {
		return "", "", errors.New("invalid code")
	}
	email := splitted[0]
	expiredAt := splitted[1]

	timeExpired, err := time.Parse(consts.DATETIME_FORMAT, expiredAt)
	if err != nil {
		return "", "", errors.New("invalid code")
	}

	if time.Now().After(timeExpired) {
		return "", "", errors.New("code expired")
	}

	return email, expiredAt, nil
}

func UserEmailResetPasswordWithLink(request *requests.EmailPasswordResetLinkRequest) error {

	// check if the user exists
	user, err := repositories.GetUserByEmail(request.Email)

	if err != nil || user == nil {
		return errors.New("invalid email")
	}

	// generate reset password code
	flowID, expiredAt, err := GenerateResetPasswordFlowID(request.Email)
	if err != nil {
		return errors.New("try again later")
	}

	// generate link
	link := os.Getenv("FRONTEND_URL") + consts.FRONTEND_RESET_PASSWORD_ROUTE + "?flow_id=" + flowID

	err = mail.SendVerificationEmail(user.Email, user.Username, mail.VerificationFormTypeLink, mail.VerificationRequestTypeResetPassword, link, expiredAt).Error
	if err != nil {
		return errors.New("try again later")
	}

	// store the code in redis
	databases.RedisSet(consts.VERIFY_EMAIL_RESET_PWD_FLOW_ID+flowID, "1", consts.VERIFY_EMAIL_RESET_PWD_CODE_EXPIRY, datetime.MINUTES)

	return nil
}

func UserEmailResetPasswordWithCode(request *requests.EmailPasswordResetCodeRequest) error {
	// check if the user exists
	user, err := repositories.GetUserByEmail(request.Email)

	if err != nil || user == nil {
		return errors.New("invalid email")
	}

	hashedPassword, err := crypto.HashPassword(request.Password)
	if err != nil {
		return errors.New("please try again later")
	}

	verficationCode, err := GenerateVerificationCode()

	if err != nil {
		return errors.New("please try again later")
	}

	expiredAt := time.Now().Add(time.Duration(consts.VERIFY_EMAIL_RESET_PWD_CODE_EXPIRY) * time.Minute).Format(consts.DATETIME_FORMAT)

	err = mail.SendVerificationEmail(user.Email, user.Username, mail.VerificationFormTypeCode, mail.VerificationRequestTypeResetPassword, verficationCode, expiredAt).Error
	if err != nil {
		return errors.New("please try again later")
	}

	databases.RedisSet(consts.VERIFY_EMAIL_RESET_PWD_CODE+request.Email, verficationCode, consts.VERIFY_EMAIL_RESET_PWD_CODE_EXPIRY, datetime.MINUTES)
	databases.RedisSet(consts.VERIFY_EMAIL_RESET_PWD_PASSWORD+request.Email, hashedPassword, consts.VERIFY_EMAIL_RESET_PWD_CODE_EXPIRY, datetime.MINUTES)

	return nil
}

func UserEmailResetPasswordLinkVerify(request *requests.EmailPasswordResetLinkVerifyRequest) error {

	exists, err := databases.RedisGet(consts.VERIFY_EMAIL_RESET_PWD_FLOW_ID + request.FlowId)

	// Maybe the email is not in the redis
	if err != nil || exists != "1" {
		return errors.New("invalid code")
	}

	// decode the reset code
	email, _, err := DecodeResetPasswordFlowID(request.FlowId)
	if err != nil {
		return err
	}

	// hash the password
	password := request.Password
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return err
	}

	user, err := repositories.GetUserByEmail(email)
	if err != nil || user == nil {
		return errors.New("invalid email")
	}

	// update the password
	err = repositories.UpdateUserPasswordByID(user.ID, hashedPassword)
	if err != nil {
		return errors.New("update password failed")
	}

	// delete the flow id from redis
	databases.RedisDel(consts.VERIFY_EMAIL_RESET_PWD_FLOW_ID + request.FlowId)

	return nil
}

func UserEmailResetPasswordCodeVerify(request *requests.EmailPasswordResetCodeVerifyRequest) error {

	// ctx := databases.GetRedisContext()
	// code, err := databases.RedisClient.Get(ctx, consts.VERIFY_EMAIL_RESET_PWD_CODE+request.Email).Result()
	code, err := databases.RedisGet(consts.VERIFY_EMAIL_RESET_PWD_CODE + request.Email)
	if err != nil || len(code) != 6 {
		return errors.New("invalid code")
	}

	if code != request.Code {
		return errors.New("invalid code")
	}

	// hashedPassword, err := databases.RedisClient.Get(ctx, consts.VERIFY_EMAIL_RESET_PWD_PASSWORD+request.Email).Result()
	hashedPassword, err := databases.RedisGet(consts.VERIFY_EMAIL_RESET_PWD_PASSWORD + request.Email)
	if err != nil || len(hashedPassword) == 0 {
		return errors.New("please try to reset password again")
	}

	// check if the user exists
	user, err := repositories.GetUserByEmail(request.Email)
	if err != nil || user == nil {
		return errors.New("invalid email")
	}

	// update the password
	err = repositories.UpdateUserPasswordByID(user.ID, hashedPassword)
	if err != nil {
		return errors.New("update password failed")
	}

	// delete the code and password from redis
	databases.RedisDel(consts.VERIFY_EMAIL_RESET_PWD_CODE + request.Email)
	databases.RedisDel(consts.VERIFY_EMAIL_RESET_PWD_PASSWORD + request.Email)

	return nil
}

// check the reset password code expired status
func CheckUserEmailResetPasswordLinkExpired(resetCode string) (map[string]string, error) {

	exists, err := databases.RedisGet(consts.VERIFY_EMAIL_RESET_PWD_FLOW_ID + resetCode)

	// Maybe the email is not in the redis
	if err != nil || exists != "1" {
		return nil, errors.New("invalid code")
	}

	email, expiredAt, err := DecodeResetPasswordFlowID(resetCode)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"email":     email,
		"expiredAt": expiredAt,
	}, nil
}
