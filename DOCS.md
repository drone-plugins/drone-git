Use the Git plugin to clone a git repository. Note that Drone uses the Git plugin
by default for all repositories, without any configuration required. You can override
the default configuration with the following parameters:

* **depth** - creates a shallow clone with truncated history
* **recursive** - recursively clones git submodules

The following is a sample Git clone configuration in your .drone.yml file:

```yaml
clone:
  depth: 50
  recursive: false
```
