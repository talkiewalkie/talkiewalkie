// Code generated by SQLBoiler 4.5.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// UserWalk is an object representing the database table.
type UserWalk struct {
	ID     int `db:"id" boil:"id" json:"id" toml:"id" yaml:"id"`
	UserID int `db:"user_id" boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	WalkID int `db:"walk_id" boil:"walk_id" json:"walk_id" toml:"walk_id" yaml:"walk_id"`

	R *userWalkR `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
	L userWalkL  `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UserWalkColumns = struct {
	ID     string
	UserID string
	WalkID string
}{
	ID:     "id",
	UserID: "user_id",
	WalkID: "walk_id",
}

// Generated where

var UserWalkWhere = struct {
	ID     whereHelperint
	UserID whereHelperint
	WalkID whereHelperint
}{
	ID:     whereHelperint{field: "\"user_walk\".\"id\""},
	UserID: whereHelperint{field: "\"user_walk\".\"user_id\""},
	WalkID: whereHelperint{field: "\"user_walk\".\"walk_id\""},
}

// UserWalkRels is where relationship names are stored.
var UserWalkRels = struct {
	User string
	Walk string
}{
	User: "User",
	Walk: "Walk",
}

// userWalkR is where relationships are stored.
type userWalkR struct {
	User *User `db:"User" boil:"User" json:"User" toml:"User" yaml:"User"`
	Walk *Walk `db:"Walk" boil:"Walk" json:"Walk" toml:"Walk" yaml:"Walk"`
}

// NewStruct creates a new relationship struct
func (*userWalkR) NewStruct() *userWalkR {
	return &userWalkR{}
}

// userWalkL is where Load methods for each relationship are stored.
type userWalkL struct{}

var (
	userWalkAllColumns            = []string{"id", "user_id", "walk_id"}
	userWalkColumnsWithoutDefault = []string{"user_id", "walk_id"}
	userWalkColumnsWithDefault    = []string{"id"}
	userWalkPrimaryKeyColumns     = []string{"id"}
)

type (
	// UserWalkSlice is an alias for a slice of pointers to UserWalk.
	// This should generally be used opposed to []UserWalk.
	UserWalkSlice []*UserWalk
	// UserWalkHook is the signature for custom UserWalk hook methods
	UserWalkHook func(context.Context, boil.ContextExecutor, *UserWalk) error

	userWalkQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	userWalkType                 = reflect.TypeOf(&UserWalk{})
	userWalkMapping              = queries.MakeStructMapping(userWalkType)
	userWalkPrimaryKeyMapping, _ = queries.BindMapping(userWalkType, userWalkMapping, userWalkPrimaryKeyColumns)
	userWalkInsertCacheMut       sync.RWMutex
	userWalkInsertCache          = make(map[string]insertCache)
	userWalkUpdateCacheMut       sync.RWMutex
	userWalkUpdateCache          = make(map[string]updateCache)
	userWalkUpsertCacheMut       sync.RWMutex
	userWalkUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var userWalkBeforeInsertHooks []UserWalkHook
var userWalkBeforeUpdateHooks []UserWalkHook
var userWalkBeforeDeleteHooks []UserWalkHook
var userWalkBeforeUpsertHooks []UserWalkHook

var userWalkAfterInsertHooks []UserWalkHook
var userWalkAfterSelectHooks []UserWalkHook
var userWalkAfterUpdateHooks []UserWalkHook
var userWalkAfterDeleteHooks []UserWalkHook
var userWalkAfterUpsertHooks []UserWalkHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *UserWalk) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userWalkBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *UserWalk) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userWalkBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *UserWalk) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userWalkBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *UserWalk) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userWalkBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *UserWalk) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userWalkAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *UserWalk) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userWalkAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *UserWalk) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userWalkAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *UserWalk) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userWalkAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *UserWalk) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userWalkAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddUserWalkHook registers your hook function for all future operations.
func AddUserWalkHook(hookPoint boil.HookPoint, userWalkHook UserWalkHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		userWalkBeforeInsertHooks = append(userWalkBeforeInsertHooks, userWalkHook)
	case boil.BeforeUpdateHook:
		userWalkBeforeUpdateHooks = append(userWalkBeforeUpdateHooks, userWalkHook)
	case boil.BeforeDeleteHook:
		userWalkBeforeDeleteHooks = append(userWalkBeforeDeleteHooks, userWalkHook)
	case boil.BeforeUpsertHook:
		userWalkBeforeUpsertHooks = append(userWalkBeforeUpsertHooks, userWalkHook)
	case boil.AfterInsertHook:
		userWalkAfterInsertHooks = append(userWalkAfterInsertHooks, userWalkHook)
	case boil.AfterSelectHook:
		userWalkAfterSelectHooks = append(userWalkAfterSelectHooks, userWalkHook)
	case boil.AfterUpdateHook:
		userWalkAfterUpdateHooks = append(userWalkAfterUpdateHooks, userWalkHook)
	case boil.AfterDeleteHook:
		userWalkAfterDeleteHooks = append(userWalkAfterDeleteHooks, userWalkHook)
	case boil.AfterUpsertHook:
		userWalkAfterUpsertHooks = append(userWalkAfterUpsertHooks, userWalkHook)
	}
}

// One returns a single userWalk record from the query.
func (q userWalkQuery) One(ctx context.Context, exec boil.ContextExecutor) (*UserWalk, error) {
	o := &UserWalk{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for user_walk")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all UserWalk records from the query.
func (q userWalkQuery) All(ctx context.Context, exec boil.ContextExecutor) (UserWalkSlice, error) {
	var o []*UserWalk

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to UserWalk slice")
	}

	if len(userWalkAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all UserWalk records in the query.
func (q userWalkQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count user_walk rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q userWalkQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if user_walk exists")
	}

	return count > 0, nil
}

// User pointed to by the foreign key.
func (o *UserWalk) User(mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(queryMods...)
	queries.SetFrom(query.Query, "\"user\"")

	return query
}

// Walk pointed to by the foreign key.
func (o *UserWalk) Walk(mods ...qm.QueryMod) walkQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.WalkID),
	}

	queryMods = append(queryMods, mods...)

	query := Walks(queryMods...)
	queries.SetFrom(query.Query, "\"walk\"")

	return query
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (userWalkL) LoadUser(ctx context.Context, e boil.ContextExecutor, singular bool, maybeUserWalk interface{}, mods queries.Applicator) error {
	var slice []*UserWalk
	var object *UserWalk

	if singular {
		object = maybeUserWalk.(*UserWalk)
	} else {
		slice = *maybeUserWalk.(*[]*UserWalk)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &userWalkR{}
		}
		args = append(args, object.UserID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &userWalkR{}
			}

			for _, a := range args {
				if a == obj.UserID {
					continue Outer
				}
			}

			args = append(args, obj.UserID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`user`),
		qm.WhereIn(`user.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load User")
	}

	var resultSlice []*User
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice User")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for user")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for user")
	}

	if len(userWalkAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.User = foreign
		if foreign.R == nil {
			foreign.R = &userR{}
		}
		foreign.R.UserWalks = append(foreign.R.UserWalks, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.UserID == foreign.ID {
				local.R.User = foreign
				if foreign.R == nil {
					foreign.R = &userR{}
				}
				foreign.R.UserWalks = append(foreign.R.UserWalks, local)
				break
			}
		}
	}

	return nil
}

// LoadWalk allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (userWalkL) LoadWalk(ctx context.Context, e boil.ContextExecutor, singular bool, maybeUserWalk interface{}, mods queries.Applicator) error {
	var slice []*UserWalk
	var object *UserWalk

	if singular {
		object = maybeUserWalk.(*UserWalk)
	} else {
		slice = *maybeUserWalk.(*[]*UserWalk)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &userWalkR{}
		}
		args = append(args, object.WalkID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &userWalkR{}
			}

			for _, a := range args {
				if a == obj.WalkID {
					continue Outer
				}
			}

			args = append(args, obj.WalkID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`walk`),
		qm.WhereIn(`walk.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Walk")
	}

	var resultSlice []*Walk
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Walk")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for walk")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for walk")
	}

	if len(userWalkAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Walk = foreign
		if foreign.R == nil {
			foreign.R = &walkR{}
		}
		foreign.R.UserWalks = append(foreign.R.UserWalks, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.WalkID == foreign.ID {
				local.R.Walk = foreign
				if foreign.R == nil {
					foreign.R = &walkR{}
				}
				foreign.R.UserWalks = append(foreign.R.UserWalks, local)
				break
			}
		}
	}

	return nil
}

// SetUser of the userWalk to the related item.
// Sets o.R.User to related.
// Adds o to related.R.UserWalks.
func (o *UserWalk) SetUser(ctx context.Context, exec boil.ContextExecutor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"user_walk\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, userWalkPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.UserID = related.ID
	if o.R == nil {
		o.R = &userWalkR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			UserWalks: UserWalkSlice{o},
		}
	} else {
		related.R.UserWalks = append(related.R.UserWalks, o)
	}

	return nil
}

// SetWalk of the userWalk to the related item.
// Sets o.R.Walk to related.
// Adds o to related.R.UserWalks.
func (o *UserWalk) SetWalk(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Walk) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"user_walk\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"walk_id"}),
		strmangle.WhereClause("\"", "\"", 2, userWalkPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.WalkID = related.ID
	if o.R == nil {
		o.R = &userWalkR{
			Walk: related,
		}
	} else {
		o.R.Walk = related
	}

	if related.R == nil {
		related.R = &walkR{
			UserWalks: UserWalkSlice{o},
		}
	} else {
		related.R.UserWalks = append(related.R.UserWalks, o)
	}

	return nil
}

// UserWalks retrieves all the records using an executor.
func UserWalks(mods ...qm.QueryMod) userWalkQuery {
	mods = append(mods, qm.From("\"user_walk\""))
	return userWalkQuery{NewQuery(mods...)}
}

// FindUserWalk retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUserWalk(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*UserWalk, error) {
	userWalkObj := &UserWalk{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"user_walk\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, userWalkObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from user_walk")
	}

	return userWalkObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *UserWalk) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no user_walk provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userWalkColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	userWalkInsertCacheMut.RLock()
	cache, cached := userWalkInsertCache[key]
	userWalkInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			userWalkAllColumns,
			userWalkColumnsWithDefault,
			userWalkColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(userWalkType, userWalkMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(userWalkType, userWalkMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"user_walk\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"user_walk\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into user_walk")
	}

	if !cached {
		userWalkInsertCacheMut.Lock()
		userWalkInsertCache[key] = cache
		userWalkInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the UserWalk.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *UserWalk) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	userWalkUpdateCacheMut.RLock()
	cache, cached := userWalkUpdateCache[key]
	userWalkUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			userWalkAllColumns,
			userWalkPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update user_walk, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"user_walk\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, userWalkPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(userWalkType, userWalkMapping, append(wl, userWalkPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update user_walk row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for user_walk")
	}

	if !cached {
		userWalkUpdateCacheMut.Lock()
		userWalkUpdateCache[key] = cache
		userWalkUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q userWalkQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for user_walk")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for user_walk")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UserWalkSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userWalkPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"user_walk\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, userWalkPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in userWalk slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all userWalk")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *UserWalk) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no user_walk provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userWalkColumnsWithDefault, o)

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

	userWalkUpsertCacheMut.RLock()
	cache, cached := userWalkUpsertCache[key]
	userWalkUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			userWalkAllColumns,
			userWalkColumnsWithDefault,
			userWalkColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			userWalkAllColumns,
			userWalkPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert user_walk, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(userWalkPrimaryKeyColumns))
			copy(conflict, userWalkPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"user_walk\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(userWalkType, userWalkMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(userWalkType, userWalkMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert user_walk")
	}

	if !cached {
		userWalkUpsertCacheMut.Lock()
		userWalkUpsertCache[key] = cache
		userWalkUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single UserWalk record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *UserWalk) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no UserWalk provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), userWalkPrimaryKeyMapping)
	sql := "DELETE FROM \"user_walk\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from user_walk")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for user_walk")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q userWalkQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no userWalkQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from user_walk")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for user_walk")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UserWalkSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(userWalkBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userWalkPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"user_walk\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userWalkPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from userWalk slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for user_walk")
	}

	if len(userWalkAfterDeleteHooks) != 0 {
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
func (o *UserWalk) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindUserWalk(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserWalkSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := UserWalkSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userWalkPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"user_walk\".* FROM \"user_walk\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userWalkPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in UserWalkSlice")
	}

	*o = slice

	return nil
}

// UserWalkExists checks if the UserWalk row exists.
func UserWalkExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"user_walk\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if user_walk exists")
	}

	return exists, nil
}
