A simple tool that executes all input files in bash and output the difference in `env` and `alias` before and after execution.

## Usage

```bash
> shell-cross-env
Usage: shell-cross-env --to <fish|bash> source <file1> <file2>
```

## Example

```bash
# In fish shell, do this to load environment variables and aliases from bash files
shell-cross-env --to fish /etc/profile | source
```
