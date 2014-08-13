runki
=====

Ankiweb.net console client. Manage Anki flashcards without a friction.


Installation
------------

Use aur package https://aur.archlinux.org/packages/runki

For other systems you can install runki through `go get`:

  1. Install go.
  2. Execute `go get github.com/seletskiy/runki`


Configuration
-------------

All command line arguments can be stored in configuration file called
`~/.runki/runkirc` with following format:
```
-[option
  [value]]
```

Example:
```
-user
  user@example.com

-pass
  password

-deck
  english
```

See `./runki --help` for complete arguments list:

Usage
-----

Add card to ankiweb.net:
```
echo test | runki
# or
echo test | runki --user=<USERNAME> --pass=<PASSWORD> --deck=<DECKNAME>
```

Add list of cards to ankiweb.net:
```
cat words-list | runki
```

Add clipboard contents to ankiweb.net:
```
echo xclip -selection clipboard -o | runki
```

Add current selection to ankiweb.net:
```
echo xclip -o | runki
```

Add current selection to ankiweb.net and show translation:
```
notify-send "$(echo $(xclip -o) | runki"
```


Add word by shortcut (i3 window manager)
----------------------------------------

Create file ~/bin/add-anki-word with following contents:
```bash
#!/bin/bash

notify-send "$(echo $(xclip -o) | runki)"
```

Execute following command:
```
echo "bindsym \$mod+Escape exec add-anki-word" >> ~/.i3/config && i3wm-msg reload
```


Kindle
------

It is possible to learn foreign words directly from Kindle. When you encounter
unknown word, just tap on it, then select "Highlight". All this highlights will
be saved to the file "/documents/My Clippings.txt" in the Kindle FS.

Afterwards this file can be filtered to extract single word highlights (see
example script here: https://github.com/seletskiy/dotfiles/blob/92a0e8fdb533106d700dc4c7415f4d41b232edff/bin/kindle-filter-words).

So, all highlighted words can be easily added to Anki like this:

```bash
kindle-filter-words "/mnt/documents/My Clippings.txt" | tac | runki --cut 3
```

Note, that `tac` command will reverse incoming list so most recent items will
came first. `--cut 3` flag tells runki to exit if it detect at least `3` words
that are already added to Anki. `tac` and `--cut` combination allows you to
sync words from Kindle in seamless way just running one single command again
and again.


Kindle + udev
-------------

Ok, let's go deeper. I want to sync words that I highlight on the kindle
automatically when I connect kindle to the computer.

So, we need to create udev rule like this:

```
ACTION=="add", SUBSYSTEM=="block", ENV{DEVTYPE}=="partition", ENV{ID_VENDOR_ID}=="1949", RUN+="/usr/bin/su <USERNAME> -lc 'DISPLAY=:0 kindle-to-anki $env{DEVNAME}'"
```

So, after kindle is connected, `kindle-to-anki` program will have to be runned.
It will add new words directly to Anki and show nice notification about how
many new words has been added to.

Example of this `kindle-to-anki` program can be found there: https://github.com/seletskiy/dotfiles/blob/1c9da6d347cc658c9d6d177a61ef94423a3c36d4/bin/kindle-to-anki

All further invocations of `runki` can be done without specifying user/pass/deck
arguments.


VERSION
-------

0.0.1
