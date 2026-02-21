{ self }:
{
  config,
  lib,
  pkgs,
  ...
}:

with lib;

let
  cfg = config.services.bpvd;
  progCfg = config.programs.bpv;
  inherit (pkgs.stdenv.hostPlatform) system;
in
{
  options = {
    services.bpvd = {
      enable = mkEnableOption "BPV music player daemon";

      package = mkOption {
        type = types.package;
        default = self.packages.${system}.default;
        description = "The BPV package to use.";
      };
    };

    programs.bpv = {
      enable = mkEnableOption "BPV music player client";

      package = mkOption {
        type = types.package;
        default = self.packages.${system}.default;
        description = "The BPV package to use.";
      };
    };
  };

  config = mkMerge [
    (mkIf cfg.enable {
      home.packages = [ cfg.package ];

      systemd.user.services.bpvd = {
        Unit = {
          Description = "BPV Music Daemon";
          After = [ "network.target" ];
        };
        Install = {
          WantedBy = [ "default.target" ];
        };
        Service = {
          ExecStart = "${cfg.package}/bin/bpvd --no-daemonize";
          Restart = "on-failure";
          Environment = "BPV_WEB_DIR=${cfg.package}/share/bpv/dist";
        };
      };
    })

    (mkIf progCfg.enable {
      home.packages = [ progCfg.package ];
    })
  ];
}
