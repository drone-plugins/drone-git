Use the Git plugin to clone a git repository. Note that Drone uses the Git plugin
by default for all repositories, without any configuration required. You can override
the default configuration with the following parameters:

* `depth` - creates a shallow clone with truncated history
* `recursive` - recursively clones git submodules
* `skip_verify` - disables ssl verification when set to `true`
* `tags` - fetches tags when set to `true`
* `submodule_override` - override submodule urls

Sample configuration:

```yaml
clone:
  depth: 50
  recursive: false
  tags: false
```

## Submodules

Sample configuration to clone submodules:

```
clone:
  recursive: true
```

Sample configuration to override submodule urls:

```
clone:
  recursive: true
  submodule_override:
    hello-world: https://github.com/octocat/hello-world.git
```

The above configuration is intended for private submodules created using ssh clone urls (i.e. `git@github.com:octocat/hello-world.git`). We recommend overriding to use https clone urls to take advantage of this plugins built-in `netrc` authentication mechanism.