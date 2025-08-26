{
  description = "A utility for generating Mermaid diagrams from Terraform configurations";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system}; in
      {
        packages.default = import ./default.nix { inherit pkgs; };
      }) // {
        overlays.default = final: prev: { 
          terramaid = import ./default.nix { pkgs = final; };
        };
      };
}
