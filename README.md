# 🩰 ascii

![Example dog in ascii](./.github/screenshots/dog.png)

Convert images and videos to ascii.

## 🥌 Usage

| Args | Description   | Default         |
| ---- | ------------- | --------------- |
| -i   | Input file    |                 |
| -w   | Ascii width   | Terminal width  |
| -h   | Ascii height  | Terminal height |
| -s   | Square aspect | False           |

```sh
ascii -i <INPUT> [-w <WIDTH>] [-h <HEIGHT>] [-s]
```

### 🧩 Supported formats

| Image          | Video    |
| -------------- | -------- |
| png, jpg, jpeg | mp4, gif |

### 🎥 Play video

When converting videos, an `.ascii` text file is generated. To play the video, run that generated file:

```
ascii -i path/to/file.ascii
```