/*
   Copyright 2000-2017 JetBrains s.r.o.
   Copyright 2021 Evan Cleary

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package extractor

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

type Sdk2Extractor struct{}

// ExportSchema should be called to export the structure
// of the provider.
func (m *Sdk2Extractor) Export(p *schema.Provider, pi *ProviderInfo) *ResourceProviderSchema {
	result := new(ResourceProviderSchema)

	result.SchemaVersion = "2"
	result.SDKType = "terraform-sdk-2"
	result.Name = pi.Name
	result.Type = "provider"
	result.Version = pi.Revision
	result.Provider = schemaMapSdk2(p.Schema).Export(m)
	result.Resources = make(map[string]SchemaInfoWithTimeouts)
	result.DataSources = make(map[string]SchemaInfoWithTimeouts)

	for k, r := range p.ResourcesMap {
		result.Resources[k] = m.ExportResourceWithTimeouts(r)
	}
	for k, ds := range p.DataSourcesMap {
		result.DataSources[k] = m.ExportResourceWithTimeouts(ds)
	}

	return result
}

func (m *Sdk2Extractor) ExportResourceWithTimeouts(r *schema.Resource) SchemaInfoWithTimeouts {
	var timeouts []string
	t := r.Timeouts
	if t != nil {
		for _, key := range timeoutKeys() {
			var timeout *time.Duration
			switch key {
			case TimeoutCreate:
				timeout = t.Create
			case TimeoutUpdate:
				timeout = t.Update
			case TimeoutRead:
				timeout = t.Read
			case TimeoutDelete:
				timeout = t.Delete
			case TimeoutDefault:
				timeout = t.Default
			default:
				panic("Unsupported timeout key, update switch statement!")
			}
			if timeout != nil {
				timeouts = append(timeouts, key)
			}
		}
	}
	result := make(SchemaInfoWithTimeouts)
	for nk, nv := range m.ExportResource(r) {
		result[nk] = nv
	}
	if len(timeouts) > 0 {
		result["__timeouts__"] = timeouts
	}
	return result
}

func (m *Sdk2Extractor) ExportResource(r *schema.Resource) SchemaInfo {
	return schemaMapSdk2(r.Schema).Export(m)
}

// schemaMap is a wrapper that adds nice functions on top of schemas.
type schemaMapSdk2 map[string]*schema.Schema

// Export exports the format of this schema.
func (m schemaMapSdk2) Export(extractor *Sdk2Extractor) SchemaInfo {
	result := make(SchemaInfo)
	for k, v := range m {
		item := extractor.export(v)
		result[k] = item
	}
	return result
}

func (m *Sdk2Extractor) export(v *schema.Schema) SchemaDefinition {
	item := SchemaDefinition{}

	item.Type = shortenType(fmt.Sprintf("%s", v.Type))
	item.Optional = v.Optional
	item.Required = v.Required
	item.Description = v.Description
	item.InputDefault = v.InputDefault
	item.Computed = v.Computed
	item.MaxItems = v.MaxItems
	item.MinItems = v.MinItems
	item.ComputedWhen = v.ComputedWhen
	item.ConflictsWith = v.ConflictsWith
	item.Deprecated = v.Deprecated

	if v.Type == schema.TypeList || v.Type == schema.TypeSet {
		if v.ConfigMode == schema.SchemaConfigModeBlock || v.ConfigMode == schema.SchemaConfigModeAuto {
			item.IsBlock = true
			item.ConfigImplicitMode = "Block"
			if v.Computed {
				item.ConfigImplicitMode = "ComputedBlock"
			}
		} else {
			item.ConfigImplicitMode = "Attr"
		}
	}

	if v.Elem != nil {
		item.Elem = m.exportValue(v.Elem, fmt.Sprintf("%T", v.Elem))
	}

	// TODO: Find better solution
	if defValue, err := v.DefaultValue(); err == nil && defValue != nil && !reflect.DeepEqual(defValue, v.Default) {
		item.Default = m.exportValue(defValue, fmt.Sprintf("%T", defValue))
	}
	return item
}

func (m *Sdk2Extractor) exportValue(value interface{}, t string) *SchemaElement {
	s2, ok := value.(*schema.Schema)
	if ok {
		return &SchemaElement{Type: "SchemaElements", ElementsType: shortenType(fmt.Sprintf("%s", s2.Type))}
	}
	r2, ok := value.(*schema.Resource)
	if ok {
		return &SchemaElement{Type: "SchemaInfo", Info: m.ExportResource(r2)}
	}
	vt, ok := value.(schema.ValueType)
	if ok {
		return &SchemaElement{Value: shortenType(fmt.Sprintf("%v", vt))}
	}
	// Unknown case
	return &SchemaElement{Type: t, Value: fmt.Sprintf("%v", value)}
}

func (m *Sdk2Extractor) Generate(provider *schema.Provider, pi *ProviderInfo, outputPath string) {
	outputFilePath := filepath.Join(outputPath, fmt.Sprintf("%s.json", pi.Name))

	if err := m.DoGenerate(provider, pi, outputFilePath); err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err.Error())
		os.Exit(255)
	}
}

func (m *Sdk2Extractor) DoGenerate(provider *schema.Provider, pi *ProviderInfo, outputFilePath string) error {
	providerJson, err := json.MarshalIndent(m.Export(provider, pi), "", "  ")

	if err != nil {
		return err
	}

	file, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(providerJson)
	if err != nil {
		return err
	}

	return file.Sync()
}

// func main() {
// 	var provider *schema.Provider
// 	// provider = __NAME__.Provider()

// 	Generate(provider, "__NAME__", "__OUT__")

// }
