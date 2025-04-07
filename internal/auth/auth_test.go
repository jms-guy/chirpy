package auth

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestPasswords(t *testing.T) {
	start := "Thisismypassword"

	hash, err := HashPassword(start)
	if err != nil {
		t.Errorf("Failed hash: %s", err)
	}

	checkErr := CheckPasswordHash(hash, start)
	if checkErr != nil {
		t.Errorf("Failed check: %s", err)
	}
}

func TestToken(t *testing.T) {
	idString := "8dc3802e-3797-4560-9391-c56cbf479ae5"
	idBytes, _ := hex.DecodeString(idString)
	id, _ := uuid.FromBytes(idBytes)

	secret := "Thisisatokensecret"

	signature, err := MakeJWT(id, secret)
	if err != nil {
		t.Error(err)
	}

	//time.Sleep(6 * time.Second)
	user, err := ValidateJWT(signature, secret)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(user)
}

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer TOKEN_STRING")
	want := "TOKEN_STRING"

	s, err := GetBearerToken(headers)
	if err != nil {
		t.Error(err)
	}
	if s == want {
		return
	} else {
		t.Error("fail")
	}
}