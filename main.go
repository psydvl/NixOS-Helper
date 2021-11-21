package main

import (
	_ "embed"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	search_nixos_packages = iota
	search_nixos_wiki
)

const (
	shell_configuration_check = iota
	shell_update
	shell_rebuild
	shell_garbage
	shell_optimize
)

//go:embed main.ui
var main_ui string

type Config struct {
	files []string
}

type Logger struct {
	logger  *gtk.Label
	logText []string
}

type LogEvent struct {
	err error
	str string
}

func (le LogEvent) Error() (result string) {
	if le.err == nil {
		result = "[OK]: " + le.str
	} else {
		result = le.err.Error() + "\t" + le.str
	}
	return
}

func (l *Logger) Init(errEvent <-chan LogEvent, wgDone func()) {
	var count int
	var logLine, outputText string
	var logText []string = []string{}

	for e := range errEvent {
		logLine = e.Error()
		logText = append(logText, logLine)
		fmt.Println(logLine)
		if l.logger != nil {
			if count = len(logText) - 4; count < 0 {
				count = 4
			}
			outputText = "log: " + strings.Join(logText[count:], "\nlog: ")
			l.logger.SetText(outputText)
		}
	}
	wgDone()
}

func shell(option int) error {
	var title string
	var text string
	var bash_pressanykey string = ` ; read -n 1 -s -r -p 'Press any key to continue'`
	switch option {
	case shell_configuration_check:
		title = "nixos-rebuild dry-build"
		text = title
	case shell_update:
		title = "sudo nix-channel --update"
		text = title
	case shell_rebuild:
		title = "sudo nixos-rebuild boot"
		text = title
	case shell_garbage:
		title = "sudo nix-collect-garbage -d"
		text = title
	case shell_optimize:
		title = "sudo nix optimise-store"
		text = title
	}
	err := exec.Command("gnome-terminal", "--title="+title, "--", "bash", "-c",
		text+bash_pressanykey,
	).Run()
	return err
}

func gedit_open(files []string) error {
	var launch []string

	for i, f := range files {
		files[i] = "admin:///etc/nixos/" + f + ".nix"
	}

	launch = append([]string{
		"-s",
		"admin:///etc/nixos/configuration.nix",
	}, files...)

	var err error = exec.Command("gedit",
		launch...,
	).Start()
	return err
}

func search(option int, text string) error {
	switch option {
	case search_nixos_packages:
		text = "https://search.nixos.org/packages?channel=21.05&from=0&size=50&sort=relevance&query=" + text
	case search_nixos_wiki:
		text = "https://nixos.wiki/index.php?go=Go&search=" + text
	}
	err := exec.Command("xdg-open", text).Start()
	return err
}

func read_config() (Config, error) {
	var files = []string{
		"packages",
		"specialisation", //fileSystems-home hardware-configuration networking-wireless-networks
	}
	var conf = Config{
		files: files,
	}
	return conf, nil
}

func main() {
	var logChan chan LogEvent = make(chan LogEvent)
	var err error
	var b *gtk.Builder
	var obj glib.IObject
	var logger Logger = Logger{
		logger: nil,
	}

	var wg sync.WaitGroup

	var config Config

	wg.Add(1)
	go logger.Init(logChan, wg.Done)

	config, err = read_config()
	logChan <- LogEvent{err, "Config read"}

	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	b, err = gtk.BuilderNew()
	logChan <- LogEvent{err, "Builder create"}
	if err != nil {
		log.Fatal("Unable to create builder:", err)
	}

	err = b.AddFromString(main_ui)
	logChan <- LogEvent{err, "Window file load"}
	if err != nil {
		log.Fatal("Unable to load window file:", err)
	}

	obj, err = b.GetObject("window_main")
	logChan <- LogEvent{err, "Window load"}
	if err != nil {
		log.Fatal("Unable to load window:", err)
	}

	win := obj.(*gtk.Window)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	obj, err = b.GetObject("log")
	logChan <- LogEvent{err, "Logger label load"}
	logger.logger = obj.(*gtk.Label)
	logChan <- LogEvent{nil, "Loaded"}

	obj, _ = b.GetObject("search_text")
	search_text := obj.(*gtk.SearchEntry)
	obj, _ = b.GetObject("wiki")
	wiki_button := obj.(*gtk.Button)

	obj, _ = b.GetObject("open_settings")
	open_settings := obj.(*gtk.Button)

	obj, _ = b.GetObject("gedit")
	gedit := obj.(*gtk.Button)
	obj, _ = b.GetObject("configuration_check")
	configuration_check := obj.(*gtk.Button)
	obj, _ = b.GetObject("update")
	update := obj.(*gtk.Button)
	obj, _ = b.GetObject("rebuild")
	rebuild := obj.(*gtk.Button)
	obj, _ = b.GetObject("garbage")
	garbage := obj.(*gtk.Button)
	obj, _ = b.GetObject("optimize")
	optimize := obj.(*gtk.Button)

	open_settings.Connect("clicked", func() {
		logChan <- LogEvent{nil, "WIP: settings"}
	})

	search_text.Connect("activate", func() {
		if text, err := search_text.GetText(); err == nil {
			err = search(search_nixos_packages, text)
			logChan <- LogEvent{err, ("Search in nix: " + text)}
		} else {
			logChan <- LogEvent{err, "GetText for Search nix"}
		}
	})

	wiki_button.Connect("clicked", func() {
		if text, err := search_text.GetText(); err == nil {
			err = search(search_nixos_wiki, text)
			logChan <- LogEvent{err, ("Search in wiki: " + text)}
		} else {
			logChan <- LogEvent{err, "GetText for Search wiki"}
		}
	})

	gedit.Connect("clicked", func() {
		err = gedit_open(config.files)
		logChan <- LogEvent{err, "Open gedit"}
	})

	configuration_check.Connect("clicked", func() {
		err = shell(shell_configuration_check)
		logChan <- LogEvent{err, "Check configuration"}
	})

	update.Connect("clicked", func() {
		err = shell(shell_update)
		logChan <- LogEvent{err, "Update nix channels"}
	})

	rebuild.Connect("clicked", func() {
		err = shell(shell_rebuild)
		logChan <- LogEvent{err, "Rebuild"}
	})

	garbage.Connect("clicked", func() {
		err = shell(shell_garbage)
		logChan <- LogEvent{err, "Collect garbage"}
	})

	optimize.Connect("clicked", func() {
		err = shell(shell_optimize)
		logChan <- LogEvent{err, "Optimize store"}
	})

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
	close(logChan)
	wg.Wait()

	// wg.Add(2), saveConfig(), saveLog(), wg.Wait()
}
