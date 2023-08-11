{ pkgs ? import <nixpkgs> }:
pkgs.mkShell {
  packages = [
    pkgs.gcc
    pkgs.go_1_18
    pkgs.gotools
    pkgs.gopls
    pkgs.go-outline
    pkgs.gocode
    pkgs.gopkgs
    pkgs.gocode-gomod
    pkgs.godef
    pkgs.golint
    pkgs.gh
    pkgs.delve
    pkgs.go-tools
  ];
}
