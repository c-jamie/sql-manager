package data

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"math/big"
	// "math/rand"
)

type ROToken struct {
	token    string
	is_valid bool
}

type ROTokenModel struct {
	DB *sql.DB
}

func (m ROTokenModel) Register() (*ROToken, error) {
	token, err := GenerateRandomStringURLSafe(35)
	if err != nil {
		return nil, err
	}
	token = "smt_" + token

	added_token := ""
	query := `
		insert into token(token, created_at)
		values 		($1 now())
		returning 	token
	`
	args := []interface{}{token}
	err = m.DB.QueryRow(query, args...).Scan(added_token)

	if err != nil {
		return nil, err
	} else {
		return &ROToken{token: added_token, is_valid: true}, nil
	}
}

func (m ROTokenModel) Remove(token string) (error) {

	query := `
		update 		token
		set			deleted_at 	= now()
		where		token = $1
	`
	args := []interface{}{token}
	result, err := m.DB.Exec(query, args...)

	if err != nil {
		return fmt.Errorf("unable to remove token %w", err)
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("unable to remove token %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (m ROTokenModel) IsValid(token string) bool {

	query := `
		select 		token
					, created_at
		from 		token as t
		where		t.token = $1
	`

	valid_token := ""
	err := m.DB.QueryRow(query, token).Scan(valid_token)

	if err != nil {
		return false
	} else {
		return true
	}

}

func GenerateRandomStringURLSafe(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
