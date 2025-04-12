{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    nixpkgs-stable.url = "nixpkgs/release-24.11";
    utils.url = "github:numtide/flake-utils";

    go = {
      url = "github:friedenberg/dev-flake-templates?dir=go";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    shell = {
      url = "github:friedenberg/dev-flake-templates?dir=shell";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = {
    self,
    nixpkgs,
    nixpkgs-stable,
    utils,
    go,
    shell,
  }:
    (utils.lib.eachDefaultSystem
      (system:
        let

          pkgs = import nixpkgs {
            inherit system;

            overlays = [
              go.overlays.default
            ];
          };

          callPackage = pkgs.darwin.apple_sdk_11_0.callPackage or pkgs.callPackage;

          gomod2nix = pkgs.gomod2nix;

          zit = pkgs.buildGoApplication {
            pname = "zit";
            version = "0.0.1";
            src = ./.;
            modules = ./gomod2nix.toml;
          };

        in {

          packages.zit = zit;
          packages.default = zit;

          docker = pkgs.dockerTools.buildImage {
            name = "zit";
            tag = "latest";
            copyToRoot = [zit];
            config = {
              Cmd = ["${zit}/bin/zit"];
              Env = [];
              ExposedPorts = {"28669/tcp" = {};};
            };
          };

          devShells.default = pkgs.mkShell {
            # inherit (gomod2nix.packages.${system}) mkGoEnv gomod2nix;

            packages = (with pkgs; [
              govulncheck
              bats
              fish
              gnumake
              just
            ]);

            inputsFrom = [
              go.devShells.${system}.default
              shell.devShells.${system}.default
            ];
          };
        })
    );
}
