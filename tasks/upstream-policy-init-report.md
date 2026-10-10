# Reusable policy initialization

Filed on 2026-10-09 as [itos#33](https://github.com/donvargax/itos/issues/33), using the
consumer-report fields and label. Reproduced on 7.0.0-rc.3. Automatic template policy setup
depends on upstream support; no initialization workaround is installed here.

### Not reported already

- [x] I searched open and closed issues, and the newest release, and this is not one of them.

Checked the newest prerelease, v7.0.0-rc.3, and stable v6.5.1 on 2026-10-09. Searches
for `init`, `policy`, and `init policy profile reusable` found no policy-initialization
request. #30 concerns flag defaults; #31 concerns external data paths and differing
branch plans. This request is about fresh made-project adoption, not shared template data.

### What was run

From the shell, in two separate empty temporary directories. The consumer repository
`donvargax/go-template-itos` pins rc.1; these probes explicitly ran rc.3.

```sh
ITOS_VERSION=7.0.0-rc.3 ITOS_NO_UPDATE_NOTICE=1 \
  itos init --plugin no --no-git-shim --no-agent-rules --json
```

In the first directory there was no configuration. In the second, the only file was
the reusable policy shown below, named `itos.yaml`; there was no `tasks/` directory.
No application code or quality checks were executed by these probes.

### Its output

The empty-project run exited 0 and created `itos.yaml`, `tasks/phase-1.yaml` and
`tasks/work-items.yaml`. Its configuration was the generic starter, with no `ci` or
`proof` settings. The relevant exact output was:

```json
  "action": "initialized",
  "git_init": true,
  "since": null,
```

The supplied-policy run exited 1 and did not create the ledger or registry. Its exact
JSON output was:

```json
{
  "schema": 1,
  "config": "itos.yaml",
  "action": "checked",
  "missing": [
    {
      "rule": "ledger-folder-missing",
      "message": "itos.yaml: the ledger's folder tasks does not exist",
      "fix": "create tasks with the ledger's files (ledger.files is tasks/phase-{group}.yaml), or point ledger.files at the folder that holds them",
      "area": "ledger"
    },
    {
      "rule": "hook-missing",
      "message": "the commit-msg hook: hook.itos-commit-msg in the git config does not run itos",
      "fix": "run itos hook install",
      "area": "hooks"
    },
    {
      "rule": "hook-missing",
      "message": "the pre-push hook: hook.itos-pre-push in the git config does not run itos",
      "fix": "run itos hook install",
      "area": "hooks"
    }
  ],
  "plugin": {
    "action": "declined",
    "scope": null,
    "excluded": false
  },
  "notes": [],
  "git_shim": {
    "action": "declined",
    "link": null,
    "on_path": false,
    "before_git": false
  },
  "agent_rules": {
    "action": "declined",
    "files": []
  }
}
```

### itos version

```sh
ITOS_VERSION=7.0.0-rc.3 ITOS_NO_UPDATE_NOTICE=1 itos version
```

```text
itos 7.0.0-rc.3
```

### The config

The small supplied policy used for the second probe:

```yaml
version: 1
ledger:
  files: "tasks/phase-{group}.yaml"
  id: "T-\\d+"
commits:
  types: [feat, fix, refactor, perf, test, build, ci, chore, docs, style, revert]
  header_lint: { use: builtin }
  footers:
    Task:
      source: ledger
      required_for: [refactor, perf, test, build, ci, chore, style, revert]
      validate_for: all
      read_at: commit
hooks:
  bin: itos
ci:
  steps:
    - run: go test ./...
```

### What was expected

A supported, explicitly selected way for itos to initialize a project from reusable
policy while creating fresh project-local ledger and registry data. This is a capability
request: the current behavior agrees with the documented generic-starter/diagnostic
contract, but neither route supplies what a real-project template needs.

The template's branches already carry custom commit rules, CI, scenario selection and
code proof. Its renders deliberately exclude the template's `itos.yaml` and `tasks`,
because a made project must not inherit completed template tasks, owners, statuses or
history-specific SHAs. A template-owned policy can be reviewed and substituted alongside
the code, but plain `itos init` does not adopt that policy; preinstalling it instead makes
init report missing work data. Config selection is not initialization.

The requested mechanism should retain applicable policy and supported pins, create new
adoption/work data using the selected paths, set the new project's history boundary,
and install/generate its hooks and agent rules as requested. Existing project data must
not be overwritten silently. It must not import template tasks, owners, proof results or
mutation caches. Permanent quality checks belong in the new project's standing CI policy,
not copies of historical template task checks.

The interface is for itos to specify; this report does not presume a profile flag or
configuration-inheritance syntax. Keep itos-template unaware of itos: a template's declared
setup steps should invoke the supported itos mechanism. The person selected this upstream
route rather than a custom initialization workaround in go-template-itos.
