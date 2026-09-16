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
          # bubbletea v2 and its siblings declare `go 1.25.0`, which raises this
          # module's own directive above the 1.24 that is default in this
          # nixpkgs. buildGo125Module comes from the same pinned channel, so
          # the toolchain moves without the rest of nixpkgs moving with it.
          specgetty = pkgs.callPackage ./package.nix {
            inherit version;
            buildGoModule = pkgs.buildGo125Module;
          };
          default = pkgs.callPackage ./package.nix {
            inherit version;
            buildGoModule = pkgs.buildGo125Module;
          };
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
          tests = pkgs.buildGo125Module {
            pname = "specgetty-tests";
            inherit version;
            src = ./.;
            vendorHash = "sha256-J5Wy/Duhbw9JXVt1XhzaxdaLhNsbre8DO0s2FalZzP4=";

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
