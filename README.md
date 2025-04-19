# JPEG Optimizer for Windows with jpegli

This is a Windows WPF desktop application designed to optimize JPEG image files using the jpegli optimizer.

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

## Purpose

The application provides a simple graphical user interface where users can drag and drop JPEG files or folders to optimize them using the jpegli optimizer, while preserving image metadata through exiftool.

## Key Features

1. **Drag and Drop Interface**
   - Users can drag and drop individual JPEG files or entire folders
   - The app automatically scans folders recursively for all JPEG files

2. **Optimization Settings**
   - **Distance (Quality)**: A slider control that adjusts the optimization level from 0.0 to 3.0 (default: 1.0)
     - Higher values produce smaller files but potentially lower quality
   - **Overwrite Files**: Toggle option to either overwrite original files or save optimized versions with a ".jpegli.jpg" extension

3. **Processing Capabilities**
   - Processes images asynchronously to keep the UI responsive
   - Uses jpegli.exe for image optimization
   - Uses exiftool.exe to preserve metadata from the original images
   - Displays real-time progress and statistics

4. **Feedback & Progress Tracking**
   - Log area showing detailed processing information for each file
   - Status bar displaying overall progress
   - File size savings calculation (both in bytes and percentage)

## Technical Implementation

- Built using .NET (targeting net9.0-windows)
- WPF (Windows Presentation Foundation) UI framework
- Asynchronous processing using Tasks
- External tool integration (jpegli and exiftool)
- File system operations for managing optimized files

## Example Workflow

Here's a practical example of how to use this tool in a photography workflow:

1. Export images from Lightroom or Capture One with JPEG 100% Quality
2. Move these files or the entire folder to this application
3. Process them with the recommended quality setting of 0.5 (distance)
   - Distance of 0.5 is for the highest quality
   - Distance of 1.0 is considered visually lossless
   - Recommended range: 0.5 to 3.0
   - Allowed range: 0.0 to 25.0
   - Lower values maintain higher quality while still achieving good compression

In a real-world example, processing 35 files (45 MP each) with a distance setting of 0.5 reduced the total size from 918.09 MB to 273.85 MB - achieving a compression ratio of 0.3 (70% size reduction) while maintaining visual quality.

## Usage Flow

1. User adjusts optimization settings (distance slider and overwrite option)
2. User drags and drops JPEG files or folders onto the application
3. The application processes each image with the selected settings
4. Progress and results are displayed in real-time
5. Optimized files are saved according to user preferences

## Dependencies

The application requires:
- jpegli.exe - The core optimization engine
- exiftool.exe - For preserving image metadata

These tools should be placed in a "tools" folder in the application's directory.

For those who prefer command-line operations, the included PowerShell script `convert_with_jpegli.ps1` demonstrates how to perform the same optimization using a script with these executables directly.

The project is compiled and ready to run, with the executable available at bin/Debug/net9.0-windows/jpegli-windows-gui.exe.