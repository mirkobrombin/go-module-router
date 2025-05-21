package swagger

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mirkobrombin/go-module-router/v1/registry"
)

type Document struct {
	OpenAPI string              `json:"openapi"`
	Info    Info                `json:"info"`
	Paths   map[string]PathItem `json:"paths"`
}

type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

type PathItem map[string]Operation

type Operation struct {
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses"`
}

type Parameter struct {
	Name   string `json:"name"`
	In     string `json:"in"`
	Schema Schema `json:"schema"`
}

type Schema struct {
	Type    string   `json:"type,omitempty"`
	Minimum *float64 `json:"minimum,omitempty"`
}

type Response struct {
	Description string `json:"description"`
}

func Build(title, version string) ([]byte, error) {
	doc := Document{
		OpenAPI: "3.0.3",
		Info: Info{
			Title:   title,
			Version: version,
		},
		Paths: map[string]PathItem{},
	}

	for _, rp := range registry.Global().RouteProviders {
		for _, r := range rp() {

			m := r.Meta
			ok := m != nil
			if !ok {
				continue
			}

			pathItem, ok := doc.Paths[r.Path]
			if !ok || pathItem == nil {
				pathItem = make(PathItem)
			}

			op := Operation{
				Summary:     m["summary"].(string),
				Description: m["description"].(string),
				Responses:   map[string]Response{},
			}

			if raw, ok := m["parameters"].([]any); ok {
				for _, entry := range raw {
					pm := entry.(map[string]any)
					sch := pm["schema"].(map[string]any)
					var minimum *float64
					if min, ok := sch["minimum"].(float64); ok {
						minimum = &min
					}
					op.Parameters = append(op.Parameters, Parameter{
						Name: pm["name"].(string),
						In:   pm["in"].(string),
						Schema: Schema{
							Type:    sch["type"].(string),
							Minimum: minimum,
						},
					})
				}
			}

			if raw, ok := m["responses"].(map[string]any); ok {
				for code, desc := range raw {
					op.Responses[code] = Response{Description: fmt.Sprint(desc)}
				}
			} else {
				op.Responses["200"] = Response{Description: "OK"}
			}

			pathItem[strings.ToLower(r.Method)] = op
			doc.Paths[r.Path] = pathItem
		}
	}

	return json.MarshalIndent(doc, "", "  ")
}
