package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want  int
	}{
		{count: 0, want: 0},
		{count: 1, want: 1},
		{count: 2, want: 2},
		{count: 100, want: min(len(cafeList["moscow"]), 100)},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city=moscow&count="+strconv.Itoa(v.count), nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)

		body := strings.TrimSpace(response.Body.String())
		if body == "" {
			assert.Equal(t, 0, v.want)
			continue
		}

		got := len(strings.Split(body, ","))
		assert.Equal(t, v.want, got)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
	}{
		{search: "фасоль", wantCount: 0},
		{search: "кофе", wantCount: 2},
		{search: "вилка", wantCount: 1},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city=moscow&search="+v.search, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)

		body := strings.TrimSpace(response.Body.String())
		if body == "" {
			assert.Equal(t, 0, v.wantCount)
			continue
		}

		cafes := strings.Split(body, ",")
		assert.Equal(t, v.wantCount, len(cafes))

		for _, cafe := range cafes {
			assert.True(t, strings.Contains(strings.ToLower(cafe), strings.ToLower(v.search)))
		}
	}
}
