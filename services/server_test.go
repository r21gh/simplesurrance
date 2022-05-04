package services

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestGetSliceOfValuesFromFile(t *testing.T) {
	var testCases = []struct {
		fileName string
		expected []int64
	}{
		{"../test_data/test_file_1.txt", []int64{1, int64(1651658067)}},
		{"../test_data/test_file_2.txt", []int64{4, int64(1651658047)}},
	}

	for _, testCase := range testCases {
		actual, err := getSliceOfValuesFromFile(testCase.fileName)
		if err != nil {
			t.Errorf("Error: %s", err)
		}
		if len(actual) != len(testCase.expected) {
			t.Errorf("Expected %v, got %v", testCase.expected, actual)
		}
		// check if the test case is equal to the actual
		for i := 0; i < len(testCase.expected); i++ {
			if actual[i] != testCase.expected[i] {
				t.Errorf("Filename: %q, Expected %v, got %v", testCase.fileName, testCase.expected, actual)
			}
		}

		// check if the test case type is equal to the actual type
		if reflect.TypeOf(actual) != reflect.TypeOf(testCase.expected) {
			t.Errorf("Filename: %q, Expected %v, got %v", testCase.fileName, testCase.expected, actual)
		}

	}
}

func Test_apiHandler(t *testing.T) {
	tests := []struct {
		name string
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApiHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApiHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_assertCounterPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "assertCounterPath",
			args: args{
				path: "/counter",
			},
			wantErr: true,
		},
		{
			name: "assertCounterPath",
			args: args{
				path: "/api/v1/counter/",
			},
			wantErr: false,
		},
		{
			name: "assertCounterPath",
			args: args{
				path: "api/counter/1",
			},
			wantErr: true,
		},
		{
			name: "assertCounterPath",
			args: args{
				path: "api/v1/counter/",
			},
			wantErr: true,
		},
		{
			name: "assertCounterPath",
			args: args{
				path: "/api/v2/counter/1/2",
			},
			wantErr: true,
		},
		{
			name: "assertCounterPath",
			args: args{
				path: "/api/v1//counter/1/2/",
			},
			wantErr: true,
		},
		{
			name: "assertCounterPath",
			args: args{
				path: "/counter/1/2/3",
			},
			wantErr: true,
		},
		{
			name: "assertCounterPath",
			args: args{
				path: "//api/v1/counter/",
			},
			wantErr: true,
		},
		{
			name: "assertCounterPath",
			args: args{
				path: "/api/v1/counter//",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := assertCounterPath(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("assertCounterPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_addTrailingSlash(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "addTrailingSlash",
			args: args{
				path: "/counter",
			},
			want: "/counter/",
		},
		{
			name: "addTrailingSlash",
			args: args{
				path: "/counter/",
			},
			want: "/counter/",
		},
		{
			name: "addTrailingSlash",
			args: args{
				path: "/counter/1",
			},
			want: "/counter/1/",
		},
		{
			name: "addTrailingSlash",
			args: args{
				path: "/counter/1/",
			},
			want: "/counter/1/",
		},
		{
			name: "addTrailingSlash",
			args: args{
				path: "/counter/1/2",
			},
			want: "/counter/1/2/",
		},
		{
			name: "addTrailingSlash",
			args: args{
				path: "/counter/1/2/",
			},
			want: "/counter/1/2/",
		},
		{
			name: "addTrailingSlash",
			args: args{
				path: "/counter/1/2/3",
			},
			want: "/counter/1/2/3/",
		},
		{
			name: "addTrailingSlash",
			args: args{
				path: "/counter/1/2/3/",
			},
			want: "/counter/1/2/3/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addTrailingSlash(tt.args.path); got != tt.want {
				t.Errorf("addTrailingSlash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateTimestampWithWindow(t *testing.T) {
	type args struct {
		window      float64
		currentTime int64
		timestamp   int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "validateTimestampWithWindow",
			args: args{
				window:      60.0,
				currentTime: 1651661792,
				timestamp:   1651661882,
			},
			want: false,
		},
		{
			name: "validateTimestampWithWindow",
			args: args{
				window:      60.0,
				currentTime: 1651661703,
				timestamp:   1651661764,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateTimestampWithWindow(tt.args.window, tt.args.currentTime, tt.args.timestamp); got != tt.want {
				t.Errorf("validateTimestampWithWindow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_counterObjectValidator(t *testing.T) {
	type args struct {
		c           *counter
		currentTime time.Time
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "counterObjectValidator",
			args: args{
				c: &counter{
					count:     1,
					timeStamp: time.Now(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counterObjectValidator(tt.args.c, tt.args.currentTime)
		})
	}
}
