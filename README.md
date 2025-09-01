# Agent Summit Bazel Workshop #

The purpose of this workshop is to show Agent developers
how to interact with bazel. We are mainly focusing on
`go` as other parts, such as dealing with native dependencies
will be mainly in scope of `Agent Build` team. It is not expected
that participants of this workshop have prior `bazel` experience. 
In fact, our intention is to show how easy it actually is to start
working with it.


## Setup ##
1. Install `bazelisk`. We want to be able to bump bazel's version whenever we
want or need. Major LTE versions are released every year and usually bring
UX and performance improvements, sometimes it is necessary to also upgrade to
new patch or minor versions to fix bugs. `bazelisk` reads `bazelversion` file
and automatically pulls the version that is needed by this particular project.
```zsh
which bazelisk

# If you have no installation then do so
brew install bazelisk
```

2. Create `.bazelversion` file. The file is read by `bazelisk` to pull
required bazel's version. This makes the process of build system update 
completely transparent to the users, we also ensure that going to old
branches doesn't break the build process:
```
# We are using latest stable version
8.3.1
```

3. Create `MODULE.bazel` file. For external dependency management `bazel`
is using a built-in system called `bzlmod`, it works similar to any other
modern dependency manager. `MODULE.bazel` file is the main source, it is
possible to define several `MODULE` files and then import them into the root one,
but **IT IS IMPORTANT TO HAVE AT LEAST ONE MODULE.bazel FILE AS IT HINTS BAZEL**
**THAT THIS IS BAZEL'S WORKSPACE**.

```python
bazel_dep(name = "rules_go", version = "0.57.0")
bazel_dep(name = "gazelle", version = "0.45.0")


# Bazel can manage toolchains and SDKs for us.
# This way we don't need to manually install the Go SDK
# and we ensure that all users have the same version of the Go SDK.
go_sdk = use_extension("@rules_go//go:extensions.bzl", "go_sdk")

# Setting Go SDK version.
# Alternative ways to set Go SDK can be found here:
# https://github.com/bazel-contrib/rules_go/blob/master/docs/go/core/bzlmod.md#go-sdks
go_sdk.download(version = "1.24.4")
```
4. Create `BUILD.bazel` file in the root of the project. BUILD files are 
letting `bazel` know what is considered a package. The presence of the BUILD file
in a directory makes the content of that directory visible to `bazel`. The content of
the `./BUILD.bazel` file:
```python
load("@gazelle//:def.bzl", "gazelle")

gazelle(name = "gazelle")
```

5. Run `gazelle`. As you may have noticed in the previous step we are adding some
magical lines mentioning gazelle. It is mentioned in `MODULE.bazel` as well as in `BUILD.bazel`.
`gazelle` is a BUILD file generator tool for `bazel`. Go projects' structure is usually very
straight forward and go's build system is modern enough to rely on it instead of re-inventing 
the wheel again. With that being said, most of the time we don't need to deal with bazel's internals
and can just let `gazelle` do its work:
```zsh
bazel run //:gazelle

# Now let's run git to show us generated files
git status

# Let's try to build and test the project
bazel build //...

# Execute unit tests
bazel test //...
```

6. Run the demo application:
```zsh
bazel run //cmd:cmd
```

