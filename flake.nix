{
  description = "manifestor";

  outputs = { self, nixpkgs, systems }:
    let
      eachSystem = f: nixpkgs.lib.genAttrs (import systems) (system: f nixpkgs.legacyPackages.${system});
    in
    {
      packages = eachSystem (pkgs: {
        default = pkgs.buildGoModule {
          pname = "manifestor";
          version = "0";
          src = ./.;
          vendorHash = "sha256-yTBPCTuZvrnBJnSjqP1KlOwK9lFoq5YllGyGrWCwJOw=";
        };
      });

      devShells = eachSystem (pkgs: {
        default = import ./shell.nix { inherit pkgs; };
      });
    };
}
