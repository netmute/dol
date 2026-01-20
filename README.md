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

## Notes

- The tool writes the `CSI ? 996 n` device status report (DSR) query to `/dev/tty`
  and expects a reply like `CSI ? 997 ; 1 n` (dark) or `CSI ? 997 ; 2 n` (light).
- If the terminal does not support this query, `dol` prints an error and exits
  with status 1.
