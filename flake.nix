{
  description = "A utility for generating Mermaid diagrams from Terraform configurations";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    nixpkgs-go.url = "github:nixos/nixpkgs/fcd023ec9e17481b4f80ffec0e9d0f36ed847b91";
    flake-utils.url = "github:numtide/flake-utils";

    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs = {
        nixpkgs.follows = "nixpkgs";
        flake-utils.follows = "flake-utils";
      };
    };
  };
  outputs = {
    self,
    nixpkgs,
    nixpkgs-go,
    flake-utils,
    gomod2nix,
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {
        inherit system;
        overlays = [
          gomod2nix.overlays.default
          (final: prev: {
            go = nixpkgs-go.packages.${system}.go;
          })
        ];
      };
      callPackage = pkgs.callPackage;
    in {
      packages.default = callPackage ./. {
        inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
      };
      devShells.default = callPackage ./shell.nix {
        inherit (gomod2nix.legacyPackages.${system}) mkGoEnv gomod2nix;
      };
    })
    // {
      overlays.default = final: prev: {
        terramaid = import ./default.nix {
          pkgs = final.extend gomod2nix.overlays.default;
        };
      };
    };
}
