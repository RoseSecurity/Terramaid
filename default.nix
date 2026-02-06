{pkgs}:
pkgs.buildGoModule rec {
  pname = "terramaid";
  version = "2.6.2";
  src = ./.;

  vendorHash = "sha256-CfUrdpNkjkmXO09eyDDWmYNy7vDGXxQCiKmPI3uo96s=";

  subPackages = ["."];

  excludedPackages = ["tools"];

  ldflags = ["-s" "-w" "-X=cmd.Version=${version}"];

  meta = with pkgs.lib; {
    description = "A utility for generating Mermaid diagrams from Terraform configurations";
    homepage = "https://github.com/RoseSecurity/Terramaid";
  };
}
