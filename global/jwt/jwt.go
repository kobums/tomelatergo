package jwt

import (
	"errors"
	"os"
	"strconv"
	"time"
	"tomelater/models"

	"github.com/dgrijalva/jwt-go/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthTokenClaims struct {
	User               models.User `json:"user"`
	jwt.StandardClaims             // 표준 토큰 Claims
}

// JWT 서명 키는 환경변수 JWT_SECRET 으로 주입한다 (소스에 하드코딩 금지).
var _secretCode string = os.Getenv("JWT_SECRET")

// 토큰 유효기간(시간). 환경변수 JWT_TTL_HOURS 로 바꿀 수 있고 기본값은 720시간(30일)이다.
// 예전엔 10년짜리 토큰을 발급해서 유출되면 사실상 영구히 쓸 수 있었다. 앱·웹 클라이언트는
// 401 을 받으면 저장된 토큰을 지우고 로그인 화면으로 보낸다(만료 시 재로그인).
const defaultTokenTTLHours = 720

var tokenTTL = loadTokenTTL()

func loadTokenTTL() time.Duration {
	hours := defaultTokenTTLHours
	if v := os.Getenv("JWT_TTL_HOURS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			hours = n
		}
	}
	return time.Duration(hours) * time.Hour
}

func init() {
	if len(_secretCode) < 32 {
		panic("JWT_SECRET environment variable must be set (at least 32 characters)")
	}
}

func Check(str string) (*AuthTokenClaims, error) {
	if len(str) < 7 || str[:7] != "Bearer " {
		err := errors.New("tokek is broken")
		return nil, err
	}

	token := str[7:]

	claims := AuthTokenClaims{}
	key := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected Signing Method")
		}
		return []byte(_secretCode), nil
	}

	_, err := jwt.ParseWithClaims(token, &claims, key)
	if err != nil {
		return nil, err
	}

	return &claims, nil
}

func MakeToken(item models.User) string {
	now := time.Now()
	at := AuthTokenClaims{
		User: item,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  jwt.At(now),
			ExpiresAt: jwt.At(now.Add(tokenTTL)),
		},
	}

	atoken := jwt.NewWithClaims(jwt.SigningMethodHS256, &at)
	signedAuthToken, _ := atoken.SignedString([]byte(_secretCode))

	return signedAuthToken
}

func CheckPasswd(dbPasswd string, inputPasswd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dbPasswd), []byte(inputPasswd))
	return err == nil
}

func GeneratePasswd(passwd string) (string, error) {
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	} else {
		return string(hashedPasswd), err
	}
}
