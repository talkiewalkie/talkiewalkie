// Code generated by SQLBoiler 4.6.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testUsers(t *testing.T) {
	t.Parallel()

	query := Users()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testUsersDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &User{}
	if err = randomize.Struct(seed, o, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Users().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUsersQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &User{}
	if err = randomize.Struct(seed, o, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Users().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Users().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUsersSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &User{}
	if err = randomize.Struct(seed, o, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := UserSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Users().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUsersExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &User{}
	if err = randomize.Struct(seed, o, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := UserExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if User exists: %s", err)
	}
	if !e {
		t.Errorf("Expected UserExists to return true, but got false.")
	}
}

func testUsersFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &User{}
	if err = randomize.Struct(seed, o, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	userFound, err := FindUser(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if userFound == nil {
		t.Error("want a record, got nil")
	}
}

func testUsersBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &User{}
	if err = randomize.Struct(seed, o, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Users().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testUsersOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &User{}
	if err = randomize.Struct(seed, o, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Users().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testUsersAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	userOne := &User{}
	userTwo := &User{}
	if err = randomize.Struct(seed, userOne, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}
	if err = randomize.Struct(seed, userTwo, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = userOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = userTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Users().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testUsersCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	userOne := &User{}
	userTwo := &User{}
	if err = randomize.Struct(seed, userOne, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}
	if err = randomize.Struct(seed, userTwo, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = userOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = userTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Users().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func userBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *User) error {
	*o = User{}
	return nil
}

func userAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *User) error {
	*o = User{}
	return nil
}

func userAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *User) error {
	*o = User{}
	return nil
}

func userBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *User) error {
	*o = User{}
	return nil
}

func userAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *User) error {
	*o = User{}
	return nil
}

func userBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *User) error {
	*o = User{}
	return nil
}

func userAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *User) error {
	*o = User{}
	return nil
}

func userBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *User) error {
	*o = User{}
	return nil
}

func userAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *User) error {
	*o = User{}
	return nil
}

func testUsersHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &User{}
	o := &User{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, userDBTypes, false); err != nil {
		t.Errorf("Unable to randomize User object: %s", err)
	}

	AddUserHook(boil.BeforeInsertHook, userBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	userBeforeInsertHooks = []UserHook{}

	AddUserHook(boil.AfterInsertHook, userAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	userAfterInsertHooks = []UserHook{}

	AddUserHook(boil.AfterSelectHook, userAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	userAfterSelectHooks = []UserHook{}

	AddUserHook(boil.BeforeUpdateHook, userBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	userBeforeUpdateHooks = []UserHook{}

	AddUserHook(boil.AfterUpdateHook, userAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	userAfterUpdateHooks = []UserHook{}

	AddUserHook(boil.BeforeDeleteHook, userBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	userBeforeDeleteHooks = []UserHook{}

	AddUserHook(boil.AfterDeleteHook, userAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	userAfterDeleteHooks = []UserHook{}

	AddUserHook(boil.BeforeUpsertHook, userBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	userBeforeUpsertHooks = []UserHook{}

	AddUserHook(boil.AfterUpsertHook, userAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	userAfterUpsertHooks = []UserHook{}
}

func testUsersInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &User{}
	if err = randomize.Struct(seed, o, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Users().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testUsersInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &User{}
	if err = randomize.Struct(seed, o, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(userColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Users().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testUserToManyAuthorMessages(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a User
	var b, c Message

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, messageDBTypes, false, messageColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, messageDBTypes, false, messageColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&b.AuthorID, a.ID)
	queries.Assign(&c.AuthorID, a.ID)
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.AuthorMessages().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if queries.Equal(v.AuthorID, b.AuthorID) {
			bFound = true
		}
		if queries.Equal(v.AuthorID, c.AuthorID) {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := UserSlice{&a}
	if err = a.L.LoadAuthorMessages(ctx, tx, false, (*[]*User)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.AuthorMessages); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.AuthorMessages = nil
	if err = a.L.LoadAuthorMessages(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.AuthorMessages); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testUserToManyUserConversations(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a User
	var b, c UserConversation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, userConversationDBTypes, false, userConversationColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, userConversationDBTypes, false, userConversationColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.UserID = a.ID
	c.UserID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.UserConversations().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.UserID == b.UserID {
			bFound = true
		}
		if v.UserID == c.UserID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := UserSlice{&a}
	if err = a.L.LoadUserConversations(ctx, tx, false, (*[]*User)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.UserConversations); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.UserConversations = nil
	if err = a.L.LoadUserConversations(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.UserConversations); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testUserToManyAddOpAuthorMessages(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a User
	var b, c, d, e Message

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Message{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, messageDBTypes, false, strmangle.SetComplement(messagePrimaryKeyColumns, messageColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*Message{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddAuthorMessages(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if !queries.Equal(a.ID, first.AuthorID) {
			t.Error("foreign key was wrong value", a.ID, first.AuthorID)
		}
		if !queries.Equal(a.ID, second.AuthorID) {
			t.Error("foreign key was wrong value", a.ID, second.AuthorID)
		}

		if first.R.Author != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Author != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.AuthorMessages[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.AuthorMessages[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.AuthorMessages().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testUserToManySetOpAuthorMessages(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a User
	var b, c, d, e Message

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Message{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, messageDBTypes, false, strmangle.SetComplement(messagePrimaryKeyColumns, messageColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err = a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.SetAuthorMessages(ctx, tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.AuthorMessages().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetAuthorMessages(ctx, tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.AuthorMessages().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if !queries.IsValuerNil(b.AuthorID) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.AuthorID) {
		t.Error("want c's foreign key value to be nil")
	}
	if !queries.Equal(a.ID, d.AuthorID) {
		t.Error("foreign key was wrong value", a.ID, d.AuthorID)
	}
	if !queries.Equal(a.ID, e.AuthorID) {
		t.Error("foreign key was wrong value", a.ID, e.AuthorID)
	}

	if b.R.Author != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Author != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Author != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Author != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.AuthorMessages[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.AuthorMessages[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testUserToManyRemoveOpAuthorMessages(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a User
	var b, c, d, e Message

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Message{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, messageDBTypes, false, strmangle.SetComplement(messagePrimaryKeyColumns, messageColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.AddAuthorMessages(ctx, tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.AuthorMessages().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveAuthorMessages(ctx, tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.AuthorMessages().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if !queries.IsValuerNil(b.AuthorID) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.AuthorID) {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.Author != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Author != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Author != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.Author != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.AuthorMessages) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.AuthorMessages[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.AuthorMessages[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testUserToManyAddOpUserConversations(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a User
	var b, c, d, e UserConversation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*UserConversation{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, userConversationDBTypes, false, strmangle.SetComplement(userConversationPrimaryKeyColumns, userConversationColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*UserConversation{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddUserConversations(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.UserID {
			t.Error("foreign key was wrong value", a.ID, first.UserID)
		}
		if a.ID != second.UserID {
			t.Error("foreign key was wrong value", a.ID, second.UserID)
		}

		if first.R.User != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.User != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.UserConversations[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.UserConversations[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.UserConversations().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testUserToOneAssetUsingProfilePictureAsset(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local User
	var foreign Asset

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, assetDBTypes, false, assetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Asset struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&local.ProfilePicture, foreign.ID)
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.ProfilePictureAsset().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if !queries.Equal(check.ID, foreign.ID) {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := UserSlice{&local}
	if err = local.L.LoadProfilePictureAsset(ctx, tx, false, (*[]*User)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.ProfilePictureAsset == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.ProfilePictureAsset = nil
	if err = local.L.LoadProfilePictureAsset(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.ProfilePictureAsset == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testUserToOneSetOpAssetUsingProfilePictureAsset(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a User
	var b, c Asset

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, assetDBTypes, false, strmangle.SetComplement(assetPrimaryKeyColumns, assetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, assetDBTypes, false, strmangle.SetComplement(assetPrimaryKeyColumns, assetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Asset{&b, &c} {
		err = a.SetProfilePictureAsset(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.ProfilePictureAsset != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.ProfilePictureUsers[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if !queries.Equal(a.ProfilePicture, x.ID) {
			t.Error("foreign key was wrong value", a.ProfilePicture)
		}

		zero := reflect.Zero(reflect.TypeOf(a.ProfilePicture))
		reflect.Indirect(reflect.ValueOf(&a.ProfilePicture)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if !queries.Equal(a.ProfilePicture, x.ID) {
			t.Error("foreign key was wrong value", a.ProfilePicture, x.ID)
		}
	}
}

func testUserToOneRemoveOpAssetUsingProfilePictureAsset(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a User
	var b Asset

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, assetDBTypes, false, strmangle.SetComplement(assetPrimaryKeyColumns, assetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = a.SetProfilePictureAsset(ctx, tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveProfilePictureAsset(ctx, tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.ProfilePictureAsset().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.ProfilePictureAsset != nil {
		t.Error("R struct entry should be nil")
	}

	if !queries.IsValuerNil(a.ProfilePicture) {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.ProfilePictureUsers) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testUsersReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &User{}
	if err = randomize.Struct(seed, o, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testUsersReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &User{}
	if err = randomize.Struct(seed, o, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := UserSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testUsersSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &User{}
	if err = randomize.Struct(seed, o, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Users().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	userDBTypes = map[string]string{`ID`: `integer`, `UUID`: `uuid`, `FirebaseUID`: `character varying`, `ProfilePicture`: `integer`, `CreatedAt`: `timestamp with time zone`, `UpdatedAt`: `timestamp with time zone`, `Bio`: `text`, `PhoneNumber`: `character varying`, `OnboardingFinished`: `boolean`, `DisplayName`: `character varying`, `Locales`: `ARRAYcharacter varying`}
	_           = bytes.MinRead
)

func testUsersUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(userPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(userAllColumns) == len(userPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &User{}
	if err = randomize.Struct(seed, o, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Users().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, userDBTypes, true, userPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testUsersSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(userAllColumns) == len(userPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &User{}
	if err = randomize.Struct(seed, o, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Users().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, userDBTypes, true, userPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(userAllColumns, userPrimaryKeyColumns) {
		fields = userAllColumns
	} else {
		fields = strmangle.SetComplement(
			userAllColumns,
			userPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := UserSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testUsersUpsert(t *testing.T) {
	t.Parallel()

	if len(userAllColumns) == len(userPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := User{}
	if err = randomize.Struct(seed, &o, userDBTypes, true); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert User: %s", err)
	}

	count, err := Users().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, userDBTypes, false, userPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert User: %s", err)
	}

	count, err = Users().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
