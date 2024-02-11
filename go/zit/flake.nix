{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/master";
    nixpkgs-master.url = "github:NixOS/nixpkgs/master";
    utils.url = "github:numtide/flake-utils";

    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, nixpkgs-master, utils, gomod2nix }:
    (utils.lib.eachDefaultSystem
      (system:
        let

          pkgs = import nixpkgs-master {
            inherit system;
            overlays = [
              (final: prev: {
                go = prev.go_1_21;
                # buildGoModule = prev.buildGo118Module;
              })
              gomod2nix.overlays.default
            ];
          };

          zit = pkgs.buildGoApplication {
            name = "zit";
            pname = "zit";
            version = "0.0.1";
            src = ./.;
            modules = ./gomod2nix.toml;
            doCheck = false;
            enableParallelBuilding = true;
          };

        in
        {
          pname = "zit";
          packages.default = zit;
          devShells.default = pkgs.mkShell {
            buildInputs = with pkgs; [
              fish
              go
              gopls
              gotools
              /* go-tools */
              gomod2nix.packages.${system}.default
            ];
          };
        })
    );
}
