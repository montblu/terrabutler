---
name: terrabutler
description: >
  Read-only workflow and guardrails for Terrabutler-managed infrastructure repos (Terraform
  wrapper: *.tf, *.tfvars, *.hcl, *.tftpl across site_* directories). ALWAYS activate this skill
  when the user says they are in a Terrabutler project, mentions "terrabutler", asks to run
  `terrabutler` commands, or works in a repo containing configs/settings.yml with sites.ordered.
  Also use when reviewing, validating, formatting, or inspecting Terraform code for any site.
  Project-agnostic: sites, environments, auth, and backends are ALWAYS discovered at runtime,
  never assumed — different Terrabutler projects have different sites, clouds, and credentials.
---

# Terrabutler Workflow (project-agnostic)

Repos managed by **Terrabutler** use a wrapper that handles multi-environment, multi-site
orchestration over Terraform. **Never run raw `terraform` commands** for state-changing work —
always use `terrabutler tf -site <site> ...`.

> This skill ships with the `terrabutler` binary itself and is version-matched to whichever
> terrabutler version this repo pins (synced via `mise skills sync`). It hardcodes **nothing**
> project-specific: site lists, environments, cloud accounts, credentials, and backend names all
> vary per project and are discovered at runtime. A repo may still add its own project-specific
> skill under a different name (e.g. `<repo>/.claude/skills/<project>-notes`) for details this
> skill can't know — auth profiles, tenant IDs, bucket names — without colliding with the
> synced `terrabutler` skill.

---

## 0. Guardrails (READ FIRST) — STRICT ALLOWLIST

These repos control **real production and staging cloud infrastructure** (AWS, GCP, Azure, GKE).
This skill operates in **read-only mode**. It may run ONLY the commands on the allowlist below.
**Everything not on this list is forbidden — no exceptions, even if the user asks.**
If a task would require a forbidden command, STOP and tell the user to run it themselves manually.

### ✅ ALLOWED commands (the ONLY commands this skill may run)

```bash
terrabutler tf -site <site> validate     # Validate configuration files
terrabutler tf -site <site> fmt          # Reformat configuration (local files only)
terrabutler tf -site <site> show         # Show current state or a saved plan
terrabutler tf -site <site> output       # Show output values
terrabutler env show                     # Show current environment name
terrabutler env list                     # List available environments
mise run tf:summarize -- <site>          # Summarize a plan (read-only), if the project defines it
mise tasks                               # List the project's available mise tasks
```

Plus non-terrabutler read-only formatting/inspection: `terraform fmt`.
That is the complete list. Nothing else may be executed.

### 🚫 FORBIDDEN — never run these (do not run even if explicitly asked)

- `apply` — in any form
- `destroy` — in any form
- `--auto-approve`
- `--destroy` planning mode
- `plan` (including `--out`)
- `init`, `import`, `refresh`, `taint`, `untaint`, `force-unlock`, `console`, `providers`
- `state rm|mv|push|pull|replace-provider|list|show` (all state subcommands)
- `env select`, `env new`, `env delete`, `env reload`
- `mise run tf:all-sites-on-*` and any `mise login*` / `az:*` tasks
- Any command that mutates infrastructure, state, environments, or cloud auth

**If the workflow needs a forbidden command, respond with the exact command the user should run
manually and wait — do NOT execute it yourself.**

### Hard scope limits
- Touch only `*.tf`, `*.tfvars`, `*.hcl`, `*.tftpl`. Do not edit application code.
- Never commit or print decrypted secrets, credentials, or the contents of personal secret files.
- Never edit files under `.terraform/` or generated backend/state files by hand.

---

## 0.5 Enabling Terrabutler — check this FIRST if any command fails

Terrabutler only runs when `TERRABUTLER_ENABLE=true` is present in the environment. That
variable — together with `TERRABUTLER_ROOT`, `TERRABUTLER_ENV`, and the project's AWS/K8s config
paths — is injected by **mise** from the repo's `.mise.toml` when mise activates in the repo.

### Symptom
Any terrabutler command exits with:
> `Terrabutler is not currently enabled on this folder. Please set 'TERRABUTLER_ENABLE' in your environment to true to enable it.`

This means one of two things:
1. **The repo has never been initialized** (no `terrabutler init` has run yet), or
2. **mise's environment isn't loaded in the current shell** — e.g. a non-interactive shell (like
   an automated tool shell) where mise's shell-activation hook never fires.

### Diagnose which one — do this FIRST, before any other troubleshooting
Don't guess. A fresh clone has **no** generated working dirs, so check for them:
```bash
ls -d site_inception/.terraform 2>/dev/null && echo INITIALIZED || echo NOT-INITIALIZED
```
- `NOT-INITIALIZED` → the repo needs `terrabutler init` **before ANY terrabutler command will
  work** (see below). This is the common case on a fresh clone — assume it whenever the "not
  enabled" message appears on a repo you haven't confirmed is set up. Do not waste steps trying
  read-only commands first; they will all fail until init has run.
- `INITIALIZED` → the working dirs exist, so it's just mise's env not loaded in this shell. The
  user needs an mise-activated shell (interactive shell, or `eval "$(mise activate bash)"`).

### How a repo is initialized (USER runs this — the skill must NOT)
Run in this exact order — init depends on being logged in first:
```bash
mise login          # authenticate to the clouds the project uses (task name may vary — see `mise tasks`)
terrabutler init    # prepare working dirs and enable terrabutler for the repo — REQUIRED before anything else
```
Both mutate cloud auth and working state, so both are on the **FORBIDDEN** list for this read-only
skill (`mise login*`, `terrabutler init`). If you hit the "not enabled" message, **do not try
to work around it** — surface these two commands, ask the user to run them, and continue once the
repo is enabled.

### Do NOT fake it
Never export `TERRABUTLER_ENABLE=true` yourself to silence the message. That enables the binary
without the rest of the mise-provided environment (project root, current env, AWS profile,
`KUBECONFIG`), so any command you then run reports misleading results. The allowed commands below
assume an **mise-activated shell**; if yours isn't, the fix is for the user to init/activate the
repo, not to hand-set env vars.

---

## 1. Discovering Sites and Environments

### Locate the project config first
Terrabutler config lives in `configs/settings.yml`, **relative to the Terrabutler project root**.
That root is not always the repo root — in some repos it is nested (e.g. `infra/configs/settings.yml`).
Find it before anything else:

```bash
find . -path '*/configs/settings.yml' -not -path '*/.terraform/*'
```

### Sites
The authoritative list of sites — and their strict deploy order — is defined in that
`configs/settings.yml` under `sites.ordered`. **Always read that file** to discover which sites
exist. The list is **never** the same across projects (different names, different count, different
order), so never hardcode or assume site names.

### Environments
Available environments are discovered at runtime. The current environment is tracked in
`site_inception/.terraform/environment` (the first/inception site).

### 🚫 Environment switching is FORBIDDEN
This skill must **never** switch environments (`terrabutler env select` is forbidden).
If the current environment is not the one needed for the task, ASK THE USER to switch it
and provide the command they should run manually, e.g.:
```
terrabutler env select <name>
```

### Allowed read-only commands for discovery
```bash
terrabutler env show                          # Show current environment name
terrabutler env list                          # List available environments
cat <project-root>/configs/settings.yml       # Read sites.ordered for site list and deploy order
```

### Full `terrabutler tf` subcommand reference

All subcommands require `-site <site>`:

| Subcommand       | Description                                              |
|------------------|----------------------------------------------------------|
| `init`           | Prepare working directory                                 |
| `validate`       | Validate configuration files                              |
| `fmt`            | Reformat configuration in standard style                  |
| `plan`           | Show changes required by current configuration            |
| `apply`          | Create or update infrastructure                           |
| `destroy`        | Destroy infrastructure                                    |
| `show`           | Show current state or a saved plan                        |
| `console`        | Interactive expression evaluator                          |
| `output`         | Show output values from root module                       |
| `providers`      | Show required providers (subcommands: `lock`, `mirror`, `schema`) |
| `refresh`        | Update state to match remote systems                      |
| `import`         | Associate existing infrastructure with Terraform          |
| `state`          | Advanced state management (`list`, `mv`, `pull`, `push`, `rm`, `show`, `replace-provider`) |
| `taint`          | Mark resource instance as not fully functional            |
| `untaint`        | Remove tainted state from resource instance               |
| `force-unlock`   | Release a stuck lock on the current workspace             |
| `generate-options` | Generate terraform options for init/plan/apply         |
| `version`        | Show current Terraform version                            |

> These are terrabutler's **own** flags (double-dash, e.g. `--no-input`), not raw Terraform's.
> Boolean flags take no value — their presence flips the behavior (e.g. `--no-input` disables
> input; it is not `--input=false` or `--input BOOLEAN`). Never pass Terraform's own single-dash
> flags directly to `terrabutler tf`.

### Flags for the allowed commands (`fmt`, `validate`, `show`, `output`)

| Command    | Flag          | Description                                                              |
|------------|---------------|---------------------------------------------------------------------------|
| `fmt`      | `--diff`      | Display diffs of formatting changes                                       |
| `fmt`      | `--recursive` | Also process files in subdirectories                                      |
| `fmt`      | `--no-color`  | Output without color                                                      |
| `validate` | `--json`      | Machine-readable JSON output                                              |
| `validate` | `--no-color`  | Output without color                                                      |
| `show`     | `[PATH]`      | Optional positional arg: path to a saved plan file, instead of current state |
| `show`     | `--json`      | Machine-readable JSON output                                              |
| `show`     | `--no-color`  | Output without color                                                      |
| `output`   | `--json`      | Machine-readable JSON output                                              |
| `output`   | `--raw`       | Print a single value's raw string instead of a formatted representation   |
| `output`   | `--no-color`  | Output without color                                                      |

### Flags for `plan` / `apply` / `destroy` — reference only, these commands are FORBIDDEN

Never run these. Listed only so that a forbidden-command message you surface to the user quotes
the correct flag names.

| Flag             | `apply` | `plan` | `destroy` | Description                                 |
|------------------|:-------:|:------:|:---------:|----------------------------------------------|
| `--auto-approve` | ✅      |        | ✅        | Skip interactive approval — never use         |
| `--destroy`      | ✅      | ✅     |           | Select destroy planning mode                  |
| `--no-input`     | ✅      | ✅     | ✅        | Disable interactive input prompts             |
| `--no-lock`      | ✅      | ✅     | ✅        | Don't hold a state lock                       |
| `--lock-timeout` | ✅      | ✅     | ✅        | Duration to retry a state lock                |
| `--no-color`     | ✅      | ✅     | ✅        | Output without color                          |
| `--refresh-only` | ✅      | ✅     | ✅        | Select refresh-only planning mode             |
| `--no-refresh`   | ✅      | ✅     | ✅        | Skip checking for external changes            |
| `--target`       | ✅      | ✅     | ✅        | Limit to given module/resource (repeatable)   |
| `--var`          | ✅      | ✅     | ✅        | Set an input variable (repeatable)            |
| `--out`          |         | ✅     |           | (plan only) Write plan file to given path     |

---

## 2. Site Deploy Order (STRICT)

Deploy order is **strict** — each site depends on resources created by prior sites.
The first site in `sites.ordered` (typically an "inception" site) creates backends for
all subsequent sites and must run first.

**Always read `configs/settings.yml` → `sites.ordered` to determine the current deploy order.**
Do not assume specific site names — every Terrabutler project may have different sites.

---

## 3. Environments

Environments are **not hardcoded**. Discover them at runtime:

```bash
terrabutler env show    # What environment is active right now?
terrabutler env list     # What environments exist?
```

The current environment is tracked in `site_inception/.terraform/environment`.
The `TERRABUTLER_ENV` variable is read from this file automatically by mise.

---

## 4. Authentication & Login (per-project — discover, don't assume)

Login flows, cloud accounts, profiles, tenant IDs, and personal-secret locations **differ per
project**. Never assume the values from one repo apply to another. Discover them:

```bash
mise tasks              # Lists the project's login / auth / helper tasks
```

Typical shape (names and values vary by project — confirm with `mise tasks`):

```bash
mise login          # Log in to all clouds the project uses
mise login:aws      # AWS SSO (profile is project-specific)
mise login:google   # gcloud auth application-default login
mise login:az       # az login (tenant is project-specific)
```

> ⚠️ All `mise login*` / `az:*` tasks are on the FORBIDDEN list for this read-only skill —
> surface the exact command for the user to run; do not run it yourself.

### Key env vars (set by mise)
| Variable            | Purpose                                                          |
|---------------------|------------------------------------------------------------------|
| `TERRABUTLER_ROOT`  | Project root (`{{config_root}}`)                                 |
| `TERRABUTLER_ENV`   | Current env (read from `site_inception/.terraform/environment`)  |

Per-project auth config (e.g. AWS config file, `KUBECONFIG`) and personal secret files live in
project-specific locations — check the repo's `docs/`, `utils/`, and `mise` config rather than
assuming a fixed path.

---

## 5. Repository Layout

The project structure is consistent across Terrabutler projects, but the **set of sites** differs
— more or fewer sites, different names. Always derive the actual site list from
`configs/settings.yml` → `sites.ordered`, never assume specific names.

```
<project-root>/
├── configs/
│   ├── settings.yml              # Terrabutler config (envs, site order)
│   ├── backends/                 # Backend tfvars per site+env
│   ├── kubernetes/               # Kubeconfig files per env (if GKE is used)
│   └── variables/                # Variable tfvars (global, env, site+env)
├── globals/                      # Shared TF code (data sources, locals, variables)
│   ├── data.tf                   # Cloud project, KMS, secrets
│   ├── locals.tf                 # resource_name_prefix, tags, AZs
│   └── variables.tf              # Global variable definitions
├── site_<name>/                  # One directory per site (names vary per project)
├── docs/                         # Manual setup guides (project-specific auth lives here)
└── utils/                        # Utility scripts and auth config
```

Each `site_*` directory follows a convention:
- `terraform.tf` — Backend config + required providers
- `provider.tf` — Provider configurations
- `variables.tf` / `variables_global.tf` — Site and global variables
- `locals_global.tf` / `locals_globals.tf` — Shared locals
- `data_global.tf` — Shared data sources

---

## 6. Conventions

### Scope
- **In scope**: `*.tf`, `*.tfvars`, `*.hcl`, `*.tftpl`
- **Out of scope**: Application code (Java/Python/TypeScript) unless explicitly asked

### Formatting
- All `.tf` files must pass `terraform fmt -check` (typically enforced by CI)
- Run `terraform fmt -recursive` before committing

### Secrets
- Secrets are commonly AWS KMS encrypted maps: `var.global_encrypted_secrets` / `var.site_encrypted_secrets`
- Decrypted at plan time via a `data.aws_kms_secrets` in `globals/data.tf`
- Personal secrets live in a per-project, non-committed dir — see the repo's docs

### Backend pattern
- Sites commonly use an **S3 backend** with DynamoDB state locking (varies by project/cloud)
- Backend configs are generated by the inception site → `configs/backends/<project>-<env>-<site>.tfvars`
- Bucket naming is project-specific (e.g. `<project>-<env>-site-<site>-tfstate`) — read the
  generated backend tfvars rather than assuming a name

---

## 7. Gotchas

1. **"Not enabled" ≠ broken** — the `TERRABUTLER_ENABLE` message means the repo isn't initialized
   or mise's env isn't loaded in this shell. Fix: user runs `mise login` + `terrabutler init`
   (see §0.5). Never hand-set `TERRABUTLER_ENABLE` to bypass it.
2. **Deploy order is strict** — the first site in `sites.ordered` creates backends for all others.
3. **`configs/settings.yml` may be nested** — it is relative to the project root, not always the repo root.
4. **Provider versions vary by site** — check each site's `terraform.tf` for required versions.
5. **Current env may be production** — always check `terrabutler env show` before running anything.
6. **Site names vary per project** — never assume specific site names; read `configs/settings.yml`.
7. **Auth varies per project** — profiles, tenants, and secret paths differ; discover via `mise tasks`.
8. **Flags are terrabutler's own, not Terraform's** — double-dash, boolean flags take no value
   (e.g. `--no-input`, not `--input=false`). Never guess Terraform's raw flag spelling.

---

## Quick Reference (allowed commands only)

```bash
find . -path '*/configs/settings.yml'   # Locate the Terrabutler project root
terrabutler env show                    # Current environment
terrabutler env list                    # Available environments
terrabutler tf -site <site> validate    # Validate a site's configuration
terrabutler tf -site <site> fmt         # Reformat a site's .tf files
terrabutler tf -site <site> show        # Show state or a saved plan
terrabutler tf -site <site> output      # Show output values
mise run tf:summarize -- <site>         # Summarize a plan (read-only), if defined
cat <project-root>/configs/settings.yml # Discover sites and deploy order
mise tasks                              # List all available mise tasks
```
