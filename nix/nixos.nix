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

      user = mkOption {
        type = types.str;
        default = "bpvd";
        description = "User to run the daemon as.";
      };

      group = mkOption {
        type = types.str;
        default = "bpvd";
        description = "Group to run the daemon as.";
      };

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
      environment.systemPackages = [ cfg.package ];

      users.users = mkIf (cfg.user == "bpvd") {
        bpvd = {
          isSystemUser = true;
          group = cfg.group;
          home = "/var/lib/bpvd";
          createHome = true;
        };
      };

      users.groups = mkIf (cfg.group == "bpvd") {
        bpvd = { };
      };

      systemd.services.bpvd = {
        description = "BPV Music Daemon";
        after = [ "network.target" ];
        wantedBy = [ "multi-user.target" ];
        serviceConfig = {
          User = cfg.user;
          Group = cfg.group;
          ExecStart = "${cfg.package}/bin/bpvd --no-daemonize";
          Restart = "on-failure";
          Environment = "BPV_WEB_DIR=${cfg.package}/share/bpv/dist";
          StateDirectory = "bpvd";
        };
      };
    })

    (mkIf progCfg.enable {
      environment.systemPackages = [ progCfg.package ];
    })
  ];
}
