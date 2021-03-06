{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/aff647e2704fa1223994604887bb78276dc57083.tar.gz") {} }:

pkgs.mkShell {
	buildInputs = [
		pkgs.git
		pkgs.nodejs-14_x
		pkgs.yarn
		pkgs.nodePackages.typescript
		pkgs.go
		pkgs.air
		pkgs.gnumake
		pkgs.watchman
    pkgs.openssl
	];
}
