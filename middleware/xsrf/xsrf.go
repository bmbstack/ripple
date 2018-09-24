package xsrf

import (
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1 << letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n - 1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// Options maintains options to manage behavior of Generate.
type Options struct {
	// false: Deny all invalid request.
	// true: Check result in your own code by (*echo.Context).Get('XSRF_ISVAILD').
	Passed     bool
	// The global secret value used to generate Tokens.
	Secret     string
	// HTTP header used to set and get token.
	Header     string
	// Form value used to set and get token.
	Form       string
	// Cookie value used to set and get token.
	Cookie     string
	// Cookie path.
	CookiePath string
	// If true, send token via X-CSRFToken header.
	SetHeader  bool
	// If true, send token via _csrf cookie.
	SetCookie  bool
	// Set the Secure flag to true on the cookie.
	Secure     bool
}

func prepareOptions(options []Options) Options {
	var opt Options
	opt.Passed = false
	if len(options) > 0 {
		opt = options[0]
	}

	// Defaults.
	if len(opt.Secret) == 0 {
		opt.Secret = RandString(10)
	}
	if len(opt.Header) == 0 {
		opt.Header = "X-CSRFToken"
	}
	if len(opt.Form) == 0 {
		opt.Form = "_xsrf"
	}
	if len(opt.Cookie) == 0 {
		opt.Cookie = "_xsrf"
	}
	if len(opt.CookiePath) == 0 {
		opt.CookiePath = "/"
	}

	return opt
}

func XsrfProtection(options ...Options) echo.MiddlewareFunc {
	opt := prepareOptions(options)
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			var cookie *http.Cookie
			var token string
			var err error
			// bypass some http method
			if ctx.Request().Method != "POST" {
				token = GenerateToken(opt.Secret, "", "POST")
				ctx.Set("XSRF_TOKEN", token)
				ctx.Set("XSRF_ISVAILD", true)
				err = h(ctx)
				if err != nil {
					return err
				}
			} else {
				cookie, err = ctx.Request().Cookie("_xsrf")
				if err != nil {
					// _xsrf cookie not found
					if opt.Passed == false {
						return echo.NewHTTPError(http.StatusForbidden)
					} else {
						ctx.Set("XSRF_ISVAILD", false)
						cookie = &http.Cookie{}
					}
				} else {
					// xsrf cookie found
					value, _ := url.QueryUnescape(cookie.Value)
					if ctx.Form(opt.Form) != value {
						if opt.Passed == false {
							return echo.NewHTTPError(http.StatusForbidden)
						} else {
							ctx.Set("XSRF_ISVAILD", false)
							cookie = &http.Cookie{}
						}
					}
					if ValidToken(value, opt.Secret, "", "POST") == true {
						ctx.Set("XSRF_ISVAILD", true)
					} else {
						if opt.Passed == false {
							return echo.NewHTTPError(http.StatusForbidden)
						} else {
							ctx.Set("XSRF_ISVAILD", false)
						}
					}
				}
				token = GenerateToken(opt.Secret, "", "POST")
				err = h(ctx)
				if err != nil {
					return err
				}
			}
			if opt.SetCookie {
				cookie.Name = opt.Cookie
				cookie.Value = token
				cookie.Path = opt.CookiePath
				ctx.Response().Header().Add(opt.Cookie, cookie.String())
			}
			if opt.SetHeader {
				ctx.Response().Header().Add(opt.Header, token)
			}
			return nil
		}
	}
}
