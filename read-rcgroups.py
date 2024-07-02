#!/usr/bin/env python
"""
# Read all rcgroups.com content from email notifications.

## Enable IMAP in Gmail:

- Go to your Gmail account settings.
- Under the "Forwarding and POP/IMAP" tab, enable IMAP access.
- Create App Password

https://knowledge.workspace.google.com/kb/how-to-create-app-passwords-000009237
https://myaccount.google.com/apppasswords

## Forum

https://www.rcgroups.com/forums/forumdisplay.php?f=172

## Thread

https://www.rcgroups.com/forums/showthread.php?t=3496155

Append this: `&goto=newpost`
"""

import imaplib
import email
import os
import re
import subprocess
import time
from email.header import decode_header

# Account credentials and IMAP host from environment variables
username = os.environ["RCMAIL_USER"]
app_password = os.environ["RCMAIL_PASSWD"]
imap_host = os.environ["RCMAIL_IMAP"]

# Settings

rcgroups_folder = "[Gmail]/rcgroups"
#rcgroups_folder = "test2"
delete_after_processing=True

# To keep track of processed threads and forums
processed_threads = set()
processed_forums = set()

def decode_email_body(part):
    """Try to decode with UTF-8, fallback to other encodings if necessary."""
    encodings = ['utf-8', 'latin-1', 'ascii']
    for enc in encodings:
        try:
            return part.get_payload(decode=True).decode(enc)
        except (UnicodeDecodeError, AttributeError):
            continue
    return part.get_payload(decode=True)  # Return as bytes if all decoding attempts fail

def open_url_in_browser(url, delay=0.0):
    """Open the URL in a browser.  Currently Mac-specific."""
    subprocess.run(["open", url])
    time.sleep(delay)

def process(mail, email_id, index, total, delete_after_processing):
    """Process a single message."""
    print(f'{index}/{total}', end=' ')

    status, msg_data = mail.fetch(email_id, "(RFC822)")

    for response_part in msg_data:
        if isinstance(response_part, tuple):
            # Parse the bytes email into a message object
            msg = email.message_from_bytes(response_part[1])

            # Decode the email subject
            subject, encoding = decode_header(msg["Subject"])[0]
            if isinstance(subject, bytes):
                # If it's bytes, decode to str
                subject = subject.decode(encoding if encoding else "utf-8")
            # Decode email sender
            from_ = msg.get("From")
            #print(f"Subject: {subject} (message {index}/{total})")
            #print(f"From: {from_}")

            # Extract the message body
            body = ""
            if msg.is_multipart():
                for part in msg.walk():
                    if part.get_content_type() == "text/plain":
                        body = decode_email_body(part)
                        break
            else:
                body = decode_email_body(msg)

            # TODO: Factor forum/thread processing into single function.
            # Check for forum or thread subscription and process the first URL appropriately
            if "You are subscribed to the forum" in body:
                forum_urls = re.findall(r'https://www\.rcgroups\.com/forums/forumdisplay\.php\?f=\d+', body)
                if forum_urls:
                    url = forum_urls[0]
                    forum_id = url.split('=')[-1]
                    if forum_id not in processed_forums:
                        processed_forums.add(forum_id)
                        print(f'     F {url}')
                        open_url_in_browser(url)
                    else:
                        print(f'skip F {url}')
                    if delete_after_processing:
                        mail.store(email_id, '+FLAGS', '\\Deleted')
                    return  # Skip remaining text in the email after handling the first URL

            if "You are subscribed to the thread" in body:
                thread_urls = re.findall(r'https://www\.rcgroups\.com/forums/showthread\.php\?t=\d+', body)
                if thread_urls:
                    url = thread_urls[0]
                    thread_id = url.split('=')[-1]
                    if thread_id not in processed_threads:
                        processed_threads.add(thread_id)
                        full_url = f"{url}&goto=newpost"
                        print(f'     T {full_url}')
                        open_url_in_browser(full_url)
                    else:
                        print(f'skip T {url}')
                    if delete_after_processing:
                        mail.store(email_id, '+FLAGS', '\\Deleted')
                    return  # Skip remaining text in the email after handling the first URL


def main():
    """Main program."""
    mail = imaplib.IMAP4_SSL(imap_host)
    mail.login(username, app_password)

    mail.select(rcgroups_folder)

    # Search for all emails in the selected mailbox
    status, messages = mail.search(None, "ALL")

    # Convert messages to a list of email IDs
    email_ids = messages[0].split()
    total_emails = len(email_ids)

    print(f"Total emails in {rcgroups_folder}: {total_emails}")

    for index, email_id in enumerate(email_ids, start=1):
        process(mail, email_id, index, total_emails, delete_after_processing)

    # Close the connection and logout
    mail.close()
    mail.logout()

if __name__ == "__main__":
    main()
