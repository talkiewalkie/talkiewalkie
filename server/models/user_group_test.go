// Code generated by SQLBoiler 4.5.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

func testUserGroups(t *testing.T) {
	t.Parallel()

	query := UserGroups()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testUserGroupsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserGroup{}
	if err = randomize.Struct(seed, o, userGroupDBTypes, true, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
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

	count, err := UserGroups().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUserGroupsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserGroup{}
	if err = randomize.Struct(seed, o, userGroupDBTypes, true, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := UserGroups().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := UserGroups().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUserGroupsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserGroup{}
	if err = randomize.Struct(seed, o, userGroupDBTypes, true, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := UserGroupSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := UserGroups().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUserGroupsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserGroup{}
	if err = randomize.Struct(seed, o, userGroupDBTypes, true, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := UserGroupExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if UserGroup exists: %s", err)
	}
	if !e {
		t.Errorf("Expected UserGroupExists to return true, but got false.")
	}
}

func testUserGroupsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserGroup{}
	if err = randomize.Struct(seed, o, userGroupDBTypes, true, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	userGroupFound, err := FindUserGroup(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if userGroupFound == nil {
		t.Error("want a record, got nil")
	}
}

func testUserGroupsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserGroup{}
	if err = randomize.Struct(seed, o, userGroupDBTypes, true, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = UserGroups().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testUserGroupsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserGroup{}
	if err = randomize.Struct(seed, o, userGroupDBTypes, true, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := UserGroups().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testUserGroupsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	userGroupOne := &UserGroup{}
	userGroupTwo := &UserGroup{}
	if err = randomize.Struct(seed, userGroupOne, userGroupDBTypes, false, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}
	if err = randomize.Struct(seed, userGroupTwo, userGroupDBTypes, false, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = userGroupOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = userGroupTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := UserGroups().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testUserGroupsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	userGroupOne := &UserGroup{}
	userGroupTwo := &UserGroup{}
	if err = randomize.Struct(seed, userGroupOne, userGroupDBTypes, false, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}
	if err = randomize.Struct(seed, userGroupTwo, userGroupDBTypes, false, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = userGroupOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = userGroupTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UserGroups().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func userGroupBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *UserGroup) error {
	*o = UserGroup{}
	return nil
}

func userGroupAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *UserGroup) error {
	*o = UserGroup{}
	return nil
}

func userGroupAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *UserGroup) error {
	*o = UserGroup{}
	return nil
}

func userGroupBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *UserGroup) error {
	*o = UserGroup{}
	return nil
}

func userGroupAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *UserGroup) error {
	*o = UserGroup{}
	return nil
}

func userGroupBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *UserGroup) error {
	*o = UserGroup{}
	return nil
}

func userGroupAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *UserGroup) error {
	*o = UserGroup{}
	return nil
}

func userGroupBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *UserGroup) error {
	*o = UserGroup{}
	return nil
}

func userGroupAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *UserGroup) error {
	*o = UserGroup{}
	return nil
}

func testUserGroupsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &UserGroup{}
	o := &UserGroup{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, userGroupDBTypes, false); err != nil {
		t.Errorf("Unable to randomize UserGroup object: %s", err)
	}

	AddUserGroupHook(boil.BeforeInsertHook, userGroupBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	userGroupBeforeInsertHooks = []UserGroupHook{}

	AddUserGroupHook(boil.AfterInsertHook, userGroupAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	userGroupAfterInsertHooks = []UserGroupHook{}

	AddUserGroupHook(boil.AfterSelectHook, userGroupAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	userGroupAfterSelectHooks = []UserGroupHook{}

	AddUserGroupHook(boil.BeforeUpdateHook, userGroupBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	userGroupBeforeUpdateHooks = []UserGroupHook{}

	AddUserGroupHook(boil.AfterUpdateHook, userGroupAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	userGroupAfterUpdateHooks = []UserGroupHook{}

	AddUserGroupHook(boil.BeforeDeleteHook, userGroupBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	userGroupBeforeDeleteHooks = []UserGroupHook{}

	AddUserGroupHook(boil.AfterDeleteHook, userGroupAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	userGroupAfterDeleteHooks = []UserGroupHook{}

	AddUserGroupHook(boil.BeforeUpsertHook, userGroupBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	userGroupBeforeUpsertHooks = []UserGroupHook{}

	AddUserGroupHook(boil.AfterUpsertHook, userGroupAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	userGroupAfterUpsertHooks = []UserGroupHook{}
}

func testUserGroupsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserGroup{}
	if err = randomize.Struct(seed, o, userGroupDBTypes, true, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UserGroups().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testUserGroupsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserGroup{}
	if err = randomize.Struct(seed, o, userGroupDBTypes, true); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(userGroupColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := UserGroups().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testUserGroupToOneGroupUsingGroup(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local UserGroup
	var foreign Group

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, userGroupDBTypes, false, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, groupDBTypes, false, groupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Group struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.GroupID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Group().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := UserGroupSlice{&local}
	if err = local.L.LoadGroup(ctx, tx, false, (*[]*UserGroup)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Group == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Group = nil
	if err = local.L.LoadGroup(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Group == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testUserGroupToOneUserUsingUser(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local UserGroup
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, userGroupDBTypes, false, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
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

	slice := UserGroupSlice{&local}
	if err = local.L.LoadUser(ctx, tx, false, (*[]*UserGroup)(&slice), nil); err != nil {
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

func testUserGroupToOneSetOpGroupUsingGroup(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a UserGroup
	var b, c Group

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userGroupDBTypes, false, strmangle.SetComplement(userGroupPrimaryKeyColumns, userGroupColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, groupDBTypes, false, strmangle.SetComplement(groupPrimaryKeyColumns, groupColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, groupDBTypes, false, strmangle.SetComplement(groupPrimaryKeyColumns, groupColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Group{&b, &c} {
		err = a.SetGroup(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Group != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.UserGroups[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.GroupID != x.ID {
			t.Error("foreign key was wrong value", a.GroupID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.GroupID))
		reflect.Indirect(reflect.ValueOf(&a.GroupID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.GroupID != x.ID {
			t.Error("foreign key was wrong value", a.GroupID, x.ID)
		}
	}
}
func testUserGroupToOneSetOpUserUsingUser(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a UserGroup
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, userGroupDBTypes, false, strmangle.SetComplement(userGroupPrimaryKeyColumns, userGroupColumnsWithoutDefault)...); err != nil {
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

		if x.R.UserGroups[0] != &a {
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

func testUserGroupsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserGroup{}
	if err = randomize.Struct(seed, o, userGroupDBTypes, true, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
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

func testUserGroupsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserGroup{}
	if err = randomize.Struct(seed, o, userGroupDBTypes, true, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := UserGroupSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testUserGroupsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserGroup{}
	if err = randomize.Struct(seed, o, userGroupDBTypes, true, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := UserGroups().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	userGroupDBTypes = map[string]string{`ID`: `integer`, `UserID`: `integer`, `GroupID`: `integer`, `CreatedAt`: `timestamp with time zone`}
	_                = bytes.MinRead
)

func testUserGroupsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(userGroupPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(userGroupAllColumns) == len(userGroupPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &UserGroup{}
	if err = randomize.Struct(seed, o, userGroupDBTypes, true, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UserGroups().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, userGroupDBTypes, true, userGroupPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testUserGroupsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(userGroupAllColumns) == len(userGroupPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &UserGroup{}
	if err = randomize.Struct(seed, o, userGroupDBTypes, true, userGroupColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UserGroups().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, userGroupDBTypes, true, userGroupPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(userGroupAllColumns, userGroupPrimaryKeyColumns) {
		fields = userGroupAllColumns
	} else {
		fields = strmangle.SetComplement(
			userGroupAllColumns,
			userGroupPrimaryKeyColumns,
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

	slice := UserGroupSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testUserGroupsUpsert(t *testing.T) {
	t.Parallel()

	if len(userGroupAllColumns) == len(userGroupPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := UserGroup{}
	if err = randomize.Struct(seed, &o, userGroupDBTypes, true); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert UserGroup: %s", err)
	}

	count, err := UserGroups().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, userGroupDBTypes, false, userGroupPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize UserGroup struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert UserGroup: %s", err)
	}

	count, err = UserGroups().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}