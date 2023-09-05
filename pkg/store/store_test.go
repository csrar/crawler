package store

import (
	"errors"
	"sync"
	"testing"

	mock_store "github.com/csrar/crawler/pkg/store/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestWasAlreadyVisited(t *testing.T) {
	cases := []struct {
		name           string
		inVisitedSite  string
		outReadBytes   []byte
		outTruncate    error
		outWriteAt     error
		ReadBytesTimes int
		TruncateTimes  int
		WriteAtTimes   int
		inWriteAt      []byte
		expectedResult bool
		expedtedErr    error
	}{
		{
			name:           "error due invalid data",
			inVisitedSite:  "",
			outReadBytes:   []byte("}"),
			expectedResult: false,
			expedtedErr:    errors.New("error decoding unmarshaling data: invalid character '}' looking for beginning of value"),
			ReadBytesTimes: 1,
			WriteAtTimes:   0,
			TruncateTimes:  0,
		},
		{
			name:           "error truncating data",
			inVisitedSite:  "mock-site.com",
			outReadBytes:   []byte(`{"sites":{"mocK-site.com": true}}`),
			expectedResult: false,
			expedtedErr:    errors.New("error truncating visited data mock-error"),
			outTruncate:    errors.New("mock-error"),
			TruncateTimes:  1,
			ReadBytesTimes: 1,
			WriteAtTimes:   0,
		},
		{
			name:           "error writting data",
			inVisitedSite:  "example-site.com",
			inWriteAt:      []byte(`{"sites":{"example-site-com":true,"mock-site-com":true}}`),
			outReadBytes:   []byte(`{"sites":{"mock-site-com": true}}`),
			expectedResult: false,
			expedtedErr:    errors.New("error writing visited data mock-error"),
			outWriteAt:     errors.New("mock-error"),
			TruncateTimes:  1,
			ReadBytesTimes: 1,
			WriteAtTimes:   1,
		},
		{
			name:           "Successful not visited",
			inVisitedSite:  "example-site.com",
			inWriteAt:      []byte(`{"sites":{"example-site-com":true,"mock-site-com":true}}`),
			outReadBytes:   []byte(`{"sites":{"mock-site-com": true}}`),
			expectedResult: false,
			TruncateTimes:  1,
			ReadBytesTimes: 1,
			WriteAtTimes:   1,
		},
		{
			name:           "Successful already visited",
			inVisitedSite:  "mock-site.com",
			inWriteAt:      []byte(`{"sites":{"example-site-com":true,"mock-site-com":true}}`),
			outReadBytes:   []byte(`{"sites":{"mock-site-com": true}}`),
			expectedResult: true,
			ReadBytesTimes: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var mu sync.Mutex
			fileMock := mock_store.NewMockImFile(ctrl)

			fileMock.EXPECT().Bytes().Return(tc.outReadBytes).Times(tc.ReadBytesTimes)
			fileMock.EXPECT().Truncate(int64(0)).Return(tc.outTruncate).Times(tc.TruncateTimes)
			fileMock.EXPECT().WriteAt(tc.inWriteAt, int64(0)).Return(int(0), tc.outWriteAt).Times(tc.WriteAtTimes)

			store := NewMemfileStore(&mu, fileMock)
			result, err := store.WasAlreadyVisited(tc.inVisitedSite)

			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expedtedErr, err)

		})

	}
}
