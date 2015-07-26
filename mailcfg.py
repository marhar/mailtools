#!/usr/bin/env python
"""
mailcfg -- some stuff to manage mail servers

Generally handy tools for dealing with mail processing programs

This is part of the mailtools package :  (c) 2015 Mark Harrison.
https://github.com/marhar/mailtools   :  Share and Enjoy!
"""

import os
import imaplib

#-----------------------------------------------------------------------
class MailtoolsError(Exception):
    """standard exception for Mailtools package"""
    def __init__(self, msg):
        self.msg = msg
    def __str__(self):
        return repr(self.msg)

#-----------------------------------------------------------------------
class MailtoolsIMAPError(MailtoolsError):
    """indicates an unusual (non-OK) IMAP response"""
    pass

#-----------------------------------------------------------------------
def ok(rc,data):
    """raises exception if return code is not 'OK'"""
    if rc != 'OK':
        raise MailtoolsIMAPError(data)

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
        raise MailtoolsError('.imap file must not be readable by others')

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
