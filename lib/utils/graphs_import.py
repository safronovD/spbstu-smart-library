import requests
import json
import os
import subprocess
import argparse
import time

def login(host, user=None, passw=None):
  session = requests.session()  
  if user:
    session.auth = (user, passw)
  session.headers.update({'Content-Type': 'application/json', 'kbn-xsrf': 'true'})

  return session

def waiting_for_kibana(session, host):
  while True:
    try:
      resp = session.get(host)
      if resp.status_code == 200:
        print("Kibana is ready")
        break
    except:
        pass
    print("Failed to connect. Trying again...")
    time.sleep(5)

def create_index_pattern(session, host, file_path):
  json_file = os.path.join(file_path, '..', 'resources', 'new_index_pattern.json')
  with open(json_file, 'r') as file:
    index_json = json.load(file)
    response = session.post(host + '/api/saved_objects/index-pattern', json=index_json)
    print(response.text)
        

def import_graphs(session, host, json_files_list, index_pattern_name):
  response = session.get(host + '/api/saved_objects/_find?type=index-pattern&search_fields=title&search={}*'
                         .format(index_pattern_name))
  print(response.text)
  index_id = response.json().get('saved_objects')[0].get('id')
  for filename in json_files_list:
    with open(filename, 'r') as file:
      graph_json = json.load(file)
      graph_json.get('references')[0].update({'id': index_id})
      response = session.post(host + '/api/saved_objects/visualization/' + filename[0:-5], json=graph_json)
      print(response.text)


def export_graphs(session, host):
  response = session.get(host + '/api/saved_objects/_find?type=visualization&search_fields=title&search=Library*')
  for obj in response.json().get('saved_objects'):
    with open(obj.get('attributes').get('title').lower().replace(' ', '_') + '.json', 'w') as file:
      template = {'attributes': obj.get('attributes'), 'references': obj.get('references')}
      json.dump(template, file)
  print(response.json())


if __name__ == "__main__":
  file_path = os.getcwd()
  host = 'http://kibana:5601'
  os.chdir(os.path.join(file_path, '..', 'resources', 'visualization'))
  json_list = next(os.walk(os.getcwd()))[2]
  current_session = login(host)
  waiting_for_kibana(current_session, host)
  create_index_pattern(current_session, host, file_path)
  import_graphs(current_session, host, json_list, 'new_data')
