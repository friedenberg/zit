{ pkgs ? import <nixpkgs> {}
}:

pkgs.mkShell {
  buildInputs = [
    pkgs.gnumake
    pkgs.bash
    pkgs.go
    pkgs.git
    pkgs.yubikey-agent
  ];

  shellHook = ''
    export GOPATH=$HOME/eng/go
  '';
}
