from Pyro5.api import Daemon, expose, Proxy, oneway
from Pyro5.core import locate_ns
from Pyro5 import nameserver
from .static import NAMES
from .states import States
from datetime import datetime
from time import time
from threading import Thread

class Peer:
    
    def __init__(self, name: str):
        self.daemon = None
        self.uri = None
        self.ns = None
        self.state = States.RELEASED
        self.last_request_timestamp = None
        self.request_queue = []
        self.active_peers = {}

        for n in NAMES:
            if n != name:
                self.active_peers[n] = time() + 60

        self.create_daemon()
        self.register_on_ns(name)

    def create_daemon(self) -> None:
        self.daemon = Daemon()
        self.uri = self.daemon.register(self)


    @oneway
    def send_heartbeat(self) -> None:
        last_pulse = time() - 60
        while True:
            if time() >= last_pulse + 30:
                self.verificate_heartbeat()

            for peer_name, peer_time in self.active_peers.items():
                try:
                    peer: Peer = Proxy(f"PYRONAME:{peer_name}")
                    peer.receive_heartbeat(self.name)
            
                except:
                    pass
        last_pulse = time()


    @oneway       
    def receive_heartbeat(self, name) -> None:
        self.active_peers[name] = time() + 60


    def verificate_heartbeat(self) -> None:
        peers_to_remove = []
        for name, peer_time in self.active_peers.items():
            if peer_time < time():
                peers_to_remove.append(name)

        for name in peers_to_remove:
            del self.active_peers[name]


    def get_resource(self) -> None:
        self.state = States.WANTED
        self.last_request_timestamp = datetime.now()

        for peer_name, peer_time in self.active_peers.items():
            peer: Peer = Proxy(f"PYRONAME:{peer_name}")
            if peer.request_resource(self, self.last_request_timestamp) == False:
                break


        self.state = States.HELD


    @expose
    def request_resource(self, who, timestamp: datetime) -> None:
        if self.state == States.HELD or (self.state == States.WANTED and self.last_request_timestamp < timestamp):
            self.request_queue.append((who, timestamp))
            return False
        else:
            return True

    @expose
    def free_resouce(self) -> None:
        self.state = States.RELEASED
        next_peer, timestamp = self.request_queue.pop(0)
        next_peer.receive_resource(self.request_queue)
        self.request_queue = []


    @oneway
    def receive_resource(self, requests: list) -> None:
        self.state = States.HELD
        self.request_queue = requests


    def register_on_ns(self, name: str) -> None:
        self.ns : nameserver.NameServer = locate_ns()
        self.ns.register(name, self.uri)


    def wait_requests(self) -> None:
        self.daemon.requestLoop()


    def list_active_peers(self) -> None:
        print("Peers ativos:")
        for name, fail_time in self.active_peers.items():
            print(f"\t{name}: Último heartbeat em {fail_time - 60}")
            print()


    def menu(self) -> None:
        while True:
            print("Digite sua preferência:")
            print("\t1) Requisitar recurso")
            print("\t2) Liberar recurso")
            print("\t3) Listar peers ativos")
            print("\t0) Sair")
            print()
            option = input()

            match option:
                case "1":
                    self.get_resource()

                case "2":
                    self.free_resource()

                case "3":
                    self.list_active_peers()

                case "0":
                    break

                case _:
                    print("Opção inexistente")
                    menu()


    def run(self) -> None:
           thread_menu = Thread(target=self.menu)
           thread_hb = Thread(target=self.send_heartbeat)
           thread_requests = Thread(target=self.wait_requests)

           thread_hb.start()
           thread_menu.start()
           thread_requests.start()

