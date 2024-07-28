package proxy

import (
	"errors"
	"goxy/internal/checker"
	"goxy/internal/config"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.uber.org/zap"
)

var (
	ErrNotEstablished = errors.New("The connection has not been established")
)

func Handle(zl *zap.Logger, cfg *config.Config, userAgent string) http.HandlerFunc {
	const fn = "goxy.proxy.handle"

	check, _ := checker.HttpCodeTargetCheck(cfg.TargetURL, zl)

	if !check {
		return func(w http.ResponseWriter, r *http.Request) {
			bUrl, err := url.Parse(cfg.BackupURL)
			//bUrl, err := url.Parse("https://vidikon.info")
			if err != nil {
				zl.Error("Backup URL parse error", zap.String("fn:", fn), zap.String("err:", err.Error()))
				http.Error(w, ErrNotEstablished.Error(), http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name: "backredirect",
				// Возможно стоит заменить таммстамп на зашифрованную строку для проверки на сервисе
				Value:  strconv.FormatInt(time.Now().Unix(), 10),
				Domain: bUrl.Host,
				MaxAge: 3600,
			})
			http.Redirect(w, &http.Request{}, cfg.BackupURL, 301)
		}
	} else {

		return func(w http.ResponseWriter, r *http.Request) {
			// Создаем URL для запроса с сохранением пути и параметров
			targetURL, err := url.Parse(cfg.TargetURL)
			if err != nil {
				zl.Error("URL parse error", zap.String("fn:", fn), zap.String("err:", err.Error()))
				http.Error(w, ErrNotEstablished.Error(), http.StatusInternalServerError)
				return
			}
			targetURL.Path = r.URL.Path
			targetURL.RawQuery = r.URL.RawQuery

			// Создаем новый запрос к целевому URL
			req, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
			if err != nil {
				zl.Error("Create request error", zap.String("fn:", fn), zap.String("err:", err.Error()))
				http.Error(w, ErrNotEstablished.Error(), http.StatusInternalServerError)
				return
			}

			// Копируем заголовки оригинального запроса
			req.Header = r.Header

			// Заменяем User-Agent и Host заголовки
			req.Header.Set("User-Agent", userAgent)
			req.Host = cfg.Host

			// Отправляем запрос к целевому серверу
			client := &http.Client{Timeout: cfg.RequestTimeout}
			resp, err := client.Do(req)
			if err != nil {
				zl.Error("Request was not executed", zap.String("fn:", fn), zap.String("err:", err.Error()))
				http.Error(w, ErrNotEstablished.Error(), http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			// Копируем статус и заголовки ответа от целевого сервера
			for key, values := range resp.Header {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.WriteHeader(resp.StatusCode)

			// Копируем тело ответа
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "write: broken pipe" {
					// Отлавливаем ошибку broken pipe
					zl.Info("Client disconnected before response was fully sent", zap.String("fn:", fn))
				} else {
					zl.Error("Error copying response body", zap.String("fn:", fn), zap.String("err:", err.Error()))
				}
				return

			}
		}
	}
}
