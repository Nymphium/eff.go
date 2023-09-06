{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    opam-repository = { url = "github:ocaml/opam-repository"; flake = false; };

    flake-utils.url = "github:numtide/flake-utils";

    opam-nix = {
      url = "github:tweag/opam-nix";
      inputs = {
        nixpkgs.follows = "nixpkgs";
        flake-utils.follows = "flake-utils";
        opam-repository.follows = "opam-repository";
      };
    };
  };
  outputs = { self, flake-utils, opam-nix, nixpkgs, opam-repository, ... }:
    flake-utils.lib.eachDefaultSystem (system: 
      let
        pkgs = nixpkgs.legacyPackages.${system};
        go = pkgs.go_1_21;
        gopls = pkgs.gopls;
        golangci-lint = pkgs.golangci-lint;
        gotest = pkgs.gotest;
      in
      {
          legacyPackages = pkgs;
          devShells.default =
            pkgs.mkShell {
              buildInputs = [ go gopls golangci-lint gotest pkgs.nil pkgs.nixpkgs-fmt ];
            };
          formatter = pkgs.nixpkgs-fmt;
      });
}
