# Agent Summit Bazel Workshop #

The purpose of this workshop is to show Agent developers
how to interact with bazel. We are mainly focusing on
`go` since other parts, such as dealing with native dependencies,
will be mainly in scope of the `Agent Build` team. It is not expected
that participants of this workshop have prior `bazel` experience. 
In fact, our intention is to show how easy it actually is to get started
working with it.


## Setup ##
1. Install `bazelisk`. We want to be able to bump bazel's version whenever we need.
Major LTE versions are released every year and usually bring
UX and performance improvements. Sometimes, it is also necessary upgrade to
new patch or minor versions to fix bugs. `bazelisk` reads the `.bazelversion` file
and automatically pulls the version that is needed by this particular project.
### Macos ###
```zsh
$> which bazelisk

# If you have no installation then do so
$> brew install bazelisk
```

### Windows ###
```powershell
$> winget install Bazel.Bazelisk
```
or
```powershell
$> choco install bazelisk
```
or
```powershell
$> scoop install bazelisk
```

### Linux ###
You can either download an executable from [releases](https://github.com/bazelbuild/bazelisk/releases)
or a debian package (if applicable) and run:
```zsh
$> dpkg -i bazelisk-<arch>.deb
```

**As a last resort you can always download executables for your platform**
**from [releases](https://github.com/bazelbuild/bazelisk/releases)**

2. Create a `.bazelversion` file at the root of the repository. The file is read by `bazelisk` to pull
required bazel's version. This frees users from having to manually manage the version of the build system
and avoids versions drifting between developer machines and with CI.
It also lets us ensure that going to old
branches doesn't break the build process:
```
# We are using latest stable version
8.3.1
```

You can confirm that this works as expected by running:

```bash
$> bazel --version
```

3. For external dependency management, `bazel`
uses a built-in system called [`bzlmod`](https://bazel.build/external/overview), which works similarly to other
modern dependency managers you may already be familiar with. The `MODULE.bazel` file at the root of a project
marks it as a [Bazel module](https://bazel.build/external/module), and it's where dependencies are defined.
It is possible to import other `MODULE` files inside the same project as a way to split the contents
across files, but **it is important to have at least one MODULE.bazel file at the root as it's what bazel**
**uses as [repository](https://bazel.build/concepts/build-ref#repositories) boundary markers**.

Edit the provided empty `MODULE.bazel` file and add the following contents:

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

4. Create a `BUILD.bazel` file in the root of the project. BUILD files let
`bazel` know what is considered a [package](https://bazel.build/concepts/build-ref#packages). The presence of the BUILD file
in a directory makes the content of that directory visible to `bazel`. Add the following to
the `./BUILD.bazel` file:
```python
load("@gazelle//:def.bzl", "gazelle")

gazelle(name = "gazelle")
```

5. Run `gazelle`. As you may have noticed in the previous step we added some
magical lines mentioning `gazelle`. It is mentioned in `MODULE.bazel` as well as in `BUILD.bazel`.
[`gazelle`](https://github.com/bazel-contrib/bazel-gazelle) is a BUILD file generator tool for `bazel`. Go projects' structure is usually very
straight forward and go's build system is modern enough to rely on it instead of re-inventing 
the wheel again. That means that most of the time we don't need to interact with BUILD files directly
and can just let `gazelle` do its work:
```zsh
$> bazel run //:gazelle

# Now let's run git to show us generated files
$> git status
```

You can take a look at the generated `BUILD.bazel` files to get a feel for what they look like, even though, as just mentioned, most of the time you won't need to make changes to them.
The contents can be fairly easy to understand most of the time even without knowing much about Bazel.

```zsh
# Let's try to build and test the project
$> bazel build //...

# Execute unit tests
$> bazel test //...
```

6. Run the demo application:
```zsh
$> bazel run //cmd:cmd
```
# Exercises #
## Adding new package ##
This is a "freestyle" exercise. You can however just follow the presenter.

1. Create a new folder within [./pkg](./pkg/) of your choice
2. Create a new `.go` source file of your choice and implement some logic. Keep it independent from the rest of the packages for now.
3. Now you need to make bazel aware of your newly created package. To do so just run gazelle.
```zsh
# This will shrink gazelle's scope to speed up the process and not go
# through the entire project. In our case we could also run gazelle
# without specifying path as this project is very simple and small.
$> bazel run //:gazelle -- /pkg/<new_pkg>
```
4. **(Optional)** Now let's try to import our new package into already existing code.
Let's see what happens if we run the same command as above:
```zsh
$> bazel run //:gazelle -- /pkg/<new_pkg>
```
As you see nothing has changed in the `<new_pkg>`. So we need to update the package
where we imported our new package:
```zsh
# Example
$> bazel run //:gazelle -- /pkg/server # In case changes were made in server

# Template command
$> bazel run //:gazelle -- /pkg/<changed_pkg>

# Now let's see if BUILD file was updated accordingly
$> git diff
```

## Flaky Tests Detection ##
`bazel` is not only powerful when it comes to building, but also very helpful running
tests. It allows us:
- to set how many times we should run certain tests or all of them by setting [--runs_per_test](https://bazel.build/reference/command-line-reference#flag--runs_per_test)
- to automatically detect flaky tests by setting [--runs_per_test_detects_flakes](https://bazel.build/reference/command-line-reference#flag--runs_per_test_detects_flakes). When used in combination with the flag above it will not fail the results of `bazel test` commands,
but will mark failing tests as `FLAKY` instead of `FAILED`
- to set timeouts and compute resources based on [size attribute](https://bazel.build/reference/be/common-definitions#common-attributes-tests)
- to mark tests known to fail now and then as flaky by setting `flaky = True` attribute to the `*_test` target. 

Now let's actually ask bazel to run tests several times to see if we have any failing tests.
```zsh
$> bazel test //... --runs_per_test=10
```

Now let's see if failing tests are failing constantly or flaky.
```zsh
$> bazel test //... --runs_per_test=10 --runs_per_test_detects_flakes
```

And now let's mark our known test as flaky in `./pkg/stock/BUILD.bazel`:
```python
go_test(
    name = "stock_test",
    srcs = [
        "client_test.go",
        "service_test.go",
    ],
    embed = [":stock"],
    deps = [
        "//internal/testutils",
        "//pkg/models",
    ],
    flaky = True,
)
```

## Worth reading ##
- [Bazel command line reference](https://bazel.build/reference/command-line-reference) - overview of all flags available in bazel. Keep in mind that flags may change
based on bazel's version so it is recommended to use versioned docs (see navigation bar)
- [Working with Go in bazel](https://github.com/bazel-contrib/rules_go/blob/master/docs/go/core/bzlmod.md#specifying-external-dependencies) - must read for anyone who is planning to use bazel with their codebase.
- [Bazel Central Registry](https://registry.bazel.build/) - even though `bzlmod` is still fresh and a lot of legacy projects are still relying on `WORKSPACE` approach
more and more modules and tools are added to the Central Registry, so before thinking on "bazelization" of a new tool or a dependency it's worth checking if it isn't
already present in BCR.
- [Reasoning and migration guide for bzlmod](https://bazel.build/external/migration) - we are starting fresh in the Agent, therefore, we will use `bzlmod` right away,
so you won't have to deal with `WORKSPACE`, however it's worth reading this article to understand the background.
