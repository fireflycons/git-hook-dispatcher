# Git Hook Dispatcher for Windows

Setting up local hooks on Windows is problematic because git is built with POSIX filesystems in mind, so does not expect hook scripts to have file extensions, merely to be executable which on such filesystems is determined by file attributes. Windows determines whether a file is executable by its extension.

This little utility helps you to run native Windows scripts e.g. batch/cmd or PowerShell when using git on Windows without needing git-bash or other such workarounds.

## Installation

You can quickly install it using the provided [install.ps1](./install.ps1) which will pull the latest release from github and set it up in `.githooks` in your user profile folder. If you run this script "As Administrator", then it will use symbolic links to create all the hook points, otherwise it will make a copy of the executable for each hook point (using more disk space). See below for more details on this.

## How it works

The binary `hook.exe` that this package builds works as a hook dispatcher. This is what is invoked by the git hooks process and will pass the environment and arguments to a Windows script file.

By having multiple copies of this utility named after each hook type, then git will find it and call it. The utility will derive the name of the repo that invoked the hook and then look for a script with extension of either `.cmd`, `.bat` or `.ps1` and a path of the format `<hookdir>\<repo>\<hookname>.<extension>` and invoke it with the working directory (typically the root of the local repo), arguments and [environment](https://git-scm.com/book/en/v2/Git-Internals-Environment-Variables) set up by git.

If no directory `<hookdir>\<repo>` is found, then directory `<hookdir>\00-githooks-shared` will be searched. If this directory is present and contains the relevant hook script, then that will be executed as a default for all repos configured to use this `<hookdir>`

### Debugging

If the environment variable `GITHOOK_DEBUG` is present with any value, the dispatcher will print information about the hook being called. This can be useful to see exactly which hooks are called in which order for any git workflow involving hooks, e.g. a commit would output the following if you have not defined any hook scripts

```
(No script found for hook post-index-change)
(No script found for hook pre-commit)
(No script found for hook prepare-commit-msg)
(No script found for hook commit-msg)
(No script found for hook reference-transaction)
(No script found for hook reference-transaction)   <- As many reference transactions as it needs.
(No script found for hook post-commit)
```

### Example

1. Create a directory to hold hooks, e.g. `C:\Users\myname\.githooks`
1. Copy this utility into it (see setup further down)
1. Init a git repo and set up `hooksPath`
    ```powershell
    cd C:\Users\myname\repos
    git init test-repo
    cd test-repo
    git config core.hooksPath C:/Users/myname/.githooks
    ```

    You can set it globally for all current and future repos with

    ```
    git config --global core.hooksPath C:/Users/myname/.githooks
    ```

1. Create a sub-directory in your hooks directory with the same name as the repo
    ```powershell
    New-Item -Type Directory -Path C:\Users\myname\.githooks -Name test-repo
    ```
1. Create a pre-commit hook for your repo in the new directory and put your logic in the script.
    ```powershell
    New-Item -Type File -Path C:\Users\myname\.githooks\test-repo -Name pre-commit.ps1
    ```
1. Try committing something to your new repo

## Manual Setup

For all [hooks](https://git-scm.com/docs/githooks#_hooks) you want to support, you need to make a copy of `hooks.exe` in the hooks directory you created.

If you have admin permission on your machine, you can use symlinks which take up no extra disk space and facilitate future upgrades:

1. Run a PowerShell command prompt As Administrator
1. CD to your hooks directory
1. For each hook, make a symlink, e.g
    ```powershell
    New-Item -Type SymbolicLink -Name pre-commit.exe -Target hook.exe
    New-Item -Type SymbolicLink -Name pre-merge-commit.exe -Target hook.exe
    ```

If you do not have admin permission, then make copies:

1. Run a command prompt
1. CD to your hooks directory
1. For each hook, make a copy, e.g
    ```dos
    copy hook.exe pre-commit.exe
    copy hook.exe pre-merge-commit.exe
    ```

    Remember that if you update `hook.exe` from a newer release, then you must redo all the copies!
1. You can delete `hook.exe` as it won't be invoked directly.


