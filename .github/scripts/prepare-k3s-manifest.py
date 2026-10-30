#!/usr/bin/env python3
"""Preserve the verified native runtime contract while delivering new application code."""
import argparse,copy,json,subprocess
from pathlib import Path

WORKLOADS={'Fdeployment','Deployment','StatefulSet','DaemonSet'}
POD_FIELDS=['nodeSelector','affinity','topologySpreadConstraints','tolerations','priorityClassName','volumes','initContainers','securityContext','serviceAccountName','automountServiceAccountToken','runtimeClassName','schedulerName']

def patch_objects(value,plan,namespace='default',live=None):
 objects=value['items'] if value.get('kind')=='List' else [value]
 lookup={(t['namespace'],t['name']):t for t in plan['workloads']}
 for obj in objects:
  kind=obj.get('kind')
  if kind not in WORKLOADS:continue
  md=obj['metadata'];ns=md.get('namespace',namespace);key=(ns,md['name']);target=lookup.get(key)
  if not target:raise RuntimeError('Workload is absent from the native placement plan; review placement before delivery: '+ns+'/'+md['name'])
  spec=obj['spec']
  if kind=='Fdeployment':
   if target.get('delivery') in ['native-offline','native']:raise RuntimeError('Native OAuth backend must remain a Deployment; competing Fdeployment ownership is refused')
   spec['replicas']=target['replicas'];spec['priorityClassName']=target['priorityClassName']
   if target.get('resources'):spec['resources']=copy.deepcopy(target['resources'].get(md['name'],spec.get('resources',{})))
   spec['placement']=({'mode':'Distributed','spread':target['spread'],'nodeSelector':{'fabkli.ch/production':'true'}} if len(target['nodes'])==2 else {'mode':'SingleNode','node':target['nodes'][0]})
   continue
  if live is None:raise RuntimeError('Native workload requires its verified current runtime contract')
  current=live(kind,ns,md['name']);old=current['spec']
  if current['metadata'].get('namespace')!=ns or current['metadata']['name']!=md['name']:raise RuntimeError('Native workload identity mismatch')
  if spec.get('selector')!=old.get('selector'):raise RuntimeError('Delivery would change an immutable native workload selector')
  if kind!='DaemonSet':spec['replicas']=target['replicas']
  if kind=='Deployment':spec['strategy']=copy.deepcopy(old['strategy'])
  pod=spec['template']['spec'];oldpod=old['template']['spec']
  for field in POD_FIELDS:
   if field in oldpod:pod[field]=copy.deepcopy(oldpod[field])
   else:pod.pop(field,None)
  existing={c['name']:c for c in oldpod['containers']}
  incoming={c['name']:c for c in pod['containers']}
  if not set(incoming)<=set(existing):raise RuntimeError('New native containers require runtime-contract review')
  for name,container in incoming.items():
   for field in ['volumeMounts','volumeDevices','securityContext','resources']:
    if field in existing[name]:container[field]=copy.deepcopy(existing[name][field])
    else:container.pop(field,None)
  for name,container in incoming.items():
   names={e['name'] for e in container.get('env',[])}
   for env in existing[name].get('env',[]):
    if env['name']=='GOMEMLIMIT' and env['name'] not in names:container.setdefault('env',[]).append(copy.deepcopy(env))
  # Retain native helper sidecars; images/environment of declared app containers remain new.
  pod['containers']=[incoming.get(c['name'],copy.deepcopy(c)) for c in oldpod['containers']]
 return value

def main():
 p=argparse.ArgumentParser();p.add_argument('manifest');p.add_argument('output');a=p.parse_args()
 def kube(*argv):return subprocess.check_output(['kubectl',*argv],text=True,timeout=30)
 raw=kube('-n','f-operator-system','get','configmap','k3s-workload-placement','--ignore-not-found','-o','json').strip()
 value=json.loads(kube('create','--dry-run=client','-f',a.manifest,'-o','json'))
 if raw:
  plan=json.loads(json.loads(raw)['data']['plan.json']);uid=json.loads(kube('get','namespace','kube-system','-o','json'))['metadata']['uid']
  if uid!=plan.get('destinationClusterUID') or uid==plan.get('sourceClusterUID'):raise RuntimeError('Placement plan is not for the destination cluster')
  namespace=kube('config','view','--minify','-o','jsonpath={.contexts[0].context.namespace}').strip() or 'default'
  value=patch_objects(value,plan,namespace,lambda kind,ns,name:json.loads(kube('-n',ns,'get',kind,name,'-o','json')))
 Path(a.output).write_text(json.dumps(value)+'\n')
if __name__=='__main__':main()
