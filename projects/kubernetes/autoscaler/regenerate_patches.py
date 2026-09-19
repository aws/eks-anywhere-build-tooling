#!/usr/bin/env python3

import argparse
import json
import os
import shutil
import subprocess
import tempfile
from dataclasses import dataclass
from pathlib import Path
from typing import Optional


@dataclass(frozen=True)
class PatchMetadata:
    author_name: str
    author_email: str
    author_date: str
    subject: str
    body: str


def run(command: list[str], cwd: Path, env: Optional[dict[str, str]] = None, check: bool = True):
    completed = subprocess.run(
        command,
        cwd=cwd,
        env=env,
        text=True,
        capture_output=True,
        check=False,
    )
    if check and completed.returncode != 0:
        raise RuntimeError(
            f"command failed ({' '.join(command)}):\n{completed.stdout}\n{completed.stderr}".strip()
        )
    return completed


def parse_patch_metadata(path: Path) -> PatchMetadata:
    patch_text = path.read_text(encoding="utf-8")
    with tempfile.TemporaryDirectory() as temp_dir:
        message_path = Path(temp_dir) / "message"
        patch_path = Path(temp_dir) / "patch"
        completed = subprocess.run(
            ["git", "mailinfo", str(message_path), str(patch_path)],
            input=patch_text,
            text=True,
            capture_output=True,
            check=False,
        )
        if completed.returncode != 0:
            raise ValueError(f"unable to parse patch metadata in {path}: {completed.stderr.strip()}")
        fields = {}
        for line in completed.stdout.splitlines():
            key, separator, value = line.partition(": ")
            if separator:
                fields[key] = value
        body = message_path.read_text(encoding="utf-8").strip()
    author_name = fields.get("Author", "")
    author_email = fields.get("Email", "")
    author_date = fields.get("Date", "")
    subject = fields.get("Subject", "")
    if not all((author_name, author_email, author_date, subject)):
        raise ValueError(f"incomplete patch metadata in {path}")
    return PatchMetadata(author_name, author_email, author_date, subject, body)


def commit_and_format(
    repo: Path,
    metadata: PatchMetadata,
    output: Path,
    paths: Optional[list[str]] = None,
) -> None:
    add_command = ["git", "add", "--all"] if paths is None else ["git", "add", "--", *paths]
    status_command = ["git", "status", "--porcelain"]
    if paths is not None:
        status_command.extend(["--", *paths])
    run(add_command, cwd=repo)
    if not run(status_command, cwd=repo).stdout.strip():
        return
    message = metadata.subject
    if metadata.body:
        message += f"\n\n{metadata.body}"
    message_path = repo / ".git" / "patch-fixer-message"
    message_path.write_text(message + "\n", encoding="utf-8")
    env = os.environ.copy()
    env.update(
        {
            "GIT_AUTHOR_NAME": metadata.author_name,
            "GIT_AUTHOR_EMAIL": metadata.author_email,
            "GIT_AUTHOR_DATE": metadata.author_date,
            "GIT_COMMITTER_NAME": "EKS Distro PR Bot",
            "GIT_COMMITTER_EMAIL": "aws-model-rocket-bots+eksdistroprbot@amazon.com",
        }
    )
    run(["git", "commit", "--file", str(message_path)], cwd=repo, env=env)
    patch = run(["git", "format-patch", "-1", "--stdout", "--no-signature"], cwd=repo).stdout
    output.write_text(patch, encoding="utf-8")


def prune_provider_files(directory: Path, prefix: str, keep: set[str]) -> None:
    for path in directory.glob(f"{prefix}*.go"):
        if path.name not in keep:
            path.unlink()


def transform_router(source: Path) -> bool:
    router = source / "cluster-autoscaler" / "cloudprovider" / "router"
    router_all = router / "router_all.go"
    if not router_all.exists():
        return False
    prune_provider_files(router, "router_", {"router_all.go", "router_clusterapi.go"})
    lines = router_all.read_text(encoding="utf-8").splitlines()
    transformed = []
    for line in lines:
        stripped = line.strip()
        if (
            stripped.startswith('_ "')
            and "/cloudprovider/" in stripped
            and "/clusterapi\"" not in stripped
        ):
            continue
        if "builder.SetDefaultCloudProvider(" in line:
            indentation = line[: len(line) - len(line.lstrip())]
            line = f"{indentation}builder.SetDefaultCloudProvider(cloudprovider.ClusterAPIProviderName)"
        transformed.append(line)
    router_all.write_text("\n".join(transformed) + "\n", encoding="utf-8")
    return True


def transform_builder(source: Path) -> bool:
    builder = source / "cluster-autoscaler" / "cloudprovider" / "builder"
    builder_all = builder / "builder_all.go"
    if not builder_all.exists():
        return False
    prune_provider_files(
        builder,
        "builder_",
        {"builder_all.go", "builder_clusterapi.go"},
    )
    lines = builder_all.read_text(encoding="utf-8").splitlines()
    transformed = []
    in_import = False
    in_providers = False
    in_switch = False
    keep_case = True
    switch_depth = 0
    for line in lines:
        stripped = line.strip()
        if stripped == "import (":
            in_import = True
            transformed.append(line)
            continue
        if in_import and stripped == ")":
            in_import = False
            transformed.append(line)
            continue
        if in_import:
            if (
                "/cloudprovider/" in stripped
                and "/clusterapi\"" not in stripped
            ):
                continue
            transformed.append(line)
            continue

        if strings_contains_any(line, ["var AvailableCloudProviders"]):
            in_providers = True
            transformed.append(line)
            continue
        if in_providers:
            if stripped == "}":
                in_providers = False
                transformed.append(line)
            elif "ClusterAPIProviderName" in line or "ProviderName" not in line:
                transformed.append(line)
            continue

        if stripped.startswith("const DefaultCloudProvider ="):
            indentation = line[: len(line) - len(line.lstrip())]
            transformed.append(
                f"{indentation}const DefaultCloudProvider = cloudprovider.ClusterAPIProviderName"
            )
            continue

        if "switch opts.CloudProviderName" in line:
            in_switch = True
            switch_depth = line.count("{") - line.count("}")
            keep_case = True
            transformed.append(line)
            continue
        if in_switch:
            if stripped.startswith("case cloudprovider."):
                keep_case = "ClusterAPIProviderName" in line
            switch_depth += line.count("{") - line.count("}")
            if switch_depth <= 0:
                in_switch = False
                transformed.append(line)
                continue
            if keep_case:
                transformed.append(line)
            continue
        transformed.append(line)
    builder_all.write_text("\n".join(transformed) + "\n", encoding="utf-8")
    return True


def strings_contains_any(value: str, needles: list[str]) -> bool:
    return any(needle in value for needle in needles)


def remove_gce_dependencies(source: Path) -> bool:
    changed = False
    for path in (source / "cluster-autoscaler").rglob("*.go"):
        original = path.read_text(encoding="utf-8")
        if "LocalSSDDiskSizeProvider" not in original and "localssdsize" not in original:
            continue
        content = original.replace("opts.GCEOptions.LocalSSDDiskSizeProvider", "nil")
        lines = [
            line
            for line in content.splitlines()
            if "localssdsize" not in line
            and "LocalSSDDiskSizeProvider" not in line
        ]
        content = "\n".join(lines) + "\n"
        if content != original:
            path.write_text(content, encoding="utf-8")
            changed = True
    return changed


def remove_cloud_provider_directories(source: Path) -> None:
    cloudprovider = source / "cluster-autoscaler" / "cloudprovider"
    keep = {"builder", "router", "clusterapi", "mocks", "test"}
    for path in cloudprovider.iterdir():
        if path.is_dir() and path.name not in keep:
            shutil.rmtree(path)


def go_binary(release_branch: str, project_root: Path) -> str:
    version_file = project_root / release_branch / "GOLANG_VERSION"
    if version_file.exists():
        version = version_file.read_text(encoding="utf-8").strip()
        versioned_go = Path(f"/go/go{version}/bin/go")
        if versioned_go.exists():
            return str(versioned_go)
    return "go"


def regenerate(args) -> list[Path]:
    source = Path(args.source).resolve()
    patches_dir = Path(args.patches).resolve()
    output_dir = Path(args.output).resolve()
    output_dir.mkdir(parents=True, exist_ok=True)
    existing = sorted(patches_dir.glob("*.patch"))
    if len(existing) < 2:
        raise ValueError("autoscaler requires at least provider and go.mod patches")

    run(["git", "am", "--abort"], cwd=source, check=False)
    run(["git", "checkout", "--force", args.revision], cwd=source)
    run(["git", "clean", "-fdx"], cwd=source)
    run(["git", "config", "user.name", "Prow Bot"], cwd=source)
    run(["git", "config", "user.email", "prow@amazonaws.com"], cwd=source)

    first_patch = existing[0]
    if not transform_router(source) and not transform_builder(source):
        raise ValueError("unsupported autoscaler source: no router_all.go or builder_all.go")
    generated = [output_dir / first_patch.name]
    commit_and_format(source, parse_patch_metadata(first_patch), generated[-1])

    for static_patch in existing[1:-1]:
        if "gce" in static_patch.name.lower():
            if remove_gce_dependencies(source):
                generated.append(output_dir / static_patch.name)
                commit_and_format(source, parse_patch_metadata(static_patch), generated[-1])
        else:
            run(["git", "am", "--committer-date-is-author-date", str(static_patch)], cwd=source)
            shutil.copy2(static_patch, output_dir / static_patch.name)
            generated.append(output_dir / static_patch.name)

    dependency_patch = existing[-1]
    cluster_autoscaler = source / "cluster-autoscaler"
    remove_cloud_provider_directories(source)
    run([go_binary(args.release_branch, patches_dir.parent.parent), "mod", "tidy"], cwd=cluster_autoscaler)
    if run(["git", "status", "--porcelain", "--", "cluster-autoscaler/go.mod", "cluster-autoscaler/go.sum"], cwd=source).stdout.strip():
        generated.append(output_dir / dependency_patch.name)
        commit_and_format(
            source,
            parse_patch_metadata(dependency_patch),
            generated[-1],
            ["cluster-autoscaler/go.mod", "cluster-autoscaler/go.sum"],
        )

    run(["git", "reset", "--hard", args.revision], cwd=source)
    run(["git", "clean", "-fdx"], cwd=source)
    for patch in generated:
        run(["git", "am", "--committer-date-is-author-date", str(patch)], cwd=source)
    return generated


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--source", required=True)
    parser.add_argument("--patches", required=True)
    parser.add_argument("--output", required=True)
    parser.add_argument("--manifest", required=True)
    parser.add_argument("--revision", required=True)
    parser.add_argument("--release-branch", required=True)
    args = parser.parse_args()

    patches = regenerate(args)
    Path(args.manifest).write_text(
        json.dumps({"patches": [str(path) for path in patches]}, indent=2),
        encoding="utf-8",
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
