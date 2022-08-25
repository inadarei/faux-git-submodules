# faux-git-submodules

Because native git submodules are buggy and unreliable

## Usage

checks-out or updates repos depicted in the current folder's fgs.json file. 
The binary can be installed anywhere that is reachable from $PATH. 

If branch is missing in the fgs.json, "main" is assumed by default.

Makefile.sample contains code that can help execute mac, linux or windows
version of the binary, depending on the system.

## Python vs. Native.

Implementation of Faux Git Submodules is written in Go and native builds
are provided for a variety of platforms. You can find them under [build](build)
folder. These are, probably, easiest to use. However it is impossible for
the authors of this utility to provide a native build for every possible
operating system/cpu architecture combinations. If you don't see the platform
you are using, you can build yourself, or you can use the python version, given
that you have python3 on your target platform.
