{
  description = "MCP (Model Context Protocol) library for Go";

  inputs = {
    nixpkgs-master.url = "github:NixOS/nixpkgs/master";
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    utils.url = "https://flakehub.com/f/numtide/flake-utils/0.1.102";
    devenv-go.url = "github:friedenberg/eng?dir=pkgs/alfa/devenv-go";
    devenv-shell.url = "github:friedenberg/eng?dir=pkgs/alfa/devenv-shell";
  };

  outputs =
    {
      self,
      nixpkgs,
      utils,
      devenv-go,
      devenv-shell,
      nixpkgs-master,
    }:
    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [
            devenv-go.overlays.default
          ];
        };

        version = "0.1.0";

        go-lib-mcp = pkgs.buildGoModule {
          pname = "go-lib-mcp";
          inherit version;
          src = ./.;
          vendorHash = null;  # Library with no dependencies

          meta = with pkgs.lib; {
            description = "MCP (Model Context Protocol) library for building MCP servers in Go";
            homepage = "https://github.com/friedenberg/go-lib-mcp";
            license = licenses.mit;
          };
        };
      in
      {
        packages = {
          default = go-lib-mcp;
          inherit go-lib-mcp;
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            just
            golangci-lint
          ];

          inputsFrom = [
            devenv-go.devShells.${system}.default
            devenv-shell.devShells.${system}.default
          ];

          shellHook = ''
            echo "go-lib-mcp: MCP library for Go - dev environment"
          '';
        };
      }
    );
}
