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
            "name": "pi",
            "image": "perl",
            "command": ["perl", "-Mbignum=bpi", "-wle", "print bpi(50)"],
        },
    ],
}


def main():
    methods = {
        "create": create,
        "delete": delete,
        "fetch": fetch,
        "create_scheduled": create_scheduled,
        "delete_scheduled": delete_scheduled,
        "fetch_scheduled": fetch_scheduled,
        "update_scheduled": update_scheduled,
    }
    methods[METHOD]()

def create():
    """
    Create a SimpleJob
    """
    r = requests.post(
        f"{URL}{api_group}/simple/",
        data=json.dumps(simple_job),
    )
    print(r.text)

def fetch():
    """
    Fetch a SimpleJob
    """
    r = requests.get(
        f"{URL}{api_group}/simple/{simple_job['namespace']}/{simple_job['name']}",
    )
    print(r.text)

def delete():
    """
    Delete a SimpleJob
    """
    r = requests.delete(
        f"{URL}{api_group}/simple/{simple_job['namespace']}/{simple_job['name']}",
    )
    print(r.text)

def create_scheduled():
    """
    Create a scheduled SimpleJob
    """
    r = requests.post(
        f"{URL}{api_group}/scheduled/",
        data=json.dumps(simple_job_cron),
    )
    print(r.text)

def fetch_scheduled():
    """
    Fetch a scheduled SimpleJob
    """
    r = requests.get(
        f"{URL}{api_group}/scheduled/{simple_job_cron['namespace']}/{simple_job_cron['name']}",
    )
    print(r.text)

def update_scheduled():
    """
    Update a scheduled SimpleJob
    """
    simple_job_cron["cron"] = "9 9 * * *"
    simple_job_cron["maxRetries"] = 99

    r = requests.put(
        f"{URL}{api_group}/scheduled/{simple_job_cron['namespace']}/{simple_job_cron['name']}",
        data=json.dumps(simple_job_cron),
    )
    print(r.text)

def delete_scheduled():
    """
    Delete a scheduled SimpleJob
    """
    r = requests.delete(
        f"{URL}{api_group}/scheduled/{simple_job_cron['namespace']}/{simple_job_cron['name']}",
    )
    print(r.text)


if __name__ == "__main__":
    main()
