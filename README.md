# Majokko

`majokko` is intended to be a lightweight ImageMagick `convert` and
`identify` alternative written in Go. The underlying library, found
in the `henshin/` directory, can be reused by other programs.

**This software is still in alpha, and as such, we do not commit to
providing a stable API.**

## Supported Formats

```
+--- (D)ecode
|+-- (E)ncode
||+- (M)etadata wrangling
|||  == Format ==
DE-  JPEG
DEM  PNG
DEM  ZNG (Zstd PNG)
DE-  GIF[1]
DEM  NetPBM (PPM, PGM, etc.)
DE-  JPEG XL[2][3]
DE-  QOI (Quite OK Image format)
D--  WEBP

[1] Animated GIFs are not yet supported.
[2] Requires external `cjxl` and `djxl` binaries. Enable with the
`--enable-external-codecs` option.
[3] Format not yet supported.
```

## License and Copyright Notice

Copyright &copy; 2022-2023 Ronsor Labs.

This software is provided to you under the terms of the MIT license.
See the included `LICENSE` file for more information.
