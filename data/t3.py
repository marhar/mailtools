#!/usr/bin/python
import sys
import email

msg = email.message_from_file(open(sys.argv[1]))

for part in msg.walk():
    if part.get_content_type() == 'text/plain':
        for line in part.get_payload(decode=True).split('\n'):
            if line.startswith('http://url5878.vsobr.com') or line.startswith('https://www.rcgroups.com'):
                print line
                break
