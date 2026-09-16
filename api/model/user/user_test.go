package userModel_test

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"sl-api/api/model/user"
	testUtil "sl-api/util/test"
)

func TestFormToModelHashesPassword(t *testing.T) {
	t.Parallel()

	const plaintext = "sup3r-s3cret-pw"
	form := &userModel.Form{
		Name:     "Ana",
		LastName: "Silva",
		Username: "ana",
		Password: plaintext,
	}

	model, err := form.ToModel()
	testUtil.NoError(t, err)

	if model.Password == plaintext {
		t.Fatal("password was stored in plaintext")
	}

	// The stored value must be a valid bcrypt hash of the original password.
	err = bcrypt.CompareHashAndPassword([]byte(model.Password), []byte(plaintext))
	testUtil.NoError(t, err)
}

func TestValidateCredentials(t *testing.T) {
	t.Parallel()

	hashed, err := bcrypt.GenerateFromPassword([]byte("right-pw"), bcrypt.DefaultCost)
	testUtil.NoError(t, err)

	creds := &userModel.UserCredentials{
		Username: "ana",
		Password: string(hashed),
	}

	if !creds.ValidateCredentials(&userModel.LoginForm{Username: "ana", Password: "right-pw"}) {
		t.Fatal("correct password was rejected")
	}

	if creds.ValidateCredentials(&userModel.LoginForm{Username: "ana", Password: "wrong-pw"}) {
		t.Fatal("wrong password was accepted")
	}
}
