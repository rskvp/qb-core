# G&G Core

![](icon.png)

Utility library.

This folder contains some utility packages.

Mose of them are pure Go, only some uses external modules but always in pure Go.

No dependencies from other external binaries suing C bindings.

## Special Modules

- [Scheduler](./qb_scheduler/readme.md): Schedule a task.
- [Updater](./qb_updater/readme.md): Update and launch a program.

## How to Use

To use just call:

```
go get bitbucket.org/digi-sense/qb-core@latest
```

### Versioning

Sources are versioned using git tags:

```
git tag v0.1.159
git push origin v0.1.159
```
