package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setup(t *testing.T) (*gin.Engine, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return nil, w
}

func TestSuccessShape(t *testing.T) {
	_, w := setup(t)
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	Success(c, gin.H{"status": "ok"})

	if w.Code != http.StatusOK {
		t.Fatalf("应返回 200，实际 %d", w.Code)
	}
	var body struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("响应不是合法 JSON: %v", err)
	}
	if body.Code != 0 || body.Msg != "success" {
		t.Fatalf("成功响应应为 code=0/msg=success，实际 %+v", body)
	}
}

func TestFailShape(t *testing.T) {
	_, w := setup(t)
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	Fail(c, ErrItemNotFound)

	if w.Code != http.StatusNotFound {
		t.Fatalf("应返回 404，实际 %d", w.Code)
	}
	var body struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("响应不是合法 JSON: %v", err)
	}
	if body.Code != 12001 || body.Msg != "物品不存在" {
		t.Fatalf("应为 code=12001/msg=物品不存在，实际 %+v", body)
	}
}
