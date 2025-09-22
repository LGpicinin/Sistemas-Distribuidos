from Pyro5.api import Daemon, expose, Proxy, oneway
from Pyro5.core import locate_ns
from Pyro5 import nameserver
from .static import NAMES
from .states import States
from datetime import datetime
from typing import Self
from time import time

class Peer:
    
    def __init__(self, name: str):
        self.daemon = None
        self.uri = None
        self.ns = None
        self.state = States.RELEASED
        self.last_request_timestamp = None
        self.request_queue = []
        self.active_peers = []

        for n in NAMES:
            if n != name:
                self.active_peers.append((n, time()+60))
        
        self.create_daemon()
        self.register_on_ns(name)
    
    def create_daemon(self) -> None:
        self.daemon = Daemon()
        self.uri = self.daemon.register(self)


    @oneway
    def send_heartbeat(self) -> None:
        for peer_name, peer_time in self.active_peers:
            peer: Peer = Pyro5.api.Proxy(f"PYRONAME:{peer_name}")
            peer.receive_heartbeat(self.name)
        

    @oneway       
    def receive_heartbeat(self, name) -> None:
        for i in range(0, len(self.active_peers)):
            peer_name = self.active_peers[i][0]
            if name == peer_name:
                self.active_peers[i] = (name, time()+60)
                break


    def verificate_heartbeat(self) -> None: 
        for i in range(0, len(self.active_peers)):
            peer_time = self.active_peers[i][1]
            if peer_time < time():
                self.active_peers.pop(i)

    
    def get_resource(self) -> None:
        self.state = States.WANTED
        self.last_request_timestamp = datetime.now()

        for peer_name, peer_time in self.active_peers:
            peer: Peer = Pyro5.api.Proxy(f"PYRONAME:{peer_name}")
            if peer.request_resource(self, self.last_request_timestamp) == False:
                break
        
        
        self.state = States.HELD
    
    
    @expose
    def request_resource(self, who: Self, timestamp: datetime) -> None:
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
        
        
        
