{ pkgs }:
pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    pkg-config
    alsa-lib
    fdk_aac
    nodejs
  ];

  CGO_ENABLED = "1";
}
