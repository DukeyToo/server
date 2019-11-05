// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/step"
	"github.com/go-vela/server/router/middleware/user"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestMiddleware_Logger(t *testing.T) {
	// setup types
	num := 1
	num64 := int64(num)
	foo := "foo"
	bar := "bar"
	foobar := "foo/bar"
	b := &library.Build{ID: &num64, RepoID: &num64, Number: &num}
	r := &library.Repo{ID: &num64, UserID: &num64, Org: &foo, Name: &bar, FullName: &foobar}
	s := &library.Step{ID: &num64, RepoID: &num64, BuildID: &num64, Number: &num, Name: &foo}
	u := &library.User{ID: &num64, Name: &foo, Token: &bar}
	payload, _ := json.Marshal(`{"foo": "bar"}`)
	wantLevel := logrus.InfoLevel
	wantMessage := ""
	logger, hook := test.NewNullLogger()
	defer hook.Reset()

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodPost, "/foobar", bytes.NewBuffer(payload))

	// setup mock server
	engine.Use(func(c *gin.Context) { build.ToContext(c, b) })
	engine.Use(func(c *gin.Context) { repo.ToContext(c, r) })
	engine.Use(func(c *gin.Context) { step.ToContext(c, s) })
	engine.Use(func(c *gin.Context) { user.ToContext(c, u) })
	engine.Use(Payload())
	engine.Use(Logger(logger, time.RFC3339, true))
	engine.POST("/foobar", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	gotLevel := hook.LastEntry().Level
	gotMessage := hook.LastEntry().Message

	if resp.Code != http.StatusOK {
		t.Errorf("Logger returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(gotLevel, wantLevel) {
		t.Errorf("Logger Level is %v, want %v", gotLevel, wantLevel)
	}

	if !reflect.DeepEqual(gotMessage, wantMessage) {
		t.Errorf("Logger Message is %v, want %v", gotMessage, wantMessage)
	}
}

func TestMiddleware_Logger_Error(t *testing.T) {
	// setup types
	wantLevel := logrus.ErrorLevel
	wantMessage := "Error #01: test error\n"
	logger, hook := test.NewNullLogger()
	defer hook.Reset()

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foobar", nil)

	// setup mock server
	engine.Use(Logger(logger, time.RFC3339, true))
	engine.GET("/foobar", func(c *gin.Context) {
		c.Error(fmt.Errorf("test error"))
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	gotLevel := hook.LastEntry().Level
	gotMessage := hook.LastEntry().Message

	if resp.Code != http.StatusOK {
		t.Errorf("Logger returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(gotLevel, wantLevel) {
		t.Errorf("Logger Level is %v, want %v", gotLevel, wantLevel)
	}

	if !reflect.DeepEqual(gotMessage, wantMessage) {
		t.Errorf("Logger Message is %v, want %v", gotMessage, wantMessage)
	}
}