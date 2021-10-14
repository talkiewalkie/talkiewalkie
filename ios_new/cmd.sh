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
      --proto_path=../protos \
      --plugin="${HOME}/.local/bin/protoc-gen-swift" \
      --swift_opt=Visibility=Public \
      --swift_out=TalkieWalkie/Clients \
      --plugin="${HOME}/.local/bin/protoc-gen-grpc-swift" \
      --grpc-swift_opt=Visibility=Public \
      --grpc-swift_out=TalkieWalkie/Clients \
      ../protos/app.proto
}

refresh_config() {
  # en0 should be the wifi on recent MacOS. YMMV.
  local_ip=$(ipconfig getifaddr en0)
  cat <<EOF > TalkieWalkie/Config.dev.plist
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>apiHost</key>
	<string>${local_ip}</string>
	<key>apiPort</key>
	<integer>8080</integer>
</dict>
</plist>
EOF
}

"$@"
