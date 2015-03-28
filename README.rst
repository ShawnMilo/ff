==
ff
==

``ff`` is meant to be a replacement for the common ``find`` idiom
``find .  -iname '*foo*'``.

All searches are recursive and case-insensitive.

Usage::

    # Find a file named monkey.py
    ff monkey.py

    # Find Python files with "mon" in the name (argument order is not important):
    ff mon .py

License
=======

BSD
