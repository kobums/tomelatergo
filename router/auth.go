package router

import (
	"crypto/rsa"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/url"
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

	if u.Path == "/api/jwt" {
		return c.Next()
	}

	if u.Path[:9] == "/api/user" {
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

	return nil
}

func JwtAuth(c *fiber.Ctx, loginid string, passwd string) map[string]interface{} {
	loginType := c.Query("type")
	loginToken := c.Query("token")

	conn := models.NewConnection()
	defer conn.Close()

	userManager := models.NewUserManager(conn)
	loginlogManager := models.NewLoginlogManager(conn)
	ipblockManager := models.NewIpblockManager(conn)

	var item *models.User

	if loginType == "kakao" {
		client := resty.New()
		client.SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true, // Set to true to skip server verification (not recommended for production)
		})

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
		client.SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true, // Set to true to skip server verification (not recommended for production)
		})

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
		client.SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true, // Set to true to skip server verification (not recommended for production)
		})

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
		response := GoogleResponse{}
		global.JsonDecode(resp.String(), &response)

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


func validateAppleToken(idToken string) (*AppleResponse, error) {
	claims := jwtgo.MapClaims{}

	// 토큰을 파싱하되 audience 에러는 무시
	token, err := jwtgo.ParseWithClaims(idToken, claims, func(token *jwtgo.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		publicKey, err := getApplePublicKey(kid)
		if err != nil {
			return nil, err
		}

		return publicKey, nil
	})

	// audience 관련 에러는 무시
	if err != nil && !strings.Contains(err.Error(), "audience") {
		return nil, err
	}

	if token != nil && !token.Valid && err != nil && !strings.Contains(err.Error(), "audience") {
		return nil, fmt.Errorf("invalid token")
	}

	response := &AppleResponse{}

	if iss, exists := claims["iss"]; exists {
		response.Iss = iss.(string)
	}

	if aud, exists := claims["aud"]; exists {
		response.Aud = aud.(string)
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
	client.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})

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
