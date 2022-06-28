#!/usr/bin/env python3

import os
import subprocess

SCRIPT_VERSION = 'v1'
OUT_DIR = f'assets/static/{SCRIPT_VERSION}'

GOOGLE_API_HOST = 'https://translate.googleapis.com/'
GOOGLE_EXTRA_HEADER = 'Google-Translate-Element-Mode: library'

ELEMENT_JS_DOWNLOAD_PATH = 'translate_a/element.js?cb=cr.googleTranslate.onTranslateElementLoad&aus=true&clc=cr.googleTranslate.onLoadCSS&jlc=cr.googleTranslate.onLoadJavascript&hl=en'
ELEMENT_JS_OUT_DIR = OUT_DIR
ELEMENT_JS_OUT_PATH = f'{ELEMENT_JS_OUT_DIR}/element.js'
os.makedirs(ELEMENT_JS_OUT_DIR, exist_ok=True)

subprocess.check_call([
    'curl', '-H', GOOGLE_EXTRA_HEADER,
    GOOGLE_API_HOST + ELEMENT_JS_DOWNLOAD_PATH, '-o', ELEMENT_JS_OUT_PATH
])

MAIN_JS_DOWNLOAD_PATH = '_/translate_http/_/js/k=translate_http.tr.en_US.oOC1Oa7Rttc.O/am=Bg/d=1/exm=el_conf/ed=1/rs=AN8SPfqPYBV0hk02iWIVCgyiPCEnQfgUdA/m=el_main'
MAIN_JS_OUT_DIR = f'{OUT_DIR}/js/element/'
MAIN_JS_OUT_PATH = MAIN_JS_OUT_DIR + 'main.js'

os.makedirs(MAIN_JS_OUT_DIR, exist_ok=True)
subprocess.check_call(
    ['curl', GOOGLE_API_HOST + MAIN_JS_DOWNLOAD_PATH, '-o', MAIN_JS_OUT_PATH])

CSS_DOWNLOAD_PATH = 'translate_static/css/translateelement.css'
CSS_OUT_DIR = f'{OUT_DIR}/css/'
CSS_OUT_PATH = CSS_OUT_DIR + 'translateelement.css'
os.makedirs(CSS_OUT_DIR, exist_ok=True)
subprocess.check_call(['curl', GOOGLE_API_HOST + CSS_DOWNLOAD_PATH, '-o', CSS_OUT_PATH])
