# talkiewalkie

TalkieWalkie monorepo: code for server, webapp, iOS app &amp; their deployment files.

[webapp](https://web.talkiewalkie.app) - [ios app](ios)

- [audio](/audio): Python micro service for audio processing
- [kube](/kube): core kubernetes files - service deployment files are in their respective subdirectories
- [ios](/ios): ios app repo
- [protos](/protos): protobuf files for microservices<->main server interfaces
- [scripts](/scripts): code for one time projects
- [server](/server): Golang REST server
- [sqlboiler](/sqlboiler): submodule for our own version of this SQL code generation tool for the server
- [web](/web): NextJS webapp code
