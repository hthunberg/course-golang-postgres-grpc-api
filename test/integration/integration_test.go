//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSomeIntegration(t *testing.T) {
	err := testDbInstance.Ping(context.Background())
	require.NoError(t, err)
}

func TestCreateUser(t *testing.T) {
	bankClient, err := newTestBankCLient(testBankBaseURL)
	require.NoError(t, err)

	userReq := UserRequest{
		UserName: "johndoe",
		Password: "qwerty",
		FullName: "John Doe",
		Email:    "john.doe@testbank.qwerty",
	}

	userReqJson, err := marshalJson(userReq)
	require.NoError(t, err)

	res, resBody, err := bankClient.createUser(bytes.NewReader(userReqJson))
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	m, err := unMarshalJson(resBody)
	require.NoError(t, err)

	assertJsonElement(t, m, "username", userReq.UserName)
	assertJsonElement(t, m, "full_name", userReq.FullName)
	assertJsonElement(t, m, "email", userReq.Email)
}

func TestCreateUserAccount(t *testing.T) {
	bankClient, err := newTestBankCLient(testBankBaseURL)
	require.NoError(t, err)

	userReq := UserRequest{
		UserName: "johndoe",
		Password: "qwerty",
		FullName: "John Doe",
		Email:    "john.doe@testbank.qwerty",
	}

	userReqJson, err := marshalJson(userReq)
	require.NoError(t, err)

	_, _, err = bankClient.createUser(bytes.NewReader(userReqJson))
	require.NoError(t, err)

	accountReq := AccountRequest{
		Owner:    "johndoe",
		Currency: "SEK",
	}

	accountReqJson, err := marshalJson(accountReq)
	require.NoError(t, err)

	res, resBody, err := bankClient.createAccount(bytes.NewReader(accountReqJson))
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	m, err := unMarshalJson(resBody)
	require.NoError(t, err)

	assertJsonElement(t, m, "owner", accountReq.Owner)
	assertJsonElement(t, m, "currency", accountReq.Currency)
	assertJsonElement(t, m, "balance", float64(0))
}

func assertJsonElement(t *testing.T, m map[string]any, key string, expected any) {
	actual, exists := m[key]
	assert.True(t, exists)
	assert.Equal(t, expected, actual)
}

func marshalJson(v any) ([]byte, error) {
	buffer, err := json.Marshal(v)
	return buffer, err
}

func unMarshalJson(v []byte) (map[string]any, error) {
	var resp map[string]any
	err := json.Unmarshal(v, &resp)
	return resp, err
}
