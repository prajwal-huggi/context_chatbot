package rag

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"time"

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

func UploadDocument() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Uploading document to RAG system")

		// 1. Parse multipart form
		err := r.ParseMultipartForm(10 << 20) // 10MB
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// 2. Extract uploaded file
		file, handler, err := r.FormFile("file")
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		defer file.Close()

		// 3. Build a new multipart form
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)

		// 4. Preserve original headers (including Content-Type)
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", handler.Filename))
		h.Set("Content-Type", handler.Header.Get("Content-Type"))

		part, err := writer.CreatePart(h)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// 5. Copy file contents into new multipart
		if _, err := io.Copy(part, file); err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// 6. Close writer to finalize body
		writer.Close()

		// 7. Send request to RAG backend
		ragBackendURL := os.Getenv("RAG_BACKEND_URL") + "/add_pdf"
		req, err := http.NewRequest("POST", ragBackendURL, &buf)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			response.WriteJson(w, http.StatusBadGateway, response.GeneralError(err))
			return
		}
		defer resp.Body.Close()

		// 8. Relay backend response
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}
