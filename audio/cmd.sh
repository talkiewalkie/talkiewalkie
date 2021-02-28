#!/usr/bin/env bash

codegen() {
  python -m grpc_tools.protoc -I=protos/ \
    --python_out=audio_proc/ --grpc_python_out=audio_proc/ \
    protos/audio_proc.proto
}

"$@"
