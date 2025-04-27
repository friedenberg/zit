{
  inputs = {
    nixpkgs-stable.url = "https://flakehub.com/f/NixOS/nixpkgs/0.2411.717296.tar.gz";
    utils.url = "https://flakehub.com/f/numtide/flake-utils/0.1.102.tar.gz";

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
              ExposedPorts = {"9000/tcp" = {};};
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
