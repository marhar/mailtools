#!/usr/bin/env python
"""
#-----------------------------------------------------------------------
mailcfg -- some stuff to manage mail servers

Generally handy tools for dealing with mail processing programs

This is part of the mailtools package : (c) 2015 Mark Harrison.
https://github.com/marhar/mailtools   :     Share and Enjoy!
"""

import os,imaplib,configparser

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
    """raises exception if IMAP return code is not 'OK'
       this typically should not occur and indicates either
       a coding problem or unexpected server operation.
    """
    if rc != 'OK':
        raise MailtoolsIMAPError(data)

#-----------------------------------------------------------------------
def imapserver(server,fname=None):
    """
    connect to the imap server specified in $HOME/.imap
    returns a standard IMAP4/IMAP4_SSL object.

    the config file is in standard ini-style.  a section looks like this:
        [gmail]
        ssl=True
        smtp=smtp.gmail.com
        imap=imap.gmail.com
        login=myname@gmail.com
        passwd=mypasswd
    """
    if fname is None:
        fname = os.environ['HOME']+'/.imap'

    # imap file can't be readable by other people!
    s=os.stat(fname)
    if s.st_mode & 0o077 != 0:
        raise MailtoolsError('%s must not be readable by others'%(fname))

    cfg=configparser.ConfigParser()
    cfg.read(fname)
    x={}
    for k,v in cfg.items(server):
        x[k]=v
    #import pprint; pprint.pprint(x)

    if x['ssl'] == 'True':
        m=imaplib.IMAP4_SSL(x['imap'])
    else:
        m=imaplib.IMAP4(x['imap'])

    m.login(x['login'],x['passwd'])
    return m

if __name__ == "__main__":
    m=imapserver('gmail')
