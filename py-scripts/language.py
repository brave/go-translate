import requests
import time
from urllib.parse import urlparse
import random

def get_response_attempt(dest, data):
  response = requests.request(method='GET',
                              url=dest,
                              timeout=5,
                              headers={
                                'Content-Type': 'application/x-www-form-urlencoded',
                                'charset': 'utf-8',
                                'Cache-Control': 'no-cache'
                              },
                              data=data)
  return response.content

url = 'https://translate-relay.brave.com/language'
print(get_response_attempt(url, {}))
