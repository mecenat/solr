package solr

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
)

type Command string

func (c Command) String() string {
	return string(c)
}

type UpdateBuilder struct {
	commands map[Command]map[string]interface{}
}

func NewUpdateBuilder() *UpdateBuilder {
	commands := make(map[Command]map[string]interface{})
	return &UpdateBuilder{commands: commands}
}

func (b *UpdateBuilder) Add(item map[string]interface{}) {
	b.commands[CommandAdd] = item
}

func (b *UpdateBuilder) Delete(doc map[string]interface{}) {
	b.commands[CommandDelete] = doc
}

func (b *UpdateBuilder) Commit() {
	b.commands[CommandCommit] = map[string]interface{}{}
}

func (b *UpdateBuilder) Rollback() {
	b.commands[CommandRollback] = map[string]interface{}{}
}

type Fields struct {
	fields map[string]interface{}
}

func NewUpdate(id string) *Fields {
	fields := make(map[string]interface{})
	fields["id"] = id
	return &Fields{fields: fields}
}

func (f *Fields) Set(key string, val interface{}) {
	f.fields[key] = map[string]interface{}{ActionSet: val}
}

func (f *Fields) Add(key string, val interface{}) {
	f.fields[key] = map[string]interface{}{ActionAdd: val}
}

func (f *Fields) AddDistinct(key string, val interface{}) {
	f.fields[key] = map[string]interface{}{ActionAddDistinct: val}
}

func (f *Fields) Remove(key string, val interface{}) {
	f.fields[key] = map[string]interface{}{ActionRemove: val}
}

func (f *Fields) RemoveRegex(key string, val interface{}) {
	f.fields[key] = map[string]interface{}{ActionRemoveRegex: val}
}

func (f *Fields) IncrementBy(key string, val int) {
	f.fields[key] = map[string]interface{}{ActionIncrement: val}
}
