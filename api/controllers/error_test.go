package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	rw := httptest.NewRecorder()
	c := gin.CreateTestContextOnly(rw, gin.New())
	errResp := "test internal error"

	handleError(c, newErrResp(500, errors.New(errResp), errResp))

	if !assert.Equal(t, http.StatusInternalServerError, rw.Result().StatusCode) {
		t.FailNow()
	}

	if !assert.Equal(t, errResp, rw.Body.String()) {
		t.FailNow()
	}
}

func TestErrRespSetDefaultMsgByStatusCode(t *testing.T) {
	err := errorResponse{
		statusCode: http.StatusBadRequest,
	}
	msg := "Bad Request"

	err.setDefaultMsgByStatusCode()
	if !assert.Equal(t, msg, err.msg) {
		t.FailNow()
	}
}

func TestNewErrResp(t *testing.T) {
	errMsg := "test"
	err := errors.New(errMsg)
	errResp := newErrResp(http.StatusBadRequest, err, errMsg)

	if !assert.EqualError(t, errResp, errMsg) {
		t.FailNow()
	}

	errResp = newErrResp(http.StatusBadRequest, err)

	if !assert.EqualError(t, errResp, "") ||
		!assert.Equal(t, http.StatusBadRequest, errResp.statusCode) ||
		!assert.Equal(t, err, errResp.err) {
		t.FailNow()
	}
}
