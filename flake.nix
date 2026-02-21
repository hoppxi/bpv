{
  description = "BPV - Music Player for your terminal and browser";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        packagesDef = import ./nix/packages.nix { inherit pkgs; };
      in
      {
        devShells.default = import ./nix/devshell.nix { inherit pkgs; };

        packages = {
          default = packagesDef.default;
          bpv = packagesDef.bpvclient;
          bpvd = packagesDef.bpvdaemon;
          web = packagesDef.bpvweb;
        };
      }
    )
    // {
      homeModules.default = import ./nix/home-manager.nix { inherit self; };
      nixosModules.default = import ./nix/nixos.nix { inherit self; };
    };
}
