package main

import (
	"bootGDExtension/internal/app"
	"flag"
)

// TODO: Add some dynamic flags to the project that make changes to assets like the COMPATIBILITY_MINIMUM in gd.gdextension
// TODO: Add Linux support
// TODO: Add Macos support
// TODO: Add SCons script to a build directory

// TODO: Make assets embeddable in the binary
func main() {
	// Read options from command line arguments
	options := app.BootstrapOptions{}

	godotProject := flag.String("project", "", "The godot project directory")
	godotVersion := flag.String("godot-version", "", "Godot version")
	godotRepoUrl := flag.String("godot-repo-url", "https://github.com/godotengine/godot-cpp", "Godot repository URL")

	flag.Parse()

	if *godotProject == "" {
		panic("Godot project directory is required")
	}

	if *godotVersion == "" {
		panic("Godot version is required")
	}

	options.GodotVersion = *godotVersion

	if *godotRepoUrl != "" {
		options.GodotRepoUrl = *godotRepoUrl
	}

	// Bootstrap the project
	options.BootstrapCpp(*godotProject)

}
