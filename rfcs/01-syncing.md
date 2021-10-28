# Syncing clients & support from the backend

## What we want

We want to offer feature parity with mainstream chat apps, including being offline friendly. That could mean in a
distant future relying on the same stack as they do - XMPP complying servers, but for now we'd rather do as little as we
can to ensure basic functionality.

## Problem

Keeping a permanent local store on clients implies the possibility of conflicts between local and remote state. We'd
want to have a system that's as reliable as possible to handle this.

## Proposed solution

We distinguish between two classes of events: those pertaining to transitory state, and persistent objects' events. A "
transitory event" is e.g. `UserIsTyping` and `UserIsConnected`. The other category of events is described below - I am
not sure how to best describe the dichotomy - but basically there are events used to compute "hard state", and some
attributes of the whole state are transitory (connected, typing). Tell me if this is wrong.

We eventify all changes:

- profile pic change is an event,
- conversation name change is an event,
- new message is an event,
- someone joined is an event

Events are stored in postgres with strong constraints and uuids. Clients keep the uuid of the last event they ingested,
and send it during `sync` calls so that the server can compile the list of events the client needs to process.

The `sync` call happens every time the network connection starts up.

Client side, we keep a local queue of these events. The queue automatically dispatches events, and when events are not
acknowledged by the server (which would return the uuid of the stored event) - for any reason, they're added up to the
queue, which always dispatches all events in its stead.

```protobuf
message Event {
  string uuid = 1;
  oneof content {
    // -- USERS
    ProfilePictureChanged profilePictureChanged = 2;
    DisplayNameChanged displayNameChanged = 8;

    // -- CONVERSATIONS
    ConversationNameChanged conversationNameChanged = 3;

    // -- MESSAGES
    // ---- Sender events
    SendingNewMessage sendingNewMessage = 4;
    SentNewMessage sentNewMessage = 5;
    ReadNewMessage readNewMessage = 6;
    // ---- Receiver events
    NewMessage newMessage = 7;
  }
}

message ProfilePictureChanged {
  string userUuid = 1;
  bytes pngImage = 2;
}

message ConversationNameChanged {
  string conversationUuid = 1;
  string newName = 2;
  string authorUuid = 3;
  google.protobuf.Timestamp occuredAt = 4;
}

message SendingNewMessage {
  string coredataLocalUuid = 1;
  Message message = 2;
}

message SentNewMessage {
  string uuid = 1;
}

message ReadNewMessage {
  string uuid = 1;
}

message NewMessage {
  string coreDataTemporaryUuid = 1;
  string uuid = 2;
}

// -- SERVICES

message SubmitEventOutput {
  repeated Event events = 1;
}


service EventService {
  rpc SyncFromServer(Empty) returns (stream Event) {};
  rpc Sync(repeated Event) returns (SubmitEventOutput) {};
}
```



