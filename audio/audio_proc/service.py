import logging
from concurrent import futures

import grpc

from .audio_proc_pb2 import FormatAndCompressInput, FormatAndCompressOutput
from .audio_proc_pb2_grpc import CompressionServicer, add_CompressionServicer_to_server

logging.basicConfig(level=logging.DEBUG)
log = logging.getLogger("server")


class Compression(CompressionServicer):
    def FormatAndCompress(self, request: FormatAndCompressInput, context) -> FormatAndCompressOutput:
        log.debug("FormatAndCompress")
        # do things
        return FormatAndCompressOutput(content=request.Content)


def build_server(port: str):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    cs = CompressionServicer()
    add_CompressionServicer_to_server(cs, server)
    server.add_insecure_port(port)
    log.debug(f"built server for port '{port}'")
    return server


if __name__ == '__main__':
    srv = build_server("[::]:50051")
    log.debug("listening to calls")
    srv.start()
    srv.wait_for_termination()
