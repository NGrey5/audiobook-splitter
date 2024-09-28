# ffmpeg Audiobook Splitter

### Overview

Takes in a compatible input audiobook file (mp3) and splits it into multiple
files without re-encoding. The tool reads a file containing labels for each chapter in the following format:

```txt
START_TIME  END_TIME  CHAPTER_NAME
```

An example file would look like so:

```txt
0.000000  100.000000  Chapter 1
100.000000  200.000000  Chapter 2
200.000000  300.000000  Chapter 3
...
```

The data for each line should be separated by a tab `/t`. This is how the program Audacity exports labels from a project and will be assumed what is being used.

### Arguments

- `-i` - the input file. DEFAULT `input.mp3`
  - This file can be an audio file of any type
- `-l` - the label file. DEFAULT `label.txt`
  - This file must be a `.txt` file in the format specified in the Overview section
- `-o` - the output directory. DEFAULT `./output`
  - If the directory is not present, it will be created

### Requirements

- `ffmpeg` - should be installed for your operating system using one of the many installation methods. Must be the command line tool as this program directly calls `ffmpeg`

### Purpose

When an audiobook is in a format that does not support chapters like a `.m4b` file does, its difficult to use. When converting a file to `.m4b` we need to split the original file up by it's chapters, before using a tool to merge them into a single chapter compatible file.

When splitting the original file into it's chapters, most common programs that do this will re-encode the file when splitting. If the original file is in a lossy format, such as `.mp3` the re-encoding will make the new individual files lose quality, even though they're in the same format as the original. Again, some programs and tools do this, others do not.

To avoid this, we use `ffmpeg` in this tool to perform as lossless split of the original file, into it's individual chapterized files. This ensures the highest quality audio result possible before merging these files into a chapterized container (`.m4b`)

### Notes

- This tool `DOES NOT` merge files into a chapterized format. It only splits an input file. You must use another tool to merge these.
- The output files of this program will always be the same file type (encoding) as the input. If your input is an `.mp3` file, then the output will be to. This is because the tool performs a lossless split and not re-encoding.

### Technical Notes

All this program does is performs the following command for each chapter in the labels file:

```console
> ffmpeg -i {INPUT_FILE} -ss {START_TIME} -to {END_TIME} -c copy -metadata track={TRACK_NO} {OUTPUT_FILE} -y
```
