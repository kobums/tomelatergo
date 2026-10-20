# tomelater 서버 (go fiber 기반)

이 프로젝트는 Go와 gofiber를 이용하여 만든 간단한 서버입니다.

mariadb 테이블의 crud 위주의 서버이다.

## 실행

1. make run
# tomelatergo

## 환경변수 (인증)

- `JWT_SECRET` — JWT 서명 키 (32자 이상, 필수)
- `APPLE_CLIENT_IDS` — Apple ID 토큰의 허용 `aud` (쉼표 구분). 기본값: iOS 번들 ID `com.tomelater.tomelater`
- `GOOGLE_CLIENT_IDS` — Google ID 토큰의 허용 `aud` = OAuth 클라이언트 ID (쉼표 구분). **비어 있으면 Google 로그인 거부**
- `JWT_TTL_HOURS` — 발급 토큰 유효기간(시간). 기본값 `720`(30일). 만료되면 클라이언트가 401 을 받고 재로그인한다
- `KAKAO_APP_ID` — 카카오 앱 ID(숫자). 카카오 액세스 토큰의 `app_id`(`/v1/user/access_token_info`)와 비교한다. **비어 있으면 카카오 로그인 거부**
- `NAVER_CLIENT_ID` — 네이버 애플리케이션 Client ID. **비어 있으면 네이버 로그인 거부**. 네이버는 토큰의 발급 client_id 를 조회하는 API 가 없어 서버는 프로필 API 로 토큰 유효성만 확인한다(다른 앱 토큰을 완전히 막으려면 인가 코드를 서버에서 교환하는 방식으로 바꿔야 함)

로그인은 `POST /api/auth/login` (JSON 본문). `GET /api/jwt?loginid=&passwd=` 는 구버전 호환용(deprecated)이다.
