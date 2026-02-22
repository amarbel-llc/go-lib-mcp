{
  description = "MCP (Model Context Protocol) library for Go";

  inputs = {
    nixpkgs-master.url = "github:NixOS/nixpkgs/5b7e21f22978c4b740b3907f3251b470f466a9a2";
    nixpkgs.url = "github:NixOS/nixpkgs/6d41bc27aaf7b6a3ba6b169db3bd5d6159cfaa47";
    utils.url = "https://flakehub.com/f/numtide/flake-utils/0.1.102";
    go.url = "github:amarbel-llc/eng?dir=devenvs/go";
    shell.url = "github:amarbel-llc/eng?dir=devenvs/shell";
  };

  outputs =
    {
      self,
      nixpkgs,
      utils,
      go,
      shell,
      nixpkgs-master,
    }:
    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [
            go.overlays.default
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
            homepage = "https://github.com/amarbel-llc/go-lib-mcp";
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
            go.devShells.${system}.default
            shell.devShells.${system}.default
          ];

          shellHook = ''
            echo "go-lib-mcp: MCP library for Go - dev environment"
          '';
        };
      }
    );
}
