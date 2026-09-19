import importlib.util
import json
import os
import subprocess
import tempfile
import unittest
from pathlib import Path


MODULE_PATH = Path(__file__).with_name("runner.py")
SPEC = importlib.util.spec_from_file_location("patch_fixer_runner", MODULE_PATH)
runner = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(runner)


class RunnerTest(unittest.TestCase):
    def test_workspace_blocks_path_escape_and_unlisted_edits(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            (root / "main.go").write_text("old\n", encoding="utf-8")
            workspace = runner.Workspace(root, {"main.go"}, "patch")

            with self.assertRaises(ValueError):
                workspace.read_file("../secret")
            with self.assertRaises(ValueError):
                workspace.list_files("../*")
            with self.assertRaises(ValueError):
                workspace.replace_text("other.go", "a", "b")

    def test_workspace_blocks_symlink_escape(self):
        with tempfile.TemporaryDirectory() as temp_dir, tempfile.TemporaryDirectory() as outside:
            root = Path(temp_dir)
            outside_path = Path(outside)
            (outside_path / "secret").write_text("secret", encoding="utf-8")
            (root / "linked").symlink_to(outside_path, target_is_directory=True)
            workspace = runner.Workspace(root, {"linked/secret"}, "patch")

            with self.assertRaises(ValueError):
                workspace.read_file("linked/secret")
            with self.assertRaises(ValueError):
                workspace.write_file("linked/new", "data")

    def test_workspace_replaces_line_range(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            path = root / "main.go"
            path.write_text("one\ntwo\nthree\nfour\n", encoding="utf-8")
            workspace = runner.Workspace(root, {"main.go"}, "patch")

            result = workspace.replace_lines("main.go", 2, 3, "new-two\nnew-three")

            self.assertEqual(result, "replaced lines 2-3 in main.go")
            self.assertEqual(
                path.read_text(encoding="utf-8"),
                "one\nnew-two\nnew-three\nfour\n",
            )

    def test_workspace_rejects_sibling_file_below_allowed_parent(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            (root / "pkg").mkdir()
            (root / "pkg" / "allowed.go").write_text("allowed\n", encoding="utf-8")
            (root / "pkg" / "sibling.go").write_text("sibling\n", encoding="utf-8")
            workspace = runner.Workspace(root, {"pkg/allowed.go"}, "patch")

            with self.assertRaises(ValueError):
                workspace.replace_text("pkg/sibling.go", "sibling", "changed")

    def test_trace_recorder_excludes_reasoning_fields(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            path = Path(temp_dir) / "trace.jsonl"
            recorder = runner.TraceRecorder(path)

            recorder(
                data="visible response",
                reasoningText="hidden reasoning",
                event={
                    "contentBlockStart": {
                        "start": {
                            "toolUse": {
                                "name": "read_file",
                                "input": {"path": "main.go"},
                            }
                        }
                    },
                    "reasoningContent": "hidden content",
                },
            )

            trace = path.read_text(encoding="utf-8")
            self.assertIn("visible response", trace)
            self.assertIn("read_file", trace)
            self.assertNotIn("hidden reasoning", trace)
            self.assertNotIn("hidden content", trace)

    def test_execute_uses_prepared_workspace_without_importing_strands(self):
        captured = {}

        class FakeMetrics:
            cycle_count = 2
            accumulated_usage = {"totalTokens": 38}

        class FakeResult:
            stop_reason = "end_turn"
            metrics = FakeMetrics()
            message = {"content": [{"text": "done"}]}

            def __str__(self):
                return "done"

        class FakeAgent:
            def __init__(self, workspace):
                self.workspace = workspace

            def __call__(self, prompt, limits):
                self.workspace.replace_text("main.go", "current\n", "new\n")
                captured["prompt"] = prompt
                captured["limits"] = limits
                return FakeResult()

        def factory(workspace, model, region, system_prompt, trace_path):
            captured["model"] = model
            captured["region"] = region
            captured["system_prompt"] = system_prompt
            captured["trace_path"] = trace_path
            return FakeAgent(workspace)

        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            (root / "main.go").write_text("current\n", encoding="utf-8")
            subprocess.run(["git", "init", "-q"], cwd=root, check=True)
            result_path = root / "agent-result.json"

            response = runner.execute(
                {
                    "schema_version": 1,
                    "workspace_path": str(root),
                    "allowed_files": ["main.go"],
                    "original_patch": "patch",
                    "system_prompt": "system",
                    "prompt": "repair",
                    "model_id": "test-model",
                    "region": "us-test-1",
                    "max_turns": 4,
                    "max_total_tokens": 1000,
                    "diagnostics": False,
                },
                result_path,
                agent_factory=factory,
            )

            self.assertEqual(response["status"], "completed")
            self.assertEqual(response["stop_reason"], "end_turn")
            self.assertEqual(response["turns"], 2)
            self.assertEqual(response["total_tokens"], 38)
            self.assertEqual((root / "main.go").read_text(encoding="utf-8"), "new\n")
            self.assertEqual(captured["model"], "test-model")
            self.assertEqual(captured["region"], "us-test-1")
            self.assertEqual(captured["system_prompt"], "system")
            self.assertEqual(captured["limits"], {"turns": 4, "total_tokens": 1000})
            self.assertEqual(json.loads(result_path.read_text())["summary"], "done")


if __name__ == "__main__":
    unittest.main()
