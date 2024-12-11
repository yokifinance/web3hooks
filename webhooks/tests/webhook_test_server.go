package tests

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"yoki.finance/common/rcommon"
)

const internalServerPort = 8085

type WebhookTestServer struct {
	server                    *http.Server
	WebhookSuccessChan        chan WebhookCompleteRes
	WebhookHeaderChan         chan http.Header // request headers
	WebhookUrlOpenedChan      chan []byte      // request body is sent
	webhookFailsBeforeSuccess int
	currentFails              int
	returnBody                []byte
}

type WebhookCompleteRes struct {
	RequestBody  []byte
	ResponseBody []byte
}

func CreateWebhookSuccessTestServer() *WebhookTestServer {
	return CreateWebhookTestServer(0)
}

func CreateWebhookTestServer(webhookFailsBeforeSuccess int) *WebhookTestServer {
	mux := http.NewServeMux()
	s := &WebhookTestServer{
		server: &http.Server{
			Addr:    ":" + strconv.Itoa(internalServerPort),
			Handler: mux,
		},
		WebhookSuccessChan:        make(chan WebhookCompleteRes, 200),
		WebhookUrlOpenedChan:      make(chan []byte, 200),
		WebhookHeaderChan:         make(chan http.Header, 200),
		webhookFailsBeforeSuccess: webhookFailsBeforeSuccess,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.handlePost(w, r)
	})

	go func() {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("createWebhookTestServer: %s", err.Error())
		}
	}()

	return s
}

// Sets body that will be returned on webhook call. By default, empty body is returned.
func (s *WebhookTestServer) SetReturnBody(body []byte) {
	s.returnBody = body
}

// Stops the server and returns once it's stopped
func (s *WebhookTestServer) Stop() {
	rcommon.Println("Server STOP called")
	if err := s.server.Shutdown(context.Background()); err != nil {
		log.Fatalf("WebhookTestServer.Stop: %s", err.Error())
	}

	rcommon.Println("Server stopped")
}

// Stops the server and returns
func (s *WebhookTestServer) ListenUrl() string {
	return fmt.Sprintf("http://localhost:%d/", internalServerPort)
}

func (s *WebhookTestServer) signalWebhookSuccess(req, resp []byte) {
	// sleep, so that test checker will be activated only after worker has already finished webhook processing
	time.Sleep(time.Millisecond * 200)
	s.WebhookSuccessChan <- WebhookCompleteRes{req, resp}
}

func (s *WebhookTestServer) handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		rcommon.Println("HTTP Server: called")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		s.WebhookHeaderChan <- r.Header
		s.WebhookUrlOpenedChan <- body

		if s.currentFails < s.webhookFailsBeforeSuccess {
			s.currentFails++
			time.Sleep(time.Millisecond * 200)
			http.Error(w, "webhool call fail limit not exceeded", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(s.returnBody)

		rcommon.Println("HTTP Server: signal webhook success")
		go s.signalWebhookSuccess(body, s.returnBody)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
