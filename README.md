qhash
=====

qhash is a small tool which can perform several hashing functions on files.

I wrote this for two reasons:
 1. I have a need to perform file hashing on systems which may not have the standard tools installed, and I wanted a cross-platform tool with as few run time dependencies as possible
 2. I want to keep experimenting with Go

This tool doesn't really exercise the run-time power of Go, but every little bit helps. I'm sure there are a ton of improvements which could be made to this code, I'll happily accept pull requests.

If you want to verify signed commits, my key is available on the public key servers or via `git show maintainer-pgp-pub | gpg --import`.
