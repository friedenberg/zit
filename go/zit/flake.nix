{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    (utils.lib.eachDefaultSystem
      (system:
        let

          pkgs = import nixpkgs {
            inherit system;
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
              gnumake
              parallel
            ];
          };
        })
    );
}
