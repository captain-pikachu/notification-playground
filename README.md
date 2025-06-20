## install dependencies

```bash
bazel mod tidy
```

## generate `BUILD` for go projects

```bash
bazel run //:gazelle_go

# specific dir
bazel run //:gazelle_go -- dir1 dir2
```

## sync bazel with go.mod

```bash
# create `deps.bzl` if it doesn't exist
touch deps.bzl

bazel run @io_bazel_rules_go//go -- mod tidy
bazel mod tidy
bazel run //:gazelle_go
```

## build

```bash
# build all
bazel build //...

# build a target (e.g. quotes_gen)
bazel build //quotes_gen
```

## run

```bash
# run a target
bazel run //quotes_gen

# or

# run binary manually
./bazel-bin/quotes_gen/quotes_gen_/quotes_gen
```

## clean

```bash
bazel clean --async

# --expunge: removes the entire working tree and stops the bazel server
bazel clean --expunge

# linux
sudo rm -rf ~/.cache/bazel

# mac
sudo rm -rf  /private/var/tmp/_bazel*
```
