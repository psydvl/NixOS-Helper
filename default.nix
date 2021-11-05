with import <nixpkgs>{}; # { buildGoPackage, fetchFromGitHub, pkg-config, gtk3, go, makeDesktopItem }:
buildGoPackage rec {
	pname = "NixOS-Helper";
	version = "0.1.1";

	src = fetchFromGitHub {
		owner = "psydvl";
		repo = "NixOS-Helper";
		rev = "v${version}";
		sha256 = "SHA256"; # nix-prefetch-git https://github.com/psydvl/NixOS-Helper --rev v0.1.1
	};

	nativeBuildInputs = [
		pkg-config
	];
	buildInputs = [
		gtk3
		go
	];

	goPackagePath = "github.com/psydvl/NixOS-Helper";

	goDeps = ./deps.nix;
	
	desktopItem = makeDesktopItem {
		name = pname;
		desktopName = "NixOS Helper";
		exec = "NixOS-Helper";
		icon = pname;
	};
	
	postInstall = ''
		cp -r ${pname}.svg $out/usr/share/icons/hicolor/scalable/apps/
	'';
}
