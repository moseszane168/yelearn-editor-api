package esutil

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/spf13/viper"
)

var INDEX = "crf-mold"

func NewClient() (*elasticsearch.Client, error) {
	index := viper.GetString("es.index")
	if index != "" {
		INDEX = index
	}

	node := viper.GetStringSlice("es.node")
	cfg := elasticsearch.Config{
		Addresses: node,
	}
	return elasticsearch.NewClient(cfg)
}

func Search(query *map[string]interface{}) (map[string]interface{}, error) {
	c, err := NewClient()
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		panic(err)
	}
	// Perform the search request.
	res, err := c.Search(
		c.Search.WithContext(context.Background()),
		c.Search.WithIndex(INDEX),
		c.Search.WithBody(&buf),
		c.Search.WithTrackTotalHits(true),
		c.Search.WithPretty(),
	)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// 结果处理
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
			return nil, err
		} else {
			// Print the response status and error information.
			msg := fmt.Sprintf("[%s] %s: %s", res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"])
			log.Fatal(msg)
			return nil, errors.New(msg)
		}
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r, nil

}

func Save(id string, doc *map[string]interface{}) error {
	c, err := NewClient()
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		return err
	}
	res, err := c.Index(INDEX, &buf, c.Index.WithDocumentID(id))
	if err != nil {
		return err
	}

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
			return err
		} else {
			// Print the response status and error information.
			msg := fmt.Sprintf("[%s] %s: %s", res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"])
			log.Fatal(msg)
			return errors.New(msg)
		}
	}

	return nil
}

func Delete(id string) error {
	c, err := NewClient()
	if err != nil {
		panic(err)
	}

	res, err := c.Delete(INDEX, id)
	if err != nil {
		return err
	}

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
			return err
		} else {
			// Print the response status and error information.
			msg := fmt.Sprintf("[%s] %s: %s", res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"])
			log.Fatal(msg)
			return errors.New(msg)
		}
	}

	return nil

}

func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile(`\<[\S\s]+?\>`)
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile(`\<style[\S\s]+?\</style\>`)
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile(`\<script[\S\s]+?\</script\>`)
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile(`\<[\S\s]+?\>`)
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile(`\s{2,}`)
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}
