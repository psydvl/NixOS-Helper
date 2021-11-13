with import <nixpkgs>{}; # { buildGoPackage, fetchFromGitHub, pkg-config, gtk3, go, makeDesktopItem }:
buildGoPackage rec {
	pname = "NixOS-Helper";
	version = "0.1.2";

	src = ./.;
	#src = github.nix;

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
		install -Dm644 \
			${desktopItem}/share/applications/${pname}.desktop \
			$out/share/applications/${pname}.desktop
		install -Dm644 \
			$src/${pname}.svg \
			$out/share/icons/hicolor/scalable/apps/${pname}.svg
	'';
}
