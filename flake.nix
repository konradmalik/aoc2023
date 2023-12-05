{
  description = "AOC2023";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs = { self, nixpkgs, nixpkgs-unstable, ... }:
    let
      mkOverlay = input: name: (final: prev: {
        "${name}" = import input {
          system = final.system;
          config = final.config;
        };
      });

      forAllSystems = function:
        nixpkgs.lib.genAttrs [
          "x86_64-linux"
          "aarch64-linux"
          "x86_64-darwin"
          "aarch64-darwin"
        ]
          (system:
            function (import nixpkgs {
              inherit system;
              config.allowUnfree = true;
              overlays = [
                (mkOverlay nixpkgs-unstable "unstable")
              ];
            }));
    in
    {
      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell
          {
            name = "aoc2023";

            packages = with pkgs; [
              gcc
              go
            ];
          };
      });
      formatter = forAllSystems (pkgs: pkgs.nixpkgs-fmt);
    };
}
