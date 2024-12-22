{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    nixpkgs-stable.url = "github:NixOS/nixpkgs/nixos-24.11";
    utils.url = "github:numtide/flake-utils";

    go = {
      url = "github:friedenberg/dev-flake-templates?dir=go";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, nixpkgs-stable, utils, go }:
    (utils.lib.eachDefaultSystem
      (system:
        let

          pkgs = import nixpkgs {
            inherit system;
            overlays = [
              go.overlays.default
            ];
          };

          # The current default sdk for macOS fails to compile go projects, so we use a newer one for now.
          # This has no effect on other platforms.
          callPackage = pkgs.darwin.apple_sdk_11_0.callPackage or pkgs.callPackage;

        in
        {
          # packages.default = callPackage ./. {
          #   inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
          # };

          packages.default = pkgs.buildGoModule rec {
            enableParallelBuilding = true;
            doCheck = false;
            pname = "zit";
            version = "0.0.0";
            src = ./.;
            # vendorHash = pkgs.lib.fakeHash;
            vendorHash = "sha256-i+RNBc8Dl8rqgQDAQ32dtvFfvwJ/YA8TTOEWm0AdO0s=";
          };

          devShells.default = pkgs.mkShell {
            packages = (with pkgs; [
              bats
              fish
              gnumake
              just
            ]);

            inputsFrom = [
              go.devShells.${system}.default
            ];
          };
        })
    );
}
