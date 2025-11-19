package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"testovoe/internal/domain"
	"testovoe/internal/http/handlers/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHTTPHandler_HandleRequest_Success(t *testing.T) {
	mockUseCase := mocks.NewMockUseCase(t)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mockUseCase.EXPECT().
		PutNumber(mock.Anything, 42).
		Return(nil).
		Once()

	mockUseCase.EXPECT().
		GetSlices(mock.Anything).
		Return([]int{1, 2, 42}, nil).
		Once()

	handler := &HTTPHandler{
		useCase: mockUseCase,
		log:     logger,
	}

	requestBody := domain.UserNum{Num: 42}
	body, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/handle", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handlerFunc := handler.HandleRequest(context.Background())
	handlerFunc(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response []int
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 42}, response)
}

func TestHTTPHandler_HandleRequest_PutNumberError(t *testing.T) {
	mockUseCase := mocks.NewMockUseCase(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mockUseCase.EXPECT().
		PutNumber(mock.Anything, 42).
		Return(errors.New("database error")).
		Once()

	handler := &HTTPHandler{
		useCase: mockUseCase,
		log:     logger,
	}

	requestBody := domain.UserNum{Num: 42}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/handle", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handlerFunc := handler.HandleRequest(context.Background())
	handlerFunc(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHTTPHandler_HandleRequest_GetSlicesError(t *testing.T) {
	mockUseCase := mocks.NewMockUseCase(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mockUseCase.EXPECT().
		PutNumber(mock.Anything, 42).
		Return(nil).
		Once()

	mockUseCase.EXPECT().
		GetSlices(mock.Anything).
		Return(nil, errors.New("database error")).
		Once()

	handler := &HTTPHandler{
		useCase: mockUseCase,
		log:     logger,
	}

	requestBody := domain.UserNum{Num: 42}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/handle", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handlerFunc := handler.HandleRequest(context.Background())
	handlerFunc(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHTTPHandler_HandleRequest_InvalidJSON(t *testing.T) {
	mockUseCase := mocks.NewMockUseCase(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	handler := &HTTPHandler{
		useCase: mockUseCase,
		log:     logger,
	}

	req := httptest.NewRequest(http.MethodPost, "/api/handle", bytes.NewBufferString("invalid json{"))
	w := httptest.NewRecorder()

	handlerFunc := handler.HandleRequest(context.Background())
	handlerFunc(w, req)
}

func TestHTTPHandler_HandleRequest_EmptyBody(t *testing.T) {
	mockUseCase := mocks.NewMockUseCase(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	handler := &HTTPHandler{
		useCase: mockUseCase,
		log:     logger,
	}

	req := httptest.NewRequest(http.MethodPost, "/api/handle", bytes.NewBufferString(""))
	w := httptest.NewRecorder()

	handlerFunc := handler.HandleRequest(context.Background())
	handlerFunc(w, req)
}

func TestHTTPHandler_HandleRequest_DifferentNumbers(t *testing.T) {
	testCases := []struct {
		name           string
		inputNumber    int
		expectedSlices []int
	}{
		{
			name:           "positive number",
			inputNumber:    100,
			expectedSlices: []int{1, 2, 100},
		},
		{
			name:           "zero",
			inputNumber:    0,
			expectedSlices: []int{0},
		},
		{
			name:           "negative number",
			inputNumber:    -5,
			expectedSlices: []int{-5, -4, -3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUseCase := mocks.NewMockUseCase(t)
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

			mockUseCase.EXPECT().
				PutNumber(mock.Anything, tc.inputNumber).
				Return(nil).
				Once()

			mockUseCase.EXPECT().
				GetSlices(mock.Anything).
				Return(tc.expectedSlices, nil).
				Once()

			handler := &HTTPHandler{
				useCase: mockUseCase,
				log:     logger,
			}

			requestBody := domain.UserNum{Num: tc.inputNumber}
			body, _ := json.Marshal(requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/handle", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handlerFunc := handler.HandleRequest(context.Background())
			handlerFunc(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response []int
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedSlices, response)
		})
	}
}

func TestHTTPHandler_HandleRequest_ContextPassed(t *testing.T) {
	mockUseCase := mocks.NewMockUseCase(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	var capturedCtx context.Context

	mockUseCase.EXPECT().
		PutNumber(mock.Anything, 42).
		Run(func(ctx context.Context, num int) {
			capturedCtx = ctx
		}).
		Return(nil).
		Once()

	mockUseCase.EXPECT().
		GetSlices(mock.Anything).
		Return([]int{42}, nil).
		Once()

	handler := &HTTPHandler{
		useCase: mockUseCase,
		log:     logger,
	}

	requestBody := domain.UserNum{Num: 42}
	body, _ := json.Marshal(requestBody)

	ctx := context.WithValue(context.Background(), "test-key", "test-value")
	req := httptest.NewRequest(http.MethodPost, "/api/handle", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handlerFunc := handler.HandleRequest(ctx)
	handlerFunc(w, req)

	assert.NotNil(t, capturedCtx)
	assert.Equal(t, "test-value", capturedCtx.Value("test-key"))
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
