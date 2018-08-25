package anilist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type AnilistAPIResposne struct {
	Data Data `json:"data"`
}

func getResponse(body []byte) (*AnilistAPIResposne, error) {
	var s = new(AnilistAPIResposne)
	err := json.Unmarshal(body, &s)
	if err != nil {
		panic(err)
	}
	return s, err
}

func runQuery(query string, variables map[string]string) interface{} {
	url := "https://graphql.anilist.co"
	values := map[string]interface{}{"query": query, "variables": variables}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		fmt.Println("1")
		panic(err.Error())
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("2")
		panic(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("3")
		panic(err.Error())
	}
	anilistData, err := getResponse([]byte(body))
	if err != nil {
		fmt.Println("4")
		panic(err.Error())
	}
	return anilistData.Data
}
