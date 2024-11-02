package httputil

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetPagination(t *testing.T) {
	tests := []struct {
		name     string
		query    url.Values
		expected Pagination
	}{
		{
			name:  "Empty Query Parameters",
			query: url.Values{},
			expected: Pagination{
				PerPage:   0,
				Cursor:    0,
				Next:      false,
				FirstPage: true,
			},
		},
		{
			name: "Valid Parameters",
			query: url.Values{
				"perPage": []string{"20"},
				"cursor":  []string{"10"},
				"next":    []string{"true"},
			},
			expected: Pagination{
				PerPage:   20,
				Cursor:    10,
				Next:      true,
				FirstPage: false,
			},
		},
		{
			name: "Exceeds Maximum PerPage",
			query: url.Values{
				"perPage": []string{"150"},
				"cursor":  []string{"5"},
			},
			expected: Pagination{
				PerPage:   100,
				Cursor:    5,
				Next:      false,
				FirstPage: false,
			},
		},
		{
			name: "Invalid Numbers",
			query: url.Values{
				"perPage": []string{"invalid"},
				"cursor":  []string{"invalid"},
				"next":    []string{"true"},
			},
			expected: Pagination{
				PerPage:   0,
				Cursor:    0,
				Next:      true,
				FirstPage: false,
			},
		},
		{
			name: "Negative Cursor",
			query: url.Values{
				"perPage": []string{"20"},
				"cursor":  []string{"-5"},
				"next":    []string{"false"},
			},
			expected: Pagination{
				PerPage:   20,
				Cursor:    -5,
				Next:      false,
				FirstPage: true,
			},
		},
		{
			name: "Zero PerPage",
			query: url.Values{
				"perPage": []string{"0"},
				"cursor":  []string{"5"},
				"next":    []string{"true"},
			},
			expected: Pagination{
				PerPage:   0,
				Cursor:    5,
				Next:      true,
				FirstPage: false,
			},
		},
		{
			name: "Multiple Values (Should Use First)",
			query: url.Values{
				"perPage": []string{"20", "30"},
				"cursor":  []string{"5", "10"},
				"next":    []string{"true", "false"},
			},
			expected: Pagination{
				PerPage:   20,
				Cursor:    5,
				Next:      true,
				FirstPage: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPagination(tt.query); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("GetPagination() = %v, want %v", got, tt.expected)
			}
		})
	}
}
