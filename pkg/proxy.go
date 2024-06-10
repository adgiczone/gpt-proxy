package proxy

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

type LittleProxy struct {
	port     int
	timeout  int
	certfile string
	keyfile  string
}

type content struct {
	Url     string `json:"url"`
	Content string `json:"content"`
}

func GetPubKeyFromFile(filePath string) (*rsa.PublicKey, error) {
	pubByte, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		log.Error("Invalid file path")
		return nil, err
	}
	pData, _ := pem.Decode(pubByte)
	pubI, err := x509.ParsePKIXPublicKey(pData.Bytes)
	if err != nil {
		log.Error("Parse public key fail")
		return nil, err
	}
	pub := pubI.(*rsa.PublicKey)
	return pub, nil
}

func getJWTStandardClaims(pubKey *rsa.PublicKey, jwtToken string) (string, error) {
	claims := &jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (i interface{}, e error) {
		return pubKey, nil
	})
	return claims.Audience, err
}

func (s *LittleProxy) sendToTarget(url string, token string, msg string) (string, error) {
	log.Info("")
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(msg))
	if err != nil {
		log.Error("Error creating request:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	if len(token) != 0 {
		req.Header.Set("Authorization", token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response body:", err)
		return "", err
	}

	return string(body), nil
}

func writeResponse(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Error("write response error:", err)
	}
}

func (s *LittleProxy) Proxy(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	log.Infof("get method:%s", method)
	if method != http.MethodPost {
		writeResponse(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		writeResponse(w, "Authorization fail", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("read request body error", err)
		writeResponse(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	c := content{}
	err = json.Unmarshal(body, &c)
	if err != nil {
		log.Error("Unmarshal json fail", err)
		writeResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := s.sendToTarget(c.url, authToken, c.content)
	if err != nil {
		log.Error("Send to Target fail", err)
		writeResponse(w, "Send to target url fail", http.StatusInternalServerError)
		return
	}

	writeResponse(w, resp, http.StatusOK)
}

func New(port, timeout int, certFile, keyFile string) (*LittleProxy, error) {
	return &LittleProxy{
		port:     port,
		timeout:  timeout,
		certfile: certFile,
		keyfile:  keyFile,
	}, nil
}

func (s *LittleProxy) Serv() error {

	address := fmt.Sprintf(":%d", s.port)
	log.Info("Start listen at ", address)

	mux := http.NewServeMux()
	mux.HandleFunc("/proxy", s.Proxy)

	srv := http.Server{
		Addr:              address,
		Handler:           mux,
		ReadHeaderTimeout: time.Duration(s.timeout) * time.Second,
	}

	if len(s.certfile) > 0 || len(s.keyfile) > 0 {
		err := srv.ListenAndServeTLS(s.certfile, s.keyfile)
		if err != nil {
			return err
		}
	} else {
		err := srv.ListenAndServe()
		if err != nil {
			return err
		}
	}
	return nil
}
