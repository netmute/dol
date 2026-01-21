# dol - dark or light

Detect dark/light mode on the CLI.

## Why?

Modern operating systems can automatically switch between dark and light mode. However, many CLI tools assume a fixed dark or light background unless told otherwise. This is an attempt to make that less painful.

The output is intentionally minimal; it just prints `dark` or `light`. This allows you to construct command lines with it. Examples:  
- `fzf --color=$(dol)`
- `difft --background $(dol) file1 file2`

## Installation

Either of these will work.

```sh
mise use -g github:netmute/dol
```

```sh
go install github.com/netmute/dol@latest
```

## How it works

- dol writes a `CSI ? 996 n` device status report (DSR) query to `/dev/tty`
  and expects a reply like `CSI ? 997 ; 1 n` (dark) or `CSI ? 997 ; 2 n` (light).
- If your terminal does not support this query, dol prints an error and exits
  with status 1.

## Terminal support

See [this document](https://github.com/contour-terminal/contour/blob/master/docs/vt-extensions/color-palette-update-notifications.md#adoption-state) to verify if your terminal currently supports this feature.
