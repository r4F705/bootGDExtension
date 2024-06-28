# Godot Bootstrap Project

This program can be used to bootstrap GDExtensions to a clean godot project allowing users to start writting C++ code faster!

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

You need to have Go installed on your machine. You can download it from [here](https://golang.org/dl/).

### Usage

To use this program, you need to provide the Godot project directory and the Godot version as command line arguments. Optionally, you can also provide a Godot repository URL.

Here is an example of how to use the script:

```bash
go run main.go -project=/path/to/your/project -godot-version=4.2 -godot-repo-url=https://github.com/godotengine/godot-cpp
```

* If you do not provide a Godot repository URL, the script will use "https://github.com/godotengine/godot-cpp" by default.* 

After the bootstraping is completed change directory to your godot project and run the command below to install the extentions to the project. This step needs to be repeated each time you make changes to the extention code. Also do not forget to reload your godot project for the changes to take effect.
```
python -m SCons
```

To verify that the installation was successful try to create a new *GDExample* node from the editor!

### Notes

The project should work for Windows, Linux and Macos but currently it has only been tested on Windows!


Authors \
**Nikos Raftogiannis** 

License
This project is licensed under the Apache-2.0 License - see the LICENSE file for details