package hibp

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"testing"
)

func TestExampleClientPwned(t *testing.T) {
	c := Client{HTTPClient: http.DefaultClient}
	randomPassword := strconv.FormatFloat(rand.Float64(), 'f', -1, 64)
	fmt.Println(c.Pwned(context.Background(), "letmein"))
	fmt.Println(c.Pwned(context.Background(), randomPassword))
	// Output:
	// true <nil>
	// false <nil>
}

func TestCheckLineMatch(t *testing.T) {
	tests := []struct {
		name          string
		sha1HexSuffix string
		line          string
		want          bool
		wantErr       bool
	}{
		{
			name:          "match",
			sha1HexSuffix: "suffix",
			line:          "suffix:1",
			want:          true,
			wantErr:       false,
		},
		{
			name:          "different suffix",
			sha1HexSuffix: "suffix",
			line:          "notsuffix:1",
			want:          false,
			wantErr:       false,
		},
		{
			name:          "padding line",
			sha1HexSuffix: "suffix",
			line:          "suffix:0",
			want:          false,
			wantErr:       false,
		},
		{
			name:          "malformed line no count",
			sha1HexSuffix: "suffix",
			line:          "suffix",
			want:          false,
			wantErr:       true,
		},
		{
			name:          "malformed line bad count",
			sha1HexSuffix: "suffix",
			line:          "suffix:badcount",
			want:          false,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := checkLineMatch(tt.sha1HexSuffix, tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkLineMatch() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("checkLineMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
