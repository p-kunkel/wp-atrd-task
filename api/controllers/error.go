package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	err        error
	msg        string
	statusCode int
}

func handleError(c *gin.Context, err error) {
	logMsg := err.Error()
	if v, ok := err.(*errorResponse); ok {

		if v.statusCode <= 0 {
			v.statusCode = 400
		}

		if v.msg == "" {
			v.setDefaultMsgByStatusCode()
			logMsg = fmt.Sprintf("%s, err: %s", v.msg, v.err.Error())
		}

		c.String(v.statusCode, v.msg)
	} else {
		c.Status(400)
	}

	c.Abort()

	log.Println(logMsg)
}

func (er *errorResponse) setDefaultMsgByStatusCode() {
	switch er.statusCode {
	case http.StatusMethodNotAllowed:
		er.msg = "Invalid input"
	case http.StatusBadRequest:
		er.msg = "Bad Request"
	case http.StatusNotFound:
		er.msg = "Secret not found"
	}
}

func newErrResp(status int, err error, msg ...string) *errorResponse {
	m := ""
	if len(msg) > 0 {
		m = msg[0]
	}
	return &errorResponse{
		err:        err,
		msg:        m,
		statusCode: status,
	}
}

func (er *errorResponse) Error() string {
	return er.msg
}
