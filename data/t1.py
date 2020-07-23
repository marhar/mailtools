import os
import sys
import mimetypes
import email

msg = email.message_from_file(open(sys.argv[1]))

for i,part in enumerate(msg.walk(),1):
    if part.get_content_maintype() == 'multipart':
        continue
    ext = mimetypes.guess_extension(part.get_content_type())
    filename='part-%03d%s'%(i, ext)
    # filename=os.path.join('settings.MEDIA_ROOT', filename)
    print(filename)
    with open(filename, 'wb') as fp:
        fp.write(part.get_payload(decode=True))
