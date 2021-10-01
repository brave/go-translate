import requests
import time
from urllib.parse import urlparse

def get_response_attempt(dest, data):
  response = requests.request(method='POST',
                              url=dest,
                              timeout=5,
                              headers={
                                'Content-Type': 'application/x-www-form-urlencoded',
                                'charset': 'utf-8',
                                'Cache-Control': 'no-cache'
                              },
                              data=data)
  return response.content

url = 'https://translate-relay.brave.com/translate?sl=es&tl=de'
data = {
  'q': [
    'Machine learning is the study of computer algorithms that can improve' \
    ' automatically through experience and by the use of data.',
    'It is seen as a part of artificial intelligence. Machine learning' \
    ' algorithms build a model based on sample data, known as "training' \
    ' data", in order to make predictions or decisions without being' \
    ' explicitly programmed to do so.',
    'Machine learning algorithms are used in a wide variety of applications,' \
    ' such as in medicine, email filtering, speech recognition, and computer' \
    ' vision, where it is difficult or unfeasible to develop conventional' \
    ' algorithms to perform the needed tasks.'
  ]
}

print(get_response_attempt(url, data))
