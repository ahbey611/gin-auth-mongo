package auth

import (
	"encoding/base64"
	"errors"

	// "log"

	// "log"
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

// generate a encoded email and username and expiration time for completion of registration
func GenerateCompletionRegistrationFlowID(email, username string) (string, string, error) {

	expiredAt := time.Now().Add(time.Duration(consts.VERIFY_EMAIL_REGISTER_LINK_EXPIRY) * time.Minute).Format(consts.DATETIME_FORMAT)
	code := email + "$" + username + "$" + expiredAt

	// base64 encode
	flowID := base64.RawURLEncoding.EncodeToString([]byte(code))
	// log.Println("regCode(after encode): ", encodedCode)
	return flowID, expiredAt, nil
}

// decode the completion registration code
func DecodeCompletionRegistrationFlowID(flowID string) (string, string, string, error) {

	decodedCode, err := base64.RawURLEncoding.DecodeString(flowID)
	if err != nil {
		return "", "", "", err
	}
	splitted := strings.Split(string(decodedCode), "$")
	if len(splitted) != 3 {
		return "", "", "", errors.New("invalid code")
	}
	email := splitted[0]
	username := splitted[1]
	expiredAt := splitted[2]

	timeExpired, err := time.Parse(consts.DATETIME_FORMAT, expiredAt)
	if err != nil {
		return "", "", "", errors.New("invalid code")
	}

	if time.Now().After(timeExpired) {
		return "", "", "", errors.New("code expired")
	}

	return email, username, expiredAt, nil
}

func UserEmailRegisterWithLink(request *requests.EmailRegisterLinkRequest) error {

	// databases.RedisSet("test04", "haha", consts.VERIFY_EMAIL_REGISTER_CODE_EXPIRY, datetime.MINUTES)
	// databases.RedisSetWithoutExpiry("test03", "haha")

	// check if the email or username is already registered
	exists, err := databases.RedisExists(consts.VERIFY_EMAIL_REGISTER_USERNAME + request.Email)
	if err != nil || exists {
		return errors.New("invalid email or username")
	}

	// check if email or username is already registered
	user, err := repositories.GetUserByEmailOrUsername(request.Email, request.Username)
	if err != nil {
		return errors.New("try again later")
	}
	if user != nil {
		return errors.New("email or username already registered")
	}

	// generate completion registration flow id
	flowID, expiredAt, err := GenerateCompletionRegistrationFlowID(request.Email, request.Username)
	if err != nil {
		return errors.New("try again later")
	}

	// generate link
	link := os.Getenv("FRONTEND_URL") + consts.FRONTEND_REGISTER_ROUTE + "?flow_id=" + flowID

	err = mail.SendVerificationEmail(request.Email, request.Username, mail.VerificationFormTypeLink, mail.VerificationRequestTypeRegister, link, expiredAt).Error
	if err != nil {
		return errors.New("try again later")
	}

	// store the code in redis
	databases.RedisSet(consts.VERIFY_EMAIL_REGISTER_FLOW_ID+flowID, "1", consts.VERIFY_EMAIL_REGISTER_LINK_EXPIRY, datetime.MINUTES)
	// to prevent the user from being registered multiple times
	databases.RedisSet(consts.VERIFY_EMAIL_REGISTER_USERNAME+request.Email, request.Username, consts.VERIFY_EMAIL_REGISTER_LINK_EXPIRY, datetime.MINUTES)

	return nil
}

func UserEmailRegisterWithCode(request *requests.EmailRegisterCodeRequest) error {

	// check if the email or username is already registered
	exists, err := databases.RedisExists(consts.VERIFY_EMAIL_REGISTER_USERNAME + request.Email)
	if err != nil || exists {
		return errors.New("invalid email or username")
	}

	// check if email or username is already registered
	user, err := repositories.GetUserByEmailOrUsername(request.Email, request.Username)
	if err != nil {
		return errors.New("try again later")
	}
	if user != nil {
		return errors.New("email or username already registered")
	}

	verificationCode, err := GenerateVerificationCode()
	if err != nil {
		return errors.New("try again later")
	}
	expiredAt := time.Now().Add(time.Duration(consts.VERIFY_EMAIL_REGISTER_CODE_EXPIRY) * time.Minute).Format(consts.DATETIME_FORMAT)

	err = mail.SendVerificationEmail(request.Email, request.Username, mail.VerificationFormTypeCode, mail.VerificationRequestTypeRegister, verificationCode, expiredAt).Error
	if err != nil {
		return errors.New("try again later")
	}

	hashedPassword, err := crypto.HashPassword(request.Password)
	if err != nil {
		return errors.New("try again later")
	}

	// store the code in redis
	databases.RedisSet(consts.VERIFY_EMAIL_REGISTER_USERNAME+request.Email, request.Username, consts.VERIFY_EMAIL_REGISTER_CODE_EXPIRY, datetime.MINUTES)
	databases.RedisSet(consts.VERIFY_EMAIL_REGISTER_CODE+request.Email, verificationCode, consts.VERIFY_EMAIL_REGISTER_CODE_EXPIRY, datetime.MINUTES)
	databases.RedisSet(consts.VERIFY_EMAIL_REGISTER_PASSWORD+request.Email, hashedPassword, consts.VERIFY_EMAIL_REGISTER_CODE_EXPIRY, datetime.MINUTES)

	// TODO: Limit the frequency
	return nil
}

func UserEmailRegisterLinkVerify(request *requests.EmailRegisterLinkVerifyRequest) error {

	// check if the registration code is existing
	flowID, err := databases.RedisGet(consts.VERIFY_EMAIL_REGISTER_FLOW_ID + request.FlowId)

	if err != nil || flowID == "" {
		return errors.New("invalid flow id")
	}

	// decode the registration code
	email, username, _, err := DecodeCompletionRegistrationFlowID(request.FlowId)
	if err != nil {
		return err
	}

	nickname := request.Nickname
	if nickname == "" {
		nickname = username
	}

	password := request.Password
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return err
	}

	// create user
	err = repositories.CreateUser(email, username, hashedPassword, nickname)
	if err != nil {
		return errors.New("create user failed")
	}

	// delete the flowid, username from redis
	databases.RedisDel(consts.VERIFY_EMAIL_REGISTER_FLOW_ID + request.FlowId)
	databases.RedisDel(consts.VERIFY_EMAIL_REGISTER_USERNAME + email)

	return nil
}

func UserEmailRegisterCodeVerify(request *requests.EmailRegisterCodeVerifyRequest) error {

	// check if the registration code is existing
	code, err := databases.RedisGet(consts.VERIFY_EMAIL_REGISTER_CODE + request.Email)

	if err != nil {
		return errors.New("invalid code")
	}

	if code != request.Code {
		return errors.New("invalid code")
	}

	// if the code is correct, get the username and email
	username, errUsername := databases.RedisGet(consts.VERIFY_EMAIL_REGISTER_USERNAME + request.Email)
	password, errPassword := databases.RedisGet(consts.VERIFY_EMAIL_REGISTER_PASSWORD + request.Email) // this is the hashed password

	if errUsername != nil || errPassword != nil || username == "" || password == "" {
		return errors.New("invalid code")
	}

	nickname := username

	// create user
	err = repositories.CreateUser(request.Email, username, password, nickname)
	if err != nil {
		return errors.New("create user failed")
	}

	// delete the code, username, password from redis
	databases.RedisDel(consts.VERIFY_EMAIL_REGISTER_CODE + request.Email)
	databases.RedisDel(consts.VERIFY_EMAIL_REGISTER_USERNAME + username)
	databases.RedisDel(consts.VERIFY_EMAIL_REGISTER_PASSWORD + request.Email)

	return nil
}

// check the registeration code expired status
func CheckUserEmailRegisterLinkExpired(flowID string) (map[string]string, error) {

	// check if the flow id is existing
	exist, err := databases.RedisExists(consts.VERIFY_EMAIL_REGISTER_FLOW_ID + flowID)
	if err != nil || !exist {
		return nil, errors.New("invalid flow id")
	}

	// decode the registration code
	email, username, expiredAt, err := DecodeCompletionRegistrationFlowID(flowID)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"email":     email,
		"username":  username,
		"expiredAt": expiredAt,
	}, nil
}
