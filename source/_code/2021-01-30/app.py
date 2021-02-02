from socketserver import StreamRequestHandler, ThreadingTCPServer
import threading

from flask import Flask


class DummySocketHandler(StreamRequestHandler):
    def handle(self):
        with open("socket.log", "a") as f:
            f.write("REQ\n")


class MyFlask(Flask):
    def __init__(self, *arg, **kwargs):
        super().__init__(*arg, **kwargs)
        try:
            print("Creating socket server")
            self.socket_server = ThreadingTCPServer(("0.0.0.0", 32470),
                                                    DummySocketHandler)
            threading.Thread(target=self.socket_server.serve_forever,
                             daemon=True).start()
            print("OK")
        except OSError:
            print('Unable to start socket server')
    

app = MyFlask(__name__)
app.debug = True


@app.route('/')
def home():
    print('hey', flush=True)
    return 'Hey'
