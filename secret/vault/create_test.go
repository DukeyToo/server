// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
)

func TestVault_Create_Org(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.PUT("/v1/secret/:type/:org/:name", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})
	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	one := int64(1)
	org := "foo"
	repo := "*"
	team := ""
	name := "bar"
	value := "baz"
	typee := "org"
	arr := []string{"foo", "bar"}
	sec := &library.Secret{
		ID:     &one,
		Org:    &org,
		Repo:   &repo,
		Team:   &team,
		Name:   &name,
		Value:  &value,
		Type:   &typee,
		Images: &arr,
		Events: &arr,
	}

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create(typee, org, repo, sec)

	if resp.Code != http.StatusOK {
		t.Errorf("Create returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Create returned err: %v", err)
	}
}

func TestVault_Create_Repo(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.PUT("/v1/secret/:type/:org/:repo/:name", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})
	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	one := int64(1)
	org := "foo"
	repo := "bar"
	team := ""
	name := "baz"
	value := "foob"
	typee := "repo"
	arr := []string{"foo", "bar"}
	sec := &library.Secret{
		ID:     &one,
		Org:    &org,
		Repo:   &repo,
		Team:   &team,
		Name:   &name,
		Value:  &value,
		Type:   &typee,
		Images: &arr,
		Events: &arr,
	}

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create(typee, org, repo, sec)

	if resp.Code != http.StatusOK {
		t.Errorf("Create returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Create returned err: %v", err)
	}
}

func TestVault_Create_Shared(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.PUT("/v1/secret/:type/:org/:team/:name", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})
	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	one := int64(1)
	org := "foo"
	repo := ""
	team := "bar"
	name := "baz"
	value := "foob"
	typee := "shared"
	arr := []string{"foo", "bar"}
	sec := &library.Secret{
		ID:     &one,
		Org:    &org,
		Repo:   &repo,
		Team:   &team,
		Name:   &name,
		Value:  &value,
		Type:   &typee,
		Images: &arr,
		Events: &arr,
	}

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create(typee, org, team, sec)

	if resp.Code != http.StatusOK {
		t.Errorf("Create returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Create returned err: %v", err)
	}
}

func TestVault_Create_InvalidSecret(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.PUT("/v1/secret/:type/:org/:repo/:name", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})
	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	one := int64(1)
	org := "foo"
	repo := "bar"
	name := "baz"
	value := ""
	typee := "repo"
	arr := []string{"foo", "bar"}
	sec := &library.Secret{
		ID:     &one,
		Org:    &org,
		Repo:   &repo,
		Name:   &name,
		Value:  &value,
		Type:   &typee,
		Images: &arr,
		Events: &arr,
	}

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create(typee, org, repo, sec)

	if resp.Code != http.StatusOK {
		t.Errorf("Create returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Errorf("Create should have returned err")
	}
}

func TestVault_Create_InvalidType(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	repo := "bar"
	name := "baz"
	value := "foob"
	typee := "invalid"
	arr := []string{"foo", "bar"}
	sec := &library.Secret{
		ID:     &one,
		Org:    &org,
		Repo:   &repo,
		Name:   &name,
		Value:  &value,
		Type:   &typee,
		Images: &arr,
		Events: &arr,
	}

	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	defer fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create(typee, org, repo, sec)
	if err == nil {
		t.Errorf("Create should have returned err")
	}
}

func TestVault_Create_ClosedServer(t *testing.T) {
	// setup types
	one := int64(1)
	org := "foo"
	repo := "bar"
	team := ""
	name := "baz"
	value := "foob"
	typee := "repo"
	arr := []string{"foo", "bar"}
	sec := &library.Secret{
		ID:     &one,
		Org:    &org,
		Repo:   &repo,
		Team:   &team,
		Name:   &name,
		Value:  &value,
		Type:   &typee,
		Images: &arr,
		Events: &arr,
	}

	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create(typee, org, repo, sec)
	if err == nil {
		t.Errorf("Create should have returned err")
	}
}