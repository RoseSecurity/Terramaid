{
  # this feels incredibly scuffed but apparently it works????
  pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
      import (fetchTree nixpkgs.locked) {
        overlays = [
          (import "${fetchTree gomod2nix.locked}/overlay.nix")
        ];
      }
  ),
  buildGoApplication ? pkgs.buildGoApplication,
}:
buildGoApplication rec {
  pname = "terramaid";
  version = "2.6.2";
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;
  CGO_ENABLED = 0;
  env.GOTOOLCHAIN = "local";

  subPackages = ["."];

  ldflags = ["-s" "-w" "-X=cmd.Version=${version}"];

  # meta = {
  #   description = "A utility for generating Mermaid diagrams from Terraform configurations";
  #   homepage = "https://github.com/RoseSecurity/Terramaid";
  # };
}
