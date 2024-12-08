package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockService struct {
	mock.Mock
}

func (s *MockService) EmailWarning(guid string) error {
	return nil
}

func (s *MockService) InsertRT(guid, rToken string) error {
	args := s.Called(guid, rToken)
	return args.Error(0)
}

func (s *MockService) CompareRT(rtb, guid string) (bool, error) {
	args := s.Called(rtb, guid)
	return args.Bool(0), args.Error(1)
}

func TestMain(m *testing.M) {
	os.Setenv("ATEXPIRES", "60")
	os.Setenv("RTEXPIRES", "60")
	os.Setenv("JWT_SECRET", "jwt_secret")
	m.Run()
}

func TestGetTokens(t *testing.T) {
	tests := []struct {
		id             int
		guid           string
		wantStatusCode int
	}{
		{
			id:             1,
			guid:           "true",
			wantStatusCode: 200,
		},
		{
			id:             2,
			guid:           "",
			wantStatusCode: 400,
		},
		{
			id:             3,
			guid:           "false",
			wantStatusCode: 401,
		},
	}
	serviceMock := new(MockService)
	serviceMock.On("InsertRT", "true", mock.Anything).Return(nil)
	serviceMock.On("InsertRT", "false", mock.Anything).Return(sql.ErrNoRows)

	for _, test := range tests {
		fmt.Printf("Тест id: %v\n", test.id)
		url := fmt.Sprintf("/auth?guid=%s", test.guid)
		req := httptest.NewRequest("GET", url, nil)
		resReqorder := httptest.NewRecorder()
		handler := http.HandlerFunc(GetTokens(serviceMock))
		handler.ServeHTTP(resReqorder, req)

		require.Equal(t, test.wantStatusCode, resReqorder.Code, "статус код не соответствует ожидаемому")

		if test.guid == "true" && test.wantStatusCode == 200 {
			cook := resReqorder.Result().Cookies()
			require.Equal(t, 2, len(cook))
			at := cook[0]
			assert.Equal(t, "at", at.Name, "название куки не соответствует")
			assert.NotEmpty(t, at.Value, "пустой access токен")

			rt := cook[1]
			assert.Equal(t, "rt", rt.Name, "название куки не соответствует")
			assert.NotEmpty(t, rt.Value, "пустой access токен")
		}
	}
}

func TestRefreshTokens(t *testing.T) {
	tests := []struct {
		id             int
		guid           string
		wantStatusCode int
		aToken         string
		rToken         string
		rTokenB64      string
	}{
		{
			id:             1,
			guid:           "true",
			wantStatusCode: 200,
			aToken:         "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjE3MzM1ODQzMDEsIkhvc3QiOiJsb2NhbGhvc3Q6ODA4MCIsIkxpbmtTdHJpbmciOiIzNk1Vb3VzVEI3SEdEN2lvYkNCTVhEYS1XTHdtWVczNyJ9.fdGUdP-3E_JCJdg1CQkB31Zquz5M3Zwiji5N_8PJs9cDL9MftFLUPefmhRJK0y24HjQ1R3vx8hqVqWLdiYzgkg",
			rToken:         "OGU0MTEzYTZhZjEzMzA4YzVhMjI4Zjk5NGEyMWFhMGVkMGU0ZTcyNjVlZmNkMGFkYzlhNTQzNGM1Y2Y4MDMzZmlZemdrZw==",
			rTokenB64:      "8e4113a6af13308c5a228f994a21aa0ed0e4e7265efcd0adc9a5434c5cf8033fiYzgkg",
		},
		{
			id:             2,
			guid:           "",
			wantStatusCode: 400,
			aToken:         "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjE3MzM1ODQzMDEsIkhvc3QiOiJsb2NhbGhvc3Q6ODA4MCIsIkxpbmtTdHJpbmciOiIzNk1Vb3VzVEI3SEdEN2lvYkNCTVhEYS1XTHdtWVczNyJ9.fdGUdP-3E_JCJdg1CQkB31Zquz5M3Zwiji5N_8PJs9cDL9MftFLUPefmhRJK0y24HjQ1R3vx8hqVqWLdiYzgkg",
			rToken:         "OGU0MTEzYTZhZjEzMzA4YzVhMjI4Zjk5NGEyMWFhMGVkMGU0ZTcyNjVlZmNkMGFkYzlhNTQzNGM1Y2Y4MDMzZmlZemdrZw==",
			rTokenB64:      "8e4113a6af13308c5a228f994a21aa0ed0e4e7265efcd0adc9a5434c5cf8033fiYzgkg",
		},
		{
			id:             3,
			guid:           "true",
			wantStatusCode: 401,
			aToken:         "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjE3MzM1ODQzMDEsIkhvc3QiOiJsb2NhbGhvc3Q6ODA4MCIsIkxpbmtTdHJpbmciOiIzNk1Vb3VzVEI3SEdEN2lvYkNCTVhEYS1XTHdtWVczNyJ9.fdGUdP-3E_JCJdg1CQkB31Zquz5M3Zwiji5N_8PJs9cDL9MftFLUPefmhRJK0y24HjQ1R3vx8hqVqWLdiYzgkg",
			rToken:         "",
			rTokenB64:      "8e4113a6af13308c5a228f994a21aa0ed0e4e7265efcd0adc9a5434c5cf8033fiYzgkg",
		},
		{
			id:             4,
			guid:           "true",
			wantStatusCode: 401,
			aToken:         "",
			rToken:         "OGU0MTEzYTZhZjEzMzA4YzVhMjI4Zjk5NGEyMWFhMGVkMGU0ZTcyNjVlZmNkMGFkYzlhNTQzNGM1Y2Y4MDMzZmlZemdrZw==",
			rTokenB64:      "8e4113a6af13308c5a228f994a21aa0ed0e4e7265efcd0adc9a5434c5cf8033fiYzgkg",
		},
		{
			id:             5,
			guid:           "true",
			wantStatusCode: 401,
			aToken:         "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjE3MzM1ODQzMDEsIkhvc3QiOiJsb2NhbGhvc3Q6ODA4MCIsIkxpbmtTdHJpbmciOiIzNk1Vb3VzVEI3SEdEN2lvYkNCTVhEYS1XTHdtWVczNyJ9.fdGUdP-3E_JCJdg1CQkB31Zquz5M3Zwiji5N_8PJs9cDL9MftFLUPefmhRJK0y24HjQ1R3vx8hqVqWLdiYzgkg",
			rToken:         "M2YxOWYwMGIxM2Q4ZDlmZTZkZWMyNDdkM2I2N2UzMGQ5MTc5NjU2ZjExYmNkMmI5Mzk3ZjU4ZTVlNGY0NmE5ZnNhUU1lQQ==",
			rTokenB64:      "3f19f00b13d8d9fe6dec247d3b67e30d9179656f11bcd2b9397f58e5e4f46a9fsaQMeA",
		},
		{
			id:             6,
			guid:           "false",
			wantStatusCode: 401,
			aToken:         "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmVzQXQiOjE3MzM1ODQzMDEsIkhvc3QiOiJsb2NhbGhvc3Q6ODA4MCIsIkxpbmtTdHJpbmciOiIzNk1Vb3VzVEI3SEdEN2lvYkNCTVhEYS1XTHdtWVczNyJ9.fdGUdP-3E_JCJdg1CQkB31Zquz5M3Zwiji5N_8PJs9cDL9MftFLUPefmhRJK0y24HjQ1R3vx8hqVqWLdiYzgkg",
			rToken:         "OGU0MTEzYTZhZjEzMzA4YzVhMjI4Zjk5NGEyMWFhMGVkMGU0ZTcyNjVlZmNkMGFkYzlhNTQzNGM1Y2Y4MDMzZmlZemdrZw==",
			rTokenB64:      "8e4113a6af13308c5a228f994a21aa0ed0e4e7265efcd0adc9a5434c5cf8033fiYzgkg",
		},
	}
	serviceMock := new(MockService)
	serviceMock.On("CompareRT", "8e4113a6af13308c5a228f994a21aa0ed0e4e7265efcd0adc9a5434c5cf8033fiYzgkg", "true").Return(true, nil)
	serviceMock.On("CompareRT", "8e4113a6af13308c5a228f994a21aa0ed0e4e7265efcd0adc9a5434c5cf8033fiYzgkg", "false").Return(false, nil)
	serviceMock.On("CompareRT", "3f19f00b13d8d9fe6dec247d3b67e30d9179656f11bcd2b9397f58e5e4f46a9fsaQMeA", "true").Return(false, nil)

	serviceMock.On("InsertRT", "true", mock.Anything).Return(nil)
	serviceMock.On("InsertRT", "false", mock.Anything).Return(sql.ErrNoRows)

	for _, test := range tests {
		fmt.Printf("Тест id: %v\n", test.id)
		url := fmt.Sprintf("/refresh?guid=%s", test.guid)
		req := httptest.NewRequest("GET", url, nil)
		resReqorder := httptest.NewRecorder()

		if test.aToken != "" {
			atExp := time.Now().Add(30 * time.Second)
			req.AddCookie(&http.Cookie{
				Name:     "at",
				Value:    test.aToken,
				Expires:  atExp,
				HttpOnly: true,
			})
		}
		if test.rToken != "" {
			rtExp := time.Now().Add(30 * time.Second)
			req.AddCookie(&http.Cookie{
				Name:     "rt",
				Value:    test.rToken,
				Expires:  rtExp,
				HttpOnly: true,
			})
		}

		handler := http.HandlerFunc(RefreshTokens(serviceMock))
		handler.ServeHTTP(resReqorder, req)

		require.Equal(t, test.wantStatusCode, resReqorder.Code, "статус код не соответствует ожидаемому")
		if test.guid == "true" && test.wantStatusCode == 200 {
			cook := resReqorder.Result().Cookies()
			require.Equal(t, 2, len(cook))
			at := cook[0]
			assert.Equal(t, "at", at.Name, "название куки не соответствует")
			assert.NotEmpty(t, at.Value, "пустой access токен")

			rt := cook[1]
			assert.Equal(t, "rt", rt.Name, "название куки не соответствует")
			assert.NotEmpty(t, rt.Value, "пустой access токен")
		}
	}

}
