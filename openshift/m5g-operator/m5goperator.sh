#!/bin/bash

#go mod init
# go get ./...

# prepare ENVs
#export KUBECONFIG=/home/agrion/kubernetes/aiyu
export OPERATOR_NAME=m5g-operator
export MYDNS="192.168.106.12"

###################################
# colorful echos
###################################

black='\E[30m'
red='\E[31m'
green='\E[32m'
yellow='\E[33m'
blue='\E[1;34m'
magenta='\E[35m'
cyan='\E[36m'
white='\E[37m'
reset_color='\E[00m'
COLORIZE=1

cecho()  {  
    # Color-echo
    # arg1 = message
    # arg2 = color
    local default_msg="No Message."
    message=${1:-$default_msg}
    color=${2:-$green}
    [ "$COLORIZE" = "1" ] && message="$color$message$reset_color"
    echo -e "$message"
    return
}

echo_error()   { cecho "$*" $red          ;}
echo_fatal()   { cecho "$*" $red; exit -1 ;}
echo_warn()    { cecho "$*" $yellow       ;}
echo_success() { cecho "$*" $green        ;}
echo_info()    { cecho "$*" $blue         ;}

run_local(){
    operator-sdk run --local --namespace=default
}

run_container(){
    case ${1} in
        start)
            kubectl apply -f deploy/service_account.yaml
            kubectl apply -f deploy/role.yaml
            kubectl apply -f deploy/role_binding.yaml
            kubectl apply -f deploy/operator.yaml
        ;;
        stop)
            kubectl delete -f deploy/service_account.yaml
            kubectl delete -f deploy/role.yaml
            kubectl delete -f deploy/role_binding.yaml
            kubectl delete -f deploy/operator.yaml
        ;;
        *)
            echo_error "Unkown option '${1}' for container"
    esac
}

apply_cr(){
    case ${1} in
        all-in-one)
            kubectl apply -f deploy/crds/allInOne/mosaic5g_v1alpha1_mosaic5g_cr.yaml
            echo "Custom Resources (CR) of monolitic Core Network is applied"
        ;;
        disaggregated-cn)
            kubectl apply -f deploy/crds/mosaic5g_v1alpha1_mosaic5g_cr.yaml
            echo "Custom Resources (CR) of disaggregated Core Network entities is applied"
        ;;
        *)
            echo_error "Unkown option '${1}' for deploy"
    esac
}

downgrade_image(){
    APISERVER=`kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}'`
    TOKEN=`kubectl get secret $(kubectl get serviceaccount default -o jsonpath='{.secrets[0].name}') -o jsonpath='{.data.token}' | base64 --decode `
    
    curl \
    -H "content-Type: application/json-patch+json" \
    -H "Authorization: Bearer ${TOKEN}"\
    --insecure \
    -X PATCH ${APISERVER}/apis/mosaic5g.com/v1alpha1/namespaces/default/mosaic5gs/mosaic5g \
    -d '[{"op":"replace","path":"/spec/cnImage","value":"arouk/oaicn:1.0"},{"op":"replace","path":"/spec/ranImage","value":"arouk/oairan:1.0"}]'
    echo " "
    echo "Core Network is downgraded to version 1.0, and RAN network is downgraded to 1.0"
}

upgrade_image(){
    APISERVER=`kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}'`
    TOKEN=`kubectl get secret $(kubectl get serviceaccount default -o jsonpath='{.secrets[0].name}') -o jsonpath='{.data.token}' | base64 --decode `

    curl \
    -H "content-Type: application/json-patch+json" \
    -H "Authorization: Bearer ${TOKEN}"\
    --insecure \
    -X PATCH ${APISERVER}/apis/mosaic5g.com/v1alpha1/namespaces/default/mosaic5gs/mosaic5g \
    -d '[{"op":"replace","path":"/spec/cnImage","value":"arouk/oaicn:1.1"},{"op":"replace","path":"/spec/ranImage","value":"arouk/oairan:1.1"}]'
    echo " "
    echo "Core Network is upgraded to version 1.1, and RAN network is upgraded to 1.1"
}

delete_cr(){
    kubectl delete -f deploy/crds/mosaic5g_v1alpha1_mosaic5g_cr.yaml
    echo "Custom Resources (CR) of the network is deleted"
}

deploy_operator_from_clean_machine(){
    echo "Start a fresh microk8s and deploy operator on it, tested with Ubuntu 18.04"
    echo "sudo without password is recommended"
    sudo snap install microk8s --classic --channel=1.14/stable
    sudo snap install kubectl --classic
    microk8s.start
    microk8s.enable dns
    # kubeconfig is used by operator
    sudo chown $USER -R $HOME/.kube
    microk8s.kubectl config view --raw > $HOME/.kube/config
    # enable privileged
    sudo bash -c 'echo "--allow-privileged=true" >> /var/snap/microk8s/current/args/kubelet'
    sudo bash -c 'echo "--allow-privileged=true" >> /var/snap/microk8s/current/args/kube-apiserver'
    # Restart kube
    sudo systemctl restart snap.microk8s.daemon-kubelet.service
    sudo systemctl restart snap.microk8s.daemon-apiserver.service
    # Configure DNS if it's not working
    # microk8s.kubectl -n kube-system edit configmap/coredns


}

init(){
    echo "Applying crd..."
    kubectl apply -f deploy/crds/mosaic5g_v1alpha1_mosaic5g_crd.yaml
    kubectl apply -f defaultRole.yaml
    sleep 3
    echo "Done, now run [local] or [container start] to create your operator"
}

clean(){
    kubectl delete -f deploy/crds/mosaic5g_v1alpha1_mosaic5g_crd.yaml
    kubectl delete -f defaultRole.yaml
}

break_down(){
    sudo snap remove microk8s 
    sudo snap remove kubectl 
}

main() {
    case ${1} in
        -n | --init)
            init
        ;;
        -c | --clean)
            clean
        ;;
        -l | --local)
            run_local
        ;;
        container)
            run_container ${2}
        ;;
        deploy)
            apply_cr ${2}
        ;;
        upgrade)
            upgrade_image
        ;;
        downgrade)
            downgrade_image
        ;;
        -d | --delete)
            delete_cr 
        ;;
        -i | --install)
            deploy_operator_from_clean_machine 
        ;;
        -r | --remove)
            break_down
        ;;
        *)
            echo_info '
This program installs the requirements to run kubernets on one machine, 
and run mosaic5g-operator as a pod inside kubernetes in order to manage
the deployments and services of 4G/5G networks in cloud native environment.
This program also allows to run mosaic5g-operator locally as Golang app.
Options:
-i | --install
    Install and run microk8s kubectl, then deploy operator on it
-n | --init
    Apply CRD to k8s cluster (Required for Operator)
-l | --local
    Run Operator as a Golang app at local
container [start|stop]
    Run Operator as a POD inside Kubernetes
deploy [all-in-one|disaggregated-cn]
    Deploy the network with:
        - all-in-one: all the core network entities (oai-hss, oai-mme, oai-spgw) in one pod
        - disaggregated-cn: the core network entities (oai-hss, oai-mme, oai-spgw) are deployed on disaggregated pods
-d | --delete 
    Stop the network by deleting the Custom Resources (CR) of the network
upgrade 
    upgrade the images of the network to the new version v1.1
downgrade 
    downgrade the images of the network to the old version v1.0
-c | --clean 
    Remove CRD from cluster
-r | --remove
    remove the snap of kubectl and microk8s
Usage:
    ./m5goperator.sh -i 
    ./m5goperator.sh container start
    ./m5goperator.sh deploy all-in-one
    ./m5goperator.sh deploy disaggregated-cn
            '
        ;;
    esac

}
main "$@"
