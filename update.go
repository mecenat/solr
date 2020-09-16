package solr

// Constants for different actions and commands used
// for the `/update` endpoint
const (
	ActionSet                 = "set"
	ActionAdd                 = "add"
	ActionAddDistinct         = "add-distinct"
	ActionRemove              = "remove"
	ActionRemoveRegex         = "removeregex"
	ActionIncrement           = "inc"
	CommandAdd        Command = "add"
	CommandDelete     Command = "delete"
	CommandCommit     Command = "commit"
	CommandRollback   Command = "rollback"
	CommandOptimize   Command = "optimize"
)

// CommitOptions are the available options to a commit update command.
type CommitOptions struct {
	DoNotWaitSearcher bool
	ExpungeDeletes    bool
}

// OptimizeOptions are the available options to an optimize update command.
type OptimizeOptions struct {
	DoNotWaitSearcher bool
	MaxSegments       int
}

// Command is used to restrict the available update commands that can
// be included in the body of a request to the `/update` endpoint.
type Command string

func (c Command) String() string {
	return string(c)
}

// UpdateBuilder is a helper struct that provides methods to
// easily populate the body of a custom `/update` request
type UpdateBuilder struct {
	additions []interface{}
	deletions []interface{}
	commands  map[Command]interface{}
}

// NewUpdateBuilder returns an initialized UpdateBuilder, a helper struct that provides methods to
// easily populate a custom request to the `/update` endpoint of the solr server, that can
// contain more than one action. Multiple additions or deletion are grouped in an array
// when sent to solr instead of a map (as seen in solr docs) for obvious reasons.
// Therefore actual action hierarchy CANNOT be achieved! It's usage is suggested
// for any cases that the methods provided by the Client does not cover.
// More info:
// https://lucene.apache.org/solr/guide/8_5/uploading-data-with-index-handlers.html#sending-json-update-commands
func NewUpdateBuilder() *UpdateBuilder {
	commands := make(map[Command]interface{})
	return &UpdateBuilder{commands: commands}
}

func (b *UpdateBuilder) add(item map[string]interface{}) {
	b.commands[CommandAdd] = formatDocEntry(item)
}

func (b *UpdateBuilder) delete(doc interface{}) {
	b.commands[CommandDelete] = doc
}

func (b *UpdateBuilder) prepare() {
	if len(b.additions) > 0 {
		b.commands[CommandAdd] = b.additions
	}
	if len(b.deletions) > 0 {
		b.commands[CommandDelete] = b.deletions
	}
}

// Add inserts an add command block to the body. The provided input
// must be valid JSON. For atomic or in-place updates it is
// recommended to use the `Update` method that is provided
// by the Client interface.
func (b *UpdateBuilder) Add(item interface{}) {
	b.additions = append(b.additions, item)
}

// DeleteByID inserts a delete command block to the body. It should
// contain a document identifying the id (uniqueKey field)
func (b *UpdateBuilder) DeleteByID(id string) {
	b.deletions = append(b.deletions, formatDeleteByID(id))
}

// DeleteByQuery inserts a delete command block to the body. It should
// contain a document identifying a query to properly work.
func (b *UpdateBuilder) DeleteByQuery(query string) {
	b.deletions = append(b.deletions, formatDeleteByQuery(query))
}

// commit inserts a commit command block to the body. Accepts options:
// DoNotWaitSearcher: By default a commit command blocks until a new
// 						searcher is opened and registered as the main query
// 						searcher, making the changes visible.
// ExpungeDeletes: Merges segments that have more than 10%
// 						deleted docs, expunging the deleted
// 						documents in the process.
// It is recommended to use the `commit` method that is provided
// by the Client interface.
func (b *UpdateBuilder) commit(opts *CommitOptions) {
	doc := map[string]interface{}{}
	if opts != nil {
		if opts.DoNotWaitSearcher {
			doc[OptionWaitSearcher] = false
		}
		if opts.ExpungeDeletes {
			doc[OptionExpungeDeletes] = true
		}
	}
	b.commands[CommandCommit] = doc
}

// rollback inserts a rollback command block to the body. The command is
// always an empty object. It is recommended to use the `rollback`
// method that is provided by the Client interface.
func (b *UpdateBuilder) rollback() {
	b.commands[CommandRollback] = map[string]interface{}{}
}

// optimize inserts an optimize command block to the body. Accepts options:
// DoNotWaitSearcher: By default a commit command blocks until a new
// 						searcher is opened and registered as the main queryja
// 						searcher, making the changes visible.
// MaxSegments: Merge the segments down to no more than this number
// 					  of segments but does not guarantee that the goal
// 						will be achieved.
// It is recommended to use the `optimize` method that is provided
// by the Client interface.
func (b *UpdateBuilder) optimize(opts *OptimizeOptions) {
	doc := map[string]interface{}{}
	if opts != nil {
		if opts.DoNotWaitSearcher {
			doc[OptionWaitSearcher] = false
		}
		if opts.MaxSegments > 1 {
			doc[OptionMaxSegments] = opts.MaxSegments
		}
	}
	b.commands[CommandOptimize] = doc
}

// UpdatedFields is a helper struct that contains the fields to be
// updated during an atomic/in-place update. It provides methods
// that allow to easily create a document to be sent to the
// `/update` endpoint using the Client's `Update` method
type UpdatedFields struct {
	fields map[string]interface{}
}

// NewUpdateDocument returns an UpdatedFields helper that is used
// to provided the fields to be updated in an atomic/in-place
// update. It requires as input the id (uniqueKey field) of
// the document to be updated in order for the update to
// be successful, if the id provided does not exist a
// new document will be created. More info:
// https://lucene.apache.org/solr/guide/8_5/updating-parts-of-documents.html
func NewUpdateDocument(id string) *UpdatedFields {
	fields := make(map[string]interface{})
	fields["id"] = id
	return &UpdatedFields{fields: fields}
}

// Set replaces or sets the field value(s) with the specified values(s).
// Takes as input a key which is the field name and a val which is
// the provided value(s) to set.
func (f *UpdatedFields) Set(key string, val interface{}) {
	f.fields[key] = map[string]interface{}{ActionSet: val}
}

// Add adds the specified value(s) to a multiValue field. Takes as input
// a key which is the field name and a val which is the provided
// value(s) to add.
func (f *UpdatedFields) Add(key string, val interface{}) {
	f.fields[key] = map[string]interface{}{ActionAdd: val}
}

// AddDistinct adds the specified value(s) to a multiValue field only if
// they are not already present. Takes as input a key which is the
// field name and a val which is the provided value(s) to add.
func (f *UpdatedFields) AddDistinct(key string, val interface{}) {
	f.fields[key] = map[string]interface{}{ActionAddDistinct: val}
}

// Remove removes the specified value(s) from a multiValue field. Takes
// as input a key which is the field name and a val which is the
// provided value(s) to remove.
func (f *UpdatedFields) Remove(key string, val interface{}) {
	f.fields[key] = map[string]interface{}{ActionRemove: val}
}

// RemoveRegex removes the specified regex(es) from a multiValue field.
// Takes as input a key which is the field name and a val which
// is the provided regex(es) to remove.
func (f *UpdatedFields) RemoveRegex(key string, val interface{}) {
	f.fields[key] = map[string]interface{}{ActionRemoveRegex: val}
}

// IncrementBy increments a numeric value by a specific amount. Takes
// as input a key which is the field name and a val which is an int
// signifying the amount to increment by.
func (f *UpdatedFields) IncrementBy(key string, val int) {
	f.fields[key] = map[string]interface{}{ActionIncrement: val}
}
