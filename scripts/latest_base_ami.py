#
# Usage: latest_base_ami.py <REGION>
# Ex: latest_base_ami.py us-west-2
#
import boto3
import sys


REGION = sys.argv[1]


def main():
    session = boto3.Session(
        region_name=REGION,
    )
    ec2 = session.client('ec2')
    images = ec2.describe_images(
        Filters=[
            {
                'Name': 'name',
                'Values': ['Kubernetes Base - Debian -*'],
            },
            {
                'Name': 'virtualization-type',
                'Values': ['hvm'],
            },
            {
                'Name': 'architecture',
                'Values': ['x86_64'],
            },
        ],
    )['Images']
    if not images:
        print("You do not have any Kubernetes Base image. Run 'make build-ami' to create one.")
        exit(1)
    base_images = sorted(images, key=lambda x: x['CreationDate'], reverse=True)
    print(base_images[0]['ImageId'])

if __name__ == "__main__":
    main()
