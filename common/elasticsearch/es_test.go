package esutil_test

import (
	esutil "crf-mold/common/elasticsearch"
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

var node = []string{"http://172.16.1.23:9200"}

func TestSearch(t *testing.T) {
	viper.Set("es.node", node)
	key := "世界"
	query := map[string]interface{}{
		"_source": []string{"content", "name"},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"title": key,
						},
					},
					{
						"match": map[string]interface{}{
							"content": key,
						},
					},
				},
			},
		},
		"highlight": map[string]interface{}{
			"pre_tags":  []string{"<em>"},
			"post_tags": []string{"</em>"},
			"fields": map[string]interface{}{
				"title":   map[string]interface{}{},
				"content": map[string]interface{}{},
			},
		},
	}

	r, err := esutil.Search(&query)
	if err != nil {
		panic(err)
	}

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		fmt.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["highlight"])
		fmt.Println(hit.(map[string]interface{})["_source"].((map[string]interface{})))
	}
}

func TestIndex(t *testing.T) {
	viper.Set("es.node", node)
	doc := map[string]interface{}{
		"name":       "3",
		"content":    "3外面的世界真的很精彩",
		"group":      "group3",
		"gmtCreated": "2022-03-20 20:58:00",
		"Created":    "3",
	}

	err := esutil.Save("25", &doc)
	if err != nil {
		panic(err)
	}

	doc = map[string]interface{}{
		"name":       "4外面的世界真的很精彩",
		"content":    "4",
		"group":      "group4",
		"gmtCreated": "2022-03-20 20:58:00",
		"Created":    "4",
	}

	err = esutil.Save("26", &doc)
	if err != nil {
		panic(err)
	}
}

func TestDelete(t *testing.T) {
	viper.Set("es.node", node)
	err := esutil.Delete("12")
	if err != nil {
		panic(err)
	}
}

func TestTrimHtml(t *testing.T) {
	h := "<html><srcript></srcript><style></style>><p>hello world</p> chengong</html>"
	fmt.Println(esutil.TrimHtml(h))
}
