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

simple_job_cron = {
    "name": "pi",
    "namespace": "default",
    "maxRetries": 1,
    "cron": "* * * * *",
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
    methods[METHOD]()

def create():
    """
    Create a SimpleJob
    """
    r = requests.post(
        f"{URL}{api_group}/",
        data=json.dumps(simple_job),
    )
    print(r.text)

def update():
    """
    Create & Update a scheduled SimpleJob
    """
    r = requests.post(
        f"{URL}{api_group}/",
        data=json.dumps(simple_job_cron),
    )
    print(r.text)

    simple_job_cron["cron"] = "9 9 * * *"
    simple_job_cron["maxRetries"] = 99

    r = requests.put(
        f"{URL}{api_group}/",
        data=json.dumps(simple_job_cron),
    )
    print(r.text)

def delete():
    """
    Delete a SimpleJob
    """
    r = requests.delete(
        f"{URL}{api_group}/",
        data=json.dumps(simple_job),
    )
    print(r.text)


if __name__ == "__main__":
    main()
