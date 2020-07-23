#!/usr/bin/python
import os
import sys
import mimetypes
import email

msg = email.message_from_file(open(sys.argv[1]))

nparts=0
for part in msg.walk():
    nparts+=1
    print 999999999999999999999999999999999
    print part.is_multipart()
    print part.keys()
    print part.get_content_maintype(),part.get_content_type(),part.get_content_subtype()
    #print type(part),dir(part)
    print part.get_payload(decode=True)

print 88888888,nparts
a="""
    if part.get_content_maintype() == 'multipart':
        continue
    ext = mimetypes.guess_extension(part.get_content_type())
    filename='part-%03d%s'%(i, ext)
    # filename=os.path.join('settings.MEDIA_ROOT', filename)
    print(filename)
    with open(filename, 'wb') as fp:
        fp.write(part.get_payload(decode=True))
"""
