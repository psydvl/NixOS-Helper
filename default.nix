{ nixpkgs ? import <nixpkgs> {  } }:

let
	pkgs = [
		nixpkgs.pkg-config
		nixpkgs.gtk3
		nixpkgs.go
	];

in
	nixpkgs.stdenv.mkDerivation {
		name = "env";
		buildInputs = pkgs;
	}
