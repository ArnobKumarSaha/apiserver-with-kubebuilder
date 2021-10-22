kubebuilder init --domain DOMAIN_NAME --repo MODULE_NAME
kubebuilder create api --group GROUP_NAME --version VERSION_NAME --kind KIND_NAME


## Concrete example (in my case)::
1) kubebuilder init --domain saha.com --repo saha.com/mycrd
2) kubebuilder create api --group webapp --version v1 --kind Neymar


-----
| Module_name will be used in go.mod file.
| Group_Name, Domain_Name & Version_Name will be used in api/$VERSION/groupverson_info.yaml file.
| Kind_Name will be used in api/$VERSION/$KIND_types.go file.
-----


## Use Makefile to make life easier :: 
- make (To generate deep copies , formatting the code, examine the source code, build the bin/manager binary)
- make install (To generate clusterrole, webhook & CRDs, apply the crds into cluster)
- make run (To run the controller)
- kubectl apply -f config/samples/webapp_v1_neymar.yaml