package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/sreesanthv/go-api-base/database"
)

const TokenTypeAccess = 0
const TokenTypeRefresh = 1

type AuthService struct {
	logger *logrus.Logger
	store  *database.Store
	redis  *database.Redis
}

func NewAuthService(log *logrus.Logger, store *database.Store, redis *database.Redis) *AuthService {
	return &AuthService{
		logger: log,
		store:  store,
		redis:  redis,
	}
}

func (s *AuthService) GetUser(email string) *database.AccountStore {
	user, _ := s.store.GetUser(email)
	return user
}

func (s *AuthService) GetUserById(id int64) *database.AccountStore {
	user, _ := s.store.GetUserById(id)
	return user
}

// validate password entered  - login
func (s *AuthService) IsValidPassword(act *database.AccountStore, password string) bool {
	return CheckPasswordHash(password, act.Password)
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

// generate tokens with expiry
// access token expiry - 30 minutes
// refresh token expiry - 7 days
func (s *AuthService) CreateToken(user *database.AccountStore) (*TokenDetails, error) {
	td := new(TokenDetails)

	uuidAcc, err := uuid.NewV4()
	if err != nil {
		s.logger.Errorf("Error generating uuid - access:", err)
		return nil, err
	}
	uuidRef, err := uuid.NewV4()
	if err != nil {
		s.logger.Errorf("Error generating uuid - refresh:", err)
		return nil, err
	}

	td.AtExpires = time.Now().Add(time.Minute * 30).Unix()
	td.AccessUuid = uuidAcc.String()
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuidRef.String()

	// access token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = user.ID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(viper.GetString("jwt_secret_access")))
	if err != nil {
		s.logger.Errorf("Error creating access token:", err)
		return nil, err
	}

	// refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = user.ID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(viper.GetString("jwt_secret_refresh")))
	if err != nil {
		s.logger.Errorf("Error creating refresh token:", err)
		return nil, err
	}

	return td, nil
}

// save token details in redis
func (s *AuthService) PersistToken(userId int64, td *TokenDetails) error {
	now := time.Now()
	at := time.Unix(td.AtExpires, 0)
	err := s.redis.Set(td.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now))
	if err != nil {
		return err
	}

	rt := time.Unix(td.RtExpires, 0)
	err = s.redis.Set(td.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now))
	if err != nil {
		return err
	}

	return nil
}

type AccessDetails struct {
	Uuid   string
	UserId int64
}

// parse and validate token
// validate against redis info
func (s *AuthService) ParseToken(tk string, tType int) (*AccessDetails, error) {
	var secret, uuidKey string
	switch tType {
	case TokenTypeAccess:
		secret = viper.GetString("jwt_secret_access")
		uuidKey = "access_uuid"
	case TokenTypeRefresh:
		secret = viper.GetString("jwt_secret_refresh")
		uuidKey = "refresh_uuid"
	}

	token, err := jwt.Parse(tk, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			s.logger.Error(err)
			return nil, err
		}
		return []byte(secret), nil
	})
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		err := fmt.Errorf("Invalid JWT token")
		s.logger.Error(err)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err := fmt.Errorf("Failed to extract JWT claim")
		s.logger.Error(err)
		return nil, err
	}

	uuid, ok := claims[uuidKey].(string)
	if !ok {
		err := fmt.Errorf("%s not present in claim", uuidKey)
		return nil, err
	}

	// fetch info from redis
	id, err := s.redis.Get(uuid)
	if err != nil || id == "" {
		err := fmt.Errorf("Invalid JWT token")
		s.logger.Error(err)
		return nil, err
	}

	tUserId, err := strconv.ParseInt(id, 10, 32)
	if err != nil || id == "" {
		err := fmt.Errorf("Invalid JWT token")
		s.logger.Error(err)
		return nil, err
	}

	// user id in claim
	userId, err := strconv.ParseInt(fmt.Sprintf("%v", claims["user_id"]), 10, 32)
	if err != nil || userId == 0 {
		err := fmt.Errorf("user_id not present in claim: %s", err)
		s.logger.Error(err)
		return nil, err
	}

	if userId != tUserId {
		err := fmt.Errorf("Token user_id mismatch")
		s.logger.Error(err)
		return nil, err
	}

	ad := &AccessDetails{
		Uuid:   uuid,
		UserId: userId,
	}

	return ad, nil
}

func (s *AuthService) DropToken(uuid string) error {
	err := s.redis.Delete(uuid)
	if err != nil {
		s.logger.Error("Failed to drop token")
	}

	return err
}
