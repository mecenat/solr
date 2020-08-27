package solr

import (
	"context"
	"errors"
	"net/http"
)

// Valid commands for the schema API
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

// Various errors returned from the schema API
var (
	ErrFieldNotFound        = errors.New("field not found")
	ErrFieldTypeNotFound    = errors.New("field type not found")
	ErrDynamicFieldNotFound = errors.New("dynamic field not found")
	ErrCopyFieldNotFound    = errors.New("copy field not found")
)

// Analyzer represents the analyzer entity. An analyzer examines the text of
// fields and generates a token stream. For more info:
// https://lucene.apache.org/solr/guide/8_5/analyzers.html
type Analyzer struct {
	Tokenizer map[string]interface{}   `json:"tokenizer"`
	Filters   []map[string]interface{} `json:"filters"`
}

// FieldType represents a solr field type. A field type defines the analysis that will occur on a field when documents
// are indexed or queries are sent to the index. Only a name and the class name are mandatory. For more info:
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

// FieldDefaultProperties represents the defualt properties shared by field types and fields. These are propertries
// that can be specified either on the field types, or on individual fields to override the values provided
// by the field types. Built according to schema version 1.6. For more info:
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

// Field represents a solr field. For more info:
// https://lucene.apache.org/solr/guide/8_5/defining-fields.html#field-properties
type Field struct {
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	Default interface{} `json:"default,omitempty"`
	FieldDefaultProperties
}

// CopyField represents a solr copy field rule, solr's mechanism for making copies of fields
// so that you can apply several distinct field types to a single piece of incoming
// information. The name of the field you want to copy is the source, and the
// name of the copy is the destination. For more info:
// https://lucene.apache.org/solr/guide/8_5/copying-fields.html
type CopyField struct {
	Source   string `json:"source"`
	Dest     string `json:"dest"`
	MaxChars int    `json:"maxChars,omitempty"`
}

// DynamicField is just like a regular field except it has a name with a wildcard in it.
// For more info: https://lucene.apache.org/solr/guide/8_5/dynamic-fields.html
type DynamicField Field

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

func (b *schemaBuilder) delCopyField(source, dest string) {
	b.commands[SchemaCommandDeleteCopyField] = map[string]string{"source": source, "dest": dest}
}

// SchemaAPI contains a connection to solr and the path to it.
type SchemaAPI struct {
	conn *Connection
	Path string
}

// NewSchemaAPI returns a new schema API, creating a connection to solr using the provided
// http client and host, core info.
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

func (s *SchemaAPI) post(ctx context.Context, body interface{}) (*Response, error) {
	bodyBytes, err := interfaceToBytes(body)
	if err != nil {
		return nil, err
	}
	return request(ctx, s.conn.httpClient, http.MethodPost, s.Path, bodyBytes)
}

// RetrieveSchema allows you to read how your schema has been defined. The output will
// include all fields, field types, dynamic rules and copy field rules in json.
// The schema name and version are also included.
func (s *SchemaAPI) RetrieveSchema(ctx context.Context) (*Response, error) {
	return request(ctx, s.conn.httpClient, http.MethodGet, s.Path, nil)
}

// AddFieldType adds a new field type to the schema. For more info:
// https://lucene.apache.org/solr/guide/8_5/schema-api.html#add-a-new-field-type
func (s *SchemaAPI) AddFieldType(ctx context.Context, ft *FieldType) (*Response, error) {
	sb := newSchemaBuilder()
	sb.add(SchemaCommandAddFieldType, ft)
	return s.post(ctx, sb.commands)
}

// ReplaceFieldType replaces a field type in your schema. Note that you must supply the full definition
// for a field type - this command will not partially modify a field type’s definition. If the field
// type does not exist in the schema an error is thrown. For more info:
// https://lucene.apache.org/solr/guide/8_5/schema-api.html#replace-a-field-type
func (s *SchemaAPI) ReplaceFieldType(ctx context.Context, ft *FieldType) (*Response, error) {
	sb := newSchemaBuilder()
	sb.add(SchemaCommandReplaceFieldType, ft)
	return s.post(ctx, sb.commands)
}

// DeleteFieldType removes a field type from your schema. If the field type does not
// exist in the schema, or if any field or dynamic field rule in the schema uses
// the field type, an error is thrown. For more info:
// https://lucene.apache.org/solr/guide/8_5/schema-api.html#delete-a-field-type
func (s *SchemaAPI) DeleteFieldType(ctx context.Context, name string) (*Response, error) {
	sb := newSchemaBuilder()
	sb.del(SchemaCommandDeleteFieldType, name)
	return s.post(ctx, sb.commands)
}

// RetrieveFieldType returns the specified field type.
func (s *SchemaAPI) RetrieveFieldType(ctx context.Context, name string) (*FieldType, error) {
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

	return nil, ErrFieldTypeNotFound
}

// Field methods

// AddField adds a new field definition to your schema. If a field with the same name exists
// an error is thrown. For more info:
// https://lucene.apache.org/solr/guide/8_5/schema-api.html#add-a-new-field
func (s *SchemaAPI) AddField(ctx context.Context, fl *Field) (*Response, error) {
	sb := newSchemaBuilder()
	sb.add(SchemaCommandAddField, fl)
	return s.post(ctx, sb.commands)
}

// ReplaceField replaces a field’s definition. Note that you must supply the full definition for a
// field - this command will not partially modify a field’s definition. If the field does not
// exist in the schema an error is thrown. For more info:
// https://lucene.apache.org/solr/guide/8_5/schema-api.html#replace-a-field
func (s *SchemaAPI) ReplaceField(ctx context.Context, fl *Field) (*Response, error) {
	sb := newSchemaBuilder()
	sb.add(SchemaCommandReplaceField, fl)
	return s.post(ctx, sb.commands)
}

// DeleteField removes a field definition from your schema. If the field does not exist in the schema,
// or if the field is the source or destination of a copy field rule, an error is thrown.
// For more info:
// https://lucene.apache.org/solr/guide/8_5/schema-api.html#delete-a-field
func (s *SchemaAPI) DeleteField(ctx context.Context, name string) (*Response, error) {
	sb := newSchemaBuilder()
	sb.del(SchemaCommandDeleteField, name)
	return s.post(ctx, sb.commands)
}

// RetrieveField returns the specified field.
func (s *SchemaAPI) RetrieveField(ctx context.Context, name string) (*Field, error) {
	res, err := s.RetrieveSchema(ctx)
	if err != nil {
		return nil, err
	}

	if res.Schema != nil && len(res.Schema.Fields) > 0 {
		for _, fl := range res.Schema.Fields {
			if fl.Name == name {
				return fl, nil
			}
		}
	}

	return nil, ErrFieldNotFound
}

// Dynamic Field Methods

// AddDynamicField adds a new dynamic field rule to your schema. For more info:
// https://lucene.apache.org/solr/guide/8_5/schema-api.html#add-a-dynamic-field-rule
func (s *SchemaAPI) AddDynamicField(ctx context.Context, df *DynamicField) (*Response, error) {
	sb := newSchemaBuilder()
	sb.add(SchemaCommandAddDynamicField, df)
	return s.post(ctx, sb.commands)
}

// ReplaceDynamicField replaces a dynamic field rule in your schema. Note that you must supply the full definition
// for a dynamic field rule - this command will not partially modify a dynamic field rule’s definition. If the
// dynamic field rule does not exist in the schema an error is thrown. For more info:
// https://lucene.apache.org/solr/guide/8_5/schema-api.html#replace-a-dynamic-field-rule
func (s *SchemaAPI) ReplaceDynamicField(ctx context.Context, df *DynamicField) (*Response, error) {
	sb := newSchemaBuilder()
	sb.add(SchemaCommandReplaceDynamicField, df)
	return s.post(ctx, sb.commands)
}

// DeleteDynamicField deletes a dynamic field rule from your schema. If the dynamic field rule does not exist
// in the schema, or if the schema contains a copy field rule with a target or destination that matches
// only this dynamic field rule, an error is thrown. For more info:
// https://lucene.apache.org/solr/guide/8_5/schema-api.html#delete-a-dynamic-field-rule
func (s *SchemaAPI) DeleteDynamicField(ctx context.Context, name string) (*Response, error) {
	sb := newSchemaBuilder()
	sb.del(SchemaCommandDeleteDynamicField, name)
	return s.post(ctx, sb.commands)
}

// RetrieveDynamicField returns the specified dynamic field.
func (s *SchemaAPI) RetrieveDynamicField(ctx context.Context, name string) (*DynamicField, error) {
	res, err := s.RetrieveSchema(ctx)
	if err != nil {
		return nil, err
	}

	if res.Schema != nil && len(res.Schema.DynamicFields) > 0 {
		for _, df := range res.Schema.DynamicFields {
			if df.Name == name {
				return df, nil
			}
		}
	}

	return nil, ErrDynamicFieldNotFound
}

// Copy Field Methods

// AddCopyField adds a new copy field rule to your schema. Source and Destination are required.
// Destination is always a string so for ease of use, unlike with the json API it is not
// possible to copy a field to multiple destinations. A different copy field rule must
// be made for each. For more info:
// https://lucene.apache.org/solr/guide/8_5/schema-api.html#add-a-new-copy-field-rule
func (s *SchemaAPI) AddCopyField(ctx context.Context, cf *CopyField) (*Response, error) {
	sb := newSchemaBuilder()
	sb.add(SchemaCommandAddCopyField, cf)
	return s.post(ctx, sb.commands)
}

// DeleteCopyField deletes a copy field rule from your schema. If the copy field rule does not exist in
// the schema an error is thrown. For more info:
// https://lucene.apache.org/solr/guide/8_5/schema-api.html#delete-a-copy-field-rule
func (s *SchemaAPI) DeleteCopyField(ctx context.Context, source, dest string) (*Response, error) {
	sb := newSchemaBuilder()
	sb.delCopyField(source, dest)
	return s.post(ctx, sb.commands)
}

// RetrieveCopyField returns the specified copy field rule.
func (s *SchemaAPI) RetrieveCopyField(ctx context.Context, source, dest string) (*CopyField, error) {
	res, err := s.RetrieveSchema(ctx)
	if err != nil {
		return nil, err
	}

	if res.Schema != nil && len(res.Schema.CopyFields) > 0 {
		for _, cf := range res.Schema.CopyFields {
			if cf.Source == source && cf.Dest == dest {
				return cf, nil
			}
		}
	}

	return nil, ErrCopyFieldNotFound
}
