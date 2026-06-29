---
name: maintenance-deps
description: Run dependency maintenance for this project — executes the zed:dependabot-sweep plugin skill.
---

# Dependency Maintenance

Run the full Dependabot maintenance sweep for this repository.

## Steps

1. Invoke the `zed:dependabot-sweep` skill and let it run to completion. It will
   unblock stuck Dependabot PRs, merge ready PRs, fix vulnerability alerts, update
   the changelog, and open a PR as appropriate for this repository.
2. Report what the sweep did.
