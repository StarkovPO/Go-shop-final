package handler

import (
	"bytes"
	"github.com/StarkovPO/Go-shop-final/internal/apperrors"
	mock_handler "github.com/StarkovPO/Go-shop-final/internal/handler/mocks"
	"github.com/StarkovPO/Go-shop-final/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	type mockBehavior func(s *mock_handler.MockServiceInterface, user models.Users)

	testTable := []struct {
		name        string
		requestBody string
		newUser     models.Users
		mockBehavior
		statusCode      int
		expectedBody    string
		expectedHeaders string
	}{
		{
			name:        "successfully registration new user",
			requestBody: `{"login":"new_user","password":"new_password"}`,
			newUser: models.Users{
				Login:    "new_user",
				Password: "new_password",
			},
			mockBehavior: func(s *mock_handler.MockServiceInterface, user models.Users) {
				s.EXPECT().CreateUser(gomock.Any(), user).Return("token", nil)
			},
			statusCode:      200,
			expectedBody:    "",
			expectedHeaders: "token",
		},
		{
			name:        "try to auth without login",
			requestBody: `{"password":"new_password"}`,
			newUser: models.Users{
				Password: "new_password",
			},
			mockBehavior: func(s *mock_handler.MockServiceInterface, user models.Users) {
				s.EXPECT().CreateUser(gomock.Any(), user).Return("", apperrors.ErrBadRequest)
			},
			statusCode:      400,
			expectedBody:    "Bad request\n",
			expectedHeaders: "",
		},
		{
			name:        "try to auth without login",
			requestBody: `{"login":"new_user"}`,
			newUser: models.Users{
				Login: "new_user",
			},
			mockBehavior: func(s *mock_handler.MockServiceInterface, user models.Users) {
				s.EXPECT().CreateUser(gomock.Any(), user).Return("", apperrors.ErrBadRequest)
			},
			statusCode:      400,
			expectedBody:    "Bad request\n",
			expectedHeaders: "",
		},
		{
			name:        "register already exist user",
			requestBody: `{"login":"new_user","password":"new_password"}`,
			newUser: models.Users{
				Login:    "new_user",
				Password: "new_password",
			},
			mockBehavior: func(s *mock_handler.MockServiceInterface, user models.Users) {
				s.EXPECT().CreateUser(gomock.Any(), user).Return("", apperrors.ErrLoginAlreadyExist)
			},
			statusCode:      409,
			expectedBody:    "Login already exist\n",
			expectedHeaders: "",
		},
		{
			name:        "unexpected error from DB",
			requestBody: `{"login":"new_user","password":"new_password"}`,
			newUser: models.Users{
				Login:    "new_user",
				Password: "new_password",
			},
			mockBehavior: func(s *mock_handler.MockServiceInterface, user models.Users) {
				s.EXPECT().CreateUser(gomock.Any(), user).Return("", apperrors.ErrCreateUser)
			},
			statusCode:      500,
			expectedBody:    "Something went wrong\n",
			expectedHeaders: "",
		},
		{
			name:        "Try to auth with empty username",
			requestBody: `{"login":"","password":"new_password"}`,
			newUser: models.Users{
				Login:    "",
				Password: "new_password",
			},
			mockBehavior: func(s *mock_handler.MockServiceInterface, user models.Users) {
				s.EXPECT().CreateUser(gomock.Any(), user).Return("", apperrors.ErrBadRequest)
			},
			statusCode:      400,
			expectedBody:    "Bad request\n",
			expectedHeaders: "",
		},
		{
			name:        "Try to auth with empty password",
			requestBody: `{"login":"new_user","password":""}`,
			newUser: models.Users{
				Login:    "new_user",
				Password: "",
			},
			mockBehavior: func(s *mock_handler.MockServiceInterface, user models.Users) {
				s.EXPECT().CreateUser(gomock.Any(), user).Return("", apperrors.ErrBadRequest)
			},
			statusCode:      400,
			expectedBody:    "Bad request\n",
			expectedHeaders: "",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//init controller
			c := gomock.NewController(t)
			defer c.Finish()

			//init mock service
			auth := mock_handler.NewMockServiceInterface(c)
			testCase.mockBehavior(auth, testCase.newUser)

			//init handler and router with mock
			handler := RegisterUser(auth)
			r := mux.NewRouter()
			r.Handle("/api/user/register", handler).Methods(http.MethodPost)

			//prepare test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBufferString(testCase.requestBody))
			req.Header.Add("Content-Type", "application/json")

			//execute request
			r.ServeHTTP(w, req)

			//assertion
			assert.Equal(t, testCase.statusCode, w.Code)
			assert.Equal(t, testCase.expectedBody, w.Body.String())
			assert.Equal(t, testCase.expectedHeaders, w.Header().Get("Authorization"))
		})
	}
}
