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

  vendorHash = "sha256-DWWzfif21IDuYdwa6PwiBQFa0gAi4NZ6YDKbE9/C4eE=";

  meta = with lib; {
    description = ''
      Find OpenSpec projects
    '';
    homepage = "https://github.com/mipmip/specgetty";
    mainProgram = "spg";
    license = licenses.mit;
  };

}
