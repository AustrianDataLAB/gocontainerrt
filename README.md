# Demo Repo for Lecture 2
(original content from the original authors, slightly modified)

## start by cloning the main branch 
Into a linux where you have root access
This is essentially a fork of https://github.com/adamgordonbell/chroot-containers
```
go build chrun.go
sudo ./chrun pull alpine
sudo ./chrun run alpine
```

## next we add namespaces
It assumes that you downloaded an alpine tar ball into the local ./assets dir which is implemented in the main branch
This is based on https://github.com/lizrice/containers-from-scratch/tree/master

Same procedure, compile, run, observe

## next we add cgroups
for cgroups-v2
