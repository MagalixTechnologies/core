package cache

import (
	"net/http"
	"net/url"
	"testing"
)

type input struct {
	Request http.Request
	Match   URLMatch
}

type testCaseURLMatch struct {
	name     string
	input    input
	expected bool
}

func TestMatchURL(t *testing.T) {
	//'/analytics/
	cases := []testCaseURLMatch{
		{
			name: "match_normal_path_with_method",
			input: input{
				Request: http.Request{URL: &url.URL{
					Scheme: "http",
					Host:   "hamda.com",
					Path:   "analytics/query",
				},
					Method: http.MethodPost,
				},
				Match: URLMatch{"analytics/query", "POST"},
			},
			expected: true,
		}, {
			name: "match_normal_path_with_wrong_method",
			input: input{
				Request: http.Request{URL: &url.URL{
					Scheme: "http",
					Host:   "hamda.com",
					Path:   "analytics/query",
				},
					Method: http.MethodPost,
				},
				Match: URLMatch{"analytics/query", http.MethodGet},
			},
			expected: false,
		},
		{
			name: "match_a_request_with_arg",
			input: input{
				Request: http.Request{URL: &url.URL{
					Scheme: "http",
					Host:   "hamda.com",
					Path:   "analytics/12345/query",
				},
					Method: http.MethodPost,
				},
				Match: URLMatch{"analytics/:arg/query", "POST"},
			},
			expected: true,
		},
		{
			name: "match_a_request_with_arg_2",
			input: input{
				Request: http.Request{URL: &url.URL{
					Scheme: "http",
					Host:   "hamda.com",
					Path:   "analytics/12345/query",
				},
					Method: http.MethodPost,
				},
				Match: URLMatch{"analytics/:arg/:arg", "POST"},
			},
			expected: true,
		},
		{
			name: "match_a_request_with_arg_3",
			input: input{
				Request: http.Request{URL: &url.URL{
					Scheme: "http",
					Host:   "hamda.com",
					Path:   "analytics/12345/query",
				},
					Method: http.MethodPost,
				},
				Match: URLMatch{"analytics/:arg/:arg/:arg", "POST"},
			},
			expected: false,
		}, {
			name: "match_a_request_with_arg_4",
			input: input{
				Request: http.Request{URL: &url.URL{
					Scheme: "http",
					Host:   "hamda.com",
					Path:   "analytics/query",
				},
					Method: http.MethodPost,
				},
				Match: URLMatch{"analytics/:arg/:arg/:arg", "POST"},
			},
			expected: false,
		},
		{
			name: "match_a_request_with_arg_5",
			input: input{
				Request: http.Request{URL: &url.URL{
					Scheme: "http",
					Host:   "hamda.com",
					Path:   "analytics/query",
				},
					Method: http.MethodPost,
				},
				Match: URLMatch{"analytics/:arg", http.MethodGet},
			},
			expected: false,
		},
	}

	// analytics/{analyticId}/execute
	for i := range cases {
		c := cases[i]
		t.Run(c.name, func(t *testing.T) {
			got := c.input.Match.MatchRequest(&c.input.Request)
			if c.expected != got {
				t.Errorf(
					"test (%v) expected %v but got %v on %v",
					c.name,
					c.expected,
					got,
					c.input,
				)
			}

		})
	}
}
