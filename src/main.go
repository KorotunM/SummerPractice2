package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type request struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}
type response struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func main() {
	_ = godotenv.Load()

	login := os.Getenv("SMSC_LOGIN")
	psw := os.Getenv("SMSC_PSW")

	http.HandleFunc("/api/send-sms", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(response{Success: false, Error: "Метод не поддерживается"})
			return
		}

		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response{Success: false, Error: "Ошибка в JSON"})
			return
		}

		apiURL := "https://smsc.ru/sys/send.php"
		params := url.Values{
			"login":   {login},
			"psw":     {psw},
			"phones":  {req.Phone},
			"mes":     {fmt.Sprintf("Код: %s", req.Code)},
			"charset": {"utf-8"},
			"fmt":     {"3"},
		}

		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(apiURL + "?" + params.Encode())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response{Success: false, Error: err.Error()})
			return
		}
		defer resp.Body.Close()

		var smsResp struct {
			ID      int    `json:"id"`
			Cnt     int    `json:"cnt"`
			Cost    string `json:"cost"`
			Balance string `json:"balance"`
			Error   string `json:"error"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&smsResp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response{Success: false, Error: "Ошибка разбора ответа"})
			return
		}

		if smsResp.ID == 0 {
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(response{Success: false, Error: smsResp.Error})
			return
		}

		log.Printf("SMSC отправил sms, id=%d, \n", smsResp.ID)
		json.NewEncoder(w).Encode(response{Success: true})
	})

	log.Println("Сервер доступен на  http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
