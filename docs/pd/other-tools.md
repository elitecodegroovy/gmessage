
# Tools

- `github.com/mattn/goreman` Manage Procfile-based applications.It Will start all commands defined in the Procfile and display their outputs. Any signals are forwarded to the processes.





## Switching remote URLs from HTTPS to SSH

Change the current working directory to your local project.

List your existing remotes in order to get the name of the remote you want to change.

```
git remote -v
origin https://github.com/USERNAME/REPOSITORY.git (fetch)
origin https://github.com/USERNAME/REPOSITORY.git (push)
```

Change your remote’s URL from HTTPS to SSH with the git remote set-url command.
```

git remote set-url origin git@github.com:USERNAME/OTHERREPOSITORY.git
```

Verify that the remote URL has changed.

