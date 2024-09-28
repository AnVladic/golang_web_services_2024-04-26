package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TestCase struct {
	Request SearchRequest
	Result  *SearchResponse
	IsError bool
}

func TimeoutSearchServer(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	SearchServer(w, r)
}

func WrongResponseSearchServer(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("wrong response json text"))
}

func WrongBadRequestErrorSearchClient(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte("wrong response json text"))
}

func StatusInternalServerErrorSearchClient(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func TestSearchClient_FindUsers(t *testing.T) {
	cases := []TestCase{
		{
			Request: SearchRequest{
				Query: "Hilda",
				Limit: 1,
			},
			Result: &SearchResponse{
				Users:    []User{{Id: 1, Name: "Hilda Mayer"}},
				NextPage: false,
			},
			IsError: false,
		},
		{
			Request: SearchRequest{
				Query: "Brooks",
				Limit: 30,
			},
			Result: &SearchResponse{
				Users:    []User{{Id: 2, Name: "Brooks Aguilar"}},
				NextPage: false,
			},
			IsError: false,
		},
		{
			Request: SearchRequest{
				Limit: -1,
			},
			Result:  nil,
			IsError: true,
		},
		{
			Request: SearchRequest{
				Limit:  12,
				Offset: -1,
			},
			IsError: true,
		},
		{
			Request: SearchRequest{
				Limit:   2,
				Offset:  0,
				OrderBy: OrderByDesc,
			},
			Result: &SearchResponse{
				Users:    []User{{Id: 13, Name: "Whitley Davidson"}, {Id: 33, Name: "Twila Snow"}},
				NextPage: true,
			},
			IsError: false,
		},
		{
			Request: SearchRequest{
				OrderBy: 3,
			},
			IsError: true,
		},
		{
			Request: SearchRequest{
				OrderField: "WrongName",
			},
			IsError: true,
		},
	}

	ParseDataset()
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	searchClient := SearchClient{
		AccessToken: "MyAccessToken",
		URL:         ts.URL,
	}

	for caseNum, item := range cases {
		result, err := searchClient.FindUsers(item.Request)

		if err != nil && !item.IsError {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && item.IsError {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}

		if item.Result == nil && result == nil {
			continue
		}

		if item.Result.NextPage != result.NextPage {
			t.Errorf(
				"[%d] wrong NextPage result, expected %#v, got %#v",
				caseNum, item.Result.NextPage, result.NextPage)
		}
		if len(item.Result.Users) != len(result.Users) {
			t.Errorf(
				"[%d] wrong Users length, expected %d, got %d",
				caseNum, len(item.Result.Users), len(result.Users))
		}
		for i, user := range item.Result.Users {
			resultUser := result.Users[i]
			if user.Name != resultUser.Name || user.Id != resultUser.Id {
				t.Errorf(
					"[%d] wrong User %d, expected %v, got %v", caseNum, i,
					user, result.Users[i])
			}
		}
	}
	searchClient.AccessToken = "wrong access token"
	_, err := searchClient.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf(
			"Excepcted StatusUnauthorized %s", err)
	}

	searchClient.URL = ""
	_, err = searchClient.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf(
			"unknown error %s", err)
	}
	ts.Close()
}

func TestTimeoutSearchClient_FindUsers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(TimeoutSearchServer))
	searchClient := SearchClient{
		URL: ts.URL,
	}
	_, err := searchClient.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf(
			"Excepcted Timeout")
	}
	ts.Close()
}

func TestWrongResponseSearchClient_FindUsers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(WrongResponseSearchServer))
	searchClient := SearchClient{
		URL: ts.URL,
	}
	_, err := searchClient.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf(
			"cant unpack result json")
	}
	ts.Close()
}

func TestWrongBadRequestErrorSearchClient_FindUsers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(WrongBadRequestErrorSearchClient))
	searchClient := SearchClient{
		URL: ts.URL,
	}
	_, err := searchClient.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf(
			"cant unpack result json")
	}
	ts.Close()
}

func TestStatusInternalServerErrorSearchClient_FindUsers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(StatusInternalServerErrorSearchClient))
	searchClient := SearchClient{
		URL: ts.URL,
	}
	_, err := searchClient.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf(
			"Excepcted SearchServer fatal error")
	}
	ts.Close()
}
