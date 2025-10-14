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
        self.who_has_resource = None

        self.create_daemon()
        self.register_on_ns(name)

        sleep(15)

        for n in NAMES:
            if n != name:
                peer: Peer = Proxy(f"PYRONAME:{n}")
                self.active_peers[n] = peer
                self.time_peers[n] = time() + 10
                self.response_peers[n] = True
                


    def create_daemon(self) -> None:
        self.daemon = Daemon()
        self.uri = self.daemon.register(self)


    def tranfer_proxy_to_this_thread(self, proxy):
        proxy._pyroClaimOwnership()
        return proxy
    
    
    def wait_response(self, peer_name) -> None:
        peer: Peer = Proxy(f"PYRONAME:{peer_name}")
        print(f'Esperando resposta do Peer {peer_name}')
        self.response_peers[peer_name] = peer.request_resource(self.name, self.last_request_timestamp)

    
    @expose
    @oneway
    def send_heartbeat(self) -> None:
        last_pulse = time() - 60
        while True:
            if time() >= last_pulse + 5:
                self.verificate_heartbeat()
                last_pulse = time()

            for peer_name, peer in self.active_peers.items():
                peer = self.tranfer_proxy_to_this_thread(peer)
                try:
                    peer.receive_heartbeat(self.name)
                    if peer.has_resource():
                        self.who_has_resource = peer_name
                except:
                    pass
                    

    @expose
    @oneway
    def has_resource(self):
        return self.state == States.HELD


    @expose
    @oneway       
    def receive_heartbeat(self, name) -> None:
        self.time_peers[name] = time() + 5


    def verificate_heartbeat(self) -> None:
        peers_to_remove = []
        for name, peer_time in self.time_peers.items():
            if peer_time < time():
                peers_to_remove.append(name)

        for name in peers_to_remove:
            print(f'Peer {name} removido')    
            del self.active_peers[name]
            del self.time_peers[name]
            del self.response_peers[name]
            print(f'{name=}, {self.who_has_resource=}, {self.state=}, {self.response_peers}')
            if self.state == States.WANTED and all(self.response_peers.values()):
                self.receive_resource()


    def get_resource(self) -> None:
        self.state = States.WANTED
        self.last_request_timestamp = datetime.now()

        peers_to_remove = []

        consegui_recurso = True

        for peer_name, peer in self.active_peers.items():
            try:
                thread_wait = Thread(target=self.wait_response, kwargs={"peer_name": peer_name})
                thread_wait.start()
                thread_wait.join(15)

                if thread_wait.is_alive():
                    peers_to_remove.append(peer_name)
                else:
                    if self.response_peers[peer_name] == False:
                        peer = self.tranfer_proxy_to_this_thread(peer)
                        if peer.has_resource():
                            self.who_has_resource = peer_name
                        consegui_recurso = False
            except Exception as E:
                print(E)
                print(f"Tentativa de requisitar recurso para {peer_name} foi mal sucedida")

        for name in peers_to_remove:
            print(f'Peer {name} removido')
            del self.active_peers[name]
            del self.time_peers[name]
            del self.response_peers[name]
            if self.state == States.WANTED and all(self.response_peers.values()):
                self.receive_resource()

        if consegui_recurso == True:
            print('Recurso obtido')
            self.receive_resource()


    @expose
    def request_resource(self, who, timestamp: datetime) -> None:
        timestamp = datetime.fromisoformat(timestamp)
        
        if (self.last_request_timestamp == None):
            return True
        
        if (
            self.state == States.HELD or 
            (
                self.state == States.WANTED and 
                self.last_request_timestamp < timestamp
            )
        ):
            self.request_queue.append((who, timestamp))
            return False
        
        return True


    def free_resource(self) -> None:
        if self.state != States.HELD:
            return
        
        self.state = States.RELEASED
        
        if len(self.request_queue):
            peer_name, timestamp = self.request_queue.pop(0)
            next_peer: Peer = Proxy(f"PYRONAME:{peer_name}")
            
            # verifica se processo esta ativo
            while peer_name not in self.active_peers.keys() and len(self.active_peers.keys()) != 0:
                peer_name, timestamp = self.request_queue.pop(0)
                next_peer: Peer = Proxy(f"PYRONAME:{peer_name}")
            
            # entrega recurso para o primeiro peer da fila
            if peer_name in self.active_peers.keys():
                next_peer.receive_resource()

            # manda ok para os outros peers
            print("Liberando request queue")
            for peer_name, timestamp in self.request_queue:
                if peer_name in self.active_peers.keys():
                    try:
                        print(f"Dando ok para o {peer_name}")
                        peer: Peer = Proxy(f"PYRONAME:{peer_name}")
                        peer.receive_ok(self.name)
                    except:
                        pass
            
            self.request_queue = []

            
        print("Recurso liberado")


    @expose
    @oneway
    def receive_resource(self) -> None:
        print('Recurso recebido')
        self.state = States.HELD
        Thread(target=lambda: (sleep(10) or self.free_resource())).start()
        self.who_has_resource = self.name

    @expose
    @oneway
    def receive_ok(self, peer_name : str) -> None:
        print(f"Recebendo ok de {peer_name}")
        self.response_peers[peer_name] = True
        print(f"Response queue atualizada: {self.response_peers.values()}")


    def register_on_ns(self, name: str) -> None:
        self.ns : nameserver.NameServer = locate_ns()
        self.ns.register(name, self.uri)


    def wait_requests(self) -> None:
        self.daemon.requestLoop()


    def list_active_peers(self) -> None:
        print("Peers ativos:")
        for name, fail_time in self.time_peers.items():
            print(f"\t{name}: Último heartbeat em {datetime.fromtimestamp(fail_time - 60).isoformat()}")
            print()


    def menu(self) -> None:
        while True:
            print(f"Estado Atual: {self.state}")
            print("Digite sua preferência:")
            print("\t1) Requisitar recurso")
            print("\t2) Liberar recurso")
            print("\t3) Listar peers ativos")
            print("\t0) Sair")
            print()
            option = input()

            match option:
                case "1":
                    if self.state == States.HELD:
                        print("Você já está com o recurso.")
                    elif self.state == States.WANTED:
                        print("Recurso já foi requisitado, aguarde.")
                    else:
                        self.get_resource()

                case "2":
                    if self.state != States.HELD:
                        print("Não é possível liberar recurso.")
                    else:
                        self.free_resource()

                case "3":
                    self.list_active_peers()

                case "0":
                    exit(0)
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

