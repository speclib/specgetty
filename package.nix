{ lib, buildGoModule, go, version ? "0.1.0" }:
buildGoModule rec {
  pname = "specgetty";

  inherit version;

  src = ./.;

  subPackages = [ "src" ];

  postInstall = ''
    mv $out/bin/src $out/bin/spg
  '';

  doCheck = false;

  vendorHash = "sha256-Lxik5/egn7vtWHcMPOBiFO/ZP5847TTEKfDXmLHHqS4=";

  meta = with lib; {
    description = ''
      Find OpenSpec projects
    '';
    homepage = "https://github.com/mipmip/specgetty";
    mainProgram = "spg";
    license = licenses.mit;
  };

}
