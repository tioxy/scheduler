#
# Usage: simplejob.py <METHOD> <HOST>
# Ex: simplejob.py create http://localhost:8080
#
import requests
import json
import sys

METHOD = sys.argv[1]
URL = sys.argv[2]


api_group = "/api/v1/jobs"

simple_job = {
    "name": "pi",
    "namespace": "default",
    "maxRetries": 4,
    "containers": [
        {
            "name": "pi",
            "image": "perl",
            "command": ["perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"],
        },
    ],
}

simple_job_updated = {
    "name": "pi",
    "namespace": "default",
    "maxRetries": 1,
    "containers": [
        {
            "name": "pi-new",
            "image": "perl",
            "command": ["perl", "-Mbignum=bpi", "-wle", "print bpi(50)"],
        },
    ],
}


def main():
    methods = {
        "create": create,
        "delete": delete,
        "update": update,
    }
    output = methods[METHOD]()
    print(output.text)

def create():
    r = requests.post(
        f"{URL}{api_group}/",
        data=json.dumps(simple_job),
    )
    return r

def update():
    r = requests.put(
        f"{URL}{api_group}/",
        data=json.dumps(simple_job_updated),
    )
    return r

def delete():
    r = requests.delete(
        f"{URL}{api_group}/",
        data=json.dumps(simple_job_updated),
    )
    return r


if __name__ == "__main__":
    main()
