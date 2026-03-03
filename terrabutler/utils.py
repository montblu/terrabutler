from terrabutler.requirements import check_requirements
from sys import exit
from semantic_version import Version
from os import getenv
import os
import subprocess
import signal

check_requirements()
ROOT_PATH = getenv("TERRABUTLER_ROOT")

paths = {
    "backends": ROOT_PATH + "/configs/backends",
    "environment": ROOT_PATH + "/site_inception/.terraform/environment",
    "inception": ROOT_PATH + "/site_inception",
    "root": ROOT_PATH,
    "settings": ROOT_PATH + "/configs/settings.yml",
    "templates": ROOT_PATH + "/configs/templates",
    "variables": ROOT_PATH + "/configs/variables"
}


def is_semantic_version(version):
    """
    Check if the version corresponds to the semantic versioning.
    """
    try:
        Version(version)
    except ValueError:
        return False
    except Exception as e:
        print(f"There was an error while parsing version: {e}")
        exit(1)
    return True


def run_subprocess(command, cwd=None, env=None, shell=False,
                   check=False, stdout=None, stderr=None,
                   capture_output=False, text=False, input=None,
                   use_popen=False, handle_sigint=False):
    """
    Run subprocess with LD_LIBRARY_PATH cleaned for PyInstaller compatibility.
    Restores original library path so system binaries use system libraries.

    Args:
        command: Command to run (string if shell=True, list otherwise)
        cwd: Working directory
        env: Environment variables dict (will be merged with cleaned env)
        shell: Whether to run command through shell
        check: Raise exception on non-zero exit code
        stdout/stderr: Redirect streams
        capture_output: Capture stdout/stderr
        text: Return text instead of bytes
        input: Input to pass to command
        use_popen: Use Popen instead of run (for streaming/interactive)
        handle_sigint: Ignore SIGINT during execution (for terraform)

    Returns:
        CompletedProcess object or Popen object if use_popen=True
    """
    # Clean environment
    clean_env = dict(os.environ)
    lp_key = 'LD_LIBRARY_PATH'
    lp_orig = clean_env.get(lp_key + '_ORIG')

    if lp_orig is not None:
        clean_env[lp_key] = lp_orig
    else:
        clean_env.pop(lp_key, None)

    # Merge user-provided env vars
    if env:
        clean_env.update(env)

    if use_popen:
        if handle_sigint:
            prev_sigint_handler = signal.getsignal(signal.SIGINT)
            try:
                p = subprocess.Popen(
                    args=command,
                    cwd=cwd,
                    env=clean_env,
                    stdout=stdout,
                    stderr=stderr
                )
                signal.signal(signal.SIGINT, signal.SIG_IGN)
                p.wait()
                return p
            finally:
                signal.signal(signal.SIGINT, prev_sigint_handler)
        else:
            return subprocess.Popen(
                args=command,
                cwd=cwd,
                env=clean_env,
                stdout=stdout,
                stderr=stderr
            )
    else:
        return subprocess.run(
            command,
            cwd=cwd,
            env=clean_env,
            shell=shell,
            check=check,
            stdout=stdout,
            stderr=stderr,
            capture_output=capture_output,
            text=text,
            input=input
        )
