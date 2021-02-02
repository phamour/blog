import socket

if __name__ == '__main__':
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        try:
            s.connect(("0.0.0.0", 32470))
            print("Connection OK")
        except Exception as e:
            print("Connection KO")
            print(e)
