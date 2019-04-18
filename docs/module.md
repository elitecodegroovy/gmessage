
# Modules

## Tidying it up

By default, Go does not remove a dependency from go.mod unless you ask it to. If you have dependencies that you no longer use and want to clean up, you can use the new tidy command:

```
go mod tidy
```


## Vendoring

Go modules ignores the vendor/ directory by default. The idea is to eventually do away with vendoring1. But if we still want to add vendored dependencies to our version control, we can still do it:

```
go mod vendor
```

This will create a vendor/ directory under the root of your project containing the source code for all of your dependencies.Still, go build will ignore the contents of this directory by default. If you want to build dependencies from the vendor/ directory, you’ll need to ask for it.

```
go build -mod vendor
```