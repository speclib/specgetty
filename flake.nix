{
  description = "Find OpenSpec projects on your local machine" ;

  inputs.nixpkgs.url = "nixpkgs/nixos-25.05";

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
      version = builtins.replaceStrings ["\n"] [""] (builtins.readFile ./src/VERSION);
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          specgetty = pkgs.callPackage ./package.nix { inherit version; };
          default = pkgs.callPackage ./package.nix { inherit version; };
        });

      checks = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          # The package must build.
          build = self.packages.${system}.specgetty;

          # go vet, the full test suite, and the coverage ratchet.
          # scripts/coverage-gate.sh is the single source of truth for the
          # floors, so `nix flake check` and a local run enforce the same thing.
          tests = pkgs.buildGoModule {
            pname = "specgetty-tests";
            inherit version;
            src = ./.;
            vendorHash = "sha256-Lxik5/egn7vtWHcMPOBiFO/ZP5847TTEKfDXmLHHqS4=";

            nativeBuildInputs = [ pkgs.bash ];

            buildPhase = "true";

            doCheck = true;
            checkPhase = ''
              runHook preCheck
              export HOME="$TMPDIR"
              bash scripts/coverage-gate.sh
              runHook postCheck
            '';

            installPhase = ''
              mkdir -p $out
              echo "tests and coverage ratchet passed" > $out/result
            '';
          };
        });

      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
            ];
          };
        });
    };
}
