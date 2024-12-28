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

        in
        {
          packages.default = pkgs.buildGoModule rec {
            enableParallelBuilding = true;
            doCheck = false;
            pname = "zit";
            version = "0.0.0";
            src = ./.;
            vendorHash = "sha256-kBa31IFXDkb1o0SO8EIuhvkvHFVqvOD3Tdn1JtP43Qc=";
          };

          devShells.default = pkgs.mkShell {
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
