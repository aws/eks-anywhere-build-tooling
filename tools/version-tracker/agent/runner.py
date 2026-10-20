#!/usr/bin/env python3

import argparse
import dataclasses
import datetime
import json
import os
import re
import sys
from pathlib import Path
from typing import Callable


SCHEMA_VERSION = 1
MAX_READ_LINES = 400
MAX_SEARCH_RESULTS = 200
MAX_TOOL_OUTPUT = 50000
REDACTED_TRACE_KEYS = ("reasoning", "thinking", "signature", "redacted")


def sanitize_trace_value(value):
    if dataclasses.is_dataclass(value):
        value = dataclasses.asdict(value)
    if isinstance(value, dict):
        return {
            str(key): sanitize_trace_value(item)
            for key, item in value.items()
            if not any(marker in str(key).lower() for marker in REDACTED_TRACE_KEYS)
        }
    if isinstance(value, (list, tuple)):
        return [sanitize_trace_value(item) for item in value]
    if isinstance(value, (str, int, float, bool)) or value is None:
        return value
    if hasattr(value, "stop_reason"):
        return {
            "stop_reason": str(getattr(value, "stop_reason", "")),
            "message": sanitize_trace_value(getattr(value, "message", None)),
            "metrics": sanitize_trace_value(getattr(value, "metrics", None)),
        }
    return str(value)


class TraceRecorder:
    def __init__(self, path: Path):
        self.path = path
        self.path.parent.mkdir(parents=True, exist_ok=True)

    def __call__(self, **kwargs) -> None:
        visible = {
            key: sanitize_trace_value(value)
            for key, value in kwargs.items()
            if not any(marker in key.lower() for marker in REDACTED_TRACE_KEYS)
        }
        record = {
            "timestamp": datetime.datetime.now(datetime.timezone.utc).isoformat(),
            "event": visible,
        }
        with self.path.open("a", encoding="utf-8") as trace:
            trace.write(json.dumps(record, sort_keys=True) + "\n")


class Workspace:
    def __init__(
        self,
        root: Path,
        allowed_files: set[str],
        original_patch: str,
    ):
        self.root = root.resolve()
        self.allowed_files = allowed_files
        self.original_patch = original_patch

    def resolve(
        self,
        relative_path: str,
        require_allowed: bool = False,
        require_existing: bool = True,
    ) -> Path:
        path = Path(relative_path)
        if path.is_absolute() or ".." in path.parts:
            raise ValueError(f"path escapes workspace: {relative_path}")
        candidate = self.root / path
        normalized = candidate.relative_to(self.root).as_posix()
        if require_allowed and not self.is_edit_allowed(normalized):
            raise ValueError(f"editing {normalized} is not allowed")
        if require_allowed and candidate.is_symlink():
            raise ValueError(f"editing symlink {normalized} is not allowed")

        check_path = candidate if require_existing else candidate.parent
        while not check_path.exists():
            parent = check_path.parent
            if parent == check_path:
                raise ValueError(f"no existing parent for {relative_path}")
            check_path = parent
        resolved = check_path.resolve()
        if resolved != self.root and self.root not in resolved.parents:
            raise ValueError(f"path escapes workspace: {relative_path}")
        return candidate

    def is_edit_allowed(self, normalized: str) -> bool:
        return normalized in self.allowed_files

    def list_files(self, pattern: str = "**/*") -> str:
        pattern_path = Path(pattern)
        if pattern_path.is_absolute() or ".." in pattern_path.parts:
            raise ValueError(f"invalid file pattern: {pattern}")
        files = []
        for path in self.root.glob(pattern):
            if not path.is_file() or ".git" in path.parts:
                continue
            resolved = path.resolve()
            if resolved != self.root and self.root not in resolved.parents:
                continue
            files.append(path.relative_to(self.root).as_posix())
            if len(files) >= MAX_SEARCH_RESULTS:
                break
        return "\n".join(sorted(files))

    def read_file(self, path: str, start_line: int = 1, end_line: int = 200) -> str:
        target = self.resolve(path)
        lines = target.read_text(encoding="utf-8").splitlines()
        start = max(start_line, 1) - 1
        end = min(max(end_line, start_line), start + MAX_READ_LINES, len(lines))
        return "\n".join(f"{index + 1}: {lines[index]}" for index in range(start, end))

    def search_repository(self, query: str, path: str = ".") -> str:
        target = self.resolve(path)
        pattern = re.compile(query)
        matches = []
        candidates = [target] if target.is_file() else target.rglob("*")
        for candidate in candidates:
            if not candidate.is_file() or ".git" in candidate.parts:
                continue
            resolved = candidate.resolve()
            if resolved != self.root and self.root not in resolved.parents:
                continue
            if candidate.stat().st_size > 2 * 1024 * 1024:
                continue
            try:
                lines = candidate.read_text(encoding="utf-8").splitlines()
            except UnicodeDecodeError:
                continue
            for line_number, line in enumerate(lines, 1):
                if pattern.search(line):
                    relative = candidate.relative_to(self.root).as_posix()
                    matches.append(f"{relative}:{line_number}:{line}")
                    if len(matches) >= MAX_SEARCH_RESULTS:
                        return "\n".join(matches)
        return "\n".join(matches)

    def replace_text(self, path: str, old_text: str, new_text: str) -> str:
        target = self.resolve(path, require_allowed=True)
        content = target.read_text(encoding="utf-8")
        count = content.count(old_text)
        if count != 1:
            raise ValueError(f"expected exactly one match in {path}, found {count}")
        target.write_text(content.replace(old_text, new_text, 1), encoding="utf-8")
        return f"updated {path}"

    def replace_lines(self, path: str, start_line: int, end_line: int, new_text: str) -> str:
        target = self.resolve(path, require_allowed=True)
        content = target.read_text(encoding="utf-8")
        lines = content.splitlines()
        if start_line < 1 or end_line < start_line or end_line > len(lines):
            raise ValueError(
                f"invalid line range for {path}: {start_line}-{end_line}, file has {len(lines)} lines"
            )
        replacement = new_text.splitlines()
        updated = lines[: start_line - 1] + replacement + lines[end_line:]
        suffix = "\n" if content.endswith("\n") else ""
        target.write_text("\n".join(updated) + suffix, encoding="utf-8")
        return f"replaced lines {start_line}-{end_line} in {path}"

    def write_file(self, path: str, content: str) -> str:
        if len(content.encode("utf-8")) > 2 * 1024 * 1024:
            raise ValueError("file content exceeds 2 MiB limit")
        target = self.resolve(path, require_allowed=True, require_existing=False)
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_text(content, encoding="utf-8")
        return f"wrote {path}"

    def delete_file(self, path: str) -> str:
        target = self.resolve(path, require_allowed=True)
        target.unlink()
        return f"deleted {path}"

    def read_original_patch(self, start_line: int = 1, end_line: int = 200) -> str:
        lines = self.original_patch.splitlines()
        start = max(start_line, 1) - 1
        end = min(max(end_line, start_line), start + MAX_READ_LINES, len(lines))
        return "\n".join(f"{index + 1}: {lines[index]}" for index in range(start, end))

    def show_diff(self) -> str:
        import subprocess

        completed = subprocess.run(
            ["git", "diff", "--"],
            cwd=self.root,
            text=True,
            capture_output=True,
            check=False,
        )
        if completed.returncode != 0:
            raise RuntimeError(completed.stderr.strip() or "git diff failed")
        return completed.stdout[:MAX_TOOL_OUTPUT]


def build_model(model_id: str, region: str):
    from strands.models import BedrockModel

    return BedrockModel(model_id=model_id, region_name=region, max_tokens=16384)


def build_agent(
    workspace: Workspace,
    model_id: str,
    region: str,
    system_prompt: str,
    trace_path: str,
):
    from strands import Agent, tool

    @tool
    def list_files(pattern: str = "**/*") -> str:
        """List repository files matching a glob pattern."""
        return workspace.list_files(pattern)

    @tool
    def read_file(path: str, start_line: int = 1, end_line: int = 200) -> str:
        """Read a bounded line range from a repository file."""
        return workspace.read_file(path, start_line, end_line)

    @tool
    def search_repository(query: str, path: str = ".") -> str:
        """Search repository text using a regular expression."""
        return workspace.search_repository(query, path)

    @tool
    def read_original_patch(start_line: int = 1, end_line: int = 200) -> str:
        """Read a bounded line range from the original failed mail patch."""
        return workspace.read_original_patch(start_line, end_line)

    @tool
    def replace_text(path: str, old_text: str, new_text: str) -> str:
        """Replace one exact text occurrence in an allowed file."""
        return workspace.replace_text(path, old_text, new_text)

    @tool
    def replace_lines(path: str, start_line: int, end_line: int, new_text: str) -> str:
        """Replace an inclusive line range in an allowed file."""
        return workspace.replace_lines(path, start_line, end_line, new_text)

    @tool
    def write_file(path: str, content: str) -> str:
        """Write complete UTF-8 content to an allowed file."""
        return workspace.write_file(path, content)

    @tool
    def delete_file(path: str) -> str:
        """Delete an allowed file."""
        return workspace.delete_file(path)

    @tool
    def show_diff() -> str:
        """Show current uncommitted changes."""
        return workspace.show_diff()

    callback_handler = TraceRecorder(Path(trace_path)) if trace_path else None
    return Agent(
        model=build_model(model_id, region),
        system_prompt=system_prompt,
        tools=[
            list_files,
            read_file,
            search_repository,
            read_original_patch,
            replace_text,
            replace_lines,
            write_file,
            delete_file,
            show_diff,
        ],
        callback_handler=callback_handler,
        context_manager="auto",
    )


def execute(request: dict, result_path: Path, agent_factory: Callable = build_agent) -> dict:
    if request.get("schema_version") != SCHEMA_VERSION:
        raise ValueError("unsupported request schema version")

    workspace = Workspace(
        Path(request["workspace_path"]),
        set(request.get("allowed_files", [])),
        request.get("original_patch", ""),
    )
    run_dir = Path(os.environ.get("PATCH_FIXER_RUN_DIR", result_path.parent)).resolve()
    trace_path = str(run_dir / "trace.jsonl") if request.get("diagnostics") else ""
    agent = agent_factory(
        workspace,
        request["model_id"],
        request["region"],
        request["system_prompt"],
        trace_path,
    )
    result = agent(
        request["prompt"],
        limits={
            "turns": request["max_turns"],
            "total_tokens": request["max_total_tokens"],
        },
    )

    metrics = getattr(result, "metrics", None)
    usage = getattr(metrics, "accumulated_usage", {}) if metrics else {}
    response = {
        "schema_version": SCHEMA_VERSION,
        "status": "completed",
        "stop_reason": str(getattr(result, "stop_reason", "")),
        "summary": str(result)[:2000],
        "turns": getattr(metrics, "cycle_count", 0) if metrics else 0,
        "total_tokens": usage.get("totalTokens", 0) if isinstance(usage, dict) else 0,
        "message": sanitize_trace_value(getattr(result, "message", None)),
    }
    result_path.write_text(json.dumps(response, indent=2), encoding="utf-8")
    return response


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--request", required=True)
    parser.add_argument("--result", required=True)
    args = parser.parse_args()

    result_path = Path(args.result)
    try:
        request = json.loads(Path(args.request).read_text(encoding="utf-8"))
        execute(request, result_path)
        return 0
    except Exception as error:
        result_path.write_text(
            json.dumps(
                {
                    "schema_version": SCHEMA_VERSION,
                    "status": "runner_failed",
                    "summary": str(error),
                },
                indent=2,
            ),
            encoding="utf-8",
        )
        print(f"patch fixer Strands runner failed: {error}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
