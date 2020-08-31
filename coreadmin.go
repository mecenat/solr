package solr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// CoreAdmin Option & Action constants
const (
	CoreAdminOptionIndexInfo         = "indexInfo"
	CoreAdminOptionName              = "name"
	CoreAdminOptionInstanceDir       = "instanceDir"
	CoreAdminOptionConfig            = "config"
	CoreAdminOptionSchema            = "schema"
	CoreAdminOptionDataDir           = "dataDir"
	CoreAdminOptionConfigSet         = "configSet"
	CoreAdminOptionCollection        = "collection"
	CoreAdminOptionShard             = "shard"
	CoreAdminOptionAsync             = "async"
	CoreAdminOptionCore              = "core"
	CoreAdminOptionOther             = "other"
	CoreAdminOptionAction            = "action"
	CoreAdminOptionDeleteIndex       = "deleteIndex"
	CoreAdminOptionDeleteDataDir     = "deleteDataDir"
	CoreAdminOptionDeleteInstanceDir = "deleteInstanceDir"
	CoreAdminOptionIndexDir          = "indexDir"
	CoreAdminOptionSourceCore        = "srcCore"
	CoreAdminOptionPath              = "path"
	CoreAdminOptionTargetCore        = "targetCore"
	CoreAdminOptionRanges            = "ranges"
	CoreAdminOptionSplitKey          = "split.key"
	CoreAdminOptionRequestID         = "requestid"
	CoreAdminActionStatus            = "STATUS"
	CoreAdminActionCreate            = "CREATE"
	CoreAdminActionReload            = "RELOAD"
	CoreAdminActionRename            = "RENAME"
	CoreAdminActionSwap              = "SWAP"
	CoreAdminActionUnload            = "UNLOAD"
	CoreAdminActionMergeIndexes      = "MERGEINDEXES"
	CoreAdminActionSplit             = "SPLIT"
	CoreAdminActionRequestStatus     = "REQUESTSTATUS"
	CoreAdminActionRecover           = "REQUESTRECOVERY"
)

// Errors that can be returned
var (
	ErrMoreParamsPath  = errors.New("only one of path, targetCore may be defined")
	ErrMoreParamsRange = errors.New("only one of range, split.key may be defined")
)

type CoreCreateOpts struct {
	InstanceDir string
	Config      string
	Schema      string
	DataDir     string
	ConfigSet   string
	Collection  string
	Shard       string
	AsyncID     string
}

type CoreUnloadOpts struct {
	DeleteIndex       bool
	DeleteDataDir     bool
	DeleteInstanceDir bool
	AsyncID           string
}

type CoreSplitOpts struct {
	Path       []string
	TargetCore []string
	Ranges     string
	SplitKey   string
	AsyncID    string
}

type CoreMergeOpts struct {
	IndexDir []string
	SrcCore  []string
	AsyncID  string
}

type CoreAdminResponse struct {
	Header       *ResponseHeader                `json:"responseHeader"`
	Error        *ResponseError                 `json:"error"`
	Status       map[string]*CoreStatusResponse `json:"status"`
	ReqStatus    string                         `json:"STATUS"`
	Response     interface{}                    `json:"response"`
	InitFailures interface{}                    `json:"initFailures"`
	Core         string                         `json:"core"`
}

type CoreStatusResponse struct {
	Name        string        `json:"name"`
	InstanceDir string        `json:"instanceDir"`
	DataDir     string        `json:"dataDir"`
	Config      string        `json:"config"`
	Schema      string        `json:"schema"`
	StartTime   time.Time     `json:"startTime"`
	Uptime      time.Duration `json:"uptime"`
	Index       *IndexData    `json:"index"`
}

type IndexData struct {
	NumDocs                 int64     `json:"numDocs"`
	MaxDoc                  int64     `json:"maxDoc"`
	DeletedDocs             int64     `json:"deletedDocs"`
	IndexHeapUsageBytes     int64     `json:"indexHeapUsageBytes"`
	Version                 int64     `json:"version"`
	SegmentCount            int64     `json:"segmentCount"`
	Current                 bool      `json:"current"`
	HasDeletions            bool      `json:"hasDeletions"`
	Directory               string    `json:"directory"`
	SegmentsFile            string    `json:"segmentsFile"`
	SegmentsFileSizeInBytes int64     `json:"segmentsFileSizeInBytes"`
	UserData                *UserData `json:"userData"`
	LastModified            time.Time `json:"lastModified"`
	SizeInBytes             int64     `json:"sizeInBytes"`
	Size                    string    `json:"size"`
}

type UserData struct {
	CommitCommandVersion string `json:"commitCommandVer"`
	CommitTimeMSec       string `json:"commitTimeMSec"`
}

// CoreAdmin contains a connectin to solr.
type CoreAdmin struct {
	conn *Connection
	Path string
}

// NewCoreAdmin returns a new core admin, creating a connection to solr using the provided
// http client and host, core info.
func NewCoreAdmin(ctx context.Context, host string, client *http.Client) (*CoreAdmin, error) {
	if host == "" {
		return nil, ErrInvalidConfig
	}

	_, err := url.ParseRequestURI(host)
	if err != nil {
		return nil, err
	}

	conn := &Connection{
		Host:       host,
		Core:       "",
		httpClient: client,
	}
	path := fmt.Sprintf("%s/solr/admin/cores?", host)

	return &CoreAdmin{conn: conn, Path: path}, nil
}

// SetBasicAuth sets the authentication credentials if needed.
func (a *CoreAdmin) SetBasicAuth(username, password string) {
	a.conn.Username = username
	a.conn.Password = password
}

func (a *CoreAdmin) request(ctx context.Context, method, url string) (*CoreAdminResponse, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	if a.conn.Username != "" && a.conn.Password != "" {
		req.SetBasicAuth(a.conn.Username, a.conn.Password)
	}

	res, err := a.conn.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	var r CoreAdminResponse
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

// Status returns the status of all running Solr cores, or status for only the named core. If the
// noIndexInfo option is true information about the index will not be returned with a core.
// For more info:
// https://lucene.apache.org/solr/guide/8_5/coreadmin-api.html#coreadmin-status
func (a *CoreAdmin) Status(ctx context.Context, core string, noIndexInfo bool) (*CoreAdminResponse, error) {
	params := url.Values{}
	params.Set(CoreAdminOptionAction, CoreAdminActionStatus)
	if core != "" {
		params.Set(CoreAdminOptionCore, core)
	}
	if noIndexInfo {
		params.Set(CoreAdminOptionIndexInfo, "false")
	}
	url := a.Path + params.Encode()
	return a.request(ctx, http.MethodGet, url)
}

// Create creates a new core and registers it. For more info:
// https://lucene.apache.org/solr/guide/8_5/coreadmin-api.html#coreadmin-create
func (a *CoreAdmin) Create(ctx context.Context, name string, opts *CoreCreateOpts) (*CoreAdminResponse, error) {
	params := url.Values{}
	params.Set(CoreAdminOptionAction, CoreAdminActionCreate)
	params.Set(CoreAdminOptionName, name)
	if opts != nil {
		if opts.AsyncID != "" {
			params.Set(CoreAdminOptionAsync, opts.AsyncID)
		}
		if opts.DataDir != "" {
			params.Set(CoreAdminOptionDataDir, opts.DataDir)
		}
		if opts.InstanceDir != "" {
			params.Set(CoreAdminOptionInstanceDir, opts.InstanceDir)
		}
		if opts.Collection != "" {
			params.Set(CoreAdminOptionCollection, opts.Collection)
		}
		if opts.ConfigSet != "" {
			params.Set(CoreAdminOptionConfigSet, opts.ConfigSet)
		}
		if opts.Config != "" {
			params.Set(CoreAdminOptionConfig, opts.Config)
		}
		if opts.Schema != "" {
			params.Set(CoreAdminOptionSchema, opts.Schema)
		}
		if opts.Shard != "" {
			params.Set(CoreAdminOptionShard, opts.Shard)
		}
	}
	url := a.Path + params.Encode()
	return a.request(ctx, http.MethodGet, url)
}

// Reload loads a new core from the configuration of an existing, registered Solr core. While the
// new core is initializing, the existing one will continue to handle requests. When the new
// Solr core is ready, it takes over and the old core is unloaded. For More info:
// https://lucene.apache.org/solr/guide/8_5/coreadmin-api.html#coreadmin-reload
func (a *CoreAdmin) Reload(ctx context.Context, core string) (*CoreAdminResponse, error) {
	params := url.Values{}
	params.Set(CoreAdminOptionAction, CoreAdminActionReload)
	params.Set(CoreAdminOptionCore, core)
	url := a.Path + params.Encode()
	return a.request(ctx, http.MethodGet, url)
}

// Rename changes the name of a Solr core. An asyncID may be provided in order to track
// this action which will be processed asynchronously. For more info:
// https://lucene.apache.org/solr/guide/8_5/coreadmin-api.html#coreadmin-rename
func (a *CoreAdmin) Rename(ctx context.Context, core, other, asyncID string) (*CoreAdminResponse, error) {
	params := url.Values{}
	params.Set(CoreAdminOptionAction, CoreAdminActionRename)
	params.Set(CoreAdminOptionCore, core)
	params.Set(CoreAdminOptionOther, other)
	if asyncID != "" {
		params.Set(CoreAdminOptionAsync, asyncID)
	}
	url := a.Path + params.Encode()
	return a.request(ctx, http.MethodGet, url)
}

// Swap atomically swaps the names used to access two existing Solr cores. This can be used
// to swap new content into production. For more info:
// https://lucene.apache.org/solr/guide/8_5/coreadmin-api.html#coreadmin-swap
func (a *CoreAdmin) Swap(ctx context.Context, core, other, asyncID string) (*CoreAdminResponse, error) {
	params := url.Values{}
	params.Set(CoreAdminOptionAction, CoreAdminActionSwap)
	params.Set(CoreAdminOptionCore, core)
	params.Set(CoreAdminOptionOther, other)
	if asyncID != "" {
		params.Set(CoreAdminOptionAsync, asyncID)
	}
	url := a.Path + params.Encode()
	return a.request(ctx, http.MethodGet, url)
}

// Unload removes a core from Solr. Requires the name of the core to be unloaded.
// For more info:
// https://lucene.apache.org/solr/guide/8_5/coreadmin-api.html#coreadmin-unload
func (a *CoreAdmin) Unload(ctx context.Context, core string, opts *CoreUnloadOpts) (*CoreAdminResponse, error) {
	params := url.Values{}
	params.Set(CoreAdminOptionAction, CoreAdminActionUnload)
	params.Set(CoreAdminOptionCore, core)
	if opts != nil {
		if opts.AsyncID != "" {
			params.Set(CoreAdminOptionAsync, opts.AsyncID)
		}
		if opts.DeleteIndex {
			params.Set(CoreAdminOptionDeleteIndex, "true")
		}
		if opts.DeleteDataDir {
			params.Set(CoreAdminOptionDeleteDataDir, "true")
		}
		if opts.DeleteInstanceDir {
			params.Set(CoreAdminOptionDeleteInstanceDir, "true")
		}
	}
	url := a.Path + params.Encode()
	return a.request(ctx, http.MethodGet, url)
}

// Merge merges one or more indexes to another index. The target core index must already exist
// and have a compatible schema with the one or more indexes that will be merged to it.
// Another commit on the target core should also be performed after the merge
// is complete. For more info:
// https://lucene.apache.org/solr/guide/8_5/coreadmin-api.html#coreadmin-mergeindexes
func (a *CoreAdmin) Merge(ctx context.Context, core string, opts *CoreMergeOpts) (*CoreAdminResponse, error) {
	params := url.Values{}
	params.Set(CoreAdminOptionAction, CoreAdminActionMergeIndexes)
	params.Set(CoreAdminOptionCore, core)
	if opts != nil {
		if opts.AsyncID != "" {
			params.Set(CoreAdminOptionAsync, opts.AsyncID)
		}
		if len(opts.IndexDir) > 0 {
			for _, d := range opts.IndexDir {
				params.Add(CoreAdminOptionIndexDir, d)
			}
		}
		if len(opts.SrcCore) > 0 {
			for _, c := range opts.SrcCore {
				params.Add(CoreAdminOptionSourceCore, c)
			}
		}
	}
	url := a.Path + params.Encode()
	return a.request(ctx, http.MethodGet, url)
}

// Split splits an index into two or more indexes. The index being split can continue to handle requests.
// The split pieces can be placed into a specified directory on the server’s filesystem or it can be
// merged into running Solr cores. For more info:
// https://lucene.apache.org/solr/guide/8_5/coreadmin-api.html#coreadmin-split
func (a *CoreAdmin) Split(ctx context.Context, core string, opts *CoreSplitOpts) (*CoreAdminResponse, error) {
	params := url.Values{}
	params.Set(CoreAdminOptionAction, CoreAdminActionSplit)
	params.Set(CoreAdminOptionCore, core)
	if opts != nil {
		if len(opts.Path) > 0 && len(opts.TargetCore) > 0 {
			return nil, ErrMoreParamsPath
		}
		if opts.Ranges != "" && opts.SplitKey != "" {
			return nil, ErrMoreParamsRange
		}
		if opts.AsyncID != "" {
			params.Set(CoreAdminOptionAsync, opts.AsyncID)
		}
		if len(opts.Path) > 0 {
			for _, p := range opts.Path {
				params.Add(CoreAdminOptionPath, p)
			}
		}
		if opts.Ranges != "" {
			params.Set(CoreAdminOptionRanges, opts.Ranges)
		}
		if len(opts.TargetCore) > 0 {
			for _, t := range opts.TargetCore {
				params.Add(CoreAdminOptionTargetCore, t)
			}
		}
		if opts.SplitKey != "" {
			params.Set(CoreAdminOptionSplitKey, opts.SplitKey)
		}
	}
	url := a.Path + params.Encode()
	return a.request(ctx, http.MethodGet, url)
}

// RequestStatus returns the status of an already submitted asynchronous CoreAdmin API call.
// For more info:
// https://lucene.apache.org/solr/guide/8_5/coreadmin-api.html#coreadmin-requeststatus
func (a *CoreAdmin) RequestStatus(ctx context.Context, id string) (*CoreAdminResponse, error) {
	params := url.Values{}
	params.Set(CoreAdminOptionAction, CoreAdminActionRequestStatus)
	params.Set(CoreAdminOptionRequestID, id)
	url := a.Path + params.Encode()
	return a.request(ctx, http.MethodGet, url)
}

// Recover manually asks a core to recover by synching with the leader. This should be considered
// an "expert" level command and should be used in situations where the node (SorlCloud replica)
// is unable to become active automatically. For more info
// https://lucene.apache.org/solr/guide/8_5/coreadmin-api.html#coreadmin-requestrecovery
func (a *CoreAdmin) Recover(ctx context.Context, core string) (*CoreAdminResponse, error) {
	params := url.Values{}
	params.Set(CoreAdminOptionAction, CoreAdminActionRecover)
	params.Set(CoreAdminOptionCore, core)
	url := a.Path + params.Encode()
	return a.request(ctx, http.MethodGet, url)
}
