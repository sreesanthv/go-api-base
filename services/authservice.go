package services

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/sreesanthv/go-api-base/database"
)

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
	td.AccessToken, err = at.SignedString([]byte(viper.GetString("jwt_secret_refresh")))
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
func (s *AuthService) PersistToken(userId int32, td *TokenDetails) error {
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
