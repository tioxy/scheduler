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
        # Simple
        "create": create,
        "delete": delete,
        "fetch": fetch,
        "list_all": list_all,
        "list_default": list_default,
        
        # Scheduled
        "create_scheduled": create_scheduled,
        "delete_scheduled": delete_scheduled,
        "fetch_scheduled": fetch_scheduled,
        "list_all_scheduled": list_all_scheduled,
        "list_default_scheduled": list_default_scheduled,
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

def list_all():
    """
    List all SimpleJobs
    """
    r = requests.get(
        f"{URL}{api_group}/simple/",
    )
    print(r.text)

def list_default():
    """
    List SimpleJobs from namespace 'default'
    """
    r = requests.get(
        f"{URL}{api_group}/simple/default",
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

def list_all_scheduled():
    """
    List all scheduled SimpleJobs
    """
    r = requests.get(
        f"{URL}{api_group}/scheduled/",
    )
    print(r.text)

def list_default_scheduled():
    """
    List scheduled SimpleJobs from namespace 'default'
    """
    r = requests.get(
        f"{URL}{api_group}/scheduled/default",
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


if __name__ == "__main__":
    main()
