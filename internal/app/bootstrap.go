package app

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

//go:embed assets/*
var assets embed.FS

type BootstrapOptions struct {
	GodotVersion string
	GodotRepoUrl string
}

// TODO: Could change the code to act like it executes steps in a pipeline
func (b *BootstrapOptions) BootstrapCpp(projectPath string) {
	err := checkDependencies()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	valid, err := verifyGodotProject(projectPath)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	if !valid {
		log.Fatalf("Invalid Godot project path: %s", projectPath)
	}

	// Change directory to the project path
	err = os.Chdir(projectPath)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// Init new git repository
	err = initGitRepository()

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// Pull the godot-cpp engine
	err = pullGodotCppEngine(b.GodotVersion, b.GodotRepoUrl)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// TODO: Optionally build the c++ bindings. It is not necessary to
	// build the bindings if engine is same version as the godot-cpp engine
	// Example: --dump-extension-api

	// If the ext directory already exists skip creating it
	// Create a new directory for the c++ extension named ext
	if _, err := os.Stat("ext"); os.IsNotExist(err) {
		err = os.Mkdir("ext", 0755)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
	}

	// If the bin directory already exists skip creating it
	if _, err := os.Stat("bin"); os.IsNotExist(err) {
		// Create a new directory for the c++ extension binaries named bin
		err = os.Mkdir("bin", 0755)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
	}

	// Add the register_types.h and register_types.cpp files to the ext directory
	err = addRegisterTypesFiles()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// Add the example class files to the ext/example directory
	err = addExampleClassFiles()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// Add the SConstruct file to the project root
	err = addSconstruct()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	err = addGDExtensionConfig(b.GodotVersion)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	log.Println("Godot C++ extension bootstrapped successfully!")
}

func checkDependencies() error {
	// Check if git is installed
	_, err := exec.LookPath("git")
	if err != nil {
		return err
	}

	// Check if python is installed
	_, err = exec.LookPath("python")
	if err != nil {
		return err
	}

	return nil
}

func initGitRepository() error {
	err := Run("git", "init")

	if err != nil {
		return err
	}

	return nil
}

func pullGodotCppEngine(godotVersion, godotRepoUrl string) error {

	err := Run("git", "submodule", "add", "-b", godotVersion, godotRepoUrl)
	if err != nil {
		return err
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Get the last part from the godotRepoUrl to use as the submodule name
	// Example: godot-cpp
	submoduleName := godotRepoUrl[strings.LastIndex(godotRepoUrl, "/")+1:]

	// Change directory to the submodule
	err = os.Chdir(submoduleName)
	if err != nil {
		return err
	}

	// Pull the submodule
	err = Run("git", "submodule", "update", "--init")
	if err != nil {
		return err
	}

	// Change directory back to the project path
	err = os.Chdir(cwd)
	if err != nil {
		return err
	}

	return nil
}

// Verifies that the project path is a valid Godot project
func verifyGodotProject(projectPath string) (bool, error) {

	// Check if the project path exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return false, err
	}

	// Check if the project path contains a project.godot file
	if _, err := os.Stat(projectPath + "/project.godot"); os.IsNotExist(err) {
		return false, err
	}

	// Check if the project path contains a .godot dir
	if _, err := os.Stat(projectPath + "/.godot"); os.IsNotExist(err) {
		return false, err
	}

	return true, nil
}

func addRegisterTypesFiles() error {
	// Read the contents of the assets/register_types.h file
	registerTypesH, err := assets.ReadFile("assets/register_types.h")
	if err != nil {
		return err
	}

	// Write the contents of the assets/register_types.h file to the ext/register_types.h file
	err = os.WriteFile("ext/register_types.h", registerTypesH, 0644)
	if err != nil {
		return err
	}

	// Read the contents of the assets/register_types.cpp file
	registerTypesCpp, err := assets.ReadFile("assets/register_types.cpp")
	if err != nil {
		return err
	}

	// Write the contents of the assets/register_types.cpp file to the ext/register_types.cpp file
	err = os.WriteFile("ext/register_types.cpp", registerTypesCpp, 0644)
	if err != nil {
		return err
	}

	return nil
}

func addExampleClassFiles() error {
	// Make a new directory for the example class if it does not exist
	if _, err := os.Stat("ext/example"); os.IsNotExist(err) {
		err = os.Mkdir("ext/example", 0755)
		if err != nil {
			return err
		}
	}

	// Read the contents of the assets/gdexample.h file
	gdExampleH, err := assets.ReadFile("assets/gdexample.h")
	if err != nil {
		return err
	}

	// Write the contents of the assets/gdexample.h file to the ext/example/gdexample.h file
	err = os.WriteFile("ext/example/gdexample.h", gdExampleH, 0644)
	if err != nil {
		return err
	}

	// Read the contents of the assets/gdexample.cpp file
	gdExampleCpp, err := assets.ReadFile("assets/gdexample.cpp")
	if err != nil {
		return err
	}

	// Write the contents of the assets/gdexample.cpp file to the ext/example/gdexample.cpp file
	err = os.WriteFile("ext/example/gdexample.cpp", gdExampleCpp, 0644)
	if err != nil {
		return err
	}

	return nil
}

// TODO: Could allow for some configuration options on the SConstruct file
func addSconstruct() error {

	// Install SCons to python if not already installed
	err := Run("pip", "install", "SCons")
	if err != nil {
		return err
	}

	// Read the contents of the assets/SConstruct file
	sconstruct, err := assets.ReadFile("assets/SConstruct")
	if err != nil {
		return err
	}

	// Write the contents of the assets/SConstruct file to the SConstruct file
	err = os.WriteFile("SConstruct", sconstruct, 0644)
	if err != nil {
		return err
	}

	// Create build scipts for the project
	var platform string
	var buildSciptExt string
	switch os := runtime.GOOS; os {
	case "darwin":
		platform = "macos"
		buildSciptExt = ".sh"
	case "linux":
		platform = "linux"
		buildSciptExt = ".sh"
	case "windows":
		platform = "windows"
		buildSciptExt = ".bat"
	default:
		// panic for unsupported platform
		panic("Unsupported platform")
	}

	buildScirptName := fmt.Sprintf("build%s", buildSciptExt)
	format := "python -m SCons platform=%s target=%s"
	debugData := fmt.Sprintf(format, platform, "template_debug")
	releaseData := fmt.Sprintf(format, platform, "template_release")

	err = os.WriteFile("debug-"+buildScirptName, []byte(debugData), 0644)
	if err != nil {
		return err
	}

	err = os.WriteFile("release-"+buildScirptName, []byte(releaseData), 0644)
	if err != nil {
		return err
	}

	return nil
}

func addGDExtensionConfig(godotVersion string) error {

	// Read the contents of the assets/SConstruct file
	gdextensionConfig, err := assets.ReadFile("assets/gd.gdextension")
	if err != nil {
		return err
	}

	// Replace the COMPATIBILITY_MINIMUM placeholder with the godot version
	gdextensionConfig = []byte(strings.ReplaceAll(string(gdextensionConfig), "COMPATIBILITY_MINIMUM", godotVersion))

	err = os.WriteFile("bin/gd.gdextension", gdextensionConfig, 0644)
	if err != nil {
		return err
	}

	return nil
}
