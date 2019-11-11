package openapi3

import (
	"fmt"
	"log"
	"net/url"
	"testing"
)

var basePathTests = []struct {
	v    string
	want string
}{
	{"https://{customerId}.saas-app.com:{port}/v2", "v2"},
}

func TestBasePath(t *testing.T) {
	for _, tt := range basePathTests {

		u1 := url.URL{
			Host: "https://{customerId}.saas-app.com:{port}",
		}
		fmt.Printf("HOST [%v]\n", u1.Host)

		tryUrl, err := url.Parse(tt.v)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("GOT PATH [%v]\n", tryUrl.Path)
		/*
			got := BasePath(tt.v)
			want := got
			if got != want {
				t.Errorf("openapi3.BasePath() Error: input [%v], want [%v], got [%v]",
					tt.v, tt.want, got)
			}*/
	}
}
