package solr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ManagedResponse represents the response from solr's managed resources API.
// Header and Error (if there is any) will always be populated. The rest
// are helpers on specific cases. Currently supported cases are when
// requesting for a list of all managed resources, and for
// a managed synonyms list.
type ManagedResponse struct {
	Header    *ResponseHeader    `json:"responseHeader"`
	Error     *ResponseError     `json:"error"`
	Resources []*ManagedResource `json:"managedResources"`
	Synonyms  *SynonymMappings   `json:"synonymMappings"`
	RawMap    map[string]interface{}
}

// UnmarshalJSON implements the unmarshaler interface
func (r *ManagedResponse) UnmarshalJSON(b []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}
	r.RawMap = m

	var h ResponseHeader
	err = json.Unmarshal(b, &h)
	if err != nil {
		return err
	}
	r.Header = &h

	_, ok := m["error"]
	if ok {
		var e ResponseError
		err = json.Unmarshal(b, &e)
		if err != nil {
			return err
		}
		r.Error = &e
	}

	resInfArr, ok := m["managedResources"]
	if ok {
		var resourceMap []*ManagedResource
		resArr, ok := resInfArr.([]interface{})
		if ok {
			for _, resInf := range resArr {
				var resource ManagedResource
				resBytes, err := interfaceToBytes(resInf)
				if err != nil {
					return err
				}
				err = json.Unmarshal(resBytes, &resource)
				if err != nil {
					return err
				}
				resourceMap = append(resourceMap, &resource)
			}
		}
		r.Resources = resourceMap
	}

	_, ok = m["synonymMappings"]
	if ok {
		var syn SynonymMappings
		err = json.Unmarshal(b, &syn)
		if err != nil {
			return err
		}
		r.Synonyms = &syn
	}

	return nil
}

// ManagedResource represents a managed resource in solr.
type ManagedResource struct {
	ID              string `json:"resourceId"`
	Class           string `json:"class"`
	ObserversNumber string `json:"numObservers"`
}

// SynonymMappings is a helper struct for navigating a synonyms managed list.
type SynonymMappings struct {
	InitArgs   *SynonymInitArgs    `json:"initArgs"`
	InitOn     time.Time           `json:"initializedOn"`
	UpdatedOn  time.Time           `json:"updatedSinceInit"`
	ManagedMap map[string][]string `json:"managedMap"`
}

// SynonymInitArgs are the initialization arguments for a synonyms
// managed list.
type SynonymInitArgs struct {
	IgnoreCase bool `json:"ignoreCase"`
}

// ManagedAPI contains a connection to solr
type ManagedAPI struct {
	conn     *Connection
	BasePath string
}

// NewManagedAPI returns a new Managed Resources API, creating a connection to solr using the provided
// http client, host and core info.
// https://lucene.apache.org/solr/guide/8_5/managed-resources.html#managed-resources-overview
func NewManagedAPI(ctx context.Context, host, core string, client *http.Client) (*ManagedAPI, error) {
	if host == "" || core == "" {
		return nil, ErrInvalidConfig
	}

	_, err := url.ParseRequestURI(host)
	if err != nil {
		return nil, err
	}

	conn := &Connection{
		Host:       host,
		Core:       core,
		httpClient: client,
	}

	path := fmt.Sprintf("%s/schema", formatBasePath(host, core))

	return &ManagedAPI{conn: conn, BasePath: path}, nil
}

// SetBasicAuth sets the authentication credentials if needed.
func (m *ManagedAPI) SetBasicAuth(username, password string) {
	m.conn.Username = username
	m.conn.Password = password
}

func (m *ManagedAPI) request(ctx context.Context, method, url string, body []byte) (*ManagedResponse, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	if m.conn.Username != "" && m.conn.Password != "" {
		req.SetBasicAuth(m.conn.Username, m.conn.Password)
	}

	res, err := m.conn.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	var r ManagedResponse
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	if r.Error != nil {
		return &r, r.Error
	}

	return &r, nil
}

func (m *ManagedAPI) formatURL(path string) string {
	if strings.HasPrefix(path, "/") {
		return m.BasePath + path
	}
	return m.BasePath + "/" + path
}

// RetrieveResource returns the specified resource. Requires the path to the resource.
func (m *ManagedAPI) RetrieveResource(ctx context.Context, path string) (*ManagedResponse, error) {
	return m.request(ctx, http.MethodGet, m.formatURL(path), nil)
}

// UpsertResource updates the specified resource. Requires the path to the resource and
// the resource to be created/updated
func (m *ManagedAPI) UpsertResource(ctx context.Context, path string, data interface{}) (*ManagedResponse, error) {
	bodyBytes, err := interfaceToBytes(data)
	if err != nil {
		return nil, err
	}
	return m.request(ctx, http.MethodPut, m.formatURL(path), bodyBytes)
}

// DeleteResource deletes the specified resource. Requires the path to the resource.
func (m *ManagedAPI) DeleteResource(ctx context.Context, path string) (*ManagedResponse, error) {
	return m.request(ctx, http.MethodDelete, m.formatURL(path), nil)
}

// SetInitArgs set the initialization arguments for a managed resource. It requires the path to the managed
// resource and a map of the init arguments to update. Attention must be given to make sure that the
// arguments provided are valid init arguments, since solr doesn't check the validity during update
// but only during core reload.
func (m *ManagedAPI) SetInitArgs(ctx context.Context, path string, args map[string]interface{}) (*ManagedResponse, error) {
	body := map[string]map[string]interface{}{"initArgs": args}
	return m.UpsertResource(ctx, path, body)
}

// RestManager returns all available managed resources on the solr core.
func (m *ManagedAPI) RestManager(ctx context.Context) (*ManagedResponse, error) {
	url := m.BasePath + "/managed"
	return m.request(ctx, http.MethodGet, url, nil)
}

// SynonymSetIgnoreCase set the desired value to the ignoreCase initialization argument for
// managed synonym resources.
func (m *ManagedAPI) SynonymSetIgnoreCase(ctx context.Context, listName string, value bool) (*ManagedResponse, error) {
	path := "/analysis/synonyms/" + listName
	ign := map[string]interface{}{"ignoreCase": value}
	return m.SetInitArgs(ctx, path, ign)
}

// SynonymList returns a map of all the synonyms in the specified list.
func (m *ManagedAPI) SynonymList(ctx context.Context, listName string) (*ManagedResponse, error) {
	path := "/analysis/synonyms/" + listName
	return m.RetrieveResource(ctx, path)
}

// SynonymGet returns the synonym mapping for the specified word in the specified list.
func (m *ManagedAPI) SynonymGet(ctx context.Context, listName string, synonym string) (*ManagedResponse, error) {
	path := fmt.Sprintf("/analysis/synonyms/%s/%s", listName, synonym)
	return m.RetrieveResource(ctx, path)
}

// SynonymAdd adds a new synonym mapping in the specified list.
func (m *ManagedAPI) SynonymAdd(ctx context.Context, listName string, synonyms map[string][]string) (*ManagedResponse, error) {
	path := "/analysis/synonyms/" + listName
	return m.UpsertResource(ctx, path, synonyms)
}

// SynonymAddSymmetric adds a list of symmetric synonyms. These are expanded into a mapping
// for each term in the list by solr. Despite what is said in the solr docs tho, currently
// (v.8.6.1) adding a synonym slice this way does not seem to remove the word for the
// mapping, ending up having a word being a synonym of itself. SynonymAddOptimal
// is therefore recommended for this purpose.
func (m *ManagedAPI) SynonymAddSymmetric(ctx context.Context, listName string, synonyms []string) (*ManagedResponse, error) {
	path := "/analysis/synonyms/" + listName
	return m.UpsertResource(ctx, path, synonyms)
}

// SynonymAddOptimal creates a mapping for each word in the given slice, just as
// solr should be doing under the hood in AddSymmetric. The conversion from
// slice to array of maps is here handled by golang.
func (m *ManagedAPI) SynonymAddOptimal(ctx context.Context, listName string, synonyms []string) (*ManagedResponse, error) {
	path := "/analysis/synonyms/" + listName
	body := map[string][]string{}
	for i, s := range synonyms {
		rest := make([]string, len(synonyms))
		copy(rest, synonyms)
		rest[i] = rest[len(rest)-1]
		rest = rest[:len(rest)-1]
		body[s] = rest
	}
	return m.UpsertResource(ctx, path, body)
}

// SynonymDelete removes the specified mapping from the specified synonyms list.
func (m *ManagedAPI) SynonymDelete(ctx context.Context, listName string, synonym string) (*ManagedResponse, error) {
	path := fmt.Sprintf("/analysis/synonyms/%s/%s", listName, synonym)
	return m.DeleteResource(ctx, path)
}
