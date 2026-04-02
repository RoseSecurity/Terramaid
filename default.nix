{pkgs}:
pkgs.buildGoModule rec {
  pname = "terramaid";
  version = "2.6.2";
  src = ./.;

  vendorHash = "sha256-rLIqrNgx8Vk4ijdSwGn5ye+6QYjiUYZ5zyogGx+fd/E=";

  subPackages = ["."];

  ldflags = ["-s" "-w" "-X=cmd.Version=${version}"];

  meta = with pkgs.lib; {
    description = "A utility for generating Mermaid diagrams from Terraform configurations";
    homepage = "https://github.com/RoseSecurity/Terramaid";
  };
}
