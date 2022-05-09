#!/usr/bin/env python
"""
#-----------------------------------------------------------------------
read-rcgroups -- read new messages on rcgroups.com

rcgroups sends a set of nightly emails (one per subscribed thread)
telling you what's new to be read.  this program grabs those messages,
opens the threads in your browser, and deletes the messages.  Don't
be alarmed if it opens a lot of tabs.

If you use gmail, add this filter:
    Matches: from:(rcgroups.com@vsobr.com)
    Do this: Skip Inbox, Apply label "[Gmail]/rcgroups"

usage: read-rcgroups [mailserver]

This is part of the mailtools package : (c) 2015-2019 Mark Harrison
https://github.com/marhar/mailtools   :     Share and Enjoy!
"""

#-----------------------------------------------------------------------
import sys,time,os,email, urllib.request, urllib.error, urllib.parse
import mailcfg

#-----------------------------------------------------------------------
def doopen(url):
    """open this url in a browser
       currently works for mac.  if you're on another system let me
       know what works!
    """
    os.system("open '%s'"%(url))
    #time.sleep(.25)  # this delay seems to makes opening smoother

#-----------------------------------------------------------------------
def url(zzemail_text):
    """look for the strings that indicate unread forums or threads.
       rcgroups keeps changing the message format, so we adjust whenever
       we notice things not working.
    """
    zzmsg = email.message_from_string(zzemail_text)
    import pudb;pudb.set_trace()
    for part in zzmsg.walk():
        print('part', part.get_content_type())
        if part.get_content_type() == 'text/plain':
            for line in part.get_payload(decode=True).split(b'\n'):
                line=line.strip()
                if line.startswith(b'http://url5878.vsobr.com'):
                    class NoRedirection(urllib2.HTTPErrorProcessor):
                        def http_response(self, request, response):
                            print(666666)
                            return ',',response
                        https_response = http_response
                    opener = urllib.request.build_opener(NoRedirection)
                    urllib.request.install_opener(opener)
                    req=urllib.request.Request(line)
                    response = urllib.request.urlopen(req)
                    for cline in response.info().headers:
                        x = cline.split()
                        if x[0] == 'Location:':
                            print(7777777)
                            return '.',x[1]+'&goto=newpost'
                    print(88888)
                    return None
                if line.startswith(b'https://www.rcgroups.com'):
                    print(99999)
                    return '.',line+'&goto=newpost'
    print(5555555)
#-----------------------------------------------------------------------
def main():
    """Our main business is not to see what lies dimly at a
       distance, but to do what lies clearly at hand.
          --Thomas Carlyle
    """

    showstart()
    logger = open("rcreader.log", "a")
    logger.write("-----------\n")
    logger.flush()
    if len(sys.argv) > 1:      # max messages to process (including dupes)
        N=int(sys.argv[1])
    else:
        N=200

    if len(sys.argv) > 2:      # default to gmail
        server=sys.argv[2]
    else:
        server='gmail'

    m=mailcfg.imapserver(server)
    rc,data=m.select('[Gmail]/rcgroups') # we filter to here on gmail
    mailcfg.ok(rc,data)

    rc,data=m.search(None, '(FROM "rcgroups.com@vsobr.com")')
    mailcfg.ok(rc,data)

    seen={}
    nshown=0
    nskipped = 0
    nerrored = 0
    totmsgs=len(data[0].split())
    #print('messages:',totmsgs)
    i = 0
    for msgno in data[0].split():
        i += 1
        rc,data=m.fetch(msgno,'(RFC822)')
        mailcfg.ok(rc,data)
        body=email.message_from_string(str(data[0][1]))
        signal,h=url(str(body))
        if h is None:
            nerrored += 1
            show(i, totmsgs, nskipped, "E", runstate)
        elif (h in seen) == False:
            nshown += 1
            seen[h] = 1
            logger.write("open '" +h + "'\n")
            logger.flush()
            if DOIT: doopen(h)
            show(i, totmsgs, nskipped, signal, runstate)
        else:
            show(i, totmsgs, nskipped, "o", runstate)
            nskipped += 1
            #print('# (skip)', h)
        if DOIT: m.store(msgno, '+FLAGS', '\Deleted')

        if nshown >= N:
            print(('*** stopped after first %d, %d left'%(N,totmsgs-N)))
            break

    sys.stdout.write("\n%d opened %d skipped\n" % (nshown, nskipped))
    sys.stdout.write("%d errored\n" % (nerrored))
    m.close()
    m.logout()

def ttt(s):
    """write string to screen"""
    sys.stdout.write(s)
    sys.stdout.flush()

def showstart():
    ttt("\033[2J\033[H")

def show(i, totmsgs, nskipped, ch, extra):
    if not DOIT:
        return
    WIDTH = 10
    ttt("\033[Htotmsgs=%d i=%d skipped=%d %s"%(totmsgs, i, nskipped, extra))
    cc = (i - 1) // WIDTH
    rr = (i - 1) % WIDTH
    rr += 1
    cc += 2
    ttt("\033[%d;%dH%s"%(cc,rr,ch))
    time.sleep(DELAY)

#-----------------------------------------------------------------------
DELAY=0  # make larger for debugging or for slow computer
DOIT=True
#DOIT=False  # uncomment for testing
if (not DOIT):
    runstate = "DRY RUN MODE"
else:
    runstate = ""

try:
    main()
except KeyboardInterrupt:
    pass
