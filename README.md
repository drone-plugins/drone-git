An alternate git clone implementation based on this discussion: https://discourse.drone.io/t/planned-change-to-git-clone-logic/1165

Build the plugin:

```
docker build -t plugins/git:next .
```

Test cloning a commit:

```
docker run --rm \
-v /tmp/test/commit:/data \
-e DRONE_BUILD_EVENT=push \
-e DRONE_REMOTE_URL=https://github.com/octocat/Hello-World.git \
-e DRONE_WORKSPACE=/data \
-e DRONE_COMMIT=762941318ee16e59dabbacb1b4049eec22f0d303 \
-e DRONE_BRANCH=master \
plugins/git:next
```

Test cloning a feature branch:

```
docker run --rm \
-v /tmp/test/branch:/data \
-e DRONE_BUILD_EVENT=push \
-e DRONE_REMOTE_URL=https://github.com/octocat/Hello-World.git \
-e DRONE_WORKSPACE=/data \
-e DRONE_COMMIT=b3cbd5bbd7e81436d2eee04537ea2b4c0cad4cdf \
-e DRONE_BRANCH=test \
plugins/git:next
```

Test cloning a tag:

```
docker run --rm \
-v /tmp/test/tag:/data \
-e DRONE_BUILD_EVENT=tag \
-e DRONE_REMOTE_URL=https://github.com/octocat/linguist.git \
-e DRONE_WORKSPACE=/data \
-e DRONE_TAG=v4.8.7 \
plugins/git:next
```

Test cloning a pull request:

```
docker run --rm \
-v /tmp/test/pr:/data \
-e DRONE_BUILD_EVENT=pull_request \
-e DRONE_REMOTE_URL=https://github.com/octocat/Spoon-Knife.git \
-e DRONE_WORKSPACE=/data \
-e DRONE_PULL_REQUEST=14596 \
-e DRONE_COMMIT=26923a8f37933ccc23943de0d4ebd53908268582 \
-e DRONE_BRANCH=master \
plugins/git:next
```

Try the plugin:

```
clone:
  git:
    image: plugins/git:next

pipeline: ...
```

This plugin is a work-in-progress and does not support the following features and / or capabilities:

* does not work with deployment events
* does not support depth
* only works with github pull requests (hard coded ref)
* cloning submodules is out of scope
