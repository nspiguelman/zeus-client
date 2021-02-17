package client

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var kahootPostURL = "http://localhost:8080/room"
var questionPostURL = "http://localhost:8080/room/:pin/question"
var answerPostURL = "http://localhost:8080/question/:id/answer"
var client = &http.Client{}

type KahootBody struct {
	Name string
}

type KahootReceived struct {
	Id int
	Pin string
	Name string
}

type QuestionBody struct {
	Question string
	Description string
}

type AnswerBody struct {
	Description string
	IsTrue bool
}

type QuestionReceived struct {
	Id          int
	Question    string
	Description string
	KahootID    int
}

type AnswerReceived struct {
	Id int
	Description string
	QuestionId int
	IsTrue bool
}

type AnswersReceived []AnswerReceived

func makePostRequest(urlStr string, jsonStr []byte) ([]byte, string, error) {
	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	status := resp.Status
	body, _ := ioutil.ReadAll(resp.Body)
	return body, status, nil
}

func processKahoots(filename string) ([]KahootReceived, error) {
	var message KahootReceived
	response := make([]KahootReceived, 0)
	csvFile, err := os.Open(filename)
	defer csvFile.Close()
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(csvFile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		bytes, err := json.Marshal(KahootBody{record[0]})
		if err != nil {
			return nil, err
		}

		byteResp, status, err := makePostRequest(kahootPostURL, bytes)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(byteResp, &message)
		if err != nil {
			return nil, err
		}
		if status == "201 Created" {
			response = append(response, message)
		}
	}
	return response, nil
}

func processQuestions(filename string, pin string) ([]QuestionReceived, error) {
	var questionReceived QuestionReceived
	response := make([]QuestionReceived, 0)
	csvFile, err := os.Open(filename)
	defer csvFile.Close()
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(csvFile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		bytes, err := json.Marshal(QuestionBody{record[0], record[1] })
		if err != nil {
			return nil, err
		}

		byteResp, status, err := makePostRequest(strings.Replace(questionPostURL, ":pin", pin, 1), bytes)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(byteResp, &questionReceived)
		if err != nil {
			return nil, err
		}

		if status == "201 Created" {
			response = append(response, questionReceived)
		}
	}
	return response, nil
}

func processAnswers(filename string, id string) (AnswersReceived, error) {
	var answersReceived AnswersReceived
	body := make([]AnswerBody, 0)
	csvFile, err := os.Open(filename)
	defer csvFile.Close()
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(csvFile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		isTrue, err := strconv.ParseBool(record[1])
		if err != nil {
			return nil, err
		}
		body = append(body, AnswerBody{record[0], isTrue })
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	byteResp, status, err := makePostRequest(strings.Replace(answerPostURL, ":id", id, 1), bytes)
	if err != nil {
		return nil, err
	}
	if status != "201 Created" {
		return nil, errors.New(status)
	}
	err = json.Unmarshal(byteResp, &answersReceived)
	if err != nil {
		return nil, err
	}
	return answersReceived, nil
}

func handlerError (wrongLength bool, messageToInvalidLength string, err error) {
	if err != nil {
		panic(err)
	}
	if wrongLength {
		panic(messageToInvalidLength)
	}
}

func ProcessCSV() int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	kahoots, err := processKahoots("./CSVs/kahoots.csv")
	handlerError(len(kahoots) == 0, "The game can't be played. At least one kahoot must have been created", err)
	pinToPlay := kahoots[0].Pin
	questions, err := processQuestions("./CSVs/questions.csv", pinToPlay)
	handlerError(len(questions) == 0, "The game can't be played. At least one question must have been created", err)
	for _, value := range questions {
		answers, err := processAnswers("./CSVs/answers.csv", strconv.Itoa(value.Id))
		handlerError(len(answers) != 4, "The game can't be played. Each question must have four answers", err)
	}
	response, err := strconv.Atoi(pinToPlay)
	return response
}
