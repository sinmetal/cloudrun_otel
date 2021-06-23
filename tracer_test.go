package main

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetCloudTraceHeader(t *testing.T) {
	cases := []struct {
		name   string
		header string
		want   *CloudTraceHeader
	}{
		{"p1", "105445aa7843bc8bf206b12000100000/1;o=1", &CloudTraceHeader{"105445aa7843bc8bf206b12000100000", "1", 1}},
		{"p2", "105445aa7843bc8bf206b12000100000/1;o=0", &CloudTraceHeader{"105445aa7843bc8bf206b12000100000", "1", 0}},
		{"p3", "105445aa7843bc8bf206b12000100000/1", &CloudTraceHeader{"105445aa7843bc8bf206b12000100000", "1", -1}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "", nil)
			if err != nil {
				t.Fatal(err)
			}
			r.Header.Set("X-Cloud-Trace-Context", tt.header)
			got, err := GetCloudTraceHeader(r)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Error(diff)
			}
		})
	}
}
