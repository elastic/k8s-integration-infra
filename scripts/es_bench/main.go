// $ go run main.go

package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	log.SetFlags(0)

	elasticQuery := `
{
  "query": {
    "bool": {
      "must": [
        {
          "range": {
            "@timestamp": {
              "gte": "now-24h",
              "lte": "now",
              "format": "strict_date_optional_time"
            }
          }
        },
        {
          "bool": {
            "must": [
              {
                "query_string": {
                  "query": "event.module:kubernetes AND metricset.name:container",
                  "analyze_wildcard": true
                }
              }
            ],
            "filter": [],
            "should": [],
            "must_not": []
          }
        }
      ],
      "filter": [],
      "should": [],
      "must_not": []
    }
  },
  "aggs": {
    "5d3692a1-2bfc-11e7-859b-f78b612cde28": {
      "terms": {
        "field": "kubernetes.pod.name",
        "order": {
          "5d3692a2-2bfc-11e7-859b-f78b612cde28-SORT": "desc"
        }
      },
      "aggs": {
        "5d3692a2-2bfc-11e7-859b-f78b612cde28-SORT": {
          "max": {
            "field": "kubernetes.container.cpu.usage.core.ns"
          }
        },
        "timeseries": {
          "auto_date_histogram": {
            "field": "@timestamp"
          },
          "aggs": {
            "5d3692a2-2bfc-11e7-859b-f78b612cde28": {
              "max": {
                "field": "kubernetes.container.cpu.usage.core.ns"
              }
            },
            "6c905240-2bfc-11e7-859b-f78b612cde28": {
              "derivative": {
                "buckets_path": "5d3692a2-2bfc-11e7-859b-f78b612cde28",
                "gap_policy": "skip",
                "unit": "1s"
              }
            },
            "9a51f710-359d-11e7-aa4a-8313a0c92a88": {
              "bucket_script": {
                "buckets_path": {
                  "value": "6c905240-2bfc-11e7-859b-f78b612cde28[normalized_value]"
                },
                "script": {
                  "source": "params.value > 0.0 ? params.value : 0.0",
                  "lang": "painless"
                },
                "gap_policy": "skip"
              }
            }
          }
        }
      }
    }
  }
}
`

	// Initialize a client
	cfgTSDB := elasticsearch.Config{
		Addresses: []string{
			"https://35.195.16.23:9200",
		},
		Username: "elastic",
		Password: "1TVP5SAGs207s3490TY2Nymo",
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	cfgSimple := elasticsearch.Config{
		Addresses: []string{
			"https://104.199.81.103:9200",
		},
		Username: "elastic",
		Password: "fgx1L9A1EQ7F3sG224S21ErB",
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	tsdbIndex := ".ds-metricbeat-tsdb-8.3.0-2022.05.23-000001"
	simpleIndex := ".ds-metricbeat-8.3.0-2022.05.23-000001"

	execute(tsdbIndex, elasticQuery, cfgTSDB)
	log.Printf("\n\n")
	execute(simpleIndex, elasticQuery, cfgSimple)

}

func execute(indexName string, elasticQuery string, cfg elasticsearch.Config) {
	log.Printf(strings.Repeat("*", 37))
	log.Printf("Executing against new ES cluster")
	log.Printf(strings.Repeat("*", 37))
	var (
		r  map[string]interface{}
		ind interface{}
		//wg sync.WaitGroup
	)

	// Initialize a client
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	// 1. Get cluster info
	//
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))


	res, err = esapi.CatIndicesRequest{Format: "json"}.Do(context.Background(), es)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&ind); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	for _, index := range ind.([]interface{}) {
		id := index.(map[string]interface{})
		if id["index"] == indexName {
			log.Printf("index name: %v", id["index"])
			log.Printf("pri.store.size: %v", id["pri.store.size"])
			log.Printf("docs.count: %v", id["docs.count"])
			log.Printf(strings.Repeat("~", 37))
		}
	}
	// 3. Search for the indexed documents
	//
	// Build the request body.
	var buf bytes.Buffer

	// Check for JSON errors
	isValid := json.Valid([]byte(elasticQuery)) // returns bool

	// Default query is "{}" if JSON is invalid
	if isValid == false {
		log.Fatalf("ERROR: query string not valid: %v", elasticQuery)
		log.Fatalf("Using default match_all query")
		elasticQuery = "{}"
	} else {
		log.Println("valid JSON:", isValid)
	}

	// Build a new string from JSON query
	var b strings.Builder
	b.WriteString(elasticQuery)

	// Instantiate a *strings.Reader object from string
	read := strings.NewReader(b.String())

	if err := json.NewEncoder(&buf).Encode(read); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	medianTime := 0
	for i := 1; i <= 20; i++ {
		res, err = es.Search(
			es.Search.WithContext(context.Background()),
			es.Search.WithIndex(indexName),
			es.Search.WithBody(&buf),
			es.Search.WithTrackTotalHits(true),
			es.Search.WithPretty(),
		)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			var e map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				log.Fatalf("Error parsing the response body: %s", err)
			} else {
				// Print the response status and error information.
				log.Fatalf("[%s] %s: %s",
					res.Status(),
					e["error"].(map[string]interface{})["type"],
					e["error"].(map[string]interface{})["reason"],
				)
			}
		}

		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}
		// Print the response status, number of results, and request duration.
		//log.Printf(
		//	"[%s] %d hits; took: %dms",
		//	res.Status(),
		//	int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		//	int(r["took"].(float64)),
		//)
		medianTime = medianTime + int(r["took"].(float64))
	}
	log.Printf("median time is: %dms", medianTime/20)
}
