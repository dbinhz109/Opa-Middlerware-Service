package service

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"encoding/json"
	"go-app/src/logger"
	"io"
	"net/http"
	"os/exec"
	"runtime/debug"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const ShellToUse = "/bin/bash"

// ShellExecute runs a shell command and return stdout and stderr
func ShellExecute(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// PanicOnError shortcut
func PanicOnError(err error, args ...zapcore.Field) {
	if err != nil {
		// stack := strings.Split(string(debug.Stack()), "\n")
		// stack = stack[10:20]
		// args = append(args, zap.Strings("trace", stack))
		logger.Warn("PanicError", args...)
		panic(err)
	}
}

// PanicRecover default panic recover
func PanicRecover() {
	err := recover()
	if err != nil {
		stack := strings.Split(string(debug.Stack()), "\n")
		stack = stack[4:]
		logger.Warn("PanicRecover", zap.Any("value", err), zap.Strings("trace", stack))
	}
}

// PanicRecoverTx panic recover with transaction rollback
func PanicRecoverTx(tx *sqlx.Tx) {
	err := recover()
	if err != nil {
		stack := strings.Split(string(debug.Stack()), "\n")
		stack = stack[4:]
		logger.Warn("PanicRecover", zap.Any("value", err), zap.Strings("trace", stack))
		if tx != nil {
			tx.Rollback()
		}
	}
}

type JsonMap = map[string]any

// JsonMapToStruct
func JsonMapToStruct[T any](m JsonMap) (*T, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	var t T
	err = json.Unmarshal(bytes, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Clone a JsonMap
func DeepCopyMap(m JsonMap) (JsonMap, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}
	var copy JsonMap
	err = dec.Decode(&copy)
	if err != nil {
		return nil, err
	}
	return copy, nil
}

func GenerateRandomDecimal(length int) string {
	var table = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	b := make([]byte, length)
	nBytes, err := io.ReadAtLeast(rand.Reader, b, length)
	if nBytes < length {
		logger.Info("GenerateRandomFailed", zap.Int("nbytes", nBytes), zap.Error(err))
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func GenerateRandomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	nBytes, err := io.ReadAtLeast(rand.Reader, b, length)
	if nBytes < length {
		logger.Info("GenerateRandomFailed", zap.Int("nbytes", nBytes), zap.Error(err))
		return nil, err
	}
	return b, nil
}

func HttpPost(url string, contentType string, body []byte) []byte {
	data, err := http.Post(url, contentType, bytes.NewBuffer(body))
	if err != nil {
		logger.Warn("http post failed", zap.Error(err))
		return nil
	}
	defer data.Body.Close()
	resp, err := io.ReadAll(data.Body)
	if err != nil {
		logger.Warn("http post response read failed", zap.Error(err), zap.ByteString("request", body))
		return nil
	}
	return resp
}

func HttpGet(url string) []byte {
	data, err := http.Get(url)
	if err != nil {
		logger.Warn("http get failed", zap.Error(err))
		return nil
	}
	defer data.Body.Close()
	resp, err := io.ReadAll(data.Body)
	if err != nil {
		logger.Warn("http get response read failed", zap.Error(err))
		return nil
	}
	return resp
}

// TimeStartOf returns start of unit m,h,d,w,M,y
func TimeStartOf(units string, t time.Time) time.Time {
	switch units {
	case "m":
		return t.Truncate(time.Minute)
	case "h":
		y, m, d := t.Date()
		return time.Date(y, m, d, t.Hour(), 0, 0, 0, t.Location())
	case "d":
		y, m, d := t.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	case "w":
		weekday := int(t.Weekday())
		return t.AddDate(0, 0, -weekday)
	case "M":
		y, m, _ := t.Date()
		return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
	case "y":
		y, _, _ := t.Date()
		return time.Date(y, time.January, 1, 0, 0, 0, 0, t.Location())
	default:
		return t
	}
}

// SliceMap maps a slice to another
func SliceMap[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}
