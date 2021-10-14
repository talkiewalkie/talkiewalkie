#!/usr/bin/env bash

ios_dev_install() {
  mkdir -p "${HOME}/.local/{share,bin}"
  echo "-- make sure $HOME/.local/bin is in your PATH (currently: [$PATH])"

  echo "-- installing protoc plugins to generate swift code for grpc"
  pushd "${HOME}/.local/share" || exit
  git clone https://github.com/grpc/grpc-swift.git

  pushd grpc-swift || exit
  make plugins
  cp protoc-gen-grpc-swift protoc-gen-swift $HOME/.local/bin/
  popd || exit

  popd || exit
}

codegen() {
  protoc -I=protos \
      --proto_path=../server/protos \
      --plugin="${HOME}/.local/bin/protoc-gen-swift" \
      --swift_opt=Visibility=Public \
      --swift_out=TalkieWalkie/Clients \
      --plugin="${HOME}/.local/bin/protoc-gen-grpc-swift" \
      --grpc-swift_opt=Visibility=Public \
      --grpc-swift_out=TalkieWalkie/Clients \
      ../server/protos/app.proto
}

"$@"
