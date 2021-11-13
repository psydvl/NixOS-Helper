{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
	nativeBuildInputs = [
		pkgs.pkg-config
		pkgs.gtk3.dev
		pkgs.go
		pkgs.gnome.gnome-terminal
	];
}
