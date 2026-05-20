package jwt

import (
	"testing"
	"time"
)

const (
	testAccessSecret  = "test-access-secret"
	testRefreshSecret = "test-refresh-secret"
)

func TestGenerateAndParseAccessToken(t *testing.T) {
	tests := []struct {
		name    string
		uid     int64
		role    string
		mid     int64
		expire  time.Duration
		wantErr bool
	}{
		{
			name:   "consumer用户正常生成",
			uid:    1001,
			role:   "consumer",
			mid:    0,
			expire: 2 * time.Hour,
		},
		{
			name:   "merchant用户含mid",
			uid:    2001,
			role:   "merchant",
			mid:    100,
			expire: 2 * time.Hour,
		},
		{
			name:   "admin用户",
			uid:    9999,
			role:   "admin",
			mid:    0,
			expire: 2 * time.Hour,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token, err := GenerateAccessToken(tc.uid, tc.role, tc.mid, testAccessSecret, tc.expire)
			if err != nil {
				t.Fatalf("GenerateAccessToken error: %v", err)
			}
			if token == "" {
				t.Fatal("token为空")
			}

			claims, err := ParseAccessToken(token, testAccessSecret)
			if err != nil {
				t.Fatalf("ParseAccessToken error: %v", err)
			}
			if claims.UID != tc.uid {
				t.Errorf("UID: got %d, want %d", claims.UID, tc.uid)
			}
			if claims.Role != tc.role {
				t.Errorf("Role: got %s, want %s", claims.Role, tc.role)
			}
			if claims.MID != tc.mid {
				t.Errorf("MID: got %d, want %d", claims.MID, tc.mid)
			}
		})
	}
}

func TestParseAccessToken_WrongSecret(t *testing.T) {
	token, err := GenerateAccessToken(1, "consumer", 0, testAccessSecret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ParseAccessToken(token, "wrong-secret")
	if err != ErrTokenInvalid {
		t.Errorf("期望 ErrTokenInvalid, 得到: %v", err)
	}
}

func TestParseAccessToken_Expired(t *testing.T) {
	token, err := GenerateAccessToken(1, "consumer", 0, testAccessSecret, -time.Second)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ParseAccessToken(token, testAccessSecret)
	if err != ErrTokenExpired {
		t.Errorf("期望 ErrTokenExpired, 得到: %v", err)
	}
}

func TestParseAccessToken_InvalidToken(t *testing.T) {
	_, err := ParseAccessToken("not.a.token", testAccessSecret)
	if err != ErrTokenInvalid {
		t.Errorf("期望 ErrTokenInvalid, 得到: %v", err)
	}
}

func TestParseAccessToken_EmptyToken(t *testing.T) {
	_, err := ParseAccessToken("", testAccessSecret)
	if err != ErrTokenInvalid {
		t.Errorf("期望 ErrTokenInvalid, 得到: %v", err)
	}
}

func TestGenerateAndParseRefreshToken(t *testing.T) {
	uid := int64(12345)
	token, err := GenerateRefreshToken(uid, testRefreshSecret, 7*24*time.Hour)
	if err != nil {
		t.Fatalf("GenerateRefreshToken error: %v", err)
	}
	if token == "" {
		t.Fatal("refresh token为空")
	}

	parsedUID, err := ParseRefreshToken(token, testRefreshSecret)
	if err != nil {
		t.Fatalf("ParseRefreshToken error: %v", err)
	}
	if parsedUID != uid {
		t.Errorf("uid: got %d, want %d", parsedUID, uid)
	}
}

func TestParseRefreshToken_Expired(t *testing.T) {
	token, err := GenerateRefreshToken(1, testRefreshSecret, -time.Second)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ParseRefreshToken(token, testRefreshSecret)
	if err != ErrTokenExpired {
		t.Errorf("期望 ErrTokenExpired, 得到: %v", err)
	}
}

func TestParseRefreshToken_WrongSecret(t *testing.T) {
	token, err := GenerateRefreshToken(1, testRefreshSecret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ParseRefreshToken(token, "wrong-secret")
	if err != ErrTokenInvalid {
		t.Errorf("期望 ErrTokenInvalid, 得到: %v", err)
	}
}

func TestGenerateAccessToken_EmptySecret(t *testing.T) {
	// 空 secret 仍然可以生成（golang-jwt允许），但解析时使用错误的secret应当失败
	token, err := GenerateAccessToken(1, "consumer", 0, "", time.Hour)
	if err != nil {
		t.Skip("空secret生成报错，跳过")
	}
	// 使用正确的空 secret 解析应当成功
	claims, err := ParseAccessToken(token, "")
	if err != nil {
		t.Fatalf("ParseAccessToken with empty secret error: %v", err)
	}
	if claims.UID != 1 {
		t.Errorf("UID: got %d, want 1", claims.UID)
	}
}

func TestTokenClaims_ExpiresAt(t *testing.T) {
	expire := 2 * time.Hour
	before := time.Now()
	token, err := GenerateAccessToken(1, "consumer", 0, testAccessSecret, expire)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := ParseAccessToken(token, testAccessSecret)
	if err != nil {
		t.Fatal(err)
	}
	// 过期时间应在 before+expire 附近（允许1秒误差）
	expectedExp := before.Add(expire)
	actualExp := claims.ExpiresAt.Time
	diff := actualExp.Sub(expectedExp)
	if diff > time.Second || diff < -time.Second {
		t.Errorf("ExpiresAt偏差过大: %v", diff)
	}
}
