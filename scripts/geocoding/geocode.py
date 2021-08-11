#!/usr/bin/env python
import json
import os
import re
from io import BytesIO

import google.cloud.storage as gs
import requests


def env(key: str) -> str:
    val = os.getenv(key, None)
    if val is None:
        raise ValueError(f"env variable '{key}' is not set.")
    return val


api_key = env("GMAPS_API_KEY")
storage_sa = env("STORAGE_SERVICE_ACCOUNT")


def geocode(q: str) -> requests.Response:
    url = "https://maps.googleapis.com/maps/api/geocode/json"
    return requests.get(url, params=dict(address=q, key=api_key))


if __name__ == '__main__':
    client = gs.client.Client.from_service_account_json(storage_sa)
    bucket = gs.bucket.Bucket(client, "talkiewalkie-dev")
    blob = bucket.blob("scraping/audiotours/mywowo-full-run.jl")
    with BytesIO() as fi:
        blob.download_to_file(fi, client)

        fi.seek(0)
        mywowo = [json.loads(line) for line in fi]

    located = []
    for i, tour in enumerate(mywowo):
        if i > 0 and i % 10 == 0:
            print(f"{i}th tour processed")

        try:
            url = tour["tour_url"]
            match = re.match(r"https://mywowo.net/en/([^/]*)/([^/*]*)/([^/]*).*", url)
            query = f"{match.group(3)} {match.group(2)} {match.group(1)}".replace("-", " ")
        except:
            print(f"failed to build query for tour #{i}: '{url}'")
            continue

        try:
            geocoded = geocode(query)
            out = geocoded.json()
            result = out["results"][0]
            located.append(
                {**tour, "location": result["geometry"]["location"], "gmaps-output": result})
        except:
            print(
                f"failed to geocode for query: '{query}' - {geocoded.status_code}, #results={len(out.get('results', []))}")
            continue

    out_blob = bucket.blob("scraping/audiotours/mywowo-full-run-located.jl")
    with BytesIO() as fo:
        fo.write("\n".join(json.dumps(t) for t in located).encode())
        fo.seek(0)
        out_blob.upload_from_file(fo)
