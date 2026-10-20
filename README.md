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

로그인은 `POST /api/auth/login` (JSON 본문). `GET /api/jwt?loginid=&passwd=` 는 구버전 호환용(deprecated)이다.
