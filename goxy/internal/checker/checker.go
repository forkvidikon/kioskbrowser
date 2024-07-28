package checker

import (
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

func HttpCodeTargetCheck(target string, zl *zap.Logger) (bool, int) {
	const fn = "http.checker.response"
	_, err := url.Parse(target)
	if err != nil {
		zl.Error("Backup URL parse error", zap.String("fn:", fn), zap.String("err:", err.Error()))
		return false, 0
	}

	resp, err := http.Get(target)
	if err != nil {
		zl.Error("Check target URL err:", zap.String("fn:", fn), zap.String("err:", err.Error()))
		return false, 0
	} else {

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			zl.Debug("Check target url completed successfully:", zap.String("fn:", fn))
			return true, resp.StatusCode
		} else {
			zl.Debug("Target url return non 200 code:", zap.String("fn:", fn))
			return false, resp.StatusCode
		}
	}

}
