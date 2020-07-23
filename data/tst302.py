

#!/usr/bin/python

url='http://url5878.vsobr.com/wf/click?upn=1HRMNt-2FVhl6L095IYnIbmvf-2Fzne19ubxd5-2FaqgqNN14s74ToAIJMkTQAvK-2B-2FbSu6FPS6gK0-2FLyHMoXScCmZp5oSet8W5jnReHpeU0op6GS8-3D_SuJVuXlBI1KDlGJOfEpQFRaFI4g-2Bl-2BgWLqTpj6kKibeWuObxjrMAkUPkyN0-2BurewIuAPvRvyvtdPsRkiIRmUjk6kk2kHCprvuYwhA1qEGB7Dt-2BMBLAJk9VnUxz3kCFEoVgK4CMJeaKnMzhnN5-2FcBvl-2FMCAdjd7BYslYbpM6Yi1Q730Vi-2B4UNTtHic-2Fv1S-2BRm7lLuCWQyR1eaYOXS-2FindrA-3D-3D'


def x1():
     import subprocess
     process = subprocess.Popen(['curl', '-Is', url], stdout=subprocess.PIPE)
     stdout, _ = process.communicate()
     return stdout.replace('\r','').split('\n')

def x2():
     import pycurl
     headers=[]
     def x2collect(s):
          headers.append(s.decode('iso-8859-1').strip())
     c = pycurl.Curl()
     c.setopt(c.URL, url)
     c.setopt(c.HEADERFUNCTION, x2collect)
     c.perform()
     c.close()
     return headers

def x3():
     import urllib2
     class NoRedirection(urllib2.HTTPErrorProcessor):
         def http_response(self, request, response):
             return response
         https_response = http_response
     opener = urllib2.build_opener(NoRedirection)
     urllib2.install_opener(opener)
     req=urllib2.Request(url)
     response = urllib2.urlopen(req)
     ret = []
     for line in response.info().headers:
         ret.append(line.strip())
     return ret

def x4():
     import httplib
     parts = url.split('/')
     site = parts[2]
     page = '/'+'/'.join(parts[3:])
     conn = httplib.HTTPConnection(site)
     conn.request('HEAD', page)
     rl = conn.getresponse()
     return rl.getheaders()
     return rl.getheader('location')

from pprint import pprint
print '----'
pprint(x1())
print '----'
pprint(x2())
print '----'
pprint(x3())
print '----'
pprint(x4())
