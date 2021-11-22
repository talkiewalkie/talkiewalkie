package pb

// MessageContentOneOf Necessary to instantiate a blank oneof instance
// https://github.com/golang/protobuf/issues/261
type MessageContentOneOf = isMessage_Content
type MessageSendInputContentOneOf = isMessageSendInput_Content
type EventContentOneOf = isEvent_Content
