import requests
import json
import argparse
import urllib3
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
import subprocess
import sys

from sys import exit
from os import listdir
from os.path import isfile, join
from requests.auth import HTTPBasicAuth

# Create the parser
my_parser = argparse.ArgumentParser(description='Install Elasticsearch watchers for k8s monitoring')

# Add the arguments
my_parser.add_argument('-elasticsearch-host',
                       '-es',
                       type=str,
                       nargs="?",
                       const="https://localhost:9200",
                       default="https://localhost:9200",
                       help='The url of Elasticsearch Watcher API')
my_parser.add_argument('-username',
                       '-u',
                       type=str,
                       nargs="?",
                       const="elastic",
                       default="elastic",
                       help='The username to access the ES API')
my_parser.add_argument('-password',
                       '-p',
                       type=str,
                       nargs="?",
                       const="changeme",
                       default="changeme",
                       help='The password to access the ES API')

args = my_parser.parse_args()

watchers_path = "./watchers"
watchers = [f for f in listdir(watchers_path) if isfile(join(watchers_path, f))]
print("Watchers to be installed: " + str(watchers))

base_url=args.elasticsearch_host+"/_watcher/watch/"
headersList = {"Content-Type": "application/json"}

for watcher in watchers:
    # Create the watcher

    f = open('./watchers/'+watcher)
    watcher_body = json.load(f)

    watcher_name = watcher.split(".json")[0]
    url = base_url+watcher_name+"?pretty"
    res = requests.post(
        url+"",
        headers=headersList,
        verify=False,
        json=watcher_body,
        auth = HTTPBasicAuth(args.username, args.password)
    )
    if res.status_code not in [200, 201]:
        print("Status code is: " + str(res.status_code))
        print("Error- Installing Watcher")
        print(res.text)
        exit(1)