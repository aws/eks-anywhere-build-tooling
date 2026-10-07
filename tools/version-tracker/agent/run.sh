#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BOOTSTRAP_DIR="${PATCH_FIXER_BOOTSTRAP_DIR:-/tmp/patch-fixer-bootstrap}"
DEPENDENCY_DIR="${PATCH_FIXER_DEPENDENCY_DIR:-/tmp/patch-fixer-python}"
PYTHON_INSTALL_DIR="${PATCH_FIXER_PYTHON_INSTALL_DIR:-/tmp/patch-fixer-python-installations}"
PYTHON_VERSION="${PATCH_FIXER_PYTHON_VERSION:-3.11.13}"
UV_CACHE_DIR="${PATCH_FIXER_UV_CACHE_DIR:-/tmp/patch-fixer-uv-cache}"
UV_BIN="$BOOTSTRAP_DIR/bin/uv"

if ! "$UV_BIN" --version >/dev/null 2>&1; then
    rm -rf "$BOOTSTRAP_DIR"
    timeout "${PATCH_FIXER_BOOTSTRAP_TIMEOUT_SECONDS:-180}" python3 -m pip install \
        --disable-pip-version-check \
        --no-input \
        --only-binary=:all: \
        --quiet \
        --require-hashes \
        --target "$BOOTSTRAP_DIR" \
        -r "$SCRIPT_DIR/bootstrap-requirements.txt"
fi

find_python() {
    UV_CACHE_DIR="$UV_CACHE_DIR" UV_PYTHON_INSTALL_DIR="$PYTHON_INSTALL_DIR" \
        "$UV_BIN" python find "$PYTHON_VERSION" \
        --managed-python \
        --no-python-downloads \
        --no-project \
        2>/dev/null
}

PYTHON_BIN="$(find_python || true)"
if [[ -z "$PYTHON_BIN" ]]; then
    timeout "${PATCH_FIXER_PYTHON_TIMEOUT_SECONDS:-300}" \
        "$UV_BIN" python install "$PYTHON_VERSION" \
        --install-dir "$PYTHON_INSTALL_DIR" \
        --cache-dir "$UV_CACHE_DIR" \
        --no-bin \
        --no-progress
    PYTHON_BIN="$(find_python)"
fi

if ! PYTHONPATH="$DEPENDENCY_DIR${PYTHONPATH:+:$PYTHONPATH}" "$PYTHON_BIN" -c \
    'import importlib.metadata as m; assert m.version("strands-agents") == "1.56.0"' \
    >/dev/null 2>&1; then
    rm -rf "$DEPENDENCY_DIR"
    UV_CACHE_DIR="$UV_CACHE_DIR" \
    timeout "${PATCH_FIXER_BOOTSTRAP_TIMEOUT_SECONDS:-180}" "$UV_BIN" pip install \
        --python "$PYTHON_BIN" \
        --target "$DEPENDENCY_DIR" \
        --only-binary=:all: \
        --require-hashes \
        --no-config \
        --no-progress \
        -r "$SCRIPT_DIR/requirements.txt" \
        >/dev/null
fi

export PYTHONPATH="$DEPENDENCY_DIR${PYTHONPATH:+:$PYTHONPATH}"
exec "$PYTHON_BIN" "$SCRIPT_DIR/runner.py" "$@"
