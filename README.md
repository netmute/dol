# dol - dark or light

Query your terminal for its dark/light appearance.

## Usage

```sh
# System uses dark appearance.
> dol
dark
```

```sh
# System uses light appearance.
> dol
light
```

## Installation

```sh
# With mise
mise use -g github:netmute/dol
```

```sh
# With go
go install github.com/netmute/dol@latest
```

## Notes

- dol writes a `CSI ? 996 n` device status report (DSR) query to `/dev/tty`
  and expects a reply like `CSI ? 997 ; 1 n` (dark) or `CSI ? 997 ; 2 n` (light).
- If your terminal does not support this query, dol prints an error and exits
  with status 1.

## Terminal support

See [this document](https://github.com/contour-terminal/contour/blob/master/docs/vt-extensions/color-palette-update-notifications.md#adoption-state) to verify if your terminal currently supports this feature.
