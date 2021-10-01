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

url = 'https://translate-relay.brave.com/translate?sl=auto&tl=en'
data = {
  'q': [
    'Il fumo di tabacco è la causa più comune della condizione insieme ad' \
    ' altri fattori, quali l\'inquinamento dell\'aria e la genetica, che' \
    ' rivestono un ruolo minore.',
    'O piloto inglês Lewis Hamilton da McLaren conquistou seu primeiro' \
    ' campeonato de pilotos com a diferença final de um ponto para o' \
    ' brasileiro Felipe Massa da Ferrari, que ficou com a segunda' \
    ' colocação, enquanto seu companheiro de equipe Kimi Räikkönen ficou' \
    ' com a terceira colocação.',
    'В последние дни на британских заправках все чаще выстраиваются очереди.' \
    ' Неожиданный дефицит бензина и дизеля отразился и на футболе, поставив' \
    ' под угрозу срыва сотни матчей в десятках низших лиг.',
    'Jaén es una ciudad y municipio español de la comunidad autónoma de' \
    ' Andalucía, capital de la provincia homónima. Ostenta el título de «Muy' \
    ' Noble y Muy Leal Ciudad de Jaén, Guarda y Defendimiento de los Reinos' \
    ' de Castilla» y es conocida como la «capital del Santo Reino».',
    'Kõige elementaarsemal tasemel asendab masintõlge ühe keele sõnad teise' \
    ' keele omadega, kuid sellest ei piisa heaks tõlkeks, sest tuleb tunda' \
    ' ära terved fraasid ja leida neile teises keeles vasted.',
    'Word-sense disambiguation concerns finding a suitable translation when' \
    ' a word can have more than one meaning.',
  ]
}

print(get_response_attempt(url, data))
