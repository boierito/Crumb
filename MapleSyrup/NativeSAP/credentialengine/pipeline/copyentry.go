package pipeline

import (
    "archive/zip"
    "io"
)

// Keep untouched IPA archive entries byte-identical; do not decompress the
// original application assets or increase memory usage during repackaging.
func copyEntry(src *zip.File, dst *zip.Writer) error {
    in, err := src.OpenRaw()
    if err != nil { return err }
    header := src.FileHeader
    out, err := dst.CreateRaw(&header)
    if err != nil { return err }
    _, err = io.Copy(out, in)
    return err
}
