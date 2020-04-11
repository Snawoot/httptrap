#!/usr/bin/env python3

from http.server import HTTPServer, BaseHTTPRequestHandler
from socketserver import ThreadingMixIn
import threading
import urllib.parse
import random

MAX_BODY_LEN = 4096
LOGIN = b'admin'
PASSWORD = b'12345678'
TRAP_PROBABILITY = 0.01

class Handler(BaseHTTPRequestHandler):
    def fail_auth(self):
        if random.random() < TRAP_PROBABILITY:
            self.send_response(454)
            self.end_headers()
        else:
            self.send_response(403)
            self.end_headers()

    def read_post_form(self):
        body_len = self.headers['Content-Length']
        if (body_len is None or int(body_len) > 4096 or
            self.headers.get('Content-Type', '').lower() != 'application/x-www-form-urlencoded'):
            return {}
        try:
            body_len = int(body_len)
        except:
            return {}
        body = self.rfile.read(body_len)
        data = urllib.parse.parse_qs(body)
        return data

    def do_POST(self):
        data = self.read_post_form()
        if data.get(b'login', [b''])[0] == LOGIN and data.get(b'password', [b''])[0] == PASSWORD:
            self.send_response(200)
            self.send_header('Content-Type', 'text/plain')
            self.end_headers()
            self.wfile.write(b"You are in!\n")
        else:
            self.fail_auth()


class ThreadingSimpleServer(ThreadingMixIn, HTTPServer):
    pass

def main():
    try:
        server = ThreadingSimpleServer(('127.0.0.1', 8080), Handler)
        server.serve_forever()
    except KeyboardInterrupt:
        pass


if __name__ == '__main__':
    main()
