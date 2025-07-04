package security

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/constant"
	"github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/http_agent"
	"github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/logger"
	"github.com/pkg/errors"
)

type NacosAuthClient struct {
	username           string
	password           string
	accessToken        *atomic.Value
	tokenTtl           int64
	lastRefreshTime    int64
	tokenRefreshWindow int64
	agent              http_agent.IHttpAgent
	clientCfg          constant.ClientConfig
	serverCfgs         []constant.ServerConfig
}

func NewNacosAuthClient(clientCfg constant.ClientConfig, serverCfgs []constant.ServerConfig, agent http_agent.IHttpAgent) *NacosAuthClient {
	client := &NacosAuthClient{
		username:    clientCfg.Username,
		password:    clientCfg.Password,
		serverCfgs:  serverCfgs,
		clientCfg:   clientCfg,
		agent:       agent,
		accessToken: &atomic.Value{},
	}

	return client
}

func (ac *NacosAuthClient) GetAccessToken() string {
	v := ac.accessToken.Load()
	if v == nil {
		return ""
	}
	return v.(string)
}

func (ac *NacosAuthClient) GetSecurityInfo(resource RequestResource) map[string]string {
	var securityInfo = make(map[string]string, 4)
	v := ac.accessToken.Load()
	if v != nil {
		securityInfo[constant.KEY_ACCESS_TOKEN] = v.(string)
	}
	return securityInfo
}

func (ac *NacosAuthClient) AutoRefresh(ctx context.Context) {

	// If the username is not set, the automatic refresh Token is not enabled

	if ac.username == "" {
		return
	}

	go func() {
		var timer *time.Timer
		if lastLoginSuccess := ac.lastRefreshTime > 0 && ac.tokenTtl > 0 && ac.tokenRefreshWindow > 0; lastLoginSuccess {
			timer = time.NewTimer(time.Second * time.Duration(ac.tokenTtl-ac.tokenRefreshWindow))
		} else {
			timer = time.NewTimer(time.Second * time.Duration(5))
		}
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				_, err := ac.Login()
				if err != nil {
					logger.Errorf("login has error %+v", err)
					timer.Reset(time.Second * time.Duration(5))
				} else {
					logger.Infof("login success, tokenTtl: %+v seconds, tokenRefreshWindow: %+v seconds", ac.tokenTtl, ac.tokenRefreshWindow)
					timer.Reset(time.Second * time.Duration(ac.tokenTtl-ac.tokenRefreshWindow))
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (ac *NacosAuthClient) Login() (bool, error) {
	var throwable error = nil
	for i := 0; i < len(ac.serverCfgs); i++ {
		result, err := ac.login(ac.serverCfgs[i])
		throwable = err
		if result {
			return true, nil
		}
	}
	return false, throwable
}

func (ac *NacosAuthClient) UpdateServerList(serverList []constant.ServerConfig) {
	ac.serverCfgs = serverList
}

func (ac *NacosAuthClient) GetServerList() []constant.ServerConfig {
	return ac.serverCfgs
}

func (ac *NacosAuthClient) login(server constant.ServerConfig) (bool, error) {
	if ac.lastRefreshTime > 0 && ac.tokenTtl > 0 {
		// We refresh 2 windows before expiration to ensure continuous availability
		tokenRefreshTime := ac.lastRefreshTime + ac.tokenTtl - 2*ac.tokenRefreshWindow
		if time.Now().Unix() < tokenRefreshTime {
			return true, nil
		}
	}
	if ac.username == "" {
		ac.lastRefreshTime = time.Now().Unix()
		return true, nil
	}

	contextPath := server.ContextPath

	if !strings.HasPrefix(contextPath, "/") {
		contextPath = "/" + contextPath
	}

	if strings.HasSuffix(contextPath, "/") {
		contextPath = contextPath[0 : len(contextPath)-1]
	}

	if server.Scheme == "" {
		server.Scheme = "http"
	}

	reqUrl := server.Scheme + "://" + server.IpAddr + ":" + strconv.FormatInt(int64(server.Port), 10) + contextPath + "/v1/auth/users/login"

	header := http.Header{
		"content-type": []string{"application/x-www-form-urlencoded"},
	}
	resp, err := ac.agent.Post(reqUrl, header, ac.clientCfg.TimeoutMs, map[string]string{
		"username": ac.username,
		"password": ac.password,
	})

	if err != nil {
		return false, err
	}

	var bytes []byte
	bytes, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return false, err
	}

	if resp.StatusCode != constant.RESPONSE_CODE_SUCCESS {
		errMsg := string(bytes)
		return false, errors.New(errMsg)
	}

	var result map[string]interface{}

	err = json.Unmarshal(bytes, &result)

	if err != nil {
		return false, err
	}

	if val, ok := result[constant.KEY_ACCESS_TOKEN]; ok {
		ac.accessToken.Store(val)
		ac.lastRefreshTime = time.Now().Unix()
		ac.tokenTtl = int64(result[constant.KEY_TOKEN_TTL].(float64))
		ac.tokenRefreshWindow = ac.tokenTtl / 10
	}

	return true, nil

}
