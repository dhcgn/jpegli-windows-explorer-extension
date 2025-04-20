# JPEG Optimizer for Windows with jpegli

> ⚠️⚠️⚠️ Very early version under active development ⚠️⚠️⚠️

This is a command-line application written in Go for optimizing JPEG image files using the jpegli optimizer. When run without arguments, it installs itself as a Windows Explorer context menu entry, allowing you to optimize JPEGs directly from the right-click menu on files or folders.

## About jpegli

> **WARNING:** jpegli is an experimental technology and is NOT production ready. Everything you do with this tool is at your own risk!

jpegli is a new JPEG coding library developed by Google Research that maintains high backward compatibility while offering enhanced capabilities and a 35% compression ratio improvement at high quality compression settings.

Key features of jpegli include:
- Full interoperability with the original JPEG standard
- Higher quality results with fewer artifacts
- Comparable coding speed to traditional approaches like libjpeg-turbo and MozJPEG
- Support for 10+ bits per component (while maintaining compatibility with 8-bit viewers)
- More efficient compression than traditional JPEG codecs

For more information, see: [Introducing Jpegli: A New JPEG Coding Library](https://opensource.googleblog.com/2024/04/introducing-jpegli-new-jpeg-coding-library.html)

## Screenshots

### Demo

![Screenshot of cli on executing on folder](docs/folder-execute.png)

![Screenshot of explorer after executing on folder](docs/folder-result.png)

### Installation

![Screenshot of explorer after executing on folder](docs/install.png)

## Purpose

This CLI application provides a simple way to optimize JPEG files or entire folders using the jpegli optimizer, while preserving image metadata through exiftool. The application can be invoked directly from the command line or via the Windows Explorer context menu after installation.

## Key Features

1. **Context Menu Integration**
   - When run without arguments, the app installs itself as a Windows Explorer context menu entry for both files and folders.
   - After installation, you can right-click any file or folder to optimize JPEGs using the context menu.

2. **Drag and Drop & CLI Usage**
   - You can run the app from the command line, passing a file or folder as an argument to optimize JPEGs.

3. **Optimization Settings**
   - Uses a default distance (quality) setting for jpegli, which can be adjusted in the code.
   - Preserves metadata using exiftool.

4. **Embedded Tools**
   - Both jpegli.exe and exiftool.exe are embedded within the application and extracted as needed. No manual download is required.

5. **Processing Capabilities**
   - Processes images asynchronously to keep the UI responsive (if used in a GUI context).
   - Displays real-time progress and statistics in the CLI.

## Example Workflow

- Run the executable without arguments to install the context menu integration.
- Right-click a file or folder in Windows Explorer and select the optimization option.
- Or, run the executable from the command line with a file or folder as an argument to optimize JPEGs.

## Dependencies

The application embeds:
- jpegli.exe - The core optimization engine
- exiftool.exe - For preserving image metadata

No manual setup of these tools is required.

For those who prefer command-line operations, the included PowerShell script `convert_with_jpegli.ps1` demonstrates how to perform the same optimization using a script with these executables directly.

The project is compiled and ready to run, with the executable available at bin/Debug/net9.0-windows/jpegli-windows-gui.exe.