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

func testConversationsUpsert(t *testing.T) {
	t.Parallel()

	if len(conversationAllColumns) == len(conversationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Conversation{}
	if err = randomize.Struct(seed, &o, conversationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Conversation: %s", err)
	}

	count, err := Conversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, conversationDBTypes, false, conversationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Conversation: %s", err)
	}

	count, err = Conversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testConversations(t *testing.T) {
	t.Parallel()

	query := Conversations()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testConversationsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Conversation{}
	if err = randomize.Struct(seed, o, conversationDBTypes, true, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
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

	count, err := Conversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testConversationsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Conversation{}
	if err = randomize.Struct(seed, o, conversationDBTypes, true, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Conversations().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Conversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testConversationsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Conversation{}
	if err = randomize.Struct(seed, o, conversationDBTypes, true, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ConversationSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Conversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testConversationsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Conversation{}
	if err = randomize.Struct(seed, o, conversationDBTypes, true, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := ConversationExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if Conversation exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ConversationExists to return true, but got false.")
	}
}

func testConversationsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Conversation{}
	if err = randomize.Struct(seed, o, conversationDBTypes, true, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	conversationFound, err := FindConversation(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if conversationFound == nil {
		t.Error("want a record, got nil")
	}
}

func testConversationsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Conversation{}
	if err = randomize.Struct(seed, o, conversationDBTypes, true, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Conversations().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testConversationsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Conversation{}
	if err = randomize.Struct(seed, o, conversationDBTypes, true, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Conversations().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testConversationsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	conversationOne := &Conversation{}
	conversationTwo := &Conversation{}
	if err = randomize.Struct(seed, conversationOne, conversationDBTypes, false, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}
	if err = randomize.Struct(seed, conversationTwo, conversationDBTypes, false, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = conversationOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = conversationTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Conversations().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testConversationsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	conversationOne := &Conversation{}
	conversationTwo := &Conversation{}
	if err = randomize.Struct(seed, conversationOne, conversationDBTypes, false, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}
	if err = randomize.Struct(seed, conversationTwo, conversationDBTypes, false, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = conversationOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = conversationTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Conversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func conversationBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *Conversation) error {
	*o = Conversation{}
	return nil
}

func conversationAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *Conversation) error {
	*o = Conversation{}
	return nil
}

func conversationAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *Conversation) error {
	*o = Conversation{}
	return nil
}

func conversationBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Conversation) error {
	*o = Conversation{}
	return nil
}

func conversationAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Conversation) error {
	*o = Conversation{}
	return nil
}

func conversationBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Conversation) error {
	*o = Conversation{}
	return nil
}

func conversationAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Conversation) error {
	*o = Conversation{}
	return nil
}

func conversationBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Conversation) error {
	*o = Conversation{}
	return nil
}

func conversationAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Conversation) error {
	*o = Conversation{}
	return nil
}

func testConversationsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &Conversation{}
	o := &Conversation{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, conversationDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Conversation object: %s", err)
	}

	AddConversationHook(boil.BeforeInsertHook, conversationBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	conversationBeforeInsertHooks = []ConversationHook{}

	AddConversationHook(boil.AfterInsertHook, conversationAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	conversationAfterInsertHooks = []ConversationHook{}

	AddConversationHook(boil.AfterSelectHook, conversationAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	conversationAfterSelectHooks = []ConversationHook{}

	AddConversationHook(boil.BeforeUpdateHook, conversationBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	conversationBeforeUpdateHooks = []ConversationHook{}

	AddConversationHook(boil.AfterUpdateHook, conversationAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	conversationAfterUpdateHooks = []ConversationHook{}

	AddConversationHook(boil.BeforeDeleteHook, conversationBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	conversationBeforeDeleteHooks = []ConversationHook{}

	AddConversationHook(boil.AfterDeleteHook, conversationAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	conversationAfterDeleteHooks = []ConversationHook{}

	AddConversationHook(boil.BeforeUpsertHook, conversationBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	conversationBeforeUpsertHooks = []ConversationHook{}

	AddConversationHook(boil.AfterUpsertHook, conversationAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	conversationAfterUpsertHooks = []ConversationHook{}
}

func testConversationsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Conversation{}
	if err = randomize.Struct(seed, o, conversationDBTypes, true, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Conversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testConversationsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Conversation{}
	if err = randomize.Struct(seed, o, conversationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(conversationColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Conversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testConversationToManyMessages(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Conversation
	var b, c Message

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, conversationDBTypes, true, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
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

	b.ConversationID = a.ID
	c.ConversationID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.Messages().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.ConversationID == b.ConversationID {
			bFound = true
		}
		if v.ConversationID == c.ConversationID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := ConversationSlice{&a}
	if err = a.L.LoadMessages(ctx, tx, false, (*[]*Conversation)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Messages); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Messages = nil
	if err = a.L.LoadMessages(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Messages); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testConversationToManyUserConversations(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Conversation
	var b, c UserConversation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, conversationDBTypes, true, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
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

	b.ConversationID = a.ID
	c.ConversationID = a.ID

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
		if v.ConversationID == b.ConversationID {
			bFound = true
		}
		if v.ConversationID == c.ConversationID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := ConversationSlice{&a}
	if err = a.L.LoadUserConversations(ctx, tx, false, (*[]*Conversation)(&slice), nil); err != nil {
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

func testConversationToManyAddOpMessages(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Conversation
	var b, c, d, e Message

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, conversationDBTypes, false, strmangle.SetComplement(conversationPrimaryKeyColumns, conversationColumnsWithoutDefault)...); err != nil {
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
		err = a.AddMessages(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.ConversationID {
			t.Error("foreign key was wrong value", a.ID, first.ConversationID)
		}
		if a.ID != second.ConversationID {
			t.Error("foreign key was wrong value", a.ID, second.ConversationID)
		}

		if first.R.Conversation != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Conversation != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Messages[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Messages[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Messages().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testConversationToManyAddOpUserConversations(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Conversation
	var b, c, d, e UserConversation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, conversationDBTypes, false, strmangle.SetComplement(conversationPrimaryKeyColumns, conversationColumnsWithoutDefault)...); err != nil {
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

		if a.ID != first.ConversationID {
			t.Error("foreign key was wrong value", a.ID, first.ConversationID)
		}
		if a.ID != second.ConversationID {
			t.Error("foreign key was wrong value", a.ID, second.ConversationID)
		}

		if first.R.Conversation != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Conversation != &a {
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

func testConversationsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Conversation{}
	if err = randomize.Struct(seed, o, conversationDBTypes, true, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
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

func testConversationsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Conversation{}
	if err = randomize.Struct(seed, o, conversationDBTypes, true, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ConversationSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testConversationsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Conversation{}
	if err = randomize.Struct(seed, o, conversationDBTypes, true, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Conversations().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	conversationDBTypes = map[string]string{`ID`: `integer`, `UUID`: `uuid`, `Name`: `character varying`, `CreatedAt`: `timestamp with time zone`}
	_                   = bytes.MinRead
)

func testConversationsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(conversationPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(conversationAllColumns) == len(conversationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Conversation{}
	if err = randomize.Struct(seed, o, conversationDBTypes, true, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Conversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, conversationDBTypes, true, conversationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testConversationsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(conversationAllColumns) == len(conversationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Conversation{}
	if err = randomize.Struct(seed, o, conversationDBTypes, true, conversationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Conversations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, conversationDBTypes, true, conversationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Conversation struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(conversationAllColumns, conversationPrimaryKeyColumns) {
		fields = conversationAllColumns
	} else {
		fields = strmangle.SetComplement(
			conversationAllColumns,
			conversationPrimaryKeyColumns,
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

	slice := ConversationSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}
