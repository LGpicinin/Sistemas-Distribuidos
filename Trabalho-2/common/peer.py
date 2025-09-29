from Pyro5.api import Daemon, expose, Proxy, oneway
from Pyro5.core import locate_ns
from Pyro5 import nameserver
from .static import NAMES
from .states import States
from datetime import datetime
from time import time, sleep
from threading import Thread

class Peer:
    
    def __init__(self, name: str):
        self.name = name
        self.daemon = None
        self.uri = None
        self.ns = None
        self.state = States.RELEASED
        self.last_request_timestamp = None
        self.request_queue = []
        self.active_peers = {}
        self.time_peers = {}
        self.response_peers = {}

        self.create_daemon()
        self.register_on_ns(name)

        sleep(30)

        for n in NAMES:
            if n != name:
                peer: Peer = Proxy(f"PYRONAME:{n}")
                self.active_peers[n] = peer
                self.time_peers[n] = time() + 60


    def create_daemon(self) -> None:
        self.daemon = Daemon()
        self.uri = self.daemon.register(self)
    
    def wait_response(self, peer): # necessário rever
        peer = self.tranfer_proxy_to_this_thread(peer)
        self.response_peers[peer.name] = peer.request_resource(self, self.last_request_timestamp)
    
    def tranfer_proxy_to_this_thread(self, proxy):
        proxy._pyroClaimOwnership()
        return proxy
    
    @expose
    @oneway
    def send_heartbeat(self) -> None:
        last_pulse = time() - 60
        while True:
            if time() >= last_pulse + 30:
                self.verificate_heartbeat()

            for peer_name, peer in self.active_peers.items():
                peer = self.tranfer_proxy_to_this_thread(peer)
                peer.receive_heartbeat(self.name)

        last_pulse = time()

    @expose
    @oneway       
    def receive_heartbeat(self, name) -> None:
        self.time_peers[name] = time() + 60


    def verificate_heartbeat(self) -> None:
        peers_to_remove = []
        for name, peer_time in self.time_peers.items():
            if peer_time < time():
                peers_to_remove.append(name)

        for name in peers_to_remove:
            del self.active_peers[name]
            del self.time_peers[name]


    def get_resource(self) -> None:
        self.state = States.WANTED
        self.last_request_timestamp = datetime.now()

        peers_to_remove = []
        

        for peer_name, peer in self.active_peers.items():
            try:
                self.response_peers[peer_name] = 0
                thread_wait = Thread(target=self.wait_response, args=(peer)) # necessário rever
                thread_wait.start()
                sleep(15)
                if self.response_peers[peer_name] == 0:
                    peers_to_remove.append(peer_name)
                else:
                    if self.response_peers[peer_name] == False:
                        for name in peers_to_remove:
                            del self.active_peers[name]
                            del self.time_peers[name]
                        return
            except:
                print(f"Tentativa de requisitar recurso para {peer_name} foi mal sucedida")

        for name in peers_to_remove:
            del self.active_peers[name]
            del self.time_peers[name]
        
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
        
        # verifica se processo esta ativo
        while next_peer.name not in self.active_peers.keys() and len(self.active_peers.keys()) != 0:
            next_peer, timestamp = self.request_queue.pop(0)
            
        if next_peer.name in self.active_peers.keys():
            next_peer = self.tranfer_proxy_to_this_thread(next_peer)
            next_peer.receive_resource()
            self.request_queue = []

    @expose
    @oneway
    def receive_resource(self) -> None:
        self.state = States.HELD


    def register_on_ns(self, name: str) -> None:
        self.ns : nameserver.NameServer = locate_ns()
        self.ns.register(name, self.uri)


    def wait_requests(self) -> None:
        self.daemon.requestLoop()


    def list_active_peers(self) -> None:
        print("Peers ativos:")
        for name, fail_time in self.time_peers.items():
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
                    self.menu()


    def run(self) -> None:
           thread_menu = Thread(target=self.menu)
           thread_hb = Thread(target=self.send_heartbeat)
           thread_requests = Thread(target=self.wait_requests)

           thread_hb.start()
           thread_menu.start()
           thread_requests.start()

