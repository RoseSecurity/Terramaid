{pkgs}:
pkgs.buildGoModule rec {
  pname = "terramaid";
  version = "2.6.2";
  src = ./.;

  vendorHash = "sha256-ZnYGyoKYnx9inSA2FryFds5jAf3L45nAsgm5ElXtv9Y=";

  subPackages = ["."];

  ldflags = ["-s" "-w" "-X=cmd.Version=${version}"];

  meta = with pkgs.lib; {
    description = "A utility for generating Mermaid diagrams from Terraform configurations";
    homepage = "https://github.com/RoseSecurity/Terramaid";
  };
}
