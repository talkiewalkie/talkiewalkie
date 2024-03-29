// Code generated by SQLBoiler 4.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

func testUserConversations(t *testing.T) {
	t.Parallel()

	query := UserConversations()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testUserConversationsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserConversation{}
	if err = randomize.Struct(seed, o, userConversationDBTypes, true, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
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

	count, err := UserConversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUserConversationsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserConversation{}
	if err = randomize.Struct(seed, o, userConversationDBTypes, true, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := UserConversations().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := UserConversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUserConversationsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserConversation{}
	if err = randomize.Struct(seed, o, userConversationDBTypes, true, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := UserConversationSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := UserConversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUserConversationsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserConversation{}
	if err = randomize.Struct(seed, o, userConversationDBTypes, true, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := UserConversationExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if UserConversation exists: %s", err)
	}
	if !e {
		t.Errorf("Expected UserConversationExists to return true, but got false.")
	}
}

func testUserConversationsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserConversation{}
	if err = randomize.Struct(seed, o, userConversationDBTypes, true, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	userConversationFound, err := FindUserConversation(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if userConversationFound == nil {
		t.Error("want a record, got nil")
	}
}

func testUserConversationsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserConversation{}
	if err = randomize.Struct(seed, o, userConversationDBTypes, true, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = UserConversations().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testUserConversationsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserConversation{}
	if err = randomize.Struct(seed, o, userConversationDBTypes, true, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := UserConversations().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testUserConversationsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	userConversationOne := &UserConversation{}
	userConversationTwo := &UserConversation{}
	if err = randomize.Struct(seed, userConversationOne, userConversationDBTypes, false, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}
	if err = randomize.Struct(seed, userConversationTwo, userConversationDBTypes, false, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = userConversationOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = userConversationTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := UserConversations().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testUserConversationsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	userConversationOne := &UserConversation{}
	userConversationTwo := &UserConversation{}
	if err = randomize.Struct(seed, userConversationOne, userConversationDBTypes, false, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}
	if err = randomize.Struct(seed, userConversationTwo, userConversationDBTypes, false, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = userConversationOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = userConversationTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UserConversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func userConversationBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *UserConversation) error {
	*o = UserConversation{}
	return nil
}

func userConversationAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *UserConversation) error {
	*o = UserConversation{}
	return nil
}

func userConversationAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *UserConversation) error {
	*o = UserConversation{}
	return nil
}

func userConversationBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *UserConversation) error {
	*o = UserConversation{}
	return nil
}

func userConversationAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *UserConversation) error {
	*o = UserConversation{}
	return nil
}

func userConversationBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *UserConversation) error {
	*o = UserConversation{}
	return nil
}

func userConversationAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *UserConversation) error {
	*o = UserConversation{}
	return nil
}

func userConversationBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *UserConversation) error {
	*o = UserConversation{}
	return nil
}

func userConversationAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *UserConversation) error {
	*o = UserConversation{}
	return nil
}

func testUserConversationsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &UserConversation{}
	o := &UserConversation{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, userConversationDBTypes, false); err != nil {
		t.Errorf("Unable to randomize UserConversation object: %s", err)
	}

	AddUserConversationHook(boil.BeforeInsertHook, userConversationBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	userConversationBeforeInsertHooks = []UserConversationHook{}

	AddUserConversationHook(boil.AfterInsertHook, userConversationAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	userConversationAfterInsertHooks = []UserConversationHook{}

	AddUserConversationHook(boil.AfterSelectHook, userConversationAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	userConversationAfterSelectHooks = []UserConversationHook{}

	AddUserConversationHook(boil.BeforeUpdateHook, userConversationBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	userConversationBeforeUpdateHooks = []UserConversationHook{}

	AddUserConversationHook(boil.AfterUpdateHook, userConversationAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	userConversationAfterUpdateHooks = []UserConversationHook{}

	AddUserConversationHook(boil.BeforeDeleteHook, userConversationBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	userConversationBeforeDeleteHooks = []UserConversationHook{}

	AddUserConversationHook(boil.AfterDeleteHook, userConversationAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	userConversationAfterDeleteHooks = []UserConversationHook{}

	AddUserConversationHook(boil.BeforeUpsertHook, userConversationBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	userConversationBeforeUpsertHooks = []UserConversationHook{}

	AddUserConversationHook(boil.AfterUpsertHook, userConversationAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	userConversationAfterUpsertHooks = []UserConversationHook{}
}

func testUserConversationsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserConversation{}
	if err = randomize.Struct(seed, o, userConversationDBTypes, true, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UserConversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testUserConversationsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserConversation{}
	if err = randomize.Struct(seed, o, userConversationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(userConversationColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := UserConversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testUserConversationToOneConversationUsingConversation(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local UserConversation
	var foreign Conversation

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, userConversationDBTypes, false, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, conversationDBTypes, false, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.ConversationID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Conversation().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := UserConversationSlice{&local}
	if err = local.L.LoadConversation(ctx, tx, false, (*[]*UserConversation)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Conversation == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Conversation = nil
	if err = local.L.LoadConversation(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Conversation == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testUserConversationToOneUserUsingUser(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local UserConversation
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, userConversationDBTypes, false, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.UserID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.User().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := UserConversationSlice{&local}
	if err = local.L.LoadUser(ctx, tx, false, (*[]*UserConversation)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.User == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.User = nil
	if err = local.L.LoadUser(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.User == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testUserConversationToOneSetOpConversationUsingConversation(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a UserConversation
	var b, c Conversation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userConversationDBTypes, false, strmangle.SetComplement(userConversationPrimaryKeyColumns, userConversationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, conversationDBTypes, false, strmangle.SetComplement(conversationPrimaryKeyColumns, conversationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, conversationDBTypes, false, strmangle.SetComplement(conversationPrimaryKeyColumns, conversationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Conversation{&b, &c} {
		err = a.SetConversation(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Conversation != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.UserConversations[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.ConversationID != x.ID {
			t.Error("foreign key was wrong value", a.ConversationID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.ConversationID))
		reflect.Indirect(reflect.ValueOf(&a.ConversationID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.ConversationID != x.ID {
			t.Error("foreign key was wrong value", a.ConversationID, x.ID)
		}
	}
}
func testUserConversationToOneSetOpUserUsingUser(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a UserConversation
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userConversationDBTypes, false, strmangle.SetComplement(userConversationPrimaryKeyColumns, userConversationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*User{&b, &c} {
		err = a.SetUser(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.User != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.UserConversations[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.UserID != x.ID {
			t.Error("foreign key was wrong value", a.UserID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.UserID))
		reflect.Indirect(reflect.ValueOf(&a.UserID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.UserID != x.ID {
			t.Error("foreign key was wrong value", a.UserID, x.ID)
		}
	}
}

func testUserConversationsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserConversation{}
	if err = randomize.Struct(seed, o, userConversationDBTypes, true, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
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

func testUserConversationsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserConversation{}
	if err = randomize.Struct(seed, o, userConversationDBTypes, true, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := UserConversationSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testUserConversationsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserConversation{}
	if err = randomize.Struct(seed, o, userConversationDBTypes, true, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := UserConversations().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	userConversationDBTypes = map[string]string{`ID`: `integer`, `UserID`: `integer`, `ConversationID`: `integer`, `CreatedAt`: `timestamp with time zone`, `ReadUntil`: `timestamp with time zone`}
	_                       = bytes.MinRead
)

func testUserConversationsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(userConversationPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(userConversationAllColumns) == len(userConversationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &UserConversation{}
	if err = randomize.Struct(seed, o, userConversationDBTypes, true, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UserConversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, userConversationDBTypes, true, userConversationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testUserConversationsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(userConversationAllColumns) == len(userConversationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &UserConversation{}
	if err = randomize.Struct(seed, o, userConversationDBTypes, true, userConversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UserConversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, userConversationDBTypes, true, userConversationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(userConversationAllColumns, userConversationPrimaryKeyColumns) {
		fields = userConversationAllColumns
	} else {
		fields = strmangle.SetComplement(
			userConversationAllColumns,
			userConversationPrimaryKeyColumns,
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

	slice := UserConversationSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testUserConversationsUpsert(t *testing.T) {
	t.Parallel()

	if len(userConversationAllColumns) == len(userConversationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := UserConversation{}
	if err = randomize.Struct(seed, &o, userConversationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert UserConversation: %s", err)
	}

	count, err := UserConversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, userConversationDBTypes, false, userConversationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize UserConversation struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert UserConversation: %s", err)
	}

	count, err = UserConversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
