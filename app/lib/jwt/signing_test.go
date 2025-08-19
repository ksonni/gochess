package jwt

import "testing"

const kTestSecret = "9Atwe3XmzcsoEZZZecnI4WeMoqAn7vlHg0B3XmGXS50="
const kBadSecret = "8Atwe3XmzcsoEZZZecnI4WeMoqAn7vlHg0B3XmGXS50="

type TestClaims struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
}

var c = TestClaims{
	Id:   "id",
	Name: "John",
}

func TestEncodesDecodes(t *testing.T) {
	token, err := CreateToken(c, []byte(kTestSecret))
	if err != nil {
		t.Fatal("Failed to encode initial JWT")
	}
	var out TestClaims
	err = ParseToken(token, []byte(kTestSecret), &out)
	if err != nil {
		t.Fatalf("Failed to decode encoded JWT: %v", err)
	}
	if out != c {
		t.Errorf("ParseToken(%s, %s, {}) got: %v, want %v", token, kTestSecret, out, c)
	}
}

func TestRejectsInvalidToken(t *testing.T) {
	token, err := CreateToken(c, []byte(kBadSecret))
	if err != nil {
		t.Fatal("Failed to encode initial JWT")
	}
	var out TestClaims
	err = ParseToken(token, []byte(kTestSecret), &out)
	if err == nil {
		t.Errorf("Invalid JWT should fail to parse, but succeeded")
	}
}
