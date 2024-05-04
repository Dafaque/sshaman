{ pkgs ? import <nixpkgs> {} }:
let 
    protoFiles = "api/remote.proto";
    outputDir = "./pkg/remote";
in
pkgs.mkShell {
  buildInputs = [
    pkgs.protobuf
    pkgs.go_1_22
    pkgs.protoc-gen-go
    pkgs.protoc-gen-go-grpc
  ];
  
  shellHook = ''
    set +e
    echo "Compiling proto files..."
    protoc\
      --go_out=${outputDir}\
      --go_opt=paths=source_relative\
      --go-grpc_out=${outputDir}\
      --go-grpc_opt=paths=source_relative\
      ${protoFiles}
    exit 0
  '';
}
