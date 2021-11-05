package main

import (
	_ "embed"
	"fmt"
	"log"
	"os/exec"

	"github.com/gotk3/gotk3/gtk"
)

//go:embed main.glade
var main_glade string

var bash_pressanykey string = ` ; read -n 1 -s -r -p 'Press any key to continue'`

func main() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	b, err := gtk.BuilderNew()
	if err != nil {
		log.Fatal("Unable to create builder:", err)
	}

	err = b.AddFromString(main_glade)
	if err != nil {
		log.Fatal("Unable to load window file:", err)
	}

	obj, err := b.GetObject("window_main")
	if err != nil {
		log.Fatal("Unable to load window:", err)
	}

	win := obj.(*gtk.Window)
	win.SetTitle("NixOS helper")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	obj, _ = b.GetObject("log")
	log_text := obj.(*gtk.Label)
	log_text.SetText("log: loaded")

	obj, _ = b.GetObject("search_text")
	search_text := obj.(*gtk.Entry)
	obj, _ = b.GetObject("search_button")
	search_button := obj.(*gtk.Button)
	obj, _ = b.GetObject("wiki")
	wiki_button := obj.(*gtk.Button)

	obj, _ = b.GetObject("gedit")
	gedit := obj.(*gtk.Button)
	obj, _ = b.GetObject("update")
	update := obj.(*gtk.Button)
	obj, _ = b.GetObject("rebuild")
	rebuild := obj.(*gtk.Button)
	obj, _ = b.GetObject("garbage")
	garbage := obj.(*gtk.Button)
	obj, _ = b.GetObject("optimize")
	optimize := obj.(*gtk.Button)

	search_text.Connect("activate", func() {
		text, err := search_text.GetText()
		text = "https://search.nixos.org/packages?channel=21.05&from=0&size=50&sort=relevance&query=" + text
		if err == nil {
			err = exec.Command("xdg-open", text).Start()
			log_text.SetText(fmt.Sprintf("log: %v", err))
		}
	})

	search_button.Connect("clicked", func() {
		text, err := search_text.GetText()
		text = "https://search.nixos.org/packages?channel=21.05&from=0&size=50&sort=relevance&query=" + text
		if err == nil {
			err = exec.Command("xdg-open", text).Start()
			log_text.SetText(fmt.Sprintf("log: %v", err))
		}
	})

	wiki_button.Connect("clicked", func() {
		text, err := search_text.GetText()
		text = "https://nixos.wiki/index.php?go=Go&search=" + text
		if err == nil {
			err = exec.Command("xdg-open", text).Start()
			log_text.SetText(fmt.Sprintf("log: %v", err))
		}
	})

	gedit.Connect("clicked", func() {
		err = exec.Command("gedit", "admin:///etc/nixos/configuration.nix").Start()
		log_text.SetText(fmt.Sprintf("log: %v", err))
	})

	update.Connect("clicked", func() {
		err = exec.Command("gnome-terminal", "--", "bash", "-c",
			"sudo nix-channel --update"+bash_pressanykey,
		).Start()
		log_text.SetText(fmt.Sprintf("log: %v", err))
	})

	rebuild.Connect("clicked", func() {
		err = exec.Command("gnome-terminal", "--", "bash", "-c",
			"sudo nixos-rebuild boot"+bash_pressanykey,
		).Start()
		log_text.SetText(fmt.Sprintf("log: %v", err))
	})

	garbage.Connect("clicked", func() {
		err = exec.Command("gnome-terminal", "--", "bash", "-c",
			"sudo nix-collect-garbage -d"+bash_pressanykey,
		).Start()
		log_text.SetText(fmt.Sprintf("log: %v", err))
	})

	optimize.Connect("clicked", func() {
		err = exec.Command("gnome-terminal", "--", "bash", "-c",
			"sudo nix optimise-store"+bash_pressanykey,
		).Start()
		log_text.SetText(fmt.Sprintf("log: %v", err))
	})

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}
