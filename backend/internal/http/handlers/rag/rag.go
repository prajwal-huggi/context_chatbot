package rag

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/prajwal-huggi/context_chatbot/internal/utils/response"
)

var validate= validator.New()

type QuestionRequest struct {
	Question string `json:"question"`
}

type AnswerResponse struct {
	Model  string `json:"model"`
	Answer string `json:"answer"`
}

func GetAnswer()http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		slog.Info("Sending the request to RAG")
		var req QuestionRequest

		// 1. Decode incoming request JSON
		err:= json.NewDecoder(r.Body).Decode(&req)
		
		//If any field is empty
		if errors.Is(err, io.EOF){
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return 
		}

		if err!= nil{
			slog.Error("invalid request body")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// 1a. Validate the user request fields
		if err= validate.Struct(req); err!= nil{
			slog.Error("validation failed", "error", err)
			validateErrs:= err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		// 2. Marshal again and send to another backend
		reqBytes, err := json.Marshal(req)//Marshal convert it into the json format
		if err != nil {
			slog.Error("failed to marshal request", "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// 3. Send the request to the RAG backend
		ragBackendURL := os.Getenv("RAG_BACKEND_URL")+"/query"
		slog.Info("The backend url of the RAG is: ", slog.String("ragUrl", ragBackendURL))
		resp, err := http.Post(ragBackendURL, "application/json", bytes.NewBuffer(reqBytes))
		if err != nil {
			// response.WriteJson(w, http.StatusBadGateway, "failed to contact backend")
			response.WriteJson(w, http.StatusBadGateway, response.GeneralError(err))
			return
		}
		defer resp.Body.Close()

		// 4. Read backend response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			// http.Error(w, "failed to read backend response", http.StatusInternalServerError)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))

			return
		}

		// 5. Forward backend response directly to client

		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(resp.StatusCode)
		// w.Write(body)
		for k, v := range resp.Header {
			w.Header()[k] = v // preserve backend headers if needed
		}
		w.WriteHeader(resp.StatusCode)

		if _, err := w.Write(body); err != nil {
			slog.Error("failed to write response", "error", err)
		}
	}
}

func ResetRAG()http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		slog.Info("Sending the request to the RAG system to reset the database")

		// 1. Build RAG backend URL
		ragBackendURL:= os.Getenv("RAG_BACKEND_URL")+"/reset"

		// 2. Make the request to the RAG backend
		resp, err:= http.Post(ragBackendURL, "application/json", nil)
		if err!= nil{
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}
		defer resp.Body.Close()

		// 3. Read backend response
		body, err:= io.ReadAll(resp.Body)
		if err!= nil{
			slog.Error("failed to read backend response", "error", err)
            response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
            return
		}

		// 4. Forward response header and status
		for k, v:= range resp.Header{
			w.Header()[k]= v
		}
		w.WriteHeader(resp.StatusCode)

		if _, err:= w.Write(body); err!= nil{
			slog.Error("failed to write resopnse", "error", err)
		}

	}
}

// func UploadDocument()http.HandlerFunc{

// }