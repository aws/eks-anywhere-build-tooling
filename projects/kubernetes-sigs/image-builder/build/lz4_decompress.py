#!/usr/bin/env python3
# Decompress an lz4 frame file with the system liblz4, for hosts without the lz4 CLI.
import ctypes
import ctypes.util
import sys

LZ4F_VERSION = 100


def main(src, dst):
    lib = ctypes.CDLL(ctypes.util.find_library("lz4") or "liblz4.so.1")
    lib.LZ4F_isError.restype = ctypes.c_uint
    lib.LZ4F_getErrorName.restype = ctypes.c_char_p
    lib.LZ4F_decompress.restype = ctypes.c_size_t
    lib.LZ4F_decompress.argtypes = [
        ctypes.c_void_p, ctypes.c_void_p, ctypes.POINTER(ctypes.c_size_t),
        ctypes.c_void_p, ctypes.POINTER(ctypes.c_size_t), ctypes.c_void_p]

    def check(code):
        if lib.LZ4F_isError(ctypes.c_size_t(code)):
            sys.exit("lz4: " + lib.LZ4F_getErrorName(ctypes.c_size_t(code)).decode())
        return code

    dctx = ctypes.c_void_p()
    check(lib.LZ4F_createDecompressionContext(ctypes.byref(dctx), LZ4F_VERSION))
    out_buf = ctypes.create_string_buffer(4 << 20)
    hint = 1
    with open(src, "rb") as fin, open(dst, "wb") as fout:
        while True:
            data = fin.read(1 << 20)
            if not data:
                break
            pos = 0
            while pos < len(data):
                chunk = data[pos:]
                in_buf = ctypes.create_string_buffer(chunk, len(chunk))
                in_size = ctypes.c_size_t(len(chunk))
                out_size = ctypes.c_size_t(len(out_buf))
                hint = check(lib.LZ4F_decompress(dctx, out_buf, ctypes.byref(out_size),
                                                 in_buf, ctypes.byref(in_size), None))
                fout.write(out_buf.raw[:out_size.value])
                pos += in_size.value
    lib.LZ4F_freeDecompressionContext(dctx)
    if hint != 0:
        sys.exit("lz4: truncated input " + src)


if __name__ == "__main__":
    if len(sys.argv) != 3:
        sys.exit("usage: lz4_decompress.py SRC DST")
    main(sys.argv[1], sys.argv[2])
