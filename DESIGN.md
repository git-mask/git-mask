There are 2 modes git-mask can be setup:

1. Plugin Mode.

User interact with git-mask with `git mask <command> [flags]` commands.

In this case, the official `git` executable is first invoked. It searches for `git-mask` executable from PATH, and invoke it, passing the rest of the user arguments to it.

* Entry executable: official `git`
* git-mask executable file name: `git-mask`

2. Wrapper Mode.

User install git-mask executable as `git` with higher priority in the PATH search list.

In this case, user can use Git commands seemlessly like `git clone git@github.com:user/repo`. Under the hood, the git-mask executable disguised with the file name `git` is first invoked. git-mask would search for the next `git` executable from the PATH to be used as the real official `git`. Then the user's command will be handled in the same way as `git mask run -- clone git@github.com:user/repo`. Unless the user is executing `git mask <command> [flags]`, then git-mask would handle it internally.

* Entry executable: git-mask
* git-mask executable file name: `git`


## git-mask initialization

Determine which mode git-mask is running in.

* If git-mask's executable name is `git`, then it's running in Wrapper Mode.
* Otherwise we can say it's running in Plugin Mode.

If git-mask is in Wrapper Mode:
    * If the top-level command is `mask`, then handle the rest of the command line with git-mask's Plugin Mode's rootCmd
    * Else, handle the rest of the command line with git-mask's runCmd




## git-mask commands

### `git mask run [flags] -- <cmd>`

Invoke `git <cmd>` using the official `git`

### `git mask connect [flags] <host> <port>`

* `--proxy` supports `socks5://`, `https://` etc

### `git mask profile add <name>`

* `--name`
* `--email`
* `--timezone`
* `--vague-time`
* `--proxy`
* `--gpg-key`
* `--ssh-key`

### `git mask profile rm <name>`

### `git mask rule add [rule_name] [flags]`

* `--url` glob for remote url matching (can have multiple)
* `--profile`

