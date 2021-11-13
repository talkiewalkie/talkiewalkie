// Code generated by SQLBoiler 4.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("DieselSchemaMigrations", testDieselSchemaMigrations)
	t.Run("Assets", testAssets)
	t.Run("Conversations", testConversations)
	t.Run("Messages", testMessages)
	t.Run("Users", testUsers)
	t.Run("UserConversations", testUserConversations)
}

func TestDelete(t *testing.T) {
	t.Run("DieselSchemaMigrations", testDieselSchemaMigrationsDelete)
	t.Run("Assets", testAssetsDelete)
	t.Run("Conversations", testConversationsDelete)
	t.Run("Messages", testMessagesDelete)
	t.Run("Users", testUsersDelete)
	t.Run("UserConversations", testUserConversationsDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("DieselSchemaMigrations", testDieselSchemaMigrationsQueryDeleteAll)
	t.Run("Assets", testAssetsQueryDeleteAll)
	t.Run("Conversations", testConversationsQueryDeleteAll)
	t.Run("Messages", testMessagesQueryDeleteAll)
	t.Run("Users", testUsersQueryDeleteAll)
	t.Run("UserConversations", testUserConversationsQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("DieselSchemaMigrations", testDieselSchemaMigrationsSliceDeleteAll)
	t.Run("Assets", testAssetsSliceDeleteAll)
	t.Run("Conversations", testConversationsSliceDeleteAll)
	t.Run("Messages", testMessagesSliceDeleteAll)
	t.Run("Users", testUsersSliceDeleteAll)
	t.Run("UserConversations", testUserConversationsSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("DieselSchemaMigrations", testDieselSchemaMigrationsExists)
	t.Run("Assets", testAssetsExists)
	t.Run("Conversations", testConversationsExists)
	t.Run("Messages", testMessagesExists)
	t.Run("Users", testUsersExists)
	t.Run("UserConversations", testUserConversationsExists)
}

func TestFind(t *testing.T) {
	t.Run("DieselSchemaMigrations", testDieselSchemaMigrationsFind)
	t.Run("Assets", testAssetsFind)
	t.Run("Conversations", testConversationsFind)
	t.Run("Messages", testMessagesFind)
	t.Run("Users", testUsersFind)
	t.Run("UserConversations", testUserConversationsFind)
}

func TestBind(t *testing.T) {
	t.Run("DieselSchemaMigrations", testDieselSchemaMigrationsBind)
	t.Run("Assets", testAssetsBind)
	t.Run("Conversations", testConversationsBind)
	t.Run("Messages", testMessagesBind)
	t.Run("Users", testUsersBind)
	t.Run("UserConversations", testUserConversationsBind)
}

func TestOne(t *testing.T) {
	t.Run("DieselSchemaMigrations", testDieselSchemaMigrationsOne)
	t.Run("Assets", testAssetsOne)
	t.Run("Conversations", testConversationsOne)
	t.Run("Messages", testMessagesOne)
	t.Run("Users", testUsersOne)
	t.Run("UserConversations", testUserConversationsOne)
}

func TestAll(t *testing.T) {
	t.Run("DieselSchemaMigrations", testDieselSchemaMigrationsAll)
	t.Run("Assets", testAssetsAll)
	t.Run("Conversations", testConversationsAll)
	t.Run("Messages", testMessagesAll)
	t.Run("Users", testUsersAll)
	t.Run("UserConversations", testUserConversationsAll)
}

func TestCount(t *testing.T) {
	t.Run("DieselSchemaMigrations", testDieselSchemaMigrationsCount)
	t.Run("Assets", testAssetsCount)
	t.Run("Conversations", testConversationsCount)
	t.Run("Messages", testMessagesCount)
	t.Run("Users", testUsersCount)
	t.Run("UserConversations", testUserConversationsCount)
}

func TestHooks(t *testing.T) {
	t.Run("DieselSchemaMigrations", testDieselSchemaMigrationsHooks)
	t.Run("Assets", testAssetsHooks)
	t.Run("Conversations", testConversationsHooks)
	t.Run("Messages", testMessagesHooks)
	t.Run("Users", testUsersHooks)
	t.Run("UserConversations", testUserConversationsHooks)
}

func TestInsert(t *testing.T) {
	t.Run("DieselSchemaMigrations", testDieselSchemaMigrationsInsert)
	t.Run("DieselSchemaMigrations", testDieselSchemaMigrationsInsertWhitelist)
	t.Run("Assets", testAssetsInsert)
	t.Run("Assets", testAssetsInsertWhitelist)
	t.Run("Conversations", testConversationsInsert)
	t.Run("Conversations", testConversationsInsertWhitelist)
	t.Run("Messages", testMessagesInsert)
	t.Run("Messages", testMessagesInsertWhitelist)
	t.Run("Users", testUsersInsert)
	t.Run("Users", testUsersInsertWhitelist)
	t.Run("UserConversations", testUserConversationsInsert)
	t.Run("UserConversations", testUserConversationsInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("MessageToUserUsingAuthor", testMessageToOneUserUsingAuthor)
	t.Run("MessageToConversationUsingConversation", testMessageToOneConversationUsingConversation)
	t.Run("MessageToAssetUsingRawAudio", testMessageToOneAssetUsingRawAudio)
	t.Run("UserToAssetUsingProfilePictureAsset", testUserToOneAssetUsingProfilePictureAsset)
	t.Run("UserConversationToConversationUsingConversation", testUserConversationToOneConversationUsingConversation)
	t.Run("UserConversationToUserUsingUser", testUserConversationToOneUserUsingUser)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("AssetToRawAudioMessages", testAssetToManyRawAudioMessages)
	t.Run("AssetToProfilePictureUsers", testAssetToManyProfilePictureUsers)
	t.Run("ConversationToMessages", testConversationToManyMessages)
	t.Run("ConversationToUserConversations", testConversationToManyUserConversations)
	t.Run("UserToAuthorMessages", testUserToManyAuthorMessages)
	t.Run("UserToUserConversations", testUserToManyUserConversations)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("MessageToUserUsingAuthorMessages", testMessageToOneSetOpUserUsingAuthor)
	t.Run("MessageToConversationUsingMessages", testMessageToOneSetOpConversationUsingConversation)
	t.Run("MessageToAssetUsingRawAudioMessages", testMessageToOneSetOpAssetUsingRawAudio)
	t.Run("UserToAssetUsingProfilePictureUsers", testUserToOneSetOpAssetUsingProfilePictureAsset)
	t.Run("UserConversationToConversationUsingUserConversations", testUserConversationToOneSetOpConversationUsingConversation)
	t.Run("UserConversationToUserUsingUserConversations", testUserConversationToOneSetOpUserUsingUser)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {
	t.Run("MessageToUserUsingAuthorMessages", testMessageToOneRemoveOpUserUsingAuthor)
	t.Run("MessageToAssetUsingRawAudioMessages", testMessageToOneRemoveOpAssetUsingRawAudio)
	t.Run("UserToAssetUsingProfilePictureUsers", testUserToOneRemoveOpAssetUsingProfilePictureAsset)
}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
	t.Run("AssetToRawAudioMessages", testAssetToManyAddOpRawAudioMessages)
	t.Run("AssetToProfilePictureUsers", testAssetToManyAddOpProfilePictureUsers)
	t.Run("ConversationToMessages", testConversationToManyAddOpMessages)
	t.Run("ConversationToUserConversations", testConversationToManyAddOpUserConversations)
	t.Run("UserToAuthorMessages", testUserToManyAddOpAuthorMessages)
	t.Run("UserToUserConversations", testUserToManyAddOpUserConversations)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
	t.Run("AssetToRawAudioMessages", testAssetToManySetOpRawAudioMessages)
	t.Run("AssetToProfilePictureUsers", testAssetToManySetOpProfilePictureUsers)
	t.Run("UserToAuthorMessages", testUserToManySetOpAuthorMessages)
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
	t.Run("AssetToRawAudioMessages", testAssetToManyRemoveOpRawAudioMessages)
	t.Run("AssetToProfilePictureUsers", testAssetToManyRemoveOpProfilePictureUsers)
	t.Run("UserToAuthorMessages", testUserToManyRemoveOpAuthorMessages)
}

func TestReload(t *testing.T) {
	t.Run("DieselSchemaMigrations", testDieselSchemaMigrationsReload)
	t.Run("Assets", testAssetsReload)
	t.Run("Conversations", testConversationsReload)
	t.Run("Messages", testMessagesReload)
	t.Run("Users", testUsersReload)
	t.Run("UserConversations", testUserConversationsReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("DieselSchemaMigrations", testDieselSchemaMigrationsReloadAll)
	t.Run("Assets", testAssetsReloadAll)
	t.Run("Conversations", testConversationsReloadAll)
	t.Run("Messages", testMessagesReloadAll)
	t.Run("Users", testUsersReloadAll)
	t.Run("UserConversations", testUserConversationsReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("DieselSchemaMigrations", testDieselSchemaMigrationsSelect)
	t.Run("Assets", testAssetsSelect)
	t.Run("Conversations", testConversationsSelect)
	t.Run("Messages", testMessagesSelect)
	t.Run("Users", testUsersSelect)
	t.Run("UserConversations", testUserConversationsSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("DieselSchemaMigrations", testDieselSchemaMigrationsUpdate)
	t.Run("Assets", testAssetsUpdate)
	t.Run("Conversations", testConversationsUpdate)
	t.Run("Messages", testMessagesUpdate)
	t.Run("Users", testUsersUpdate)
	t.Run("UserConversations", testUserConversationsUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("DieselSchemaMigrations", testDieselSchemaMigrationsSliceUpdateAll)
	t.Run("Assets", testAssetsSliceUpdateAll)
	t.Run("Conversations", testConversationsSliceUpdateAll)
	t.Run("Messages", testMessagesSliceUpdateAll)
	t.Run("Users", testUsersSliceUpdateAll)
	t.Run("UserConversations", testUserConversationsSliceUpdateAll)
}
