package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var es, _ = elasticsearch.NewDefaultClient()

func main() {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("0) Exit")
		fmt.Println("1) Load spacecraft")
		fmt.Println("2) Get spacecraft")
		fmt.Println("3) Search spacecraft by key and value")
		fmt.Println("4) Search spacecraft by key and prefix")
		option := ReadText(reader, "Enter option")
		if option == "0" {
			Exit()
		} else if option == "1" {
			LoadData()
		} else if option == "2" {
			Get(reader)
		} else if option == "3" {
			Search(reader, "match")
		} else if option == "4" {
			Search(reader, "prefix")
		} else {
			fmt.Println("Invalid option")
		}
	}
}

func LoadData() {
	var spacecrafts []map[string]interface{}
	pageNumber := 0
	for {
		response, _ :=
			http.Get("http://stapi.co/api/v1/rest/spacecraft/search?pageSize=100&pageNumber=
		" + strconv.Itoa(pageNumber))
		body, _ := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		page := result["page"].(map[string]interface{})
		totalPages := int(page["totalPages"].(float64))
		crafts := result["spacecrafts"].([]interface{})
		for _, craftInterface := range crafts {
			craft := craftInterface.(map[string]interface{})
			spacecrafts = append(spacecrafts, craft)
		}

		pageNumber++
		if pageNumber >= totalPages {
			break
		}
	}
	for _, data := range spacecrafts {
		uid, _ := data["uid"].(string)
		jsonString, _ := json.Marshal(data)
		request := esapi.IndexRequest{Index: "stsc", DocumentID: uid, Body:
		strings.NewReader(string(jsonString))}
		request.Do(context.Background(), es)
	}
	print(len(spacecrafts), " spacecraft read\n")
}

func Print(spacecraft map[string]interface{}) {
	name := spacecraft["name"]
	status := ""
	if spacecraft["status"] != nil {
		status = "- " + spacecraft["status"].(string)
	}
	registry := ""
	if spacecraft["registry"] != nil {
		registry = "- " + spacecraft["registry"].(string)
	}
	class := ""
	if spacecraft["spacecraftClass"] != nil {
		class = "- " +
			spacecraft["spacecraftClass"].(map[string]interface{})["name"].(string)
	}

	fmt.Println(name, registry, class, status)
}

func Get(reader *bufio.Scanner) {
	id := ReadText(reader, "Enter spacecraft ID")
	request := esapi.GetRequest{Index: "stsc", DocumentID: id}
	response, _ := request.Do(context.Background(), es)
	var results map[string]interface{}
	json.NewDecoder(response.Body).Decode(&results)
	Print(results["_source"].(map[string]interface{}))
}

func Search(reader *bufio.Scanner, querytype string) {
	key := ReadText(reader, "Enter key")
	value := ReadText(reader, "Enter value")
	var buffer bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			querytype: map[string]interface{}{
				key: value,
			},
		},
	}
	json.NewEncoder(&buffer).Encode(query)
	response, _ := es.Search(es.Search.WithIndex("stsc"),
		es.Search.WithBody(&buffer))
	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)
	for _, hit := range
		result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		craft :=
			hit.(map[string]interface{})["_source"].(map[string]interface{})
		Print(craft)
	}
}