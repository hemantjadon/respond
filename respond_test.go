package respond_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hemantjadon/respond"
)

func TestWith(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		status int
		data   []byte
	}
	tests := []struct {
		name       string
		recorder   *httptest.ResponseRecorder
		args       args
		wantErr    bool
		wantStatus int
		wantBody   []byte
	}{
		{
			name:     "nil data ok status",
			recorder: httptest.NewRecorder(),
			args: args{
				w:      httptest.NewRecorder(),
				status: http.StatusOK,
				data:   nil,
			},
			wantStatus: http.StatusOK,
			wantBody:   nil,
		},
		{
			name: "nil data different status",
			args: args{
				w:      httptest.NewRecorder(),
				status: http.StatusBadGateway,
				data:   nil,
			},
			wantStatus: http.StatusBadGateway,
			wantBody:   nil,
		},
		{
			name: "simple data",
			args: args{
				w:      httptest.NewRecorder(),
				status: http.StatusOK,
				data:   []byte(`hello`),
			},
			wantStatus: http.StatusOK,
			wantBody:   []byte(`hello`),
		},
		{
			name: "error writer",
			args: args{
				w:      &errResponseRecorder{},
				status: http.StatusOK,
				data:   nil,
			},
			wantErr: true,
		},
	}
	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := respond.With(tt.args.w, tt.args.status, tt.args.data); (err != nil) != tt.wantErr {
				t.Fatalf("With() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr {
				return
			}
			rec := tt.args.w.(*httptest.ResponseRecorder)
			if rec.Code != tt.wantStatus {
				t.Fatalf("code = %#v, wantCode = %#v", rec.Code, tt.wantStatus)
			}
			if !bytes.Equal(rec.Body.Bytes(), tt.wantBody) {
				t.Fatalf("body = %#v, wantBody = %#v", rec.Body.Bytes(), tt.wantBody)
			}
		})
	}
}

func TestWithJSON(t *testing.T) {
	const jsonContentType = "application/json; utf-8"

	type dummy struct {
		F1 string `json:"f_1"`
		F2 int    `json:"f_2"`
	}

	type args struct {
		w      http.ResponseWriter
		status int
		data   interface{}
	}
	tests := []struct {
		name       string
		recorder   *httptest.ResponseRecorder
		args       args
		wantErr    bool
		wantStatus int
		wantBody   []byte
	}{
		{
			name:     "nil data ok status",
			recorder: httptest.NewRecorder(),
			args: args{
				w:      httptest.NewRecorder(),
				status: http.StatusOK,
				data:   nil,
			},
			wantStatus: http.StatusOK,
			wantBody:   nil,
		},
		{
			name: "nil data different status",
			args: args{
				w:      httptest.NewRecorder(),
				status: http.StatusBadGateway,
				data:   nil,
			},
			wantStatus: http.StatusBadGateway,
			wantBody:   nil,
		},
		{
			name: "simple map data",
			args: args{
				w:      httptest.NewRecorder(),
				status: http.StatusOK,
				data:   map[string]interface{}{"f_1": "v1", "f_2": 2},
			},
			wantStatus: http.StatusOK,
			wantBody:   ignoreSerErr(json.Marshal(map[string]interface{}{"f_1": "v1", "f_2": 2})),
		},
		{
			name: "simple struct data",
			args: args{
				w:      httptest.NewRecorder(),
				status: http.StatusOK,
				data:   dummy{F1: "v1", F2: 2},
			},
			wantStatus: http.StatusOK,
			wantBody:   ignoreSerErr(json.Marshal(dummy{F1: "v1", F2: 2})),
		},
		{
			name: "simple pointer to struct data",
			args: args{
				w:      httptest.NewRecorder(),
				status: http.StatusOK,
				data:   &dummy{F1: "v1", F2: 2},
			},
			wantStatus: http.StatusOK,
			wantBody:   ignoreSerErr(json.Marshal(dummy{F1: "v1", F2: 2})),
		},
		{
			name: "non json data",
			args: args{
				w:      httptest.NewRecorder(),
				status: http.StatusOK,
				data:   func() {},
			},
			wantErr: true,
		},
		{
			name: "error writer",
			args: args{
				w:      &errResponseRecorder{},
				status: http.StatusOK,
				data:   nil,
			},
			wantErr: true,
		},
	}
	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := respond.WithJSON(tt.args.w, tt.args.status, tt.args.data); (err != nil) != tt.wantErr {
				t.Fatalf("With() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr {
				return
			}
			rec := tt.args.w.(*httptest.ResponseRecorder)
			if rec.Code != tt.wantStatus {
				t.Fatalf("code = %#v, wantCode = %#v", rec.Code, tt.wantStatus)
			}
			if !bytes.Equal(rec.Body.Bytes(), tt.wantBody) {
				t.Fatalf("body = %#v, wantBody = %#v", rec.Body.Bytes(), tt.wantBody)
			}
			if tt.wantBody != nil {
				contentType := rec.Header().Get("Content-Type")
				if contentType != jsonContentType {
					t.Fatalf("header Content-Type = %#v, want header Content-Type = %#v", contentType, jsonContentType)
				}
			}
		})
	}
}

type errResponseRecorder struct {
	httptest.ResponseRecorder
}

func (w errResponseRecorder) Write([]byte) (int, error) {
	return 0, fmt.Errorf("error")
}

func ignoreSerErr(data []byte, err error) []byte {
	return data
}
