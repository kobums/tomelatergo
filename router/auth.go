package router

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"strings"
	"tomelater/global"
	"tomelater/global/jwt"
	"tomelater/global/log"
	"tomelater/global/setting"
	"tomelater/global/time"
	"tomelater/models"
	"tomelater/models/ipblock"
	"tomelater/models/systemlog"
	"tomelater/models/user"

	jwtgo "github.com/dgrijalva/jwt-go/v4"
	"github.com/gofiber/fiber/v2"
	"resty.dev/v3"
)

type KakaoResponse struct {
	ID          int64  `json:"id"`
	ConnectedAt string `json:"connected_at"`
	Properties  struct {
		Nickname string `json:"nickname"`
	} `json:"properties"`
	KakaoAccount struct {
		ProfileNicknameNeedsAgreement bool `json:"profile_nickname_needs_agreement"`
		Profile                       struct {
			Nickname          string `json:"nickname"`
			IsDefaultNickname bool   `json:"is_default_nickname"`
		} `json:"profile"`
		HasEmail            bool   `json:"has_email"`
		EmailNeedsAgreement bool   `json:"email_needs_agreement"`
		IsEmailValid        bool   `json:"is_email_valid"`
		IsEmailVerified     bool   `json:"is_email_verified"`
		Email               string `json:"email"`
	} `json:"kakao_account"`
}

type NaverResponse struct {
	Response struct {
		ID           string `json:"id"`
		Email        string `json:"email"`
		Name         string `json:"name"`
		Nickname     string `json:"nickname"`
		ProfileImage string `json:"profile_image"`
	} `json:"response"`
}

type GoogleResponse struct {
	Iss           string `json:"iss"`
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Nonce         string `json:"nonce"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Iat           string `json:"iat"`
	Exp           string `json:"exp"`
	Alg           string `json:"alg"`
	Kid           string `json:"kid"`
	Type          string `json:"type"`
}

type AppleResponse struct {
	Iss            string `json:"iss"`
	Aud            string `json:"aud"`
	Exp            int64  `json:"exp"`
	Iat            int64  `json:"iat"`
	Sub            string `json:"sub"`
	CHash          string `json:"c_hash"`
	Email          string `json:"email"`
	EmailVerified  string `json:"email_verified"`
	AuthTime       int64  `json:"auth_time"`
	NonceSupported bool   `json:"nonce_supported"`
}

type ApplePublicKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type AppleKeys struct {
	Keys []ApplePublicKey `json:"keys"`
}

var JwtAuthRequired = func(c *fiber.Ctx) error {
	if values := c.Get("Authorization"); len(values) > 0 {
		str := values

		claims, err := jwt.Check(str)
		if err == nil {
			user := (*claims).User
			user.Passwd = ""
			c.Locals("user", &user)
			return c.Next()
		}
	}

	path := c.Path()
	u, _ := url.Parse(path)

	if u.Path == "/api/jwt" || u.Path == "/api/auth/login" {
		return c.Next()
	}

	// 인증 없이 허용하는 것은 회원가입(POST /api/user) 뿐이다.
	// 이전에는 "/api/user" 로 시작하는 모든 경로(목록/조회/수정/삭제)가 인증 없이 열려 있었다.
	if u.Path == "/api/user" && c.Method() == fiber.MethodPost {
		return c.Next()
	}

	// if u.Path == "/api/user/changepasswd" {
	// 	return c.Next()
	// }

	// if u.Path == "/api/integration/approve" {
	// 	return c.Next()
	// }

	// if u.Path == "/api/upload/sync" {
	// 	return c.Next()
	// }

	log.Info().Msg("Jwt header is broken")

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "error", "message": "unauthorized"})
}

func JwtAuth(c *fiber.Ctx, loginType string, loginToken string, loginid string, passwd string) map[string]interface{} {

	conn := models.NewConnection()
	defer conn.Close()

	userManager := models.NewUserManager(conn)
	loginlogManager := models.NewLoginlogManager(conn)
	ipblockManager := models.NewIpblockManager(conn)

	var item *models.User

	if loginType == "kakao" {
		client := resty.New()

		headers := map[string]string{
			"Authorization": fmt.Sprintf("Bearer %v", loginToken),
			"Content-Type":  "application/json",
		}

		url := "https://kapi.kakao.com/v2/user/me"
		resp, err := client.R().
			EnableTrace().
			SetHeaders(headers).
			Get(url)

		if err != nil {
			log.Println(err)
			return map[string]interface{}{
				"code":    "error",
				"message": "request api error",
			}
		}
		response := KakaoResponse{}
		global.JsonDecode(resp.String(), &response)

		connectid := fmt.Sprintf("%v", response.ID)
		if connectid == "" {
			log.Println(response)
			return map[string]interface{}{
				"code":    "error",
				"message": "connect id error",
				"item":    response,
			}
		}
		item = userManager.GetByConnectid(connectid)
		if item == nil || item.Id == 0 {
			return map[string]interface{}{
				"code":    "error",
				"message": "join",
				"item":    response,
			}
		}
	} else if loginType == "naver" {
		client := resty.New()

		headers := map[string]string{
			"Authorization": fmt.Sprintf("Bearer %v", loginToken),
			"Content-Type":  "application/json",
		}

		url := "https://openapi.naver.com/v1/nid/me"
		resp, err := client.R().
			EnableTrace().
			SetHeaders(headers).
			Get(url)

		if err != nil {
			log.Println(err)
			return map[string]interface{}{
				"code":    "error",
				"message": "request api error",
			}
		}
		response := NaverResponse{}
		global.JsonDecode(resp.String(), &response)

		connectid := fmt.Sprintf("%v", response.Response.ID)
		if connectid == "" {
			log.Println(response)
			return map[string]interface{}{
				"code":    "error",
				"message": "connect id error",
				"item":    response,
			}
		}
		item = userManager.GetByConnectid(connectid)
		if item == nil || item.Id == 0 {
			return map[string]interface{}{
				"code":    "error",
				"message": "join",
				"item":    response,
			}
		}
	} else if loginType == "google" {
		client := resty.New()

		url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%v", loginToken)
		resp, err := client.R().
			EnableTrace().
			Get(url)

		if err != nil {
			log.Println(err)
			return map[string]interface{}{
				"code":    "error",
				"message": "request api error",
			}
		}
		if resp.StatusCode() != 200 {
			return map[string]interface{}{
				"code":    "error",
				"message": "google token validation failed",
			}
		}
		response := GoogleResponse{}
		global.JsonDecode(resp.String(), &response)

		// tokeninfo 는 서명·만료만 확인한다. 다른 앱용으로 발급된 ID 토큰을 그대로 들고 와도
		// 로그인되지 않도록 aud(우리 OAuth 클라이언트 ID)와 iss 를 반드시 확인한다.
		if err := checkGoogleClaims(response); err != nil {
			log.Println(err)
			return map[string]interface{}{
				"code":    "error",
				"message": "google token validation failed",
			}
		}

		connectid := fmt.Sprintf("%v", response.Sub)
		if connectid == "" {
			log.Println(response)

			return map[string]interface{}{
				"code":    "error",
				"message": "connect id error",
				"item":    response,
			}
		}
		item = userManager.GetByConnectid(connectid)
		if item == nil || item.Id == 0 {
			return map[string]interface{}{
				"code":    "error",
				"message": "join",
				"item":    response,
			}
		}

	} else if loginType == "apple" {
		response, err := validateAppleToken(loginToken)
		if err != nil {
			log.Println(err)
			return map[string]interface{}{
				"code":    "error",
				"message": "apple token validation failed",
			}
		}

		connectid := fmt.Sprintf("%v", response.Sub)
		if connectid == "" {
			log.Println(response)

			return map[string]interface{}{
				"code":    "error",
				"message": "connect id error",
				"item":    response,
			}
		}
		item = userManager.GetByConnectid(connectid)
		if item == nil || item.Id == 0 {
			return map[string]interface{}{
				"code":    "error",
				"message": "join",
				"item":    response,
			}
		}

	} else {
		item = userManager.GetByLoginid(loginid)

		if item == nil {
			return map[string]interface{}{
				"code":    "error",
				"message": "user not found",
			}
		}

		if !jwt.CheckPasswd(item.Passwd, passwd) {
			return map[string]interface{}{
				"code":    "error",
				"message": "wrong password",
			}
		}

		if item.Use != user.UseUse {
			return map[string]interface{}{
				"code":    "error",
				"message": "status error",
			}
		}

		// instance := setting.GetInstance()
		// days := instance.SettingInt("user.passwd.changeday")

		// if days > 0 {
		// 	now := time.Now()
		// 	now = now.AddDate(0, 0, -1*days)

		// 	lastlogin := time.Parse(item.Lastchangepasswddate)
		// 	if now.After(lastlogin) {
		// 		return map[string]interface{}{
		// 			"code":    "error",
		// 			"message": "expired passwd",
		// 		}
		// 	}
		// }
	}

	typeid := ipblock.TypeAdmin
	if item.Level < user.LevelAdmin {
		typeid = ipblock.TypeNormal
	}

	ipblocks := ipblockManager.Find([]interface{}{
		models.Where{Column: "type", Value: typeid, Compare: "="},
		models.Where{Column: "use", Value: ipblock.UseUse, Compare: "="},
		models.Ordering("ib_order,ib_id"),
	})

	access := ipblock.PolicyGrant

	for _, v := range ipblocks {
		check, err := setting.NewIP(v.Address)
		if err != nil {
			continue
		}

		ip, _ := setting.NewIP(c.IP())

		if check.Contains(ip) {
			access = v.Policy
		}
	}

	if access == ipblock.PolicyDeny {
		systemlogManager := models.NewSystemlogManager(conn)
		item := models.Systemlog{
			Type:    systemlog.TypeLogin,
			Content: fmt.Sprintf("%v IP에서 %v 계정 접근이 거부되었습니다", c.IP(), loginid),
		}

		err := systemlogManager.Insert(&item)
		if err != nil {
			log.Error().Msg(err.Error())
		}

		return map[string]interface{}{
			"code":    "error",
			"message": "access denied",
		}
	}

	now := time.Now()
	ip, _ := setting.NewIP(c.IP())
	loginlog := models.Loginlog{User: item.Id, Ip: c.IP(), Ipvalue: ip.StartInt, Date: now.Datetime()}
	err := loginlogManager.Insert(&loginlog)
	if err != nil {
		log.Error().Msg(err.Error())
		return map[string]interface{}{
			"code":    "error",
			"message": "db error",
		}
	}

	err = userManager.UpdateLogindateById(now.Datetime(), item.Id)
	if err != nil {
		log.Error().Msg(err.Error())
		return map[string]interface{}{
			"code":    "error",
			"message": "db error",
		}
	}

	item.Passwd = ""

	token := jwt.MakeToken(*item)

	return map[string]interface{}{
		"code":  "ok",
		"token": token,
		"user":  item,
	}
}

// ---- OAuth audience (client ID) allow-lists ----
//
// APPLE_CLIENT_IDS / GOOGLE_CLIENT_IDS (쉼표 구분) 로 덮어쓸 수 있다.
//   - Apple: 기본값은 tomelaterapp 의 iOS 번들 ID (ios/Runner.xcodeproj PRODUCT_BUNDLE_IDENTIFIER).
//     웹(Sign in with Apple JS)용 Services ID 를 쓰면 여기에 추가해야 한다.
//   - Google: 앱 저장소에 OAuth 클라이언트 ID 가 없어 기본값이 없다. 비어 있으면 Google 로그인을
//     거부한다(aud 를 확인할 수 없으면 다른 앱용 토큰으로 로그인할 수 있으므로 fail-closed).

const defaultAppleClientIDs = "com.tomelater.tomelater"

var (
	appleClientIDs  = envList("APPLE_CLIENT_IDS", defaultAppleClientIDs)
	googleClientIDs = envList("GOOGLE_CLIENT_IDS", "")
)

func init() {
	if len(googleClientIDs) == 0 {
		log.Info().Msg("GOOGLE_CLIENT_IDS is not set: Google login is disabled until it is configured")
	}
}

func envList(name string, def string) []string {
	v, ok := os.LookupEnv(name)
	if !ok {
		v = def
	}
	var out []string
	for _, s := range strings.Split(v, ",") {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// audAllowed reports whether the aud claim (string or array) names one of allowed.
func audAllowed(aud interface{}, allowed []string) bool {
	var auds []string
	switch v := aud.(type) {
	case string:
		auds = []string{v}
	case []interface{}:
		for _, a := range v {
			if s, ok := a.(string); ok {
				auds = append(auds, s)
			}
		}
	case []string:
		auds = v
	}
	for _, a := range auds {
		for _, ok := range allowed {
			if a == ok {
				return true
			}
		}
	}
	return false
}

func checkGoogleClaims(r GoogleResponse) error {
	if len(googleClientIDs) == 0 {
		return fmt.Errorf("google login disabled: GOOGLE_CLIENT_IDS not set")
	}
	if r.Iss != "accounts.google.com" && r.Iss != "https://accounts.google.com" {
		return fmt.Errorf("google token: unexpected iss %q", r.Iss)
	}
	if !audAllowed(r.Aud, googleClientIDs) {
		return fmt.Errorf("google token: unexpected aud %q", r.Aud)
	}
	return nil
}

func validateAppleToken(idToken string) (*AppleResponse, error) {
	claims := jwtgo.MapClaims{}

	// 서명(Apple 공개키)·만료·iss 는 파서가, aud 는 아래에서 허용 목록으로 직접 확인한다.
	// 예전엔 audience 에러를 무시해서 다른 앱용으로 발급된 Apple 토큰으로도 로그인할 수 있었다.
	token, err := jwtgo.ParseWithClaims(idToken, claims, func(token *jwtgo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtgo.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		publicKey, err := getApplePublicKey(kid)
		if err != nil {
			return nil, err
		}

		return publicKey, nil
	}, jwtgo.WithoutAudienceValidation(), jwtgo.WithIssuer("https://appleid.apple.com"))

	if err != nil {
		return nil, err
	}
	if token == nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	if !audAllowed(claims["aud"], appleClientIDs) {
		return nil, fmt.Errorf("apple token: unexpected aud")
	}

	response := &AppleResponse{}

	if iss, exists := claims["iss"]; exists {
		response.Iss = iss.(string)
	}

	if aud, ok := claims["aud"].(string); ok {
		response.Aud = aud
	}

	if sub, exists := claims["sub"]; exists {
		response.Sub = sub.(string)
	}

	if exp, exists := claims["exp"]; exists {
		response.Exp = int64(exp.(float64))
	}

	if iat, exists := claims["iat"]; exists {
		response.Iat = int64(iat.(float64))
	}

	if authTime, exists := claims["auth_time"]; exists {
		response.AuthTime = int64(authTime.(float64))
	}

	if email, exists := claims["email"]; exists {
		response.Email = email.(string)
	}

	if emailVerified, exists := claims["email_verified"]; exists {
		if verified, ok := emailVerified.(string); ok {
			response.EmailVerified = verified
		} else if verified, ok := emailVerified.(bool); ok {
			if verified {
				response.EmailVerified = "true"
			} else {
				response.EmailVerified = "false"
			}
		}
	}

	if cHash, exists := claims["c_hash"]; exists {
		response.CHash = cHash.(string)
	}

	return response, nil
}

func getApplePublicKey(kid string) (*rsa.PublicKey, error) {
	client := resty.New()

	resp, err := client.R().
		Get("https://appleid.apple.com/auth/keys")

	if err != nil {
		return nil, err
	}

	var appleKeys AppleKeys
	global.JsonDecode(resp.String(), &appleKeys)

	for _, key := range appleKeys.Keys {
		if key.Kid == kid {
			return parseRSAPublicKey(key.N, key.E)
		}
	}

	return nil, fmt.Errorf("public key not found for kid: %s", kid)
}

func parseRSAPublicKey(n, e string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(n)
	if err != nil {
		return nil, err
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(e)
	if err != nil {
		return nil, err
	}

	var eInt int
	for _, b := range eBytes {
		eInt = eInt*256 + int(b)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: eInt,
	}, nil
}
