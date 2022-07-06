#!/bin/bash
eval "$(jq -r '@sh "export run_script=\(.run_script) cluster_name=\(.cluster_name) region=\(.region) project=\(.project)"')"

if [ $run_script == "false" ];
then
  jq -n '{}'
  exit 0
fi
gcloud container clusters get-credentials ${cluster_name} --region ${region} --project ${project} &> /dev/null

sleep 20

until kubectl get svc ha-es-http &> /dev/null
do
  sleep 1
done

until [ "$(kubectl get svc ha-es-http | sed -n 2p | awk '{print $4}')" != "<pending>" ]
do
    sleep 1
done
es_ip=$(kubectl get svc ha-es-http | sed -n 2p | awk '{print $4}')


until kubectl get svc ha-kb-http &> /dev/null
do
  sleep 1
done

until [ "$(kubectl get svc ha-kb-http | sed -n 2p | awk '{print $4}')" != "<pending>" ]
do
    sleep 1
done
kibana_ip=$(kubectl get svc ha-kb-http| sed -n 2p | awk '{print $4}')

until kubectl get secret ha-es-elastic-user -o go-template='{{.data.elastic | base64decode}}' &> /dev/null
do
    sleep 1
done
es_password=$(kubectl get secret ha-es-elastic-user -o go-template='{{.data.elastic | base64decode}}')

jq -n \
    --arg es_ip "https://$es_ip:9200" \
    --arg es_password "$es_password" \
    --arg kibana_ip "https://$kibana_ip:5601" \
    '{"es_ip":$es_ip,"es_password":$es_password,"kibana_ip":$kibana_ip}'