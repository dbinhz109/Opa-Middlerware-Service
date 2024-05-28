package impl

import (
	"bytes"
	"encoding/json"
	"go-app/src/logger"
	"go-app/src/service/dto"
	"go-app/src/service/spec/apperror"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func LinkSafePost(endpoint string, param any) (r *dto.ApiResponse, e error) {
	defer func() {
		rec := recover()
		if rec != nil {
			err, ok := rec.(error)
			if ok {
				logger.Warn("PanicRecover", zap.Error(err))
				e = err
			} else {
				logger.Warn("PanicRecover2", zap.Any("value", err))
				e = apperror.ErrInternal
			}
		}
	}()
	payload, _ := json.Marshal(param)
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		logger.Warn("PostFailed", zap.Error(err), zap.String("endpoint", endpoint), zap.ByteString("param", payload))
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bodyResp := dto.ApiResponse{}
	err = json.Unmarshal(body, &bodyResp)
	if err != nil {
		logger.Warn("InvalidApiResponse", zap.ByteString("response", body))
		return nil, err
	}
	return &bodyResp, nil
}

func LinkSafeGet(endpoint string) (r *dto.ApiResponse, e error) {
	defer func() {
		rec := recover()
		if rec != nil {
			err, ok := rec.(error)
			if ok {
				logger.Warn("PanicRecover", zap.Error(err))
				e = err
			} else {
				logger.Warn("PanicRecover2", zap.Any("value", err))
				e = apperror.ErrInternal
			}
		}
	}()
	resp, err := http.Get(endpoint)
	if err != nil {
		logger.Warn("GetFailed", zap.Error(err), zap.String("endpoint", endpoint))
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bodyResp := dto.ApiResponse{}
	err = json.Unmarshal(body, &bodyResp)
	if err != nil {
		logger.Warn("InvalidApiResponse", zap.ByteString("response", body))
		return nil, err
	}
	return &bodyResp, nil
}
