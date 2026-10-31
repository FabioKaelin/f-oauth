import copy,importlib.util,sys,unittest,json,tempfile
from unittest.mock import patch
from pathlib import Path
spec=importlib.util.spec_from_file_location('delivery',Path(__file__).with_name('prepare-k3s-manifest.py'))
delivery=importlib.util.module_from_spec(spec);spec.loader.exec_module(delivery)

class NativeDeliveryTests(unittest.TestCase):
 def setUp(self):
  self.target={'namespace':'apps','name':'app','replicas':2,'priorityClassName':'essential','nodes':['raspy5','raspy4'],'spread':'Required','resources':{'app':{'requests':{'cpu':'10m','memory':'32Mi'},'limits':{'cpu':'100m','memory':'128Mi'}}}}
  self.plan={'workloads':[self.target]}
  self.obj={'kind':'Deployment','metadata':{'namespace':'apps','name':'app'},'spec':{'selector':{'matchLabels':{'app':'app'}},'template':{'spec':{'containers':[{'name':'app','image':'new-image','env':[{'name':'CODE_SETTING','value':'new'}]}]}}}}
  self.live=copy.deepcopy(self.obj)
  self.live['spec']['strategy']={'type':'Recreate'}
  self.live['spec']['template']['spec'].update({'volumes':[{'name':'state','hostPath':{'path':'/verified/native'}}],'initContainers':[{'name':'guard','image':'verified'}],'nodeSelector':{'node':'raspy5'}})
  self.live['spec']['template']['spec']['containers'][0].update({'image':'old-image','volumeMounts':[{'name':'state','mountPath':'/data'}],'env':[{'name':'GOMEMLIMIT','value':'100MiB'}]})
 def test_native_delivery_keeps_storage_guards_and_new_code(self):
  out=delivery.patch_objects(self.obj,self.plan,live=lambda *args:self.live)
  pod=out['spec']['template']['spec'];container=pod['containers'][0]
  self.assertEqual(container['image'],'new-image')
  self.assertEqual(pod['volumes'][0]['hostPath']['path'],'/verified/native')
  self.assertEqual(container['volumeMounts'][0]['mountPath'],'/data')
  self.assertEqual(pod['initContainers'][0]['name'],'guard')
  self.assertIn({'name':'GOMEMLIMIT','value':'100MiB'},container['env'])
  self.assertIn({'name':'CODE_SETTING','value':'new'},container['env'])
 def test_unknown_retired_workload_refused(self):
  self.obj['metadata']['name']='retired'
  with self.assertRaises(RuntimeError):delivery.patch_objects(self.obj,self.plan,live=lambda *args:self.live)
 def test_immutable_selector_refused(self):
  self.obj['spec']['selector']={'matchLabels':{'app':'different'}}
  with self.assertRaises(RuntimeError):delivery.patch_objects(self.obj,self.plan,live=lambda *args:self.live)
 def test_competing_operator_owner_refused(self):
  self.target['delivery']='native-offline';self.obj['kind']='Fdeployment'
  with self.assertRaises(RuntimeError):delivery.patch_objects(self.obj,self.plan)
 def test_operator_placement_uses_both_nodes(self):
  self.obj['kind']='Fdeployment'
  out=delivery.patch_objects(self.obj,self.plan)
  self.assertEqual(out['spec']['replicas'],2)
  self.assertEqual(out['spec']['placement']['spread'],'Required')
 def test_cli_parses_incoming_image_instead_of_returning_live_apply_object(self):
  self.plan.update({'destinationClusterUID':'destination','sourceClusterUID':'source'})
  def kube(args,**kwargs):
   if args[1]=='create':
    self.assertIn('--dry-run=client',args)
    return json.dumps(self.obj)
   if 'apply'in args:raise AssertionError('apply dry-run can return the existing live image')
   if 'configmap'in args:return json.dumps({'data':{'plan.json':json.dumps(self.plan)}})
   if 'namespace'in args:return json.dumps({'metadata':{'uid':'destination'}})
   if 'view'in args:return 'apps'
   if 'Deployment'in args:return json.dumps(self.live)
   raise AssertionError(args)
  with tempfile.TemporaryDirectory()as directory:
   output=Path(directory)/'output.json'
   with patch.object(sys,'argv',['prepare','input.yaml',str(output)]),patch.object(delivery.subprocess,'check_output',side_effect=kube):delivery.main()
   actual=json.loads(output.read_text())
   self.assertEqual(actual['spec']['template']['spec']['containers'][0]['image'],'new-image')
 def test_multi_resource_create_output_is_parsed_without_losing_items(self):
  config={'apiVersion':'v1','kind':'ConfigMap','metadata':{'name':'config'}}
  result=delivery.parse_manifest_json(json.dumps(self.obj)+'\n'+json.dumps(config)+'\n')
  self.assertEqual(result['kind'],'List');self.assertEqual(result['items'],[self.obj,config])
  self.assertEqual(delivery.parse_manifest_json(json.dumps(result)),result)
  with self.assertRaises(RuntimeError):delivery.parse_manifest_json('')
 def test_application_bundle_cannot_replace_native_storage_bindings(self):
  claim={'apiVersion':'v1','kind':'PersistentVolumeClaim','metadata':{'name':'state','namespace':'apps'},'spec':{'volumeName':'old-minikube-pv'}}
  volume={'apiVersion':'v1','kind':'PersistentVolume','metadata':{'name':'old-minikube-pv'},'spec':{'hostPath':{'path':'/old'}}}
  bundle={'kind':'List','items':[self.obj,claim,volume]}
  def live(kind,*args):
   return {'status':{'phase':'Bound'},'spec':{'storageClassName':'fabkli-local','volumeName':'k3s-native'}} if kind=='PersistentVolumeClaim'else self.live
  result=delivery.patch_objects(bundle,self.plan,live=live)
  self.assertEqual([o['kind']for o in result['items']],['Deployment'])
  with self.assertRaises(RuntimeError):delivery.patch_objects({'kind':'List','items':[claim]},self.plan,live=lambda *args:{'status':{'phase':'Pending'},'spec':{}})

if __name__=='__main__':unittest.main()
