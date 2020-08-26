package solr

import (
	"context"
	"fmt"
	"net/http"
)

const (
	SchemaCommandAddField            SchemaCommand = "add-field"
	SchemaCommandDeleteField         SchemaCommand = "delete-field"
	SchemaCommandReplaceField        SchemaCommand = "replace-field"
	SchemaCommandAddDynamicField     SchemaCommand = "add-dynamic-field"
	SchemaCommandDeleteDynamicField  SchemaCommand = "delete-dynamic-field"
	SchemaCommandReplaceDynamicField SchemaCommand = "replace-dynamic-field"
	SchemaCommandAddFieldType        SchemaCommand = "add-field-type"
	SchemaCommandDeleteFieldType     SchemaCommand = "delete-field-type"
	SchemaCommandReplaceFieldType    SchemaCommand = "replace-field-type"
	SchemaCommandAddCopyField        SchemaCommand = "add-copy-field"
	SchemaCommandDeleteCopyField     SchemaCommand = "delete-copy-field"
)

type SchemaAPI struct {
	conn *Connection
	Path string
}

type Analyzer struct {
	Tokenizer map[string]interface{}   `json:"tokenizer"`
	Filters   []map[string]interface{} `json:"filters"`
}

// https://lucene.apache.org/solr/guide/8_5/field-type-definitions-and-properties.html#general-properties
type FieldType struct {
	Name                      string    `json:"name"`
	CLass                     string    `json:"class"`
	PositionIncrementGap      string    `json:"positionIncrementGap,omitempty"`
	AutoGeneratePhraseQueries string    `json:"autoGeneratePhraseQueries,omitempty"`
	SynonymQueryStyle         string    `json:"synonymQueryStyle,omitempty"`
	EnableGraphQueries        bool      `json:"enableGraphQueries,omitempty"`
	DocValuesFormat           string    `json:"docValuesFormat,omitempty"`
	PostingsFormat            string    `json:"postingsFormat,omitempty"`
	Analyzer                  *Analyzer `json:"analyzer,omitempty"`
	IndexAnalyzer             *Analyzer `json:"indexAnalyzer,omitempty"`
	QueryAnalyzer             *Analyzer `json:"queryAnalyzer,omitempty"`
	FieldDefaultProperties
}

// schema  v1.6
//https://lucene.apache.org/solr/guide/8_5/field-type-definitions-and-properties.html#field-default-properties
type FieldDefaultProperties struct {
	Indexed                  *bool `json:"indexed,omitempty"`
	Stored                   *bool `json:"stored,omitempty"`
	DocValues                *bool `json:"docValues,omitempty"`
	SortMissingFirst         *bool `json:"sortMissingFirst,omitempty"`
	SortMissingLast          *bool `json:"sortMissingLast,omitempty"`
	MultiValued              *bool `json:"multiValued,omitempty"`
	Uninvertible             *bool `json:"uninvertible,omitempty"`
	OmitNorms                *bool `json:"omitNorms,omitempty"`
	OmitTermFreqAndPositions *bool `json:"omitTermFreqAndPositions,omitempty"`
	OmitPositions            *bool `json:"omitPositions,omitempty"`
	TermVectors              *bool `json:"termVectors,omitempty"`
	TermPositions            *bool `json:"termPositions,omitempty"`
	TermOffsets              *bool `json:"termOffsets,omitempty"`
	TermPayloads             *bool `json:"termPayloads,omitempty"`
	Required                 *bool `json:"required,omitempty"`
	UseDocValuesAsStored     *bool `json:"useDocValuesAsStored,omitempty"`
	Large                    *bool `json:"large,omitempty"`
}

type Field struct {
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	Default interface{} `json:"default,omitempty"`
	FieldDefaultProperties
}

type CopyField struct {
	Source   string `json:"source"`
	Dest     string `json:"dest"`
	MaxChars int    `json:"maxChars,omitempty"`
}

// DynamicField is just like a regular field except it has a name with a wildcard in it.
// For more info: https://lucene.apache.org/solr/guide/8_5/dynamic-fields.html
type DynamicField struct {
	Field
}

// SchemaCommand is used to restrict the available update commands that can
// be included in the body of a request to the `/update` endpoint.
type SchemaCommand string

func (c SchemaCommand) String() string {
	return string(c)
}

type schemaBuilder struct {
	commands map[SchemaCommand]interface{}
}

func newSchemaBuilder() *schemaBuilder {
	commands := make(map[SchemaCommand]interface{})
	return &schemaBuilder{commands: commands}
}

func (b *schemaBuilder) add(command SchemaCommand, item interface{}) {
	b.commands[command] = item
}

func (b *schemaBuilder) del(command SchemaCommand, name string) {
	b.commands[command] = map[string]string{"name": name}
}

// NewSchemaAPI
func NewSchemaAPI(ctx context.Context, host, core string, client *http.Client) (*SchemaAPI, error) {
	if host == "" || core == "" {
		return nil, ErrInvalidConfig
	}
	conn := &Connection{
		Host:       host,
		Core:       core,
		httpClient: client,
	}
	path := formatBasePath(host, core) + "/schema"
	return &SchemaAPI{conn: conn, Path: path}, nil
}

func (s *SchemaAPI) post() (*Response, error) {
	return nil, nil
}

// RetrieveSchema allows you to read how your schema has been defined. The output will
// include all fields, field types, dynamic rules and copy field rules in json.
// The schema name and version are also included.
func (s *SchemaAPI) RetrieveSchema(ctx context.Context) (*Response, error) {
	return request(ctx, s.conn.httpClient, http.MethodGet, s.Path, nil)
}

func (s *SchemaAPI) AddFieldType(ctx context.Context, ft *FieldType) (*Response, error) {
	sb := newSchemaBuilder()
	sb.add(SchemaCommandAddFieldType, ft)

	bodyBytes, err := interfaceToBytes(sb.commands)
	if err != nil {
		return nil, err
	}
	return request(ctx, s.conn.httpClient, http.MethodPost, s.Path, bodyBytes)
}

func (s *SchemaAPI) ReplaceFieldType(ctx context.Context, ft *FieldType) (*Response, error) {
	sb := newSchemaBuilder()
	sb.add(SchemaCommandReplaceFieldType, ft)

	bodyBytes, err := interfaceToBytes(sb.commands)
	if err != nil {
		return nil, err
	}
	return request(ctx, s.conn.httpClient, http.MethodPost, s.Path, bodyBytes)
}

func (s *SchemaAPI) DeleteFieldType(ctx context.Context, name string) (*Response, error) {
	sb := newSchemaBuilder()
	sb.del(SchemaCommandDeleteFieldType, name)

	bodyBytes, err := interfaceToBytes(sb.commands)
	if err != nil {
		return nil, err
	}
	return request(ctx, s.conn.httpClient, http.MethodPost, s.Path, bodyBytes)
}

func (s *SchemaAPI) GetFieldType(ctx context.Context, name string) (*FieldType, error) {
	res, err := s.RetrieveSchema(ctx)
	if err != nil {
		return nil, err
	}

	if res.Schema != nil && len(res.Schema.FieldTypes) > 0 {
		for _, ft := range res.Schema.FieldTypes {
			if ft.Name == name {
				return ft, nil
			}
		}
	}

	return nil, fmt.Errorf("not found")
}
