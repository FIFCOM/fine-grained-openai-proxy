package ctrl

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fine-grained-openai-proxy/conf"
	"fine-grained-openai-proxy/svc"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Proxy proxy request to OpenAI
//
// Note: this function is http.HandlerFunc, not gin.HandlerFunc,
// so use gin.WrapF(Proxy) to convert it to gin.HandlerFunc
func Proxy(w http.ResponseWriter, r *http.Request) {
	fg_token := r.Header.Get("Authorization")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// r.Body is io.ReadCloser, it will be closed or wiped after read, so we need to re-write it
	r.Body = io.NopCloser(strings.NewReader(string(body)))

	// json parse body to var param
	var params map[string]interface{}
	
	// For requests without parameters or no "model" in parameters, 
	// the request will forwarded directly to OpenAI
	params_ok := true
	if err := json.Unmarshal(body, &params); err != nil {
		params_ok = false
	}

	models, models_ok := params["model"]

	var OpenAIApiKey string
	if !params_ok || !models_ok {
		OpenAIApiKey = fg_token
	} else {
		OpenAIApiKey = ConvertToken(models, fg_token) // convert fine-grained token to real openai token
	}

	client := http.DefaultClient
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	client.Transport = tr

	req, err := http.NewRequest(r.Method, conf.OpenAIAPIBaseUrl + r.URL.Path, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header = r.Header

	req.Header.Set("Authorization", OpenAIApiKey)
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("Connection", "keep-alive")

	rsp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}(rsp.Body)

	for name, values := range rsp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	rsp.Header.Del("Transfer-Encoding")

	head := map[string]string{
		"Cache-Control":                    "no-store",
		"access-control-allow-origin":      "*",
		"access-control-allow-credentials": "true",
		"Connection":                       "keep-alive",
	}
	content_type := rsp.Header.Get("Content-Type")
	if content_type == "text/event-stream" {
		head["Transfer-Encoding"] = "chunked"
	}

	for k, v := range head {
		if _, ok := rsp.Header[k]; !ok {
			w.Header().Set(k, v)
		}
	}

	rsp.Header.Del("content-security-policy")
	rsp.Header.Del("content-security-policy-report-only")
	rsp.Header.Del("clear-site-data")
	w.Header().Set("Accept-Encoding", "gzip")

	w.WriteHeader(rsp.StatusCode)

	if content_type == "text/event-stream" {
		// stream output to client
		scanner := bufio.NewScanner(rsp.Body)
		for scanner.Scan() {
			_, _ = w.Write([]byte(strconv.Itoa(len(scanner.Text())) + "\r\n"))
			_, _ = w.Write(scanner.Bytes())
			_, _ = w.Write([]byte("\r\n"))
			w.(http.Flusher).Flush()
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("failed to read response: %v", err)
		}
	} else {
		// simple copy response body to client
		_, err = io.Copy(w, rsp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// ConvertToken convert fine-grained token to real openai token
func ConvertToken(models interface{}, fg_token string) string {
    var token string
    fgSvc := svc.FineGrainedKeySvc{}
    var fgKey *svc.FineGrainedKey
    fgKeyOK := false
    if fg_token != "" {
        parts := strings.SplitN(fg_token, " ", 2)
        if len(parts) == 2 && parts[0] == "Bearer" {
            token = parts[1]
        } else {
            return fg_token // no token
        }
        var err error
        fgKey, err = fgSvc.ByHash(svc.Hash(token))
        timeNow := time.Now().Unix()
        if err != nil || fgKey.RemainCalls == 0 || fgKey.Expire <= timeNow {
            return fg_token // invalid token
        }
        fgKeyOK = true
    } else {
        return fg_token // No Authorization header
    }

    modelSvc := svc.ModelSvc{}
    model, err := modelSvc.ByName(models.(string))
    if err != nil {
        return fg_token // model not in the list of allowed models
    }
    modelOK := false
    var arr []int64

    // json_decode fgKey.List to arr
    _ = json.Unmarshal([]byte(fgKey.List), &arr)
    if fgKey.Type == "whitelist" {
        for _, mid := range arr {
            if mid == model.ID {
                modelOK = true
            }
        }
    } else if fgKey.Type == "blacklist" {
        modelOK = true
        for _, mid := range arr {
            if mid == model.ID {
                modelOK = false
            }
        }
    }

    if !fgKeyOK || !modelOK {
        return fg_token // fgKey invalid or model not allowed
    }
    fgKey.RemainCalls -= 1
    _ = fgSvc.Update(fgKey)

    apikeySvc := svc.ApiKeySvc{}
    apiKey, _ := apikeySvc.ByID(fgKey.ParentID)

	return "Bearer " + apiKey.Key
}