import json
from uuid import UUID
from fastapi import FastAPI, HTTPException
from cache import r


app = FastAPI()

users = {
    "12345": {
        "id": "1",
        "email": "asd@a.dk",
        "token": "12345"
    },
    "67890": {
        "id": "1",
        "email": "f@a.dk",
        "token": "12345"
    }
}


messages = {
    "1": [
        {"id": "6b653741-5418-48fe-9bc6-5b989e8a8401", "message": "m2", "access": "*", "created_at": 1650450307},
        {"id": "23576d4e-e510-4141-973e-24ab16c82996", "message": "m1", "access": "*", "created_at": 1650450306},
        {"id": "0a78b2b5-402f-4bac-823c-feef04cf6c80", "message": "m3", "access": "*", "created_at": 1650450308},
        {"id": "5476d32e-6cbd-414c-a4c1-e9755492cc41", "message": "m4", "access": "*", "created_at": 1650450309},
    ]
}


@app.on_event("startup")
async def _():
    r.set("foo", "bar")
    

@app.get("/")
async def index():
    return "(◡ ‿ ◡ ✿)"


@app.get("/redis-test")
async def _():
    return r.get("foo")


@app.get("/provider/{id}/from/{last_message_id}/limit/{limit}/token/{token}")
async def _(id, last_message_id, limit: int, token):
    if not limit:
        raise HTTPException(status_code=400, detail="Limit is 0")
    if token not in users:
        raise HTTPException(status_code=400, detail="Invalid token")
    try:
        messages[id]
    except:
        raise HTTPException(status_code=204)

    # sorts the list after timestamp
    sorted_messages = sorted(messages[id], key=lambda x: x["created_at"])
    print(sorted_messages)

    # Takes the index of the first last_message_id 
    offset = 0
    if last_message_id:
        offset = [i for i, _ in enumerate(sorted_messages) if _["id"] == last_message_id][0]
        print(offset)
    # to make sure the limit also happens with offset
    limit = limit + offset

    if msg := messages[id][offset:limit]:
        return json.dumps(msg)
    raise HTTPException(status_code=204, detail="No messages")

