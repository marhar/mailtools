#!/usr/bin/env python
"""
mailcfg -- some stuff to manage mail servers
"""

import os
import imaplib

#-----------------------------------------------------------------------
def imapserver(fname=None):
    """
    connect to the imap server specified in $HOME/.imap
    imap file looks like this:

        $ cat ~/.mailcfg 
        SSL=True
        smtp='smtp.example.com'
        imap='imap.example.com'
        login='myname'
        passwd='mypass'
    """
    if fname is None:
        fname = os.environ['HOME']+'/.imap'

    # imap file can't be readable by other people!
    s=os.stat(fname)
    if s.st_mode & 0077 != 0:
        raise RuntimeError('.imap file must not be readable by others')

    x={}
    execfile(fname,x)
    # del x['__builtins__']; print x

    if x['SSL']:
        m=imaplib.IMAP4_SSL(x['imap'])
    else:
        m=imaplib.IMAP4(x['imap'])

    m.login(x['login'],x['passwd'])
    return m

if __name__ == "__main__":
    m=imapserver()
