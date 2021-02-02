from socketserver import StreamRequestHandler, ThreadingTCPServer
import threading

from flask import Flask


class DummySocketHandler(StreamRequestHandler):
    def handle(self):
        with open("socket.log", "a") as f:
            f.write("REQ\n")


def init():
    socket_server = ThreadingTCPServer(("0.0.0.0", 32470),
                                       DummySocketHandler)
    threading.Thread(target=socket_server.serve_forever).start()


app = Flask(__name__)
app.debug = True

@app.route('/')
def home():
    print("Hey", flush=True)
    return "Hey"


if __name__ == '__main__':
    init()
    app.run(host="0.0.0.0", port=8888)
