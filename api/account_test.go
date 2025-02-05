package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/mrafid01/simplebank/db/mock"
	db "github.com/mrafid01/simplebank/db/sqlc"
	"github.com/mrafid01/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct{
		name string
		accountId int64
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			accountId: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "AccountNotFound",
			accountId: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			accountId: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidID",
			accountId: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
		
			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)
		
			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()
		
			url := fmt.Sprintf("/accounts/%d", tc.accountId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
		
			server.router.ServeHTTP(recorder, request)
		
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListAccountsAPI(t *testing.T) {
	var listAccounts []db.Account
	for i := 0; i < 10; i++ {
		listAccounts = append(listAccounts, randomAccount())
	}
	testCases := []struct{
		name string
		pageSize int32
		pageId int32
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			pageSize: 10,
			pageId: 1,
			buildStubs: func(store *mockdb.MockStore) {
				expectedArgs := db.ListAccountsParams{
					Limit: 10,
					Offset: 0,
				}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(expectedArgs)).
					Times(1).
					Return(listAccounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchListAccounts(t, recorder.Body, listAccounts)
			},
		},
		{
			name: "BadRequest",
			pageSize: 100,
			pageId: 1,
			buildStubs: func(store *mockdb.MockStore) {
				expectedArgs := db.ListAccountsParams{
					Limit: 100,
					Offset: 0,
				}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(expectedArgs)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			pageSize: 10,
			pageId: 1,
			buildStubs: func(store *mockdb.MockStore) {
				expectedArgs := db.ListAccountsParams{
					Limit: 10,
					Offset: 0,
				}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(expectedArgs)).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
		
			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)
		
			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()
		
			url := fmt.Sprintf("/accounts?page_size=%d&page_id=%d", tc.pageSize, tc.pageId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
		
			server.router.ServeHTTP(recorder, request)
		
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}
func TestCreateAccountAPI(t *testing.T) {
	account := randomAccount()
	testCases := []struct{
		name string
		body gin.H
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			body: gin.H{
				"Owner": account.Owner,
				"Currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				expectedArgs := db.CreateAccountParams{
					Owner: account.Owner,
					Balance: 0,
					Currency: account.Currency,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(expectedArgs)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"Owner": "",
				"Balance": 0,
				"Currency": "",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"Owner": account.Owner,
				"Balance": account.Balance,
				"Currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				expectedArgs := db.CreateAccountParams{
					Owner: account.Owner,
					Balance: 0,
					Currency: account.Currency,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(expectedArgs)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
		
			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)
		
			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()
		
			url := "/accounts"
			
			// Convert the struct to JSON
			jsonData, err := json.Marshal(tc.body)
			require.NoError(t, err)
			
			// Create an io.Reader from the JSON data
			body := bytes.NewBuffer(jsonData)
			request, err := http.NewRequest(http.MethodPost, url, body)
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")
		
			server.router.ServeHTTP(recorder, request)
		
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount() db.Account{
	return db.Account{
		ID: util.RandomInt(1,1000),
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}



func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func requireBodyMatchListAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}