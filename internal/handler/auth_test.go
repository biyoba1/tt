package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"valera/internal/service/mocks"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"valera/internal/service/auth"
	"valera/models"
)

func TestSignUp(t *testing.T) {
	t.Parallel()

	type authServiceMockFunc func(mc *minimock.Controller) auth.AuthService

	type args struct {
		ctx context.Context
		req *models.LoginReq
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		guid         = gofakeit.UUID()
		clientIP     = "127.0.0.1"
		accessToken  = gofakeit.UUID()
		refreshToken = gofakeit.UUID()

		serviceErr = fmt.Errorf("service error")

		req = &models.LoginReq{
			Guid: guid,
		}

		authResponse = &models.AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		wantCode        int
		wantBody        map[string]interface{}
		authServiceMock authServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			wantCode: http.StatusCreated,
			wantBody: map[string]interface{}{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
			authServiceMock: func(mc *minimock.Controller) auth.AuthService {
				mock := mocks.NewAuthServiceMock(mc)
				mock.LoginMock.Expect(guid, clientIP).Return(authResponse, nil)
				return mock
			},
		},
		{
			name: "invalid JSON",
			args: args{
				ctx: ctx,
				req: nil,
			},
			wantCode: http.StatusBadRequest,
			wantBody: map[string]interface{}{
				"error": "invalid request body",
			},
			authServiceMock: func(mc *minimock.Controller) auth.AuthService {
				return mocks.NewAuthServiceMock(mc)
			},
		},
		{
			name: "missing Guid",
			args: args{
				ctx: ctx,
				req: &models.LoginReq{},
			},
			wantCode: http.StatusBadRequest,
			wantBody: map[string]interface{}{
				"error": "guid is required",
			},
			authServiceMock: func(mc *minimock.Controller) auth.AuthService {
				return mocks.NewAuthServiceMock(mc)
			},
		},
		{
			name: "service error",
			args: args{
				ctx: ctx,
				req: req,
			},
			wantCode: http.StatusInternalServerError,
			wantBody: map[string]interface{}{
				"error": serviceErr.Error(),
			},
			authServiceMock: func(mc *minimock.Controller) auth.AuthService {
				mock := mocks.NewAuthServiceMock(mc)
				mock.LoginMock.Expect(guid, clientIP).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			authServiceMock := tt.authServiceMock(mc)
			handler := NewHandler(authServiceMock)
			var reqBody []byte
			if tt.args.req != nil {
				var err error
				reqBody, err = json.Marshal(tt.args.req)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
			req.RemoteAddr = clientIP
			w := httptest.NewRecorder()

			handler.signUp(w, req)
			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)
			require.Equal(t, tt.wantBody, responseBody)
		})
	}
}
