{ pkgs }:
let
  sharedGoAttrs = {
    version = "0.1.0";
    src = pkgs.lib.cleanSource ../.;
    vendorHash = "sha256-DXFZW0+Kixpfn+lObvrKzAni5VP7Wm0LX5IFyZmjnmE=";
    ldflags = [
      "-s"
      "-w"
    ];
  };

  bpvclient = pkgs.buildGoModule (
    sharedGoAttrs
    // {
      pname = "bpv-client";

      nativeBuildInputs = with pkgs; [
        pkg-config
        installShellFiles
        copyDesktopItems
      ];

      buildInputs = with pkgs; [
        alsa-lib
        fdk_aac
      ];

      desktopItems = [
        (pkgs.makeDesktopItem {
          name = "bpv-tui";
          exec = "bpv";
          icon = "audio-x-generic";
          comment = "Music Player for terminal and browser";
          desktopName = "BPV TUI";
          categories = [
            "Audio"
            "Music"
            "Player"
          ];
          terminal = true;
        })
        (pkgs.makeDesktopItem {
          name = "bpv-web";
          exec = "bpv --client web";
          icon = "audio-x-generic";
          comment = "Music Player for terminal and browser";
          desktopName = "BPV Web";
          categories = [
            "Audio"
            "Music"
            "Player"
          ];
        })
      ];

      subPackages = [ "cmd/bpv" ];

      postInstall = ''
        installShellCompletion --cmd bpv \
          --bash <($out/bin/bpv completion bash) \
          --fish <($out/bin/bpv completion fish) \
          --zsh <($out/bin/bpv completion zsh) || true
      '';
    }
  );

  bpvdaemon = pkgs.buildGoModule (
    sharedGoAttrs
    // {
      pname = "bpv-daemon";

      nativeBuildInputs = with pkgs; [
        pkg-config
        copyDesktopItems
      ];

      desktopItems = [
        (pkgs.makeDesktopItem {
          name = "bpvd";
          exec = "bpvd";
          icon = "audio-x-generic";
          comment = "Music Player Daemon";
          desktopName = "BPV daemon";
          categories = [
            "Audio"
            "Music"
            "Player"
          ];
        })
      ];

      subPackages = [ "cmd/bpvd" ];
    }
  );

  bpvweb = pkgs.buildNpmPackage {
    pname = "bpv-web";
    version = "0.1.0";
    src = ../web;
    npmDepsHash = "sha256-cyuJTwOP/9qgx35SKmje94x1jqkI4iwUbOhMg3NZthk=";

    installPhase = ''
      mkdir -p $out/share/bpv
      cp -r ./dist $out/share/bpv/
    '';
  };

in
{
  inherit bpvclient bpvdaemon bpvweb;
  default = pkgs.stdenv.mkDerivation {
    pname = "bpv";
    version = "0.1.0";

    buildInputs = [
      bpvclient
      bpvdaemon
      bpvweb
    ];

    phases = [ "installPhase" ];

    installPhase = ''
      mkdir -p $out/bin
      cp -r ${bpvclient}/bin/* $out/bin/ || true
      cp -r ${bpvdaemon}/bin/* $out/bin/ || true
      mkdir -p $out/share/bpv
      cp -r ${bpvweb}/share/bpv/dist $out/share/bpv/ || true

      mkdir -p $out/share
      cp -r ${bpvclient}/share/bash-completion $out/share/ || true
      cp -r ${bpvclient}/share/fish $out/share/ || true
      cp -r ${bpvclient}/share/zsh $out/share/ || true

      mkdir $out/share/applications
      cp -r ${bpvclient}/share/applications/* $out/share/applications/ || true
      cp -r ${bpvdaemon}/share/applications/* $out/share/applications/ || true
    '';
  };
}
