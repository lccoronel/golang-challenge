package main

import (
	"io"
	"net/http"
	"time"
)

type ResponseCEP struct {
	Api  string `json:"api"`
	Body string `json:"body"`
}

func (responseCEP *ResponseCEP) MarshalJSON() ([]byte, error) {
	return []byte(`{"api": "` + responseCEP.Api + `", "body": ` + responseCEP.Body + `}`), nil
}

// Thread 1
func main() {
	canal := make(chan ResponseCEP)
	var cep string = "06855-330"

	API1 := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
	API2 := "https://viacep.com.br/ws/" + cep + "/json/"

	// Thread 2
	go func() {
		time.Sleep(time.Second * 2)
		canal <- *httpGetGEP(API1, 1)
	}()

	// Thread 3
	go func() {
		time.Sleep(time.Second * 2)
		canal <- *httpGetGEP(API2, 2)
	}()

	response := <-canal
	responseJSON, _ := response.MarshalJSON()
	println(string(responseJSON))
}

func httpGetGEP(API string, APINumber int) *ResponseCEP {
	client := http.Client{Timeout: time.Second}

	response, err := client.Get(API)
	if err != nil {
		println("Erro ao buscar CEP", APINumber, err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		println("Erro ao ler Body", err)
	}

	return &ResponseCEP{Api: API, Body: string(body)}
}
