// Code generated by SQLBoiler 4.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/satori/go.uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Conversation is an object representing the database table.
type Conversation struct {
	ID        int         `db:"id" boil:"id" json:"id" toml:"id" yaml:"id"`
	UUID      uuid.UUID   `db:"uuid" boil:"uuid" json:"uuid" toml:"uuid" yaml:"uuid"`
	Name      null.String `db:"name" boil:"name" json:"name,omitempty" toml:"name" yaml:"name,omitempty"`
	CreatedAt time.Time   `db:"created_at" boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *conversationR `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
	L conversationL  `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ConversationColumns = struct {
	ID        string
	UUID      string
	Name      string
	CreatedAt string
}{
	ID:        "id",
	UUID:      "uuid",
	Name:      "name",
	CreatedAt: "created_at",
}

var ConversationTableColumns = struct {
	ID        string
	UUID      string
	Name      string
	CreatedAt string
}{
	ID:        "conversation.id",
	UUID:      "conversation.uuid",
	Name:      "conversation.name",
	CreatedAt: "conversation.created_at",
}

// Generated where

var ConversationWhere = struct {
	ID        whereHelperint
	UUID      whereHelperuuid_UUID
	Name      whereHelpernull_String
	CreatedAt whereHelpertime_Time
}{
	ID:        whereHelperint{field: "\"conversation\".\"id\""},
	UUID:      whereHelperuuid_UUID{field: "\"conversation\".\"uuid\""},
	Name:      whereHelpernull_String{field: "\"conversation\".\"name\""},
	CreatedAt: whereHelpertime_Time{field: "\"conversation\".\"created_at\""},
}

// ConversationRels is where relationship names are stored.
var ConversationRels = struct {
	Messages          string
	UserConversations string
}{
	Messages:          "Messages",
	UserConversations: "UserConversations",
}

// conversationR is where relationships are stored.
type conversationR struct {
	Messages          MessageSlice          `db:"Messages" boil:"Messages" json:"Messages" toml:"Messages" yaml:"Messages"`
	UserConversations UserConversationSlice `db:"UserConversations" boil:"UserConversations" json:"UserConversations" toml:"UserConversations" yaml:"UserConversations"`
}

// NewStruct creates a new relationship struct
func (*conversationR) NewStruct() *conversationR {
	return &conversationR{}
}

// conversationL is where Load methods for each relationship are stored.
type conversationL struct{}

var (
	conversationAllColumns            = []string{"id", "uuid", "name", "created_at"}
	conversationColumnsWithoutDefault = []string{"name"}
	conversationColumnsWithDefault    = []string{"id", "uuid", "created_at"}
	conversationPrimaryKeyColumns     = []string{"id"}
)

type (
	// ConversationSlice is an alias for a slice of pointers to Conversation.
	// This should almost always be used instead of []Conversation.
	ConversationSlice []*Conversation
	// ConversationHook is the signature for custom Conversation hook methods
	ConversationHook func(context.Context, boil.ContextExecutor, *Conversation) error

	conversationQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	conversationType                 = reflect.TypeOf(&Conversation{})
	conversationMapping              = queries.MakeStructMapping(conversationType)
	conversationPrimaryKeyMapping, _ = queries.BindMapping(conversationType, conversationMapping, conversationPrimaryKeyColumns)
	conversationInsertCacheMut       sync.RWMutex
	conversationInsertCache          = make(map[string]insertCache)
	conversationUpdateCacheMut       sync.RWMutex
	conversationUpdateCache          = make(map[string]updateCache)
	conversationUpsertCacheMut       sync.RWMutex
	conversationUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var conversationBeforeInsertHooks []ConversationHook
var conversationBeforeUpdateHooks []ConversationHook
var conversationBeforeDeleteHooks []ConversationHook
var conversationBeforeUpsertHooks []ConversationHook

var conversationAfterInsertHooks []ConversationHook
var conversationAfterSelectHooks []ConversationHook
var conversationAfterUpdateHooks []ConversationHook
var conversationAfterDeleteHooks []ConversationHook
var conversationAfterUpsertHooks []ConversationHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Conversation) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range conversationBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Conversation) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range conversationBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Conversation) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range conversationBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Conversation) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range conversationBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Conversation) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range conversationAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Conversation) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range conversationAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Conversation) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range conversationAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Conversation) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range conversationAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Conversation) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range conversationAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddConversationHook registers your hook function for all future operations.
func AddConversationHook(hookPoint boil.HookPoint, conversationHook ConversationHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		conversationBeforeInsertHooks = append(conversationBeforeInsertHooks, conversationHook)
	case boil.BeforeUpdateHook:
		conversationBeforeUpdateHooks = append(conversationBeforeUpdateHooks, conversationHook)
	case boil.BeforeDeleteHook:
		conversationBeforeDeleteHooks = append(conversationBeforeDeleteHooks, conversationHook)
	case boil.BeforeUpsertHook:
		conversationBeforeUpsertHooks = append(conversationBeforeUpsertHooks, conversationHook)
	case boil.AfterInsertHook:
		conversationAfterInsertHooks = append(conversationAfterInsertHooks, conversationHook)
	case boil.AfterSelectHook:
		conversationAfterSelectHooks = append(conversationAfterSelectHooks, conversationHook)
	case boil.AfterUpdateHook:
		conversationAfterUpdateHooks = append(conversationAfterUpdateHooks, conversationHook)
	case boil.AfterDeleteHook:
		conversationAfterDeleteHooks = append(conversationAfterDeleteHooks, conversationHook)
	case boil.AfterUpsertHook:
		conversationAfterUpsertHooks = append(conversationAfterUpsertHooks, conversationHook)
	}
}

// One returns a single conversation record from the query.
func (q conversationQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Conversation, error) {
	o := &Conversation{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for conversation")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Conversation records from the query.
func (q conversationQuery) All(ctx context.Context, exec boil.ContextExecutor) (ConversationSlice, error) {
	var o []*Conversation

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Conversation slice")
	}

	if len(conversationAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Conversation records in the query.
func (q conversationQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count conversation rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q conversationQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if conversation exists")
	}

	return count > 0, nil
}

// Messages retrieves all the message's Messages with an executor.
func (o *Conversation) Messages(mods ...qm.QueryMod) messageQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"message\".\"conversation_id\"=?", o.ID),
	)

	query := Messages(queryMods...)
	queries.SetFrom(query.Query, "\"message\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"message\".*"})
	}

	return query
}

// UserConversations retrieves all the user_conversation's UserConversations with an executor.
func (o *Conversation) UserConversations(mods ...qm.QueryMod) userConversationQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"user_conversation\".\"conversation_id\"=?", o.ID),
	)

	query := UserConversations(queryMods...)
	queries.SetFrom(query.Query, "\"user_conversation\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"user_conversation\".*"})
	}

	return query
}

// LoadMessages allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (conversationL) LoadMessages(ctx context.Context, e boil.ContextExecutor, singular bool, maybeConversation interface{}, mods queries.Applicator) error {
	var slice []*Conversation
	var object *Conversation

	if singular {
		object = maybeConversation.(*Conversation)
	} else {
		slice = *maybeConversation.(*[]*Conversation)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &conversationR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &conversationR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`message`),
		qm.WhereIn(`message.conversation_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load message")
	}

	var resultSlice []*Message
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice message")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on message")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for message")
	}

	if len(messageAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Messages = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &messageR{}
			}
			foreign.R.Conversation = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.ConversationID {
				local.R.Messages = append(local.R.Messages, foreign)
				if foreign.R == nil {
					foreign.R = &messageR{}
				}
				foreign.R.Conversation = local
				break
			}
		}
	}

	return nil
}

// LoadUserConversations allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (conversationL) LoadUserConversations(ctx context.Context, e boil.ContextExecutor, singular bool, maybeConversation interface{}, mods queries.Applicator) error {
	var slice []*Conversation
	var object *Conversation

	if singular {
		object = maybeConversation.(*Conversation)
	} else {
		slice = *maybeConversation.(*[]*Conversation)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &conversationR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &conversationR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`user_conversation`),
		qm.WhereIn(`user_conversation.conversation_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load user_conversation")
	}

	var resultSlice []*UserConversation
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice user_conversation")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on user_conversation")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for user_conversation")
	}

	if len(userConversationAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.UserConversations = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &userConversationR{}
			}
			foreign.R.Conversation = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.ConversationID {
				local.R.UserConversations = append(local.R.UserConversations, foreign)
				if foreign.R == nil {
					foreign.R = &userConversationR{}
				}
				foreign.R.Conversation = local
				break
			}
		}
	}

	return nil
}

// AddMessages adds the given related objects to the existing relationships
// of the conversation, optionally inserting them as new records.
// Appends related to o.R.Messages.
// Sets related.R.Conversation appropriately.
func (o *Conversation) AddMessages(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Message) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.ConversationID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"message\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"conversation_id"}),
				strmangle.WhereClause("\"", "\"", 2, messagePrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.ConversationID = o.ID
		}
	}

	if o.R == nil {
		o.R = &conversationR{
			Messages: related,
		}
	} else {
		o.R.Messages = append(o.R.Messages, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &messageR{
				Conversation: o,
			}
		} else {
			rel.R.Conversation = o
		}
	}
	return nil
}

// AddUserConversations adds the given related objects to the existing relationships
// of the conversation, optionally inserting them as new records.
// Appends related to o.R.UserConversations.
// Sets related.R.Conversation appropriately.
func (o *Conversation) AddUserConversations(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*UserConversation) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.ConversationID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"user_conversation\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"conversation_id"}),
				strmangle.WhereClause("\"", "\"", 2, userConversationPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.ConversationID = o.ID
		}
	}

	if o.R == nil {
		o.R = &conversationR{
			UserConversations: related,
		}
	} else {
		o.R.UserConversations = append(o.R.UserConversations, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &userConversationR{
				Conversation: o,
			}
		} else {
			rel.R.Conversation = o
		}
	}
	return nil
}

// Conversations retrieves all the records using an executor.
func Conversations(mods ...qm.QueryMod) conversationQuery {
	mods = append(mods, qm.From("\"conversation\""))
	return conversationQuery{NewQuery(mods...)}
}

// FindConversation retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindConversation(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*Conversation, error) {
	conversationObj := &Conversation{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"conversation\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, conversationObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from conversation")
	}

	if err = conversationObj.doAfterSelectHooks(ctx, exec); err != nil {
		return conversationObj, err
	}

	return conversationObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Conversation) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no conversation provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(conversationColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	conversationInsertCacheMut.RLock()
	cache, cached := conversationInsertCache[key]
	conversationInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			conversationAllColumns,
			conversationColumnsWithDefault,
			conversationColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(conversationType, conversationMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(conversationType, conversationMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"conversation\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"conversation\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into conversation")
	}

	if !cached {
		conversationInsertCacheMut.Lock()
		conversationInsertCache[key] = cache
		conversationInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Conversation.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Conversation) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	conversationUpdateCacheMut.RLock()
	cache, cached := conversationUpdateCache[key]
	conversationUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			conversationAllColumns,
			conversationPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update conversation, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"conversation\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, conversationPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(conversationType, conversationMapping, append(wl, conversationPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update conversation row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for conversation")
	}

	if !cached {
		conversationUpdateCacheMut.Lock()
		conversationUpdateCache[key] = cache
		conversationUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q conversationQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for conversation")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for conversation")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ConversationSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), conversationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"conversation\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, conversationPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in conversation slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all conversation")
	}
	return rowsAff, nil
}

// Delete deletes a single Conversation record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Conversation) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Conversation provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), conversationPrimaryKeyMapping)
	sql := "DELETE FROM \"conversation\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from conversation")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for conversation")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q conversationQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no conversationQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from conversation")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for conversation")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ConversationSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(conversationBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), conversationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"conversation\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, conversationPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from conversation slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for conversation")
	}

	if len(conversationAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Conversation) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindConversation(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ConversationSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ConversationSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), conversationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"conversation\".* FROM \"conversation\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, conversationPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ConversationSlice")
	}

	*o = slice

	return nil
}

// ConversationExists checks if the Conversation row exists.
func ConversationExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"conversation\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if conversation exists")
	}

	return exists, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Conversation) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no conversation provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(conversationColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	conversationUpsertCacheMut.RLock()
	cache, cached := conversationUpsertCache[key]
	conversationUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			conversationAllColumns,
			conversationColumnsWithDefault,
			conversationColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			conversationAllColumns,
			conversationPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert conversation, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(conversationPrimaryKeyColumns))
			copy(conflict, conversationPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"conversation\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(conversationType, conversationMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(conversationType, conversationMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert conversation")
	}

	if !cached {
		conversationUpsertCacheMut.Lock()
		conversationUpsertCache[key] = cache
		conversationUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}
