"""Stub server for the CI crawler-timing benchmark.

Answers every request with one fixed JSON body containing every field any
Thunder ID or backend endpoint reads (access_token, expires_in, id,
failed_jobs), so CrawlerManager can get past _start_run and _flush_batch
without a real backend or Thunder ID instance. This is only meant to
measure how long the 5 crawlers take to scrape their real external job
sites — it is not a behavioral stand-in for the real services.
"""

import json
import sys
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

STUB_BODY = json.dumps(
    {
        "access_token": "stub-token",
        "expires_in": 3600,
        "id": 1,
        "failed_jobs": [],
    }
).encode("utf-8")


class StubHandler(BaseHTTPRequestHandler):
    def _respond(self):
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(STUB_BODY)))
        self.end_headers()
        self.wfile.write(STUB_BODY)

    def do_GET(self):
        self._respond()

    def do_POST(self):
        self._respond()

    def log_message(self, format, *args):
        pass  # keep CI logs focused on the crawler, not stub request noise


if __name__ == "__main__":
    port = int(sys.argv[1]) if len(sys.argv) > 1 else 9000
    server = ThreadingHTTPServer(("127.0.0.1", port), StubHandler)
    print(f"Stub server listening on 127.0.0.1:{port}")
    server.serve_forever()
