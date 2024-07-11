{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
    utils.url = "github:numtide/flake-utils";

    go = {
      url = "github:friedenberg/dev-flake-templates?dir=go";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, utils, go }:
    (utils.lib.eachDefaultSystem
      (system:
        let

          pkgs = import nixpkgs {
            inherit system;
            overlays = [
              go.overlays.default
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
            packages = (with pkgs; [
              fish
              gnumake
            ]);

            inputsFrom = [
              go.devShells.${system}.default
            ];
          };
        })
    );
}
